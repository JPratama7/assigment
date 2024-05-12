package repository

import (
	"be-assignment/model"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Transaction struct {
	coll *mongo.Collection
}

func NewTransaction(db *mongo.Database, coll string) *Transaction {
	return &Transaction{coll: db.Collection(coll)}
}

func (r *Transaction) Create(ctx context.Context, data model.Transaction) (id primitive.ObjectID, err error) {
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

func (r *Transaction) FindByAccount(ctx context.Context, account string) (data []model.Transaction, err error) {
	cur, err := r.coll.Find(ctx, bson.M{
		"account": account,
	})

	if err != nil {
		return
	}

	err = cur.All(ctx, &data)
	return

}
