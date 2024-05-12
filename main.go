package main

import (
	"be-assignment/helper"
	"be-assignment/repository"
	"be-assignment/router"
	"context"
	"github.com/JPratama7/util/token/paseto"
	"github.com/gin-gonic/gin"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

func main() {

	mainCtx := context.Background()

	client, err := mongo.Connect(mainCtx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		panic(err)
	}

	gi := gin.Default()

	engine := fizz.NewFromEngine(gi)

	// Use global error handler middleware
	engine.Use(helper.GlobalErrorHandler())

	// token generator
	token := paseto.NewPASETO(os.Getenv("PASETO_PUBLIC"), os.Getenv("PASETO_SECRET"))
	engine.Use(helper.SetContextData("token", &token))

	// init repository
	userRepo := repository.NewUser(client.Database("be-assignment"), "users")
	transactionRepo := repository.NewTransaction(client.Database("be-assignment"), "transactions")
	paymentRepo := repository.NewPaymentAccount(client.Database("be-assignment"), "payment_accounts")

	// put repo to context
	engine.Use(helper.SetContextData("userRepo", userRepo))
	engine.Use(helper.SetContextData("transactionRepo", transactionRepo))
	engine.Use(helper.SetContextData("paymentRepo", paymentRepo))

	router.UserRoute(engine)
	router.Transaction(engine)

	infos := &openapi.Info{
		Title:       "Assignment",
		Description: `An Assignment REST API.`,
		Version:     "1.0.0",
	}

	engine.GET("/openapi.json", nil, engine.OpenAPI(infos, "json"))
	gi.Run()
}
