package godobot

// ColorPort 颜色传感器端口
type ColorPort uint8

const (
	CL_PORT_GP1 ColorPort = iota
	CL_PORT_GP2
	CL_PORT_GP4
	CL_PORT_GP5
)

// InfraredPort 红外传感器端口
type InfraredPort uint8

const (
	IF_PORT_GP1 InfraredPort = iota
	IF_PORT_GP2
	IF_PORT_GP4
	IF_PORT_GP5
)

// Pose 实时位姿
type Pose struct {
	X          float32
	Y          float32
	Z          float32
	R          float32
	JointAngle [4]float32
}

// Kinematics 运动学参数
type Kinematics struct {
	Velocity     float32
	Acceleration float32
}

// HOMEParams HOME相关参数
type HOMEParams struct {
	X float32
	Y float32
	Z float32
	R float32
}

// HOMECmd HOME命令
type HOMECmd struct {
	Reserved uint32
}

// ParallelOutputCmd 并行输出命令结构
type ParallelOutputCmd struct {
	Ratio   uint8
	Address uint16
	Level   uint8
}

// AutoLevelingCmd 自动调平命令
type AutoLevelingCmd struct {
	ControlFlag uint8
	Precision   float32
}

// HHTTrigMode 手持示教触发模式
type HHTTrigMode uint8

const (
	TriggeredOnKeyReleased HHTTrigMode = iota
	TriggeredOnPeriodicInterval
)

// EndEffectorParams 末端执行器参数
type EndEffectorParams struct {
	XBias float32
	YBias float32
	ZBias float32
}

// ArmOrientation 机械臂方向
type ArmOrientation uint8

const (
	LeftyArmOrientation ArmOrientation = iota
	RightyArmOrientation
)

// JOGJointParams JOG关节参数
type JOGJointParams struct {
	Velocity     [4]float32
	Acceleration [4]float32
}

// JOGCoordinateParams JOG坐标参数
type JOGCoordinateParams struct {
	Velocity     [4]float32
	Acceleration [4]float32
}

// JOGLParams JOGL参数
type JOGLParams struct {
	Velocity     float32
	Acceleration float32
}

// JOGCommonParams JOG通用参数
type JOGCommonParams struct {
	VelocityRatio     float32
	AccelerationRatio float32
}

// JOGCmd JOG命令
type JOGCmd struct {
	IsJoint uint8
	Cmd     uint8
}

const (
	JogIdle = iota
	JogAPPressed
	JogANPressed
	JogBPPressed
	JogBNPressed
	JogCPPressed
	JogCNPressed
	JogDPPressed
	JogDNPressed
	JogEPPressed
	JogENPressed
)

// PTPJointParams PTP关节参数
type PTPJointParams struct {
	Velocity     [4]float32
	Acceleration [4]float32
}

// PTPCoordinateParams PTP坐标参数
type PTPCoordinateParams struct {
	XYZVelocity     float32
	RVelocity       float32
	XYZAcceleration float32
	RAcceleration   float32
}

// PTPLParams PTPL参数
type PTPLParams struct {
	Velocity     float32
	Acceleration float32
}

// PTPJumpParams PTP跳跃参数
type PTPJumpParams struct {
	JumpHeight float32
	ZLimit     float32
}

// PTPJump2Params PTP跳跃2参数
type PTPJump2Params struct {
	StartJumpHeight float32
	EndJumpHeight   float32
	ZLimit          float32
}

// PTPCommonParams PTP通用参数
type PTPCommonParams struct {
	VelocityRatio     float32
	AccelerationRatio float32
}

// PTPMode PTP运动模式
type PTPMode uint8

const (
	PTPJUMPXYZMode PTPMode = iota
	PTPMOVJXYZMode
	PTPMOVLXYZMode
	PTPJUMPANGLEMode
	PTPMOVJANGLEMode
	PTPMOVLANGLEMode
	PTPMOVJANGLEINCMode
	PTPMOVLXYZINCMode
	PTPMOVJXYZINCMode
	PTPJUMPMOVLXYZMode
)

// PTPCmd PTP命令
type PTPCmd struct {
	PTPMode PTPMode
	X       float32
	Y       float32
	Z       float32
	R       float32
}

// PTPWithLCmd 带L轴的PTP命令
type PTPWithLCmd struct {
	PTPMode PTPMode
	X       float32
	Y       float32
	Z       float32
	R       float32
	L       float32
}

// DeviceCountInfo 设备计数信息
type DeviceCountInfo struct {
	DeviceRunTime  uint64
	DevicePowerOn  uint32
	DevicePowerOff uint32
}

// 连接错误码
const (
	DobotConnectNoError = iota
	DobotConnectNotFound
	DobotConnectOccupied
)

// 通信错误码
const (
	DobotCommunicateNoError = iota
	DobotCommunicateBufferFull
	DobotCommunicateTimeout
	DobotCommunicateInvalidParams
)

// CPParams CP参数
type CPParams struct {
	PlanAcc       float32
	JuncitionVel  float32
	AccOrPeriod   float32
	RealTimeTrack uint8
}

// CPMode CP模式
type CPMode uint8

const (
	CPRelativeMode CPMode = iota // 相对模式
	CPAbsoluteMode               // 绝对模式
)

// CPCmd CP命令
type CPCmd struct {
	CPMode   CPMode
	X        float32
	Y        float32
	Z        float32
	Velocity float32
}

// CPCommonParams CP通用参数
type CPCommonParams struct {
	VelocityRatio     float32
	AccelerationRatio float32
}

// ARCParams ARC参数
type ARCParams struct {
	XYZVelocity     float32
	RVelocity       float32
	XYZAcceleration float32
	RAcceleration   float32
}

// ARCCommonParams ARC通用参数
type ARCCommonParams struct {
	VelocityRatio     float32
	AccelerationRatio float32
}

// ARCCmd ARC命令
type ARCCmd struct {
	CirPoint struct {
		X float32
		Y float32
		Z float32
		R float32
	}
	ToPoint struct {
		X float32
		Y float32
		Z float32
		R float32
	}
}

// CircleCmd 圆弧命令
type CircleCmd struct {
	CirPoint struct {
		X float32
		Y float32
		Z float32
		R float32
	}
	ToPoint struct {
		X float32
		Y float32
		Z float32
		R float32
	}
	Count uint32
}

// WAITCmd 等待命令
type WAITCmd struct {
	Timeout uint32
}

// TRIGMode 触发模式
type TRIGMode uint8

const (
	TRIGInputIOMode TRIGMode = iota
	TRIGADCMode
)

// TRIGInputIOCondition 触发IO条件
type TRIGInputIOCondition uint8

const (
	TRIGInputIOEqual TRIGInputIOCondition = iota
	TRIGInputIONotEqual
)

// TRIGADCCondition 触发ADC条件
type TRIGADCCondition uint8

const (
	TRIGADCLT TRIGADCCondition = iota // Lower than
	TRIGADCLE                         // Lower than or Equal
	TRIGADCGE                         // Greater than or Equal
	TRIGADCGT                         // Greater Than
)

// TRIGCmd 触发命令
type TRIGCmd struct {
	Address   uint8
	Mode      uint8
	Condition uint8
	Threshold float32
}

// IOFunction IO功能
type IOFunction uint8

const (
	IOFunctionDummy IOFunction = iota
	IOFunctionDO
	IOFunctionPWM
	IOFunctionDI
	IOFunctionADC
)

// IOMultiplexing IO复用
type IOMultiplexing struct {
	Address   uint8
	Multiplex uint8
}

// IODO IO数字输出
type IODO struct {
	Address uint8
	Level   uint8
}

// IOPWM IO PWM输出
type IOPWM struct {
	Address   uint8
	Frequency float32
	DutyCycle float32
}

// IODI IO数字输入
type IODI struct {
	Address uint8
	Level   uint8
}

// IOADC IO模拟输入
type IOADC struct {
	Address uint8
	Value   uint16
}

// EMotor 扩展电机
type EMotor struct {
	Index     uint8
	IsEnabled bool
	Speed     int32
}

// EMotorS 扩展步进电机
type EMotorS struct {
	Index     uint8
	IsEnabled bool
	Speed     int32
	Distance  uint32
}
