package domain

type Position struct {
	Product Product
	Count   uint32
}

type Cart struct {
	Positions  map[ProductSKU]Position
	TotalPrice uint32
}

type PostProductRequest struct {
	UserID     int64
	ProductSKU int64
	Count      uint32
}

type GetProductRequest struct {
	UserID int64
}

type Item struct {
	Sku   int64
	Name  string
	Count uint32
	Price uint32
}

type GetProductResponse struct {
	Items      []Item
	TotalPrice uint32
}
