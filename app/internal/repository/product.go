package repository

import (
	"context"
	"strconv"
	"strings"
	"time"

	"app/internal/model"
)

type ProductRepository interface {
	GetProductList(ctx context.Context, userId string) ([]*model.ProductListResponse, error)
	GetRecommendProductList(ctx context.Context) ([]*model.ProductListItemDTO, error)
	GetProductByID(ctx context.Context, id int, userId string) (*model.ProductListItemDTO, error)
	GetProductDetailsByCartIDs(ctx context.Context, cartIDs []uint64) ([]model.ProductListItemDTO, error)
}

func NewProductRepository(
	repository *Repository,
) ProductRepository {
	return &productRepository{
		Repository: repository,
	}
}

type productRepository struct {
	*Repository
}

func (r *productRepository) GetProductList(ctx context.Context, userId string) ([]*model.ProductListResponse, error) {
	// Get all categories (level 1)
	var category1List []*model.Category
	if err := r.DB(ctx).Where("parent_id = 0").Find(&category1List).Error; err != nil {
		return nil, err
	}

	// 根据userId获取购物车
	var cart []*model.UserCart
	if err := r.DB(ctx).Where("user_id = ? AND status = 0", userId).Find(&cart).Error; err != nil {
		return nil, err
	}

	// 获取购物车中的商品ID
	productCardMap := make(map[uint64]int)
	for _, c := range cart {
		productCardMap[c.ProductID] = c.Quantity
	}

	var result []*model.ProductListResponse
	// For each category1, get its category2 and products
	for _, category1 := range category1List {
		// Get category2 list
		var category2List []*model.Category
		if err := r.DB(ctx).Where("parent_id = ?", category1.ID).Find(&category2List).Error; err != nil {
			return nil, err
		}

		category1Response := &model.ProductListResponse{
			Category1ID:   int(category1.ID),
			Category1Name: category1.CategoryName,
			Category2List: []model.ProductCategory2DTO{},
		}

		// For each category2, get its products
		for _, category2 := range category2List {
			// Get products in this category
			var products []*model.Product
			if err := r.DB(ctx).Where("category2_id = ?", category2.ID).Find(&products).Error; err != nil {
				return nil, err
			}

			category2Response := model.ProductCategory2DTO{
				Category2ID:   int(category2.ID),
				Category2Name: category2.CategoryName,
				ProductList:   []model.ProductListItemDTO{},
			}

			// For each product, get coupons and evaluations
			for _, product := range products {
				// Process product images (split banner_img by comma)
				var productImages []string
				if product.BannerImg != "" {
					productImages = strings.Split(product.BannerImg, ",")
				}

				productQuantity := 0
				if quantity, ok := productCardMap[uint64(product.ID)]; ok {
					productQuantity = quantity
				}
				// Create product DTO
				productDTO := model.ProductListItemDTO{
					ProductID:            product.ID,
					ProductName:          product.ProductName,
					ProductCurrentPrice:  product.CurrentPrice,
					ProductOriginalPrice: product.OriginPrice,
					ProductUnit:          "", // You may need to add this field to the Product model
					ProductSpec:          product.Specification,
					ProductSales:         product.Sales,
					ProductSpecification: "¥" + strconv.FormatFloat(product.CourierFeeMin, 'f', -1, 64) + "-" + strconv.FormatFloat(product.CourierFeeMax, 'f', -1, 64),
					HeaderImg:            product.HeaderImg,
					ProductImages:        productImages,
					ProductIsSpecial:     product.IsSpecial,
					CourierFeeMin:        product.CourierFeeMin,
					CourierFeeMax:        product.CourierFeeMax,
					ProductContent:       product.Content,
					ProductCoupons:       []model.ProductCouponDTO{},
					ProductEvaluate:      []model.ProductEvaluateDTO{},
					ProductQuantity:      productQuantity,
				}

				category2Response.ProductList = append(category2Response.ProductList, productDTO)
			}

			category1Response.Category2List = append(category1Response.Category2List, category2Response)
		}

		result = append(result, category1Response)
	}

	return result, nil
}

// Helper function to safely get string value from pointer
func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func (r *productRepository) GetRecommendProductList(ctx context.Context) ([]*model.ProductListItemDTO, error) {
	var products []*model.Product
	if err := r.DB(ctx).Where("is_recommend = 1").Order("recommend_sort ASC").Find(&products).Error; err != nil {
		return nil, err
	}

	result := []*model.ProductListItemDTO{}
	for _, product := range products {
		result = append(result, &model.ProductListItemDTO{
			ProductID:            product.ID,
			ProductName:          product.ProductName,
			ProductCurrentPrice:  product.CurrentPrice,
			ProductOriginalPrice: product.OriginPrice,
			ProductUnit:          "", // You may need to add this field to the Product model
			ProductSpec:          product.Specification,
		})
	}

	return result, nil
}

// 通过商品ID获取商品
func (r *productRepository) GetProductByID(ctx context.Context, id int, userId string) (*model.ProductListItemDTO, error) {
	var product model.Product
	if err := r.DB(ctx).Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}
	// Get coupons
	var coupons []*model.ProductCoupon
	if err := r.DB(ctx).Where("product_id = ?", product.ID).Find(&coupons).Error; err != nil {
		return nil, err
	}
	// Get evaluations
	var evaluations []*model.ProductReview
	if err := r.DB(ctx).Where("product_id = ?", product.ID).Find(&evaluations).Error; err != nil {
		return nil, err
	}
	// 用户是否已经领取优惠券
	userCouponList := []*model.UserCoupon{}
	if err := r.DB(ctx).Where("user_id = ? AND product_id = ?", userId, id).Find(&userCouponList).Error; err != nil {
		return nil, err
	}
	userCouponMap := make(map[uint64]uint8)
	for _, userCoupon := range userCouponList {
		userCouponMap[userCoupon.CouponID] = userCoupon.Status + 1
	}

	// Add coupons
	productCouponsDTO := []model.ProductCouponDTO{}
	for _, coupon := range coupons {
		deadlineStr := ""
		if coupon.Deadline != nil {
			deadlineStr = coupon.Deadline.Format("2006-01-02")
		}
		isReceived := 0
		if status, ok := userCouponMap[uint64(coupon.ID)]; ok {
			isReceived = int(status)
		}
		couponDTO := model.ProductCouponDTO{
			CouponID:          coupon.ID,
			CouponName:        coupon.CouponName,
			CouponPrice:       coupon.CouponPrice,
			AvailableMinPrice: coupon.AvailableMinPrice,
			Deadline:          deadlineStr,
			IsReceived:        isReceived,
		}
		productCouponsDTO = append(productCouponsDTO, couponDTO)
	}

	// Add evaluations
	productEvaluateDTO := []model.ProductEvaluateDTO{}
	for _, eval := range evaluations {
		// Process evaluation images
		var evalImages []string
		if len(eval.Images) == 0 && eval.Images != "" {
			evalImages = strings.Split(eval.Images, ",")
		}

		evalDTO := model.ProductEvaluateDTO{
			EvaluateID:      uint(eval.ID),
			EvaluateContent: getStringValue(&eval.Content),
			EvaluateTime:    time.Now().Format("2006-01-02"),                                    // You might want to use a field from eval
			Nickname:        eval.UserID,                                                        // Consider adding a nickname field
			Avatar:          "https://txtimages.oss-cn-beijing.aliyuncs.com/product/avatar.jpg", // Default or from DB
			IsReturn:        eval.IsReturn,
			EvaluateImages:  evalImages,
			ViewNums:        eval.ViewNums,
			EvaluateNums:    eval.EvaluateNums,
			PraiseNums:      eval.PraiseNums,
		}
		productEvaluateDTO = append(productEvaluateDTO, evalDTO)
	}
	return &model.ProductListItemDTO{
		ProductID:            product.ID,
		ProductName:          product.ProductName,
		ProductCurrentPrice:  product.CurrentPrice,
		ProductOriginalPrice: product.OriginPrice,
		ProductUnit:          "", // You may need to add this field to the Product model
		ProductSpec:          product.Specification,
		ProductSales:         product.Sales,
		ProductSpecification: "¥" + strconv.FormatFloat(product.CourierFeeMin, 'f', -1, 64) + "-" + strconv.FormatFloat(product.CourierFeeMax, 'f', -1, 64),
		HeaderImg:            product.HeaderImg,
		ProductImages:        strings.Split(product.BannerImg, ","),
		ProductIsSpecial:     product.IsSpecial,
		CourierFeeMin:        product.CourierFeeMin,
		CourierFeeMax:        product.CourierFeeMax,
		MemberDiscount:       product.MemberDiscount,
		ProductContent:       product.Content,
		ProductCoupons:       productCouponsDTO,
		ProductEvaluate:      productEvaluateDTO,
		ProductEvaluateNums:  len(productEvaluateDTO),
		UserID:               userId,
		Category1ID:          *product.Category1ID,
		Category2ID:          *product.Category2ID,
	}, nil
}

// GetProductDetailsByCartIDs 通过购物车ID查询商品详情
func (r *productRepository) GetProductDetailsByCartIDs(ctx context.Context, cartIDs []uint64) ([]model.ProductListItemDTO, error) {
	// 首先获取购物车信息
	var carts []model.UserCart
	if err := r.DB(ctx).Where("id IN (?)", cartIDs).Find(&carts).Error; err != nil {
		return nil, err
	}

	// 提取购物车中的商品ID
	productIDs := make([]uint64, 0, len(carts))
	userCouponMap := make(map[uint64]model.UserCart)
	for _, cart := range carts {
		productIDs = append(productIDs, cart.ProductID)
		userCouponMap[cart.ProductID] = cart
	}

	// 没有商品ID，直接返回空列表
	if len(productIDs) == 0 {
		return []model.ProductListItemDTO{}, nil
	}

	// 查询商品详情
	var products []model.Product
	if err := r.DB(ctx).Model(&model.Product{}).Where("id IN (?)", productIDs).Find(&products).Error; err != nil {
		return nil, err
	}
	var productsDTO []model.ProductListItemDTO
	for _, product := range products {
		productsDTO = append(productsDTO, model.ProductListItemDTO{
			ProductID:            product.ID,
			ProductName:          product.ProductName,
			ProductCurrentPrice:  product.CurrentPrice,
			ProductOriginalPrice: product.OriginPrice,
			ProductUnit:          "", // You may need to add this field to the Product model
			ProductSpec:          product.Specification,
			ProductSales:         product.Sales,
			ProductSpecification: "¥" + strconv.FormatFloat(product.CourierFeeMin, 'f', -1, 64) + "-" + strconv.FormatFloat(product.CourierFeeMax, 'f', -1, 64),
			HeaderImg:            product.HeaderImg,
			ProductImages:        strings.Split(product.BannerImg, ","),
			ProductIsSpecial:     product.IsSpecial,
			CourierFeeMin:        product.CourierFeeMin,
			CourierFeeMax:        product.CourierFeeMax,
			MemberDiscount:       product.MemberDiscount,
			ProductContent:       product.Content,
		})
	}
	for k, product := range products {
		// 是否存在优惠券
		if cart, ok := userCouponMap[uint64(product.ID)]; ok {
			productsDTO[k].CouponID = int(cart.CouponID)
			productsDTO[k].CouponPrice = cart.CouponPrice
		}
	}

	return productsDTO, nil
}
