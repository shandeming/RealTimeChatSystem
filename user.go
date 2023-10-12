package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	Server *Server
}

// Create a User API
func NewUser(conn net.Conn, server *Server) *User {
	UserAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   UserAddr,
		Addr:   UserAddr,
		C:      make(chan string),
		conn:   conn,
		Server: server,
	}
	//
	go user.ListenMessage()
	return user
}

// Online
func (this *User) Online() {
	//Add the New User to OnlineMap
	this.Server.maplock.Lock()
	this.Server.OnlineMap[this.Name] = this
	this.Server.maplock.Unlock()
	//Broadcasting the user's online message
	this.Server.Broadcast(this, "已上线")
}

// Offline
func (this *User) Offline() {
	this.Server.maplock.Lock()
	delete(this.Server.OnlineMap, this.Name)
	this.Server.maplock.Unlock()
	this.Server.Broadcast(this, "下线")
}

func (this *User) SendMessage(msg string) {
	this.conn.Write([]byte(msg))
}

// Handle message
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		for _, cli := range this.Server.OnlineMap {
			onlineUser := "[" + cli.Addr + "]" + cli.Name + ":Online!\n"
			this.SendMessage(onlineUser)
		}
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		//Message format: rename|name
		newName := strings.Split(msg, "|")[1]
		_, ok := this.Server.OnlineMap[newName]
		if ok {
			this.SendMessage("已存在，请重新输入\n")
		} else {
			this.Server.maplock.Lock()
			delete(this.Server.OnlineMap, this.Name)
			this.Server.OnlineMap[newName] = this
			this.Server.maplock.Unlock()
			this.Name = newName
			this.SendMessage("用户名修改成功\n")
		}
	} else if len(msg) > 3 && msg[:3] == "to|" {
		//Msg Format: to|张三|Message
		//1 Get the other party's username
		toUserName := strings.Split(msg, "|")[1]
		if toUserName == "" {
			this.SendMessage("用户名输入格式错误，请重新输入\n")
			return
		}
		// 2 Get the other party's User Object
		toUser, ok := this.Server.OnlineMap[toUserName]
		if !ok {
			this.SendMessage("找不到对应User，请重新输入\n")
			return
		}
		//3 Send message
		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.SendMessage("无消息内容，请重发")
			return
		}
		toUser.SendMessage(this.Name + "对你说:" + content + "\n")
	} else {
		this.Server.Broadcast(this, msg)
	}
}

// Monitor User's channel
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
