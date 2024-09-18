package Socket

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	Dblocal "socket-project/dblocal"
// 	"sync"
// 	"time"

// 	"github.com/gorilla/websocket"
// )

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }
// var (
// 	clients   = make(map[*websocket.Conn]bool)
// 	clientsMu sync.Mutex
// 	broadcast = make(chan Message)
// )

// // var clients = make(map[*websocket.Conn]bool)
// // var broadcast = make(chan Message)

// type Message struct {
// 	MsgType      string `json:"MsgType"`
// 	Created_Time string `json:"Created_Time"`
// 	ToUid        string `json:"ToUid,omitempty"`
// 	ToGroupID    string `json:"ToGroupID,omitempty"`
// 	FromUid      string `json:"FromUid"`
// 	FromUserName string `json:"FromUserName,omitempty"`
// 	Message      string `json:"Message,omitempty"`
// }

// func SockethandleConnections(w http.ResponseWriter, r *http.Request) {
// 	// the Upgrade method refers to the process of upgrading an HTTP connection to a WebSocket connection.
// 	//  This is a key part of establishing a WebSocket connection,
// 	//  which starts as an HTTP request and is then upgraded to a WebSocket protocol if the server supports it.
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	//Ensure to close the conenction once the function is over.
// 	defer conn.Close()

// 	//Making connection true to Enable socket connection
// 	clientsMu.Lock()
// 	clients[conn] = true
// 	clientsMu.Unlock()
// 	//This conction is exicuted as goroutine to handle messages

// 	for {
// 		var msg Message

// 		//from the connection data is stored in the Message struct pointer
// 		err := conn.ReadJSON(&msg)

// 		if err != nil {
// 			fmt.Println(err)
// 			//deleting socket connection from Map of socket connection
// 			clientsMu.Lock()
// 			delete(clients, conn)
// 			clientsMu.Unlock()
// 			return
// 		}
// 		//broadcasting the the Message Struct
// 		broadcast <- msg
// 	}
// }

// func handleMessages() {

// 	for {
// 		//Re-Assigning the Broadcast to Message struct
// 		lMsg := <-broadcast
// 		//Assinging DateTime for each message
// 		lMsg.Created_Time = time.Now().Format("2006-01-02 03:04:05")
// 		saveMsg(lMsg)
// 		fmt.Println(lMsg)
// 		clientsMu.Lock()
// 		for client := range clients {
// 			//from the connection data is stored in the Message struct pointer
// 			err := client.WriteJSON(lMsg)

// 			if err != nil {
// 				fmt.Println(err)
// 				//Close the connection
// 				client.Close()
// 				//deleting socket connection from Map of socket connection
// 				delete(clients, client)
// 			}
// 		}
// 		clientsMu.Unlock()
// 	}
// }
// func StartServer() {
// 	go handleMessages()
// }
import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	Dblocal "socket-project/dblocal"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Customize origin checking as needed for your application
		return true
	},
}

var (
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
	broadcast = make(chan Message)
)

type Message struct {
	MsgType      string `json:"MsgType"`
	Created_Time string `json:"Created_Time"`
	ToUid        string `json:"ToUid,omitempty"`
	ToGroupID    string `json:"ToGroupID,omitempty"`
	FromUid      string `json:"FromUid"`
	FromUserName string `json:"FromUserName,omitempty"`
	Message      string `json:"Message,omitempty"`
}

func SockethandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()

	for {
		var msg Message
		lErr := conn.ReadJSON(&msg)
		if lErr != nil {
			log.Println("ReadJSON error:", lErr)
			clientsMu.Lock()
			delete(clients, conn)
			clientsMu.Unlock()
			break
		}
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		msg.Created_Time = time.Now().Format("2006-01-02 15:04:05")
		// msg.CurrentTime = time.Now().Format("2006-01-02 15:04:05")
		saveMsg(msg)
		log.Println("Broadcasting message:", msg)

		clientsMu.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println("WriteJSON error:", err)
				client.Close()
				delete(clients, client)
			}
		}
		clientsMu.Unlock()
	}
}

func StartServer() {
	go handleMessages()
}

func saveMsg(pMsg Message) {
	log.Println("saveMsg(+)")
	lDb, lErr := Dblocal.LocalDbConnect()
	if lErr != nil {
		//Log Error massage If Database Connection fails

		fmt.Println("SKT-001", lErr.Error())
	} else {
		defer lDb.Close()

		lsqlString := `insert into previous_chat (MsgType, ToUid, ToGroupID, FromUid, FromUserName, Message, Created_Time)values (?, ?, ?, ?, ?, ?, NOW())`

		lExecResult, lErr := lDb.Exec(lsqlString, pMsg.MsgType, pMsg.ToUid, pMsg.ToGroupID, pMsg.FromUid, pMsg.FromUserName, pMsg.Message)

		if lErr != nil {
			// Log an error message if the query execution fails
			log.Println("SKT-002", lErr.Error())
		} else {
			// Check the number of rows affected by the insert query
			rowsAffected, lErr := lExecResult.RowsAffected()
			if lErr != nil {
				// Log an error message if fetching the affected rows count fails
				log.Println("SKT-003", lErr.Error())
			} else {
				// Log the number of rows affected and a success message
				log.Printf("InsertRecords Rows affected: %d\n", rowsAffected)
				log.Println("Record Inserted successfully")
			}
		}

	}
}
