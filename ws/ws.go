package ws

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	pongWait   = 20 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

var (
	clients       = make(map[string]Client)
	messageBus    = make(chan Message)
	messageWriter = defaultMessageWriter
)

type Message struct {
	Type      string    `json:"type"`
	Payload   any       `json:"payload"`
	Timestamp time.Time `json:"timestamp"`
}

func ServeWS(w http.ResponseWriter, r *http.Request) {
	fmt.Println("client connected to websocket")

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		var handShakeError *websocket.HandshakeError
		if errors.As(err, &handShakeError) {
			fmt.Printf("ERR: error during websocket handshake: %v \n", err)
		} else {
			fmt.Printf("ERR: websocket error: %v \n", err)
		}

		return
	}

	clients[r.RemoteAddr] = Client{
		Addr: r.RemoteAddr,
		Conn: ws,
	}

	readFromWs(ws, r.RemoteAddr)
}

func readFromWs(ws *websocket.Conn, remoteAddr string) {
	defer ws.Close()
	ws.SetReadLimit(512)

	err := ws.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		fmt.Printf("ERR: websockets error, failed to set deadline: %v \n", err)
	}

	ws.SetPongHandler(func(string) error {
		err = ws.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			fmt.Printf("ERR: websockets error, failed to set handler: %v \n", err)
		}
		return nil
	})

	for {
		_, data, err := ws.ReadMessage()
		if err != nil {
			fmt.Printf("ERR: websockets error, failed to read: %v \n", err)
			break
		}

		fmt.Printf("Received data from WS client (%v): %v \n", remoteAddr, string(data))
	}
}

func WriteToWS() {
	pingTicker := time.NewTicker(pingPeriod)

	defer func() {
		pingTicker.Stop()

		for _, client := range clients {
			closeConnection(client.Addr)
		}
	}()

	for {
		var err error
		select {
		case msg := <-messageBus:
			for _, client := range clients {
				err = client.Conn.WriteJSON(msg)

				if err != nil {
					closeConnection(client.Addr)
				}
			}
		case <-pingTicker.C:
			for _, client := range clients {
				if err = client.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					closeConnection(client.Addr)
				}
			}
		}
	}
}

func closeConnection(addr string) {
	fmt.Printf("WS closing client connection. Remaining clients: %d \n", len(clients)-1)
	_ = clients[addr].Conn.Close()
	delete(clients, addr)
}

func WriteMessage(msgType string, payload any) {
	writeMessage(Message{msgType, payload, time.Now()})
}

func writeMessage(msg Message) {
	fmt.Printf("WS message: %q \n", msg)
	messageWriter(msg)
}

func SetMessageWriter(writer func(msg Message)) {
	messageWriter = writer
}

func defaultMessageWriter(msg Message) {
	if len(clients) > 0 {
		messageBus <- msg
	}
}
