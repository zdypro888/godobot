package godobot

import (
	"context"
	"encoding/binary"
	"errors"
	"math"
)

// Dobot 机械臂控制结构
type Dobot struct {
	connector *Connector
}

// NewDobot 创建新的Dobot实例
func NewDobot() *Dobot {
	return &Dobot{connector: &Connector{}}
}

// ConnectDobot 连接到Dobot设备
func (dobot *Dobot) Connect(portName string, baudrate uint32) error {
	ctx := context.Background()
	err := dobot.connector.Open(ctx, portName, baudrate)
	if err != nil {
		return err
	}
	return nil
}

// DobotExec 执行Dobot命令
func (dobot *Dobot) DobotExec() error {
	message := &Message{
		Id:       ProtocolQueuedCmdStartExec,
		RW:       true,
		IsQueued: false,
	}
	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// DisconnectDobot 断开与Dobot设备的连接
func (dobot *Dobot) DisconnectDobot() error {
	if dobot.connector.port != nil {
		return dobot.connector.port.Close()
	}
	return nil
}

// GetMarlinVersion 获取Marlin版本
func (dobot *Dobot) GetMarlinVersion() (uint8, error) {
	message := &Message{
		Id:       ProtocolDeviceVersion,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}
	if len(resp.Params) < 1 {
		return 0, errors.New("invalid response")
	}
	return resp.Params[0], nil
}

// SetCmdTimeout 设置命令超时时间
func (dobot *Dobot) SetCmdTimeout(timeout uint32) error {
	message := &Message{
		Id:       ProtocolId(0xFE), // 特殊命令
		RW:       true,
		IsQueued: false,
		Params:   make([]byte, 4),
	}
	binary.LittleEndian.PutUint32(message.Params, timeout)
	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// SetDeviceSN 设置设备序列号
func (dobot *Dobot) SetDeviceSN(sn string) error {
	if sn == "" {
		return errors.New("invalid params: empty sn")
	}
	message := &Message{
		Id:       ProtocolDeviceSN,
		RW:       true,
		IsQueued: false,
		Params:   []byte(sn),
	}
	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// GetDeviceSN 获取设备序列号
func (dobot *Dobot) GetDeviceSN() (string, error) {
	message := &Message{
		Id:       ProtocolDeviceSN,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return "", err
	}
	return string(resp.Params), nil
}

// SetDeviceName 设置设备名称
func (dobot *Dobot) SetDeviceName(name string) error {
	if name == "" {
		return errors.New("invalid params: empty name")
	}
	message := &Message{
		Id:       ProtocolDeviceName,
		RW:       true,
		IsQueued: false,
		Params:   []byte(name),
	}
	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// GetDeviceName 获取设备名称
func (dobot *Dobot) GetDeviceName() (string, error) {
	message := &Message{
		Id:       ProtocolDeviceName,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return "", err
	}
	return string(resp.Params), nil
}

// GetDeviceVersion 获取设备版本信息
func (dobot *Dobot) GetDeviceVersion() (majorVersion, minorVersion, revision, hwVersion uint8, err error) {
	message := &Message{
		Id:       ProtocolDeviceVersion,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	if len(resp.Params) < 4 {
		return 0, 0, 0, 0, errors.New("invalid response")
	}
	return resp.Params[0], resp.Params[1], resp.Params[2], resp.Params[3], nil
}

// SetDeviceWithL 设置设备L轴
func (dobot *Dobot) SetDeviceWithL(isWithL bool, version uint8) (queuedCmdIndex uint64, err error) {
	message := &Message{
		Id:       ProtocolDeviceWithL,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 2),
	}
	if isWithL {
		message.Params[0] = 1
	}
	message.Params[1] = version
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}
	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// GetDeviceWithL 获取设备L轴状态
func (dobot *Dobot) GetDeviceWithL() (bool, error) {
	message := &Message{
		Id:       ProtocolDeviceWithL,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return false, err
	}
	if len(resp.Params) < 1 {
		return false, errors.New("invalid response")
	}
	return resp.Params[0] != 0, nil
}

// GetDeviceTime 获取设备运行时间
func (dobot *Dobot) GetDeviceTime() (uint32, error) {
	message := &Message{
		Id:       ProtocolDeviceTime,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}
	if len(resp.Params) < 4 {
		return 0, errors.New("invalid response")
	}
	return binary.LittleEndian.Uint32(resp.Params), nil
}

// GetDeviceInfo 获取设备信息
func (dobot *Dobot) GetDeviceInfo() (*DeviceCountInfo, error) {
	message := &Message{
		Id:       ProtocolDeviceInfo,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return nil, err
	}
	if len(resp.Params) < 16 {
		return nil, errors.New("invalid response")
	}
	info := &DeviceCountInfo{
		DeviceRunTime:  binary.LittleEndian.Uint64(resp.Params[0:8]),
		DevicePowerOn:  binary.LittleEndian.Uint32(resp.Params[8:12]),
		DevicePowerOff: binary.LittleEndian.Uint32(resp.Params[12:16]),
	}
	return info, nil
}

// GetPose 获取当前位姿
func (dobot *Dobot) GetPose() (*Pose, error) {
	message := &Message{
		Id:       ProtocolGetPose,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return nil, err
	}
	if len(resp.Params) < 32 {
		return nil, errors.New("invalid response")
	}
	pose := &Pose{
		X: float32(binary.LittleEndian.Uint32(resp.Params[0:4])),
		Y: float32(binary.LittleEndian.Uint32(resp.Params[4:8])),
		Z: float32(binary.LittleEndian.Uint32(resp.Params[8:12])),
		R: float32(binary.LittleEndian.Uint32(resp.Params[12:16])),
	}
	for i := 0; i < 4; i++ {
		pose.JointAngle[i] = float32(binary.LittleEndian.Uint32(resp.Params[16+i*4 : 20+i*4]))
	}
	return pose, nil
}

// ResetPose 重置位姿
func (dobot *Dobot) ResetPose(manual bool, rearArmAngle, frontArmAngle float32) error {
	message := &Message{
		Id:       ProtocolResetPose,
		RW:       true,
		IsQueued: false,
		Params:   make([]byte, 9),
	}
	if manual {
		message.Params[0] = 1
	}
	binary.LittleEndian.PutUint32(message.Params[1:5], math.Float32bits(rearArmAngle))
	binary.LittleEndian.PutUint32(message.Params[5:9], math.Float32bits(frontArmAngle))
	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// GetKinematics 获取运动学参数
func (dobot *Dobot) GetKinematics() (*Kinematics, error) {
	message := &Message{
		Id: ProtocolGetKinematics,
		RW: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return nil, err
	}
	if len(resp.Params) < 8 {
		return nil, errors.New("invalid response")
	}
	kinematics := &Kinematics{
		Velocity:     float32(binary.LittleEndian.Uint32(resp.Params[0:4])),
		Acceleration: float32(binary.LittleEndian.Uint32(resp.Params[4:8])),
	}
	return kinematics, nil
}

// GetPoseL 获取L轴位置
func (dobot *Dobot) GetPoseL() (float32, error) {
	message := &Message{
		Id: ProtocolGetPoseL,
		RW: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}
	if len(resp.Params) < 4 {
		return 0, errors.New("invalid response")
	}
	return float32(binary.LittleEndian.Uint32(resp.Params)), nil
}

// GetAlarmsState 获取报警状态
func (dobot *Dobot) GetAlarmsState() ([]uint8, error) {
	message := &Message{
		Id: ProtocolAlarmsState,
		RW: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return nil, err
	}
	return resp.Params, nil
}

// ClearAllAlarmsState 清除所有报警状态
func (dobot *Dobot) ClearAllAlarmsState() error {
	message := &Message{
		Id: ProtocolAlarmsState,
		RW: true,
	}
	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// SetHOMEParams 设置HOME参数
func (dobot *Dobot) SetHOMEParams(params *HOMEParams) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}
	message := &Message{
		Id:       ProtocolHOMEParams,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 16),
	}

	binary.LittleEndian.PutUint32(message.Params[0:4], math.Float32bits(params.X))
	binary.LittleEndian.PutUint32(message.Params[4:8], math.Float32bits(params.Y))
	binary.LittleEndian.PutUint32(message.Params[8:12], math.Float32bits(params.Z))
	binary.LittleEndian.PutUint32(message.Params[12:16], math.Float32bits(params.R))

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}
	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// GetHOMEParams 获取HOME参数
func (dobot *Dobot) GetHOMEParams() (*HOMEParams, error) {
	message := &Message{
		Id:       ProtocolHOMEParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return nil, err
	}
	if len(resp.Params) < 16 {
		return nil, errors.New("invalid response")
	}

	params := &HOMEParams{
		X: math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[0:4])),
		Y: math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[4:8])),
		Z: math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[8:12])),
		R: math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[12:16])),
	}
	return params, nil
}

// SetHOMECmd 设置HOME命令
func (dobot *Dobot) SetHOMECmd(cmd *HOMECmd) (queuedCmdIndex uint64, err error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}

	message := &Message{
		Id:       ProtocolHOMECmd,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 4),
	}
	binary.LittleEndian.PutUint32(message.Params[0:4], cmd.Reserved)

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}
	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// SetAutoLevelingCmd 设置自动调平命令
func (dobot *Dobot) SetAutoLevelingCmd(cmd *AutoLevelingCmd) (queuedCmdIndex uint64, err error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}
	message := &Message{
		Id:       ProtocolAutoLeveling,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 5),
	}
	message.Params[0] = cmd.ControlFlag
	binary.LittleEndian.PutUint32(message.Params[1:5], uint32(cmd.Precision))

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}
	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// GetAutoLevelingResult 获取自动调平结果
func (dobot *Dobot) GetAutoLevelingResult() (float32, error) {
	message := &Message{
		Id:       ProtocolAutoLeveling,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}
	if len(resp.Params) < 4 {
		return 0, errors.New("invalid response")
	}
	return math.Float32frombits(binary.LittleEndian.Uint32(resp.Params)), nil
}

// SetHHTTrigMode 设置手持示教触发模式
func (dobot *Dobot) SetHHTTrigMode(mode HHTTrigMode) error {
	message := &Message{
		Id:       ProtocolHHTTrigMode,
		RW:       true,
		IsQueued: false,
		Params:   []byte{uint8(mode)},
	}
	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// GetHHTTrigMode 获取手持示教触发模式
func (dobot *Dobot) GetHHTTrigMode() (HHTTrigMode, error) {
	message := &Message{
		Id:       ProtocolHHTTrigMode,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}
	if len(resp.Params) < 1 {
		return 0, errors.New("invalid response")
	}
	return HHTTrigMode(resp.Params[0]), nil
}

// SetHHTTrigOutputEnabled 设置手持示教触发输出使能
func (dobot *Dobot) SetHHTTrigOutputEnabled(enabled bool) error {
	message := &Message{
		Id:       ProtocolHHTTrigOutputEnabled,
		RW:       true,
		IsQueued: false,
		Params:   make([]byte, 1),
	}
	if enabled {
		message.Params[0] = 1
	}
	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// GetHHTTrigOutputEnabled 获取手持示教触发输出使能状态
func (dobot *Dobot) GetHHTTrigOutputEnabled() (bool, error) {
	message := &Message{
		Id:       ProtocolHHTTrigOutputEnabled,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return false, err
	}
	if len(resp.Params) < 1 {
		return false, errors.New("invalid response")
	}
	return resp.Params[0] != 0, nil
}

// GetHHTTrigOutput 获取手持示教触发输出
func (dobot *Dobot) GetHHTTrigOutput() (bool, error) {
	message := &Message{
		Id:       ProtocolHHTTrigOutput,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return false, err
	}
	if len(resp.Params) < 1 {
		return false, errors.New("invalid response")
	}
	return resp.Params[0] != 0, nil
}

// SetEndEffectorParams 设置末端执行器参数
func (dobot *Dobot) SetEndEffectorParams(params *EndEffectorParams) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolEndEffectorParams,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 12),
	}

	binary.LittleEndian.PutUint32(message.Params[0:4], math.Float32bits(params.XBias))
	binary.LittleEndian.PutUint32(message.Params[4:8], math.Float32bits(params.YBias))
	binary.LittleEndian.PutUint32(message.Params[8:12], math.Float32bits(params.ZBias))

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// GetEndEffectorParams 获取末端执行器参数
func (dobot *Dobot) GetEndEffectorParams() (*EndEffectorParams, error) {
	message := &Message{
		Id:       ProtocolEndEffectorParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return nil, err
	}
	if len(resp.Params) < 12 {
		return nil, errors.New("invalid response")
	}

	params := &EndEffectorParams{
		XBias: math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[0:4])),
		YBias: math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[4:8])),
		ZBias: math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[8:12])),
	}
	return params, nil
}

// SetEndEffectorLaser 设置末端执行器激光
func (dobot *Dobot) SetEndEffectorLaser(enableCtrl bool, on bool) (queuedCmdIndex uint64, err error) {
	message := &Message{
		Id:       ProtocolEndEffectorLaser,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 2),
	}
	if enableCtrl {
		message.Params[0] = 1
	}
	if on {
		message.Params[1] = 1
	}

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// GetEndEffectorLaser 获取末端执行器激光状态
func (dobot *Dobot) GetEndEffectorLaser() (isCtrlEnabled bool, isOn bool, err error) {
	message := &Message{
		Id:       ProtocolEndEffectorLaser,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return false, false, err
	}
	if len(resp.Params) < 2 {
		return false, false, errors.New("invalid response")
	}
	return resp.Params[0] != 0, resp.Params[1] != 0, nil
}

// SetEndEffectorSuctionCup 设置末端执行器吸盘
func (dobot *Dobot) SetEndEffectorSuctionCup(enableCtrl bool, suck bool) (queuedCmdIndex uint64, err error) {
	message := &Message{
		Id:       ProtocolEndEffectorSuctionCup,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 2),
	}
	if enableCtrl {
		message.Params[0] = 1
	}
	if suck {
		message.Params[1] = 1
	}

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// GetEndEffectorSuctionCup 获取末端执行器吸盘状态
func (dobot *Dobot) GetEndEffectorSuctionCup() (isCtrlEnabled bool, isSucked bool, err error) {
	message := &Message{
		Id:       ProtocolEndEffectorSuctionCup,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return false, false, err
	}
	if len(resp.Params) < 2 {
		return false, false, errors.New("invalid response")
	}
	return resp.Params[0] != 0, resp.Params[1] != 0, nil
}

// SetEndEffectorGripper 设置末端执行器夹爪
func (dobot *Dobot) SetEndEffectorGripper(enableCtrl bool, grip bool) (queuedCmdIndex uint64, err error) {
	message := &Message{
		Id:       ProtocolEndEffectorGripper,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 2),
	}
	if enableCtrl {
		message.Params[0] = 1
	}
	if grip {
		message.Params[1] = 1
	}

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// GetEndEffectorGripper 获取末端执行器夹爪状态
func (dobot *Dobot) GetEndEffectorGripper() (isCtrlEnabled bool, isGripped bool, err error) {
	message := &Message{
		Id:       ProtocolEndEffectorGripper,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return false, false, err
	}
	if len(resp.Params) < 2 {
		return false, false, errors.New("invalid response")
	}
	return resp.Params[0] != 0, resp.Params[1] != 0, nil
}

// SetArmOrientation 设置机械臂方向
func (dobot *Dobot) SetArmOrientation(armOrientation ArmOrientation) (queuedCmdIndex uint64, err error) {
	message := &Message{
		Id:       ProtocolArmOrientation,
		RW:       true,
		IsQueued: true,
		Params:   []byte{uint8(armOrientation)},
	}

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// GetArmOrientation 获取机械臂方向
func (dobot *Dobot) GetArmOrientation() (ArmOrientation, error) {
	message := &Message{
		Id:       ProtocolArmOrientation,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}
	if len(resp.Params) < 1 {
		return 0, errors.New("invalid response")
	}
	return ArmOrientation(resp.Params[0]), nil
}

// SetJOGJointParams 设置JOG关节参数
func (dobot *Dobot) SetJOGJointParams(params *JOGJointParams) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolJOGJointParams,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 32),
	}

	for i := 0; i < 4; i++ {
		binary.LittleEndian.PutUint32(message.Params[i*4:i*4+4], math.Float32bits(params.Velocity[i]))
		binary.LittleEndian.PutUint32(message.Params[16+i*4:16+i*4+4], math.Float32bits(params.Acceleration[i]))
	}

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// GetJOGJointParams 获取JOG关节参数
func (dobot *Dobot) GetJOGJointParams() (*JOGJointParams, error) {
	message := &Message{
		Id:       ProtocolJOGJointParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return nil, err
	}
	if len(resp.Params) < 32 {
		return nil, errors.New("invalid response")
	}

	params := &JOGJointParams{}
	for i := 0; i < 4; i++ {
		params.Velocity[i] = math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[i*4 : i*4+4]))
		params.Acceleration[i] = math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[16+i*4 : 16+i*4+4]))
	}
	return params, nil
}

// SetJOGCoordinateParams 设置JOG坐标参数
func (dobot *Dobot) SetJOGCoordinateParams(params *JOGCoordinateParams) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolJOGCoordinateParams,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 32),
	}

	for i := 0; i < 4; i++ {
		binary.LittleEndian.PutUint32(message.Params[i*4:i*4+4], math.Float32bits(params.Velocity[i]))
		binary.LittleEndian.PutUint32(message.Params[16+i*4:16+i*4+4], math.Float32bits(params.Acceleration[i]))
	}

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// GetJOGCoordinateParams 获取JOG坐标参数
func (dobot *Dobot) GetJOGCoordinateParams() (*JOGCoordinateParams, error) {
	message := &Message{
		Id:       ProtocolJOGCoordinateParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return nil, err
	}
	if len(resp.Params) < 32 {
		return nil, errors.New("invalid response")
	}

	params := &JOGCoordinateParams{}
	for i := 0; i < 4; i++ {
		params.Velocity[i] = math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[i*4 : i*4+4]))
		params.Acceleration[i] = math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[16+i*4 : 16+i*4+4]))
	}
	return params, nil
}

// SetJOGLParams 设置JOGL参数
func (dobot *Dobot) SetJOGLParams(params *JOGLParams) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolJOGLParams,
		RW:       true,
		IsQueued: false, // C++实现中强制为false
		Params:   make([]byte, 8),
	}

	binary.LittleEndian.PutUint32(message.Params[0:4], math.Float32bits(params.Velocity))
	binary.LittleEndian.PutUint32(message.Params[4:8], math.Float32bits(params.Acceleration))

	_, err = dobot.connector.SendMessage(context.Background(), message)
	return 0, err // 由于IsQueued为false，所以不返回queuedCmdIndex
}

func (dobot *Dobot) GetJOGLParams() (*JOGLParams, error) {
	message := &Message{
		Id:       ProtocolJOGLParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return nil, err
	}
	if len(resp.Params) < 8 {
		return nil, errors.New("invalid response")
	}

	params := &JOGLParams{
		Velocity:     math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[0:4])),
		Acceleration: math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[4:8])),
	}
	return params, nil
}

// SetJOGCommonParams 设置JOG通用参数
func (dobot *Dobot) SetJOGCommonParams(params *JOGCommonParams) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolJOGCommonParams,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 8),
	}

	binary.LittleEndian.PutUint32(message.Params[0:4], math.Float32bits(params.VelocityRatio))
	binary.LittleEndian.PutUint32(message.Params[4:8], math.Float32bits(params.AccelerationRatio))

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// GetJOGCommonParams 获取JOG通用参数
func (dobot *Dobot) GetJOGCommonParams() (*JOGCommonParams, error) {
	message := &Message{
		Id:       ProtocolJOGCommonParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return nil, err
	}
	if len(resp.Params) < 8 {
		return nil, errors.New("invalid response")
	}

	params := &JOGCommonParams{
		VelocityRatio:     math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[0:4])),
		AccelerationRatio: math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[4:8])),
	}
	return params, nil
}

// SetJOGCmd 设置JOG命令
func (dobot *Dobot) SetJOGCmd(cmd *JOGCmd) (queuedCmdIndex uint64, err error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}

	message := &Message{
		Id:       ProtocolJOGCmd,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 2),
	}

	message.Params[0] = cmd.IsJoint
	message.Params[1] = cmd.Cmd

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// SetPTPJointParams 设置PTP关节参数
func (dobot *Dobot) SetPTPJointParams(params *PTPJointParams) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolPTPJointParams,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 32),
	}

	for i := 0; i < 4; i++ {
		binary.LittleEndian.PutUint32(message.Params[i*4:i*4+4], math.Float32bits(params.Velocity[i]))
		binary.LittleEndian.PutUint32(message.Params[16+i*4:16+i*4+4], math.Float32bits(params.Acceleration[i]))
	}

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// GetPTPJointParams 获取PTP关节参数
func (dobot *Dobot) GetPTPJointParams() (*PTPJointParams, error) {
	message := &Message{
		Id:       ProtocolPTPJointParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return nil, err
	}
	if len(resp.Params) < 32 {
		return nil, errors.New("invalid response")
	}

	params := &PTPJointParams{}
	for i := 0; i < 4; i++ {
		params.Velocity[i] = math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[i*4 : i*4+4]))
		params.Acceleration[i] = math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[16+i*4 : 16+i*4+4]))
	}
	return params, nil
}

// SetPTPCoordinateParams 设置PTP坐标参数
func (dobot *Dobot) SetPTPCoordinateParams(params *PTPCoordinateParams) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolPTPCoordinateParams,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 16),
	}

	binary.LittleEndian.PutUint32(message.Params[0:4], math.Float32bits(params.XYZVelocity))
	binary.LittleEndian.PutUint32(message.Params[4:8], math.Float32bits(params.RVelocity))
	binary.LittleEndian.PutUint32(message.Params[8:12], math.Float32bits(params.XYZAcceleration))
	binary.LittleEndian.PutUint32(message.Params[12:16], math.Float32bits(params.RAcceleration))

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// GetPTPCoordinateParams 获取PTP坐标参数
func (dobot *Dobot) GetPTPCoordinateParams() (*PTPCoordinateParams, error) {
	message := &Message{
		Id:       ProtocolPTPCoordinateParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return nil, err
	}
	if len(resp.Params) < 16 {
		return nil, errors.New("invalid response")
	}

	params := &PTPCoordinateParams{
		XYZVelocity:     math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[0:4])),
		RVelocity:       math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[4:8])),
		XYZAcceleration: math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[8:12])),
		RAcceleration:   math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[12:16])),
	}
	return params, nil
}

// SetPTPLParams 设置PTPL参数
func (dobot *Dobot) SetPTPLParams(params *PTPLParams) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolPTPLParams,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 8),
	}

	binary.LittleEndian.PutUint32(message.Params[0:4], math.Float32bits(params.Velocity))
	binary.LittleEndian.PutUint32(message.Params[4:8], math.Float32bits(params.Acceleration))

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// GetPTPLParams 获取PTPL参数
func (dobot *Dobot) GetPTPLParams() (*PTPLParams, error) {
	message := &Message{
		Id:       ProtocolPTPLParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return nil, err
	}
	if len(resp.Params) < 8 {
		return nil, errors.New("invalid response")
	}

	params := &PTPLParams{
		Velocity:     math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[0:4])),
		Acceleration: math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[4:8])),
	}
	return params, nil
}

// SetPTPJumpParams 设置PTP跳跃参数
func (dobot *Dobot) SetPTPJumpParams(params *PTPJumpParams) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolPTPJumpParams,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 8),
	}

	binary.LittleEndian.PutUint32(message.Params[0:4], math.Float32bits(params.JumpHeight))
	binary.LittleEndian.PutUint32(message.Params[4:8], math.Float32bits(params.ZLimit))

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// GetPTPJumpParams 获取PTP跳跃参数
func (dobot *Dobot) GetPTPJumpParams() (*PTPJumpParams, error) {
	message := &Message{
		Id:       ProtocolPTPJumpParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return nil, err
	}
	if len(resp.Params) < 8 {
		return nil, errors.New("invalid response")
	}

	params := &PTPJumpParams{
		JumpHeight: math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[0:4])),
		ZLimit:     math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[4:8])),
	}
	return params, nil
}

// SetPTPJump2Params 设置PTP跳跃2参数
func (dobot *Dobot) SetPTPJump2Params(params *PTPJump2Params) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolPTPJump2Params,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 12),
	}

	binary.LittleEndian.PutUint32(message.Params[0:4], math.Float32bits(params.StartJumpHeight))
	binary.LittleEndian.PutUint32(message.Params[4:8], math.Float32bits(params.EndJumpHeight))
	binary.LittleEndian.PutUint32(message.Params[8:12], math.Float32bits(params.ZLimit))

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// GetPTPJump2Params 获取PTP跳跃2参数
func (dobot *Dobot) GetPTPJump2Params() (*PTPJump2Params, error) {
	message := &Message{
		Id:       ProtocolPTPJump2Params,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return nil, err
	}
	if len(resp.Params) < 12 {
		return nil, errors.New("invalid response")
	}

	params := &PTPJump2Params{
		StartJumpHeight: math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[0:4])),
		EndJumpHeight:   math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[4:8])),
		ZLimit:          math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[8:12])),
	}
	return params, nil
}

// SetPTPCommonParams 设置PTP通用参数
func (dobot *Dobot) SetPTPCommonParams(params *PTPCommonParams) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolPTPCommonParams,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 8),
	}

	binary.LittleEndian.PutUint32(message.Params[0:4], math.Float32bits(params.VelocityRatio))
	binary.LittleEndian.PutUint32(message.Params[4:8], math.Float32bits(params.AccelerationRatio))

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

func (dobot *Dobot) GetPTPCommonParams() (*PTPCommonParams, error) {
	message := &Message{
		Id:       ProtocolPTPCommonParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return nil, err
	}
	if len(resp.Params) < 8 {
		return nil, errors.New("invalid response")
	}

	params := &PTPCommonParams{
		VelocityRatio:     math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[0:4])),
		AccelerationRatio: math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[4:8])),
	}
	return params, nil
}

// SetPTPCmd 设置PTP命令
func (dobot *Dobot) SetPTPCmd(cmd *PTPCmd) (queuedCmdIndex uint64, err error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}

	message := &Message{
		Id:       ProtocolPTPCmd,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 20),
	}

	message.Params[0] = uint8(cmd.PTPMode)
	binary.LittleEndian.PutUint32(message.Params[4:8], math.Float32bits(cmd.X))
	binary.LittleEndian.PutUint32(message.Params[8:12], math.Float32bits(cmd.Y))
	binary.LittleEndian.PutUint32(message.Params[12:16], math.Float32bits(cmd.Z))
	binary.LittleEndian.PutUint32(message.Params[16:20], math.Float32bits(cmd.R))

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// SetPTPWithLCmd 设置带L轴的PTP命令
func (dobot *Dobot) SetPTPWithLCmd(cmd *PTPWithLCmd) (queuedCmdIndex uint64, err error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}

	message := &Message{
		Id:       ProtocolPTPWithLCmd,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 24),
	}

	message.Params[0] = uint8(cmd.PTPMode)
	binary.LittleEndian.PutUint32(message.Params[4:8], math.Float32bits(cmd.X))
	binary.LittleEndian.PutUint32(message.Params[8:12], math.Float32bits(cmd.Y))
	binary.LittleEndian.PutUint32(message.Params[12:16], math.Float32bits(cmd.Z))
	binary.LittleEndian.PutUint32(message.Params[16:20], math.Float32bits(cmd.R))
	binary.LittleEndian.PutUint32(message.Params[20:24], math.Float32bits(cmd.L))

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// SetCPParams 设置CP参数
func (dobot *Dobot) SetCPParams(params *CPParams) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolCPParams,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 13),
	}

	binary.LittleEndian.PutUint32(message.Params[0:4], math.Float32bits(params.PlanAcc))
	binary.LittleEndian.PutUint32(message.Params[4:8], math.Float32bits(params.JuncitionVel))
	binary.LittleEndian.PutUint32(message.Params[8:12], math.Float32bits(params.Acc))
	message.Params[12] = params.RealTimeTrack

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// SetCPCmd 设置CP命令
func (dobot *Dobot) SetCPCmd(cmd *CPCmd) (queuedCmdIndex uint64, err error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}

	message := &Message{
		Id:       ProtocolCPCmd,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 17),
	}

	message.Params[0] = cmd.CPMode
	binary.LittleEndian.PutUint32(message.Params[1:5], uint32(cmd.X))
	binary.LittleEndian.PutUint32(message.Params[5:9], uint32(cmd.Y))
	binary.LittleEndian.PutUint32(message.Params[9:13], uint32(cmd.Z))
	binary.LittleEndian.PutUint32(message.Params[13:17], uint32(cmd.Velocity))

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// SetCPLECmd 设置CPLE命令
func (dobot *Dobot) SetCPLECmd(cpMode uint8, x, y, z, power float32) (queuedCmdIndex uint64, err error) {
	message := &Message{
		Id:       ProtocolCPLECmd,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 17),
	}

	message.Params[0] = cpMode
	binary.LittleEndian.PutUint32(message.Params[1:5], math.Float32bits(x))
	binary.LittleEndian.PutUint32(message.Params[5:9], math.Float32bits(y))
	binary.LittleEndian.PutUint32(message.Params[9:13], math.Float32bits(z))
	binary.LittleEndian.PutUint32(message.Params[13:17], math.Float32bits(power))

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// SetCPRHoldEnable 设置CPR保持使能
func (dobot *Dobot) SetCPRHoldEnable(isEnable bool) error {
	message := &Message{
		Id:       ProtocolCPRHoldEnable,
		RW:       true,
		IsQueued: false,
		Params:   make([]byte, 1),
	}
	if isEnable {
		message.Params[0] = 1
	}
	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// GetCPRHoldEnable 获取CPR保持使能状态
func (dobot *Dobot) GetCPRHoldEnable() (bool, error) {
	message := &Message{
		Id:       ProtocolCPRHoldEnable,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return false, err
	}
	if len(resp.Params) < 1 {
		return false, errors.New("invalid response")
	}
	return resp.Params[0] != 0, nil
}

// SetCPCommonParams 设置CP通用参数
func (dobot *Dobot) SetCPCommonParams(params *CPCommonParams) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolCPCommonParams,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 8),
	}

	binary.LittleEndian.PutUint32(message.Params[0:4], uint32(params.VelocityRatio))
	binary.LittleEndian.PutUint32(message.Params[4:8], uint32(params.AccelerationRatio))

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// GetCPCommonParams 获取CP通用参数
func (dobot *Dobot) GetCPCommonParams() (*CPCommonParams, error) {
	message := &Message{
		Id:       ProtocolCPCommonParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return nil, err
	}
	if len(resp.Params) < 8 {
		return nil, errors.New("invalid response")
	}

	params := &CPCommonParams{
		VelocityRatio:     float32(binary.LittleEndian.Uint32(resp.Params[0:4])),
		AccelerationRatio: float32(binary.LittleEndian.Uint32(resp.Params[4:8])),
	}
	return params, nil
}

// SetARCParams 设置ARC参数
func (dobot *Dobot) SetARCParams(params *ARCParams) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolARCParams,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 16),
	}

	binary.LittleEndian.PutUint32(message.Params[0:4], math.Float32bits(params.XYZVelocity))
	binary.LittleEndian.PutUint32(message.Params[4:8], math.Float32bits(params.RVelocity))
	binary.LittleEndian.PutUint32(message.Params[8:12], math.Float32bits(params.XYZAcceleration))
	binary.LittleEndian.PutUint32(message.Params[12:16], math.Float32bits(params.RAcceleration))

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// GetARCParams 获取ARC参数
func (dobot *Dobot) GetARCParams() (*ARCParams, error) {
	message := &Message{
		Id:       ProtocolARCParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return nil, err
	}
	if len(resp.Params) < 16 {
		return nil, errors.New("invalid response")
	}

	params := &ARCParams{
		XYZVelocity:     math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[0:4])),
		RVelocity:       math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[4:8])),
		XYZAcceleration: math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[8:12])),
		RAcceleration:   math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[12:16])),
	}
	return params, nil
}

// SetARCCmd 设置ARC命令
func (dobot *Dobot) SetARCCmd(cmd *ARCCmd) (queuedCmdIndex uint64, err error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}

	message := &Message{
		Id:       ProtocolARCCmd,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 32),
	}

	binary.LittleEndian.PutUint32(message.Params[0:4], uint32(cmd.CirPoint.X))
	binary.LittleEndian.PutUint32(message.Params[4:8], uint32(cmd.CirPoint.Y))
	binary.LittleEndian.PutUint32(message.Params[8:12], uint32(cmd.CirPoint.Z))
	binary.LittleEndian.PutUint32(message.Params[12:16], uint32(cmd.CirPoint.R))
	binary.LittleEndian.PutUint32(message.Params[16:20], uint32(cmd.ToPoint.X))
	binary.LittleEndian.PutUint32(message.Params[20:24], uint32(cmd.ToPoint.Y))
	binary.LittleEndian.PutUint32(message.Params[24:28], uint32(cmd.ToPoint.Z))
	binary.LittleEndian.PutUint32(message.Params[28:32], uint32(cmd.ToPoint.R))

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// SetCircleCmd 设置圆弧命令
func (dobot *Dobot) SetCircleCmd(cmd *CircleCmd) (queuedCmdIndex uint64, err error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}

	message := &Message{
		Id:       ProtocolCircleCmd,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 36),
	}

	binary.LittleEndian.PutUint32(message.Params[0:4], uint32(cmd.CirPoint.X))
	binary.LittleEndian.PutUint32(message.Params[4:8], uint32(cmd.CirPoint.Y))
	binary.LittleEndian.PutUint32(message.Params[8:12], uint32(cmd.CirPoint.Z))
	binary.LittleEndian.PutUint32(message.Params[12:16], uint32(cmd.CirPoint.R))
	binary.LittleEndian.PutUint32(message.Params[16:20], uint32(cmd.ToPoint.X))
	binary.LittleEndian.PutUint32(message.Params[20:24], uint32(cmd.ToPoint.Y))
	binary.LittleEndian.PutUint32(message.Params[24:28], uint32(cmd.ToPoint.Z))
	binary.LittleEndian.PutUint32(message.Params[28:32], uint32(cmd.ToPoint.R))
	binary.LittleEndian.PutUint32(message.Params[32:36], cmd.Count)

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// SetARCCommonParams 设置ARC通用参数
func (dobot *Dobot) SetARCCommonParams(params *ARCCommonParams) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolARCCommonParams,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 8),
	}

	binary.LittleEndian.PutUint32(message.Params[0:4], uint32(params.VelocityRatio))
	binary.LittleEndian.PutUint32(message.Params[4:8], uint32(params.AccelerationRatio))

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// GetARCCommonParams 获取ARC通用参数
func (dobot *Dobot) GetARCCommonParams() (*ARCCommonParams, error) {
	message := &Message{
		Id:       ProtocolARCCommonParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return nil, err
	}
	if len(resp.Params) < 8 {
		return nil, errors.New("invalid response")
	}

	params := &ARCCommonParams{
		VelocityRatio:     float32(binary.LittleEndian.Uint32(resp.Params[0:4])),
		AccelerationRatio: float32(binary.LittleEndian.Uint32(resp.Params[4:8])),
	}
	return params, nil
}

// SetWAITCmd 设置等待命令
func (dobot *Dobot) SetWAITCmd(cmd *WAITCmd) (queuedCmdIndex uint64, err error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}

	message := &Message{
		Id:       ProtocolWAITCmd,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 4),
	}

	binary.LittleEndian.PutUint32(message.Params[0:4], cmd.Timeout)

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// SetTRIGCmd 设置触发命令
func (dobot *Dobot) SetTRIGCmd(cmd *TRIGCmd) (queuedCmdIndex uint64, err error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}

	message := &Message{
		Id:       ProtocolTRIGCmd,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 7),
	}

	message.Params[0] = cmd.Address
	message.Params[1] = cmd.Mode
	message.Params[2] = cmd.Condition
	binary.LittleEndian.PutUint32(message.Params[3:7], uint32(cmd.Threshold))

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// SetIOMultiplexing 设置IO复用
func (dobot *Dobot) SetIOMultiplexing(params *IOMultiplexing) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolIOMultiplexing,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 2),
	}

	message.Params[0] = params.Address
	message.Params[1] = params.Multiplex

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// SetIODO 设置IO数字输出
func (dobot *Dobot) SetIODO(params *IODO) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolIODO,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 2),
	}

	message.Params[0] = params.Address
	message.Params[1] = params.Level

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// SetIOPWM 设置IO PWM输出
func (dobot *Dobot) SetIOPWM(params *IOPWM) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolIOPWM,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 9),
	}

	message.Params[0] = params.Address
	binary.LittleEndian.PutUint32(message.Params[1:5], math.Float32bits(params.Frequency))
	binary.LittleEndian.PutUint32(message.Params[5:9], math.Float32bits(params.DutyCycle))

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// SetEMotor 设置扩展电机
func (dobot *Dobot) SetEMotor(params *EMotor) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolEMotor,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 6),
	}

	message.Params[0] = params.Index
	if params.IsEnabled {
		message.Params[1] = 1
	}
	binary.LittleEndian.PutUint32(message.Params[2:6], uint32(params.Speed))

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// SetEMotorS 设置扩展步进电机
func (dobot *Dobot) SetEMotorS(params *EMotorS) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}
	message := &Message{
		Id:       ProtocolEMotorS,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 10),
	}

	message.Params[0] = params.Index
	if params.IsEnabled {
		message.Params[1] = 1
	}
	binary.LittleEndian.PutUint32(message.Params[2:6], uint32(params.Speed))
	binary.LittleEndian.PutUint32(message.Params[6:10], params.Distance)

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}
	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// SetColorSensor 设置颜色传感器
func (dobot *Dobot) SetColorSensor(enable bool, colorPort ColorPort, version uint8) error {
	message := &Message{
		Id:       ProtocolColorSensor,
		RW:       true,
		IsQueued: false,
		Params:   make([]byte, 3),
	}

	if enable {
		message.Params[0] = 1
	}
	message.Params[1] = uint8(colorPort)
	message.Params[2] = version

	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// GetColorSensor 获取颜色传感器
func (dobot *Dobot) GetColorSensor() (r, g, b uint8, err error) {
	message := &Message{
		Id:       ProtocolColorSensor,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, 0, 0, err
	}
	if len(resp.Params) < 3 {
		return 0, 0, 0, errors.New("invalid response")
	}
	return resp.Params[0], resp.Params[1], resp.Params[2], nil
}

// SetInfraredSensor 设置红外传感器
func (dobot *Dobot) SetInfraredSensor(enable bool, infraredPort InfraredPort, version uint8) error {
	message := &Message{
		Id:       ProtocolInfraredSensor,
		RW:       true,
		IsQueued: false,
		Params:   make([]byte, 3),
	}

	if enable {
		message.Params[0] = 1
	}
	message.Params[1] = uint8(infraredPort)
	message.Params[2] = version

	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// GetInfraredSensor 获取红外传感器
func (dobot *Dobot) GetInfraredSensor(port InfraredPort) (uint8, error) {
	message := &Message{
		Id:       ProtocolInfraredSensor,
		RW:       false,
		IsQueued: false,
		Params:   []byte{uint8(port)},
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}
	if len(resp.Params) < 1 {
		return 0, errors.New("invalid response")
	}
	return resp.Params[0], nil
}

// SetAngleSensorStaticError 设置角度传感器静态误差
func (dobot *Dobot) SetAngleSensorStaticError(rearArmAngleError, frontArmAngleError float32) error {
	message := &Message{
		Id:       ProtocolAngleSensorStaticError,
		RW:       true,
		IsQueued: false,
		Params:   make([]byte, 8),
	}

	binary.LittleEndian.PutUint32(message.Params[0:4], math.Float32bits(rearArmAngleError))
	binary.LittleEndian.PutUint32(message.Params[4:8], math.Float32bits(frontArmAngleError))

	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// GetAngleSensorStaticError 获取角度传感器静态误差
func (dobot *Dobot) GetAngleSensorStaticError() (rearArmAngleError, frontArmAngleError float32, err error) {
	message := &Message{
		Id:       ProtocolAngleSensorStaticError,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, 0, err
	}
	if len(resp.Params) < 8 {
		return 0, 0, errors.New("invalid response")
	}
	return math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[0:4])),
		math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[4:8])), nil
}

// SetAngleSensorCoef 设置角度传感器系数
func (dobot *Dobot) SetAngleSensorCoef(rearArmAngleCoef, frontArmAngleCoef float32) error {
	message := &Message{
		Id:       ProtocolAngleSensorCoef,
		RW:       true,
		IsQueued: false,
		Params:   make([]byte, 8),
	}

	binary.LittleEndian.PutUint32(message.Params[0:4], math.Float32bits(rearArmAngleCoef))
	binary.LittleEndian.PutUint32(message.Params[4:8], math.Float32bits(frontArmAngleCoef))

	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// GetAngleSensorCoef 获取角度传感器系数
func (dobot *Dobot) GetAngleSensorCoef() (rearArmAngleCoef, frontArmAngleCoef float32, err error) {
	message := &Message{
		Id:       ProtocolAngleSensorCoef,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, 0, err
	}
	if len(resp.Params) < 8 {
		return 0, 0, errors.New("invalid response")
	}
	return math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[0:4])),
		math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[4:8])), nil
}

// SetBaseDecoderStaticError 设置基座解码器静态误差
func (dobot *Dobot) SetBaseDecoderStaticError(baseDecoderError float32) error {
	message := &Message{
		Id:       ProtocolBaseDecoderStaticError,
		RW:       true,
		IsQueued: false,
		Params:   make([]byte, 4),
	}

	binary.LittleEndian.PutUint32(message.Params[0:4], math.Float32bits(baseDecoderError))

	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// GetBaseDecoderStaticError 获取基座解码器静态误差
func (dobot *Dobot) GetBaseDecoderStaticError() (float32, error) {
	message := &Message{
		Id:       ProtocolBaseDecoderStaticError,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}
	if len(resp.Params) < 4 {
		return 0, errors.New("invalid response")
	}
	return math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[0:4])), nil
}

// SetLRHandCalibrateValue 设置左右手校准值
func (dobot *Dobot) SetLRHandCalibrateValue(lrHandCalibrateValue float32) error {
	message := &Message{
		Id:       ProtocolLRHandCalibrateValue,
		RW:       true,
		IsQueued: false,
		Params:   make([]byte, 4),
	}

	binary.LittleEndian.PutUint32(message.Params[0:4], math.Float32bits(lrHandCalibrateValue))

	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// GetLRHandCalibrateValue 获取左右手校准值
func (dobot *Dobot) GetLRHandCalibrateValue() (float32, error) {
	message := &Message{
		Id:       ProtocolLRHandCalibrateValue,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}
	if len(resp.Params) < 4 {
		return 0, errors.New("invalid response")
	}
	return math.Float32frombits(binary.LittleEndian.Uint32(resp.Params[0:4])), nil
}

// SetQueuedCmdStopExec 停止执行队列命令
func (dobot *Dobot) SetQueuedCmdStopExec() error {
	message := &Message{
		Id:       ProtocolQueuedCmdStopExec,
		RW:       true,
		IsQueued: false,
	}
	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// SetQueuedCmdForceStopExec 强制停止执行队列命令
func (dobot *Dobot) SetQueuedCmdForceStopExec() error {
	message := &Message{
		Id:       ProtocolQueuedCmdForceStopExec,
		RW:       true,
		IsQueued: false,
	}
	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// SetQueuedCmdStartDownload 开始下载队列命令
func (dobot *Dobot) SetQueuedCmdStartDownload(totalLoop uint32, linePerLoop uint32) error {
	message := &Message{
		Id:       ProtocolQueuedCmdStartDownload,
		RW:       true,
		IsQueued: false,
		Params:   make([]byte, 8),
	}
	binary.LittleEndian.PutUint32(message.Params[0:4], totalLoop)
	binary.LittleEndian.PutUint32(message.Params[4:8], linePerLoop)
	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// SetQueuedCmdStopDownload 停止下载队列命令
func (dobot *Dobot) SetQueuedCmdStopDownload() error {
	message := &Message{
		Id:       ProtocolQueuedCmdStopDownload,
		RW:       true,
		IsQueued: false,
	}
	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// SetQueuedCmdClear 清除队列命令
func (dobot *Dobot) SetQueuedCmdClear() error {
	message := &Message{
		Id:       ProtocolQueuedCmdClear,
		RW:       true,
		IsQueued: false,
	}
	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// GetQueuedCmdCurrentIndex 获取当前队列命令索引
func (dobot *Dobot) GetQueuedCmdCurrentIndex() (uint64, error) {
	message := &Message{
		Id:       ProtocolQueuedCmdCurrentIndex,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}
	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	return binary.LittleEndian.Uint64(resp.Params), nil
}

// GetQueuedCmdMotionFinish 获取队列命令运动是否完成
func (dobot *Dobot) GetQueuedCmdMotionFinish() (bool, error) {
	message := &Message{
		Id:       ProtocolQueuedCmdMotionFinish,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return false, err
	}
	if len(resp.Params) < 1 {
		return false, errors.New("invalid response")
	}
	return resp.Params[0] != 0, nil
}

// SetPTPPOCmd 设置PTP并行输出命令
func (dobot *Dobot) SetPTPPOCmd(ptpCmd *PTPCmd, parallelCmd []ParallelOutputCmd) (queuedCmdIndex uint64, err error) {
	if ptpCmd == nil {
		return 0, errors.New("invalid params: ptpCmd is nil")
	}

	// 计算参数总长度
	paramsLen := 20 // PTPCmd size
	paramsLen += 1  // parallelCmdCount
	if len(parallelCmd) > 0 {
		paramsLen += len(parallelCmd) * 4 // ParallelOutputCmd size
	}

	message := &Message{
		Id:       ProtocolPTPPOCmd,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, paramsLen),
	}

	// 设置PTPCmd
	message.Params[0] = uint8(ptpCmd.PTPMode)
	binary.LittleEndian.PutUint32(message.Params[4:8], math.Float32bits(ptpCmd.X))
	binary.LittleEndian.PutUint32(message.Params[8:12], math.Float32bits(ptpCmd.Y))
	binary.LittleEndian.PutUint32(message.Params[12:16], math.Float32bits(ptpCmd.Z))
	binary.LittleEndian.PutUint32(message.Params[16:20], math.Float32bits(ptpCmd.R))

	// 设置ParallelOutputCmd
	message.Params[20] = uint8(len(parallelCmd))
	for i, cmd := range parallelCmd {
		offset := 21 + i*4
		message.Params[offset] = cmd.Ratio
		binary.LittleEndian.PutUint16(message.Params[offset+1:offset+3], cmd.Address)
		message.Params[offset+3] = cmd.Level
	}

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// SetPTPPOWithLCmd 设置带L轴的PTP并行输出命令
func (dobot *Dobot) SetPTPPOWithLCmd(ptpWithLCmd *PTPWithLCmd, parallelCmd []ParallelOutputCmd) (queuedCmdIndex uint64, err error) {
	if ptpWithLCmd == nil {
		return 0, errors.New("invalid params: ptpWithLCmd is nil")
	}

	// 计算参数总长度
	paramsLen := 24 // PTPWithLCmd size
	paramsLen += 1  // parallelCmdCount
	if len(parallelCmd) > 0 {
		paramsLen += len(parallelCmd) * 4 // ParallelOutputCmd size
	}

	message := &Message{
		Id:       ProtocolPTPPOWithLCmd,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, paramsLen),
	}

	// 设置PTPWithLCmd
	message.Params[0] = uint8(ptpWithLCmd.PTPMode)
	binary.LittleEndian.PutUint32(message.Params[4:8], math.Float32bits(ptpWithLCmd.X))
	binary.LittleEndian.PutUint32(message.Params[8:12], math.Float32bits(ptpWithLCmd.Y))
	binary.LittleEndian.PutUint32(message.Params[12:16], math.Float32bits(ptpWithLCmd.Z))
	binary.LittleEndian.PutUint32(message.Params[16:20], math.Float32bits(ptpWithLCmd.R))
	binary.LittleEndian.PutUint32(message.Params[20:24], math.Float32bits(ptpWithLCmd.L))

	// 设置ParallelOutputCmd
	message.Params[24] = uint8(len(parallelCmd))
	for i, cmd := range parallelCmd {
		offset := 25 + i*4
		message.Params[offset] = cmd.Ratio
		binary.LittleEndian.PutUint16(message.Params[offset+1:offset+3], cmd.Address)
		message.Params[offset+3] = cmd.Level
	}

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// SetWIFIConfigMode 设置WIFI配置模式
func (dobot *Dobot) SetWIFIConfigMode(enable bool) error {
	message := &Message{
		Id:       ProtocolWIFIConfigMode,
		RW:       true,
		IsQueued: false,
		Params:   make([]byte, 1),
	}

	if enable {
		message.Params[0] = 1
	}

	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// GetWIFIConfigMode 获取WIFI配置模式
func (dobot *Dobot) GetWIFIConfigMode() (bool, error) {
	message := &Message{
		Id:       ProtocolWIFIConfigMode,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return false, err
	}
	if len(resp.Params) < 1 {
		return false, errors.New("invalid response")
	}
	return resp.Params[0] != 0, nil
}

// SetWIFISSID 设置WIFI SSID
func (dobot *Dobot) SetWIFISSID(ssid string) error {
	if ssid == "" {
		return errors.New("invalid params: empty ssid")
	}
	message := &Message{
		Id:       ProtocolWIFISSID,
		RW:       true,
		IsQueued: false,
		Params:   []byte(ssid),
	}
	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// GetWIFISSID 获取WIFI SSID
func (dobot *Dobot) GetWIFISSID() (string, error) {
	message := &Message{
		Id:       ProtocolWIFISSID,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return "", err
	}
	return string(resp.Params), nil
}

// SetWIFIPassword 设置WIFI密码
func (dobot *Dobot) SetWIFIPassword(password string) error {
	if password == "" {
		return errors.New("invalid params: empty password")
	}
	message := &Message{
		Id:       ProtocolWIFIPassword,
		RW:       true,
		IsQueued: false,
		Params:   []byte(password),
	}
	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// GetWIFIPassword 获取WIFI密码
func (dobot *Dobot) GetWIFIPassword() (string, error) {
	message := &Message{
		Id:       ProtocolWIFIPassword,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return "", err
	}
	return string(resp.Params), nil
}

// SetLostStepParams 设置丢步参数
func (dobot *Dobot) SetLostStepParams(threshold float32) (queuedCmdIndex uint64, err error) {
	message := &Message{
		Id:       ProtocolLostStepSet,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 4),
	}

	binary.LittleEndian.PutUint32(message.Params[0:4], math.Float32bits(threshold))

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}

// SetLostStepCmd 设置丢步命令
func (dobot *Dobot) SetLostStepCmd() (queuedCmdIndex uint64, err error) {
	message := &Message{
		Id:       ProtocolLostStepDetect,
		RW:       true,
		IsQueued: true,
	}

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	if len(resp.Params) < 8 {
		return 0, errors.New("invalid response")
	}
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Params)
	return queuedCmdIndex, nil
}
