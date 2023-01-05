package main

import (
	"context"
	"flashy-product/backend/web/controllers"
	"flashy-product/common"
	"flashy-product/repositories"
	"flashy-product/services"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/opentracing/opentracing-go/log"
)

func main() {
	//1. initialize Iris
	app := iris.New()
	//2. set error mode, show Error message in MVC mode
	app.Logger().SetLevel("debug")
	//3. register template
	tmplate := iris.HTML("./backend/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(tmplate)
	//4. set template goal
	app.StaticWeb("/assets", "./backend/web/assets")
	// jump to certain page then there's an error
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "Page error！"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})
	// connect to mysql database
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//5. register controler
	productRepository := repositories.NewProductManager("product", db)
	productService := services.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productService)
	product.Handle(new(controllers.ProductController))

	orderRepository := repositories.NewOrderMangerRepository("order", db)
	orderService := services.NewOrderService(orderRepository)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx, orderService)
	order.Handle(new(controllers.OrderController))

	//6. activate the service
	app.Run(
		iris.Addr("localhost:8080"),
		iris.WithoutVersionChecker,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)

}
