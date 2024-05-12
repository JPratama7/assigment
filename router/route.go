package router

import (
	"be-assignment/controller"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"
)

func UserRoute(e *fizz.Fizz) {
	g := e.Group("/user", "User", "User Login")
	g.POST("/login", []fizz.OperationOption{
		fizz.Description("Login to the system"),
		fizz.Summary("Login"),
	}, tonic.Handler(controller.Login, 200))
	g.POST("/register", []fizz.OperationOption{
		fizz.Description("Register to the system"),
		fizz.Summary("Register"),
	}, tonic.Handler(controller.Register, 200))

	g.POST("/payment", []fizz.OperationOption{
		fizz.Description("Add Payment"),
		fizz.Summary("Payment"),
	}, tonic.Handler(controller.AddPayment, 200))

	g.GET("/payment", []fizz.OperationOption{
		fizz.Description("Get All Payment and Transaction"),
		fizz.Summary("Payment"),
	}, tonic.Handler(controller.AllTransaction, 200))
}

func Transaction(e *fizz.Fizz) {
	g := e.Group("/transaction", "User", "User Login")
	g.POST("/withdraw", []fizz.OperationOption{
		fizz.Description("Withdraw from account"),
		fizz.Summary("Withdraw"),
	}, tonic.Handler(controller.Withdraw, 200))
	g.POST("/send", []fizz.OperationOption{
		fizz.Description("Send to account"),
		fizz.Summary("Send"),
	}, tonic.Handler(controller.Send, 200))
}
