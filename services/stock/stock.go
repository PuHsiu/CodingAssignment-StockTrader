package stock

import (
	"fmt"
	"trader/model"
	"trader/services"
	"trader/utils"
)

const (
	// TODO: buffer size by different stock
	traderBufferSize int = 100
)

type instance struct {
	// stockTraders write only when init, ignore mutex
	stockTraders         map[string]*stockTrader
	stockTradeNotifyChan chan string
}

func init() {
	stocks, err := model.GetAllStocks()
	if err != nil {
		panic(fmt.Sprintf("unable to get stock lists, err: %s", err.Error()))
	}

	instance := instance{
		stockTraders:         make(map[string]*stockTrader, len(stocks)),
		stockTradeNotifyChan: make(chan string),
	}

	for _, stock := range stocks {
		instance.stockTraders[stock.StockID] = &stockTrader{
			stock:        stock,
			serviceChan:  make(chan model.StockHalfTrade, traderBufferSize),
			tradeIDCache: make(map[string]struct{}),
		}
	}

	go instance.startTrader()

	services.Register(&instance)
}

func (inst *instance) startTrader() {
	fmt.Println("Start trader ... ")

	// TODO: different consumer on different stock

	for stockID := range inst.stockTradeNotifyChan {
		stockTrader := inst.stockTraders[stockID]
		trade := <-stockTrader.serviceChan

		if _, ok := stockTrader.tradeIDCache[trade.TradeID().String()]; !ok {
			stockTrader.tradeIDCache[trade.TradeID().String()] = struct{}{}
			stockTrader.matchOrAddHalfTrade(trade)
		}
	}
}

func (inst *instance) matchOrAddHalfTrade(stock model.Stock, trade model.StockHalfTrade) error {
	stockTrader, ok := inst.stockTraders[stock.StockID]
	if !ok {
		return fmt.Errorf("no such stock")
	}

	fmt.Printf("create trade %d * %s with id %s\n", trade.RemainingQuantity(), stock.StockID, trade.TradeID())

	inst.stockTradeNotifyChan <- stock.StockID
	stockTrader.serviceChan <- trade

	return nil
}

type stockTrader struct {
	stock model.Stock

	serviceChan chan model.StockHalfTrade

	// TODO: GC if timeout
	tradeIDCache map[string]struct{}

	// TODO: Enhance struct below for high freq insert & remove (e.g GC, tree, ...)
	// buyer: from highest price to lowest, from oldest to latest
	buyer []model.StockHalfTrade

	// seller: from lowest price to highest, from oldest to latest
	seller []model.StockHalfTrade
}

func (trader *stockTrader) matchOrAddHalfTrade(trade model.StockHalfTrade) error {
	var (
		matcheeHalfTradeList []model.StockHalfTrade
		matcherHalfTradeList []model.StockHalfTrade
	)

	switch trade.StrategyType() {
	case model.StockTradeOrderTypeLimitPriceBuy, model.StockTradeOrderTypeMarketPriceBuy:
		matcheeHalfTradeList = trader.seller
		matcherHalfTradeList = trader.buyer
	case model.StockTradeOrderTypeLimitPriceSell, model.StockTradeOrderTypeMarketPriceSell:
		matcheeHalfTradeList = trader.buyer
		matcherHalfTradeList = trader.seller
	}

	index := findFirstCurrentScopeIndex(matcheeHalfTradeList, trade)

	if index < len(matcheeHalfTradeList) {
		for i := index; i < len(matcheeHalfTradeList); i++ {
			if matcheeHalfTradeList[i].Completed() {
				fmt.Printf("%s is completed, ignored\n", matcheeHalfTradeList[i].TradeID())
				continue
			}

			trade, matcheeHalfTradeList[i] = trader.trade(trade, matcheeHalfTradeList[i])

			if trade.Completed() {
				break
			}
		}
	}

	if !trade.Completed() {
		index := findFirstNextScopeIndex(matcherHalfTradeList, trade)
		matcherHalfTradeList = utils.Insert(matcherHalfTradeList, index, trade)

		fmt.Printf("add trade into list:\n")
		for i := range matcherHalfTradeList {
			fmt.Printf("- %v\n", matcherHalfTradeList[i])
		}
	}

	switch trade.StrategyType() {
	case model.StockTradeOrderTypeLimitPriceBuy, model.StockTradeOrderTypeMarketPriceBuy:
		trader.seller = matcheeHalfTradeList
		trader.buyer = matcherHalfTradeList
	case model.StockTradeOrderTypeLimitPriceSell, model.StockTradeOrderTypeMarketPriceSell:
		trader.buyer = matcheeHalfTradeList
		trader.seller = matcherHalfTradeList
	}

	return nil
}

func (trader *stockTrader) trade(this model.StockHalfTrade, that model.StockHalfTrade) (model.StockHalfTrade, model.StockHalfTrade) {
	quantity := utils.Min(this.RemainingQuantity(), that.RemainingQuantity())

	fmt.Printf("trade %d * %s between %s and %s\n", quantity, trader.stock.StockID, this.TradeID().String(), that.TradeID().String())

	newThis := this.Trade(quantity)
	newThat := that.Trade(quantity)

	if newThis.Completed() {
		fmt.Printf("%s completed\n", newThis.TradeID())
	}

	if newThat.Completed() {
		fmt.Printf("%s completed\n", newThat.TradeID())
	}

	return newThis, newThat
}

// TODO: faster way to find index
func findFirstCurrentScopeIndex[T interface{ Match(T) bool }](slice []T, value T) int {
	for i := range slice {
		if slice[i].Match(value) {
			return i
		}
	}

	return len(slice)
}

// TODO: faster way to find index
func findFirstNextScopeIndex[T interface{ Equal(T) bool }](slice []T, value T) int {
	equal := false

	for i := range slice {
		if slice[i].Equal(value) {
			equal = true
		} else {
			if equal {
				return i
			}
		}
	}

	return len(slice)
}
