package api

// import (
// 	"crypto/sha256"
// 	"encoding/hex"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"socket-project/common"
// 	Dblocal "socket-project/dblocal"
// 	"strconv"
// 	"strings"
// )

// type Signup struct {
// 	UserId   string `json:"Uid"`
// 	Username string `json:"User_Name"`
// 	Password string `json:"Pass"`
// }
// type Status struct {
// 	Uid        string `json:"Uid,omitempty"`
// 	UserName   string `json:"User_Name,omitempty"`
// 	SuccessMsg string `json:"SuccessMsg,omitempty"`
// 	ErrMsg     string `json:"errMsg,omitempty"`
// 	Status     string `json:"Status"`
// }

// func CreateAccount(w http.ResponseWriter, r *http.Request) {
// 	// Set CORS headers
// 	(w).Header().Set("Access-Control-Allow-Origin", "*")
// 	(w).Header().Set("Access-Control-Allow-Credentials", "true")
// 	(w).Header().Set("Access-Control-Allow-Methods", "POST")
// 	(w).Header().Set("Access-Control-Allow-Headers", "USER,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
// 	// Log the start of the CreateAccount function
// 	log.Println("CreateAccount(+)")
// 	// Handle POST requests
// 	if strings.EqualFold(r.Method, http.MethodPost) {
// 		var lsignUpRec Signup
// 		var lstatus Status
// 		//Read the date from the Body
// 		lBody, lErr := io.ReadAll(r.Body)
// 		//Log the Body
// 		log.Println(string(lBody))
// 		if lErr != nil {
// 			log.Println("ACG-001 " + lErr.Error())
// 			lstatus.ErrMsg = "ACG-001 " + lErr.Error()
// 		} else {
// 			// Unmarshal response data
// 			lErr = json.Unmarshal(lBody, &lsignUpRec)
// 			if lErr != nil {
// 				//Log Error massage If unmarshal fails
// 				log.Println("ACG-002 " + lErr.Error())
// 				lstatus.ErrMsg = "ACG-002 " + lErr.Error()
// 			} else {
// 				// Executing POST Query in with This Function and Storing the Data In Local Veriable to Pass on.
// 				lErr = CreateFunction(&lsignUpRec)
// 				// Log the end of the CreateFunction function
// 				log.Println("CreateFunction(-)")
// 				if lErr != nil {
// 					lstatus.Status = common.ErrorCode
// 					lstatus.ErrMsg = "CreateFunction() Error:" + lErr.Error()

// 				} else {
// 					lstatus.Status = common.SuccessCode
// 					lstatus.SuccessMsg = "SignUp Successs"
// 					lstatus.Uid = lsignUpRec.UserId
// 					//Log the Success Message Of Updation.
// 					log.Println("Insert Successfully")
// 				}
// 			}

// 		}
// 		// Marshal response data
// 		lData, lErr := json.Marshal(lstatus)
// 		if lErr != nil {
// 			fmt.Fprintf(w, "Error taking data "+lErr.Error())
// 		} else {
// 			// Write response data
// 			fmt.Fprint(w, string(lData))
// 		}
// 		// Log the end of the InsertAmc function
// 		log.Println("CreateAccount(-)")

// 	}

// }

// func CreateFunction(psignUpRec *Signup) error {
// 	// Log the start of the CreateFunction function
// 	log.Println("CreateFunction(+)")
// 	lDb, lErr := Dblocal.LocalDbConnect()
// 	if lErr != nil {
// 		//Log Error massage If Database Connection fails

// 		fmt.Println("ACG-003", lErr.Error())
// 		return fmt.Errorf("ACG-003 " + lErr.Error())
// 	} else {
// 		defer lDb.Close()

// 		var uidData string
// 		lErr = lDb.QueryRow(`SELECT Uid FROM account_create ORDER BY Uid DESC LIMIT ?`, 1).Scan(
// 			&uidData,
// 		)
// 		if lErr != nil && uidData == "" {
// 			// Log an error message if scanning the row fails
// 			log.Println("Error: AGAM-003", lErr)
// 			psignUpRec.UserId = "RA001"
// 			// return fmt.Errorf("Error: AGAM-003 " + lErr.Error())
// 		} else {

// 			log.Println("Fetching success")

// 			num, lErr := strconv.Atoi((strings.Split(uidData, "RA"))[1])

// 			if lErr != nil {
// 				fmt.Println("ACG-006", lErr)
// 				return fmt.Errorf("ACG-006 " + lErr.Error())
// 			} else {
// 				num++
// 				newStr := fmt.Sprintf("RA%03d", num) // Format the number with leading zeros
// 				fmt.Println("New string:", newStr)
// 				psignUpRec.UserId = newStr
// 			}
// 		}
// 		input := psignUpRec.Password
// 		hashed := sha256.Sum256([]byte(input))
// 		hashedString := hex.EncodeToString(hashed[:])
// 		psignUpRec.Password = hashedString
// 		fmt.Println("SHA256 Hash:", hashedString)
// 		fmt.Println(psignUpRec.Password)
// 		lsqlString := `insert into account_create ( Uid, User_Name, Pass, Created_Time)values (?, ?, ?, NOW())`

// 		lExecResult, lErr := lDb.Exec(lsqlString, psignUpRec.UserId, psignUpRec.Username, psignUpRec.Password)
// 		if lErr != nil {
// 			// Log an error message if the query execution fails
// 			log.Println("ACG-007", lErr.Error())
// 			// Return an error with a specific code and the error message

// 			return fmt.Errorf("ACG-007 " + lErr.Error())
// 		} else {
// 			// Check the number of rows affected by the insert query
// 			rowsAffected, lErr := lExecResult.RowsAffected()
// 			if lErr != nil {
// 				// Log an error message if fetching the affected rows count fails
// 				log.Println("ACG-008", lErr.Error())
// 			} else {
// 				// Log the number of rows affected and a success message
// 				log.Printf("InsertRecords Rows affected: %d\n", rowsAffected)
// 				log.Println("Record Inserted successfully")
// 			}
// 		}

// 	}
// 	return nil
// }

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"socket-project/common"
	Dblocal "socket-project/dblocal"
	"strings"
)

type Signup struct {
	UserId   string `json:"Uid"`
	Username string `json:"User_Name"`
	Password string `json:"Pass"`
}
type Status struct {
	Uid        string `json:"Uid,omitempty"`
	UserName   string `json:"User_Name,omitempty"`
	SuccessMsg string `json:"SuccessMsg,omitempty"`
	ErrMsg     string `json:"errMsg,omitempty"`
	Status     string `json:"Status"`
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "POST")
	(w).Header().Set("Access-Control-Allow-Headers", "USER,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
	// Log the start of the CreateAccount function
	log.Println("CreateAccount(+)")
	// Handle POST requests
	if strings.EqualFold(r.Method, http.MethodPost) {
		var lsignUpRec Signup
		var lstatus Status
		//Read the date from the Body
		lBody, lErr := io.ReadAll(r.Body)
		//Log the Body
		log.Println(string(lBody))
		if lErr != nil {
			log.Println("ACG-001 " + lErr.Error())
			lstatus.ErrMsg = "ACG-001 " + lErr.Error()
		} else {
			// Unmarshal response data
			lErr = json.Unmarshal(lBody, &lsignUpRec)
			if lErr != nil {
				//Log Error massage If unmarshal fails
				log.Println("ACG-002 " + lErr.Error())
				lstatus.ErrMsg = "ACG-002 " + lErr.Error()
			} else {
				// Executing POST Query in with This Function and Storing the Data In Local Veriable to Pass on.
				lErr = CreateFunction(&lsignUpRec)
				// Log the end of the CreateFunction function
				log.Println("CreateFunction(-)")
				if lErr != nil {
					lstatus.Status = common.ErrorCode
					lstatus.ErrMsg = "CreateFunction() Error:" + lErr.Error()

				} else {
					lstatus.Status = common.SuccessCode
					lstatus.SuccessMsg = "SignUp Successs"
					lstatus.Uid = lsignUpRec.UserId
					//Log the Success Message Of Updation.
					log.Println("Insert Successfully")
				}
			}

		}
		// Marshal response data
		lData, lErr := json.Marshal(lstatus)
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

func CreateFunction(psignUpRec *Signup) error {
	// Log the start of the CreateFunction function
	log.Println("CreateFunction(+)")
	lDb, lErr := Dblocal.LocalDbConnect()
	if lErr != nil {
		//Log Error massage If Database Connection fails

		fmt.Println("ACG-003", lErr.Error())
		return fmt.Errorf("ACG-003 " + lErr.Error())
	} else {
		defer lDb.Close()

		input := psignUpRec.Password
		hashed := sha256.Sum256([]byte(input))
		hashedString := hex.EncodeToString(hashed[:])
		psignUpRec.Password = hashedString
		fmt.Println("SHA256 Hash:", hashedString)
		fmt.Println(psignUpRec.Password)
		lsqlString := `insert into account_create ( Uid, User_Name, Pass, Created_Time)values (?, ?, ?, NOW())`

		lExecResult, lErr := lDb.Exec(lsqlString, psignUpRec.UserId, psignUpRec.Username, psignUpRec.Password)
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
