# QuantBot 服务

## 项目地址

https://github.com/phonegapX/QuantBot

## 登录

QuantBot运行后，打开 `http://localhost:9876`。

默认的用户名和密码都是`admin`，请在登录后立即修改！

## 支持的交易所

| 交易所 | 货币类型 |
| -------- | ----- |
| zb | `BTC/USDT`, `ETH/USDT`, `EOS/USDT`, `LTC/USDT`, `QTUM/USDT` |
| okex | `BTC/USDT`, `ETH/USDT`, `EOS/USDT`, `ONT/USDT`, `QTUM/USDT`, `ONT/ETH` |
| 火币网 | `BTC/USDT`, `ETH/USDT`, `EOS/USDT`, `ONT/USDT`, `QTUM/USDT` |
| 比特儿国际 | `BTC/USDT`, `ETH/USDT`, `EOS/USDT`, `ONT/USDT`, `QTUM/USDT` |
| 币安 | `BTC/USDT`, `ETH/USDT`, `EOS/USDT`, `ONT/USDT`, `QTUM/USDT` |
| poloniex | `ETH/BTC`, `XMR/BTC`, `BTC/USDT`, `LTC/BTC`, `ETC/BTC`, `XRP/BTC`, `ETH/USDT`, `ETC/ETH`, ... |
| okex 期货 | `BTC.WEEK/USD`, `BTC.WEEK2/USD`, `BTC.MONTH3/USD`, `LTC.WEEK/USD`, ... |
| BigONE | `BTC/USDT`, `ONE/USDT`, `EOS/USDT`, `ETH/USDT`, `BCH/USDT`, `EOS/ETH` |

# 算法策略编写说明

## 语法规则

### 全局常量

| 名称 | 类型 | 说明 |
| ---- | ---- | ---- |
| Global/G| Object | 一个拥有各种全局方法的结构体 |
| Exchange/E | Object | 一个拥有各种交易所方法的结构体 |
| Exchanges/Es | Object List | 一个 `Exchange/E` 列表 |

### 交易类型

| 名称 | 类型 | 说明 |
| ---- | ---- | ---- |
| BUY | String | 买入交易 |
| SELL | String | 卖出交易 |
| LONG | String | 做多合约交易 |
| SHORT | String | 做空合约交易 |
| LONG_CLOSE | String | 平多合约交易 |
| SHORT_CLOSE | String | 平空合约交易 |

### K线周期

| 名称 | 类型 | 说明 |
| ---- | ---- | ---- |
| M | String | 1 分钟 |
| M5 | String | 5 分钟 |
| M15 | String | 15 分钟 |
| M30 | String | 30 分钟 |
| H | String | 1 小时 |
| D | String | 1 天 |
| W | String | 1 周 |

## 数据结构

### Account

| 名称 | 类型 | 说明 |
| ---- | ---- | ---- | 
| USDT | Number | 可用的 USDT 数量 |
| FrozenUSDT | Number | 冻结的 USDT 数量 |
| BTC | Number | 可用的 BTC 数量 |
| FrozenBTC | Number | 冻结的 BTC 数量 |
| LTC | Number | 可用的 LTC 数量 |
| FrozenLTC | Number | 冻结的 LTC 数量 |
| ... | Number | 可用的 ... 数量 |
| Frozen... | Number | 冻结的 ... 数量 |

### Position

| 名称 | 类型 | 说明 |
| ---- | ---- | ---- |
| Price | Number | 价格 |
| Leverage | Number | 杠杆比例 |
| Amount | Number | 总合约数量 |
| FrozenAmount | Number | 冻结的合约数量 |
| Profit | Number | 收益 |
| ContractType | String | 合约类型 |
| TradeType | String | 交易类型 |
| StockType | String | 货币类型 |

### Order

| 名称 | 类型 | 说明 |
| ---- | ---- | ---- |
| ID | String | 唯一 ID |
| Price | Number | 价格 |
| Amount | Number | 总量 |
| DealAmount | Number | 成交量 |
| Fee | Number | 这个订单的交易费 |
| TradeType | String | 交易类型 |
| StockType | String | 货币类型 |

### Record

| 名称 | 类型 | 说明 |
| ---- | ---- | ---- |
| Time | Number | unix 时间戳 |
| Open | Number | 开盘价 |
| High | Number | 最高价 |
| Low | Number | 最低价 |
| Close | Number | 收盘价 |
| Volume | Number | 交易量 |

### OrderBook

| 名称 | 类型 | 说明 |
| ---- | ---- | ---- |
| Price | Number | 价格 |
| Amount | Number | 市场深度量 |

### Ticker

| 名称 | 类型 | 说明 |
| ---- | ---- | ---- |
| Bids | OrderBook List | 买单市场深度列表 |
| Buy | Number | 买一价, `Bids[0].Price` |
| Mid | Number | `(Buy + Sell) / 2` |
| Sell | Number | 卖一价, `Asks[0].Price` |
| Asks | OrderBook List | 卖单市场深度列表 |

## Global/G

`Global`/`G` 是一个拥有各种全局方法的结构体。

### Sleep

> G.Sleep(Interval: *Any*) => *No Return*

```javascript
// 程序将休眠 5 秒
// 如果 Interval <= 0, 将自动执行所有交易所的 AutoSleep() 方法
G.Sleep(5000);
```

### Log

> G.Log(Message: *Any*) => *No Return*

```javascript
// 向管理台发送打印信息
G.Log("I'm running…");
```

### Console

> G.Console(Message: *Any*) => *No Return*

```javascript
// 向控制台发送打印信息
G.Console("I'm running…");
```

### LogProfit

> G.LogProfit(Profit: *Number*, Message: *Any*) => *No Return*

```javascript
// 向管理台发送收益信息，用来生成收益图表
G.LogProfit(12.345, 'Round 1 end');
```

### LogStatus

> G.LogStatus(Message: *Any*) => *No Return*

```javascript
// 向管理台发送实时状态信息
G.LogStatus('Latest BTC Ticker: ', E.GetTicker('BTC/USD'));
```

### AddTask

> G.AddTask(group: *String*, FunctionName: *String*, Arguments: *Any*) => *Boolean*

```javascript
// 和 G.ExecTasks() 配合使用
```

### BindTaskParam

> G.BindTaskParam(group: *String*, FunctionName: *String*, Arguments: *Any*) => *Boolean*

```javascript
// 和 G.ExecTasks() 配合使用
```

### ExecTasks

> G.ExecTasks(group: *String*) => *List*

```javascript
// 添加几个任务到任务列表里面
G.AddTask("myGroup", "function1");
G.AddTask("myGroup", "function2", 'param1');

// 可以随时为函数绑定新的参数
G.BindTaskParam("myGroup", "function1", 'param1');
G.BindTaskParam("myGroup", "function2", 'param2', 'param3');

// 执行一组任务列表里面的所有数据并返回所有的执行结果
var results = G.ExecTasks("myGroup");
var r1 = results[0];
var r2 = results[1];
```

## Exchange/E

`Exchange`/`E` 是一个拥有各种交易所方法的结构体。

### Log

> E.Log(Message: *Any*) => *No Return*

```javascript
// 向管理台发送这个交易所的打印信息
E.Log("I'm running…");
```

### GetType

> E.GetType() => *String*

```javascript
// 获取交易所类型
var thisType = E.GetType();
```

### GetName

> E.GetName() => *String*

```javascript
// 获取交易所名称
var thisName = E.GetName();
```

### GetMainStock

> E.GetMainStock() => *String*

```javascript
// 获取交易所的默认货币类型
var thisMainStock = E.GetMainStock();
```

### SetMainStock

> E.SetMainStock(StockType: *String*) => *String*

```javascript
// 设置交易所的默认货币类型
var newMainStockType = E.SetMainStock('LTC/USD');
```

### SetLimit

> E.SetLimit(times: *Number*) => *Number*

```javascript
// 设置交易所的API访问频率
// 和 E.AutoSleep() 配合使用
var newLimit = E.SetLimit(6);
```

### AutoSleep

> E.AutoSleep() => *No Return*

```javascript
// 自动休眠以满足设置的交易所的API访问频率
E.AutoSleep();
```

### GetAccount

> E.GetAccount() => *Account*

```javascript
// 获取交易所的账户资金信息
var thisAccount = E.GetAccount();
```

### GetPositions

> E.GetPositions(StockType: *String*) => *Position List*

```javascript
// 获取交易所的合约列表
var thisPositions = E.GetPositions('BTC/USD');
```

### GetMinAmount

> E.GetMinAmount(StockType: *String*) => *Number*

```javascript
// 获取交易所的最小交易数量
var thisMinAmount = E.GetMinAmount('BTC/USD');
```

### Trade

> E.Trade(TradeType: [*String*](#trade-type), StockType: *String*, Price: *Number*, Amount: *Number*, Message: *Any*) => *String*/*Boolean*

```javascript
// 买入示例
// 如果 Price <= 0 自动设置为市价单，数量参数也有所不同
// 如果成功返回订单的 ID
// 如果失败返回 false
E.Trade('BUY', 'BTC/USD', 600, 0.5, 'I paid $300'); // 限价单
E.Trade('BUY', 'BTC/USD', 0, 300, 'I also paid $300'); // 市价单

// 卖出示例
// 如果 Price <= 0 自动设置为市价单
// 如果成功返回订单的 ID
// 如果失败返回 false
E.Trade('SELL', 'BTC/USD', 600, 0.5); // 限价单
E.Trade('SELL', 'BTC/USD', 0, 0.5); // 市价单
```

### GetOrder

> E.GetOrder(StockType: *String*, ID: *String*) => *Order*/*Boolean*

```javascript
// 如果成功返回订单信息
// 如果失败返回 false
var thisOrder = E.GetOrder('BTC/USD', 'XXXXXX');
```

### GetOrders

> E.GetOrders(StockType: *String*) => *Order List*

```javascript
// 返回所有的未完成订单列表
var thisOrders = E.GetOrders('BTC/USD');
```

### GetTrades

> E.GetTrades(StockType: *String*) => *Order List*

```javascript
// 返回最近的已完成订单列表
var thisTrades = E.GetTrades('BTC/USD');
```

### CancelOrder

> E.CancelOrder(Order: *Order*) => *Boolean*

```javascript
var thisOrders = E.GetOrders('BTC/USD');
for (var i = 0; i < thisOrders.length; i++) {
    // 返回是否取消成功的结果
    var isCanceled = E.CancelOrder(thisOrders[i]);
}
```

### GetTicker

> E.GetTicker(StockType: *String*, Size: *Any*) => *Ticker*

```javascript
// 获取交易所的最新市场行情数据
var thisTicker = E.GetTicker('BTC/USD');
```

### GetRecords

> E.GetRecords(StockType: *String*, Period: [*String*](#records-period), Size: *Any*) => *Record List*

```javascript
// 返回交易所的最新K线数据列表
var thisRecords = E.GetRecords('BTC/USD', 'M5');
```
