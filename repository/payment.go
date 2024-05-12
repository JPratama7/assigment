package repository

import (
	"be-assignment/model"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentAccount struct {
	coll *mongo.Collection
}

func NewPaymentAccount(db *mongo.Database, coll string) *PaymentAccount {
	return &PaymentAccount{coll: db.Collection(coll)}
}

func (r *PaymentAccount) Create(ctx context.Context, data model.PaymentAccount) (id primitive.ObjectID, err error) {
	cur, err := r.coll.InsertOne(ctx, data)
	if err != nil {
		return
	}

	id, ok := cur.InsertedID.(primitive.ObjectID)
	if !ok {
		err = errors.New("failed to get inserted id")
		return
	}
	return
}

func (r *PaymentAccount) FindByAccount(ctx context.Context, account string) (data []model.PaymentAccount, err error) {
	cur, err := r.coll.Find(ctx, bson.M{
		"account": account,
	})

	if err != nil {
		return
	}

	err = cur.All(ctx, &data)
	return
}

func (r *PaymentAccount) FindByIdUser(ctx context.Context, account, user string) (data model.PaymentAccount, err error) {
	id, err := primitive.ObjectIDFromHex(account)
	if err != nil {
		return
	}
	userId, err := primitive.ObjectIDFromHex(user)

	m := bson.M{
		"_id":        id,
		"account_id": userId,
	}
	err = r.coll.FindOne(ctx, m).Decode(&data)

	return
}
func (r *PaymentAccount) FindById(ctx context.Context, account string) (data model.PaymentAccount, err error) {
	id, err := primitive.ObjectIDFromHex(account)
	if err != nil {
		return
	}

	m := bson.M{
		"_id": id,
	}
	err = r.coll.FindOne(ctx, m).Decode(&data)
	return
}

func (r *PaymentAccount) Update(ctx context.Context, data model.PaymentAccount) (err error) {
	_, err = r.coll.UpdateOne(ctx, bson.M{
		"_id":        data.Id,
		"account_id": data.AccountID,
	}, bson.M{"$set": data})
	return
}
