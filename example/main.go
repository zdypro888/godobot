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
		return
	}
	dobot.ClearAllAlarmsState(ctx)
	dobot.SetQueuedCmdClear(ctx)
	dobot.SetQueuedCmdStartExec(ctx)
	defer dobot.SetQueuedCmdStopExec(ctx)

	notify := make(chan os.Signal, 1)
	signal.Notify(notify, os.Interrupt, syscall.SIGTERM)
	<-notify
}
