package gotcpool


import (
	"github.com/blackspace/goserver"
	"github.com/blackspace/goserver/command"
	"github.com/blackspace/goserver/client"
	"net"
	"sync"
	"testing"
	"os"
)

func TestMain(m *testing.M) {
	command.Commands.RegistCommand("hello", func(clt *client.Client, args ...string) string {
		return "hello,I am a robot"
	}, "say hello")

	s := goserver.NewLineServer()
	s.Start("127.0.0.1", "5050")
	defer s.Stop()

	os.Exit(m.Run())
}

func TestTcpool_Do(t *testing.T) {
	max_len:=1000
	tcpool:=NewTcpool("127.0.0.1:5050",10,max_len)

	if tcpool._PoolLen()!=10 {
		t.Log(tcpool._PoolLen())
		t.Fail()
	}

	wg:=sync.WaitGroup{}

	wg.Add(max_len*2)

	for i:=0;i<max_len*2;i++ {
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


	if tcpool._PoolLen()!=max_len {
		t.Log(tcpool._PoolLen())
		t.Fail()
	}

	tcpool.Do(func(c *net.TCPConn) {
		c.Write([]byte("hello\r\n"))

		buf := make([]byte, 256)
		_, err := c.Read(buf)

		if err != nil  {
			panic(err)
		}
		if tcpool._PoolLen()!=max_len-1 {
			t.Log(tcpool._PoolLen())
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
		if tcpool._PoolLen()!=max_len-1{
			t.Log(tcpool._PoolLen())
			t.Fail()
		}
	})

	if tcpool._PoolLen()!=max_len {
		t.Log(tcpool._PoolLen())
		t.Fail()
	}
}

