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

const prefix = "/"

func main() {

	var comand string = ""
	var IPPort string = ""
	var msg string = ""
	fmt.Println("Bem vindo ao chat TCP")
	fmt.Println("Deseja entrar no chat? Digite /ENTRAR para começar")
	for comand != "/ENTRAR" {
		fmt.Scan(&comand)
	}

	fmt.Println("Informe o ip e a porta que deseja tentar se conectar no formato IP:PORTA")
	fmt.Scan(&IPPort)
	p := make([]byte, 2048)
	conn, err := net.Dial("tcp", IPPort)
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}

	conn.Read((p))
	fmt.Printf("%s", p)
	var username string
	fmt.Scan(&username)
	var address []string = strings.Split(conn.LocalAddr().String(), ":")
	p, err = jsonBuilder(username, msg, address[0], address[1])
	conn.Write([]byte(p))
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
		//var msg string
		//fmt.Scan(&msg)
		in := bufio.NewReader(os.Stdin)

		line, err := in.ReadString('\n')
		p, err := jsonBuilder(username, line, ip, port)
		if err != nil {
			fmt.Printf("Some error %v", err)
			return
		}
		conn.Write(p)
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
