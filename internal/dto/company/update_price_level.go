package company

type UpdatePriceLevelReq struct {
	PriceLevel string `json:"price_level" form:"price_level" binding:"required,oneof=price_1 price_2 price_3 price_4" example:"price_1或price_2或price_3或price_4"`
}
