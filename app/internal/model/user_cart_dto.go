package model

// CartResponse defines the overall cart response structure
type CartResponse []CartStoreDTO

// CartStoreDTO defines the store information and its products
type CartStoreDTO struct {
	StoreID   uint64           `json:"store_id"`
	StoreName string           `json:"store_name"`
	StoreURL  string           `json:"store_url"`
	List      []CartProductDTO `json:"list"`
}

// CartProductDTO defines the product information in the cart
type CartProductDTO struct {
	CartID         uint    `json:"cart_id"`
	ProductID      uint64  `json:"product_id"`
	ProductName    string  `json:"product_name"`
	CurrentPrice   float64 `json:"current_price"`
	Quantity       int     `json:"quantity"`
	CourierFeeMin  float64 `json:"courier_fee_min"`
	MemberDiscount float64 `json:"member_discount"`
	CouponID       uint64  `json:"coupon_id"`
	CouponPrice    float64 `json:"coupon_price"`
}
