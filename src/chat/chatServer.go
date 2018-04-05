package chat

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

//每一个连接的控制器
type WebClient struct {
	Conn    *websocket.Conn
	mt      int
	OnRead  chan []byte
	OnWrite chan []byte
	in      chan bool
	out     chan bool
}

//全部连接
type ClientManager struct {
	Clients map[string]WebClient //所有客户端的列表
	newMsg  chan []byte          //有新的消息
}

func (cm *ClientManager) start() {

	for {
		select {
		//接收到新的消息，那么就通知所有的客户端
		case newMsg := <-cm.newMsg:
			for _, client := range cm.Clients {
				client.OnWrite <- newMsg
			}
		}
	}
}

func (cm *ClientManager) addClient(client WebClient) {
	cm.Clients[client.Conn.RemoteAddr().String()] = client
	go client.init()
}

func (client *WebClient) init() {
	defer func() {
		client.out <- true
		client.Conn.Close()
	}()

	//开启循环监听
	go client.monitor()

	client.in <- true

	for {
		mt, msg, err := client.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			client.out <- true
			client.Conn.Close()
			break
		}
		client.mt = mt
		client.OnRead <- msg
	}

}

func (client *WebClient) monitor() {
	for {
		select {
		case inMsg := <-client.OnRead:
			fmt.Println(client.Conn.RemoteAddr().String())
			cm.newMsg <- []byte(fmt.Sprintf("[%s]:%s",client.Conn.RemoteAddr().String(),string(inMsg)))
		case outMsg := <-client.OnWrite:
			reply := outMsg
			client.Conn.WriteMessage(client.mt, reply)
		case <-client.in:
			fmt.Printf("client[%s] in！\n", client.Conn.RemoteAddr().String())
		case <-client.out:
			fmt.Printf("client[%s] out！\n", client.Conn.RemoteAddr().String())
		}
	}
}

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


var cm = ClientManager{make(map[string]WebClient), make(chan []byte)}


func Echo(c *gin.Context) {
	//获取连接
	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := WebClient{conn, 0, make(chan []byte), make(chan []byte), make(chan bool), make(chan bool)}
	go cm.start()
	cm.addClient(client)
}
