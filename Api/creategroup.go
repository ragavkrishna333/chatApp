package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"socket-project/common"
	Dblocal "socket-project/dblocal"
	"strconv"
	"strings"
)

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "POST")
	(w).Header().Set("Access-Control-Allow-Headers", "USER,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
	// Log the start of the Login function
	log.Println("CreateGroup(+)")
	// Handle POST requests
	if strings.EqualFold(r.Method, http.MethodPost) {
		var lCreateGroupRec MembersList
		var lStatus Status
		//Read the date from the Body
		lBody, lErr := io.ReadAll(r.Body)
		//Log the Body
		log.Println(string(lBody))
		if lErr != nil {
			log.Println("ACG-001 " + lErr.Error())
			lStatus.ErrMsg = "ACG-001 " + lErr.Error()
		} else {
			// Unmarshal response data
			lErr = json.Unmarshal(lBody, &lCreateGroupRec)
			if lErr != nil {
				//Log Error massage If unmarshal fails
				log.Println("ACG-002 " + lErr.Error())
				lStatus.ErrMsg = "ACG-002 " + lErr.Error()
			} else {
				// Executing POST Query in with This Function and Storing the Data In Local Veriable to Pass on.
				lErr = CreateGroupFunction(lCreateGroupRec)
				// Log the end of the CreateFunction function
				log.Println("CreateFunction(-)")
				if lErr != nil {
					lStatus.Status = common.ErrorCode
					lStatus.ErrMsg = "CreateFunction() Error:" + lErr.Error()

				} else {
					lStatus.Status = common.SuccessCode
					lStatus.SuccessMsg = "Group Created"

					//Log the Success Message Of Updation.
					log.Println("Insert Successfully")
				}
			}

		}
		// Marshal response data
		lData, lErr := json.Marshal(lStatus)
		if lErr != nil {
			fmt.Fprintf(w, "Error taking data "+lErr.Error())
		} else {
			// Write response data
			fmt.Fprint(w, string(lData))
		}
		// Log the end of the InsertAmc function
		log.Println("CreateAccount(-)")

	}
}

func CreateGroupFunction(pCreateGroupRec MembersList) error {
	// Log the start of the CreateFunction function
	log.Println("CreateFunction(+)")
	lDb, lErr := Dblocal.LocalDbConnect()
	if lErr != nil {
		//Log Error massage If Database Connection fails

		fmt.Println("ACG-003", lErr.Error())
		return fmt.Errorf("ACG-003 " + lErr.Error())
	} else {
		defer lDb.Close()

		var lGridData string
		lErr = lDb.QueryRow(`SELECT Group_ID FROM group_chat_create ORDER BY Group_ID DESC LIMIT ?`, 1).Scan(
			&lGridData,
		)
		if lErr != nil && lGridData == "" {
			// Log an error message if scanning the row fails
			log.Println("Error: AGAM-003", lErr)
			pCreateGroupRec.GroupID = "GR001"
			// return fmt.Errorf("Error: AGAM-003 " + lErr.Error())
		} else {

			log.Println("Fetching success")

			num, lErr := strconv.Atoi((strings.Split(lGridData, "GR"))[1])

			if lErr != nil {
				fmt.Println("ACG-006", lErr)
				return fmt.Errorf("ACG-006 " + lErr.Error())
			} else {
				num++
				newStr := fmt.Sprintf("GR%03d", num) // Format the number with leading zeros
				fmt.Println("New string:", newStr)
				pCreateGroupRec.GroupID = newStr
			}
		}
		lGroupMembersData, err := json.Marshal(pCreateGroupRec.GroupMembersArr)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return fmt.Errorf("Error: AFG-003 " + lErr.Error())
		}
		lsqlString := `insert into group_chat_create (Group_ID, Group_Name, Group_Members, Created_Time)values (?, ?, ?, NOW())`

		lExecResult, lErr := lDb.Exec(lsqlString, pCreateGroupRec.GroupID, pCreateGroupRec.GroupName, lGroupMembersData)
		fmt.Println(pCreateGroupRec)
		if lErr != nil {
			// Log an error message if the query execution fails
			log.Println("ACG-007", lErr.Error())
			// Return an error with a specific code and the error message

			return fmt.Errorf("ACG-007 " + lErr.Error())
		} else {
			// Check the number of rows affected by the insert query
			rowsAffected, lErr := lExecResult.RowsAffected()
			if lErr != nil {
				// Log an error message if fetching the affected rows count fails
				log.Println("ACG-008", lErr.Error())
			} else {
				// Log the number of rows affected and a success message
				log.Printf("InsertRecords Rows affected: %d\n", rowsAffected)
				log.Println("Record Inserted successfully")
			}
		}

	}
	return nil
}
