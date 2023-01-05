package main

import (
	"flashy-product/common"
	"flashy-product/rabbitmq"
	"flashy-product/repositories"
	"flashy-product/services"
	"fmt"
)

func main() {
	db, err := common.NewMysqlConn()
	if err != nil {
		fmt.Println(err)
	}
	// Create a product database operation instance
	product := repositories.NewProductManager("product", db)
	//Create product serivce
	productService := services.NewProductService(product)
	// Create an Order database instance
	order := repositories.NewOrderMangerRepository("order", db)
	//Create order Service
	orderService := services.NewOrderService(order)

	rabbitmqConsumeSimple := rabbitmq.NewRabbitMQSimple("flashyProduct")
	rabbitmqConsumeSimple.ConsumeSimple(orderService, productService)
}
