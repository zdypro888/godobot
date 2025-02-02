package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zdypro888/godobot"
)

func main() {
	ctx := context.Background()
	dobot := godobot.NewDobot()
	if err := dobot.Connect(ctx, "/dev/cu.usbserial-840", 115200); err != nil {
		fmt.Println(err)
	}
	dobot.SetQueuedCmdClear()
	dobot.SetHOMEParams(&godobot.HOMEParams{X: 200, Y: 200, Z: 200, R: 200})
	jointParams := &godobot.PTPJointParams{}
	jointParams.Velocity[0] = 200
	jointParams.Velocity[1] = 200
	jointParams.Velocity[2] = 200
	jointParams.Velocity[3] = 200
	jointParams.Acceleration[0] = 200
	jointParams.Acceleration[1] = 200
	jointParams.Acceleration[2] = 200
	jointParams.Acceleration[3] = 200
	dobot.SetPTPJointParams(jointParams)
	dobot.SetPTPCommonParams(&godobot.PTPCommonParams{VelocityRatio: 100, AccelerationRatio: 100})

	dobot.SetHOMECmd(&godobot.HOMECmd{Reserved: 0})

	var lastIndex uint64
	for i := 0; i < 5; i++ {
		var offset float32
		if i%2 == 0 {
			offset = 50
		} else {
			offset = -50
		}
		lastIndex, _ = dobot.SetPTPCmd(&godobot.PTPCmd{PTPMode: godobot.PTPJUMPXYZMode, X: 200 + offset, Y: 200 + offset, Z: 200 + offset, R: 200 + offset})
	}

	dobot.SetQueuedCmdStartExec()
	for {
		index, _ := dobot.GetQueuedCmdCurrentIndex()
		if lastIndex > index {
			time.Sleep(100 * time.Millisecond)
		} else {
			break
		}
	}

	dobot.SetQueuedCmdStopExec()

	notify := make(chan os.Signal, 1)
	signal.Notify(notify, os.Interrupt, syscall.SIGTERM)
	<-notify
}
