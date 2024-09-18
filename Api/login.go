package api

import (
	"crypto/sha256"
	"database/sql"
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

func Login(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "POST")
	(w).Header().Set("Access-Control-Allow-Headers", "USER,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
	// Log the start of the Login function
	log.Println("Login(+)")
	// Handle POST requests
	if strings.EqualFold(r.Method, http.MethodPost) {
		var lLoginRec Signup
		var lStatusRec Status
		// Read the date from the Body
		lBody, lErr := io.ReadAll(r.Body)
		// Log the Body
		log.Println(string(lBody))
		if lErr != nil {
			log.Println("ALN-001 " + lErr.Error())
			lStatusRec.ErrMsg = "ALN-001 " + lErr.Error()
		} else {
			// Unmarshal response data
			lErr = json.Unmarshal(lBody, &lLoginRec)
			if lErr != nil {
				//Log Error massage If unmarshal fails
				log.Println("ALN-002 " + lErr.Error())
				lStatusRec.ErrMsg = "ALN-002 " + lErr.Error()
			} else {
				lErr = LoginFunction(&lLoginRec)
				// Log the end of the LoginFunction function
				log.Println("LoginFunction(-)")
				if lErr != nil {
					lStatusRec.Status = common.ErrorCode
					log.Println("CreateFunction() Error:" + lErr.Error())
					lStatusRec.ErrMsg = lErr.Error()

				} else {
					// sessionKey, err := GenerateSHA256SessionKey()
					// if err != nil {
					// 	fmt.Println("Error generating session key:", err)
					// 	return
					// }
					// fmt.Println("Generated SHA256 session key:", sessionKey)
					lStatusRec.Status = common.SuccessCode
					//Log the Success Message Of Updation.
					lStatusRec.SuccessMsg = "Login Successs"
					lStatusRec.Uid = lLoginRec.UserId
					lStatusRec.UserName = lLoginRec.Username
					log.Println("Fetching Successfully")
				}
			}
		}
		// Marshal response data
		lData, lErr := json.Marshal(lStatusRec)
		if lErr != nil {
			fmt.Fprintf(w, "Error taking data"+lErr.Error())
		} else {
			// Write response data
			fmt.Fprint(w, string(lData))
		}
		// Log the end of the UpdateAmcMaster function
		log.Println("Login(-)")
	}
}

func LoginFunction(pLoginRec *Signup) error {
	// Log the start of the LoginFunction function
	log.Println("LoginFunction(+)")
	lDb, lErr := Dblocal.LocalDbConnect()
	if lErr != nil {
		//Log Error massage If Database Connection fails
		log.Println("ALN-003", lErr.Error())
		return fmt.Errorf("ALN-003 " + lErr.Error())
	} else {
		// Ensure the database connection is closed when the function exits
		defer lDb.Close()
		//Convert the Password to SHA256 formate
		input := pLoginRec.Password
		hashed := sha256.Sum256([]byte(input))
		hashedString := hex.EncodeToString(hashed[:])
		pLoginRec.Password = hashedString
		log.Println("SHA256 Hash:", hashedString)
		// Prepare the SQL Query to retrieve data from the table
		lquery := "SELECT NVL(Uid, ''), NVL(User_Name, ''), NVL(Pass, '') FROM account_create where Uid=? && Pass=?"
		lErr = lDb.QueryRow(lquery, pLoginRec.UserId, pLoginRec.Password).Scan(
			&pLoginRec.UserId, &pLoginRec.Username, &pLoginRec.Password,
		)
		if lErr == sql.ErrNoRows {
			return fmt.Errorf("No User Found")
		}
		if lErr != nil || pLoginRec.UserId == "" {
			log.Println("No Matching Data" + pLoginRec.UserId + pLoginRec.Password)
			log.Println("ALN-003", lErr.Error())
			return fmt.Errorf("ALN-003 " + lErr.Error())
		}
	}
	return nil
}

// // GenerateRandomBytes generates a random byte slice of a given length
// func GenerateRandomBytes(length int) ([]byte, error) {
// 	bytes := make([]byte, length)
// 	_, err := rand.Read(bytes)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return bytes, nil
// }

// // GenerateSHA256SessionKey generates a random SHA256 session key
// func GenerateSHA256SessionKey() (string, error) {
// 	// Generate 32 random bytes
// 	randomBytes, err := GenerateRandomBytes(32)
// 	if err != nil {
// 		return "", err
// 	}

// 	// Hash the random bytes with SHA256
// 	hash := sha256.New()
// 	if _, err := io.Copy(hash, bytes.NewReader(randomBytes)); err != nil {
// 		return "", err
// 	}
// 	sessionKey := hash.Sum(nil)

// 	// Return the session key as a hex string
// 	return hex.EncodeToString(sessionKey), nil
// }
