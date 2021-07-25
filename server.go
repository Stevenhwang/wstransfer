package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

var Folder = "C:/Users/90hua/"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func sendFile(conn *websocket.Conn, filepath string) {
	bf, err := os.Stat(filepath)
	if err != nil {
		log.Printf("1===%v", err)
	}
	if bf.IsDir() {
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("dir %s", bf.Name())))
		fs, err := ioutil.ReadDir(filepath)
		if err != nil {
			log.Printf("2===%v", err)
		}
		for _, f := range fs {
			sendFile(conn, filepath+"/"+f.Name())
		}
	} else {
		file, err := os.Open(filepath)
		if err != nil {
			log.Printf("3===%v", err)
		}
		defer file.Close()
		buf := make([]byte, 4096)
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("file %s", file.Name())))
		// read and write
		for {
			n, err := file.Read(buf)
			if err == io.EOF {
				fmt.Println("read finished")
				return
			}
			if err != nil {
				fmt.Println("read err:", err)
				return
			}
			err = conn.WriteMessage(websocket.BinaryMessage, buf[:n])
			if err != nil {
				fmt.Println("conn.Write err:", err)
				return
			}
		}
	}
}

func handleTransfer(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	files, _ := ioutil.ReadDir(Folder)
	filesMap := map[string]string{}
	for _, f := range files {
		if f.IsDir() {
			filesMap[f.Name()] = "Dir"
		} else {
			filesMap[f.Name()] = "File"
		}
	}
	prettyMap, _ := json.MarshalIndent(filesMap, "", "  ")
	c.WriteMessage(websocket.TextMessage, prettyMap)
	c.WriteMessage(websocket.TextMessage, []byte("usage: get {filename}|{dirname}"))

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		if !strings.HasPrefix(string(message), "get ") {
			c.WriteMessage(websocket.TextMessage, []byte("usage: get {filename}|{dirname}"))
		} else {
			transName := strings.Replace(string(message), "get ", "", -1)
			if _, ok := filesMap[transName]; !ok {
				c.WriteMessage(mt, []byte("no such file"))
			} else {
				// begin transfer
				sendFile(c, Folder+transName)
			}
		}
	}
}

func main() {
	http.HandleFunc("/", handleTransfer)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
