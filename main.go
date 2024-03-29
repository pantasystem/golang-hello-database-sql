package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DbConnection *sql.DB
var sc = bufio.NewScanner(os.Stdin)
var messageRepository MessageRepository

func init() {
	var err error
	DbConnection, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", "test", "secret", "127.0.0.1:3308", "app-database"))
	if err != nil {
		fmt.Println("open error", err)

	}
	_, err = DbConnection.Query("CREATE TABLE IF NOT EXISTS message(id INT NOT NULL PRIMARY KEY AUTO_INCREMENT, text VARCHAR(255) NOT NULL)")
	if err != nil {
		fmt.Println("create table error:", err)
	}
	messageRepository = CreateMessageRepository()

}

func handleGreet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello")
}

func handleSearchMessages(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//query := r.FormValue("query")

	result, e := DbConnection.Query(`SELECT id, text FROM message WHERE text LIKE ?`, "%hoge%")
	if e != nil {
		fmt.Println("エラー:", e)
	}

	var messages []Message
	for result.Next() {
		var msg Message
		result.Scan(&msg.Id, &msg.Text)
		messages = append(messages, msg)
	}

	w.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal(messages)
	w.Write(res)
}

func handleMessages(w http.ResponseWriter, r *http.Request) {

	messages, _ := messageRepository.FindAll()
	w.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal(messages)
	w.Write(res)

}

func handleCreateMessage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	text := r.FormValue("text")
	if len(text) == 0 {
		http.Error(w, http.StatusText(403), 403)
		return
	}

	msg, e := messageRepository.Add(Message{
		Text: text,
	})

	if e != nil {
		fmt.Println("find error", e)
		return
	}
	res, _ := json.Marshal(msg)
	w.Write(res)
}

func main() {

	fmt.Println("start")

	defer DbConnection.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", handleGreet)
	mux.HandleFunc("/search-messages", handleSearchMessages)
	mux.HandleFunc("/messages", handleMessages)
	mux.HandleFunc("/messages/create", handleCreateMessage)

	http.ListenAndServe(":8080", mux)

}
