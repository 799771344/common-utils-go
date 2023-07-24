package websocket

import (
	"flag"     // flag包用于解析命令行参数
	"log"      // log包用于记录服务器日志信息
	"net/http" // net/http包用于实现HTTP服务器和客户端
	"strconv"  // strconv包用于数字和字符串之间的转换
	"sync"     // sync包提供了并发原语，如WaitGroup和Mutex
	"time"     // time包提供了时间相关的功能
)

var (
	addr       = flag.String("addr", ":8080", "http service address") // 端口号参数，缺省值为8080
	hub        = NewHub()                                             // 创建一个Hub对象，用于管理WebSocket连接
	portNumber int                                                    // 服务器的端口号
	once       sync.Once                                              // 用于确保端口号只被获取一次的同步锁
)

func main() {
	flag.Parse()                                  // 解析命令行参数
	http.HandleFunc("/ws", WebSocketHandler(hub)) // 在"/ws"路径下注册WebSocket处理函数
	go hub.Run()                                  // 启动Hub对象的Run方法，用于管理WebSocket连接

	once.Do(func() { // 用once确保端口号只被获取一次
		portNumber, _ = strconv.Atoi(*addr) // 将端口号字符串转换为数字
	})

	log.Printf("Server started on port %d\n", portNumber) // 记录服务器启动日志信息
	err := http.ListenAndServe(*addr, nil)                // 启动HTTP服务器，开始监听客户端连接
	if err != nil {
		log.Fatal("ListenAndServe: ", err) // 如果服务器启动失败，则记录错误日志并退出程序
	}
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) // 设置日志格式，包括日期、时间和文件名
	log.SetPrefix("[GoWebSocket] ")                      // 设置日志前缀为"[GoWebSocket] "
	log.Println("Starting server...")                    // 记录服务器启动日志信息
	time.Sleep(2 * time.Second)                          // 等待2秒，让用户看到服务器启动日志信息
}
