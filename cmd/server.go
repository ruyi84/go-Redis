package main

import "C"

import (
	"bufio"
	"fmt"
	"goRedis/model"
	"net"
	"strconv"
	"strings"
	"time"
)

func main() {
	addr := &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 3384,
		Zone: "",
	}

	tcp, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		panic(err)
	}

	defer tcp.Close()

	go func() {
		for  {
			model.SaveToFile()
			time.Sleep(time.Second)
		}
	}()

	for {
		acceptTCP, err := tcp.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go tcpRead(acceptTCP)
	}

}

func tcpRead(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		readString, err := reader.ReadString(byte('\n'))
		if err != nil {
			fmt.Println("close")
			return
		}

		trimRight := strings.TrimRight(readString, "\n")

		split := strings.Split(trimRight, " ")
		fmt.Println(split)
		switch split[0] {
		case "Ping":
			conn.Write([]byte("Pong \n"))
			continue
		case "Set":
			if len(split) != 3 {
				conn.Write([]byte("set parsing is wrong \n"))
				continue
			}
			var expireDate int
			if len(split) > 3 {
				expireDate, err = strconv.Atoi(split[4])
				if err != nil {
					failed(conn, "parsing is wrong.")
				}
			}
			model.Set(split[1], split[2], expireDate)
			conn.Write([]byte("Ok \n"))
			continue
		case "Del":
			if len(split) != 2 {
				conn.Write([]byte("set parsing is wrong \n"))
				continue
			}
			model.Del(split[1])
			continue
		case "Get":
			if len(split) != 2 {
				conn.Write([]byte("get parsing is wrong \n"))
				continue
			}
			value := model.Get(split[1])
			success(conn, value)
			continue
		}
		conn.Write([]byte("wrong \n"))
	}
}

func failed(conn net.Conn, failedMsg string) error {
	msg := fmt.Sprintf("command failed:%s \n", failedMsg)
	_, err := conn.Write([]byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func success(conn net.Conn, value interface{}) {
	msg := fmt.Sprintf("%s \n", value)
	_, err := conn.Write([]byte(msg))
	if err != nil {
		return
	}
	return
}
