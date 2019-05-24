package main

import (
	"fmt"
	"log"
	"net/http"

	envv "./env"
	"./globals"

	"golang.org/x/net/websocket"
)

func main() {

	//GlobalVariables := make(map[string]string)
	var interval string
	var NumberOfCandles string
	var PriceSelection string

	AskDefaultValues(&interval, &NumberOfCandles, &PriceSelection)

	SetDefaultEnvironmentValues(interval, NumberOfCandles, PriceSelection)

	HandleSocketConnection()
}

func SetDefaultEnvironmentValues(interval string, NumberOfCandles string, PriceSelection string) {
	envv.Setenv("DefaultInterval", interval, globals.GlobalVariables)
	envv.Setenv("DefaultNumberOfCandles", NumberOfCandles, globals.GlobalVariables)
	envv.Setenv("DefaultPriceSelection", PriceSelection, globals.GlobalVariables)
}

func AskDefaultValues(interval *string, NumberOfCandles *string, PriceSelection *string) {

	log.Println("Please enter default interval in standard format if not provided by client: ")
	fmt.Scanf("%s", interval)

	log.Println("Please enter default No.of candles in integer format if not provided by client: ")
	fmt.Scanf("%s", NumberOfCandles)

	log.Println("Please enter default Price Slection (h , l, o or c) if not provided by client: ")
	fmt.Scanf("%s", PriceSelection)

	log.Printf("Default values Set ---- Interval: %s -------Number Of Candles: %s ---------- Price Selection: %s", *interval, *NumberOfCandles, *PriceSelection)
	log.Println()
}

func HandleSocketConnection() {

	log.Println("Socket server listening for clients on ws://localhost:8000/socket")

	http.Handle("/socket", websocket.Handler(socket))

	http.ListenAndServe(":8000", nil)
}
