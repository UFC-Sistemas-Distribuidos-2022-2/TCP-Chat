package main

import (
	"encoding/json"

	"bufio"
	"fmt"
	"net"
)

const welcomeMsg = "Olá, por favor, me informa seu nome de usuário:"

var comands [3]string = [3]string{"/ENTRAR", "/SAIR", "/USUARIOS"}

func main() {
	var conns []*net.TCPConn
	fmt.Println("server test")

	p := make([]byte, 2048)
	addr := net.TCPAddr{
		Port: 1234,
		IP:   net.ParseIP("127.0.0.1"),
	}
	ser, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	//defer ser.Close()

	for {
		conn, err := ser.AcceptTCP()
		if err != nil {
			return
		}
		conns = append(conns, conn)
		conn.Write([]byte(welcomeMsg))
		//for _, c := range conns {
		//fmt.Println(c, c.RemoteAddr())
		//}
		go handleUserConnection(conn, p, conns)
	}
}

func handleUserConnection(c net.Conn, p []byte, conns []*net.TCPConn) {

	for {
		n, err := bufio.NewReader(c).Read(p)
		if err != nil {
			return
		}
		var userMsg UserMsg

		err = json.Unmarshal(p[:n], &userMsg)
		completAddress := userMsg.IP + ":" + userMsg.Port
		fmt.Println(completAddress)
		if userMsg.User != "" {
			for _, c := range conns {
				fmt.Println(userMsg.Body)
				if userMsg.Body != "" {
					if c.RemoteAddr().String() != completAddress {
						c.Write([]byte(userMsg.Body))
						//fmt.Println(c, c.RemoteAddr())
					}
				}

			}
		}
	}
}

type UserMsg struct {
	User string
	Body string
	IP   string
	Port string
}
type User struct {
	User string
	IP   string
	Port string
}
