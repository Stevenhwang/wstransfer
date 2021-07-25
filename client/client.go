package client

import (
	"bufio"
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

func Start() {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial error:", err)
	}
	defer c.Close()
	go func() {
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}
			if mt == websocket.TextMessage {
				log.Println(string(message))
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		c.WriteMessage(websocket.TextMessage, []byte(scanner.Text()))
	}
}
