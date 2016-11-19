package gotcpool

import (
	"net"
	"sync"
)

type Tcpool struct {
	_addr           *net.TCPAddr
	_connects       chan *net.TCPConn
	_max_len        int
	_created_number int
	_creat_mutex    sync.Mutex
}

func NewTcpool(addr string, min_len int,max_len int) *Tcpool {
	if min_len>max_len {
		panic("min_len must be smaller than max_len")
	}

	if a,err:=net.ResolveTCPAddr("tcp",addr);err!=nil {
		panic(err)
	} else {
		p:=&Tcpool{_addr:a,_connects:make(chan *net.TCPConn, max_len), _max_len:max_len}

		for i:=0;i<min_len;i++ {
			c,_:=p._CreateConnect()

			p._connects<-c
		}

		return p
	}
}


func (p *Tcpool)_NewConnect() *net.TCPConn {
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

func (p *Tcpool)_PoolLen() int {
	return len(p._connects)
}

func (p *Tcpool)_CreateConnect()  (c *net.TCPConn,all_done bool) {
	p._creat_mutex.Lock()
	defer p._creat_mutex.Unlock()
	if p._created_number<p._max_len {
		c=p._NewConnect()
		p._created_number++
		return c,false
	}
	return nil,true
}


func (p *Tcpool)_Take() *net.TCPConn {
	select {
	case c:=<-p._connects:
		return c
	default:
		c,all_done:=p._CreateConnect()
		if all_done {
			return <-p._connects
		} else {
			return c
		}
	}

	panic("Taking conntect occurs a unkown wrong")
}

func (p *Tcpool)_Revert(c *net.TCPConn) {
	select {
		case p._connects<-c:
		default:
			panic("Take too many connect,can't bring back the connect")
	}

}