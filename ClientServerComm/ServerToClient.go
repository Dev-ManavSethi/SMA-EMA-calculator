package ClientServerComm

import (
	"golang.org/x/net/websocket"
)

func SendToClient(ws *websocket.Conn, Channel_ServerToClient chan []byte) {
	for {

		select {
		case <-Channel_ServerToClient:
			MessageToClient_byte := <-Channel_ServerToClient
			//log.Println("Sending message to client")
			MessageToClient := string(MessageToClient_byte)
			websocket.Message.Send(ws, MessageToClient)

		}

	}

}
