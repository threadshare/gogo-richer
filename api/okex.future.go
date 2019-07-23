package api

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/miaolz123/conver"
	"github.com/phonegapX/QuantBot/constant"
	"github.com/phonegapX/QuantBot/model"
)

// OkexFuture the exchange struct of okex.com future
type OkexFuture struct {
	stockTypeMap        map[string][2]string
	tradeTypeMap        map[string]string
	tradeTypeAntiMap    map[int]string
	tradeTypeLogMap     map[string]string
	contractTypeAntiMap map[string]string
	leverageMap         map[string]string
	recordsPeriodMap    map[string]string
	records             map[string][]Record
	host                string
	logger              model.Logger
	option              Option

	limit     float64
	lastSleep int64
	lastTimes int64
}

// NewOkexFuture create an exchange struct of okex.com
func NewOkexFuture(opt Option) Exchange {
	return &OkexFuture{
		stockTypeMap: map[string][2]string{
			"BTC.WEEK/USD":   {"btc_usd", "this_week"},
			"BTC.WEEK2/USD":  {"btc_usd", "next_week"},
			"BTC.MONTH3/USD": {"btc_usd", "quarter"},
			"LTC.WEEK/USD":   {"ltc_usd", "this_week"},
			"LTC.WEEK2/USD":  {"ltc_usd", "next_week"},
			"LTC.MONTH3/USD": {"ltc_usd", "quarter"},
		},
		tradeTypeMap: map[string]string{
			constant.TradeTypeLong:       "1",
			constant.TradeTypeShort:      "2",
			constant.TradeTypeLongClose:  "3",
			constant.TradeTypeShortClose: "4",
		},
		tradeTypeAntiMap: map[int]string{
			1: constant.TradeTypeLong,
			2: constant.TradeTypeShort,
			3: constant.TradeTypeLongClose,
			4: constant.TradeTypeShortClose,
		},
		tradeTypeLogMap: map[string]string{
			constant.TradeTypeLong:       constant.LONG,
			constant.TradeTypeShort:      constant.SHORT,
			constant.TradeTypeLongClose:  constant.LONGCLOSE,
			constant.TradeTypeShortClose: constant.SHORTCLOSE,
		},
		contractTypeAntiMap: map[string]string{
			"this_week": "WEEK",
			"next_week": "WEEK2",
			"quarter":   "MONTH3",
		},
		leverageMap: map[string]string{
			"10": "10",
			"20": "20",
		},
		recordsPeriodMap: map[string]string{
			"M":   "1min",
			"M3":  "3min",
			"M5":  "5min",
			"M15": "15min",
			"M30": "30min",
			"H":   "1hour",
			"H2":  "2hour",
			"H4":  "4hour",
			"H6":  "6hour",
			"H12": "12hour",
			"D":   "1day",
			"D3":  "3day",
			"W":   "1week",
		},
		records: make(map[string][]Record),
		host:    "https://www.okex.com/api/v1/",
		logger:  model.Logger{TraderID: opt.TraderID, ExchangeType: opt.Type},
		option:  opt,

		limit:     10.0,
		lastSleep: time.Now().UnixNano(),
	}
}

// Log print something to console
func (e *OkexFuture) Log(msgs ...interface{}) {
	e.logger.Log(constant.INFO, "", 0.0, 0.0, msgs...)
}

// GetType get the type of this exchange
func (e *OkexFuture) GetType() string {
	return e.option.Type
}

// GetName get the name of this exchange
func (e *OkexFuture) GetName() string {
	return e.option.Name
}

// SetLimit set the limit calls amount per second of this exchange
func (e *OkexFuture) SetLimit(times interface{}) float64 {
	e.limit = conver.Float64Must(times)
	return e.limit
}

// AutoSleep auto sleep to achieve the limit calls amount per second of this exchange
func (e *OkexFuture) AutoSleep() {
	now := time.Now().UnixNano()
	interval := 1e+9/e.limit*conver.Float64Must(e.lastTimes) - conver.Float64Must(now-e.lastSleep)
	if interval > 0.0 {
		time.Sleep(time.Duration(conver.Int64Must(interval)))
	}
	e.lastTimes = 0
	e.lastSleep = now
}

// GetMinAmount get the min trade amonut of this exchange
func (e *OkexFuture) GetMinAmount(stock string) float64 {
	return 1.0
}

func (e *OkexFuture) getAuthJSON(url string, params []string) (json *simplejson.Json, err error) {
	e.lastTimes++
	params = append(params, "api_key="+e.option.AccessKey)
	sort.Strings(params)
	params = append(params, "secret_key="+e.option.SecretKey)
	params = append(params, "sign="+strings.ToUpper(signMd5(params)))
	resp, err := post(url, params)
	if err != nil {
		return
	}
	return simplejson.NewJson(resp)
}

// GetAccount get the account detail of this exchange
func (e *OkexFuture) GetAccount() interface{} {
	json, err := e.getAuthJSON(e.host+"future_userinfo.do", []string{})
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetAccount() error, ", err)
		return false
	}
	if result := json.Get("result").MustBool(); !result {
		err = fmt.Errorf("GetAccount() error, the error number is %v", json.Get("error_code").MustInt())
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetAccount() error, ", err)
		return false
	}
	return map[string]float64{
		"BTC":       conver.Float64Must(json.GetPath("info", "btc", "account_rights").Interface()),
		"FrozenBTC": 0.0,
		"LTC":       conver.Float64Must(json.GetPath("info", "ltc", "account_rights").Interface()),
		"FrozenLTC": 0.0,
	}
}

// GetPositions get the positions detail of this exchange
func (e *OkexFuture) GetPositions(stockType string) interface{} {
	stockType = strings.ToUpper(stockType)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetPositions() error, unrecognized stockType: ", stockType)
		return false
	}
	positions := []Position{}
	params := []string{
		"symbol=" + e.stockTypeMap[stockType][0],
		"contract_type=" + e.stockTypeMap[stockType][1],
	}
	json, err := e.getAuthJSON(e.host+"future_position.do", params)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetPositions() error, ", err)
		return false
	}
	if result := json.Get("result").MustBool(); !result {
		err = fmt.Errorf("GetPositions() error, the error number is %v", json.Get("error_code").MustInt())
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetPositions() error, ", err)
		return false
	}
	positionsJSON := json.Get("holding")
	count := len(positionsJSON.MustArray())
	for i := 0; i < count; i++ {
		positionJSON := positionsJSON.GetIndex(i)
		side := "sell"
		tradeType := constant.TradeTypeShort
		amount := conver.Float64Must(positionJSON.Get("buy_amount").Interface())
		if amount > 0 {
			side = "buy"
			tradeType = constant.TradeTypeLong
		} else if amount = conver.Float64Must(positionJSON.Get("sell_amount").Interface()); amount == 0.0 {
			continue
		}
		positions = append(positions, Position{
			Price:         conver.Float64Must(positionJSON.Get(side + "_price_avg").Interface()),
			Leverage:      conver.IntMust(positionJSON.Get("lever_rate").Interface()),
			Amount:        conver.Float64Must(positionJSON.Get(side + "_amount").Interface()),
			ConfirmAmount: conver.Float64Must(positionJSON.Get(side + "_available").Interface()),
			FrozenAmount:  0.0,
			Profit:        conver.Float64Must(positionJSON.Get(side + "_profit_real").Interface()),
			ContractType:  e.contractTypeAntiMap[positionJSON.Get("contract_type").MustString()],
			TradeType:     tradeType,
			StockType:     stockType,
		})
	}
	return positions
}

// Trade place an order
func (e *OkexFuture) Trade(tradeType string, stockType string, _price, _amount interface{}, msgs ...interface{}) interface{} {
	tradeType = strings.ToUpper(tradeType)
	stockType = strings.ToUpper(stockType)
	price := conver.Float64Must(_price)
	amount := conver.Float64Must(_amount)
	if _, ok := e.tradeTypeMap[tradeType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Trade() error, unrecognized tradeType: ", tradeType)
		return false
	}
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Trade() error, unrecognized stockType: ", stockType)
		return false
	}
	if len(msgs) < 1 {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Trade() error, unrecognized leverage")
		return false
	}
	leverage := fmt.Sprint(msgs[0])
	if _, ok := e.leverageMap[leverage]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Trade() error, unrecognized leverage: ", leverage)
		return false
	}
	matchPrice := "match_price=1"
	if price > 0.0 {
		matchPrice = "match_price=0"
	} else {
		price = 0.0
	}
	params := []string{
		"symbol=" + e.stockTypeMap[stockType][0],
		"contract_type=" + e.stockTypeMap[stockType][1],
		fmt.Sprintf("price=%f", price),
		fmt.Sprintf("amount=%f", amount),
		"type=" + e.tradeTypeMap[tradeType],
		matchPrice,
		"lever_rate=" + leverage,
	}
	json, err := e.getAuthJSON(e.host+"future_trade.do", params)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Trade() error, ", err)
		return false
	}
	if result := json.Get("result").MustBool(); !result {
		err = fmt.Errorf("Trade() error, the error number is %v", json.Get("error_code").MustInt())
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Trade() error, ", err)
		return false
	}
	e.logger.Log(e.tradeTypeLogMap[tradeType], stockType, price, amount, msgs[2:]...)
	return fmt.Sprint(json.Get("order_id").Interface())
}

// GetOrder get details of an order
func (e *OkexFuture) GetOrder(stockType, id string) interface{} {
	stockType = strings.ToUpper(stockType)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrder() error, unrecognized stockType: ", stockType)
		return false
	}
	params := []string{
		"symbol=" + e.stockTypeMap[stockType][0],
		"contract_type=" + e.stockTypeMap[stockType][1],
		"order_id=" + id,
	}
	json, err := e.getAuthJSON(e.host+"future_orders_info.do", params)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrder() error, ", err)
		return false
	}
	if result := json.Get("result").MustBool(); !result {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrder() error, the error number is ", json.Get("error_code").MustInt())
		return false
	}
	ordersJSON := json.Get("orders")
	if len(ordersJSON.MustArray()) > 0 {
		orderJSON := ordersJSON.GetIndex(0)
		return Order{
			ID:         fmt.Sprint(orderJSON.Get("order_id").Interface()),
			Price:      orderJSON.Get("price").MustFloat64(),
			Amount:     orderJSON.Get("amount").MustFloat64(),
			DealAmount: orderJSON.Get("deal_amount").MustFloat64(),
			Fee:        orderJSON.Get("fee").MustFloat64(),
			TradeType:  e.tradeTypeAntiMap[orderJSON.Get("type").MustInt()],
			StockType:  stockType,
		}
	}
	return false
}

// GetOrders get all unfilled orders
func (e *OkexFuture) GetOrders(stockType string) interface{} {
	stockType = strings.ToUpper(stockType)
	orders := []Order{}
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrders() error, unrecognized stockType: ", stockType)
		return false
	}
	params := []string{
		"symbol=" + e.stockTypeMap[stockType][0],
		"contract_type=" + e.stockTypeMap[stockType][1],
		"status=1",
		"order_id=-1",
		"current_page=1",
		"page_length=50",
	}
	json, err := e.getAuthJSON(e.host+"future_order_info.do", params)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrders() error, ", err)
		return false
	}
	if result := json.Get("result").MustBool(); !result {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrders() error, the error number is ", json.Get("error_code").MustInt())
		return false
	}
	ordersJSON := json.Get("orders")
	count := len(ordersJSON.MustArray())
	for i := 0; i < count; i++ {
		orderJSON := ordersJSON.GetIndex(i)
		orders = append(orders, Order{
			ID:         fmt.Sprint(orderJSON.Get("order_id").Interface()),
			Price:      orderJSON.Get("price").MustFloat64(),
			Amount:     orderJSON.Get("amount").MustFloat64(),
			DealAmount: orderJSON.Get("deal_amount").MustFloat64(),
			Fee:        orderJSON.Get("fee").MustFloat64(),
			TradeType:  e.tradeTypeAntiMap[orderJSON.Get("type").MustInt()],
			StockType:  stockType,
		})
	}
	return orders
}

// GetTrades get all filled orders recently
func (e *OkexFuture) GetTrades(stockType string) interface{} {
	stockType = strings.ToUpper(stockType)
	orders := []Order{}
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetTrades() error, unrecognized stockType: ", stockType)
		return false
	}
	params := []string{
		"symbol=" + e.stockTypeMap[stockType][0],
		"contract_type=" + e.stockTypeMap[stockType][1],
		"status=2",
		"order_id=-1",
		"current_page=1",
		"page_length=50",
	}
	json, err := e.getAuthJSON(e.host+"future_order_info.do", params)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetTrades() error, ", err)
		return false
	}
	if result := json.Get("result").MustBool(); !result {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetTrades() error, the error number is ", json.Get("error_code").MustInt())
		return false
	}
	ordersJSON := json.Get("orders")
	count := len(ordersJSON.MustArray())
	for i := 0; i < count; i++ {
		orderJSON := ordersJSON.GetIndex(i)
		orders = append(orders, Order{
			ID:         fmt.Sprint(orderJSON.Get("order_id").Interface()),
			Price:      orderJSON.Get("price").MustFloat64(),
			Amount:     orderJSON.Get("amount").MustFloat64(),
			DealAmount: orderJSON.Get("deal_amount").MustFloat64(),
			Fee:        orderJSON.Get("fee").MustFloat64(),
			TradeType:  e.tradeTypeAntiMap[orderJSON.Get("type").MustInt()],
			StockType:  stockType,
		})
	}
	return orders
}

// CancelOrder cancel an order
func (e *OkexFuture) CancelOrder(order Order) bool {
	params := []string{
		"symbol=" + e.stockTypeMap[order.StockType][0],
		"order_id=" + order.ID,
		"contract_type=" + e.stockTypeMap[order.StockType][1],
	}
	json, err := e.getAuthJSON(e.host+"future_cancel.do", params)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "CancelOrder() error, ", err)
		return false
	}
	if result := json.Get("result").MustBool(); !result {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "CancelOrder() error, the error number is ", json.Get("error_code").MustInt())
		return false
	}
	e.logger.Log(constant.CANCEL, order.StockType, order.Price, order.Amount-order.DealAmount, order)
	return true
}

// getTicker get market ticker & depth
func (e *OkexFuture) getTicker(stockType string, sizes ...interface{}) (ticker Ticker, err error) {
	stockType = strings.ToUpper(stockType)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		err = fmt.Errorf("GetTicker() error, unrecognized stockType: %+v", stockType)
		return
	}
	size := 20
	if len(sizes) > 0 && conver.IntMust(sizes[0]) > 0 {
		size = conver.IntMust(sizes[0])
	}
	resp, err := get(fmt.Sprintf("%vfuture_depth.do?symbol=%v&contract_type=%v&size=%v", e.host, e.stockTypeMap[stockType][0], e.stockTypeMap[stockType][1], size))
	if err != nil {
		err = fmt.Errorf("GetTicker() error, %+v", err)
		return
	}
	json, err := simplejson.NewJson(resp)
	if err != nil {
		err = fmt.Errorf("GetTicker() error, %+v", err)
		return
	}
	depthsJSON := json.Get("bids")
	for i := 0; i < len(depthsJSON.MustArray()); i++ {
		depthJSON := depthsJSON.GetIndex(i)
		ticker.Bids = append(ticker.Bids, OrderBook{
			Price:  depthJSON.GetIndex(0).MustFloat64(),
			Amount: depthJSON.GetIndex(1).MustFloat64(),
		})
	}
	depthsJSON = json.Get("asks")
	for i := len(depthsJSON.MustArray()); i > 0; i-- {
		depthJSON := depthsJSON.GetIndex(i - 1)
		ticker.Asks = append(ticker.Asks, OrderBook{
			Price:  depthJSON.GetIndex(0).MustFloat64(),
			Amount: depthJSON.GetIndex(1).MustFloat64(),
		})
	}
	if len(ticker.Bids) < 1 || len(ticker.Asks) < 1 {
		err = fmt.Errorf("GetTicker() error, can not get enough Bids or Asks")
		return
	}
	ticker.Buy = ticker.Bids[0].Price
	ticker.Sell = ticker.Asks[0].Price
	ticker.Mid = (ticker.Buy + ticker.Sell) / 2
	return
}

// GetTicker get market ticker & depth
func (e *OkexFuture) GetTicker(stockType string, sizes ...interface{}) interface{} {
	ticker, err := e.getTicker(stockType, sizes...)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, err)
		return false
	}
	return ticker
}

// GetRecords get candlestick data
func (e *OkexFuture) GetRecords(stockType, period string, sizes ...interface{}) interface{} {
	stockType = strings.ToUpper(stockType)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetRecords() error, unrecognized stockType: ", stockType)
		return false
	}
	if _, ok := e.recordsPeriodMap[period]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetRecords() error, unrecognized period: ", period)
		return false
	}
	size := 200
	if len(sizes) > 0 && conver.IntMust(sizes[0]) > 0 {
		size = conver.IntMust(sizes[0])
	}
	resp, err := get(fmt.Sprintf("%vfuture_kline.do?symbol=%v&contract_type=%v&type=%v&size=%v", e.host, e.stockTypeMap[stockType][0], e.stockTypeMap[stockType][1], e.recordsPeriodMap[period], size))
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetRecords() error, ", err)
		return false
	}
	json, err := simplejson.NewJson(resp)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetRecords() error, ", err)
		return false
	}
	timeLast := int64(0)
	if len(e.records[period]) > 0 {
		timeLast = e.records[period][len(e.records[period])-1].Time
	}
	recordsNew := []Record{}
	for i := len(json.MustArray()); i > 0; i-- {
		recordJSON := json.GetIndex(i - 1)
		recordTime := recordJSON.GetIndex(0).MustInt64() / 1000
		if recordTime > timeLast {
			recordsNew = append([]Record{{
				Time:   recordTime,
				Open:   recordJSON.GetIndex(1).MustFloat64(),
				High:   recordJSON.GetIndex(2).MustFloat64(),
				Low:    recordJSON.GetIndex(3).MustFloat64(),
				Close:  recordJSON.GetIndex(4).MustFloat64(),
				Volume: recordJSON.GetIndex(5).MustFloat64(),
			}}, recordsNew...)
		} else if timeLast > 0 && recordTime == timeLast {
			e.records[period][len(e.records[period])-1] = Record{
				Time:   recordTime,
				Open:   recordJSON.GetIndex(1).MustFloat64(),
				High:   recordJSON.GetIndex(2).MustFloat64(),
				Low:    recordJSON.GetIndex(3).MustFloat64(),
				Close:  recordJSON.GetIndex(4).MustFloat64(),
				Volume: recordJSON.GetIndex(5).MustFloat64(),
			}
		} else {
			break
		}
	}
	e.records[period] = append(e.records[period], recordsNew...)
	if len(e.records[period]) > size {
		e.records[period] = e.records[period][len(e.records[period])-size : len(e.records[period])]
	}
	return e.records[period]
}
