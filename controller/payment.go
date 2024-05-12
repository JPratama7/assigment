package controller

import (
	"be-assignment/helper"
	"be-assignment/model"
	"be-assignment/repository"
	"github.com/JPratama7/util/token/paseto"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddPayment(ctx *gin.Context, req *model.PaymentAccountData) (res *model.ReturnData[primitive.ObjectID], err error) {
	tokenHeader, err := helper.GetAuthHeader(ctx)
	if err != nil {
		err = model.NewError(401, "unauthorized "+err.Error())
		return
	}

	tokenGen, err := helper.GetContextData[paseto.PASETO](ctx, "token")
	if err != nil {
		err = model.NewError(500, "token generator not found")
		return
	}

	tokenized := new(paseto.Payload[string])
	err = tokenGen.Decode(tokenHeader, tokenized)
	if err != nil {
		err = model.NewError(401, "unauthorized "+err.Error())
		return
	}

	repo, err := helper.GetContextData[repository.PaymentAccount](ctx, "paymentRepo")
	if err != nil {
		err = model.NewError(500, "repository not found")
		return
	}

	repoUser, err := helper.GetContextData[repository.User](ctx, "userRepo")
	if err != nil {
		err = model.NewError(500, "repository not found")
		return
	}

	user, err := repoUser.FindByUser(ctx, tokenized.Id)
	if err != nil {
		err = model.NewError(404, "user not found")
		return
	}

	id, err := repo.Create(ctx, model.PaymentAccount{
		AccountID: user.Id,
		Balance:   req.Balance,
	})

	if err != nil {
		err = model.NewError(500, "failed to create payment account")
		return
	}

	res = model.NewReturnData(200, true, "Payment Account Created", &id)
	return
}
