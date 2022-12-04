package ws

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jsawo/colin-server/internal/model"

	"github.com/gorilla/websocket"
)

const (
	pongWait   = 20 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

var (
	clients    = make(map[string]Client)
	messageBus = make(chan Message)
)

type Message struct {
	Topic      string    `json:"topic"`
	Payload    any       `json:"payload"`
	Timestamp  time.Time `json:"timestamp"`
	Recipients []string  `json:"recipients"`
}

func ServeWS(w http.ResponseWriter, r *http.Request) {
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

		fmt.Printf("<< WS client (%v): %v \n", remoteAddr, string(data))
		handleClientCommands(remoteAddr, string(data))
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
		select {
		case msg := <-messageBus:
			for _, clientAddr := range msg.Recipients {
				if _, ok := clients[clientAddr]; !ok {
					model.RemoveClient(clientAddr)
				} else {
					client := clients[clientAddr]
					err := client.Conn.WriteJSON(msg)
					if err != nil {
						fmt.Printf("ERR: closing connection due to websockets write error: %v \n", err)
						closeConnection(client.Addr)
					}
				}
			}
		case <-pingTicker.C:
			for _, client := range clients {
				if err := client.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					closeConnection(client.Addr)
				}
			}
		}
	}
}

func closeConnection(addr string) {
	fmt.Printf("Client left - remaining clients: %d \n", len(clients)-1)
	_ = clients[addr].Conn.Close()
	delete(clients, addr)
}

func SendMessageToSubscribers(msg Message) {
	var recipients []string
	for _, clientAddr := range model.CollectorInstances[msg.Topic].Clients {
		if _, ok := clients[clientAddr]; !ok {
			model.RemoveClient(clientAddr)
		} else {
			recipients = append(recipients, clientAddr)
		}
	}
	msg.Recipients = recipients
	messageBus <- msg
}

func SendMessageToAll(msg Message) {
	var recipients []string
	for _, client := range clients {
		recipients = append(recipients, client.Addr)
	}
	msg.Recipients = recipients
	messageBus <- msg
}

func SendMessageToRecipients(msg Message) {
	messageBus <- msg
}
