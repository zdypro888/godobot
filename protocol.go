package godobot

// ProtocolId 定义了所有协议命令的ID
type ProtocolId uint8

const (
	// Device information
	ProtocolFunctionDeviceInfoBase ProtocolId = iota
	ProtocolDeviceSN                          // Base + 0
	ProtocolDeviceName                        // Base + 1
	ProtocolDeviceVersion                     // Base + 2
	ProtocolDeviceWithL                       // Base + 3
	ProtocolDeviceTime                        // Base + 4
	_                                         // Base + 5 (跳过)
	ProtocolDeviceInfo                        // Base + 6
)

const (
	// Pose
	ProtocolFunctionPoseBase ProtocolId = 10 + iota
	ProtocolGetPose                     // Base + 0
	ProtocolResetPose                   // Base + 1
	ProtocolGetKinematics               // Base + 2
	ProtocolGetPoseL                    // Base + 3
)

const (
	// Alarm
	ProtocolFunctionALARMBase ProtocolId = 20 + iota
	ProtocolAlarmsState                  // Base + 0
)

const (
	// HOME
	ProtocolFunctionHOMEBase ProtocolId = 30 + iota
	ProtocolHOMEParams                  // Base + 0
	ProtocolHOMECmd                     // Base + 1
	ProtocolAutoLeveling                // Base + 2
)

const (
	// HHT
	ProtocolFunctionHHTBase      ProtocolId = 40 + iota
	ProtocolHHTTrigMode                     // Base + 0
	ProtocolHHTTrigOutputEnabled            // Base + 1
	ProtocolHHTTrigOutput                   // Base + 2
)

const (
	// Arm Orientation
	ProtocolFunctionArmOrientationBase ProtocolId = 50 + iota
	ProtocolArmOrientation                        // Base + 0
)

const (
	// End effector
	ProtocolFunctionEndEffectorBase ProtocolId = 60 + iota
	ProtocolEndEffectorParams                  // Base + 0
	ProtocolEndEffectorLaser                   // Base + 1
	ProtocolEndEffectorSuctionCup              // Base + 2
	ProtocolEndEffectorGripper                 // Base + 3
)

const (
	// JOG
	ProtocolFunctionJOGBase     ProtocolId = 70 + iota
	ProtocolJOGJointParams                 // Base + 0
	ProtocolJOGCoordinateParams            // Base + 1
	ProtocolJOGCommonParams                // Base + 2
	ProtocolJOGCmd                         // Base + 3
	ProtocolJOGLParams                     // Base + 4
)

const (
	// PTP
	ProtocolFunctionPTPBase     ProtocolId = 80 + iota
	ProtocolPTPJointParams                 // Base + 0
	ProtocolPTPCoordinateParams            // Base + 1
	ProtocolPTPJumpParams                  // Base + 2
	ProtocolPTPCommonParams                // Base + 3
	ProtocolPTPCmd                         // Base + 4
	ProtocolPTPLParams                     // Base + 5
	ProtocolPTPWithLCmd                    // Base + 6
	ProtocolPTPJump2Params                 // Base + 7
	ProtocolPTPPOCmd                       // Base + 8
	ProtocolPTPPOWithLCmd                  // Base + 9
)

const (
	// CP
	ProtocolFunctionCPBase ProtocolId = 90 + iota
	ProtocolCPParams                  // Base + 0
	ProtocolCPCmd                     // Base + 1
	ProtocolCPLECmd                   // Base + 2
	ProtocolCPRHoldEnable             // Base + 3
	ProtocolCPCommonParams            // Base + 4
)

const (
	// ARC
	ProtocolFunctionARCBase ProtocolId = 100 + iota
	ProtocolARCParams                  // Base + 0
	ProtocolARCCmd                     // Base + 1
	ProtocolCircleCmd                  // Base + 2
	ProtocolARCCommonParams            // Base + 3
)

const (
	// WAIT
	ProtocolFunctionWAITBase ProtocolId = 110 + iota
	ProtocolWAITCmd                     // Base + 0
)

const (
	// TRIG
	ProtocolFunctionTRIGBase ProtocolId = 120 + iota
	ProtocolTRIGCmd                     // Base + 0
)

const (
	// EIO
	ProtocolFunctionEIOBase ProtocolId = 130 + iota
	ProtocolIOMultiplexing             // Base + 0
	ProtocolIODO                       // Base + 1
	ProtocolIOPWM                      // Base + 2
	ProtocolIODI                       // Base + 3
	ProtocolIOADC                      // Base + 4
	ProtocolEMotor                     // Base + 5
	ProtocolEMotorS                    // Base + 6
	ProtocolColorSensor                // Base + 7
	ProtocolIRSwitch                   // Base + 8
)

const (
	// CAL
	ProtocolFunctionCALBase        ProtocolId = 140 + iota
	ProtocolAngleSensorStaticError            // Base + 0
	ProtocolAngleSensorCoef                   // Base + 1
	ProtocolBaseDecoderStaticError            // Base + 2
	ProtocolLRHandCalibrateValue              // Base + 3
)

const (
	// WIFI
	ProtocolFunctionWIFIBase  ProtocolId = 150 + iota
	ProtocolWIFIConfigMode               // Base + 0
	ProtocolWIFISSID                     // Base + 1
	ProtocolWIFIPassword                 // Base + 2
	ProtocolWIFIIPAddress                // Base + 3
	ProtocolWIFINetmask                  // Base + 4
	ProtocolWIFIGateway                  // Base + 5
	ProtocolWIFIDNS                      // Base + 6
	ProtocolWIFIConnectStatus            // Base + 7
)

const (
	// Firmware
	ProtocolFunctionFirmware ProtocolId = 160 + iota
	ProtocolFirmwareSwitch              // Base + 0
	ProtocolFirmwareMode                // Base + 1
)

const (
	// LostStep
	ProtocolFunctionLostStepBase ProtocolId = 170 + iota
	ProtocolLostStepSet                     // Base + 0
	ProtocolLostStepDetect                  // Base + 1
)

const (
	// UART4 Peripherals
	ProtocolFunctionCheckModelBase     ProtocolId = 180 + iota
	ProtocolCheckUART4PeripheralsModel            // Base + 1
	ProtocolUART4PeripheralsEnabled               // Base + 2
)

const (
	// Pulse Mode
	ProtocolFunctionPulseModeBase ProtocolId = 190 + iota
	ProtocolFunctionPulseMode                // Base + 1
)

const (
	// TEST
	ProtocolTESTBase   ProtocolId = 220 + iota
	ProtocolUserParams            // Base + 0
	ProtocolPTPTime               // Base + 1
)

const (
	// QueuedCmd
	ProtocolFunctionQueuedCmdBase  ProtocolId = 240 + iota
	ProtocolQueuedCmdStartExec                // Base + 0
	ProtocolQueuedCmdStopExec                 // Base + 1
	ProtocolQueuedCmdForceStopExec            // Base + 2
	ProtocolQueuedCmdStartDownload            // Base + 3
	ProtocolQueuedCmdStopDownload             // Base + 4
	ProtocolQueuedCmdClear                    // Base + 5
	ProtocolQueuedCmdCurrentIndex             // Base + 6
	ProtocolQueuedCmdLeftSpace                // Base + 7
	ProtocolQueuedCmdMotionFinish             // Base + 8
)

const (
// ProtocolMax ProtocolID = 256
)
