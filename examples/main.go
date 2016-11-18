package main

import (
	"github.com/blackspace/goserver"
	"github.com/blackspace/goserver/command"
	"github.com/blackspace/goserver/client"
	"github.com/blackspace/gotcpool"
	"net"
	"log"
)

func main() {
	command.Commands.RegistCommand("hello", func(clt *client.Client, args ...string) string {
		return "hello,I am a robot"
	}, "say hello")

	s := goserver.NewLineServer()
	s.Start("127.0.0.1", "5050")
	defer s.Stop()

	tcpool:=gotcpool.NewTcpool("127.0.0.1:5050")

	tcpool.Do(func(c *net.TCPConn) {
		c.Write([]byte("hello\r\n"))

		buf := make([]byte, 256)
		n, err := c.Read(buf)

		if err == nil {
			log.Println(string(buf[:n]))
		} else {
			panic(err)
		}
	})
}
