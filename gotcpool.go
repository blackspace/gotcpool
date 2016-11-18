package gotcpool

import (
	"net"
)

type Tcpool struct {
	_addr *net.TCPAddr
	_connects chan *net.TCPConn


}

func NewTcpool(addr string) *Tcpool {
	if a,err:=net.ResolveTCPAddr("tcp",addr);err!=nil {
		panic(err)
	} else {
		p:=&Tcpool{_addr:a,_connects:make(chan *net.TCPConn,1<<10)}

		c:=p.NewConnect()

		p._connects<-c

		return p
	}

}

func (p *Tcpool)NewConnect() *net.TCPConn {
	if c,err:=net.DialTCP("tcp",nil,p._addr);err==nil {
		return c
	} else {
		panic(err)
	}
}

func (p *Tcpool)Do(f func(c *net.TCPConn)) {
	c:=p._Take()
	defer p._Revert(c)
	f(c)
}

func (p *Tcpool)_Take() *net.TCPConn{
	select {
	case c :=<-p._connects:
		return c
	default:
		return p.NewConnect()
	}

	return nil
}

func (p *Tcpool)_Revert(c *net.TCPConn) {
	p._connects<-c
}