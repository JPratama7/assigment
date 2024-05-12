package repository

import (
	"be-assignment/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type User struct {
	coll *mongo.Collection
}

func NewUser(db *mongo.Database, coll string) *User {
	return &User{coll: db.Collection(coll)}
}

func (u *User) FindByUser(ctx context.Context, username string) (data model.Account, err error) {
	err = u.coll.FindOne(ctx, bson.M{
		"username": username,
	}).Decode(&data)

	if err == nil {
		data.LastLogin = primitive.Timestamp{T: uint32(time.Now().Unix())}
		_, err = u.coll.UpdateOne(ctx, bson.M{"username": username}, bson.M{"$set": bson.M{"last_login": data.LastLogin}})
	}
	return
}

func (u *User) FindByEmail(ctx context.Context, email string) (data model.Account, err error) {
	err = u.coll.FindOne(ctx, bson.M{
		"email": email,
	}).Decode(&data)

	return
}

func (u *User) Create(ctx context.Context, data model.Account) (id primitive.ObjectID, err error) {
	data.CreatedAt = primitive.Timestamp{T: uint32(time.Now().Unix())}

	cur, err := u.coll.InsertOne(ctx, data)
	if err != nil {
		return
	}
	id = cur.InsertedID.(primitive.ObjectID)
	return
}

func (u *User) Update(ctx context.Context, username string, data model.Account) (err error) {
	data.UpdatedAt = primitive.Timestamp{T: uint32(time.Now().Unix())}

	_, err = u.coll.UpdateOne(ctx, bson.M{"username": username}, bson.M{"$set": data})
	return
}

func (u *User) FindWithTransactionAccount(ctx context.Context, username string) (data []model.FullData, err error) {

	cur, err := u.coll.Aggregate(ctx, bson.A{
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "payment_accounts"},
					{"localField", "_id"},
					{"foreignField", "account_id"},
					{"as", "payment"},
				},
			},
		},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$payment"},
					{"preserveNullAndEmptyArrays", true},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "transactions"},
					{"localField", "payment._id"},
					{"foreignField", "from_account"},
					{"as", "transaction"},
				},
			},
		},
		bson.D{{"$match", bson.D{{"username", username}}}},
	})

	if err != nil {
		err = model.NewError(500, "failed to get user transaction")
		return
	}

	err = cur.All(ctx, &data)
	return
}
