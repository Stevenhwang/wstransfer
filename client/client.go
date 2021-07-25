package client

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

var nf string

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
				if strings.HasPrefix(string(message), "dir") {
					needDir := strings.ReplaceAll(string(message), "dir ", "")
					os.Mkdir(needDir, 0755)
					c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("dir %s ok", needDir)))
				}
				if strings.HasPrefix(string(message), "file") {
					needFile := strings.ReplaceAll(string(message), "file ", "")
					nf = needFile
					c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("file %s ok", needFile)))
				}
			}
			if mt == websocket.BinaryMessage {
				f, _ := os.Create(nf)
				f.Write(message)
				defer f.Close()
			}
			log.Println(string(message))
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		c.WriteMessage(websocket.TextMessage, []byte(scanner.Text()))
	}
}
