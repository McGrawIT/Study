package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type User struct {
	Id           uint32 `json:"id"`
	Username     string `json:"username"`
	MoneyBalance uint32 `json:"balance"`
}

type UserParams struct {
	Username     string `json:"username"`
	MoneyBalance uint32 `json:"balance"`
}

var userIdCounter uint32 = 0

var userStore = []User{}

func createUserHandler(w http.ResponseWriter, r *http.Request) {

	p := UserParams{}
//	p := []UserParams{}

//	fmt.Println( "Creating User")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("ReadAll() Error: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
//	fmt.Println( "Read User Input")
//	fmt.Println( "Input:", string ( body ) )
//	fmt.Println( "Unmarshalling")

	err = json.Unmarshal(body, &p)
	if err != nil {
		fmt.Printf("Unmarshal() Error: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
//	fmt.Println( "Checking for Unique User(s)")
//	fmt.Println( "User List:", p )

//	for _, u := range p {

//		fmt.Println( "User:", p )
		err = validateUniqueness( p.Username )
		if err != nil {

			fmt.Printf("Error: %s\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		u1 := User{
			Id:           userIdCounter,
			Username:     p.Username,
			MoneyBalance: p.MoneyBalance,
		}

//		fmt.Println( "Adding Record:", u1 )
		userStore = append(userStore, u1 )

		userIdCounter += 1

//	}

//	fmt.Println( "All Users:", userStore )
	w.WriteHeader(http.StatusCreated)
}

func validateUniqueness(username string) error {

	for _, u := range userStore {

		if u.Username == username {
			return errors.New("Username is already used")
		}
	}

//	fmt.Println( "User:", username, "is New" )

	return nil
}

func listUsersHandler(w http.ResponseWriter, r *http.Request) {

	users, err := json.Marshal(userStore)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(users)
}

func Handlers() *mux.Router {

	r := mux.NewRouter()

	r.HandleFunc("/users", createUserHandler).Methods("POST")

	r.HandleFunc("/users", listUsersHandler).Methods("GET")

	return r
}
