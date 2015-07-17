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
		fmt.Println("detected css", title)
		w.Header().Set("Content-Type", "text/css")
	}
	
	p, err := ioutil.ReadFile(title)
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
    if username == "" {
        fmt.Fprintf(w, "%s", "please re-authenticate")
        return
    }
        cookie := http.Cookie{Name:"Id", Value:db.CreateSessionId(username), Expires:time.Now().Add(time.Duration(15)*time.Minute), Path:"/"}
        http.SetCookie(w, &cookie)
    fmt.Println("opening page for " + username)
	title := r.URL.Path[len("/secure/"):]
    body, err := ioutil.ReadFile(title)
	if err != nil {
	fmt.Println(err)
	}
    fmt.Fprintf(w, "%s", body)
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
	messtype, mess, err := connect.ReadMessage()
	if err != nil {
		fmt.Println(err)
	}
	err = connect.WriteMessage(messtype, mess)
	if err != nil {
		fmt.Println(err)
	}
	client := chat.Client{}
	client.Username = user
	client.Conn = connect
	fmt.Println("testing", client.Username)
	chater.Users = append(chater.Users, client)
	
}

func main(){
	chater.SessionName = "test sess"
	http.HandleFunc("/view/", viewHandler)
    http.HandleFunc("/secure/", secureHandler)
	http.HandleFunc("/login/", loginHandler)
	http.HandleFunc("/sock", sockHandler)
	
	if err := http.ListenAndServe(":8090", nil); err != nil {
	fmt.Println(err)
	}
	
}