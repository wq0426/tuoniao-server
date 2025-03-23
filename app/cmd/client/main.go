package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	serverURL := "ws://127.0.0.1:50053"

	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	go readMessages(conn)

	fmt.Println("Connected to the server! Type your message:")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			fmt.Println("Error sending message:", err)
			break
		}
	}
}

func readMessages(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		fmt.Printf("Message from server: %s\n", string(msg))
		time.Sleep(100 * time.Millisecond) // 控制输出节奏
	}
}
