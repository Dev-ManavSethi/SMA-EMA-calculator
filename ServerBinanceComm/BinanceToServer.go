package ServerBinanceComm

import (
	"../env"
)

func ListenFromBinance(Channel_BinanceToServer chan *env.ResponseFromBinance, Channel_ServerComputation chan *env.ResponseFromBinance) {

	for {

		select {
		case <-Channel_BinanceToServer:
			Responses := <-Channel_BinanceToServer
			Channel_ServerComputation <- Responses
			//recieve data from *socket

		}

	}
}
