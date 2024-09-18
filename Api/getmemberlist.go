package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"socket-project/common"
	Dblocal "socket-project/dblocal"
	"strings"
)

type Members struct {
	Uid      string `json:"Uid"`
	UserName string `json:"User_Name"`
}
type MembersStatus struct {
	MembersArr []Members `jsonasd:"Members,omitempty"`
	ErrMsg     string    `json:"errMsg,omitempty"`
	Status     string    `json:"Status"`
}

func GetMembersList(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET")
	(w).Header().Set("Access-Control-Allow-Headers", "USER,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
	// Log the start of the GetMembersList function
	log.Println("GetMembersList(+)")
	// Handle GET requests
	if strings.EqualFold(r.Method, http.MethodGet) {
		var lMembersRec Members
		var lMembersStatusRec MembersStatus
		lMembersStatusRec.Status = common.SuccessCode
		// Connect to the database
		lDb, lErr := Dblocal.LocalDbConnect()
		if lErr != nil {
			// Log an error message if the database connection fails
			log.Println("Error: AGML-001", lErr)
			lMembersStatusRec.ErrMsg = "Error: AGML-001" + lErr.Error()
			lMembersStatusRec.Status = common.ErrorCode

		} else {
			//Executing GET Query in with This Function and Storing the Data In Local Veriable to Pass on.
			lErr = GetMembersFunction(lDb, &lMembersStatusRec, lMembersRec)
			// Log the end of the execute function
			log.Println("GetMembersFunction(-)")
			if lErr != nil {
				lMembersStatusRec.Status = common.ErrorCode
				lMembersStatusRec.ErrMsg = "GetMembersFunction() Error: " + lErr.Error()
			}

		}
		// Marshal response data
		lData, lErr := json.Marshal(lMembersStatusRec)
		if lErr != nil {
			fmt.Fprintf(w, "Error taking data"+lErr.Error())

		} else {
			// Write response data
			fmt.Fprint(w, string(lData))
		}
		// Log the end of the GetMembersList function
		log.Println("GetMembersList(-)")
	}
}

func GetMembersFunction(pDb *sql.DB, pMembersStatusRec *MembersStatus, pMembersRec Members) error {
	// Log the start of the GetDataFunction function
	log.Println("GetMembersFunction(+)")
	// Ensure the database connection is closed when the function exits
	defer pDb.Close()
	// Prepare the SQL Query to retrieve data from the table
	lCoreString := `select NVL(Uid , '') , NVL(User_Name , '')  from account_create`
	pRows, lErr := pDb.Query(lCoreString)
	if lErr != nil {
		// Log an error message if the SQL Query fails to retrieve data from the table
		log.Println("Error: AGML-003", lErr)
		return fmt.Errorf("Error: AGML-003 " + lErr.Error())

	} else {
		// Ensure the result set is closed when the function exits
		defer pRows.Close()
		// Process the result set
		for pRows.Next() {
			// Scan the row and store the result in the local variable
			lErr := pRows.Scan(&pMembersRec.Uid, &pMembersRec.UserName)
			if lErr != nil {
				// Log an error message if scanning the row fails
				log.Println("Error: AGML-004", lErr)
				return fmt.Errorf("Error: AGML-004 " + lErr.Error())
			} else {
				//Stores the Data fetched from Data base
				pMembersStatusRec.MembersArr = append(pMembersStatusRec.MembersArr, pMembersRec)
			}
		}
	}
	return nil
}
