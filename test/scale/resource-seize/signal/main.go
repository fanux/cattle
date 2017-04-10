package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("signal test")
	go func() {
		for {
			c := make(chan os.Signal, 1)
			//signal.Notify(c, syscall.SIGTERM)
			signal.Notify(c)
			s := <-c
			fmt.Println("Got signal:", s)
		}
	}()

	time.Sleep(time.Second * 3)
	syscall.Kill(os.Getpid(), syscall.SIGTERM)

	time.Sleep(time.Second * 100)
}
