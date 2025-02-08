package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/zdypro888/godobot/draw"
)

func main() {
	robot, err := draw.NewRobot("/dev/cu.usbserial-840", 115200)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer robot.Close()
	robot.DrawInit()
	// robot.Capture(context.Background(), true)

	notify := make(chan os.Signal, 1)
	signal.Notify(notify, os.Interrupt, syscall.SIGTERM)
	<-notify
}
