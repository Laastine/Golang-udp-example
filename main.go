package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

type Message struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func logger(str string) {
	var time = time.Now().Format(time.RFC3339)
	fmt.Sprintln("%s %s", time, str)
}

func main() {
	udpServer, err := net.ListenPacket("udp", ":3000")
	if err != nil {
		log.Fatal(err)
	}
	defer udpServer.Close()

	for {
		requestBuf := make([]byte, 4096)
		len, addr, err := udpServer.ReadFrom(requestBuf)
		var buf = requestBuf[:len]
		if err != nil {
			continue
		}
		var str = parseMessage(buf)
		go response(udpServer, addr, str)
	}
}

func parseMessage(msgBuf []byte) Message {
	var msg Message
	err := json.Unmarshal([]byte(msgBuf), &msg)
	if err != nil {
		fmt.Println("JSON parse error %s", err)
	}
	return msg
}

func response(udpServer net.PacketConn, addr net.Addr, msg Message) {
	time := time.Now().Format(time.RFC3339)
	responseStr := fmt.Sprintf("%v msg: hello %v!", time, msg.Value)

	udpServer.WriteTo([]byte(responseStr), addr)
	logger(fmt.Sprintf("Message %+v sent to client", msg))
}
