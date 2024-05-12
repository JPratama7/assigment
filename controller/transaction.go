package controller

import (
	"be-assignment/helper"
	"be-assignment/model"
	"be-assignment/repository"
	"github.com/JPratama7/util/token/paseto"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Withdraw(ctx *gin.Context, req *model.WithdrawData) (res *model.ReturnData[primitive.ObjectID], err error) {
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

	repoUser, err := helper.GetContextData[repository.User](ctx, "userRepo")
	if err != nil {
		err = model.NewError(500, "repository not found")
		return
	}

	repo, err := helper.GetContextData[repository.Transaction](ctx, "transactionRepo")
	if err != nil {
		err = model.NewError(500, "repository not found")
		return
	}

	repoPayment, err := helper.GetContextData[repository.PaymentAccount](ctx, "paymentRepo")
	if err != nil {
		err = model.NewError(500, "repository not found")
		return
	}

	tokenized := new(paseto.Payload[string])
	err = tokenGen.Decode(tokenHeader, tokenized)
	if err != nil {
		err = model.NewError(401, "unauthorized "+err.Error())
		return
	}

	user, err := repoUser.FindByUser(ctx, tokenized.Id)
	if err != nil {
		err = model.NewError(404, "user not found")
		return
	}

	paymentData, err := repoPayment.FindByIdUser(ctx, req.Account, user.Id.Hex())
	if err != nil {
		err = model.NewError(404, "user payment not found")
		return
	}

	if paymentData.Balance < req.Amount {
		err = model.NewError(400, "insufficient balance")
		return
	}

	paymentData.Balance -= req.Amount
	err = repoPayment.Update(ctx, paymentData)
	if err != nil {
		err = model.NewError(500, "failed to update balance")
		return
	}

	id, err := repo.Create(ctx, model.Transaction{
		FromAccountID: paymentData.Id,
		ToAccountID:   paymentData.Id,
		Amount:        req.Amount,
		Type:          model.Withdraw,
		Status:        model.Success,
		CreatedAt:     helper.NewTimestamp(),
		UpdatedAt:     primitive.Timestamp{},
	})
	if err != nil {
		err = model.NewError(500, "failed to create transaction")
		return
	}

	res = model.NewReturnData(200, true, "Transaction Success", &id)
	return
}

func Send(ctx *gin.Context, req *model.WithdrawData) (res *model.ReturnData[primitive.ObjectID], err error) {
	if req.Destination == "" {
		err = model.NewError(400, "destination account is required")
		return
	}

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

	repoUser, err := helper.GetContextData[repository.User](ctx, "userRepo")
	if err != nil {
		err = model.NewError(500, "repository not found")
		return
	}

	repo, err := helper.GetContextData[repository.Transaction](ctx, "transactionRepo")
	if err != nil {
		err = model.NewError(500, "repository not found")
		return
	}

	repoPayment, err := helper.GetContextData[repository.PaymentAccount](ctx, "paymentRepo")
	if err != nil {
		err = model.NewError(500, "repository not found")
		return
	}

	tokenized := new(paseto.Payload[string])
	err = tokenGen.Decode(tokenHeader, tokenized)
	if err != nil {
		err = model.NewError(401, "unauthorized "+err.Error())
		return
	}

	user, err := repoUser.FindByUser(ctx, tokenized.Id)
	if err != nil {
		err = model.NewError(404, "user not found")
		return
	}

	paymentData, err := repoPayment.FindByIdUser(ctx, req.Account, user.Id.Hex())
	if err != nil {
		err = model.NewError(404, "user payment not found")
		return
	}

	destData, err := repoPayment.FindById(ctx, req.Destination)
	if err != nil {
		err = model.NewError(404, "destination payment not found")
		return

	}

	if paymentData.Balance < req.Amount {
		err = model.NewError(400, "insufficient balance")
		return
	}

	paymentData.Balance -= req.Amount
	destData.Balance += req.Amount

	err = repoPayment.Update(ctx, paymentData)
	if err != nil {
		err = model.NewError(500, "failed to update balance")
		return
	}

	err = repoPayment.Update(ctx, destData)
	if err != nil {
		err = model.NewError(500, "failed to update balance")
		return
	}

	id, err := repo.Create(ctx, model.Transaction{
		FromAccountID: paymentData.Id,
		ToAccountID:   destData.Id,
		Amount:        req.Amount,
		Type:          model.Transfer,
		Status:        model.Success,
		CreatedAt:     helper.NewTimestamp(),
		UpdatedAt:     primitive.Timestamp{},
	})
	if err != nil {
		err = model.NewError(500, "failed to create transaction")
		return
	}

	res = model.NewReturnData(200, true, "Transaction Success", &id)
	return
}
