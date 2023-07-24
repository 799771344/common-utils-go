package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// WebSocket客户端的连接结构体
type connection struct {
	// WebSocket连接
	ws *websocket.Conn
	// 消息的缓冲区
	send chan []byte
}

// WebSocket服务器的结构体
type server struct {
	// 所有的连接
	connections map[*connection]bool
	// 接收新连接的通道
	register chan *connection
	// 断开连接的通道
	unregister chan *connection
	// 广播消息的通道
	broadcast chan []byte
}

// 向WebSocket客户端发送消息
func (c *connection) write() {
	for {
		select {
		// 从缓冲区读取消息并发送
		case message, ok := <-c.send:
			if !ok {
				// 通道已经关闭
				c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.ws.WriteMessage(websocket.TextMessage, message)
		}
	}
}

// 从WebSocket客户端接收消息
func (c *connection) read(s *server) {
	defer func() {
		s.unregister <- c
		c.ws.Close()
	}()

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			s.unregister <- c
			c.ws.Close()
			break
		}
		// 将接收到的消息放入广播通道
		s.broadcast <- message
	}
}

// WebSocket服务器的请求处理函数
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 升级HTTP连接为WebSocket连接
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	s.register <- c
	go c.write()
	c.read(s)
}

// 启动WebSocket服务器
func startServer() *server {
	s := &server{
		connections: make(map[*connection]bool),
		register:    make(chan *connection),
		unregister:  make(chan *connection),
		broadcast:   make(chan []byte),
	}
	go func() {
		for {
			select {
			// 接收新连接
			case c := <-s.register:
				s.connections[c] = true
			// 断开连接
			case c := <-s.unregister:
				if _, ok := s.connections[c]; ok {
					delete(s.connections, c)
					close(c.send)
				}
			// 广播消息
			case message := <-s.broadcast:
				for c := range s.connections {
					select {
					case c.send <- message:
					default:
						close(c.send)
						delete(s.connections, c)
					}
				}
			}
		}
	}()
	return s
}
