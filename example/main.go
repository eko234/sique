package main

import (
	"fmt"
	"github.com/eko234/sique"
	"time"
)

func main() {
	done := make(chan any)
	println("STARTING SERVER...")
	go sique.Serve(":8096", nil)
	time.Sleep(time.Second * 2)
	println("STARTING CLIENT...")

	c1, err := sique.NewMQClient(":8096", nil)
	if err != nil {
		panic(err)
	}

	c2, err := sique.NewMQClient(":8096", nil)
	if err != nil {
		panic(err)
	}

	go func() {
		for ix := 0; ix < 30; ix++ {
			c2.Spit([]byte(fmt.Sprint(ix)))
		}
	}()

	go func() {
		mc, err := c1.Consoom()
		if err != nil {
			panic(err)
		}
		for m := range mc {
			fmt.Printf("1::id:%s,got:%s\n", m.ID(), m.Read())
			if err := m.Ack(); err != nil {
				fmt.Printf("1::ERR ACKING:%s\n", err.Error())
			}
		}
	}()

	go func() {
		mc, err := c2.Consoom()
		if err != nil {
			panic(err)
		}
		for m := range mc {
			fmt.Printf("2::id:%s,got:%s\n", m.ID(), m.Read())
		}
	}()

	<-done
}
