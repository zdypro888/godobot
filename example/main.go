package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

//go:embed darw.html
var drawHtml []byte

func main() {
	robot, err := NewRobot("/dev/cu.usbserial-840", 115200)
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
	ListTrajectories(":8080", robot)

	notify := make(chan os.Signal, 1)
	signal.Notify(notify, os.Interrupt, syscall.SIGTERM)
	<-notify
}
