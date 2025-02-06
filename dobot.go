package godobot

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/zdypro888/godobot/internal"
)

var ErrLeftSpace = internal.ErrLeftSpace

// Dobot 机械臂控制结构
type Dobot struct {
	conn *internal.Connector
}

// NewDobot 创建新的Dobot实例
func NewDobot() *Dobot {
	return &Dobot{}
}

// ConnectDobot 连接到Dobot设备
func (dobot *Dobot) Connect(portName string, baudrate uint32) error {
	dobot.conn = &internal.Connector{}
	err := dobot.conn.Open(portName, baudrate)
	if err != nil {
		return err
	}
	return nil
}

func (dobot *Dobot) Close() error {
	return dobot.conn.Close()
}

type QueuedCommander func() (uint64, error)

func (dobot *Dobot) QueuedSend(command QueuedCommander) (uint64, error) {
	for {
		cmdIndex, err := command()
		if err != nil {
			if err == ErrLeftSpace {
				time.Sleep(time.Millisecond * 10)
				continue
			}
			return 0, err
		}
		return cmdIndex, nil
	}
}

func (dobot *Dobot) QueuedComplete(command QueuedCommander) error {
	var err error
	var cmdIndex uint64
	for {
		if cmdIndex, err = command(); err != nil {
			if err == ErrLeftSpace {
				time.Sleep(time.Millisecond * 10)
				continue
			}
			return err
		}
		break
	}
	for {
		nowIndex, err := dobot.GetQueuedCmdCurrentIndex()
		if err != nil {
			return err
		}
		if cmdIndex <= nowIndex {
			break
		}
		time.Sleep(time.Millisecond * 10)
	}
	return nil
}

// SetDeviceSN 设置设备序列号
func (dobot *Dobot) SetDeviceSN(sn string) error {
	if sn == "" {
		return errors.New("invalid params: empty sn")
	}
	message := &internal.Message{
		Id:       internal.ProtocolDeviceSN,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	writer.WriteString(sn)
	writer.WriteByte(0)
	message.Params = writer.Bytes()
	_, err := dobot.conn.SendMessage(message)
	return err
}

// GetDeviceSN 获取设备序列号
func (dobot *Dobot) GetDeviceSN() (string, error) {
	message := &internal.Message{
		Id:       internal.ProtocolDeviceSN,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return "", err
	}
	return string(resp.Data()), nil
}

// SetDeviceName 设置设备名称
func (dobot *Dobot) SetDeviceName(name string) error {
	if name == "" {
		return errors.New("invalid params: empty name")
	}
	message := &internal.Message{
		Id:       internal.ProtocolDeviceName,
		RW:       true,
		IsQueued: false,
	}
	message.Params = []byte(name)
	message.Params = append(message.Params, 0) // 添加一个字节 0x00 作为校验字节
	_, err := dobot.conn.SendMessage(message)
	return err
}

// GetDeviceName 获取设备名称
func (dobot *Dobot) GetDeviceName() (string, error) {
	message := &internal.Message{
		Id:       internal.ProtocolDeviceName,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return "", err
	}
	return string(resp.Data()), nil
}

// GetDeviceVersion 获取设备版本信息
func (dobot *Dobot) GetDeviceVersion() (majorVersion, minorVersion, revision, hwVersion uint8, err error) {
	message := &internal.Message{
		Id:       internal.ProtocolDeviceVersion,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	return resp.Params[0], resp.Params[1], resp.Params[2], resp.Params[3], nil
}

// SetDeviceWithL 设置设备L轴
func (dobot *Dobot) SetDeviceWithL(isWithL bool, version uint8) (uint64, error) {
	message := &internal.Message{
		Id:       internal.ProtocolDeviceWithL,
		RW:       true,
		IsQueued: true,
		Params:   make([]byte, 2),
	}
	if isWithL {
		message.Params[0] = 1
	}
	message.Params[1] = version
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}
	return resp.Uint64(), nil
}

// GetDeviceWithL 获取设备L轴状态
func (dobot *Dobot) GetDeviceWithL() (bool, error) {
	message := &internal.Message{
		Id:       internal.ProtocolDeviceWithL,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return false, err
	}
	return resp.Bool(), nil
}

// GetDeviceTime 获取设备运行时间
func (dobot *Dobot) GetDeviceTime() (uint32, error) {
	message := &internal.Message{
		Id:       internal.ProtocolDeviceTime,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}
	return resp.Uint32(), nil
}

// GetDeviceInfo 获取设备信息
func (dobot *Dobot) GetDeviceInfo() (*DeviceCountInfo, error) {
	message := &internal.Message{
		Id:       internal.ProtocolDeviceInfo,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) GetPose() (*Pose, error) {
	message := &internal.Message{
		Id:       internal.ProtocolGetPose,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) ResetPose(manual bool, rearArmAngle, frontArmAngle float32) error {
	message := &internal.Message{
		Id:       internal.ProtocolResetPose,
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
	_, err := dobot.conn.SendMessage(message)
	return err
}

// GetKinematics 获取运动学参数
func (dobot *Dobot) GetKinematics() (*Kinematics, error) {
	message := &internal.Message{
		Id: internal.ProtocolGetKinematics,
		RW: false,
	}
	resp, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) GetPoseL() (float32, error) {
	message := &internal.Message{
		Id: internal.ProtocolGetPoseL,
		RW: false,
	}
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}
	return resp.Float32(), nil
}

// GetAlarmsState 获取报警状态
func (dobot *Dobot) GetAlarmsState() ([]uint8, error) {
	message := &internal.Message{
		Id:       internal.ProtocolAlarmsState,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return nil, err
	}
	return resp.Data(), nil
}

// ClearAllAlarmsState 清除所有报警状态
func (dobot *Dobot) ClearAllAlarmsState() error {
	message := &internal.Message{
		Id:       internal.ProtocolAlarmsState,
		RW:       true,
		IsQueued: false,
	}
	_, err := dobot.conn.SendMessage(message)
	return err
}

// SetHOMEParams 设置HOME参数
func (dobot *Dobot) SetHOMEParams(params *HOMEParams, isQueued bool) (uint64, error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}
	message := &internal.Message{
		Id:       internal.ProtocolHOMEParams,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}
	return resp.Uint64(), nil
}

// GetHOMEParams 获取HOME参数
func (dobot *Dobot) GetHOMEParams() (*HOMEParams, error) {
	message := &internal.Message{
		Id:       internal.ProtocolHOMEParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) SetHOMECmd(cmd *HOMECmd, isQueued bool) (uint64, error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}
	message := &internal.Message{
		Id:       internal.ProtocolHOMECmd,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, cmd)
	message.Params = writer.Bytes()
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}
	return resp.Uint64(), nil
}

// SetAutoLevelingCmd 执行自动调平
func (dobot *Dobot) SetAutoLevelingCmd(cmd *AutoLevelingCmd, isQueued bool) (queuedCmdIndex uint64, err error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}
	message := &internal.Message{
		Id:       internal.ProtocolAutoLeveling,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, cmd)
	message.Params = writer.Bytes()

	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}
	return resp.Uint64(), nil
}

// GetAutoLevelingResult 获取自动调平结果
func (dobot *Dobot) GetAutoLevelingResult() (float32, error) {
	message := &internal.Message{
		Id:       internal.ProtocolAutoLeveling,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}
	return resp.Float32(), nil
}

// SetHHTTrigMode 设置手持示教触发模式
func (dobot *Dobot) SetHHTTrigMode(mode HHTTrigMode) error {
	message := &internal.Message{
		Id:       internal.ProtocolHHTTrigMode,
		RW:       true,
		IsQueued: false,
		Params:   []byte{uint8(mode)},
	}
	_, err := dobot.conn.SendMessage(message)
	return err
}

// GetHHTTrigMode 获取手持示教触发模式
func (dobot *Dobot) GetHHTTrigMode() (HHTTrigMode, error) {
	message := &internal.Message{
		Id:       internal.ProtocolHHTTrigMode,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}
	return HHTTrigMode(resp.Params[0]), nil
}

// SetHHTTrigOutputEnabled 设置手持示教触发输出使能
func (dobot *Dobot) SetHHTTrigOutputEnabled(enabled bool) error {
	message := &internal.Message{
		Id:       internal.ProtocolHHTTrigOutputEnabled,
		RW:       true,
		IsQueued: false,
		Params:   make([]byte, 1),
	}
	if enabled {
		message.Params[0] = 1
	}
	_, err := dobot.conn.SendMessage(message)
	return err
}

// GetHHTTrigOutputEnabled 获取手持示教触发输出使能状态
func (dobot *Dobot) GetHHTTrigOutputEnabled() (bool, error) {
	message := &internal.Message{
		Id:       internal.ProtocolHHTTrigOutputEnabled,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return false, err
	}
	return resp.Bool(), nil
}

// GetHHTTrigOutput 获取手持示教触发输出状态
func (dobot *Dobot) GetHHTTrigOutput() (bool, error) {
	message := &internal.Message{
		Id:       internal.ProtocolHHTTrigOutput,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return false, err
	}
	return resp.Bool(), nil
}

// SetEndEffectorParams 设置末端执行器参数
func (dobot *Dobot) SetEndEffectorParams(params *EndEffectorParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &internal.Message{
		Id:       internal.ProtocolEndEffectorParams,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}
	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetEndEffectorParams 获取末端执行器参数
func (dobot *Dobot) GetEndEffectorParams() (*EndEffectorParams, error) {
	message := &internal.Message{
		Id:       internal.ProtocolEndEffectorParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) SetEndEffectorLaser(enableCtrl bool, on bool, isQueued bool) (queuedCmdIndex uint64, err error) {
	message := &internal.Message{
		Id:       internal.ProtocolEndEffectorLaser,
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
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetEndEffectorLaser 获取末端激光状态
func (dobot *Dobot) GetEndEffectorLaser() (isCtrlEnabled bool, isOn bool, err error) {
	message := &internal.Message{
		Id:       internal.ProtocolEndEffectorLaser,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return false, false, err
	}
	return resp.Params[0] != 0, resp.Params[1] != 0, nil
}

// SetEndEffectorSuctionCup 设置末端吸盘状态
func (dobot *Dobot) SetEndEffectorSuctionCup(enableCtrl bool, suck bool, isQueued bool) (queuedCmdIndex uint64, err error) {
	message := &internal.Message{
		Id:       internal.ProtocolEndEffectorSuctionCup,
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
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetEndEffectorSuctionCup 获取末端执行器吸盘状态
func (dobot *Dobot) GetEndEffectorSuctionCup() (isCtrlEnabled bool, isSucked bool, err error) {
	message := &internal.Message{
		Id:       internal.ProtocolEndEffectorSuctionCup,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return false, false, err
	}
	return resp.Params[0] != 0, resp.Params[1] != 0, nil
}

// SetEndEffectorGripper 设置末端夹爪状态
func (dobot *Dobot) SetEndEffectorGripper(enableCtrl bool, grip bool, isQueued bool) (queuedCmdIndex uint64, err error) {
	message := &internal.Message{
		Id:       internal.ProtocolEndEffectorGripper,
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
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetEndEffectorGripper 获取末端夹爪状态
func (dobot *Dobot) GetEndEffectorGripper() (isCtrlEnabled bool, isGripped bool, err error) {
	message := &internal.Message{
		Id:       internal.ProtocolEndEffectorGripper,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return false, false, err
	}
	return resp.Params[0] != 0, resp.Params[1] != 0, nil
}

// SetArmOrientation 设置机械臂方向
func (dobot *Dobot) SetArmOrientation(armOrientation ArmOrientation, isQueued bool) (queuedCmdIndex uint64, err error) {
	message := &internal.Message{
		Id:       internal.ProtocolArmOrientation,
		RW:       true,
		IsQueued: isQueued,
		Params:   []byte{uint8(armOrientation)},
	}

	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetArmOrientation 获取机械臂方向
func (dobot *Dobot) GetArmOrientation() (ArmOrientation, error) {
	message := &internal.Message{
		Id:       internal.ProtocolArmOrientation,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}
	return ArmOrientation(resp.Params[0]), nil
}

// SetJOGJointParams 设置关节点动参数
func (dobot *Dobot) SetJOGJointParams(params *JOGJointParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}
	message := &internal.Message{
		Id:       internal.ProtocolJOGJointParams,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}
	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetJOGJointParams 获取关节点动参数
func (dobot *Dobot) GetJOGJointParams() (*JOGJointParams, error) {
	message := &internal.Message{
		Id:       internal.ProtocolJOGJointParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) SetJOGCoordinateParams(params *JOGCoordinateParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &internal.Message{
		Id:       internal.ProtocolJOGCoordinateParams,
		RW:       true,
		IsQueued: isQueued,
	}

	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()

	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetJOGCoordinateParams 获取坐标点动参数
func (dobot *Dobot) GetJOGCoordinateParams() (*JOGCoordinateParams, error) {
	message := &internal.Message{
		Id:       internal.ProtocolJOGCoordinateParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) SetJOGLParams(params *JOGLParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &internal.Message{
		Id:       internal.ProtocolJOGLParams,
		RW:       true,
		IsQueued: isQueued,
	}

	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()

	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, err
}

// GetJOGLParams 获取JOGL参数
func (dobot *Dobot) GetJOGLParams() (*JOGLParams, error) {
	message := &internal.Message{
		Id:       internal.ProtocolJOGLParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) SetJOGCommonParams(params *JOGCommonParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}
	message := &internal.Message{
		Id:       internal.ProtocolJOGCommonParams,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetJOGCommonParams 获取JOG通用参数
func (dobot *Dobot) GetJOGCommonParams() (*JOGCommonParams, error) {
	message := &internal.Message{
		Id:       internal.ProtocolJOGCommonParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) SetJOGCmd(cmd *JOGCmd, isQueued bool) (queuedCmdIndex uint64, err error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}

	message := &internal.Message{
		Id:       internal.ProtocolJOGCmd,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, cmd)
	message.Params = writer.Bytes()

	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// SetPTPJointParams 设置PTP关节参数
func (dobot *Dobot) SetPTPJointParams(params *PTPJointParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}
	message := &internal.Message{
		Id:       internal.ProtocolPTPJointParams,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}
	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetPTPJointParams 获取PTP关节参数
func (dobot *Dobot) GetPTPJointParams() (*PTPJointParams, error) {
	message := &internal.Message{
		Id:       internal.ProtocolPTPJointParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) SetPTPCoordinateParams(params *PTPCoordinateParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &internal.Message{
		Id:       internal.ProtocolPTPCoordinateParams,
		RW:       true,
		IsQueued: isQueued,
	}

	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()

	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetPTPCoordinateParams 获取PTP坐标运动参数
func (dobot *Dobot) GetPTPCoordinateParams() (*PTPCoordinateParams, error) {
	message := &internal.Message{
		Id:       internal.ProtocolPTPCoordinateParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) SetPTPLParams(params *PTPLParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &internal.Message{
		Id:       internal.ProtocolPTPLParams,
		RW:       true,
		IsQueued: isQueued,
	}

	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()

	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetPTPLParams 获取PTPL运动参数
func (dobot *Dobot) GetPTPLParams() (*PTPLParams, error) {
	message := &internal.Message{
		Id:       internal.ProtocolPTPLParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) SetPTPJumpParams(params *PTPJumpParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid para dms: params is nil")
	}

	message := &internal.Message{
		Id:       internal.ProtocolPTPJumpParams,
		RW:       true,
		IsQueued: isQueued,
	}

	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()

	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetPTPJumpParams 获取PTP跳跃参数
func (dobot *Dobot) GetPTPJumpParams() (*PTPJumpParams, error) {
	message := &internal.Message{
		Id:       internal.ProtocolPTPJumpParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) SetPTPJump2Params(params *PTPJump2Params, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &internal.Message{
		Id:       internal.ProtocolPTPJump2Params,
		RW:       true,
		IsQueued: isQueued,
	}

	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()

	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetPTPJump2Params 获取PTP跳跃2参数
func (dobot *Dobot) GetPTPJump2Params() (*PTPJump2Params, error) {
	message := &internal.Message{
		Id:       internal.ProtocolPTPJump2Params,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) SetPTPCommonParams(params *PTPCommonParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &internal.Message{
		Id:       internal.ProtocolPTPCommonParams,
		RW:       true,
		IsQueued: isQueued,
	}

	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()

	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

func (dobot *Dobot) GetPTPCommonParams() (*PTPCommonParams, error) {
	message := &internal.Message{
		Id:       internal.ProtocolPTPCommonParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) SetPTPCmd(cmd *PTPCmd, isQueued bool) (queuedCmdIndex uint64, err error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}
	message := &internal.Message{
		Id:       internal.ProtocolPTPCmd,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, cmd)
	message.Params = writer.Bytes()
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}
	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// SetPTPWithLCmd 设置带L轴的PTP运动指令
func (dobot *Dobot) SetPTPWithLCmd(cmd *PTPWithLCmd, isQueued bool) (queuedCmdIndex uint64, err error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}

	message := &internal.Message{
		Id:       internal.ProtocolPTPWithLCmd,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, cmd)
	message.Params = writer.Bytes()

	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}
	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// SetCPParams 设置CP参数
func (dobot *Dobot) SetCPParams(params *CPParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}
	message := &internal.Message{
		Id:       internal.ProtocolCPParams,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// SetCPCmd 设置连续运动命令
func (dobot *Dobot) SetCPCmd(cmd *CPCmd, isQueued bool) (queuedCmdIndex uint64, err error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}

	message := &internal.Message{
		Id:       internal.ProtocolCPCmd,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, cmd)
	message.Params = writer.Bytes()

	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// SetCPLECmd 设置连续运动扩展命令
func (dobot *Dobot) SetCPLECmd(cpMode uint8, x, y, z, power float32, isQueued bool) (queuedCmdIndex uint64, err error) {
	message := &internal.Message{
		Id:       internal.ProtocolCPLECmd,
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
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// SetCPRHoldEnable 设置CPR保持使能
func (dobot *Dobot) SetCPRHoldEnable(isEnable bool) error {
	message := &internal.Message{
		Id:       internal.ProtocolCPRHoldEnable,
		RW:       true,
		IsQueued: false,
		Params:   make([]byte, 1),
	}
	if isEnable {
		message.Params[0] = 1
	}
	_, err := dobot.conn.SendMessage(message)
	return err
}

// GetCPRHoldEnable 获取CP运动保持使能状态
func (dobot *Dobot) GetCPRHoldEnable() (bool, error) {
	message := &internal.Message{
		Id:       internal.ProtocolCPRHoldEnable,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return false, err
	}
	return resp.Params[0] != 0, nil
}

// SetCPCommonParams 设置CP通用参数
func (dobot *Dobot) SetCPCommonParams(params *CPCommonParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}
	message := &internal.Message{
		Id:       internal.ProtocolCPCommonParams,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetCPCommonParams 获取CP通用参数
func (dobot *Dobot) GetCPCommonParams() (*CPCommonParams, error) {
	message := &internal.Message{
		Id:       internal.ProtocolCPCommonParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) SetARCParams(params *ARCParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &internal.Message{
		Id:       internal.ProtocolARCParams,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}
	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetARCParams 获取ARC参数
func (dobot *Dobot) GetARCParams() (*ARCParams, error) {
	message := &internal.Message{
		Id:       internal.ProtocolARCParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) SetARCCmd(cmd *ARCCmd, isQueued bool) (queuedCmdIndex uint64, err error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}

	message := &internal.Message{
		Id:       internal.ProtocolARCCmd,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, cmd)
	message.Params = writer.Bytes()

	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// SetCircleCmd 设置圆周运动命令
func (dobot *Dobot) SetCircleCmd(cmd *CircleCmd, isQueued bool) (queuedCmdIndex uint64, err error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}

	message := &internal.Message{
		Id:       internal.ProtocolCircleCmd,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, cmd)
	message.Params = writer.Bytes()

	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// SetARCCommonParams 设置ARC通用参数
func (dobot *Dobot) SetARCCommonParams(params *ARCCommonParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &internal.Message{
		Id:       internal.ProtocolARCCommonParams,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetARCCommonParams 获取ARC通用参数
func (dobot *Dobot) GetARCCommonParams() (*ARCCommonParams, error) {
	message := &internal.Message{
		Id:       internal.ProtocolARCCommonParams,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) SetWAITCmd(cmd *WAITCmd, isQueued bool) (queuedCmdIndex uint64, err error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}

	message := &internal.Message{
		Id:       internal.ProtocolWAITCmd,
		RW:       true,
		IsQueued: true,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, cmd)
	message.Params = writer.Bytes()
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// SetTRIGCmd 设置触发指令
func (dobot *Dobot) SetTRIGCmd(cmd *TRIGCmd, isQueued bool) (queuedCmdIndex uint64, err error) {
	if cmd == nil {
		return 0, errors.New("invalid params: cmd is nil")
	}

	message := &internal.Message{
		Id:       internal.ProtocolTRIGCmd,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, cmd)
	message.Params = writer.Bytes()

	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// SetIOMultiplexing 设置IO复用功能
func (dobot *Dobot) SetIOMultiplexing(params *IOMultiplexing, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &internal.Message{
		Id:       internal.ProtocolIOMultiplexing,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()

	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// SetIODO 设置IO数字输出
func (dobot *Dobot) SetIODO(params *IODO, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &internal.Message{
		Id:       internal.ProtocolIODO,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()

	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// SetIOPWM 设置IO PWM输出
func (dobot *Dobot) SetIOPWM(params *IOPWM, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &internal.Message{
		Id:       internal.ProtocolIOPWM,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()

	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// GetIODI 获取IO数字输入
func (dobot *Dobot) GetIODI(ioDI *IODI) (*IODI, error) {
	message := &internal.Message{
		Id:       internal.ProtocolIODI,
		RW:       false,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, ioDI)
	message.Params = writer.Bytes()
	resp, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) GetIOADC(ioDI *IODI) (*IOADC, error) {
	message := &internal.Message{
		Id:       internal.ProtocolIOADC,
		RW:       false,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, ioDI)
	message.Params = writer.Bytes()
	resp, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) SetEMotor(params *EMotor, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &internal.Message{
		Id:       internal.ProtocolEMotor,
		RW:       true,
		IsQueued: isQueued,
	}

	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()

	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}
	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// SetEMotorS 设置扩展步进电机参数
func (dobot *Dobot) SetEMotorS(params *EMotorS, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}
	message := &internal.Message{
		Id:       internal.ProtocolEMotorS,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()

	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}
	if isQueued {
		return resp.Uint64(), nil
	}
	return 0, nil
}

// SetColorSensor 设置颜色传感器
func (dobot *Dobot) SetColorSensor(enable bool, colorPort ColorPort, version uint8) error {
	message := &internal.Message{
		Id:       internal.ProtocolColorSensor,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, enable)
	binary.Write(writer, binary.LittleEndian, colorPort)
	binary.Write(writer, binary.LittleEndian, version)
	message.Params = writer.Bytes()

	_, err := dobot.conn.SendMessage(message)
	return err
}

// GetColorSensor 获取颜色传感器数据
func (dobot *Dobot) GetColorSensor() (r, g, b uint8, err error) {
	message := &internal.Message{
		Id:       internal.ProtocolColorSensor,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, 0, 0, err
	}
	return resp.Params[0], resp.Params[1], resp.Params[2], nil
}

// SetAngleSensorStaticError 设置角度传感器静态误差
func (dobot *Dobot) SetAngleSensorStaticError(rearArmAngleError, frontArmAngleError float32) error {
	message := &internal.Message{
		Id:       internal.ProtocolAngleSensorStaticError,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, rearArmAngleError)
	binary.Write(writer, binary.LittleEndian, frontArmAngleError)
	message.Params = writer.Bytes()

	_, err := dobot.conn.SendMessage(message)
	return err
}

// GetAngleSensorStaticError 获取角度传感器静态误差
func (dobot *Dobot) GetAngleSensorStaticError() (float32, float32, error) {
	message := &internal.Message{
		Id:       internal.ProtocolAngleSensorStaticError,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) SetAngleSensorCoef(rearArmAngleCoef, frontArmAngleCoef float32) error {
	message := &internal.Message{
		Id:       internal.ProtocolAngleSensorCoef,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, rearArmAngleCoef)
	binary.Write(writer, binary.LittleEndian, frontArmAngleCoef)
	message.Params = writer.Bytes()

	_, err := dobot.conn.SendMessage(message)
	return err
}

// GetAngleSensorCoef 获取角度传感器系数
func (dobot *Dobot) GetAngleSensorCoef() (float32, float32, error) {
	message := &internal.Message{
		Id:       internal.ProtocolAngleSensorCoef,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) SetBaseDecoderStaticError(baseDecoderError float32) error {
	message := &internal.Message{
		Id:       internal.ProtocolBaseDecoderStaticError,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, baseDecoderError)
	message.Params = writer.Bytes()

	_, err := dobot.conn.SendMessage(message)
	return err
}

// GetBaseDecoderStaticError 获取底座解码器静态误差
func (dobot *Dobot) GetBaseDecoderStaticError() (float32, error) {
	message := &internal.Message{
		Id:       internal.ProtocolBaseDecoderStaticError,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}
	return resp.Float32(), nil
}

// SetLRHandCalibrateValue 设置左右手校准值
func (dobot *Dobot) SetLRHandCalibrateValue(lrHandCalibrateValue float32) error {
	message := &internal.Message{
		Id:       internal.ProtocolLRHandCalibrateValue,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, lrHandCalibrateValue)
	message.Params = writer.Bytes()

	_, err := dobot.conn.SendMessage(message)
	return err
}

// GetLRHandCalibrateValue 获取左右手校准值
func (dobot *Dobot) GetLRHandCalibrateValue() (float32, error) {
	message := &internal.Message{
		Id:       internal.ProtocolLRHandCalibrateValue,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}
	return resp.Float32(), nil
}

// SetQueuedCmdStartExec 开始执行指令队列
func (dobot *Dobot) SetQueuedCmdStartExec() error {
	message := &internal.Message{
		Id:       internal.ProtocolQueuedCmdStartExec,
		RW:       true,
		IsQueued: false,
	}
	_, err := dobot.conn.SendMessage(message)
	return err
}

// SetQueuedCmdStopExec 停止执行队列命令
func (dobot *Dobot) SetQueuedCmdStopExec() error {
	message := &internal.Message{
		Id:       internal.ProtocolQueuedCmdStopExec,
		RW:       true,
		IsQueued: false,
	}
	_, err := dobot.conn.SendMessage(message)
	return err
}

// SetQueuedCmdForceStopExec 强制停止执行队列命令
func (dobot *Dobot) SetQueuedCmdForceStopExec() error {
	message := &internal.Message{
		Id:       internal.ProtocolQueuedCmdForceStopExec,
		RW:       true,
		IsQueued: false,
	}
	_, err := dobot.conn.SendMessage(message)
	return err
}

// SetQueuedCmdStartDownload 开始下载队列命令
func (dobot *Dobot) SetQueuedCmdStartDownload(totalLoop uint32, linePerLoop uint32) error {
	message := &internal.Message{
		Id:       internal.ProtocolQueuedCmdStartDownload,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, totalLoop)
	binary.Write(writer, binary.LittleEndian, linePerLoop)
	message.Params = writer.Bytes()

	_, err := dobot.conn.SendMessage(message)
	return err
}

// SetQueuedCmdStopDownload 停止下载队列命令
func (dobot *Dobot) SetQueuedCmdStopDownload() error {
	message := &internal.Message{
		Id:       internal.ProtocolQueuedCmdStopDownload,
		RW:       true,
		IsQueued: false,
	}
	_, err := dobot.conn.SendMessage(message)
	return err
}

// SetQueuedCmdClear 清除队列命令
func (dobot *Dobot) SetQueuedCmdClear() error {
	message := &internal.Message{
		Id:       internal.ProtocolQueuedCmdClear,
		RW:       true,
		IsQueued: false,
	}
	_, err := dobot.conn.SendMessage(message)
	return err
}

// GetQueuedCmdCurrentIndex 获取当前队列命令索引
func (dobot *Dobot) GetQueuedCmdCurrentIndex() (uint64, error) {
	message := &internal.Message{
		Id:       internal.ProtocolQueuedCmdCurrentIndex,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}
	return resp.Uint64(), nil
}

// GetQueuedCmdMotionFinish 获取队列命令运动是否完成
func (dobot *Dobot) GetQueuedCmdMotionFinish() (bool, error) {
	message := &internal.Message{
		Id:       internal.ProtocolQueuedCmdMotionFinish,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return false, err
	}
	return resp.Bool(), nil
}

// SetPTPPOCmd 设置PTP并行输出命令
func (dobot *Dobot) SetPTPPOCmd(ptpCmd *PTPCmd, parallelCmd []ParallelOutputCmd) (queuedCmdIndex uint64, err error) {
	if ptpCmd == nil {
		return 0, errors.New("invalid params: ptpCmd is nil")
	}
	message := &internal.Message{
		Id:       internal.ProtocolPTPPOCmd,
		RW:       true,
		IsQueued: true,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, ptpCmd)
	writer.WriteByte(uint8(len(parallelCmd)))
	binary.Write(writer, binary.LittleEndian, parallelCmd)
	message.Params = writer.Bytes()

	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	return resp.Uint64(), nil
}

// SetPTPPOWithLCmd 设置带并行输出和L轴的PTP运动指令
func (dobot *Dobot) SetPTPPOWithLCmd(ptpWithLCmd *PTPWithLCmd, parallelCmd []ParallelOutputCmd) (queuedCmdIndex uint64, err error) {
	if ptpWithLCmd == nil {
		return 0, errors.New("invalid params: ptpWithLCmd is nil")
	}
	message := &internal.Message{
		Id:       internal.ProtocolPTPPOWithLCmd,
		RW:       true,
		IsQueued: true,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, ptpWithLCmd)
	writer.WriteByte(uint8(len(parallelCmd)))
	binary.Write(writer, binary.LittleEndian, parallelCmd)
	message.Params = writer.Bytes()

	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	return resp.Uint64(), nil
}

// SetWIFIConfigMode 设置WIFI配置模式
func (dobot *Dobot) SetWIFIConfigMode(enable bool) error {
	message := &internal.Message{
		Id:       internal.ProtocolWIFIConfigMode,
		RW:       true,
		IsQueued: false,
		Params:   make([]byte, 1),
	}

	if enable {
		message.Params[0] = 1
	}

	_, err := dobot.conn.SendMessage(message)
	return err
}

// GetWIFIConfigMode 获取WIFI配置模式状态
func (dobot *Dobot) GetWIFIConfigMode() (bool, error) {
	message := &internal.Message{
		Id:       internal.ProtocolWIFIConfigMode,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return false, err
	}
	return resp.Bool(), nil
}

// SetWIFISSID 设置WIFI SSID
func (dobot *Dobot) SetWIFISSID(ssid string) error {
	if ssid == "" {
		return errors.New("invalid params: empty ssid")
	}
	message := &internal.Message{
		Id:       internal.ProtocolWIFISSID,
		RW:       true,
		IsQueued: false,
	}

	writer := &bytes.Buffer{}
	writer.WriteString(ssid)
	writer.WriteByte(0) // 添加一个字节 0x00 作为校验字节
	message.Params = writer.Bytes()

	_, err := dobot.conn.SendMessage(message)
	return err
}

// GetWIFISSID 获取WIFI SSID
func (dobot *Dobot) GetWIFISSID() (string, error) {
	message := &internal.Message{
		Id:       internal.ProtocolWIFISSID,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return "", err
	}
	return string(resp.Data()), nil
}

// SetWIFIPassword 设置WIFI密码
func (dobot *Dobot) SetWIFIPassword(password string) error {
	if password == "" {
		return errors.New("invalid params: empty password")
	}
	message := &internal.Message{
		Id:       internal.ProtocolWIFIPassword,
		RW:       true,
		IsQueued: false,
	}

	writer := &bytes.Buffer{}
	writer.WriteString(password)
	writer.WriteByte(0)
	message.Params = writer.Bytes()

	_, err := dobot.conn.SendMessage(message)
	return err
}

// GetWIFIPassword 获取WIFI密码
func (dobot *Dobot) GetWIFIPassword() (string, error) {
	message := &internal.Message{
		Id:       internal.ProtocolWIFIPassword,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return "", err
	}
	return string(resp.Data()), nil
}

// GetPTPTime 获取PTP运动时间
func (dobot *Dobot) GetPTPTime() (float32, error) {
	message := &internal.Message{
		Id:       internal.ProtocolPTPTime,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) GetFirmwareMode() (FirmwareMode, error) {
	message := &internal.Message{
		Id:       internal.ProtocolFirmwareMode,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}
	return FirmwareMode(resp.Params[0]), nil
}

// SetLostStepParams 设置丢步参数
func (dobot *Dobot) SetLostStepParams(threshold float32, isQueued bool) (uint64, error) {
	message := &internal.Message{
		Id:       internal.ProtocolLostStepSet,
		RW:       true,
		IsQueued: isQueued,
	}

	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, threshold)
	message.Params = writer.Bytes()

	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	return resp.Uint64(), nil
}

// SetLostStepCmd 设置丢步命令
func (dobot *Dobot) SetLostStepCmd(isQueued bool) (uint64, error) {
	message := &internal.Message{
		Id:       internal.ProtocolLostStepDetect,
		RW:       true,
		IsQueued: isQueued,
	}

	resp, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	return resp.Uint64(), nil
}

// GetUART4PeripheralsType 获取UART4外设类型
func (dobot *Dobot) GetUART4PeripheralsType() (uint8, error) {
	message := &internal.Message{
		Id:       internal.ProtocolCheckUART4PeripheralsModel,
		RW:       false,
		IsQueued: false,
	}

	response, err := dobot.conn.SendMessage(message)
	if err != nil {
		return 0, err
	}

	return response.Params[0], nil
}

// SetUART4PeripheralsEnable 设置UART4外设使能状态
func (dobot *Dobot) SetUART4PeripheralsEnable(isEnable bool) error {
	message := &internal.Message{
		Id:       internal.ProtocolUART4PeripheralsEnabled,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, isEnable)
	message.Params = writer.Bytes()

	_, err := dobot.conn.SendMessage(message)
	return err
}

// GetUART4PeripheralsEnable 获取UART4外设使能状态
func (dobot *Dobot) GetUART4PeripheralsEnable() (bool, error) {
	message := &internal.Message{
		Id:       internal.ProtocolUART4PeripheralsEnabled,
		RW:       false,
		IsQueued: false,
		Params:   []byte{},
	}

	response, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) SendPluse(pluseCmd *PluseCmd, isQueued bool) (uint64, error) {
	message := &internal.Message{
		Id:       internal.ProtocolFunctionPulseMode,
		RW:       true,
		IsQueued: isQueued,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, pluseCmd)
	message.Params = writer.Bytes()
	response, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) SetWIFIIPAddress(wifiIPAddress *WIFIIPAddress) error {
	message := &internal.Message{
		Id:       internal.ProtocolWIFIIPAddress,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, wifiIPAddress)
	message.Params = writer.Bytes()

	_, err := dobot.conn.SendMessage(message)
	return err
}

// GetWIFIIPAddress 获取WIFI IP地址
func (dobot *Dobot) GetWIFIIPAddress() (*WIFIIPAddress, error) {
	message := &internal.Message{
		Id:       internal.ProtocolWIFIIPAddress,
		RW:       false,
		IsQueued: false,
	}
	response, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) SetWIFINetmask(wifiNetmask *WIFINetmask) error {
	message := &internal.Message{
		Id:       internal.ProtocolWIFINetmask,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, wifiNetmask)
	message.Params = writer.Bytes()

	_, err := dobot.conn.SendMessage(message)
	return err
}

// GetWIFINetmask 获取WIFI子网掩码
func (dobot *Dobot) GetWIFINetmask() (*WIFINetmask, error) {
	message := &internal.Message{
		Id:       internal.ProtocolWIFINetmask,
		RW:       false,
		IsQueued: false,
	}

	response, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) SetWIFIGateway(wifiGateway *WIFIGateway) error {
	message := &internal.Message{
		Id:       internal.ProtocolWIFIGateway,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, wifiGateway)
	message.Params = writer.Bytes()

	_, err := dobot.conn.SendMessage(message)
	return err
}

// GetWIFIGateway 获取WIFI网关
func (dobot *Dobot) GetWIFIGateway() (*WIFIGateway, error) {
	message := &internal.Message{
		Id:       internal.ProtocolWIFIGateway,
		RW:       false,
		IsQueued: false,
	}

	response, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) SetWIFIDNS(wifiDNS *WIFIDNS) error {
	message := &internal.Message{
		Id:       internal.ProtocolWIFIDNS,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, wifiDNS)
	message.Params = writer.Bytes()

	_, err := dobot.conn.SendMessage(message)
	return err
}

// GetWIFIDNS 获取WIFI DNS
func (dobot *Dobot) GetWIFIDNS() (*WIFIDNS, error) {
	message := &internal.Message{
		Id:       internal.ProtocolWIFIDNS,
		RW:       false,
		IsQueued: false,
	}

	response, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) GetWIFIConnectStatus() (bool, error) {
	message := &internal.Message{
		Id:       internal.ProtocolWIFIConnectStatus,
		RW:       false,
		IsQueued: false,
	}

	response, err := dobot.conn.SendMessage(message)
	if err != nil {
		return false, err
	}
	return response.Bool(), nil
}

// GetIOMultiplexing 获取IO复用状态
func (dobot *Dobot) GetIOMultiplexing(ioMultiplexing *IOMultiplexing) (*IOMultiplexing, error) {
	message := &internal.Message{
		Id:       internal.ProtocolIOMultiplexing,
		RW:       false,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, ioMultiplexing)
	message.Params = writer.Bytes()
	resp, err := dobot.conn.SendMessage(message)
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
func (dobot *Dobot) GetIODO(ioDO *IODO) (*IODO, error) {
	message := &internal.Message{
		Id:       internal.ProtocolIODO,
		RW:       false,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, ioDO)
	message.Params = writer.Bytes()

	resp, err := dobot.conn.SendMessage(message)
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
