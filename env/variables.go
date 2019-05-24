package env

//var Env_variables map[string]string

type ClientRequest struct {
	Symbol         string `json:"s"`
	Interval       string `json:"i"`
	NoOfCandles    string `json:"noc"`
	PriceSelection string `json:"ps"`
}

type KlineData struct {
	StartTime         int64  `json:"t"` // Kline start time
	CloseTime         int64  `json:"T"` // Kline close time
	Symbol            string `json:"s"`
	Interval          string `json:"i"` // Interval
	FirstTradeID      int64  `json:"f"` // First trade ID
	LastTradeID       int64  `json:"L"` // Last trade ID
	OpenPrice         string `json:"o"`
	ClosePrice        string `json:"c"` // Close price
	HighPrice         string `json:"h"`
	LowPrice          string `json:"l"` // Low price
	BaseAssetVolume   string `json:"v"` // Base asset volume
	NumberOfTrades    int64  `json:"n"` // Number of trades
	Closed            bool   `json:"x"` // Is this kline closed?
	QuoteAssetVolume  string `json:"q"` // Quote asset volume
	TakerBuyBaseAsset string `json:"V"` // Taker buy base asset volume
	TakerBuyQuote     string `json:"Q"` // Taker buy quote asset volume
	Ignore            string `json:"B"`
}

type Data struct {
	Kline     string     `json:"e"`
	EventTime int64      `json:"E"`
	Symbol    string     `json:"s"`
	KlineData *KlineData `json:"k"`
}

type SocketResponseFromBinance struct {
	Stream string `json:"stream"`
	Data   *Data  `json:"data"`
}

type ResponseFromBinance struct {
	HttpRespArrayQueue []float64
	SocketResponse     *SocketResponseFromBinance
}

type HTTPresponseFromBinance struct {
	OpenTime         int64
	OpenPrice        string
	HighPrice        string
	LowPrice         string
	ClosePrice       string
	Volume           string
	CloseTime        int64
	QuoteAssetVolume string
	NumberOfTrades   int64
	TBBAV            string
	TBQAV            string
	Ignore           string
}

type Empty struct{}

func Setenv(key string, value string, e map[string]string) {
	e[key] = value

}

func Getenv(key string, e map[string]string) string {

	return e[key]

}
