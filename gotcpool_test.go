package gotcpool


import (
	"github.com/blackspace/goserver"
	"github.com/blackspace/goserver/command"
	"github.com/blackspace/goserver/client"
	"net"
	"sync"
	"testing"
)

func TestTcpool_Do(t *testing.T) {
	command.Commands.RegistCommand("hello", func(clt *client.Client, args ...string) string {
		return "hello,I am a robot"
	}, "say hello")

	s := goserver.NewLineServer()
	s.Start("127.0.0.1", "5050")
	defer s.Stop()

	tcpool:=NewTcpool("127.0.0.1:5050",100)

	wg:=sync.WaitGroup{}

	wg.Add(1000)

	for i:=0;i<1000;i++ {
		go tcpool.Do(func(c *net.TCPConn) {
			c.Write([]byte("hello\r\n"))

			buf := make([]byte, 256)
			_, err := c.Read(buf)

			if err != nil {
				panic(err)
			}
			wg.Done()
		})
	}

	wg.Wait()

	if tcpool._PoolLen()!=100 {
		t.Fail()
	}

	tcpool.Do(func(c *net.TCPConn) {
		c.Write([]byte("hello\r\n"))

		buf := make([]byte, 256)
		_, err := c.Read(buf)

		if err != nil  {
			panic(err)
		}
		if tcpool._PoolLen()!=99 {
			t.Fail()
		}
	})

	tcpool.Do(func(c *net.TCPConn) {
		c.Write([]byte("hello\r\n"))

		buf := make([]byte, 256)
		_, err := c.Read(buf)

		if err != nil  {
			panic(err)
		}
		if tcpool._PoolLen()!=99 {
			t.Fail()
		}
	})

	if tcpool._PoolLen()!=100 {
		t.Fail()
	}
}

