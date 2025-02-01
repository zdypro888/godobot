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
	if err := connector.Open(ctx, "/dev/cu.usbserial-840"); err != nil {
		fmt.Println(err)
	}
	connector.SendMessage(ctx, &godobot.Message{Id: godobot.ProtocolQueuedCmdLeftSpace, RW: false, IsQueued: false})
	notify := make(chan os.Signal, 1)
	signal.Notify(notify, os.Interrupt, syscall.SIGTERM)
	<-notify
}
