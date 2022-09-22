package main

import (
	"encoding/json"

	"bufio"
	"fmt"
	"net"
)

const welcomeMsg = "Olá, por favor, me informa seu nome de usuário:"

var comands [3]string = [3]string{"/ENTRAR", "/SAIR", "/USUARIOS"}
var conns []*net.TCPConn

func main() {
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
	defer ser.Close()

	for {
		conn, err := ser.AcceptTCP()
		if err != nil {
			return
		}
		conn.Write([]byte(welcomeMsg))
		//for _, c := range conns {
		//fmt.Println(c, c.RemoteAddr())
		//}
		if notIn(conn) {
			go handleUserConnection(conn, p)

		}
	}
}

func handleUserConnection(c *net.TCPConn, p []byte) {

	for {
		n, err := bufio.NewReader(c).Read(p)
		if err != nil {
			fmt.Printf("Some error %v\n", err)
		}
		var userMsg UserMsg

		err = json.Unmarshal(p[:n], &userMsg)
		completAddress := userMsg.IP + ":" + userMsg.Port
		if userMsg.User != "" {
			if notIn(c) {
				conns = append(conns, c)
				for _, conn := range conns {
					if conn != c {
						conn.Write([]byte(userMsg.User + " entrou no chat\n"))
					}

				}
			}
			for _, c := range conns {
				if userMsg.Body != "" {
					if c.RemoteAddr().String() != completAddress {
						c.Write([]byte(userMsg.User + ": " + userMsg.Body))
					}
				}

			}
		}
	}
}

func notIn(c *net.TCPConn) bool {
	var aux int = 0
	for _, conn := range conns {
		if conn.RemoteAddr() != c.RemoteAddr() {
			aux++
		}
	}
	return aux == len(conns)
}

type UserMsg struct {
	User string
	Body string
	IP   string
	Port string
}
type Users struct {
	User string
	IP   string
	Port string
}
