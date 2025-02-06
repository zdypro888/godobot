package internal

// ProtocolId 定义了所有协议命令的ID
type ProtocolId uint8

const (
	// Device information
	ProtocolFunctionDeviceInfoBase ProtocolId = 0
	ProtocolDeviceSN               ProtocolId = ProtocolFunctionDeviceInfoBase + 0 // Base + 0
	ProtocolDeviceName             ProtocolId = ProtocolFunctionDeviceInfoBase + 1 // Base + 1
	ProtocolDeviceVersion          ProtocolId = ProtocolFunctionDeviceInfoBase + 2 // Base + 2
	ProtocolDeviceWithL            ProtocolId = ProtocolFunctionDeviceInfoBase + 3 // Base + 3
	ProtocolDeviceTime             ProtocolId = ProtocolFunctionDeviceInfoBase + 4 // Base + 4
	_                                                                              // Base + 5 (跳过)
	ProtocolDeviceInfo             ProtocolId = ProtocolFunctionDeviceInfoBase + 6 // Base + 6
)

const (
	// Pose
	ProtocolFunctionPoseBase ProtocolId = 10
	ProtocolGetPose          ProtocolId = ProtocolFunctionPoseBase + 0 // Base + 0
	ProtocolResetPose        ProtocolId = ProtocolFunctionPoseBase + 1 // Base + 1
	ProtocolGetKinematics    ProtocolId = ProtocolFunctionPoseBase + 2 // Base + 2
	ProtocolGetPoseL         ProtocolId = ProtocolFunctionPoseBase + 3 // Base + 3
)

const (
	// Alarm
	ProtocolFunctionALARMBase ProtocolId = 20
	ProtocolAlarmsState       ProtocolId = ProtocolFunctionALARMBase + 0 // Base + 0
)

const (
	// HOME
	ProtocolFunctionHOMEBase ProtocolId = 30
	ProtocolHOMEParams       ProtocolId = ProtocolFunctionHOMEBase + 0 // Base + 0
	ProtocolHOMECmd          ProtocolId = ProtocolFunctionHOMEBase + 1 // Base + 1
	ProtocolAutoLeveling     ProtocolId = ProtocolFunctionHOMEBase + 2 // Base + 2
)

const (
	// HHT
	ProtocolFunctionHHTBase      ProtocolId = 40
	ProtocolHHTTrigMode          ProtocolId = ProtocolFunctionHHTBase + 0 // Base + 0
	ProtocolHHTTrigOutputEnabled ProtocolId = ProtocolFunctionHHTBase + 1 // Base + 1
	ProtocolHHTTrigOutput        ProtocolId = ProtocolFunctionHHTBase + 2 // Base + 2
)

const (
	// Arm Orientation
	ProtocolFunctionArmOrientationBase ProtocolId = 50
	ProtocolArmOrientation             ProtocolId = ProtocolFunctionArmOrientationBase + 0 // Base + 0
)

const (
	// End effector
	ProtocolFunctionEndEffectorBase ProtocolId = 60
	ProtocolEndEffectorParams       ProtocolId = ProtocolFunctionEndEffectorBase + 0 // Base + 0
	ProtocolEndEffectorLaser        ProtocolId = ProtocolFunctionEndEffectorBase + 1 // Base + 1
	ProtocolEndEffectorSuctionCup   ProtocolId = ProtocolFunctionEndEffectorBase + 2 // Base + 2
	ProtocolEndEffectorGripper      ProtocolId = ProtocolFunctionEndEffectorBase + 3 // Base + 3
)

const (
	// JOG
	ProtocolFunctionJOGBase     ProtocolId = 70
	ProtocolJOGJointParams      ProtocolId = ProtocolFunctionJOGBase + 0 // Base + 0
	ProtocolJOGCoordinateParams ProtocolId = ProtocolFunctionJOGBase + 1 // Base + 1
	ProtocolJOGCommonParams     ProtocolId = ProtocolFunctionJOGBase + 2 // Base + 2
	ProtocolJOGCmd              ProtocolId = ProtocolFunctionJOGBase + 3 // Base + 3
	ProtocolJOGLParams          ProtocolId = ProtocolFunctionJOGBase + 4 // Base + 4
)

const (
	// PTP
	ProtocolFunctionPTPBase     ProtocolId = 80
	ProtocolPTPJointParams      ProtocolId = ProtocolFunctionPTPBase + 0 // Base + 0
	ProtocolPTPCoordinateParams ProtocolId = ProtocolFunctionPTPBase + 1 // Base + 1
	ProtocolPTPJumpParams       ProtocolId = ProtocolFunctionPTPBase + 2 // Base + 2
	ProtocolPTPCommonParams     ProtocolId = ProtocolFunctionPTPBase + 3 // Base + 3
	ProtocolPTPCmd              ProtocolId = ProtocolFunctionPTPBase + 4 // Base + 4
	ProtocolPTPLParams          ProtocolId = ProtocolFunctionPTPBase + 5 // Base + 5
	ProtocolPTPWithLCmd         ProtocolId = ProtocolFunctionPTPBase + 6 // Base + 6
	ProtocolPTPJump2Params      ProtocolId = ProtocolFunctionPTPBase + 7 // Base + 7
	ProtocolPTPPOCmd            ProtocolId = ProtocolFunctionPTPBase + 8 // Base + 8
	ProtocolPTPPOWithLCmd       ProtocolId = ProtocolFunctionPTPBase + 9 // Base + 9
)

const (
	// CP
	ProtocolFunctionCPBase ProtocolId = 90
	ProtocolCPParams       ProtocolId = ProtocolFunctionCPBase + 0 // Base + 0
	ProtocolCPCmd          ProtocolId = ProtocolFunctionCPBase + 1 // Base + 1
	ProtocolCPLECmd        ProtocolId = ProtocolFunctionCPBase + 2 // Base + 2
	ProtocolCPRHoldEnable  ProtocolId = ProtocolFunctionCPBase + 3 // Base + 3
	ProtocolCPCommonParams ProtocolId = ProtocolFunctionCPBase + 4 // Base + 4
)

const (
	// ARC
	ProtocolFunctionARCBase ProtocolId = 100
	ProtocolARCParams       ProtocolId = ProtocolFunctionARCBase + 0 // Base + 0
	ProtocolARCCmd          ProtocolId = ProtocolFunctionARCBase + 1 // Base + 1
	ProtocolCircleCmd       ProtocolId = ProtocolFunctionARCBase + 2 // Base + 2
	ProtocolARCCommonParams ProtocolId = ProtocolFunctionARCBase + 3 // Base + 3
)

const (
	// WAIT
	ProtocolFunctionWAITBase ProtocolId = 110
	ProtocolWAITCmd          ProtocolId = ProtocolFunctionWAITBase + 0 // Base + 0
)

const (
	// TRIG
	ProtocolFunctionTRIGBase ProtocolId = 120
	ProtocolTRIGCmd          ProtocolId = ProtocolFunctionTRIGBase + 0 // Base + 0
)

const (
	// EIO
	ProtocolFunctionEIOBase ProtocolId = 130
	ProtocolIOMultiplexing  ProtocolId = ProtocolFunctionEIOBase + 0 // Base + 0
	ProtocolIODO            ProtocolId = ProtocolFunctionEIOBase + 1 // Base + 1
	ProtocolIOPWM           ProtocolId = ProtocolFunctionEIOBase + 2 // Base + 2
	ProtocolIODI            ProtocolId = ProtocolFunctionEIOBase + 3 // Base + 3
	ProtocolIOADC           ProtocolId = ProtocolFunctionEIOBase + 4 // Base + 4
	ProtocolEMotor          ProtocolId = ProtocolFunctionEIOBase + 5 // Base + 5
	ProtocolEMotorS         ProtocolId = ProtocolFunctionEIOBase + 6 // Base + 6
	ProtocolColorSensor     ProtocolId = ProtocolFunctionEIOBase + 7 // Base + 7
	ProtocolIRSwitch        ProtocolId = ProtocolFunctionEIOBase + 8 // Base + 8
)

const (
	// CAL
	ProtocolFunctionCALBase        ProtocolId = 140
	ProtocolAngleSensorStaticError ProtocolId = ProtocolFunctionCALBase + 0 // Base + 0
	ProtocolAngleSensorCoef        ProtocolId = ProtocolFunctionCALBase + 1 // Base + 1
	ProtocolBaseDecoderStaticError ProtocolId = ProtocolFunctionCALBase + 2 // Base + 2
	ProtocolLRHandCalibrateValue   ProtocolId = ProtocolFunctionCALBase + 3 // Base + 3
)

const (
	// WIFI
	ProtocolFunctionWIFIBase  ProtocolId = 150
	ProtocolWIFIConfigMode    ProtocolId = ProtocolFunctionWIFIBase + 0 // Base + 0
	ProtocolWIFISSID          ProtocolId = ProtocolFunctionWIFIBase + 1 // Base + 1
	ProtocolWIFIPassword      ProtocolId = ProtocolFunctionWIFIBase + 2 // Base + 2
	ProtocolWIFIIPAddress     ProtocolId = ProtocolFunctionWIFIBase + 3 // Base + 3
	ProtocolWIFINetmask       ProtocolId = ProtocolFunctionWIFIBase + 4 // Base + 4
	ProtocolWIFIGateway       ProtocolId = ProtocolFunctionWIFIBase + 5 // Base + 5
	ProtocolWIFIDNS           ProtocolId = ProtocolFunctionWIFIBase + 6 // Base + 6
	ProtocolWIFIConnectStatus ProtocolId = ProtocolFunctionWIFIBase + 7 // Base + 7
)

const (
	// Firmware
	ProtocolFunctionFirmware ProtocolId = 160
	ProtocolFirmwareSwitch   ProtocolId = ProtocolFunctionFirmware + 0 // Base + 0
	ProtocolFirmwareMode     ProtocolId = ProtocolFunctionFirmware + 1 // Base + 1
)

const (
	// LostStep
	ProtocolFunctionLostStepBase ProtocolId = 170
	ProtocolLostStepSet          ProtocolId = ProtocolFunctionLostStepBase + 0 // Base + 0
	ProtocolLostStepDetect       ProtocolId = ProtocolFunctionLostStepBase + 1 // Base + 1
)

const (
	// UART4 Peripherals
	ProtocolFunctionCheckModelBase     ProtocolId = 180
	ProtocolCheckUART4PeripheralsModel ProtocolId = ProtocolFunctionCheckModelBase + 0 // Base + 0
	ProtocolUART4PeripheralsEnabled    ProtocolId = ProtocolFunctionCheckModelBase + 1 // Base + 1
)

const (
	// Pulse Mode
	ProtocolFunctionPulseModeBase ProtocolId = 190
	ProtocolFunctionPulseMode     ProtocolId = ProtocolFunctionPulseModeBase + 0 // Base + 0
)

const (
	// TEST
	ProtocolTESTBase   ProtocolId = 220
	ProtocolUserParams ProtocolId = ProtocolTESTBase + 0 // Base + 0
	ProtocolPTPTime    ProtocolId = ProtocolTESTBase + 1 // Base + 1
)

const (
	// QueuedCmd
	ProtocolFunctionQueuedCmdBase  ProtocolId = 240
	ProtocolQueuedCmdStartExec     ProtocolId = ProtocolFunctionQueuedCmdBase + 0 // Base + 0
	ProtocolQueuedCmdStopExec      ProtocolId = ProtocolFunctionQueuedCmdBase + 1 // Base + 1
	ProtocolQueuedCmdForceStopExec ProtocolId = ProtocolFunctionQueuedCmdBase + 2 // Base + 2
	ProtocolQueuedCmdStartDownload ProtocolId = ProtocolFunctionQueuedCmdBase + 3 // Base + 3
	ProtocolQueuedCmdStopDownload  ProtocolId = ProtocolFunctionQueuedCmdBase + 4 // Base + 4
	ProtocolQueuedCmdClear         ProtocolId = ProtocolFunctionQueuedCmdBase + 5 // Base + 5
	ProtocolQueuedCmdCurrentIndex  ProtocolId = ProtocolFunctionQueuedCmdBase + 6 // Base + 6
	ProtocolQueuedCmdLeftSpace     ProtocolId = ProtocolFunctionQueuedCmdBase + 7 // Base + 7
	ProtocolQueuedCmdMotionFinish  ProtocolId = ProtocolFunctionQueuedCmdBase + 8 // Base + 8
)

const (
// ProtocolMax ProtocolID = 256
)
