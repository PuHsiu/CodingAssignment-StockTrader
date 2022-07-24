package stock

import (
	"testing"
	"trader/model"

	"github.com/shopspring/decimal"
)

func TestMatch(t *testing.T) {
	trader := stockTrader{
		stock:        model.Stock{"AAPL", decimal.NewFromInt(150), decimal.Zero, decimal.Zero},
		serviceChan:  make(chan model.StockHalfTrade),
		tradeIDCache: make(map[string]struct{}),
		buyer: []model.StockHalfTrade{
			model.NewStockHalfTrade(
				model.StockTradeOrderTypeLimitPriceBuy,
				model.StockHalfTradeMeta{
					Quantity: 2,
					Price:    decimal.NewFromInt(80),
				},
			),
			model.NewStockHalfTrade(
				model.StockTradeOrderTypeLimitPriceBuy,
				model.StockHalfTradeMeta{
					Quantity: 1,
					Price:    decimal.NewFromInt(80),
				},
			),
		},
		seller: []model.StockHalfTrade{},
	}

	trade := model.NewStockHalfTrade(
		model.StockTradeOrderTypeLimitPriceSell,
		model.StockHalfTradeMeta{
			Quantity: 1,
			Price:    decimal.NewFromInt(80),
		},
	)

	trader.matchOrAddHalfTrade(trade)

	if trader.buyer[0].RemainingQuantity() != 1 {
		t.Error("stock.matchOrAddHalfTrade FAIL")
		return
	}

	if len(trader.seller) != 0 {
		t.Error("stock.matchOrAddHalfTrade FAIL")
		return
	}

	t.Log("stock.matchOrAddHalfTrade PASS")
}

func TestMatch2(t *testing.T) {
	trader := stockTrader{
		stock:        model.Stock{"AAPL", decimal.NewFromInt(150), decimal.Zero, decimal.Zero},
		serviceChan:  make(chan model.StockHalfTrade),
		tradeIDCache: make(map[string]struct{}),
		buyer: []model.StockHalfTrade{
			model.NewStockHalfTrade(
				model.StockTradeOrderTypeLimitPriceBuy,
				model.StockHalfTradeMeta{
					Quantity: 2,
					Price:    decimal.NewFromInt(80),
				},
			),
			model.NewStockHalfTrade(
				model.StockTradeOrderTypeLimitPriceBuy,
				model.StockHalfTradeMeta{
					Quantity: 1,
					Price:    decimal.NewFromInt(80),
				},
			),
		},
		seller: []model.StockHalfTrade{},
	}

	trade := model.NewStockHalfTrade(
		model.StockTradeOrderTypeLimitPriceSell,
		model.StockHalfTradeMeta{
			Quantity: 1,
			Price:    decimal.NewFromInt(90),
		},
	)

	trader.matchOrAddHalfTrade(trade)

	if trader.buyer[0].RemainingQuantity() != 2 {
		t.Errorf("stock.matchOrAddHalfTrade FAIL")
		return
	}

	if trader.buyer[1].RemainingQuantity() != 1 {
		t.Error("stock.matchOrAddHalfTrade FAIL")
		return
	}

	if len(trader.seller) != 1 {
		t.Error("stock.matchOrAddHalfTrade FAIL")
		return
	}

	if trader.seller[0].RemainingQuantity() != 1 {
		t.Error("stock.matchOrAddHalfTrade FAIL")
		return
	}

	t.Log("stock.matchOrAddHalfTrade PASS")
}
