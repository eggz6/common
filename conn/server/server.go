package server

import (
	"context"
	"encoding/json"
	"errors"
	"net"

	"github.com/eggz6/common/conn"
	"github.com/eggz6/common/conn/protocol"
)

var emtpyAuth = EmptyAuth("empty_auth")

type option struct {
	port    string
	session conn.Session
	conn.Auth
	listener net.Listener
	handles  []handler
}

type svr struct {
	option
}

type EmptyAuth string

func (ea EmptyAuth) Pass(id, key string) (bool, error) {
	return true, nil
}

type Option func(opt *option)

func Use(handles ...handler) Option {
	return func(opt *option) {
		opt.handles = append(opt.handles, handles...)
	}
}

func Serve(ctx context.Context, options ...Option) error {
	op := &option{Auth: emtpyAuth}

	for _, opt := range options {
		opt(op)
	}

	s := &svr{option: *op}
	l, err := net.Listen("tcp", s.port)

	if err != nil {
		return err
	}

	s.listener = l

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				//TODO log here

				return
			}
			go s.handle(conn)
		}
	}()

	select {
	case <-ctx.Done():
		l.Close()
	}

	return nil
}

type handler func(ctx context.Context, frame *protocol.Frame)

func (s *svr) handle(conn net.Conn, handles ...handler) error {
	frame, err := protocol.DecodeFrame(conn)
	if err != nil {
		// TODO 返回错误码
		conn.Close()

		return err
	}

	if frame.Header.Channel() != 0 {
		// TODO 返回错误码
		conn.Close()
		return err
	}

	if s.auth != emtpyAuth {
		access := map[string]string{}
		err = json.Unmarshal(frame.Extra(), &access)
		if err != nil {
			// TODO 返回错误码
			conn.Close()
			return err
		}

		accessID, _ := access["access_id"]
		accessKey, _ := access["access_key"]

		ok, err := s.auth.Pass(accessID, accessKey)
		if err != nil {
			// TODO 写回错误码
			conn.Close()
			return err
		}

		if !ok {
			// TODO 写回错误码
			conn.Close()
			return errors.New("auth failed.")
		}
	}

	for {
		frame, err := protocol.DecodeFrame(conn)
		if err != nil {
			conn.Close()

			return
		}

		for _, h := range handles {
			h(ctx, frame)
		}
	}
}

type Connection struct {
	net.Conn
	Host string
	Stat uint32 // checkin, checkout, waiting, failed
}
