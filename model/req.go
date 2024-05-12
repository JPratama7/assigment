package model

type Login struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type WithdrawData struct {
	Account     string  `json:"account" validate:"required"`
	Destination string  `json:"destination"`
	Amount      float64 `json:"amount" validate:"required"`
}

type PaymentAccountData struct {
	Type    TransactionType `json:"type" validate:"required" enum:"credit,debit,loan"`
	Balance float64         `json:"balance" validate:"required"`
}
