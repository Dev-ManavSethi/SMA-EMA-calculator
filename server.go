package main

import (
	"bytes"
	"fmt"

	//standard packages
	"encoding/binary"
	"log"
	"strconv"

	//Local packages
	"./ClientServerComm"
	"./ServerBinanceComm"
	"./env"

	//external packages
	ews "github.com/sacOO7/GoWebsocket"
	//"ews" means "external web socket" (library)
	"golang.org/x/net/websocket"
)

func socket(ws *websocket.Conn) {

	log.Printf("Client connected: %s", getClientIP(ws))

	//Keep a record of number of requests from client:
	var NumberOfRequests int = 0

	//channels are needed for communication between go-routines
	Channel_NoOfRequests := make(chan int, 100)
	Channel_NoOfRequests <- NumberOfRequests

	Channel_ClientToServer := make(chan env.ClientRequest, 1)
	Channel_ServerToClient := make(chan []byte, 1)

	//Channel_BinanceToServer := make(chan *env.ResponseFromBinance, 100000)

	Channel_CurrentSocketInfo := make(chan *ews.Socket, 1)

	Channel_ServerComputation := make(chan *env.ResponseFromBinance, 10000)

	//to instantly get SMA and EMA from anywhere
	SMAinfo := make(chan float64, 1)
	EMAinfo := make(chan float64, 1)

	Channel_RecievedFromBinance := make(chan bool, 100)
	Channel_ClientReqInfo := make(chan env.ClientRequest, 10)

	//these are go-routines that run con-currently.
	//Parameters:
	//ws : The Client web socket connection object.
	// First channel parameter in each function is the channel, from where function reads a value.
	// Second channel parameter in each function is the channel, onto where the function writes a value.
	// From third channel param onwards, are the channels passing info like no. of requests from client, socket info etc.
	go ClientServerComm.ListenFromClient(ws, Channel_ClientToServer, Channel_NoOfRequests)

	//go ServerBinanceComm.ListenFromBinance(Channel_BinanceToServer, Channel_ServerComputation)

	go calculateSMAEMA(Channel_ClientReqInfo, Channel_ServerComputation, Channel_ServerToClient, Channel_RecievedFromBinance, SMAinfo, EMAinfo)

	go ClientServerComm.SendToClient(ws, Channel_ServerToClient)

	ServerBinanceComm.ConnectToBinance(Channel_ClientToServer, Channel_ServerComputation, Channel_NoOfRequests, Channel_CurrentSocketInfo, Channel_RecievedFromBinance, Channel_ClientReqInfo)

	//if SMA and EMA are needed to be sent to client, then un-comment the below line

}

func calculateSMAEMA(Channel_ClientReqInfo chan env.ClientRequest, Channel_ServerComputation chan *env.ResponseFromBinance, Channel_ServerToClient chan []byte, Channel_RecievedFromBinance chan bool, SMAinfo chan float64, EMAinfo chan float64) {

	for {
		select {

		case CombinedResponseFromBinance := <-Channel_ServerComputation:

			RecievedFromClient := <-Channel_ClientReqInfo

			SocketResponse, q := getResponsesFromChannel(CombinedResponseFromBinance)
			//	log.Println("responses seperated")
			//fathering the http response and websocket response from combined-response

			var Sum float64 = 0

			if SocketResponse.Data.KlineData.Closed {

				for i := range q {
					Sum = Sum + q[i]
				}
				//		log.Println("sum found")
				Channel_ServerToClient <- []byte("Kline Closed")

			} else {

				for i := 1; i < len(q); i++ {
					Sum = Sum + q[i]
				}

				var ppprice string

				switch RecievedFromClient.PriceSelection {
				case "o":

					ppprice = SocketResponse.Data.KlineData.OpenPrice
				case "h":

					ppprice = SocketResponse.Data.KlineData.HighPrice
				case "l":

					ppprice = SocketResponse.Data.KlineData.LowPrice
				case "c":

					ppprice = SocketResponse.Data.KlineData.ClosePrice

				}

				PriceSelectionPrice, _ := strconv.ParseFloat(ppprice, 64)

				Sum = Sum + PriceSelectionPrice

				//log.Println("sum found")
				//log.Printf("Sum: %f", Sum)
			}

			//	log.Println("************************************************************************")

			interval, _ := strconv.ParseFloat(RecievedFromClient.NoOfCandles, 64)
			//log.Printf("Interval: %f", interval)
			sma := Sum / interval
			//log.Printf("SMA:  %v", sma)

			WeightingMultiplier := 2 / (interval + 1)
			//log.Printf("m: %f", WeightingMultiplier)

			cp, _ := strconv.ParseFloat(SocketResponse.Data.KlineData.ClosePrice, 64)

			var PrevEMA float64
			select {
			case EMAvalue, ok := <-EMAinfo:
				if ok {
					PrevEMA = EMAvalue
				} else {
					log.Println("Unable to get Previous EMA info")
				}
			default:
				PrevEMA = sma
			}

			ema := (cp-(PrevEMA))*WeightingMultiplier + (PrevEMA)
			//log.Printf("EMA: %v", ema)
			EMAinfo <- ema

			//	log.Println("************************************************************************")

			//in case you want to send data to client
			sma_byte := []byte(strconv.FormatFloat(sma, 'f', 20, 64))
			ema_byte := []byte(strconv.FormatFloat(ema, 'f', 20, 64))

			sma_byte = append(sma_byte, 32)

			for i := range ema_byte {
				sma_byte = append(sma_byte, ema_byte[i])

			}

			//send message to client

			Channel_ServerToClient <- sma_byte

			//log.Println("sum passed")

		}

	}
}

func getClientIP(ws *websocket.Conn) string {
	return ws.RemoteAddr().String()
}

func float64ToByte(f float64) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, f)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

func getResponsesFromChannel(CombinedResponseFromBinance *env.ResponseFromBinance) (*env.SocketResponseFromBinance, []float64) {

	//fathering the http response and websocket response from combined-response
	SocketResponse := *CombinedResponseFromBinance.SocketResponse
	q := CombinedResponseFromBinance.HttpRespArrayQueue

	return &SocketResponse, q

}
