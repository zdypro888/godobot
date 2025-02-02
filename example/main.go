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
	dobot := godobot.NewDobot()
	if err := dobot.Connect(ctx, "/dev/cu.usbserial-840", 115200); err != nil {
		fmt.Println(err)
	}
	leftSpace, err := dobot.GetQueuedCmdLeftSpace()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("leftSpace:", leftSpace)
	dobot.SetQueuedCmdClear()
	dobot.SetHOMEParams(&godobot.HOMEParams{X: 200, Y: 200, Z: 200, R: 200}, false)
	dobot.SetQueuedCmdStartExec()

	dobot.SetQueuedCmdStopExec()

	notify := make(chan os.Signal, 1)
	signal.Notify(notify, os.Interrupt, syscall.SIGTERM)
	<-notify
}
