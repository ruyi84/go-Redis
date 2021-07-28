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
		for {
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
	disk := model.GetDisk()
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
			success(conn, []byte("Pong \n"))
			continue
		case "Keys":
			var search string
			if len(split) > 1 {
				search = split[1]
			}
			keys := disk.Keys(search)
			success(conn, keys)
			continue
		case "Set":
			if len(split) != 3 {
				failed(conn, "set parsing is wrong.")
				continue
			}
			var expireDate int
			if len(split) > 3 {
				expireDate, err = strconv.Atoi(split[4])
				if err != nil {
					failed(conn, "parsing is wrong.")
				}
			}
			disk.Set(split[1], split[2], expireDate)
			success(conn, "Ok")
			continue
		case "Del":
			if len(split) != 2 {
				success(conn, "set parsing is wrong")
				continue
			}
			disk.Del(split[1])
			continue
		case "Get":
			if len(split) != 2 {
				failed(conn, "get parsing is wrong")
				continue
			}
			value := disk.Get(split[1])
			success(conn, value)
			continue
		}
		failed(conn, "wrong")
	}
}

func failed(conn net.Conn, failedMsg string) {
	msg := fmt.Sprintf("command failed:%s \n", failedMsg)
	_, err := conn.Write([]byte(msg))
	if err != nil {
		return
	}
	return
}

func success(conn net.Conn, value interface{}) {
	msg := fmt.Sprintf("%s \n", value)
	_, err := conn.Write([]byte(msg))
	if err != nil {
		return
	}
	return
}
