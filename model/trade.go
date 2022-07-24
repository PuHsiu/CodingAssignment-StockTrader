package model

import (
	"time"

	"github.com/google/uuid"

	"github.com/shopspring/decimal"
)

type StockTradeOrderType int

const (
	StockTradeOrderTypeLimitPriceSell StockTradeOrderType = iota
	StockTradeOrderTypeLimitPriceBuy
	StockTradeOrderTypeMarketPriceSell
	StockTradeOrderTypeMarketPriceBuy
)

type StockHalfTradeMeta struct {
	HalfTradeID uuid.UUID
	IssuerID    uuid.UUID

	CreatedAt time.Time
	UpdatedAt time.Time

	Price             decimal.Decimal
	Quantity          int
	CompletedQuantity int
}

func NewStockHalfTrade(typ StockTradeOrderType, meta StockHalfTradeMeta) StockHalfTrade {
	// TODO: calc market price for spec type
	switch typ {
	case StockTradeOrderTypeLimitPriceSell:
		return StockStrategySellLimit{StockHalfTradeMeta: meta}
	case StockTradeOrderTypeLimitPriceBuy:
		return StockStrategyBuyLimit{StockHalfTradeMeta: meta}
	case StockTradeOrderTypeMarketPriceSell:
		return StockStrategySellMarket{StockHalfTradeMeta: meta}
	case StockTradeOrderTypeMarketPriceBuy:
		return StockStrategyBuyMarket{StockHalfTradeMeta: meta}
	}

	return nil
}

type StockHalfTrade interface {
	TradeID() uuid.UUID
	TradePrice() decimal.Decimal

	Match(StockHalfTrade) bool
	Equal(StockHalfTrade) bool
	StrategyType() StockTradeOrderType

	Completed() bool

	RemainingQuantity() int
	Trade(int) StockHalfTrade
}

type StockStrategyBuyMarket struct {
	StockHalfTradeMeta
	StockHalfTrade
}

func (this StockStrategyBuyMarket) TradeID() uuid.UUID {
	return this.StockHalfTradeMeta.HalfTradeID
}

func (this StockStrategyBuyMarket) TradePrice() decimal.Decimal {
	return this.Price
}

func (this StockStrategyBuyMarket) Match(that StockHalfTrade) bool {
	switch that := that.(type) {
	case StockStrategySellLimit:
		if this.Price.GreaterThanOrEqual(that.Price) {
			return true
		}
	case StockStrategySellMarket:
		if this.Price.GreaterThanOrEqual(that.Price) {
			return true
		}
	}

	return false
}

func (this StockStrategyBuyMarket) Equal(that StockHalfTrade) bool {
	switch that := that.(type) {
	case StockStrategySellLimit:
		if this.Price.Equal(that.Price) {
			return true
		}
	case StockStrategySellMarket:
		if this.Price.Equal(that.Price) {
			return true
		}
	}

	return false
}

func (this StockStrategyBuyMarket) StrategyType() StockTradeOrderType {
	return StockTradeOrderTypeMarketPriceBuy
}

func (this StockStrategyBuyMarket) Completed() bool {
	return this.CompletedQuantity == this.Quantity
}

func (this StockStrategyBuyMarket) RemainingQuantity() int {
	return this.Quantity - this.CompletedQuantity
}

func (this StockStrategyBuyMarket) Trade(quantity int) StockHalfTrade {
	this.CompletedQuantity += quantity
	return this
}

type StockStrategySellMarket struct {
	StockHalfTradeMeta
	StockHalfTrade
}

func (this StockStrategySellMarket) TradeID() uuid.UUID {
	return this.StockHalfTradeMeta.HalfTradeID
}

func (this StockStrategySellMarket) TradePrice() decimal.Decimal {
	return this.Price
}

func (this StockStrategySellMarket) Match(that StockHalfTrade) bool {
	switch that := that.(type) {
	case StockStrategyBuyLimit:
		if this.Price.LessThanOrEqual(that.Price) {
			return true
		}
	case StockStrategyBuyMarket:
		if this.Price.LessThanOrEqual(that.Price) {
			return true
		}
	}

	return false
}

func (this StockStrategySellMarket) Equal(that StockHalfTrade) bool {
	switch that := that.(type) {
	case StockStrategyBuyLimit:
		if this.Price.Equal(that.Price) {
			return true
		}
	case StockStrategyBuyMarket:
		if this.Price.Equal(that.Price) {
			return true
		}
	}

	return false
}

func (this StockStrategySellMarket) StrategyType() StockTradeOrderType {
	return StockTradeOrderTypeMarketPriceSell
}

func (this StockStrategySellMarket) Completed() bool {
	return this.CompletedQuantity == this.Quantity
}

func (this StockStrategySellMarket) RemainingQuantity() int {
	return this.Quantity - this.CompletedQuantity
}

func (this StockStrategySellMarket) Trade(quantity int) StockHalfTrade {
	this.CompletedQuantity += quantity
	return this
}

type StockStrategyBuyLimit struct {
	StockHalfTradeMeta
	StockHalfTrade
}

func (this StockStrategyBuyLimit) TradeID() uuid.UUID {
	return this.StockHalfTradeMeta.HalfTradeID
}

func (this StockStrategyBuyLimit) TradePrice() decimal.Decimal {
	return this.Price
}

func (this StockStrategyBuyLimit) Match(that StockHalfTrade) bool {
	switch that := that.(type) {
	case StockStrategySellLimit:
		if this.Price.GreaterThanOrEqual(that.Price) {
			return true
		}
	case StockStrategySellMarket:
		if this.Price.GreaterThanOrEqual(that.Price) {
			return true
		}
	}

	return false
}

func (this StockStrategyBuyLimit) Equal(that StockHalfTrade) bool {
	switch that := that.(type) {
	case StockStrategySellLimit:
		if this.Price.Equal(that.Price) {
			return true
		}
	case StockStrategySellMarket:
		if this.Price.Equal(that.Price) {
			return true
		}
	}

	return false
}

func (this StockStrategyBuyLimit) StrategyType() StockTradeOrderType {
	return StockTradeOrderTypeLimitPriceBuy
}

func (this StockStrategyBuyLimit) Completed() bool {
	return this.CompletedQuantity == this.Quantity
}

func (this StockStrategyBuyLimit) RemainingQuantity() int {
	return this.Quantity - this.CompletedQuantity
}

func (this StockStrategyBuyLimit) Trade(quantity int) StockHalfTrade {
	this.CompletedQuantity += quantity
	return this
}

type StockStrategySellLimit struct {
	StockHalfTradeMeta
	StockHalfTrade
}

func (this StockStrategySellLimit) TradeID() uuid.UUID {
	return this.StockHalfTradeMeta.HalfTradeID
}

func (this StockStrategySellLimit) TradePrice() decimal.Decimal {
	return this.Price
}

func (this StockStrategySellLimit) Match(that StockHalfTrade) bool {
	switch that := that.(type) {
	case StockStrategyBuyLimit:
		if that.Price.GreaterThanOrEqual(this.Price) {
			return true
		}
	case StockStrategyBuyMarket:
		if that.Price.GreaterThanOrEqual(this.Price) {
			return true
		}
	}

	return false
}

func (this StockStrategySellLimit) Equal(that StockHalfTrade) bool {
	switch that := that.(type) {
	case StockStrategyBuyLimit:
		if that.Price.Equal(this.Price) {
			return true
		}
	case StockStrategyBuyMarket:
		if that.Price.Equal(this.Price) {
			return true
		}
	}

	return false
}

func (this StockStrategySellLimit) StrategyType() StockTradeOrderType {
	return StockTradeOrderTypeLimitPriceSell
}

func (this StockStrategySellLimit) Completed() bool {
	return this.CompletedQuantity == this.Quantity
}

func (this StockStrategySellLimit) RemainingQuantity() int {
	return this.Quantity - this.CompletedQuantity
}

func (this StockStrategySellLimit) Trade(quantity int) StockHalfTrade {
	this.CompletedQuantity += quantity
	return this
}
