package model

import (
	"github.com/goccy/go-json"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionType string
type StatusTransaction string
type PaymentType string

const (
	Credit PaymentType = "credit"
	Debit  PaymentType = "debit"
	Loan   PaymentType = "loan"
)

const (
	Deposit  TransactionType = "deposit"
	Withdraw TransactionType = "withdraw"
	Transfer TransactionType = "transfer"
)

const (
	Success StatusTransaction = "success"
	Failed  StatusTransaction = "failed"
)

type Account struct {
	Id        primitive.ObjectID  `bson:"_id,omitempty"`
	Username  string              `bson:"username"`
	Hash      string              `bson:"hash"`
	Password  string              `bson:"password"`
	LastLogin primitive.Timestamp `bson:"last_login" json:"-"`
	CreatedAt primitive.Timestamp `bson:"created_at" json:"-"`
	UpdatedAt primitive.Timestamp `bson:"updated_at" json:"-"`
}

type PaymentAccount struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	AccountID primitive.ObjectID `bson:"account_id"`
	Balance   float64            `bson:"balance"`
}

type Transaction struct {
	Id            primitive.ObjectID  `bson:"_id,omitempty"`
	FromAccountID primitive.ObjectID  `bson:"from_account"`
	ToAccountID   primitive.ObjectID  `bson:"to_account"`
	Amount        float64             `bson:"amount"`
	Type          TransactionType     `bson:"type"`
	Status        StatusTransaction   `bson:"status"`
	CreatedAt     primitive.Timestamp `bson:"created_at" json:"-"`
	UpdatedAt     primitive.Timestamp `bson:"updated_at" json:"-"`
}

type FullData struct {
	Id          primitive.ObjectID  `bson:"_id,omitempty"`
	Username    string              `bson:"username"`
	Hash        string              `bson:"hash"`
	Password    string              `bson:"-"`
	LastLogin   primitive.Timestamp `bson:"last_login" json:"-"`
	CreatedAt   primitive.Timestamp `bson:"created_at" json:"-"`
	UpdatedAt   primitive.Timestamp `bson:"updated_at" json:"-"`
	Payment     PaymentAccount      `bson:"payment"`
	Transaction []Transaction       `bson:"transaction"`
}

func (fd FullData) MarshalJSON() ([]byte, error) {
	mapper := make(map[string]any)

	trans := make([]any, 0)
	for _, k := range fd.Transaction {
		trans = append(trans, map[string]any{
			"id":           k.Id.Hex(),
			"from_account": k.FromAccountID.Hex(),
			"to_account":   k.ToAccountID.Hex(),
			"amount":       k.Amount,
			"type":         k.Type,
			"status":       k.Status,
		})
	}

	mapper["id"] = fd.Id.Hex()
	mapper["username"] = fd.Username
	mapper["hash"] = fd.Hash

	mapper["payment"] = map[string]any{
		"id":         fd.Payment.Id.Hex(),
		"account_id": fd.Payment.AccountID.Hex(),
		"balance":    fd.Payment.Balance,
	}

	mapper["transaction"] = trans

	return json.Marshal(mapper)
}
