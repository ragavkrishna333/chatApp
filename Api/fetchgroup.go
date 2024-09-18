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

type MembersData struct {
	Uid      string `json:"Uid"`
	UserName string `json:"User_Name"`
	Admin    string `json:"Admin,omitempty"`
	// Status
}
type GroupCreate struct {
	GroupdataArr []MembersList `json:"Group_Data"`
	Status
}
type MembersList struct {
	GroupID         string        `json:"Group_ID`
	GroupName       string        `json:"Group_Name"`
	GroupMembersArr []MembersData `json:"Group_Members"`
}

func FetchGroup(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET")
	(w).Header().Set("Access-Control-Allow-Headers", "USER,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Autharization")
	// Log the start of the GetAmcMaster function
	log.Println("FetchGroup(+)")
	// Handle GET requests
	if strings.EqualFold(r.Method, http.MethodGet) {
		var lFetchGroupRec GroupCreate

		// var lAddMembersRec AddMembers

		lDb, lErr := Dblocal.LocalDbConnect()
		if lErr != nil {
			log.Println("Error: AFG-001", lErr.Error())
			lFetchGroupRec.ErrMsg = "Error: AFG-001" + lErr.Error()
			lFetchGroupRec.Status.Status = common.ErrorCode
		} else { //Executing GET Query in with This Function and Storing the Data In Local Veriable to Pass on.
			lErr = FetchData(lDb, &lFetchGroupRec)
			// Ensure the database connection is closed when the function exits
			defer lDb.Close()
			// Log the end of the execute function
			log.Println("FetchData(-)")
			if lErr != nil {
				lFetchGroupRec.Status.Status = common.ErrorCode
				lFetchGroupRec.ErrMsg = "FetchData() Error: AGAM-002" + lErr.Error()
			} else {
				lFetchGroupRec.Status.Status = common.SuccessCode
				lFetchGroupRec.SuccessMsg = "Fetching Success"
			}

		}
		// Marshal response data
		lData, lErr := json.Marshal(lFetchGroupRec)
		if lErr != nil {
			fmt.Fprintf(w, "Error taking data"+lErr.Error())

		} else {
			// Write response data
			fmt.Fprint(w, string(lData))
		}
		// Log the end of the GetAmcMaster function
		log.Println("FetchGroup(-)")

	}
}

func FetchData(pDb *sql.DB, pFetchGroupRec *GroupCreate) error {
	// Log the start of the GetDataFunction function
	log.Println("FetchData(+)")

	// Prepare the SQL Query to retrieve data from the table
	lCoreString := `select NVL(Group_ID , ''), NVL(Group_Name , ''), NVL(Group_Members , '') from group_chat_create`
	pRows, lErr := pDb.Query(lCoreString)
	if lErr != nil {
		// Log an error message if the SQL Query fails to retrieve data from the table
		log.Println("Error: AFG-002", lErr)
		return fmt.Errorf("Error: AFG-002 " + lErr.Error())

	} else {
		// Ensure the result set is closed when the function exits
		defer pRows.Close()
		var val string
		var lFetchMemberRec []MembersData
		var lMembersList MembersList
		// Process the result set
		for pRows.Next() {

			// Scan the row and store the result in the local variable
			lErr := pRows.Scan(&lMembersList.GroupID, &lMembersList.GroupName, &val)
			// Unmarshal the JSON string into the slice
			err := json.Unmarshal([]byte(val), &lFetchMemberRec)
			if err != nil {
				fmt.Println("Error unmarshaling JSON:", err)
				return fmt.Errorf("Error: AFG-003 " + lErr.Error())
			}
			lMembersList.GroupMembersArr = lFetchMemberRec
			if lErr != nil {
				// Log an error message if scanning the row fails
				log.Println("Error: AFG-003", lErr)
				return fmt.Errorf("Error: AFG-003 " + lErr.Error())
			} else {
				//Stores the Data fetched from Data base
				pFetchGroupRec.GroupdataArr = append(pFetchGroupRec.GroupdataArr, lMembersList)
				log.Println(pFetchGroupRec.GroupdataArr)

			}
		}
	}
	return nil

}
