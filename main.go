package main

import(
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) //connected clients
var broadcast = make(chan Message)           //broadcast channel

//configure the upgrader
var upgrader = websocket.Upgrader{}

//Define our message object
type Message struct{
	Email string `json:"email"`
	Username string `json:"username"`
	Message string `json:"message"`
}

func main() {
	//simple file fileserver
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/",fs)

	
	http.HandleFunc("/ws",handleConnections)

	go handleMessages()

	log.Println("Server running on :8080")
	err := http.ListenAndServe(":8080",nil)

	if err!=nil{
		log.Fatal("ListenAndServe: ",err)
	}	
}

func handleConnections(w http.ResponseWriter,r *http.Request){
	//upgrade initial Get request to websocket

	ws,err := upgrader.Upgrade(w,r,nil)

	if err!=nil{
		log.Fatal("Upgrade: ",err)
	}

	//close the connection when funcion returns
	defer ws.Close()

	//Register new client

	clients[ws] = true

	for{
		var msg Message

		err:= ws.ReadJSON(&msg)

		if err!=nil{
			log.Printf("error: %v",err)
			delete(clients,ws)
			break
		}

		broadcast <-msg
	}
}

func handleMessages(){
	for{
		msg := <-broadcast

		for client := range clients{
			err := client.WriteJSON(msg)
			if err!=nil{
				log.Printf("error: %v",err)
				client.Close()
				delete(clients,client)
			}
		}
	}
}