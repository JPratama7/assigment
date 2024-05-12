package controller

import (
	"be-assignment/helper"
	"be-assignment/model"
	"be-assignment/repository"
	"fmt"
	"github.com/JPratama7/util/token/paseto"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func Login(ctx *gin.Context, req *model.Login) (res *model.ReturnData[string], err error) {

	repo, err := helper.GetContextData[repository.User](ctx, "userRepo")
	if err != nil {
		err = model.NewError(500, "repository not found")
		return
	}

	tokenGen, err := helper.GetContextData[paseto.PASETO](ctx, "token")
	if err != nil {
		err = model.NewError(500, "token generator not found")
		return
	}

	user, err := repo.FindByUser(ctx, req.Username)
	if err != nil {
		err = model.NewError(404, "user not found")
		return
	}

	valid := helper.VerifyPassword(req.Password, user.Password, user.Hash)
	if !valid {
		err = model.NewError(401, "invalid password")
		return
	}

	token, err := tokenGen.Encode(user.Username)
	if err != nil {
		err = model.NewError(500, "failed to generate token "+err.Error())
		return
	}

	res = model.NewReturnData(200, true, "Login Success", &token)
	return

}

func Register(ctx *gin.Context, req *model.Login) (res *model.ReturnData[primitive.ObjectID], err error) {

	repo, err := helper.GetContextData[repository.User](ctx, "userRepo")
	if err != nil {
		err = model.NewError(500, "repository not found")
		return
	}

	password, hash, err := helper.HashPassword(req.Password)
	if err != nil {
		err = model.NewError(500, "failed to hash password")
		return
	}

	_, err = repo.FindByUser(ctx, req.Username)
	if err == nil {
		err = model.NewError(404, "username duplicate")
		return
	}

	id, err := repo.Create(ctx, model.Account{
		Username:  req.Username,
		Hash:      hash,
		Password:  password,
		LastLogin: primitive.Timestamp{},
		CreatedAt: primitive.Timestamp{T: uint32(time.Now().Unix())},
		UpdatedAt: primitive.Timestamp{},
	})

	if err != nil {
		fmt.Println(err)
		err = model.NewError(404, "user not found")
		return
	}

	res = model.NewReturnData(200, true, "Login Success", &id)
	return
}

func AllTransaction(ctx *gin.Context) (res *model.ReturnData[[]model.FullData], err error) {
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

	repo, err := helper.GetContextData[repository.User](ctx, "userRepo")
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

	user, err := repo.FindByUser(ctx, tokenized.Id)
	if err != nil {
		err = model.NewError(404, "user not found")
		return
	}

	data, err := repo.FindWithTransactionAccount(ctx, user.Username)
	if err != nil {
		err = model.NewError(404, "transaction not found")
		return
	}

	res = model.NewReturnData(200, true, "Transaction Success", &data)
	return
}
