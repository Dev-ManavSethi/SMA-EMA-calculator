package ServerBinanceComm

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	envv "../env"
	ews "github.com/sacOO7/GoWebsocket"
)

func ConnectToBinance(Channel_ClientToServer chan envv.ClientRequest, Channel_ServerComputation chan *envv.ResponseFromBinance, Channel_NoOfRequests chan int, Channel_CurrentSocketInfo chan *ews.Socket, Channel_RecievedFromBinance chan bool, Channel_ClientReqInfo chan envv.ClientRequest) {

	for {

		select {

		case RecievedFromClient := <-Channel_ClientToServer:

			//Send GET request to Binance to get Noc number of data
			resp := SendGETrequestToBinance(RecievedFromClient)

			HTTPresponseArray := makeResponseInArrayFormat(resp)

			q := AddAllPricesInQueue(RecievedFromClient, HTTPresponseArray)

			//Subscribing to WebSocket Connection of Binance

			socket := SubscribeToWebSocket(RecievedFromClient, Channel_NoOfRequests, Channel_CurrentSocketInfo)

			//recieving data from websocket
			RecieveDataAndPassOn(socket, q, RecievedFromClient, Channel_RecievedFromBinance, Channel_ClientReqInfo, Channel_ServerComputation)

		}

	}

}

func RecieveDataAndPassOn(socket *ews.Socket, q []float64, RecievedFromClient envv.ClientRequest, Channel_RecievedFromBinance chan bool, Channel_ClientReqInfo chan envv.ClientRequest, Channel_ServerComputation chan *envv.ResponseFromBinance) {
	//creating instance of formatted response from binance through websocket
	RecievedFromBinanceSocket_JSON := &envv.SocketResponseFromBinance{

		Data: &envv.Data{
			KlineData: &envv.KlineData{},
		},
	}

	//Combined-response contains both HTTP response data (array of relavant prices) and Websocket data
	CombinedResponseFromBinance := &envv.ResponseFromBinance{
		HttpRespArrayQueue: q,
	}

	var ppprice string

	go func() {

		for {

			socket.OnPingReceived = func(data string, socket ews.Socket) {
				socket.SendText("PONG")
			}

			socket.OnTextMessage = func(RecievedFromBinance string, socket ews.Socket) {

				Channel_ClientReqInfo <- RecievedFromClient
				//log.Printf("Received from Binance!")

				json.Unmarshal([]byte(RecievedFromBinance), &RecievedFromBinanceSocket_JSON)

				//preparing which vprice is relavant to us, H, L, C or O price.
				switch RecievedFromClient.PriceSelection {
				case "o":

					ppprice = RecievedFromBinanceSocket_JSON.Data.KlineData.OpenPrice
				case "h":

					ppprice = RecievedFromBinanceSocket_JSON.Data.KlineData.HighPrice
				case "l":

					ppprice = RecievedFromBinanceSocket_JSON.Data.KlineData.LowPrice
				case "c":

					ppprice = RecievedFromBinanceSocket_JSON.Data.KlineData.ClosePrice

				default:

					ppprice = RecievedFromBinanceSocket_JSON.Data.KlineData.HighPrice

				}

				//just string -> float64, nothing else!
				n, _ := strconv.ParseFloat(ppprice, 64)

				if RecievedFromBinanceSocket_JSON.Data.KlineData.Closed == true {
					log.Println()
					//log.Printf("-------------------------> Kline closed! (%s)<-------------------------------", envv.Getenv("ClientSymbol", globals.GlobalVariables))
					log.Println()
					q = q[1:]
					q = append(q, n) //ADD CURRENT PRICE SELECTION

					//add updated array to combined-response variable
					CombinedResponseFromBinance.HttpRespArrayQueue = q

					//add websocket response to combined-response variable
					CombinedResponseFromBinance.SocketResponse = RecievedFromBinanceSocket_JSON

					//just ensuring everything con-current starts, only after we have recieved complete data from Binance
					Channel_RecievedFromBinance <- true

					//sending combined response to the channel
					Channel_ServerComputation <- CombinedResponseFromBinance
					//log.Println("passed")
				} else {
					//if x!=true
					//then dont touch the array
					CombinedResponseFromBinance.SocketResponse = RecievedFromBinanceSocket_JSON

					//you must know this!
					Channel_RecievedFromBinance <- true

					Channel_ServerComputation <- CombinedResponseFromBinance
					//	log.Println("passed")
				}

			}
		}
	}()
}

func SendGETrequestToBinance(RecievedFromClient envv.ClientRequest) *http.Response {

	log.Println()

	log.Println("Sending GET request to Binance!")

	httpURL := "https://api.binance.com/api/v1/klines?symbol=" + RecievedFromClient.Symbol + "&interval=" + RecievedFromClient.Interval + "&limit=" + RecievedFromClient.NoOfCandles

	log.Println("GET req url: " + httpURL)

	req, err := http.NewRequest("GET", httpURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// if resp.StatusCode != http.StatusOK {
	// 	log.Fatalf("Status code is not OK: %v (%s)", resp.StatusCode, resp.Status)
	// }

	log.Println("Binance GET request successfull")

	log.Println()

	return resp
}

func makeResponseInArrayFormat(resp *http.Response) [][]string {

	var HTTPresponseArray [][]string
	//reading body of the response
	Body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	//converting body (un-formatted, byte) to formatted, 2-D array form where each row represents individual response
	json.Unmarshal(Body, &HTTPresponseArray)
	return HTTPresponseArray
}

func AddAllPricesInQueue(ClientRequest envv.ClientRequest, HTTPresponseArray [][]string) []float64 {

	//converting string to int
	noc, _ := strconv.ParseInt(ClientRequest.NoOfCandles, 64, 64)

	//making an array to store relevant prices to calculate SMA
	q := make([]float64, noc)

	//ind is used for price selection. The user inputs price selection as o, h, l or c. ind is used to extract respective prices from 2D array
	var ind int

	switch ClientRequest.PriceSelection {
	case "o":
		ind = 1

	case "h":
		ind = 2
	case "l":
		ind = 3
	case "c":
		ind = 4

	default:
		ind = 2

	}

	for i := range HTTPresponseArray {
		//converting string to float64

		a, _ := strconv.ParseFloat(HTTPresponseArray[i][ind], 64)

		//adding the price to queue
		q = append(q, a)
	}
	return q

}

func SubscribeToWebSocket(RecievedFromClient envv.ClientRequest, Channel_NoOfRequests chan int, Channel_CurrentSocketInfo chan *ews.Socket) *ews.Socket {
	StreamName := strings.ToLower(RecievedFromClient.Symbol) + "@kline_" + RecievedFromClient.Interval

	url := "wss://stream.binance.com:9443/stream?streams=" + StreamName

	log.Println("Connecting to socket for current candle: " + url)

	socket := ews.New(url)

	socket.OnConnected = func(socket ews.Socket) {
		log.Println("Connected to Binance socket!")
		log.Println()
	}

	//Ensuring previous socket is closed before a new socket is started
	NumberOfRequests := <-Channel_NoOfRequests
	if NumberOfRequests > 1 {

		log.Println("New request recieved. Closing previous socket Connection")

		currentSocket := <-Channel_CurrentSocketInfo
		currentSocket.Close()

		socket.Connect()

		//sending current socket info to channel
		Channel_CurrentSocketInfo <- &socket

	} else {

		socket.Connect()
		Channel_CurrentSocketInfo <- &socket
	}

	return &socket
}
