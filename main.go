// package main

// import (
// 	"fmt"
// 	"net/http"

// 	"github.com/gorilla/websocket"
// )

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// }

// func handler(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer conn.Close()

// 	for {
// 		messageType, p, err := conn.ReadMessage()
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 		message := string(p)
// 		var response string
// 		if message == "hi" {
// 			response = "hello"
// 		} else {
// 			response = message
// 		}
// 		if err := conn.WriteMessage(messageType, []byte(response)); err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 	}
// }

// func main() {
// 	http.HandleFunc("/ws", handler)
// 	fmt.Println("Server started on :29096")
// 	if err := http.ListenAndServe(":29096", nil); err != nil {
// 		fmt.Println("ListenAndServe:", err)
// 	}
// }

package main

import (
	"log"
	"net/http"
	"os"
	api "socket-project/Api"
	Socket "socket-project/socket"
	"time"
)

func main() {
	log.Println("Server started on :29096")

	f, err := os.OpenFile("./log/logfile"+time.Now().Format("02012006.15.04.05.000000000")+".txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	// http.HandleFunc("/", homePage)
	http.HandleFunc("/chatAppws", Socket.SockethandleConnections)
	Socket.StartServer()
	http.HandleFunc("/signUp", api.CreateAccount)
	http.HandleFunc("/login", api.Login)
	http.HandleFunc("/getMembers", api.GetMembersList)
	http.HandleFunc("/fetchGroup", api.FetchGroup)
	http.HandleFunc("/createGroup", api.CreateGroup)
	http.HandleFunc("/deleteGroup", api.DeleteGroup)
	http.HandleFunc("/fetchChat", api.FetchChatFunction)
	// go handleMessages()

	lerr := http.ListenAndServe(":29096", nil)
	if lerr != nil {
		log.Fatalf("Server failed to start: %v", lerr)
	}
}
