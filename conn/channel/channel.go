package channel

import (
	"encoding/json"
	"github.com/eggz6/common/conn"
	"github.com/eggz6/common/conn/protocol"
	"net"
	"sync"
	"time"
)

const (
	CheckInCH uint32 = 0
	SysCH     uint32 = 1
)

type Channel struct {
	conn  net.Conn
	ID    uint32
	Start time.Time
}

type Connection struct {
	net.Conn
	Host string
	Stat uint32 // checkin, checkout, waiting, failed
}

type listener struct {
	mu      sync.RWMutex
	session map[string]*Connection
}

func (l *listener) SetSession(id string, c *Connection) {
	l.mu.Lock()
	defer l.mu.Unlock()
	//TODO 定时同步remote session 如果需要
	//TODO 检查session是否存在 还需remote session的情况

	cur, ok := l.session[id]
	if !ok {
		// TODO 踢掉信息
		cur.Close()
	}

	l.session[id] = c

}

func (c *Channel) WriteJSON(from uint32, val interface{}) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	frame := protocol.EncodeFrame(conn.JSONProtocol, c.ID, from, c.Start, data)
	_, err = c.conn.Write(frame)
	return err
}

func NewChannel(id uint32, conn net.Conn) *Channel {
	return &Channel{ID: id, conn: conn, Start: time.Now()}
}

func Listen(port string) (*Receiver, error) {
	l, err := net.Listen("tcp", port)
	if err != nil {
		return nil, err
	}

	my := &listener{mu: sync.RWMutex{}, session: map[string]*Connection{}}

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}

			go my.handShake(conn)
		}
	}()

	return nil, nil
}

func (l *listener) handShake(conn net.Conn) {
	frame, err := protocol.DecodeFrame(conn)
	if err != nil {
		// TODO 返回错误码
		conn.Close()

		return
	}

	if frame.Header.Channel() != CheckInCH {
		// TODO 返回错误码
		conn.Close()
		return
	}

	access := map[string]string{}
	err = json.Unmarshal(frame.Extra(), &access)
	if err != nil {
		// TODO 返回错误码
		conn.Close()
		return
	}

	accessID, _ := access["access_id"]
	accessKey, _ := access["access_key"]

	ok, err := auth(accessID, accessKey)
	if err != nil {
		// TODO 返回错误码
		conn.Close()
		return
	}

	if !ok {
		// TODO 返回错误码
		conn.Close()
		return
	}

	// 保存本地session
	connection := &Connection{
		Conn:      conn,
		Stat:
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// TODO 增加Remote Session
	l.SetSession(accessID, connection)

}

func auth(accessID string, accessKey string) (bool, error) {

	return true, nil
}
