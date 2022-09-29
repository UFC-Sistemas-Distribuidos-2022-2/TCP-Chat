package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

func main() {

	var comand string = ""
	var IPPort string = ""
	var try = false
	fmt.Println("Bem vindo ao chat TCP")
	fmt.Println("Deseja entrar no chat? Digite /ENTRAR para começar")
	for comand != "/ENTRAR" {
		fmt.Scanln(&comand)
	}
	fmt.Println("Informe o ip e a porta que deseja tentar se conectar no formato IP:PORTA (default:127.0.0.1:1234)")
	fmt.Scanln(&IPPort)
	if IPPort == "" {
		IPPort = "127.0.0.1:1234"
	}
	//Para debuggar
	//IPPort = "127.0.0.1:1234"
	p := make([]byte, 2048)
	conn, err := net.Dial("tcp", IPPort)
	if err != nil {
		try = true
		for try {
			fmt.Printf("Some error %v\n", err)
			fmt.Println("Informe o ip e a porta que deseja tentar se conectar no formato IP:PORTA")
			fmt.Scanln(&IPPort)
			conn, err = net.Dial("tcp", IPPort)
			if err == nil {
				try = false
			}
		}
	}
	var address []string = strings.Split(conn.LocalAddr().String(), ":")
	fmt.Println("Me informa também seu nome de usuário: (default:guest " + address[1] + ")")
	var username string
	fmt.Scanln(&username)
	if username == "" {
		username = "guest " + address[1]
	}
	conn.Write([]byte(username))
	wg.Add(1)
	go Read(conn, p)
	go Write(conn, username, address[0], address[1])
	wg.Wait()

}

var wg sync.WaitGroup

func Read(conn net.Conn, p []byte) {
	for {
		str, err := conn.Read(p)
		if err != nil {
			wg.Done()
			return
		}
		fmt.Printf("%s", p[:str])
	}
}

func Write(conn net.Conn, username string, ip string, port string) {
	for {
		var in *bufio.Reader
		in = bufio.NewReader(os.Stdin)

		line, _ := in.ReadString('\n')
		if line != "\n" {
			p, err := jsonBuilder(username, line, ip, port)
			if err != nil {
				fmt.Printf("Some error %v", err)
				return
			}
			conn.Write(p)
		}
	}

}

func jsonBuilder(username string, msg string, ip string, port string) ([]byte, error) {
	data := map[string]interface{}{
		"User": username,
		"Body": msg,
		"IP":   ip,
		"Port": port,
	}

	return json.Marshal(data)
}
