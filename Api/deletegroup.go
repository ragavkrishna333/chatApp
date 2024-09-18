package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"socket-project/common"
	Dblocal "socket-project/dblocal"
	"strings"
)

type DeleteGroupRec struct {
	GroupID string `json:"Group_ID"`
	Uid     string `json:"Uid"`
	Admin   string `json:"Admin"`
}

// This API used to delete the Group
func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "USER,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
	// Log the start of the DeleteGroup function
	log.Println("DeleteGroup(+)")
	// Handle POST requests
	if strings.EqualFold(r.Method, http.MethodDelete) {

		var lGroupData MembersList
		var lStatus Status
		var lDeleteRec DeleteGroupRec
		//Read the date from the Body
		lBody, lErr := io.ReadAll(r.Body)
		//Log the Body
		log.Println(string(lBody))
		if lErr != nil {
			log.Println("Error: ", lErr)
			lStatus.Status = common.ErrorCode
			lStatus.ErrMsg = "Error: " + lErr.Error()
		} else {
			// Unmarshal response data
			lErr = json.Unmarshal(lBody, &lDeleteRec)
			if lErr != nil {
				log.Println("Error: ", lErr)
				lStatus.Status = common.ErrorCode
				lStatus.ErrMsg = "Error: " + lErr.Error()
			} else {
				//Connect with Database
				lDb, lErr := Dblocal.LocalDbConnect()
				if lErr != nil {
					log.Println("Error: AFG-001", lErr.Error())
					lStatus.ErrMsg = "Error: AFG-001" + lErr.Error()
					lStatus.Status = common.ErrorCode
				} else {
					//Ensure to close the Database connection after the function ends
					defer lDb.Close()
					//To fetch the data from Database where Group id Matches
					lErr := FetchDataApi(lDb, &lGroupData, &lDeleteRec)

					log.Println("FetchDataApi(-)")

					if lErr != nil {
						lStatus.Status = common.ErrorCode
						lStatus.ErrMsg = "FetchDataApi()Error:" + lErr.Error()
					} else {
						lStatus.Status = common.SuccessCode
						lStatus.SuccessMsg = "Group Deleted Successfully"
						log.Println("Group Deleted Successfully")
					}

				}

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
		// Log the end of the DeleteGroup function
		log.Println("DeleteGroup(-)")
	}
}

//	func filterIPv4Addresses() net.IP {
//		addrs, lErr := net.InterfaceAddrs()
//		if lErr != nil {
//			//Log the Error message If failed to fetch ipv4
//			log.Println("ADG-003 " + lErr.Error())
//			//returns the Error Message Failed to Fetch ipv4
//			// return nil, fmt.Errorf("Error: ADG-003 " + lErr.Error())
//		}
//		var ipv4Addrs net.IP
//		for _, addr := range addrs {
//			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
//				ipv4Addrs = ipNet.IP
//			}
//		}
//		return ipv4Addrs
//	}
func FetchDataApi(pDb *sql.DB, pGroupDataRec *MembersList, pDeleteRec *DeleteGroupRec) error {

	// Log the start of the FetchDataApi function
	log.Println("FetchDataApi(+)")
	// Ensure the database connection is closed when the function exits

	// Prepare the SQL Query to retrieve data from the table
	lSqlQuery := `select Group_Members FROM group_chat_create WHERE Group_ID = ?`
	// lData, lErr := pDb.Exec(lSqlQuery, pDeleteRec.GroupID)
	// Execute the SQL query
	rows, lErr := pDb.Query(lSqlQuery, pDeleteRec.GroupID)
	if lErr != nil {
		fmt.Println("Error executing query:", lErr)
		return fmt.Errorf("Error: AFG-003 " + lErr.Error())
	} else {
		defer rows.Close()
		var GroupMembers string
		// Iterate over the result set
		for rows.Next() {
			//Store the data from database
			if lErr := rows.Scan(&GroupMembers); lErr != nil {
				fmt.Println("Error scanning row:", lErr)
				continue
			}
			// fmt.Println("Group Members:", GroupMembers)
		}

		// Unmarshal the JSON string into the slice
		lErr = json.Unmarshal([]byte(GroupMembers), &pGroupDataRec.GroupMembersArr)
		if lErr != nil {
			fmt.Println("Error unmarshaling JSON:", lErr)
			return fmt.Errorf("Error: AFG-003 " + lErr.Error())
		} else {
			//Verify it is Admin and Uid Matches before deleting
			for index := range pGroupDataRec.GroupMembersArr {
				if pGroupDataRec.GroupMembersArr[index].Admin == pDeleteRec.Admin && pGroupDataRec.GroupMembersArr[index].Uid == pDeleteRec.Uid {
					//Deletefunction is used to Delete the mentioned Group from the table
					lErr := Deletefunction(pDb, pDeleteRec)
					if lErr != nil {
						fmt.Println("Error executing query:", lErr.Error())
						return fmt.Errorf("Error: AFG-003 " + lErr.Error())
					} else {
						return nil
					}
				} else {
					//Log the Error if user data not match
					log.Println("Error: ADG-003 This User is Not Admin")

					return errors.New(" ADG-003: This User is Not Admin")
				}
			}

		}
	}

	return nil

}

// DELETE FROM group_chat_create
// WHERE Group_ID = 'GR001' AND JSON_CONTAINS(json_column, '{"Uid": "RA002"}') AND JSON_CONTAINS_PATH(Group_Members, 'one', '$[*].Admin');
func Deletefunction(pDb *sql.DB, pDeleteRec *DeleteGroupRec) error {
	//Log the start of the Deletefunction function
	log.Println("Deletefunction(+)")
	lSqlString := `DELETE FROM group_chat_create WHERE Group_ID = ?`
	//Execute the query to delete the row
	_, lErr := pDb.Exec(lSqlString, pDeleteRec.GroupID)
	if lErr != nil {
		//Log the Error message if the query fails
		log.Println("Error: ", lErr)
		fmt.Println("Error executing query:", lErr.Error())
		//Log the End of the Deletefunction function
		log.Println("Deletefunction(-)")
		return fmt.Errorf("Error: AFG-003 " + lErr.Error())
	} else {
		//Log the End of the Deletefunction function
		log.Println("Deletefunction(-)")
		return nil
	}
}
