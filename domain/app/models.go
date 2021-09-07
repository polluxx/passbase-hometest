package app

type Project struct {
	Name string `json:"name" binding:"required"`
	CustomerEmail string `json:"customer_email" binding:"required"`
}