package main

import (
    "github.com/titmuscody/decision/db"
	"github.com/titmuscody/decision/chat"
	"github.com/gorilla/websocket"
	"fmt"
	"time"
	"net/http" 
	"io/ioutil"
    "strings"
    )

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

var chater chat.ChatRoom

func viewHandler(w http.ResponseWriter, r *http.Request){
	
	
	title := r.URL.Path[len("/view/"):]
	if title[len(title)-4:] == ".css" {
		//fmt.Println("detected css", title)
		w.Header().Set("Content-Type", "text/css")
	}
	
	p, err := ioutil.ReadFile("public/" + title)
	if err != nil {
		fmt.Println(err, title)
		return
	}
	fmt.Fprintf(w, "%s", p)
	//for _, c := range r.Cookies(){
	//	fmt.Println(c.Name)
	//}
	//fmt.Println(w, r.Cookies())	
}

func loginHandler(w http.ResponseWriter, r *http.Request){
    auth := r.Header["Authorization"][0]
    if strings.Contains(auth, ":") {
        userPass := strings.Split(auth, ":")
        username := userPass[0]
        pass := userPass[1]
        userHash := db.GetUserHash(username)
		fmt.Println("checking diff submit, local", pass)
		fmt.Println(userHash)
        if userHash == pass {
            fmt.Println("user authenticated")
            expires := time.Now().Add(24 * time.Hour)
            cookie := http.Cookie{Name:"Id", Value:db.CreateSessionId(username), Expires:expires, Path:"/"}
            http.SetCookie(w, &cookie)
        } else {
            fmt.Fprintf(w, "%s", "no log in for you")

        }
    } else if auth != "" {
        key := db.GetUserKey(auth)
        fmt.Fprintf(w, "%s", key)
    } else {
        fmt.Println("in else")
        fmt.Fprintf(w, "%s", "unable to determine intentions")
    }
}

func secureHandler(w http.ResponseWriter, r *http.Request){
    id, err := r.Cookie("Id")
	if err != nil {
	fmt.Println(err)
	fmt.Fprintf(w, "%s", "sorry you must be logged in to access this page")
	}
    username := db.Validate(id.Value)
	fmt.Println("id username", id.Value, username)
    if username == "" {
        fmt.Fprintf(w, "%s", "please re-authenticate")
        return
    }
    cookie := http.Cookie{Name:"Id", Value:db.CreateSessionId(username), Expires:time.Now().Add(time.Duration(15)*time.Minute), Path:"/"}
    http.SetCookie(w, &cookie)
    fmt.Println("opening secure page for " + username)
	title := r.URL.Path[len("/secure/"):]
    body, err := ioutil.ReadFile("secure/" + title)
	if err != nil {
	fmt.Println(err)
	}
    fmt.Fprintf(w, "%s", body)
}

func createUser(w http.ResponseWriter, r *http.Request){
	id, err := r.Cookie("Id")
	if err != nil {
	fmt.Println(err)
	}
	if db.GetUser(id.Value) == "tisourit" {
		
		r.ParseForm()
		fmt.Printf("%+v\n", r.Form)
		username := r.FormValue("username")
		password := r.FormValue("password")
		if db.CreateUser(username, password) {
			fmt.Fprintf(w, "%s", "User created")
		}
	} else {
		fmt.Fprintf(w, "%s", "you are not admin")
	}
}

func sockHandler(w http.ResponseWriter, r *http.Request) {
	
	userId, _ := r.Cookie("Id")
	fmt.Println(userId)
	user := db.GetUser(string(userId.Value))
	connect, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	return
	}
	//messtype, mess, err := connect.ReadMessage()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println("message type is ", messtype)
	
	client := chat.Client{Username:user, Conn:connect}
	//client.Username = user
	//client.Conn = connect
	go chater.StartClient(&client)
	
	//users := ""
	//fmt.Println("numberr=", len(chater.Users))
	
	//for i, _ := range chater.Users {
	//	users += chater.Users[i].Username + " "
	//	fmt.Println("in loop", chater.Users[i].Username)
	//}
	//err = connect.WriteMessage(1, []byte(users))
	//if err != nil {
	//	fmt.Println(err)
	//}
}

func main(){
	chater.SessionName = "test sess"
	go chater.Start()
	http.HandleFunc("/view/", viewHandler)
    http.HandleFunc("/secure/", secureHandler)
	http.HandleFunc("/login/", loginHandler)
	http.HandleFunc("/sock", sockHandler)
	http.HandleFunc("/user", createUser)
	
	if err := http.ListenAndServe(":8090", nil); err != nil {
	fmt.Println(err)
	}
	
}