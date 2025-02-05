package main

import (
	"context"

	"github.com/zdypro888/godobot"
)

type Robot struct {
	botctx context.Context
	cancel context.CancelFunc
	dobot  *godobot.Dobot
}

func NewRobot(ctx context.Context, port string, baudrate uint32) (*Robot, error) {
	botctx, cancel := context.WithCancel(ctx)
	dobot := godobot.NewDobot()
	if err := dobot.Connect(botctx, port, baudrate); err != nil {
		cancel()
		return nil, err
	}
	if err := dobot.ClearAllAlarmsState(botctx); err != nil {
		cancel()
		return nil, err
	}
	if err := dobot.SetQueuedCmdClear(botctx); err != nil {
		cancel()
		return nil, err
	}
	if err := dobot.SetQueuedCmdStartExec(botctx); err != nil {
		cancel()
		return nil, err
	}
	return &Robot{botctx: botctx, cancel: cancel, dobot: dobot}, nil
}

func (robot *Robot) Close() {
	if robot.cancel != nil {
		robot.cancel()
	}
}

func (robot *Robot) Stop(ctx context.Context) error {
	if err := robot.dobot.SetQueuedCmdStopExec(ctx); err != nil {
		return err
	}
	return nil
}

func (robot *Robot) EnmergyStop(ctx context.Context) error {
	if err := robot.dobot.SetQueuedCmdForceStopExec(ctx); err != nil {
		return err
	}
	return nil
}

func (robot *Robot) HomeZero(ctx context.Context) error {
	if _, err := robot.dobot.SetHOMEParams(ctx, &godobot.HOMEParams{X: 160, Y: 0, Z: 0, R: 0}, true); err != nil {
		return err
	}
	if err := robot.dobot.QueuedComplete(ctx, func() (uint64, error) {
		return robot.dobot.SetHOMECmd(ctx, &godobot.HOMECmd{}, true)
	}); err != nil {
		return err
	}
	return nil
}

func (robot *Robot) Capture(ctx context.Context, cancel context.Context) ([]*godobot.Pose, error) {
	robot.dobot.SetHHTTrigMode(ctx, godobot.TriggeredOnPeriodicInterval)
	robot.dobot.SetHHTTrigOutputEnabled(ctx, true)
	var postions []*godobot.Pose
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-cancel.Done():
			return postions, nil
		default:
			if enabled, err := robot.dobot.GetHHTTrigOutput(ctx); err != nil {
				return nil, err
			} else if enabled {
				pos, err := robot.dobot.GetPose(ctx)
				if err != nil {
					return nil, err
				}
				postions = append(postions, pos)
			}
		}
	}
}

func (robot *Robot) DrawInit(ctx context.Context) error {
	// 设置精度和速度
	ptpCommonParams := godobot.PTPCommonParams{
		VelocityRatio:     50.0, // 速度百分比
		AccelerationRatio: 50.0, // 加速度百分比
	}
	if _, err := robot.dobot.SetPTPCommonParams(ctx, &ptpCommonParams, false); err != nil {
		return err
	}
	ptpCoordinateParams := godobot.PTPCoordinateParams{
		XYZVelocity:     50.0, // X, Y, Z 速度
		RVelocity:       30.0, // 旋转轴速度
		XYZAcceleration: 20.0, // 加速度
		RAcceleration:   10.0, // 旋转轴加速度
	}
	if _, err := robot.dobot.SetPTPCoordinateParams(ctx, &ptpCoordinateParams, false); err != nil {
		return err
	}
	return nil
}
