package websocket

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

func main() {
	// 启动WebSocket服务器
	startServer()

	// 打印聊天界面
	fmt.Println("Welcome to the chat room!")
	fmt.Println("Enter your name:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	name := scanner.Text()
	fmt.Println("Type 'exit' to quit.")

	// 连接WebSocket服务器
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080", nil)
	if err != nil {
		log.Fatal("dial error:", err)
	}
	defer c.Close()

	// 向WebSocket服务器发送消息
	go func() {
		for {
			scanner.Scan()
			text := scanner.Text()
			if strings.ToLower(text) == "exit" {
				os.Exit(0)
			}
			message := fmt.Sprintf("%s: %s", name, text)
			err := c.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Println("write error:", err)
				return
			}
		}
	}()

	// 从WebSocket服务器接收消息并打印
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			return
		}
		fmt.Println(string(message))
	}
}
