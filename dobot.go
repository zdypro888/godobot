package godobot

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
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
func (dobot *Dobot) Connect(ctx context.Context, portName string, baudrate uint32) error {
	err := dobot.connector.Open(ctx, portName, baudrate)
	if err != nil {
		return err
	}
	return nil
}
func (dobot *Dobot) QueuedComplete(ctx context.Context, command func() (uint64, error)) error {
	cmdIndex, err := command()
	if err != nil {
		return err
	}
	for {
		nowIndex, err := dobot.GetQueuedCmdCurrentIndex(ctx)
		if err != nil {
			return err
		}
		if cmdIndex >= nowIndex {
			break
		}
		time.Sleep(time.Millisecond * 10)
	}
	return nil
}

// SetDeviceSN 设置设备序列号
func (dobot *Dobot) SetDeviceSN(ctx context.Context, sn string) error {
	if sn == "" {
		return errors.New("invalid params: empty sn")
	}
	message := &Message{
		Id:       ProtocolDeviceSN,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	writer.WriteString(sn)
	writer.WriteByte(0)
	message.Params = writer.Bytes()
	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// GetDeviceSN 获取设备序列号
func (dobot *Dobot) GetDeviceSN(ctx context.Context) (string, error) {
	message := &Message{
		Id:       ProtocolDeviceSN,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return "", err
	}
	return string(resp.Data()), nil
}

// SetDeviceName 设置设备名称
func (dobot *Dobot) SetDeviceName(ctx context.Context, name string) error {
	if name == "" {
		return errors.New("invalid params: empty name")
	}
	message := &Message{
		Id:       ProtocolDeviceName,
		RW:       true,
		IsQueued: false,
	}
	message.Params = []byte(name)
	message.Params = append(message.Params, 0) // 添加一个字节 0x00 作为校验字节
	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// GetDeviceName 获取设备名称
func (dobot *Dobot) GetDeviceName(ctx context.Context) (string, error) {
	message := &Message{
		Id:       ProtocolDeviceName,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return "", err
	}
	return string(resp.Data()), nil
}

// GetDeviceVersion 获取设备版本信息
func (dobot *Dobot) GetDeviceVersion(ctx context.Context) (majorVersion, minorVersion, revision, hwVersion uint8, err error) {
	message := &Message{
		Id:       ProtocolDeviceVersion,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	return resp.Params[0], resp.Params[1], resp.Params[2], resp.Params[3], nil
}

// SetDeviceWithL 设置设备L轴
func (dobot *Dobot) SetDeviceWithL(ctx context.Context, isWithL bool, version uint8) (uint64, error) {
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
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}
	return resp.Uint64(), nil
}

// GetDeviceWithL 获取设备L轴状态
func (dobot *Dobot) GetDeviceWithL(ctx context.Context) (bool, error) {
	message := &Message{
		Id:       ProtocolDeviceWithL,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return false, err
	}
	return resp.Bool(), nil
}

// GetDeviceTime 获取设备运行时间
func (dobot *Dobot) GetDeviceTime(ctx context.Context) (uint32, error) {
	message := &Message{
		Id:       ProtocolDeviceTime,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}
	return resp.Uint32(), nil
}

// GetDeviceInfo 获取设备信息
func (dobot *Dobot) GetDeviceInfo(ctx context.Context) (*DeviceCountInfo, error) {
	message := &Message{
		Id:       ProtocolDeviceInfo,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}
	info := &DeviceCountInfo{}
	if err := resp.Read(info); err != nil {
		return nil, err
	}
	return info, nil
}

// GetPose 获取当前位姿信息
func (dobot *Dobot) GetPose(ctx context.Context) (*Pose, error) {
	message := &Message{
		Id:       ProtocolGetPose,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}
	pose := &Pose{}
	if err := resp.Read(pose); err != nil {
		return nil, err
	}
	return pose, nil
}

// ResetPose 重置位姿到指定状态
func (dobot *Dobot) ResetPose(ctx context.Context, manual bool, rearArmAngle, frontArmAngle float32) error {
	message := &Message{
		Id:       ProtocolResetPose,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	if manual {
		writer.WriteByte(1)
	} else {
		writer.WriteByte(0)
	}
	binary.Write(writer, binary.LittleEndian, rearArmAngle)
	binary.Write(writer, binary.LittleEndian, frontArmAngle)
	message.Params = writer.Bytes()
	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// GetKinematics 获取运动学参数
func (dobot *Dobot) GetKinematics(ctx context.Context) (*Kinematics, error) {
	message := &Message{
		Id: ProtocolGetKinematics,
		RW: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}
	kinematics := &Kinematics{}
	if err := resp.Read(kinematics); err != nil {
		return nil, err
	}
	return kinematics, nil
}

// GetPoseL 获取L轴位置
func (dobot *Dobot) GetPoseL(ctx context.Context) (float32, error) {
	message := &Message{
		Id: ProtocolGetPoseL,
		RW: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}
	return resp.Float32(), nil
}

// GetAlarmsState 获取报警状态
func (dobot *Dobot) GetAlarmsState(ctx context.Context) ([]uint8, error) {
	message := &Message{
		Id: ProtocolAlarmsState,
		RW: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}
	return resp.Data(), nil
}

// ClearAllAlarmsState 清除所有报警状态
func (dobot *Dobot) ClearAllAlarmsState(ctx context.Context) error {
	message := &Message{
		Id: ProtocolAlarmsState,
		RW: true,
	}
	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// SetHOMEParams 设置HOME参数
func (dobot *Dobot) SetHOMEParams(ctx context.Context, params *HOMEParams, isQueued bool) (uint64, error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}
	message := &Message{
		Id:       ProtocolHOMEParams,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}
	return resp.Uint64(), nil
}

// GetHOMEParams 获取HOME参数
func (dobot *Dobot) GetHOMEParams(ctx context.Context) (*HOMEParams, error) {
	message := &Message{
		Id:       ProtocolHOMEParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}
	params := &HOMEParams{}
	if err := resp.Read(params); err != nil {
		return nil, err
	}
	return params, nil
}

// SetHOMECmd 执行回零操作
func (dobot *Dobot) SetHOMECmd(ctx context.Context, cmd *HOMECmd, isQueued bool) (uint64, error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}
	message := &Message{
		Id:       ProtocolHOMECmd,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, cmd)
	message.Params = writer.Bytes()
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}
	return resp.Uint64(), nil
}

// SetAutoLevelingCmd 执行自动调平
func (dobot *Dobot) SetAutoLevelingCmd(ctx context.Context, cmd *AutoLevelingCmd, isQueued bool) (queuedCmdIndex uint64, err error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}
	message := &Message{
		Id:       ProtocolAutoLeveling,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, cmd)
	message.Params = writer.Bytes()

	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}
	return resp.Uint64(), nil
}

// GetAutoLevelingResult 获取自动调平结果
func (dobot *Dobot) GetAutoLevelingResult(ctx context.Context) (float32, error) {
	message := &Message{
		Id:       ProtocolAutoLeveling,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}
	return resp.Float32(), nil
}

// SetHHTTrigMode 设置手持示教触发模式
func (dobot *Dobot) SetHHTTrigMode(ctx context.Context, mode HHTTrigMode) error {
	message := &Message{
		Id:       ProtocolHHTTrigMode,
		RW:       true,
		IsQueued: false,
		Params:   []byte{uint8(mode)},
	}
	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// GetHHTTrigMode 获取手持示教触发模式
func (dobot *Dobot) GetHHTTrigMode(ctx context.Context) (HHTTrigMode, error) {
	message := &Message{
		Id:       ProtocolHHTTrigMode,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}
	return HHTTrigMode(resp.Params[0]), nil
}

// SetHHTTrigOutputEnabled 设置手持示教触发输出使能
func (dobot *Dobot) SetHHTTrigOutputEnabled(ctx context.Context, enabled bool) error {
	message := &Message{
		Id:       ProtocolHHTTrigOutputEnabled,
		RW:       true,
		IsQueued: false,
		Params:   make([]byte, 1),
	}
	if enabled {
		message.Params[0] = 1
	}
	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// GetHHTTrigOutputEnabled 获取手持示教触发输出使能状态
func (dobot *Dobot) GetHHTTrigOutputEnabled(ctx context.Context) (bool, error) {
	message := &Message{
		Id:       ProtocolHHTTrigOutputEnabled,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return false, err
	}
	return resp.Bool(), nil
}

// GetHHTTrigOutput 获取手持示教触发输出状态
func (dobot *Dobot) GetHHTTrigOutput(ctx context.Context) (bool, error) {
	message := &Message{
		Id:       ProtocolHHTTrigOutput,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return false, err
	}
	return resp.Bool(), nil
}

// SetEndEffectorParams 设置末端执行器参数
func (dobot *Dobot) SetEndEffectorParams(ctx context.Context, params *EndEffectorParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolEndEffectorParams,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}
	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetEndEffectorParams 获取末端执行器参数
func (dobot *Dobot) GetEndEffectorParams(ctx context.Context) (*EndEffectorParams, error) {
	message := &Message{
		Id:       ProtocolEndEffectorParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}
	params := &EndEffectorParams{}
	if err := resp.Read(params); err != nil {
		return nil, err
	}
	return params, nil
}

// SetEndEffectorLaser 设置末端激光状态
func (dobot *Dobot) SetEndEffectorLaser(ctx context.Context, enableCtrl bool, on bool, isQueued bool) (queuedCmdIndex uint64, err error) {
	message := &Message{
		Id:       ProtocolEndEffectorLaser,
		RW:       true,
		IsQueued: isQueued,
		Params:   make([]byte, 2),
	}
	if enableCtrl {
		message.Params[0] = 1
	}
	if on {
		message.Params[1] = 1
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetEndEffectorLaser 获取末端激光状态
func (dobot *Dobot) GetEndEffectorLaser(ctx context.Context) (isCtrlEnabled bool, isOn bool, err error) {
	message := &Message{
		Id:       ProtocolEndEffectorLaser,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return false, false, err
	}
	return resp.Params[0] != 0, resp.Params[1] != 0, nil
}

// SetEndEffectorSuctionCup 设置末端吸盘状态
func (dobot *Dobot) SetEndEffectorSuctionCup(ctx context.Context, enableCtrl bool, suck bool, isQueued bool) (queuedCmdIndex uint64, err error) {
	message := &Message{
		Id:       ProtocolEndEffectorSuctionCup,
		RW:       true,
		IsQueued: isQueued,
		Params:   make([]byte, 2),
	}
	if enableCtrl {
		message.Params[0] = 1
	}
	if suck {
		message.Params[1] = 1
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetEndEffectorSuctionCup 获取末端执行器吸盘状态
func (dobot *Dobot) GetEndEffectorSuctionCup(ctx context.Context) (isCtrlEnabled bool, isSucked bool, err error) {
	message := &Message{
		Id:       ProtocolEndEffectorSuctionCup,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return false, false, err
	}
	return resp.Params[0] != 0, resp.Params[1] != 0, nil
}

// SetEndEffectorGripper 设置末端夹爪状态
func (dobot *Dobot) SetEndEffectorGripper(ctx context.Context, enableCtrl bool, grip bool, isQueued bool) (queuedCmdIndex uint64, err error) {
	message := &Message{
		Id:       ProtocolEndEffectorGripper,
		RW:       true,
		IsQueued: isQueued,
		Params:   make([]byte, 2),
	}
	if enableCtrl {
		message.Params[0] = 1
	}
	if grip {
		message.Params[1] = 1
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetEndEffectorGripper 获取末端夹爪状态
func (dobot *Dobot) GetEndEffectorGripper(ctx context.Context) (isCtrlEnabled bool, isGripped bool, err error) {
	message := &Message{
		Id:       ProtocolEndEffectorGripper,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return false, false, err
	}
	return resp.Params[0] != 0, resp.Params[1] != 0, nil
}

// SetArmOrientation 设置机械臂方向
func (dobot *Dobot) SetArmOrientation(ctx context.Context, armOrientation ArmOrientation, isQueued bool) (queuedCmdIndex uint64, err error) {
	message := &Message{
		Id:       ProtocolArmOrientation,
		RW:       true,
		IsQueued: isQueued,
		Params:   []byte{uint8(armOrientation)},
	}

	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetArmOrientation 获取机械臂方向
func (dobot *Dobot) GetArmOrientation(ctx context.Context) (ArmOrientation, error) {
	message := &Message{
		Id:       ProtocolArmOrientation,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}
	return ArmOrientation(resp.Params[0]), nil
}

// SetJOGJointParams 设置关节点动参数
func (dobot *Dobot) SetJOGJointParams(ctx context.Context, params *JOGJointParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}
	message := &Message{
		Id:       ProtocolJOGJointParams,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}
	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetJOGJointParams 获取关节点动参数
func (dobot *Dobot) GetJOGJointParams(ctx context.Context) (*JOGJointParams, error) {
	message := &Message{
		Id:       ProtocolJOGJointParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}
	params := &JOGJointParams{}
	if err := resp.Read(params); err != nil {
		return nil, fmt.Errorf("failed to read JOG joint params: %v", err)
	}
	return params, nil
}

// SetJOGCoordinateParams 设置坐标点动参数
func (dobot *Dobot) SetJOGCoordinateParams(ctx context.Context, params *JOGCoordinateParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolJOGCoordinateParams,
		RW:       true,
		IsQueued: isQueued,
	}

	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()

	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetJOGCoordinateParams 获取坐标点动参数
func (dobot *Dobot) GetJOGCoordinateParams(ctx context.Context) (*JOGCoordinateParams, error) {
	message := &Message{
		Id:       ProtocolJOGCoordinateParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	params := &JOGCoordinateParams{}
	if err := resp.Read(params); err != nil {
		return nil, fmt.Errorf("failed to read JOG coordinate params: %v", err)
	}
	return params, nil
}

// SetJOGLParams 设置JOGL参数
func (dobot *Dobot) SetJOGLParams(ctx context.Context, params *JOGLParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolJOGLParams,
		RW:       true,
		IsQueued: isQueued,
	}

	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()

	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, err
}

// GetJOGLParams 获取JOGL参数
func (dobot *Dobot) GetJOGLParams(ctx context.Context) (*JOGLParams, error) {
	message := &Message{
		Id:       ProtocolJOGLParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	params := &JOGLParams{}
	if err := resp.Read(params); err != nil {
		return nil, fmt.Errorf("failed to read JOG L params: %v", err)
	}
	return params, nil
}

// SetJOGCommonParams 设置JOG通用参数
func (dobot *Dobot) SetJOGCommonParams(ctx context.Context, params *JOGCommonParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}
	message := &Message{
		Id:       ProtocolJOGCommonParams,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetJOGCommonParams 获取JOG通用参数
func (dobot *Dobot) GetJOGCommonParams(ctx context.Context) (*JOGCommonParams, error) {
	message := &Message{
		Id:       ProtocolJOGCommonParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	params := &JOGCommonParams{}
	if err := resp.Read(params); err != nil {
		return nil, fmt.Errorf("failed to read JOG common params: %v", err)
	}
	return params, nil
}

// SetJOGCmd 设置JOG运动指令
func (dobot *Dobot) SetJOGCmd(ctx context.Context, cmd *JOGCmd, isQueued bool) (queuedCmdIndex uint64, err error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}

	message := &Message{
		Id:       ProtocolJOGCmd,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, cmd)
	message.Params = writer.Bytes()

	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// SetPTPJointParams 设置PTP关节参数
func (dobot *Dobot) SetPTPJointParams(ctx context.Context, params *PTPJointParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}
	message := &Message{
		Id:       ProtocolPTPJointParams,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}
	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetPTPJointParams 获取PTP关节参数
func (dobot *Dobot) GetPTPJointParams(ctx context.Context) (*PTPJointParams, error) {
	message := &Message{
		Id:       ProtocolPTPJointParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	params := &PTPJointParams{}
	if err := resp.Read(params); err != nil {
		return nil, fmt.Errorf("failed to read PTP joint params: %v", err)
	}
	return params, nil
}

// SetPTPCoordinateParams 设置PTP坐标运动参数
func (dobot *Dobot) SetPTPCoordinateParams(ctx context.Context, params *PTPCoordinateParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolPTPCoordinateParams,
		RW:       true,
		IsQueued: isQueued,
	}

	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()

	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetPTPCoordinateParams 获取PTP坐标运动参数
func (dobot *Dobot) GetPTPCoordinateParams(ctx context.Context) (*PTPCoordinateParams, error) {
	message := &Message{
		Id:       ProtocolPTPCoordinateParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	params := &PTPCoordinateParams{}
	if err := resp.Read(params); err != nil {
		return nil, fmt.Errorf("failed to read PTP coordinate params: %v", err)
	}
	return params, nil
}

// SetPTPLParams 设置PTPL运动参数
func (dobot *Dobot) SetPTPLParams(ctx context.Context, params *PTPLParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolPTPLParams,
		RW:       true,
		IsQueued: isQueued,
	}

	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()

	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetPTPLParams 获取PTPL运动参数
func (dobot *Dobot) GetPTPLParams(ctx context.Context) (*PTPLParams, error) {
	message := &Message{
		Id:       ProtocolPTPLParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	params := &PTPLParams{}
	if err := resp.Read(params); err != nil {
		return nil, fmt.Errorf("failed to read PTP L params: %v", err)
	}
	return params, nil
}

// SetPTPJumpParams 设置PTP跳跃参数
func (dobot *Dobot) SetPTPJumpParams(ctx context.Context, params *PTPJumpParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid para dms: params is nil")
	}

	message := &Message{
		Id:       ProtocolPTPJumpParams,
		RW:       true,
		IsQueued: isQueued,
	}

	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()

	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetPTPJumpParams 获取PTP跳跃参数
func (dobot *Dobot) GetPTPJumpParams(ctx context.Context) (*PTPJumpParams, error) {
	message := &Message{
		Id:       ProtocolPTPJumpParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	params := &PTPJumpParams{}
	if err := resp.Read(params); err != nil {
		return nil, fmt.Errorf("failed to read PTP jump params: %v", err)
	}
	return params, nil
}

// SetPTPJump2Params 设置PTP跳跃2参数
func (dobot *Dobot) SetPTPJump2Params(ctx context.Context, params *PTPJump2Params, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolPTPJump2Params,
		RW:       true,
		IsQueued: isQueued,
	}

	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()

	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetPTPJump2Params 获取PTP跳跃2参数
func (dobot *Dobot) GetPTPJump2Params(ctx context.Context) (*PTPJump2Params, error) {
	message := &Message{
		Id:       ProtocolPTPJump2Params,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	params := &PTPJump2Params{}
	if err := resp.Read(params); err != nil {
		return nil, fmt.Errorf("failed to read PTP jump2 params: %v", err)
	}
	return params, nil
}

// SetPTPCommonParams 设置PTP通用参数
func (dobot *Dobot) SetPTPCommonParams(ctx context.Context, params *PTPCommonParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolPTPCommonParams,
		RW:       true,
		IsQueued: isQueued,
	}

	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()

	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

func (dobot *Dobot) GetPTPCommonParams(ctx context.Context) (*PTPCommonParams, error) {
	message := &Message{
		Id:       ProtocolPTPCommonParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	params := &PTPCommonParams{}
	if err := resp.Read(params); err != nil {
		return nil, fmt.Errorf("failed to read PTP common params: %v", err)
	}
	return params, nil
}

// SetPTPCmd 设置PTP命令
func (dobot *Dobot) SetPTPCmd(ctx context.Context, cmd *PTPCmd, isQueued bool) (queuedCmdIndex uint64, err error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}
	message := &Message{
		Id:       ProtocolPTPCmd,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, cmd)
	message.Params = writer.Bytes()
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}
	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// SetPTPWithLCmd 设置带L轴的PTP运动指令
func (dobot *Dobot) SetPTPWithLCmd(ctx context.Context, cmd *PTPWithLCmd, isQueued bool) (queuedCmdIndex uint64, err error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}

	message := &Message{
		Id:       ProtocolPTPWithLCmd,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, cmd)
	message.Params = writer.Bytes()

	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}
	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// SetCPParams 设置CP参数
func (dobot *Dobot) SetCPParams(ctx context.Context, params *CPParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}
	message := &Message{
		Id:       ProtocolCPParams,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// SetCPCmd 设置连续运动命令
func (dobot *Dobot) SetCPCmd(ctx context.Context, cmd *CPCmd, isQueued bool) (queuedCmdIndex uint64, err error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}

	message := &Message{
		Id:       ProtocolCPCmd,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, cmd)
	message.Params = writer.Bytes()

	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// SetCPLECmd 设置连续运动扩展命令
func (dobot *Dobot) SetCPLECmd(ctx context.Context, cpMode uint8, x, y, z, power float32, isQueued bool) (queuedCmdIndex uint64, err error) {
	message := &Message{
		Id:       ProtocolCPLECmd,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, cpMode)
	binary.Write(writer, binary.LittleEndian, x)
	binary.Write(writer, binary.LittleEndian, y)
	binary.Write(writer, binary.LittleEndian, z)
	binary.Write(writer, binary.LittleEndian, power)
	message.Params = writer.Bytes()
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// SetCPRHoldEnable 设置CPR保持使能
func (dobot *Dobot) SetCPRHoldEnable(ctx context.Context, isEnable bool) error {
	message := &Message{
		Id:       ProtocolCPRHoldEnable,
		RW:       true,
		IsQueued: false,
		Params:   make([]byte, 1),
	}
	if isEnable {
		message.Params[0] = 1
	}
	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// GetCPRHoldEnable 获取CP运动保持使能状态
func (dobot *Dobot) GetCPRHoldEnable(ctx context.Context) (bool, error) {
	message := &Message{
		Id:       ProtocolCPRHoldEnable,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return false, err
	}
	return resp.Params[0] != 0, nil
}

// SetCPCommonParams 设置CP通用参数
func (dobot *Dobot) SetCPCommonParams(ctx context.Context, params *CPCommonParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}
	message := &Message{
		Id:       ProtocolCPCommonParams,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetCPCommonParams 获取CP通用参数
func (dobot *Dobot) GetCPCommonParams(ctx context.Context) (*CPCommonParams, error) {
	message := &Message{
		Id:       ProtocolCPCommonParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	params := &CPCommonParams{}
	if err := resp.Read(params); err != nil {
		return nil, fmt.Errorf("failed to read CP common params: %v", err)
	}
	return params, nil
}

// SetARCParams 设置ARC参数
func (dobot *Dobot) SetARCParams(ctx context.Context, params *ARCParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolARCParams,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}
	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetARCParams 获取ARC参数
func (dobot *Dobot) GetARCParams(ctx context.Context) (*ARCParams, error) {
	message := &Message{
		Id:       ProtocolARCParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	params := &ARCParams{}
	if err := resp.Read(params); err != nil {
		return nil, fmt.Errorf("failed to read ARC params: %v", err)
	}
	return params, nil
}

// SetARCCmd 设置ARC命令
func (dobot *Dobot) SetARCCmd(ctx context.Context, cmd *ARCCmd, isQueued bool) (queuedCmdIndex uint64, err error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}

	message := &Message{
		Id:       ProtocolARCCmd,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, cmd)
	message.Params = writer.Bytes()

	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// SetCircleCmd 设置圆周运动命令
func (dobot *Dobot) SetCircleCmd(ctx context.Context, cmd *CircleCmd, isQueued bool) (queuedCmdIndex uint64, err error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}

	message := &Message{
		Id:       ProtocolCircleCmd,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, cmd)
	message.Params = writer.Bytes()

	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// SetARCCommonParams 设置ARC通用参数
func (dobot *Dobot) SetARCCommonParams(ctx context.Context, params *ARCCommonParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolARCCommonParams,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetARCCommonParams 获取ARC通用参数
func (dobot *Dobot) GetARCCommonParams(ctx context.Context) (*ARCCommonParams, error) {
	message := &Message{
		Id:       ProtocolARCCommonParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	params := &ARCCommonParams{}
	if err := resp.Read(params); err != nil {
		return nil, fmt.Errorf("failed to read ARC common params: %v", err)
	}
	return params, nil
}

// SetWAITCmd 设置等待指令
func (dobot *Dobot) SetWAITCmd(ctx context.Context, cmd *WAITCmd, isQueued bool) (queuedCmdIndex uint64, err error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}

	message := &Message{
		Id:       ProtocolWAITCmd,
		RW:       true,
		IsQueued: true,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, cmd)
	message.Params = writer.Bytes()
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// SetTRIGCmd 设置触发指令
func (dobot *Dobot) SetTRIGCmd(ctx context.Context, cmd *TRIGCmd, isQueued bool) (queuedCmdIndex uint64, err error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}

	message := &Message{
		Id:       ProtocolTRIGCmd,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, cmd)
	message.Params = writer.Bytes()

	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// SetIOMultiplexing 设置IO复用功能
func (dobot *Dobot) SetIOMultiplexing(ctx context.Context, params *IOMultiplexing, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolIOMultiplexing,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()

	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// SetIODO 设置IO数字输出
func (dobot *Dobot) SetIODO(ctx context.Context, params *IODO, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolIODO,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()

	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// SetIOPWM 设置IO PWM输出
func (dobot *Dobot) SetIOPWM(ctx context.Context, params *IOPWM, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolIOPWM,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()

	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetIODI 获取IO数字输入
func (dobot *Dobot) GetIODI(ctx context.Context, ioDI *IODI) (*IODI, error) {
	message := &Message{
		Id:       ProtocolIODI,
		RW:       false,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, ioDI)
	message.Params = writer.Bytes()
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	params := &IODI{}
	if err := resp.Read(params); err != nil {
		return nil, fmt.Errorf("failed to read IO DI: %v", err)
	}
	return params, nil
}

// GetIOADC 获取IO模拟输入
func (dobot *Dobot) GetIOADC(ctx context.Context, ioDI *IODI) (*IOADC, error) {
	message := &Message{
		Id:       ProtocolIOADC,
		RW:       false,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, ioDI)
	message.Params = writer.Bytes()
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}
	params := &IOADC{}
	if err := resp.Read(params); err != nil {
		return nil, fmt.Errorf("failed to read IO ADC: %v", err)
	}
	return params, nil
}

// SetEMotor 设置扩展电机参数
func (dobot *Dobot) SetEMotor(ctx context.Context, params *EMotor, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolEMotor,
		RW:       true,
		IsQueued: isQueued,
	}

	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()

	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}
	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// SetEMotorS 设置扩展步进电机参数
func (dobot *Dobot) SetEMotorS(ctx context.Context, params *EMotorS, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}
	message := &Message{
		Id:       ProtocolEMotorS,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()

	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}
	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// SetColorSensor 设置颜色传感器
func (dobot *Dobot) SetColorSensor(ctx context.Context, enable bool, colorPort ColorPort, version uint8) error {
	message := &Message{
		Id:       ProtocolColorSensor,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, enable)
	binary.Write(writer, binary.LittleEndian, colorPort)
	binary.Write(writer, binary.LittleEndian, version)
	message.Params = writer.Bytes()

	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// GetColorSensor 获取颜色传感器数据
func (dobot *Dobot) GetColorSensor(ctx context.Context) (r, g, b uint8, err error) {
	message := &Message{
		Id:       ProtocolColorSensor,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, 0, 0, err
	}
	return resp.Params[0], resp.Params[1], resp.Params[2], nil
}

// SetAngleSensorStaticError 设置角度传感器静态误差
func (dobot *Dobot) SetAngleSensorStaticError(ctx context.Context, rearArmAngleError, frontArmAngleError float32) error {
	message := &Message{
		Id:       ProtocolAngleSensorStaticError,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, rearArmAngleError)
	binary.Write(writer, binary.LittleEndian, frontArmAngleError)
	message.Params = writer.Bytes()

	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// GetAngleSensorStaticError 获取角度传感器静态误差
func (dobot *Dobot) GetAngleSensorStaticError(ctx context.Context) (float32, float32, error) {
	message := &Message{
		Id:       ProtocolAngleSensorStaticError,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, 0, err
	}
	reader := resp.Reader()
	var rearArmAngleError, frontArmAngleError float32
	binary.Read(reader, binary.LittleEndian, &rearArmAngleError)
	binary.Read(reader, binary.LittleEndian, &frontArmAngleError)
	return rearArmAngleError, frontArmAngleError, nil
}

// SetAngleSensorCoef 设置角度传感器系数
func (dobot *Dobot) SetAngleSensorCoef(ctx context.Context, rearArmAngleCoef, frontArmAngleCoef float32) error {
	message := &Message{
		Id:       ProtocolAngleSensorCoef,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, rearArmAngleCoef)
	binary.Write(writer, binary.LittleEndian, frontArmAngleCoef)
	message.Params = writer.Bytes()

	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// GetAngleSensorCoef 获取角度传感器系数
func (dobot *Dobot) GetAngleSensorCoef(ctx context.Context) (float32, float32, error) {
	message := &Message{
		Id:       ProtocolAngleSensorCoef,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, 0, err
	}
	reader := resp.Reader()
	var rearArmAngleCoef, frontArmAngleCoef float32
	binary.Read(reader, binary.LittleEndian, &rearArmAngleCoef)
	binary.Read(reader, binary.LittleEndian, &frontArmAngleCoef)
	return rearArmAngleCoef, frontArmAngleCoef, nil
}

// SetBaseDecoderStaticError 设置底座解码器静态误差
func (dobot *Dobot) SetBaseDecoderStaticError(ctx context.Context, baseDecoderError float32) error {
	message := &Message{
		Id:       ProtocolBaseDecoderStaticError,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, baseDecoderError)
	message.Params = writer.Bytes()

	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// GetBaseDecoderStaticError 获取底座解码器静态误差
func (dobot *Dobot) GetBaseDecoderStaticError(ctx context.Context) (float32, error) {
	message := &Message{
		Id:       ProtocolBaseDecoderStaticError,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}
	return resp.Float32(), nil
}

// SetLRHandCalibrateValue 设置左右手校准值
func (dobot *Dobot) SetLRHandCalibrateValue(ctx context.Context, lrHandCalibrateValue float32) error {
	message := &Message{
		Id:       ProtocolLRHandCalibrateValue,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, lrHandCalibrateValue)
	message.Params = writer.Bytes()

	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// GetLRHandCalibrateValue 获取左右手校准值
func (dobot *Dobot) GetLRHandCalibrateValue(ctx context.Context) (float32, error) {
	message := &Message{
		Id:       ProtocolLRHandCalibrateValue,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}
	return resp.Float32(), nil
}

// SetQueuedCmdStartExec 开始执行指令队列
func (dobot *Dobot) SetQueuedCmdStartExec(ctx context.Context) error {
	message := &Message{
		Id:       ProtocolQueuedCmdStartExec,
		RW:       true,
		IsQueued: false,
	}
	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// SetQueuedCmdStopExec 停止执行队列命令
func (dobot *Dobot) SetQueuedCmdStopExec(ctx context.Context) error {
	message := &Message{
		Id:       ProtocolQueuedCmdStopExec,
		RW:       true,
		IsQueued: false,
	}
	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// SetQueuedCmdForceStopExec 强制停止执行队列命令
func (dobot *Dobot) SetQueuedCmdForceStopExec(ctx context.Context) error {
	message := &Message{
		Id:       ProtocolQueuedCmdForceStopExec,
		RW:       true,
		IsQueued: false,
	}
	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// SetQueuedCmdStartDownload 开始下载队列命令
func (dobot *Dobot) SetQueuedCmdStartDownload(ctx context.Context, totalLoop uint32, linePerLoop uint32) error {
	message := &Message{
		Id:       ProtocolQueuedCmdStartDownload,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, totalLoop)
	binary.Write(writer, binary.LittleEndian, linePerLoop)
	message.Params = writer.Bytes()

	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// SetQueuedCmdStopDownload 停止下载队列命令
func (dobot *Dobot) SetQueuedCmdStopDownload(ctx context.Context) error {
	message := &Message{
		Id:       ProtocolQueuedCmdStopDownload,
		RW:       true,
		IsQueued: false,
	}
	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// SetQueuedCmdClear 清除队列命令
func (dobot *Dobot) SetQueuedCmdClear(ctx context.Context) error {
	message := &Message{
		Id:       ProtocolQueuedCmdClear,
		RW:       true,
		IsQueued: false,
	}
	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// GetQueuedCmdCurrentIndex 获取当前队列命令索引
func (dobot *Dobot) GetQueuedCmdCurrentIndex(ctx context.Context) (uint64, error) {
	message := &Message{
		Id:       ProtocolQueuedCmdCurrentIndex,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}
	return resp.Uint64(), nil
}

// GetQueuedCmdMotionFinish 获取队列命令运动是否完成
func (dobot *Dobot) GetQueuedCmdMotionFinish(ctx context.Context) (bool, error) {
	message := &Message{
		Id:       ProtocolQueuedCmdMotionFinish,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return false, err
	}
	return resp.Bool(), nil
}

// SetPTPPOCmd 设置PTP并行输出命令
func (dobot *Dobot) SetPTPPOCmd(ctx context.Context, ptpCmd *PTPCmd, parallelCmd []ParallelOutputCmd) (queuedCmdIndex uint64, err error) {
	if ptpCmd == nil {
		return 0, errors.New("invalid params: ptpCmd is nil")
	}
	message := &Message{
		Id:       ProtocolPTPPOCmd,
		RW:       true,
		IsQueued: true,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, ptpCmd)
	writer.WriteByte(uint8(len(parallelCmd)))
	binary.Write(writer, binary.LittleEndian, parallelCmd)
	message.Params = writer.Bytes()

	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	return resp.Uint64(), nil
}

// SetPTPPOWithLCmd 设置带并行输出和L轴的PTP运动指令
func (dobot *Dobot) SetPTPPOWithLCmd(ctx context.Context, ptpWithLCmd *PTPWithLCmd, parallelCmd []ParallelOutputCmd) (queuedCmdIndex uint64, err error) {
	if ptpWithLCmd == nil {
		return 0, errors.New("invalid params: ptpWithLCmd is nil")
	}
	message := &Message{
		Id:       ProtocolPTPPOWithLCmd,
		RW:       true,
		IsQueued: true,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, ptpWithLCmd)
	writer.WriteByte(uint8(len(parallelCmd)))
	binary.Write(writer, binary.LittleEndian, parallelCmd)
	message.Params = writer.Bytes()

	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	return resp.Uint64(), nil
}

// SetWIFIConfigMode 设置WIFI配置模式
func (dobot *Dobot) SetWIFIConfigMode(ctx context.Context, enable bool) error {
	message := &Message{
		Id:       ProtocolWIFIConfigMode,
		RW:       true,
		IsQueued: false,
		Params:   make([]byte, 1),
	}

	if enable {
		message.Params[0] = 1
	}

	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// GetWIFIConfigMode 获取WIFI配置模式状态
func (dobot *Dobot) GetWIFIConfigMode(ctx context.Context) (bool, error) {
	message := &Message{
		Id:       ProtocolWIFIConfigMode,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return false, err
	}
	return resp.Bool(), nil
}

// SetWIFISSID 设置WIFI SSID
func (dobot *Dobot) SetWIFISSID(ctx context.Context, ssid string) error {
	if ssid == "" {
		return errors.New("invalid params: empty ssid")
	}
	message := &Message{
		Id:       ProtocolWIFISSID,
		RW:       true,
		IsQueued: false,
	}

	writer := &bytes.Buffer{}
	writer.WriteString(ssid)
	writer.WriteByte(0) // 添加一个字节 0x00 作为校验字节
	message.Params = writer.Bytes()

	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// GetWIFISSID 获取WIFI SSID
func (dobot *Dobot) GetWIFISSID(ctx context.Context) (string, error) {
	message := &Message{
		Id:       ProtocolWIFISSID,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return "", err
	}
	return string(resp.Data()), nil
}

// SetWIFIPassword 设置WIFI密码
func (dobot *Dobot) SetWIFIPassword(ctx context.Context, password string) error {
	if password == "" {
		return errors.New("invalid params: empty password")
	}
	message := &Message{
		Id:       ProtocolWIFIPassword,
		RW:       true,
		IsQueued: false,
	}

	writer := &bytes.Buffer{}
	writer.WriteString(password)
	writer.WriteByte(0)
	message.Params = writer.Bytes()

	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// GetWIFIPassword 获取WIFI密码
func (dobot *Dobot) GetWIFIPassword(ctx context.Context) (string, error) {
	message := &Message{
		Id:       ProtocolWIFIPassword,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return "", err
	}
	return string(resp.Data()), nil
}

// GetPTPTime 获取PTP运动时间
func (dobot *Dobot) GetPTPTime(ctx context.Context) (float32, error) {
	message := &Message{
		Id:       ProtocolPTPTime,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}
	return resp.Float32(), nil
}

// FirmwareMode 固件模式
type FirmwareMode uint8

const (
	FirmwareModeNormal    FirmwareMode = 0 // 正常模式
	FirmwareModeUpgrade   FirmwareMode = 1 // 升级模式
	FirmwareModeCalibrate FirmwareMode = 2 // 校准模式
)

// GetFirmwareMode 获取固件模式
func (dobot *Dobot) GetFirmwareMode(ctx context.Context) (FirmwareMode, error) {
	message := &Message{
		Id:       ProtocolFirmwareMode,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}
	return FirmwareMode(resp.Params[0]), nil
}

// SetLostStepParams 设置丢步参数
func (dobot *Dobot) SetLostStepParams(ctx context.Context, threshold float32, isQueued bool) (uint64, error) {
	message := &Message{
		Id:       ProtocolLostStepSet,
		RW:       true,
		IsQueued: isQueued,
	}

	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, threshold)
	message.Params = writer.Bytes()

	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	return resp.Uint64(), nil
}

// SetLostStepCmd 设置丢步命令
func (dobot *Dobot) SetLostStepCmd(ctx context.Context, isQueued bool) (uint64, error) {
	message := &Message{
		Id:       ProtocolLostStepDetect,
		RW:       true,
		IsQueued: isQueued,
	}

	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	return resp.Uint64(), nil
}

// GetUART4PeripheralsType 获取UART4外设类型
func (dobot *Dobot) GetUART4PeripheralsType(ctx context.Context) (uint8, error) {
	message := &Message{
		Id:       ProtocolCheckUART4PeripheralsModel,
		RW:       false,
		IsQueued: false,
	}

	response, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	return response.Params[0], nil
}

// SetUART4PeripheralsEnable 设置UART4外设使能状态
func (dobot *Dobot) SetUART4PeripheralsEnable(ctx context.Context, isEnable bool) error {
	message := &Message{
		Id:       ProtocolUART4PeripheralsEnabled,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, isEnable)
	message.Params = writer.Bytes()

	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// GetUART4PeripheralsEnable 获取UART4外设使能状态
func (dobot *Dobot) GetUART4PeripheralsEnable(ctx context.Context) (bool, error) {
	message := &Message{
		Id:       ProtocolUART4PeripheralsEnabled,
		RW:       false,
		IsQueued: false,
		Params:   []byte{},
	}

	response, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return false, err
	}

	return response.Bool(), nil
}

// PluseCmd 脉冲控制命令
type PluseCmd struct {
	ControlMode uint8   // 控制模式
	Port        uint8   // 端口号
	Speed       float32 // 速度（脉冲个数每秒）
	Distance    float32 // 移动距离（脉冲个数）
}

// SendPluse 发送脉冲控制命令
func (dobot *Dobot) SendPluse(ctx context.Context, pluseCmd *PluseCmd, isQueued bool) (uint64, error) {
	message := &Message{
		Id:       ProtocolFunctionPulseMode,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, pluseCmd)
	message.Params = writer.Bytes()
	response, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}
	if isQueued {
		return response.Uint64(), nil
	}
	return 0, nil
}

// WIFIIPAddress WIFI IP地址结构
type WIFIIPAddress struct {
	DHCP bool     // 是否开启DHCP
	Addr [4]uint8 // IP地址，分成四段
}

// WIFINetmask WIFI子网掩码结构
type WIFINetmask struct {
	Addr [4]uint8 // 子网掩码，分成四段
}

// WIFIGateway WIFI网关结构
type WIFIGateway struct {
	Addr [4]uint8 // 网关地址，分成四段
}

// WIFIDNS WIFI DNS结构
type WIFIDNS struct {
	Addr [4]uint8 // DNS地址，分成四段
}

// SetWIFIIPAddress 设置WIFI IP地址
func (dobot *Dobot) SetWIFIIPAddress(ctx context.Context, wifiIPAddress *WIFIIPAddress) error {
	message := &Message{
		Id:       ProtocolWIFIIPAddress,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, wifiIPAddress)
	message.Params = writer.Bytes()

	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// GetWIFIIPAddress 获取WIFI IP地址
func (dobot *Dobot) GetWIFIIPAddress(ctx context.Context) (*WIFIIPAddress, error) {
	message := &Message{
		Id:       ProtocolWIFIIPAddress,
		RW:       false,
		IsQueued: false,
	}
	response, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}
	result := &WIFIIPAddress{}
	if err := response.Read(result); err != nil {
		return nil, fmt.Errorf("failed to read WIFI IP address: %v", err)
	}
	return result, nil
}

// SetWIFINetmask 设置WIFI子网掩码
func (dobot *Dobot) SetWIFINetmask(ctx context.Context, wifiNetmask *WIFINetmask) error {
	message := &Message{
		Id:       ProtocolWIFINetmask,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, wifiNetmask)
	message.Params = writer.Bytes()

	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// GetWIFINetmask 获取WIFI子网掩码
func (dobot *Dobot) GetWIFINetmask(ctx context.Context) (*WIFINetmask, error) {
	message := &Message{
		Id:       ProtocolWIFINetmask,
		RW:       false,
		IsQueued: false,
	}

	response, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	result := &WIFINetmask{}
	if err := response.Read(result); err != nil {
		return nil, fmt.Errorf("failed to read WIFI netmask: %v", err)
	}
	return result, nil
}

// SetWIFIGateway 设置WIFI网关
func (dobot *Dobot) SetWIFIGateway(ctx context.Context, wifiGateway *WIFIGateway) error {
	message := &Message{
		Id:       ProtocolWIFIGateway,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, wifiGateway)
	message.Params = writer.Bytes()

	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// GetWIFIGateway 获取WIFI网关
func (dobot *Dobot) GetWIFIGateway(ctx context.Context) (*WIFIGateway, error) {
	message := &Message{
		Id:       ProtocolWIFIGateway,
		RW:       false,
		IsQueued: false,
	}

	response, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	result := &WIFIGateway{}
	if err := response.Read(result); err != nil {
		return nil, fmt.Errorf("failed to read WIFI gateway: %v", err)
	}
	return result, nil
}

// SetWIFIDNS 设置WIFI DNS
func (dobot *Dobot) SetWIFIDNS(ctx context.Context, wifiDNS *WIFIDNS) error {
	message := &Message{
		Id:       ProtocolWIFIDNS,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, wifiDNS)
	message.Params = writer.Bytes()

	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// GetWIFIDNS 获取WIFI DNS
func (dobot *Dobot) GetWIFIDNS(ctx context.Context) (*WIFIDNS, error) {
	message := &Message{
		Id:       ProtocolWIFIDNS,
		RW:       false,
		IsQueued: false,
	}

	response, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	result := &WIFIDNS{}
	if err := response.Read(result); err != nil {
		return nil, fmt.Errorf("failed to read WIFI DNS: %v", err)
	}
	return result, nil
}

// GetWIFIConnectStatus 获取WIFI连接状态
func (dobot *Dobot) GetWIFIConnectStatus(ctx context.Context) (bool, error) {
	message := &Message{
		Id:       ProtocolWIFIConnectStatus,
		RW:       false,
		IsQueued: false,
	}

	response, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return false, err
	}
	return response.Bool(), nil
}

// GetIOMultiplexing 获取IO复用状态
func (dobot *Dobot) GetIOMultiplexing(ctx context.Context, ioMultiplexing *IOMultiplexing) (*IOMultiplexing, error) {
	message := &Message{
		Id:       ProtocolIOMultiplexing,
		RW:       false,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, ioMultiplexing)
	message.Params = writer.Bytes()
	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}
	var multiplexing IOMultiplexing
	if err := resp.Read(&multiplexing); err != nil {
		return nil, fmt.Errorf("failed to read IOMultiplexing: %v", err)
	}
	return &multiplexing, nil
}

// GetIODO 获取IO数字输出状态
func (dobot *Dobot) GetIODO(ctx context.Context, ioDO *IODO) (*IODO, error) {
	message := &Message{
		Id:       ProtocolIODO,
		RW:       false,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, ioDO)
	message.Params = writer.Bytes()

	resp, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}
	reader := resp.Reader()
	var iodo IODO
	if err := binary.Read(reader, binary.LittleEndian, &iodo); err != nil {
		return nil, fmt.Errorf("failed to read IODO: %v", err)
	}
	return &iodo, nil
}
