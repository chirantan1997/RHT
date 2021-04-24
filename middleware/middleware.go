package middleware

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("foxtrot")
var adminKey = []byte("nineleaps")

// IsAuthorized ...
func IsAuthorized(endpoint func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	//var res models.ResponseResult

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println(r.Header["Token"])
		fmt.Println(r.Header["Admin"])
		var token *jwt.Token
		var err error

		admin := r.Header["Admin"]
		element := admin[len(admin)-1]

		if r.Header["Token"] != nil {
			if element == "true" {

				token, err = jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return adminKey, nil
				})
			} else {
				token, err = jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return mySigningKey, nil
				})
			}

			if err != nil {
				fmt.Fprintf(w, err.Error())
				fmt.Println("trello")
				fmt.Println(err.Error())

			}

			if token.Valid {
				//fmt.Fprintf(w, "Authorized")
				endpoint(w, r)
				fmt.Println("nineleaps")
			}
		} else {
			fmt.Println("hiko")
			fmt.Fprintf(w, "Not Authorized")
		}
	})

}
