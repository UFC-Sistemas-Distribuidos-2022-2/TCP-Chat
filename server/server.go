package main

import (
	"encoding/json"
	"strings"

	"bufio"
	"fmt"
	"net"
)

const welcomeMsg = "Olá, por favor, me informa seu nome de usuário:"

var comands [3]string = [3]string{"/HELP", "/SAIR", "/USUARIOS"}
var conns []*net.TCPConn

func main() {
	fmt.Println("server test")

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
	var users []Users
	var user Users
	for {
		p := make([]byte, 2048)
		conn, err := ser.AcceptTCP()
		if err != nil {
			return
		} else {
			conn.Read(p)
			user.User = string(p)
			user.Conn = conn
			users = append(users, user)
			fmt.Println(users)
			conn.Write([]byte("Bem vindo " + user.User + " digite /HELP para exibir a lista de comandos\n"))
		}
		if notIn(conn) {
			conns = append(conns, conn)
			for _, c := range conns {
				if conn != c {
					c.Write([]byte(user.User + " entrou no chat\n"))
				}
			}
			go handleUserConnection(conn, p, &users)
		}

	}
}

func handleUserConnection(c *net.TCPConn, p []byte, users *[]Users) {
	for {
		n, err := bufio.NewReader(c).Read(p)
		if err != nil {
			fmt.Printf("Some error %v\n", err)
			conns = removeConnection(conns, c)
			var username string
			*users, username = removeUser(users, c)
			// Se o usarname for vazio, quer dizer que a pessoa usou o comando para sair
			if username != "" {
				for _, user := range *users {
					if user.Conn != c {
						user.Conn.Write([]byte(username + " desconectou-se\n"))
					}

				}
			}

			return
		}
		var userMsg UserMsg
		err = json.Unmarshal(p[:n], &userMsg)
		completAddress := userMsg.IP + ":" + userMsg.Port
		if userMsg.User != "" {
			for _, c := range conns {
				if userMsg.Body != "" {
					wasInCommandList, index := wasCommandList(userMsg.Body, comands[:])
					if wasInCommandList {
						switch index {
						// /HELP
						case 0:
							if c.RemoteAddr().String() == completAddress {
								c.Write([]byte("Lista de comandos:\n" + comands[1] + "\n" + comands[2] + "\n"))
							}
							break
							// /SAIR
						case 1:
							if c.RemoteAddr().String() == completAddress {
								conns = removeConnection(conns, c)
								*users, _ = removeUser(users, c)
								c.Close()
							} else {
								c.Write([]byte(userMsg.User + " desconectou-se\n"))
							}
							break
							// /USUARIOS
						case 2:
							if c.RemoteAddr().String() == completAddress {
								fmt.Println(users)
								for _, user := range *users {
									c.Write([]byte(user.User + "\n"))
								}
							}
						}
					} else {
						if c.RemoteAddr().String() != completAddress {
							c.Write([]byte(userMsg.User + ": " + userMsg.Body))
						}
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

func wasCommandList(msg string, comandList []string) (bool, int) {
	for index := range comandList {
		if strings.HasPrefix(msg, comandList[index]) {
			return true, index
		}
	}

	return false, -1
}

type UserMsg struct {
	User string
	Body string
	IP   string
	Port string
}
type Users struct {
	User string
	Conn *net.TCPConn
}

func removeConnection(conns []*net.TCPConn, c *net.TCPConn) []*net.TCPConn {
	var aux []*net.TCPConn
	for _, conn := range conns {
		if c != conn {
			aux = append(aux, conn)
		}
	}
	return aux
}

func removeUser(users *[]Users, conn *net.TCPConn) ([]Users, string) {
	var aux []Users
	var username string
	for _, u := range *users {
		if u.Conn != conn {
			aux = append(aux, u)
		} else {
			username = u.User
		}
	}
	return aux, username
}
