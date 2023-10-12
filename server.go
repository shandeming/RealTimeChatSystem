package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*User
	maplock   sync.RWMutex
	Message   chan string
}

// 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

func (this *Server) Handler(conn net.Conn) {
	//...当前链接的业务
	fmt.Println("链接建立成功")

	user := NewUser(conn, this)
	//User online
	user.Online()
	//Used to determined activity
	isLive := make(chan bool)
	//Receiving Message from User
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read error", err)
				return
			}
			// Extract User's Message and remove the trailing '\n'
			msg := string(buf[:n-1])
			//Broadcast
			user.DoMessage(msg)
			//
			isLive <- true
		}
	}()
	//
	for {
		select {
		case <-isLive:
			//
		case <-time.After(time.Second * 300):
			//if it is timeout
			//force logout
			user.SendMessage("你被踢出")
			//Destory resource
			close(user.C)
			//Close connection
			conn.Close()
			//
			return
		}
	}

}
func (this *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

// 启动服务器的接口
func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	//close listen socket
	defer listener.Close()
	go this.MessageListener()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		//do handler
		go this.Handler(conn)
	}
}

// Send message to ALL users from the Message Channel
func (this *Server) MessageListener() {
	for {
		msg := <-this.Message
		//
		this.maplock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.maplock.Unlock()
	}
}
