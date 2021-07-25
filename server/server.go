package server

import (
	"encoding/json"
	"fmt"
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

func sendFile(conn *websocket.Conn, path string) {
	base := strings.Split(path, "/")
	basePath := strings.Join(base[:len(base)-1], "/") + "/"
	bf, err := os.Stat(path)
	if err != nil {
		log.Printf("1===%v", err)
	}
	if bf.IsDir() {
		needDir := strings.ReplaceAll(basePath+bf.Name(), Folder, "")
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("dir %s", needDir)))
		for {
			_, message, _ := conn.ReadMessage()
			if string(message) == "dir "+needDir+" ok" {
				break
			}
		}
		fs, err := ioutil.ReadDir(path)
		if err != nil {
			log.Printf("2===%v", err)
		}
		for _, f := range fs {
			sendFile(conn, path+"/"+f.Name())
		}
	} else {
		file, err := os.Open(path)
		if err != nil {
			log.Printf("3===%v", err)
		}
		defer file.Close()
		needFile := strings.ReplaceAll(file.Name(), Folder, "")
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("file %s", needFile)))
		for {
			_, message, _ := conn.ReadMessage()
			if string(message) == "file "+needFile+" ok" {
				break
			}
		}
		// read and write
		// buf := make([]byte, 4096)
		// for {
		// 	n, err := file.Read(buf)
		// 	if err == io.EOF {
		// 		fmt.Println("read finished")
		// 		return
		// 	}
		// 	if err != nil {
		// 		fmt.Println("read err:", err)
		// 		return
		// 	}
		// 	err = conn.WriteMessage(websocket.BinaryMessage, buf[:n])
		// 	if err != nil {
		// 		fmt.Println("conn.Write err:", err)
		// 		return
		// 	}
		// }
		buf, _ := ioutil.ReadAll(file)
		err = conn.WriteMessage(websocket.BinaryMessage, buf)
		if err != nil {
			fmt.Println("conn.Write err:", err)
			return
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

func Start() {
	http.HandleFunc("/", handleTransfer)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
