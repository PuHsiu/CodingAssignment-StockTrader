package stock

import (
	"time"
	"trader/model"
	"trader/presentation/rest"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func (inst *instance) RestPrefix() string {
	return "stock"
}

func (inst *instance) RestRouter(subRouter *rest.Router) {
	subRouter.Post("/stock/order", inst.createOrder)
}

func (inst *instance) createOrder(ctx rest.Context) {
	var (
		// TODO: validate input to match valid enum
		req struct {
			TradeID  string                    `json:"tradeID"`
			StockID  string                    `json:"stockID"`
			Type     model.StockTradeOrderType `json:"type"`
			Quantity int                       `json:"quantity"`
			Price    decimal.Decimal           `json:"price"`
		}

		// TODO: issuer from context
		issuerID = uuid.MustParse("00000000-0000-0000-0000-000000000000")
	)

	ctx.BindBody(&req)
	// TODO: validate stock ownership

	stock, err := model.GetStockByID(req.StockID)
	if err != nil {
		ctx.BadRequest("NOT_FOUND", req.StockID)
		return
	}

	trade := model.NewStockHalfTrade(req.Type, model.StockHalfTradeMeta{
		HalfTradeID: uuid.MustParse(req.TradeID),
		IssuerID:    issuerID,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),

		Price:             req.Price,
		Quantity:          req.Quantity,
		CompletedQuantity: 0,
	})

	if trade == nil {
		ctx.BadRequest("INVALID")
		return
	}

	if err := inst.matchOrAddHalfTrade(stock, trade); err != nil {
		ctx.InternalError(err)
		return
	}

	ctx.Ok()
}

func (inst *instance) cancelOrder(ctx rest.Context) {
	// TBD
}
