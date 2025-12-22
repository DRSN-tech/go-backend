package converter

type ProductInfoRedisModel struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	CategoryName string `json:"category_name"`
	Price        int64  `json:"price"`
}
