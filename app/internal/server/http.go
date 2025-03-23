package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"app/docs"
	"app/internal/handler"
	"app/internal/middleware"
	"app/pkg/jwt"
	"app/pkg/log"
	"app/pkg/server/http"
)

func NewHTTPServer(
	logger *log.Logger,
	conf *viper.Viper,
	jwt *jwt.JWT,
	userHandler *handler.AccountHandler,
	resourceHandler *handler.ResourceHandler,
	settingsHandler *handler.SettingsHandler,
	smsHandler *handler.SmsHandler,
	bannerHandler *handler.BannerHandler,
	monitorHandler *handler.MonitorHandler,
	newsHandler *handler.NewsHandler,
	productHandler *handler.ProductHandler,
	freeMarketMineHandler *handler.FreeMarketMineHandler,
	userCartHandler *handler.UserCartHandler,
	userOrderHandler *handler.UserOrderHandler,
	userAddressHandler *handler.UserAddressHandler,
	userAssetHandler *handler.UserAssetHandler,
	userCouponHandler *handler.UserCouponHandler,
	pointExchangeConfigHandler *handler.PointExchangeConfigHandler,
	refundOrderHandler *handler.RefundOrderHandler,
	withdrawOrderHandler *handler.WithdrawOrderHandler,
	productReviewHandler *handler.ProductReviewHandler,
	productEvaluateHandler *handler.ProductEvaluateHandler,
	userEarningHandler *handler.UserEarningHandler,
) *http.Server {
	gin.SetMode(gin.DebugMode)
	s := http.NewServer(
		gin.Default(),
		logger,
		http.WithServerHost(conf.GetString("http.host")),
		http.WithServerPort(conf.GetInt("http.port")),
	)
	s.GET("/", resourceHandler.GetResource)
	s.GET("/resource/*path", resourceHandler.GetResource)
	// swagger doc
	docs.SwaggerInfo.BasePath = "/api/v1"
	s.GET(
		"/swagger/*any", ginSwagger.WrapHandler(
			swaggerfiles.Handler,
			// ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", conf.GetInt("app.http.port"))),
			ginSwagger.DefaultModelsExpandDepth(-1),
			ginSwagger.PersistAuthorization(true),
		),
	)

	s.Use(
		middleware.CORSMiddleware(),
		middleware.ResponseLogMiddleware(logger),
		middleware.RequestLogMiddleware(logger),
	)
	v1 := s.Group("/api/v1")
	{
		// No route group has permission
		userAuthRouter := v1.Group("/account")
		{
			userAuthRouter.POST("/code", smsHandler.SendCode)
			userAuthRouter.POST("/login", userHandler.Login)
			userAuthRouter.GET("/profile", middleware.SignMiddleware(logger, conf), userHandler.GetProfile)
			userAuthRouter.POST("/reset", middleware.SignMiddleware(logger, conf), userHandler.ResetPassword)
			userAuthRouter.PUT("/profile", middleware.SignMiddleware(logger, conf), userHandler.UpdateProfileSettings)
			userAuthRouter.GET("/wechat/login/callback", userHandler.WeChatLoginCallback)
			userAuthRouter.POST("/wechat_mini_login", userHandler.WeChatMiniLogin)
			userAuthRouter.POST("/wechat/decrypt", userHandler.WeChatDecrypt)
		}
		fileRouter := v1.Group("/file").Use(middleware.SignMiddleware(logger, conf))
		{
			fileRouter.POST("/avatar", settingsHandler.UploadAvatar)
		}
		// Banner - this can be accessed without authentication
		bannerRouter := v1.Group("/banner")
		{
			bannerRouter.GET("/list", bannerHandler.GetBannerList)
		}
		// Monitor - this can be accessed without authentication
		monitorRouter := v1.Group("/monitor")
		{
			monitorRouter.GET("/list", monitorHandler.GetMonitorList)
		}
		// News - this can be accessed without authentication
		newsRouter := v1.Group("/news")
		{
			newsRouter.GET("/list", newsHandler.GetNewsList)
		}
		// 商品
		productRouter := v1.Group("/product").Use(middleware.SignMiddleware(logger, conf))
		{
			productRouter.GET("/list", productHandler.GetProductList)
			productRouter.GET("/recommend", productHandler.GetRecommendProductList)
			productRouter.GET("/detail", productHandler.GetProductByID)
			productRouter.GET("/details", productHandler.GetProductDetailsByCartIDs)
		}
		// 自由市场
		freeMarketRouter := v1.Group("/market").Use(middleware.SignMiddleware(logger, conf))
		{
			freeMarketRouter.GET("/mine", freeMarketMineHandler.GetUserEggsSummary)
			freeMarketRouter.POST("/update-price", freeMarketMineHandler.UpdateEggPrice)
		}
		// 购物车
		cartRouter := v1.Group("/cart").Use(middleware.SignMiddleware(logger, conf))
		{
			cartRouter.POST("/add", userCartHandler.AddToCart)
			cartRouter.GET("/list", userCartHandler.GetUserCartItems)
			cartRouter.POST("/delete", userCartHandler.DeleteCartItems)
		}
		// 订单
		orderRouter := v1.Group("/order").Use(middleware.SignMiddleware(logger, conf))
		{
			orderRouter.POST("/create", userOrderHandler.CreateOrders)
			orderRouter.GET("/list", userOrderHandler.GetOrderList)
			orderRouter.GET("/products", userOrderHandler.GetOrderProductDetails)
			orderRouter.POST("/status", userOrderHandler.UpdateOrderStatus)
			orderRouter.GET("/detail/:order_item_id", userOrderHandler.GetOrderDetail)
		}
		// 用户地址
		addressRouter := v1.Group("/address").Use(middleware.SignMiddleware(logger, conf))
		{
			addressRouter.POST("/add", userAddressHandler.AddAddress)
			addressRouter.POST("/update", userAddressHandler.UpdateAddress)
			addressRouter.GET("/list", userAddressHandler.GetUserAddresses)
			addressRouter.POST("/delete", userAddressHandler.DeleteAddress)
		}
		// 用户资产
		assetRouter := v1.Group("/asset").Use(middleware.SignMiddleware(logger, conf))
		{
			assetRouter.GET("/info", userAssetHandler.GetUserAsset)
			assetRouter.POST("/recharge", userAssetHandler.RechargeBalance)
			assetRouter.GET("/balance/records", userAssetHandler.GetBalanceRecords)
			assetRouter.GET("/withdraw/records", userAssetHandler.GetWithdrawRecords)
			assetRouter.GET("/exchange/records", userAssetHandler.GetExchangeRecords)
		}
		// 优惠券
		couponRouter := v1.Group("/coupon").Use(middleware.SignMiddleware(logger, conf))
		{
			couponRouter.GET("/product", userCouponHandler.GetUserCoupons)
			couponRouter.GET("/list", userCouponHandler.GetAllUserCoupons)
			couponRouter.POST("/claim", userCouponHandler.ClaimCoupon)
			couponRouter.GET("/detail", userCouponHandler.GetUserCouponDetail)
		}
		// 积分兑换配置
		pointExchangeRouter := v1.Group("/point/exchange").Use(middleware.SignMiddleware(logger, conf))
		{
			pointExchangeRouter.GET("/list", pointExchangeConfigHandler.GetPointExchangeConfigList)
			pointExchangeRouter.POST("/exchange", pointExchangeConfigHandler.ExchangePoints)
		}
		// 退款相关路由
		refundRouter := v1.Group("/refund").Use(middleware.SignMiddleware(logger, conf))
		{
			refundRouter.POST("/create", refundOrderHandler.CreateRefund)
			refundRouter.GET("/list", refundOrderHandler.GetRefundList)
			refundRouter.GET("/detail/:refund_id", refundOrderHandler.GetRefundDetail)
			refundRouter.POST("/cancel", refundOrderHandler.CancelRefund)
			refundRouter.DELETE("/delete/:refund_id", refundOrderHandler.DeleteRefund)
		}
		// 提现相关路由
		withdrawRouter := v1.Group("/withdraw").Use(middleware.SignMiddleware(logger, conf))
		{
			withdrawRouter.POST("/create", withdrawOrderHandler.CreateWithdraw)
			withdrawRouter.GET("/list", withdrawOrderHandler.GetWithdrawList)
			withdrawRouter.GET("/detail/:withdraw_id", withdrawOrderHandler.GetWithdrawDetail)
		}
		// 评价相关路由
		reviewRouter := v1.Group("/review").Use(middleware.SignMiddleware(logger, conf))
		{
			reviewRouter.POST("/create", productReviewHandler.CreateReview)
			reviewRouter.GET("/product/:product_id", productReviewHandler.GetProductReviews)
			reviewRouter.GET("/user", productReviewHandler.GetUserReviews)
			reviewRouter.GET("/user/tab", productReviewHandler.GetUserReviewsByTab)
			reviewRouter.DELETE("/delete/:review_id", productReviewHandler.DeleteReview)
			reviewRouter.GET("/detail/:review_id", productReviewHandler.GetReviewDetail)
			reviewRouter.POST("/increment_counter", productReviewHandler.IncrementReviewCounter)
			reviewRouter.POST("/reply", productEvaluateHandler.CreateEvaluateReply)
			reviewRouter.POST("/comment", productEvaluateHandler.CreateEvaluate)
			reviewRouter.POST("/update_anonymous", productEvaluateHandler.UpdateEvaluateAnonymous)
		}
		// 用户收益相关路由
		earningRouter := v1.Group("/earning").Use(middleware.SignMiddleware(logger, conf))
		{
			earningRouter.POST("/add", userEarningHandler.AddEarning)
			earningRouter.GET("/list", userEarningHandler.GetEarningList)
		}
	}

	return s
}
