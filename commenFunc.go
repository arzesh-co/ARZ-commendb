package CommenDb

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

func ConvertorToInt64(n any) int64 {
	switch n := n.(type) {
	case int:
		return int64(n)
	case int8:
		return int64(n)
	case int16:
		return int64(n)
	case int32:
		return int64(n)
	case int64:
		return int64(n)
	case string:
		num, err := strconv.ParseInt(n, 10, 64)
		if err == nil {
			return int64(0)
		}
		return num
	}
	return int64(0)
}
func today(days any) any {
	year, month, day := time.Now().Date()
	theTime := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	return theTime.Unix()
}
func nDayBefore(days any) any {
	intDay := ConvertorToInt64(days)
	return time.Now().Add((time.Duration(intDay) * -24) * time.Hour).Unix()
}
func nDayAfter(days any) any {
	intDay := ConvertorToInt64(days)
	return time.Now().Add((time.Duration(intDay) * -24) * time.Hour).Unix()
}
func BeginOfThisYear(days any) any {
	year, _, _ := time.Now().Date()
	theTime := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	return theTime.Unix()
}
func BeginOfThisMonth(days any) any {
	year, month, _ := time.Now().Date()
	theTime := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	return theTime.Unix()
}

//func beginOfThisWeek(local int8) int64 {
//	year, month, day := time.Now().Date()
//	weekday := time.Now().Weekday()
//	switch weekday.String() {
//	case "Monday":
//		if local == 1 {
//			day = day - 2
//		} else if local == 3 {
//
//		}
//	}
//	theTime := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
//	return theTime.Unix()
//}

func SetCommenFunc() map[string]func(days any) any {
	funcs := make(map[string]func(days any) any)
	funcs["today"] = today
	funcs["nday_before"] = nDayBefore
	funcs["nday_after"] = nDayAfter
	funcs["begin_of_this_year"] = BeginOfThisYear
	funcs["begin_of_this_month"] = BeginOfThisMonth
	return funcs
}

func FindValue(token, key string) interface{} {
	// Read the RSA public key from file
	publicKey, _ := ioutil.ReadFile(os.Getenv("publicRsaKeyPath"))

	// Parse the RSA public key
	rsaKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		// Return nil if parsing fails
		return nil
	}

	// Parse the JWT token using the RSA public key
	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			// Return an error if the token's signing method is unexpected
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}
		// Return the RSA public key for signature verification
		return rsaKey, nil
	})
	if err != nil {
		// Return nil if token parsing or signature verification fails
		return nil
	}
	// Extract the claims from the parsed token
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		// If the token is invalid or the claims can't be extracted, return nil
		return nil
	}
	// Try to get the value associated with the given key from the token's claims
	value, ok := claims[key]
	if !ok {
		// If the key doesn't exist in the claims, return nil
		return nil
	}
	return value
}

func GetStringValueFromToken(token string, key string) (value string) {
	TokenAccount := FindValue(token, key)
	if TokenAccount == nil {
		return
	}
	value, ok := TokenAccount.(string)
	if !ok {
		return ""
	}
	return
}

func GetUserRolesFromToken(userToken string) []string {
	tokenRoles := FindValue(userToken, "role_keys")
	if tokenRoles == nil {
		return []string{}
	}
	mapRole := make([]map[string]interface{}, len(tokenRoles.([]interface{})))

	jsonData, err := json.Marshal(tokenRoles)
	if err != nil {
		fmt.Println("error to json data", err.Error())
		return []string{}
	}

	err = json.Unmarshal(jsonData, &mapRole)
	if err != nil {
		fmt.Println("error to unmarshal :", err.Error())
		return []string{}
	}

	stringRoles := make([]string, len(mapRole))

	for i, m := range mapRole {
		stringRoles[i] = m["role_key"].(string)
	}

	return stringRoles
}
