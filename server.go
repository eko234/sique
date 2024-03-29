package sique

import (
	"net"
	"time"
)

func Serve(host string, onErr func(error)) error {
	l, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}

	go func() {
		for {
			now := time.Now()
			ackpending.Range(func(k any, v any) bool {
				m, ok := v.(internalMsg)
				if !ok {
					return false
				}
				if m.consumedAt.Add(time.Minute * 3).Before(now) {
					mebuf <- m
					ackpending.Delete(k)
				}

				return true
			})

			time.Sleep(time.Second * 5)
		}
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			onErr(err)
			continue
		}

		go handle(conn, onErr)
	}
}

func handle(conn net.Conn, onErr func(error)) {
	defer conn.Close()

	overhear := make(chan *bool)

	go func() {
		for {
			if msg, err := readMsg(conn); err == nil {
				switch msg.Op() {
				case "send":
					id := uid()
					mebuf <- internalMsg{id, msg, time.Now()}
				case "ackn":
					ackpending.Delete(string(msg.Read()))
				case "subs":
					overhear <- nil
				}
			} else {
				onErr(err)
			}
		}
	}()

	<-overhear

	for {
		msg := <-mebuf
		ackpending.Store(msg.id, msg)
		if _, err := conn.Write(fmtMsg("recv", msg.id, string(msg.msg.Read()))); err != nil {
			onErr(err)
		}
	}
}
