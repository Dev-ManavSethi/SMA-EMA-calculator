SMA-EMA

Calculates SMA and EMA

To run the app:
1. Install go
2. Install glide
3. Install all dependencies using glide
```
glide install
```
4. Run the app
```
go run *.go
```
5. Enter the default interval, number of candles and price selection to be set, in case, the client doesn't provide.

Note: Format for interval (standard as provided by binance) (Ex: 1m, 3m, 5m, 1h etc...)

Format for number of candles: int

Format for price selection: h (high), l (low), o (open) or c (close)

6. Send query in JSON format as shown below, from Websocket client
```
{   "s" : "<your Symbol in Upper case>",
    "i" : "<interval in standard format Ex: 3m, 5m etc..>",
    "noc" : "<number of candles>",
    "ps" : "<h , l, o or c>"

}
```
(Both the key and value are strings in JSON)

Note: ClientServerComm package contains functions which handle communication between Client and the Server.

Note: BinanceServerComm package contains functions which handle communication between Binance and the Server.

env package contains structures for incoming and outgoing messages and functions to set env variables

global package contains a map which contains the default choices and the Client choices