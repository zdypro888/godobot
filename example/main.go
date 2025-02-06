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
	// trajectories, err := LoadTrajectories("trajectories.json")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	robot.DrawInit()
	// robot.Capture(context.Background(), true)
	draw.ListTrajectories(":8080", robot)

	notify := make(chan os.Signal, 1)
	signal.Notify(notify, os.Interrupt, syscall.SIGTERM)
	<-notify
}
