package sique

import (
	"net"
)

type MQClient interface {
	Consoom() (chan Msg, error)
	Spit([]byte) error
}

type mqClient struct {
	conn net.Conn
	qch  chan Msg
}

func NewMQClient(host string, onError func(error)) (MQClient, error) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	qch := make(chan Msg)

	go func() {
		for {
			if msg, err := readMsg(conn); err == nil {
				qch <- msg
			} else {
				onError(err)
			}
		}
	}()

	return &mqClient{conn, qch}, nil
}

func (m *mqClient) Consoom() (chan Msg, error) {
	return m.qch, m.overheaaaar()
}

func (m *mqClient) overheaaaar() error {
	_, err := m.conn.Write(fmtMsg("subs", uid(), ""))
	return err
}

func (m *mqClient) Spit(in []byte) error {
	_, err := m.conn.Write(fmtMsg("send", uid(), string(in)))
	return err
}
