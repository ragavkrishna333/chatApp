package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"socket-project/common"
	Dblocal "socket-project/dblocal"
	Socket "socket-project/socket"
	"strings"
)

type ChatFetch struct {
	SocketMsgArr []Socket.Message `json:"SocketMsg"`
	Status
}

func FetchChatFunction(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET")
	(w).Header().Set("Access-Control-Allow-Headers", "USER,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Autharization")
	// Log the start of the FetchChatFunction function
	log.Println("FetchChatFunction(+)")
	// Handle GET requests
	if strings.EqualFold(r.Method, http.MethodGet) {
		var lFetchMsg Socket.Message
		var lStatus ChatFetch
		// var lAddMembersRec AddMembers

		lDb, lErr := Dblocal.LocalDbConnect()
		if lErr != nil {
			log.Println("Error: AFG-001", lErr.Error())
			lStatus.ErrMsg = "Error: AFG-001" + lErr.Error()
			lStatus.Status.Status = common.ErrorCode
		} else { //Executing GET Query in with This Function and Storing the Data In Local Veriable to Pass on.
			lErr = FetchMsgData(lDb, lFetchMsg, &lStatus)
			// Ensure the database connection is closed when the function exits
			defer lDb.Close()
			// Log the end of the execute function
			log.Println("FetchData(-)")
			if lErr != nil {
				lStatus.Status.Status = common.ErrorCode
				lStatus.ErrMsg = "FetchData() Error: AGAM-002" + lErr.Error()
			} else {
				lStatus.Status.Status = common.SuccessCode
				lStatus.SuccessMsg = "Fetching Success"

			}

		}
		// Marshal response data
		lData, lErr := json.Marshal(lStatus)
		if lErr != nil {
			fmt.Fprintf(w, "Error taking data"+lErr.Error())

		} else {
			// Write response data
			fmt.Fprint(w, string(lData))
		}
		// Log the end of the GetAmcMaster function
		log.Println("FetchChatFunction(-)")

	}
}

func FetchMsgData(pDb *sql.DB, pMsg Socket.Message, pStatus *ChatFetch) error {
	// Log the start of the GetDataFunction function
	log.Println("FetchMsgData(+)")

	// Prepare the SQL Query to retrieve data from the table
	lCoreString := `select NVL(MsgType , ''), NVL(ToUid , ''), NVL(ToGroupID , ''), NVL(FromUid , ''), NVL(FromUserName , ''), NVL(Message , ''), NVL(Created_Time , '') from previous_chat`
	pRows, lErr := pDb.Query(lCoreString)
	if lErr != nil {
		// Log an error message if the SQL Query fails to retrieve data from the table
		log.Println("Error: AFG-002", lErr)
		return fmt.Errorf("Error: AFG-002 " + lErr.Error())

	} else {
		// Ensure the result set is closed when the function exits
		defer pRows.Close()
		// Process the result set
		for pRows.Next() {
			// Scan the row and store the result in the local variable
			lErr := pRows.Scan(&pMsg.MsgType, &pMsg.ToUid, &pMsg.ToGroupID, &pMsg.FromUid, &pMsg.FromUserName, &pMsg.Message, &pMsg.Created_Time)
			if lErr != nil {
				// Log an error message if scanning the row fails
				log.Println("Error: AFG-003", lErr)
				return fmt.Errorf("Error: AFG-003 " + lErr.Error())
			} else {
				//Stores the Data fetched from Data base
				pStatus.SocketMsgArr = append(pStatus.SocketMsgArr, pMsg)

			}
		}

	}
	return nil

}
