package datamodels

type Product struct {
	ID           int64  `json:"id" sql:"ID" flashy:"ID"`
	ProductName  string `json:"ProductName" sql:"productName" flashy:"ProductName"`
	ProductNum   int64  `json:"ProductNum" sql:"productNum" flashy:"ProductNum"`
	ProductImage string `json:"ProductImage" sql:"productImage" flashy:"ProductImage"`
	ProductUrl   string `json:"ProductUrl" sql:"productUrl" flashy:"ProductUrl"`
}
