package model

import "github.com/shopspring/decimal"

// TODO: interfacize model

type Stock struct {
	StockID          string
	Price            decimal.Decimal
	OpeningPrice     decimal.Decimal
	LastClosingPrice decimal.Decimal
}

var stocks = []Stock{
	{"AAPL", decimal.NewFromInt(150), decimal.Zero, decimal.Zero},
	{"MSFT", decimal.NewFromInt(260), decimal.Zero, decimal.Zero},
	{"GOOG", decimal.NewFromInt(110), decimal.Zero, decimal.Zero},
}

func GetAllStocks() ([]Stock, error) {
	return stocks, nil
}

func GetStockByID(id string) (Stock, error) {
	for i := range stocks {
		if stocks[i].StockID == id {
			return stocks[i], nil
		}
	}

	return Stock{}, ErrNotFound
}
