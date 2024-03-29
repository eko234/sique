package sique

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
)

var mebuf = make(chan internalMsg)
var ackpending = sync.Map{}
var boxes = map[string](chan []byte){}
var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

type internalMsg struct {
	id         string
	msg        Msg
	consumedAt time.Time
}

func uid() string {
	rand.Seed(time.Now().UnixNano())

	length := 10
	id := make([]rune, length)

	for i := range id {
		id[i] = chars[rand.Intn(len(chars))]
	}

	return string(id)
}

type Msg interface {
	Read() []byte
	Ack() error
	Op() string
	ID() string
}

type msg struct {
	id   string
	op   string
	m    []byte
	conn net.Conn
}

func (m *msg) Ack() error {
	_, err := m.conn.Write(fmtMsg("ackn", uid(), m.id))
	return err
}

func (m *msg) Read() []byte {
	return m.m
}

func (m *msg) ID() string {
	return m.id
}

func (m *msg) Op() string {
	return m.op
}

func fmtMsg(op string, id string, msg string) []byte {
	if !(len(op) == 0 || len(op) == 4 || len(id) != 10) {
		panic("Invalid message")
	}
	return []byte(fmt.Sprintf("%s%s%010d%s", op, id, len(msg), msg))
}

func readMsg(c net.Conn) (Msg, error) {

	buff := make([]byte, 24)
	if _, err := c.Read(buff); err != nil {
		return nil, err
	}

	op := string(buff[:4])
	id := string(buff[4:14])
	l, err := strconv.Atoi(string(buff[14:24]))
	if err != nil {
		return nil, err
	}

	mbuff := make([]byte, l)

	if _, err := c.Read(mbuff); err != nil {
		return nil, err
	}

	return &msg{string(id), string(op), mbuff, c}, nil
}
