package ClientServerComm

import (
	//standard package
	"log"

	//local package
	"../env"
	"../globals"

	//external package
	"golang.org/x/net/websocket"
)

func ListenFromClient(ws *websocket.Conn, Channel_ClientToServer chan env.ClientRequest, Channel_NoOfRequests chan int) {

	for {
		log.Println("Listening to client...")

		RecievedFromClient := &env.ClientRequest{}

		RecieveRequestFromClient(ws, RecievedFromClient)
		IncrementNumberOfRequests(Channel_NoOfRequests)

		SetEnvironmentVariables(RecievedFromClient, globals.GlobalVariables)

		Channel_ClientToServer <- *RecievedFromClient

	}

}

func SetEnvironmentVariables(RecievedFromClient *env.ClientRequest, globalsGlobalVariables map[string]string) {

	//env.Setenv("ClientSymbol", RecievedFromClient.Symbol, globalsGlobalVariables)

	if RecievedFromClient.Symbol == "" {
		log.Println("Sorry, atleast a SYMBOL is required.")
		return
	}

	if RecievedFromClient.Interval == "" {
		RecievedFromClient.Interval = env.Getenv("DefaultInterval", globalsGlobalVariables)
	}

	if RecievedFromClient.NoOfCandles == "" {
		RecievedFromClient.NoOfCandles = env.Getenv("DefaultNumberOfCandles", globalsGlobalVariables)
	}

	if RecievedFromClient.PriceSelection == "" {
		RecievedFromClient.PriceSelection = env.Getenv("DefaultPriceSelection", globalsGlobalVariables)
	}

	// log.Println("Environment map set: ")
	// for key, value := range globalsGlobalVariables {

	// 	log.Println(key + ": " + value)

	// }
	// log.Println()
}

func RecieveRequestFromClient(ws *websocket.Conn, RecievedFromClient *env.ClientRequest) {

	websocket.JSON.Receive(ws, &RecievedFromClient)

	log.Println("Recieved query from client!")
	log.Println()

}

func IncrementNumberOfRequests(Channel_NoOfRequests chan int) {

	NumberOfRequests := <-Channel_NoOfRequests
	NumberOfRequests++
	Channel_NoOfRequests <- NumberOfRequests

}
