package main

import (
	"context"
	"fmt"
	"math"

	"github.com/gomlx/bsplines"
	"github.com/zdypro888/godobot"
)

type Robot struct {
	dobot *godobot.Dobot
}

func NewRobot(port string, baudrate uint32) (*Robot, error) {
	dobot := godobot.NewDobot()
	if err := dobot.Connect(port, baudrate); err != nil {
		return nil, err
	}
	if err := dobot.ClearAllAlarmsState(); err != nil {
		return nil, err
	}
	if err := dobot.SetQueuedCmdClear(); err != nil {
		return nil, err
	}
	if err := dobot.SetQueuedCmdStartExec(); err != nil {
		return nil, err
	}
	return &Robot{dobot: dobot}, nil
}

func (robot *Robot) Close() error {
	if robot.dobot == nil {
		return nil
	}
	return robot.dobot.Close()
}

func (robot *Robot) Stop() error {
	if err := robot.dobot.SetQueuedCmdStopExec(); err != nil {
		return err
	}
	return nil
}

func (robot *Robot) EnmergyStop() error {
	if err := robot.dobot.SetQueuedCmdForceStopExec(); err != nil {
		return err
	}
	return nil
}

func (robot *Robot) HomeZero() error {
	if _, err := robot.dobot.SetHOMEParams(&godobot.HOMEParams{X: 160, Y: 0, Z: 0, R: 0}, true); err != nil {
		return err
	}
	if err := robot.dobot.QueuedComplete(func() (uint64, error) {
		return robot.dobot.SetHOMECmd(&godobot.HOMECmd{}, true)
	}); err != nil {
		return err
	}
	return nil
}

func (robot *Robot) Capture(ctx context.Context, debug bool) ([]*godobot.Pose, error) {
	robot.dobot.SetHHTTrigMode(godobot.TriggeredOnPeriodicInterval)
	robot.dobot.SetHHTTrigOutputEnabled(true)
	var postions []*godobot.Pose
	for {
		select {
		case <-ctx.Done():
			return postions, nil
		default:
			if enabled, err := robot.dobot.GetHHTTrigOutput(); err != nil {
				return nil, err
			} else if enabled {
				pos, err := robot.dobot.GetPose()
				if err != nil {
					return nil, err
				}
				if debug {
					fmt.Printf("current pos: %f, %f, %f\n", pos.X, pos.Y, pos.Z)
				}
				postions = append(postions, pos)
			}
		}
	}
}

func (robot *Robot) DrawInit() error {
	// 设置精度和速度
	ptpCommonParams := godobot.PTPCommonParams{
		VelocityRatio:     200.0, // 速度百分比
		AccelerationRatio: 200.0, // 加速度百分比
	}
	if _, err := robot.dobot.SetPTPCommonParams(&ptpCommonParams, false); err != nil {
		return err
	}
	// 设置 PTP 关节模式的速度和加速度
	ptpJointParams := &godobot.PTPJointParams{
		Velocity:     [4]float32{200.0, 200.0, 200.0, 200.0},
		Acceleration: [4]float32{200.0, 200.0, 200.0, 200.0},
	}
	if _, err := robot.dobot.SetPTPJointParams(ptpJointParams, true); err != nil {
		return err
	}
	// 设置 PTP 坐标模式的速度和加速度
	ptpCoordinateParams := godobot.PTPCoordinateParams{
		XYZVelocity:     200.0, // X, Y, Z 速度
		RVelocity:       200.0, // 旋转轴速度
		XYZAcceleration: 200.0, // 加速度
		RAcceleration:   200.0, // 旋转轴加速度
	}
	if _, err := robot.dobot.SetPTPCoordinateParams(&ptpCoordinateParams, false); err != nil {
		return err
	}
	// 设置 CP 的速度和加速度
	cpParams := godobot.CPParams{
		PlanAcc:      40, // 轨迹加速度（平缓）
		JuncitionVel: 10, // 拐点速度（减少顿挫）
		AccOrPeriod:  80, // 全局加速度比率
	}
	if _, err := robot.dobot.SetCPParams(&cpParams, true); err != nil {
		return err
	}
	return nil
}

// 计算曲率
func curvatureAt(points []Point, i int) float64 {
	if i <= 0 || i >= len(points)-1 {
		return 0 // 边界点不计算
	}
	// 取相邻的三个点
	p0 := points[i-1]
	p1 := points[i]
	p2 := points[i+1]

	// 计算一阶导数（速度）
	dx1 := p1.X - p0.X
	dy1 := p1.Y - p0.Y
	dx2 := p2.X - p1.X
	dy2 := p2.Y - p1.Y

	// 计算二阶导数（加速度）
	ddx := dx2 - dx1
	ddy := dy2 - dy1

	// 计算曲率
	curvature := math.Abs(float64(dx1*ddy-dy1*ddx)) / math.Pow(float64(dx1*dx1+dy1*dy1), 1.5)
	return curvature
}

func (robot *Robot) Draw(trajectories *Trajectories, z float32, scale float64, bspline bool) error {
	for _, stroke := range trajectories.Strokes {
		var smoothPoints []Point
		if len(stroke) < 3 || !bspline {
			for _, point := range stroke {
				smoothPoints = append(smoothPoints, Point{X: point.X / float32(scale), Y: point.Y / float32(scale)})
			}
		} else {
			// 提取 X 和 Y 坐标
			var xKnots, yKnots []float64
			for _, point := range stroke {
				xKnots = append(xKnots, float64(point.X)/scale)
				yKnots = append(yKnots, float64(point.Y)/scale)
			}
			// 创建 B-Spline（3 阶）
			degree := 3
			bSplineX := bsplines.NewRegular(degree, len(xKnots)).WithControlPoints(xKnots)
			bSplineY := bsplines.NewRegular(degree, len(yKnots)).WithControlPoints(yKnots)
			// 生成平滑曲线上的点（采样 50 个点）
			numSamples := len(xKnots) * 5
			smoothPoints = make([]Point, numSamples)
			// 应用 Kalman 过滤器（减少误差）
			alpha := float32(0.9)
			var lastX, lastY float32
			for i := 0; i < numSamples; i++ {
				t := float64(i) / float64(numSamples-1) // ✅ 修正 t
				x := float32(bSplineX.Evaluate(t))
				y := float32(bSplineY.Evaluate(t))
				if i > 0 { // ✅ 只在 i > 0 时应用 Kalman 滤波
					x = alpha*x + (1-alpha)*lastX
					y = alpha*y + (1-alpha)*lastY
				}
				smoothPoints[i] = Point{X: x, Y: y}
				lastX, lastY = x, y
			}
		}
		firstPoint := smoothPoints[0]
		goFirstPoint := &godobot.PTPCmd{
			PTPMode: godobot.PTPJUMPXYZMode,
			X:       200 - firstPoint.Y,
			Y:       50 - firstPoint.X,
			Z:       z,
			R:       0,
		}
		if err := robot.dobot.QueuedComplete(func() (uint64, error) {
			return robot.dobot.SetPTPCmd(goFirstPoint, true)
		}); err != nil {
			return err
		}
		prevPoint := firstPoint
		for i, currPoint := range smoothPoints {
			var curvature float32
			if bspline {
				curvature = float32(10 + 40*(1-curvatureAt(smoothPoints, i)))
			} else {
				curvature = 200.0
			}
			deltaX := prevPoint.Y - currPoint.Y
			deltaY := prevPoint.X - currPoint.X
			prevPoint = currPoint
			movePoint := &godobot.CPCmd{
				CPMode:   godobot.CPRelativeMode,
				X:        deltaX,
				Y:        deltaY,
				Z:        0,
				Velocity: curvature,
			}
			if _, err := robot.dobot.QueuedSend(func() (uint64, error) {
				return robot.dobot.SetCPCmd(movePoint, true)
			}); err != nil {
				return err
			}
		}
	}
	goHomePoint := &godobot.PTPCmd{
		PTPMode: godobot.PTPJUMPXYZMode,
		X:       160,
		Y:       0,
		Z:       0,
		R:       0,
	}
	if err := robot.dobot.QueuedComplete(func() (uint64, error) {
		return robot.dobot.SetPTPCmd(goHomePoint, true)
	}); err != nil {
		return err
	}
	return nil
}
