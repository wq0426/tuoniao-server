package model

// ProductListResponse represents the top level response structure
type ProductListResponse struct {
	Category1ID   int                   `json:"category1_id"`
	Category1Name string                `json:"category1_name"`
	Category2List []ProductCategory2DTO `json:"category2_list"`
}

// ProductCategory2DTO represents the second level category with products
type ProductCategory2DTO struct {
	Category2ID   int                  `json:"category2_id"`
	Category2Name string               `json:"category2_name"`
	ProductList   []ProductListItemDTO `json:"product_list"`
}

// ProductListItemDTO represents a product in the list
type ProductListItemDTO struct {
	ItemID               uint                 `json:"item_id"`
	ProductID            uint                 `json:"product_id"`
	ProductName          string               `json:"product_name"`
	ProductCurrentPrice  float64              `json:"product_current_price"`
	ProductOriginalPrice float64              `json:"product_original_price"`
	ProductUnit          string               `json:"product_unit"`
	ProductSpec          string               `json:"product_spec"`
	ProductSales         int                  `json:"product_sales"`
	ProductSpecification string               `json:"product_specification"`
	HeaderImg            string               `json:"header_img"`
	ProductImages        []string             `json:"product_images"`
	ProductIsSpecial     *int                 `json:"product_is_special"`
	CourierFeeMin        float64              `json:"courier_fee_min"`
	CourierFeeMax        float64              `json:"courier_fee_max"`
	MemberDiscount       float64              `json:"member_discount"`
	ProductContent       string               `json:"product_content"`
	ProductCoupons       []ProductCouponDTO   `json:"product_coupons"`
	ProductEvaluate      []ProductEvaluateDTO `json:"product_evaluate"`
	ProductEvaluateNums  int                  `json:"product_evaluate_nums"`
	ProductQuantity      int                  `json:"product_quantity"`
	UserID               string               `json:"user_id"`
	Category1ID          int                  `json:"category1_id"`
	Category2ID          int                  `json:"category2_id"`
	CouponID             int                  `json:"coupon_id"`
	CouponPrice          float64              `json:"coupon_price"`
}

// ProductCouponDTO represents a coupon for a product
type ProductCouponDTO struct {
	CouponID          uint    `json:"coupon_id"`
	CouponName        string  `json:"coupon_name"`
	CouponPrice       float64 `json:"coupon_price"`
	AvailableMinPrice float64 `json:"available_min_price"`
	Deadline          string  `json:"deadline"`
	IsReceived        int     `json:"is_received"`
	ProductID         uint    `json:"product_id"`
}

// ProductEvaluateDTO represents an evaluation for a product
type ProductEvaluateDTO struct {
	EvaluateID      uint     `json:"evaluate_id"`
	EvaluateContent string   `json:"evaluate_content"`
	EvaluateTime    string   `json:"evaluate_time"`
	Nickname        string   `json:"nickname"`
	Avatar          string   `json:"avatar"`
	IsReturn        int8     `json:"is_return"`
	EvaluateImages  []string `json:"evaluate_images"`
	ViewNums        uint     `json:"view_nums"`
	EvaluateNums    uint     `json:"evaluate_nums"`
	PraiseNums      uint     `json:"praise_nums"`
}
