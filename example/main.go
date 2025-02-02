package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/zdypro888/godobot"
)

func main() {
	ctx := context.Background()
	connector := &godobot.Connector{}
	if err := connector.Open(ctx, "/dev/cu.usbserial-840", 115200); err != nil {
		fmt.Println(err)
	}

	notify := make(chan os.Signal, 1)
	signal.Notify(notify, os.Interrupt, syscall.SIGTERM)
	<-notify
}
