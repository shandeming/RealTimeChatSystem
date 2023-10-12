package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, port int) *Client {
	//Create a Client
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: port,
		flag:       999,
	}
	//Connect to Server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, port))
	if err != nil {
		fmt.Println("net.Dial error", err)
		return nil
	}
	client.conn = conn
	//
	return client
}

// Function: Menu
func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>请重新输入有效信息<<<<")
		return false
	}
}

func (client *Client) Run() {
	for client.flag != 0 {
		//如果输入数字非法，那么重新打印消息，并等待输入
		for client.menu() != true {
		}
		//Handle different businesses according to different modes
		switch client.flag {
		case 1:
			//公聊模式
			client.BroadMessage()
			break
		case 2:
			//私聊模式
			client.PrivateChat()
			break
		case 3:
			//更新用户信息
			client.UpdateName()
			break
		}
	}
}

// Print the message from Server
func (client *Client) HandleResponse() {
	io.Copy(os.Stdout, client.conn)
}

// Update Username
func (client *Client) UpdateName() bool {
	fmt.Println(">>>>请输入用户名:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write() err", err)
		return false
	}
	return true
}

// Broadcast message
func (client *Client) BroadMessage() {
	var msg string
	fmt.Println(">>>>请输入聊天内容,exit退出:")
	fmt.Scanln(&msg)
	for msg != "exit" {

		if len(msg) != 0 {
			sendMsg := msg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write() err", err)
				break
			}
		}
		msg = ""
		fmt.Println("请输入聊天内容,exit退出:")
		fmt.Scanln(&msg)
	}

}

// Query current online User
func (client *Client) QueryOnlineUser() {
	client.conn.Write([]byte("who\n"))
}

// Private Chat
func (client *Client) PrivateChat() {
	client.QueryOnlineUser()
	fmt.Println("输入对方名字,输入exit退出")
	var targetUser string
	fmt.Scanln(&targetUser)
	for targetUser != "exit" {
		var msg string
		fmt.Println("输入消息，或输入exit退出：")
		fmt.Scanln(&msg)
		for msg != "exit" {
			sendMsg := "to|" + targetUser + "|" + msg + "\n"
			client.conn.Write([]byte(sendMsg))
			fmt.Scanln(&msg)
		}
		client.QueryOnlineUser()
		fmt.Println("输入对方名字,输入exit退出")
		fmt.Scanln(&targetUser)
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "Ip", "127.0.0.1", "Configure IPAddress")
	flag.IntVar(&serverPort, "Port", 8888, "Configure Port")
}

func main() {
	//Parse the commandline
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>>>>链接错误...")
		return
	}
	fmt.Println(">>>>>>>>链接成功")

	//Start a goroutine to print message from the Server
	go client.HandleResponse()
	//Start the service of Client
	client.Run()
}
