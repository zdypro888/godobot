package godobot

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

// Dobot 机械臂控制结构
// @Description Dobot Magician 机械臂的主控制结构体，封装了与机械臂通信和控制的所有功能
type Dobot struct {
	connector *Connector
}

// NewDobot 创建新的Dobot实例
// @Summary 创建新的Dobot控制实例
// @Description 初始化并返回一个新的Dobot控制实例。该实例包含了所有控制机械臂所需的方法。
// @Description 注意：创建实例后需要调用Connect方法与实际的机械臂建立连接才能使用。
// @Return *Dobot "返回Dobot实例指针"
// @Example
//
//	dobot := NewDobot()
//	err := dobot.Connect(ctx, "/dev/ttyUSB0", 115200)
func NewDobot() *Dobot {
	return &Dobot{connector: &Connector{}}
}

// ConnectDobot 连接到Dobot设备
// @Summary 连接到Dobot设备
// @Description 通过指定的串口和波特率连接到Dobot设备。
// @Description 该函数会尝试打开指定的串口并建立与机械臂的通信连接。
// @Description 连接成功后才能执行后续的控制命令。
// @Param ctx context.Context true "上下文对象，用于超时控制和取消操作"
// @Param portName string true "串口设备名称：
//   - Windows系统：'COM1', 'COM2', ...
//   - Linux系统：'/dev/ttyUSB0', '/dev/ttyACM0', ...
//   - macOS系统：'/dev/cu.usbserial-*'"
//
// @Param baudrate uint32 true "串口通信波特率，标准值：115200"
// @Success 200 {string} "连接成功"
// @Failure 400 {error} "连接失败，可能的错误：
//   - 串口不存在
//   - 串口被占用
//   - 波特率不支持
//   - 通信超时"
//
// @Example
//
//	ctx := context.Background()
//	err := dobot.Connect(ctx, "/dev/ttyUSB0", 115200)
//	if err != nil {
//	    log.Fatal("连接失败:", err)
//	}
func (dobot *Dobot) Connect(ctx context.Context, portName string, baudrate uint32) error {
	err := dobot.connector.Open(ctx, portName, baudrate)
	if err != nil {
		return err
	}
	return nil
}

// SetDeviceSN 设置设备序列号
// @Summary 设置Dobot设备序列号
// @Description 为Dobot设备设置唯一的序列号标识。序列号用于区分不同的设备。
// @Description 设置后的序列号将被永久保存在设备中。
// @Param sn string true "要设置的序列号：
//   - 长度：必须小于64字符
//   - 格式：字母、数字和特殊字符的组合
//   - 建议：使用有意义的标识，如'DOBOT_001'"
//
// @Success 200 {string} "设置成功"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 序列号为空
//   - 序列号格式无效
//   - 通信错误"
//
// @Example
//
//	err := dobot.SetDeviceSN("DOBOT_LAB_001")
//	if err != nil {
//	    log.Printf("设置序列号失败: %v", err)
//	}
func (dobot *Dobot) SetDeviceSN(sn string) error {
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
	writer.WriteByte(0) // 添加一个字节 0x00 作为校验字节
	message.Params = writer.Bytes()
	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// GetDeviceSN 获取设备序列号
// @Summary 获取Dobot设备序列号
// @Description 获取当前连接的Dobot设备的序列号。
// @Description 序列号是设备的唯一标识符。
// @Return string "设备序列号"
// @Return error "获取失败的错误信息"
// @Success 200 {string} "返回设备序列号"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 设备未连接
//   - 通信错误"
//
// @Example
//
//	sn, err := dobot.GetDeviceSN()
//	if err != nil {
//	    log.Printf("获取序列号失败: %v", err)
//	} else {
//	    log.Printf("设备序列号: %s", sn)
//	}
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
// @Summary 设置Dobot设备名称
// @Description 为Dobot设备设置一个友好的名称，便于识别和管理。
// @Description 设置后的名称将被保存在设备中。
// @Param name string true "要设置的设备名称：
//   - 长度：必须小于64字符
//   - 格式：支持任意可打印字符
//   - 建议：使用易于理解的描述性名称"
//
// @Success 200 {string} "设置成功"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 名称为空
//   - 名称格式无效
//   - 通信错误"
//
// @Example
//
//	err := dobot.SetDeviceName("实验室机械臂1号")
//	if err != nil {
//	    log.Printf("设置设备名称失败: %v", err)
//	}
func (dobot *Dobot) SetDeviceName(name string) error {
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
	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// GetDeviceName 获取设备名称
// @Summary 获取Dobot设备名称
// @Description 获取当前连接的Dobot设备的名称。
// @Description 设备名称是一个便于人类识别的标识符。
// @Return string "设备名称"
// @Return error "获取失败的错误信息"
// @Success 200 {string} "返回设备名称"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 设备未连接
//   - 通信错误"
//
// @Example
//
//	name, err := dobot.GetDeviceName()
//	if err != nil {
//	    log.Printf("获取设备名称失败: %v", err)
//	} else {
//	    log.Printf("设备名称: %s", name)
//	}
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
// @Summary 获取Dobot设备版本信息
// @Description 获取当前连接的Dobot设备的固件版本信息
// @Return majorVersion uint8 主版本号
// @Return minorVersion uint8 次版本号
// @Return revision uint8 修订版本号
// @Return hwVersion uint8 硬件版本号
// @Return error 获取过程中的错误信息，如果成功则返回 nil
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
	data := resp.Data()
	return data[0], data[1], data[2], data[3], nil
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
	data := resp.Data()
	queuedCmdIndex = binary.LittleEndian.Uint64(data)
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
	data := resp.Data()
	return data[0] != 0, nil
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
	data := resp.Data()
	return binary.LittleEndian.Uint32(data), nil
}

// GetDeviceInfo 获取设备信息
// @Summary 获取机械臂设备信息
// @Description 获取机械臂的设备信息，包括设备类型、运行时间、错误次数等统计信息
// @Return *DeviceCountInfo 设备信息结构体指针，包含:
// @Return *DeviceCountInfo.deviceType uint8 设备类型
// @Return *DeviceCountInfo.runTime uint32 运行时间，单位秒
// @Return *DeviceCountInfo.powerOnCount uint32 开机次数
// @Return *DeviceCountInfo.errorCount uint32 错误次数
// @Return *DeviceCountInfo.warningCount uint32 警告次数
// @Return error 获取过程中的错误信息，如果成功则返回 nil
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
	info := &DeviceCountInfo{}
	binary.Read(resp.Reader(), binary.LittleEndian, info)
	return info, nil
}

// GetPose 获取当前位姿
// @Summary 获取机械臂当前位姿信息
// @Description 获取机械臂末端在笛卡尔坐标系下的位置和姿态信息，以及各关节的角度值
// @Return *Pose 位姿结构体指针，包含:
// @Return *Pose.x float32 末端在x轴方向的位置，单位mm
// @Return *Pose.y float32 末端在y轴方向的位置，单位mm
// @Return *Pose.z float32 末端在z轴方向的位置，单位mm
// @Return *Pose.r float32 末端的旋转角度，单位度
// @Return *Pose.jointAngle [4]float32 各关节的角度值，单位度
// @Return error 获取过程中的错误信息，如果成功则返回 nil
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
	pose := &Pose{}
	binary.Read(resp.Reader(), binary.LittleEndian, pose)
	return pose, nil
}

// ResetPose 重置位姿
// @Summary 重置机械臂位姿
// @Description 根据是否手动模式和指定的后臂、前臂角度重置机械臂位姿
// @Param manual query bool true "是否手动重置"
// @Param rearArmAngle query float32 true "后臂角度，单位度"
// @Param frontArmAngle query float32 true "前臂角度，单位度"
// @Success 200 {string} "重置成功返回空字符串"
// @Failure 400 {object} error "重置失败时返回错误信息"
func (dobot *Dobot) ResetPose(manual bool, rearArmAngle, frontArmAngle float32) error {
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
	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// GetKinematics 获取运动学参数
// @Summary 获取机器人运动学信息
// @Description 获取机械臂的运动学参数信息，包括各关节位置等信息
// @Success 200 {object} *Kinematics "返回运动学信息结构体"
// @Failure 400 {object} error "获取运动学信息失败时返回错误信息"
func (dobot *Dobot) GetKinematics() (*Kinematics, error) {
	message := &Message{
		Id: ProtocolGetKinematics,
		RW: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return nil, err
	}
	kinematics := &Kinematics{}
	binary.Read(resp.Reader(), binary.LittleEndian, kinematics)
	return kinematics, nil
}

// GetPoseL 获取L轴位置
// @Summary 获取长臂位姿参数
// @Description 获取机械臂中与长臂相关的参数（例如长臂延伸距离），单位可能为mm
// @Success 200 {number} float32 "返回长臂位姿参数"
// @Failure 400 {object} error "获取长臂位姿参数失败时返回错误信息"
func (dobot *Dobot) GetPoseL() (float32, error) {
	message := &Message{
		Id: ProtocolGetPoseL,
		RW: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}
	data := resp.Data()
	var value float32
	binary.Decode(data, binary.LittleEndian, &value)
	return value, nil
}

// GetAlarmsState 获取报警状态
// @Summary 获取机械臂报警状态
// @Description 获取机械臂当前的报警状态列表，每个uint8代表一种报警代码
// @Success 200 {array} uint8 "返回报警状态列表"
// @Failure 400 {object} error "获取报警状态失败时返回错误信息"
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
// @Summary 清除机械臂报警信息
// @Description 清除机械臂所有报警状态，恢复至正常状态
// @Success 200 {string} "清除成功返回空字符串"
// @Failure 400 {object} error "清除报警状态失败时返回错误信息"
func (dobot *Dobot) ClearAllAlarmsState() error {
	message := &Message{
		Id: ProtocolAlarmsState,
		RW: true,
	}
	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// SetHOMEParams 设置HOME参数
// @Summary 设置HOME点参数
// @Description 设置机械臂返回HOME位置的参数，通过HOMEParams结构体指定各轴的HOME位置
// @Param params body *HOMEParams true "HOME参数结构体"
// @Param isQueued query bool true "是否队列执行"
// @Success 200 {number} uint64 "返回命令队列索引"
// @Failure 400 {object} error "设置HOME参数失败时返回错误信息"
func (dobot *Dobot) SetHOMEParams(params *HOMEParams, isQueued bool) (queuedCmdIndex uint64, err error) {
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
// @Summary 获取HOME点参数
// @Description 获取机械臂的HOME参数，返回包含各轴HOME位置的结构体
// @Success 200 {object} *HOMEParams "返回HOME参数结构体"
// @Failure 400 {object} error "获取HOME参数失败时返回错误信息"
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
	params := &HOMEParams{}
	if err := binary.Read(bytes.NewReader(resp.Params), binary.LittleEndian, params); err != nil {
		return nil, fmt.Errorf("failed to read HOME params: %v", err)
	}
	return params, nil
}

// SetHOMECmd 设置HOME命令
// @Summary 执行HOME命令
// @Description 执行机械臂的HOME命令，通过HOMECmd结构体详细指定归位参数，并决定是否队列执行
// @Param cmd body *HOMECmd true "HOMECmd结构体，包含归位命令参数"
// @Param isQueued query bool true "是否队列执行"
// @Success 200 {number} uint64 "返回命令队列索引"
// @Failure 400 {object} error "执行HOME命令失败时返回错误信息"
func (dobot *Dobot) SetHOMECmd(cmd *HOMECmd, isQueued bool) (queuedCmdIndex uint64, err error) {
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
// @Summary 执行机械臂自动调平
// @Description 控制机械臂执行自动调平操作。自动调平功能用于确保机械臂的工作平面
// @Description 保持水平，这对于某些精密操作（如3D打印、点胶等）非常重要。调平过程
// @Description 会自动测量和补偿工作平面的倾斜度。
// @Param cmd *AutoLevelingCmd true "自动调平命令参数结构体，包含：
//   - enable: 是否启用自动调平
//   - accuracy: 调平精度要求，单位mm
//   - maxHeight: 最大调平高度，单位mm"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 将指令加入队列，按顺序执行
//   - false: 立即执行该指令
//   - 建议设置为true以确保运动顺序"
//
// @Return uint64 "指令队列索引（当isQueued为true时有效）"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "执行失败，可能的错误：
//   - 参数无效
//   - 精度要求过高
//   - 超出最大调平高度
//   - 机械臂被锁定
//   - 机械臂处于报警状态
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	cmd := &AutoLevelingCmd{
//	    enable: true,
//	    accuracy: 0.02,    // 精度0.02mm
//	    maxHeight: 20.0,   // 最大调平高度20mm
//	}
//	index, err := dobot.SetAutoLevelingCmd(cmd, true)
//	if err != nil {
//	    log.Printf("执行自动调平失败: %v", err)
//	} else {
//	    log.Printf("正在执行自动调平，指令索引: %d", index)
//	}
func (dobot *Dobot) SetAutoLevelingCmd(cmd *AutoLevelingCmd, isQueued bool) (queuedCmdIndex uint64, err error) {
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
// @Summary 获取自动调平执行结果
// @Description 获取最近一次自动调平操作的执行结果。返回的精度值表示当前工作平面
// @Description 相对于理想水平面的偏差程度。此结果可用于评估调平是否达到要求。
// @Return float32 "调平精度值：
//   - 单位：mm
//   - 值越小表示调平效果越好
//   - 0表示完全水平"
//
// @Return error "错误信息"
// @Success 200 {number} float32 "返回调平精度值"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 未执行过调平操作
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	precision, err := dobot.GetAutoLevelingResult()
//	if err != nil {
//	    log.Printf("获取调平结果失败: %v", err)
//	} else {
//	    log.Printf("调平精度: %.3fmm", precision)
//	    if precision < 0.05 {
//	        log.Printf("调平效果良好")
//	    } else {
//	        log.Printf("建议重新调平")
//	    }
//	}
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
// @Summary 设置手持示教器触发模式
// @Description 设置机械臂手持示教器(HHT)的触发模式。手持示教器是一个用于手动控制
// @Description 机械臂的设备，通过不同的触发模式可以实现不同的控制效果。
// @Description 触发模式会影响示教器按键的响应方式。
// @Param mode HHTTrigMode true "触发模式：
//   - HHTTrigMode_DISABLE: 禁用触发
//   - HHTTrigMode_IMMEDIATELY: 立即触发
//   - HHTTrigMode_DELAY: 延时触发
//   - 具体行为参见HHTTrigMode枚举定义"
//
// @Return error "错误信息"
// @Success 200 {string} "设置成功返回空字符串"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 无效的触发模式
//   - 示教器未连接
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	err := dobot.SetHHTTrigMode(HHTTrigMode_IMMEDIATELY)
//	if err != nil {
//	    log.Printf("设置触发模式失败: %v", err)
//	} else {
//	    log.Printf("已设置为立即触发模式")
//	}
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
// @Summary 获取手持示教器当前触发模式
// @Description 获取机械臂手持示教器(HHT)当前设置的触发模式。可以用于确认
// @Description 示教器的工作状态和响应方式。
// @Return HHTTrigMode "当前触发模式：
//   - HHTTrigMode_DISABLE: 禁用状态
//   - HHTTrigMode_IMMEDIATELY: 立即触发状态
//   - HHTTrigMode_DELAY: 延时触发状态"
//
// @Return error "错误信息"
// @Success 200 {string} HHTTrigMode "返回当前触发模式"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 示教器未连接
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	mode, err := dobot.GetHHTTrigMode()
//	if err != nil {
//	    log.Printf("获取触发模式失败: %v", err)
//	} else {
//	    switch mode {
//	    case HHTTrigMode_DISABLE:
//	        log.Printf("当前为禁用状态")
//	    case HHTTrigMode_IMMEDIATELY:
//	        log.Printf("当前为立即触发状态")
//	    case HHTTrigMode_DELAY:
//	        log.Printf("当前为延时触发状态")
//	    }
//	}
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
// @Summary 设置手持示教器触发输出使能状态
// @Description 控制机械臂手持示教器(HHT)的触发输出功能的启用状态。启用后，
// @Description 示教器的触发动作可以产生相应的输出信号，用于控制外部设备或
// @Description 协调其他动作。
// @Param enabled bool true "使能状态：
//   - true: 启用触发输出
//   - false: 禁用触发输出"
//
// @Return error "错误信息"
// @Success 200 {string} "设置成功"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 示教器未连接
//   - 通信错误
//   - 设备未连接
//   - 当前模式不支持该操作"
//
// @Example
//
//	err := dobot.SetHHTTrigOutputEnabled(true)
//	if err != nil {
//	    log.Printf("设置触发输出使能失败: %v", err)
//	} else {
//	    log.Printf("已启用触发输出功能")
//	}
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
// @Summary 获取手持示教器触发输出使能状态
// @Description 获取机械臂手持示教器(HHT)的触发输出功能当前是否处于启用状态。
// @Description 可用于确认示教器是否能够产生触发输出信号。
// @Return bool "使能状态：
//   - true: 触发输出已启用
//   - false: 触发输出已禁用"
//
// @Return error "错误信息"
// @Success 200 {boolean} bool "返回使能状态"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 示教器未连接
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	enabled, err := dobot.GetHHTTrigOutputEnabled()
//	if err != nil {
//	    log.Printf("获取触发输出使能状态失败: %v", err)
//	} else {
//	    if enabled {
//	        log.Printf("触发输出功能已启用")
//	    } else {
//	        log.Printf("触发输出功能已禁用")
//	    }
//	}
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

// GetHHTTrigOutput 获取手持示教触发输出状态
// @Summary 获取手持示教器触发输出状态
// @Description 获取机械臂手持示教器(HHT)当前的触发输出信号状态。此状态表示
// @Description 示教器是否正在输出触发信号。注意这与使能状态不同，使能表示
// @Description 功能是否可用，而输出状态表示当前是否有实际的触发信号输出。
// @Return bool "输出状态：
//   - true: 当前有触发信号输出
//   - false: 当前无触发信号输出"
//
// @Return error "错误信息"
// @Success 200 {boolean} bool "返回触发输出状态"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 示教器未连接
//   - 触发输出未使能
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	output, err := dobot.GetHHTTrigOutput()
//	if err != nil {
//	    log.Printf("获取触发输出状态失败: %v", err)
//	} else {
//	    if output {
//	        log.Printf("当前有触发信号输出")
//	    } else {
//	        log.Printf("当前无触发信号输出")
//	    }
//	}
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
// @Summary 设置末端执行器参数
// @Description 设置机械臂末端执行器的基本参数。末端执行器是机械臂的工作工具，
// @Description 如夹爪、吸盘等。正确设置这些参数对确保末端执行器的精确控制
// @Description 和安全运行至关重要。
// @Param params *EndEffectorParams true "末端执行器参数结构体，包含：
//   - xBias: X轴偏移量，单位mm
//   - yBias: Y轴偏移量，单位mm
//   - zBias: Z轴偏移量，单位mm
//   - 注意：偏移量是相对于机械臂标准工具坐标系的补偿值"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 将指令加入队列，按顺序执行
//   - false: 立即执行该指令"
//
// @Return uint64 "指令队列索引（当isQueued为true时有效）"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数无效
//   - 偏移量超出范围
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	params := &EndEffectorParams{
//	    xBias: 0,      // X轴无偏移
//	    yBias: 10.0,   // Y轴偏移10mm
//	    zBias: 20.0,   // Z轴偏移20mm
//	}
//	index, err := dobot.SetEndEffectorParams(params, true)
//	if err != nil {
//	    log.Printf("设置末端执行器参数失败: %v", err)
//	} else {
//	    log.Printf("设置成功，指令索引: %d", index)
//	}
func (dobot *Dobot) SetEndEffectorParams(params *EndEffectorParams, isQueued bool) (queuedCmdIndex uint64, err error) {
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
// @Summary 获取末端执行器参数
// @Description 获取机械臂末端执行器当前的参数设置。这些参数定义了末端执行器
// @Description 相对于机械臂标准工具坐标系的位置补偿值。可用于确认或验证
// @Description 末端执行器的设置是否正确。
// @Return *EndEffectorParams "末端执行器参数结构体，包含：
//   - xBias: X轴偏移量，单位mm
//   - yBias: Y轴偏移量，单位mm
//   - zBias: Z轴偏移量，单位mm
//   - 注意：偏移量是相对于机械臂标准工具坐标系的补偿值"
//
// @Return error "错误信息"
// @Success 200 {object} EndEffectorParams "返回末端执行器参数结构体"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效
//   - 数据解析错误"
//
// @Example
//
//	params, err := dobot.GetEndEffectorParams()
//	if err != nil {
//	    log.Printf("获取末端执行器参数失败: %v", err)
//	} else {
//	    log.Printf("末端执行器偏移量: X=%.2f, Y=%.2f, Z=%.2f",
//	        params.xBias, params.yBias, params.zBias)
//	}
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
	params := &EndEffectorParams{}
	if err := binary.Read(bytes.NewReader(resp.Params), binary.LittleEndian, params); err != nil {
		return nil, fmt.Errorf("failed to read END_EFFECTOR_PARAMS: %v", err)
	}
	return params, nil
}

// SetEndEffectorLaser 设置末端执行器激光
// @Summary 设置末端激光器状态
// @Description 控制机械臂末端激光器的开关状态。激光器可用于切割、雕刻、标记等
// @Description 工作。使用激光器时需要特别注意安全，确保激光不会对人员或设备
// @Description 造成伤害。建议在使用完毕后立即关闭激光器。
// @Param enableCtrl bool true "激光器控制使能：
//   - true: 启用激光器控制
//   - false: 禁用激光器控制
//   - 注意：必须先启用控制才能操作激光器"
//
// @Param on bool true "激光器开关状态：
//   - true: 打开激光
//   - false: 关闭激光
//   - 注意：只有在控制使能时才有效"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 将指令加入队列，按顺序执行
//   - false: 立即执行该指令
//   - 建议使用队列模式以确保操作顺序"
//
// @Return uint64 "指令队列索引（当isQueued为true时有效）"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 激光器未连接或未识别
//   - 控制未使能
//   - 安全锁定状态
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 启用激光器控制并打开激光
//	index, err := dobot.SetEndEffectorLaser(true, true, true)
//	if err != nil {
//	    log.Printf("设置激光器失败: %v", err)
//	} else {
//	    log.Printf("激光器已开启，指令索引: %d", index)
//	    // 执行激光操作（建议添加安全延时）
//	    time.Sleep(2 * time.Second)
//	    // 完成后关闭激光
//	    _, err = dobot.SetEndEffectorLaser(true, false, true)
//	}
func (dobot *Dobot) SetEndEffectorLaser(enableCtrl bool, on bool, isQueued bool) (queuedCmdIndex uint64, err error) {
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
// @Summary 获取末端激光器状态
// @Description 获取当前末端激光器的控制和开关状态
// @Success 200 {boolean} isCtrlEnabled "是否使能激光器控制"
// @Success 200 {boolean} isOn "激光器是否开启"
// @Failure 400 {object} error "获取激光器状态失败时返回错误信息"
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
// @Summary 设置末端吸盘
// @Description 控制末端吸盘的使能和吸取状态
// @Param enableCtrl query bool true "是否启用吸盘控制"
// @Param suck query bool true "吸盘状态，true为吸附"
// @Param isQueued query bool true "是否队列执行"
// @Success 200 {number} uint64 "返回命令队列索引"
// @Failure 400 {object} error "设置吸盘状态失败时返回错误信息"
func (dobot *Dobot) SetEndEffectorSuctionCup(enableCtrl bool, suck bool, isQueued bool) (queuedCmdIndex uint64, err error) {
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
// @Summary 获取末端吸盘状态
// @Description 获取末端吸盘的控制及吸附状态
// @Success 200 {boolean} isCtrlEnabled "是否使能吸盘控制"
// @Success 200 {boolean} isSucked "吸盘是否已吸取"
// @Failure 400 {object} error "获取吸盘状态失败时返回错误信息"
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

// SetEndEffectorGripper 设置末端夹爪状态
// @Summary 设置末端夹爪的控制和夹持状态
// @Description 控制机械臂末端夹爪的使能和夹持状态。夹爪可用于抓取、搬运
// @Description 物品。使用夹爪时需要注意物品尺寸和重量不要超过夹爪规格，并确保
// @Description 夹持力度适中。建议在完成操作后将夹爪恢复到安全位置。
//
// @Param enableCtrl bool true "夹爪控制使能：
//   - true: 启用夹爪控制
//   - false: 禁用夹爪控制
//   - 注意：必须先启用控制才能操作夹爪"
//
// @Param grip bool true "夹爪状态：
//   - true: 闭合夹爪
//   - false: 打开夹爪
//   - 注意：只有在控制使能时才有效"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 将指令加入队列，按顺序执行
//   - false: 立即执行该指令
//   - 建议使用队列模式以确保操作顺序"
//
// @Return uint64 "指令队列索引（当isQueued为true时有效）"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 夹爪未连接或未识别
//   - 控制未使能
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 启用夹爪控制并闭合夹爪
//	index, err := dobot.SetEndEffectorGripper(true, true, true)
//	if err != nil {
//	    log.Printf("设置夹爪失败: %v", err)
//	} else {
//	    log.Printf("夹爪已闭合，指令索引: %d", index)
//	    // 执行夹持操作（建议添加适当延时确保夹持稳定）
//	    time.Sleep(500 * time.Millisecond)
//	    // 移动物品后打开夹爪
//	    _, err = dobot.SetEndEffectorGripper(true, false, true)
//	}
func (dobot *Dobot) SetEndEffectorGripper(enableCtrl bool, grip bool, isQueued bool) (queuedCmdIndex uint64, err error) {
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
// @Summary 获取末端夹爪状态
// @Description 获取机械臂末端夹爪的当前状态，包括控制使能状态和夹持状态。
// @Description 此函数可用于确认夹爪的工作状态，以及在执行夹爪相关操作前进行
// @Description 状态检查，确保操作安全。
//
// @Return isCtrlEnabled bool "夹爪控制使能状态：
//   - true: 控制已使能
//   - false: 控制未使能"
//
// @Return isGripped bool "夹爪状态：
//   - true: 夹爪已闭合
//   - false: 夹爪已打开"
//
// @Return error "错误信息"
// @Success 200 {object} struct{isCtrlEnabled bool, isGripped bool} "成功返回夹爪状态"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 夹爪未连接或未识别
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	// 获取夹爪状态
//	isEnabled, isGripped, err := dobot.GetEndEffectorGripper()
//	if err != nil {
//	    log.Printf("获取夹爪状态失败: %v", err)
//	} else {
//	    log.Printf("夹爪状态 - 控制使能: %v, 夹持状态: %v", isEnabled, isGripped)
//	    if isEnabled && isGripped {
//	        log.Printf("注意：夹爪当前处于闭合状态")
//	    }
//	}
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
// @Summary 设置机械臂的工作臂方向
// @Description 设置机械臂的工作臂方向，用于控制机械臂的姿态。不同的臂方向
// @Description 会影响机械臂的工作空间和运动特性。在某些应用场景下，合适的
// @Description 臂方向可以提高工作效率和避免干涉。
//
// @Param armOrientation ArmOrientation true "机械臂方向：
//   - ArmOrientation_LeftyArmOrientation: 左手方向
//   - ArmOrientation_RightyArmOrientation: 右手方向
//     注意：改变臂方向可能导致机械臂运动到新位置"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 将指令加入队列，按顺序执行
//   - false: 立即执行该指令
//   - 建议使用队列模式以确保操作顺序"
//
// @Return uint64 "指令队列索引（当isQueued为true时有效）"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 无效的臂方向参数
//   - 当前姿态无法切换到目标臂方向
//   - 机械臂被锁定
//   - 机械臂处于报警状态
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置为左手方向
//	index, err := dobot.SetArmOrientation(ArmOrientation_LeftyArmOrientation, true)
//	if err != nil {
//	    log.Printf("设置机械臂方向失败: %v", err)
//	} else {
//	    log.Printf("正在切换到左手方向，指令索引: %d", index)
//	    // 等待切换完成
//	    time.Sleep(2 * time.Second)
//	}
func (dobot *Dobot) SetArmOrientation(armOrientation ArmOrientation, isQueued bool) (queuedCmdIndex uint64, err error) {
	message := &Message{
		Id:       ProtocolArmOrientation,
		RW:       true,
		IsQueued: isQueued,
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
// @Summary 获取机械臂当前的工作臂方向
// @Description 获取机械臂当前的工作臂方向设置。此函数可用于确认机械臂的
// @Description 当前姿态配置，在执行运动指令前进行状态检查，或在切换臂方向
// @Description 后验证设置是否生效。
//
// @Return ArmOrientation "机械臂方向：
//   - ArmOrientation_LeftyArmOrientation: 左手方向
//   - ArmOrientation_RightyArmOrientation: 右手方向"
//
// @Return error "错误信息"
// @Success 200 {string} ArmOrientation "返回当前臂方向"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	// 获取当前臂方向
//	orientation, err := dobot.GetArmOrientation()
//	if err != nil {
//	    log.Printf("获取机械臂方向失败: %v", err)
//	} else {
//	    if orientation == ArmOrientation_LeftyArmOrientation {
//	        log.Printf("当前为左手方向")
//	    } else {
//	        log.Printf("当前为右手方向")
//	    }
//	}
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
// @Summary 设置JOG模式下的关节运动参数
// @Description 设置机械臂在JOG（点动）模式下各关节的运动参数。这些参数
// @Description 包括各关节的速度和加速度，直接影响机械臂在手动点动时的运动
// @Description 特性。合理的参数设置可以确保运动平稳且可控。
//
// @Param params *JOGJointParams true "JOG关节运动参数：
//   - velocity: 各关节速度数组[4]float32（单位：°/s）
//   - acceleration: 各关节加速度数组[4]float32（单位：°/s²）
//     注意：参数设置不合理可能导致运动不稳定"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 将指令加入队列，按顺序执行
//   - false: 立即执行该指令
//   - 建议使用队列模式以确保参数设置顺序"
//
// @Return uint64 "指令队列索引（当isQueued为true时有效）"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数为空
//   - 参数值超出范围
//   - 机械臂被锁定
//   - 机械臂处于报警状态
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置JOG关节运动参数
//	params := &JOGJointParams{
//	    Velocity:     [4]float32{10, 10, 10, 10},      // 各关节速度10°/s
//	    Acceleration: [4]float32{50, 50, 50, 50},      // 各关节加速度50°/s²
//	}
//	index, err := dobot.SetJOGJointParams(params, true)
//	if err != nil {
//	    log.Printf("设置JOG关节参数失败: %v", err)
//	} else {
//	    log.Printf("JOG关节参数设置成功，指令索引: %d", index)
//	}
func (dobot *Dobot) SetJOGJointParams(params *JOGJointParams, isQueued bool) (queuedCmdIndex uint64, err error) {
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
// @Summary 获取JOG模式下的关节运动参数
// @Description 获取机械臂在JOG（点动）模式下各关节的当前运动参数设置。
// @Description 可用于确认当前的运动参数配置，或在修改参数前获取原始值作为
// @Description 参考。返回的参数包括各关节的速度和加速度设置。
//
// @Return *JOGJointParams "JOG关节运动参数：
//   - velocity: 各关节速度数组[4]float32（单位：°/s）
//   - acceleration: 各关节加速度数组[4]float32（单位：°/s²）"
//
// @Return error "错误信息"
// @Success 200 {object} *JOGJointParams "返回JOG关节运动参数结构体"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效
//   - 数据解析错误"
//
// @Example
//
//	// 获取当前JOG关节运动参数
//	params, err := dobot.GetJOGJointParams()
//	if err != nil {
//	    log.Printf("获取JOG关节参数失败: %v", err)
//	} else {
//	    log.Printf("当前JOG关节参数：")
//	    log.Printf("  速度: %.2f, %.2f, %.2f, %.2f °/s",
//	        params.Velocity[0], params.Velocity[1],
//	        params.Velocity[2], params.Velocity[3])
//	    log.Printf("  加速度: %.2f, %.2f, %.2f, %.2f °/s²",
//	        params.Acceleration[0], params.Acceleration[1],
//	        params.Acceleration[2], params.Acceleration[3])
//	}
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
	if err := binary.Read(bytes.NewReader(resp.Params), binary.LittleEndian, params); err != nil {
		return nil, fmt.Errorf("failed to read JOG joint params: %v", err)
	}
	return params, nil
}

// SetJOGCoordinateParams 设置JOG坐标参数
// @Summary 设置JOG模式下的笛卡尔坐标运动参数
// @Description 设置机械臂在JOG（点动）模式下笛卡尔坐标系中的运动参数。
// @Description 这些参数包括X、Y、Z轴的速度和加速度，以及姿态（R轴）的运动
// @Description 参数，直接影响机械臂在手动点动时的运动特性。
//
// @Param params *JOGCoordinateParams true "JOG坐标运动参数：
//   - velocity: 各轴速度数组[4]float32
//   - [0-2]: X、Y、Z轴速度（单位：mm/s）
//   - [3]: R轴速度（单位：°/s）
//   - acceleration: 各轴加速度数组[4]float32
//   - [0-2]: X、Y、Z轴加速度（单位：mm/s²）
//   - [3]: R轴加速度（单位：°/s²）
//     注意：参数设置不合理可能导致运动不稳定"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 将指令加入队列，按顺序执行
//   - false: 立即执行该指令
//   - 建议使用队列模式以确保参数设置顺序"
//
// @Return uint64 "指令队列索引（当isQueued为true时有效）"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数为空
//   - 参数值超出范围
//   - 机械臂被锁定
//   - 机械臂处于报警状态
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置JOG坐标运动参数
//	params := &JOGCoordinateParams{
//	    Velocity: [4]float32{
//	        50, 50, 50,  // XYZ轴速度50mm/s
//	        30,          // R轴速度30°/s
//	    },
//	    Acceleration: [4]float32{
//	        100, 100, 100,  // XYZ轴加速度100mm/s²
//	        50,             // R轴加速度50°/s²
//	    },
//	}
//	index, err := dobot.SetJOGCoordinateParams(params, true)
//	if err != nil {
//	    log.Printf("设置JOG坐标参数失败: %v", err)
//	} else {
//	    log.Printf("JOG坐标参数设置成功，指令索引: %d", index)
//	}
func (dobot *Dobot) SetJOGCoordinateParams(params *JOGCoordinateParams, isQueued bool) (queuedCmdIndex uint64, err error) {
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
// @Summary 获取JOG模式下的笛卡尔坐标运动参数
// @Description 获取机械臂在JOG（点动）模式下笛卡尔坐标系中的当前运动参数
// @Description 设置。可用于确认当前的运动参数配置，或在修改参数前获取原始
// @Description 值作为参考。返回的参数包括各轴的速度和加速度设置。
//
// @Return *JOGCoordinateParams "JOG坐标运动参数：
//   - velocity: 各轴速度数组[4]float32
//   - [0-2]: X、Y、Z轴速度（单位：mm/s）
//   - [3]: R轴速度（单位：°/s）
//   - acceleration: 各轴加速度数组[4]float32
//   - [0-2]: X、Y、Z轴加速度（单位：mm/s²）
//   - [3]: R轴加速度（单位：°/s²）"
//
// @Return error "错误信息"
// @Success 200 {object} *JOGCoordinateParams "返回JOG坐标运动参数结构体"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效
//   - 数据解析错误"
//
// @Example
//
//	// 获取当前JOG坐标运动参数
//	params, err := dobot.GetJOGCoordinateParams()
//	if err != nil {
//	    log.Printf("获取JOG坐标参数失败: %v", err)
//	} else {
//	    log.Printf("当前JOG坐标参数：")
//	    log.Printf("  XYZ轴速度: %.2f, %.2f, %.2f mm/s",
//	        params.Velocity[0], params.Velocity[1], params.Velocity[2])
//	    log.Printf("  R轴速度: %.2f °/s", params.Velocity[3])
//	    log.Printf("  XYZ轴加速度: %.2f, %.2f, %.2f mm/s²",
//	        params.Acceleration[0], params.Acceleration[1], params.Acceleration[2])
//	    log.Printf("  R轴加速度: %.2f °/s²", params.Acceleration[3])
//	}
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
	if err := binary.Read(bytes.NewReader(resp.Params), binary.LittleEndian, params); err != nil {
		return nil, fmt.Errorf("failed to read JOG coordinate params: %v", err)
	}
	return params, nil
}

// SetJOGLParams 设置JOGL参数
// @Summary 设置JOG模式下的连续关节角度增量参数
// @Description 设置机械臂在JOG（点动）模式下的连续关节角度增量运动参数。
// @Description 这些参数用于控制机械臂在连续运动时的角度增量变化，影响运动
// @Description 的平滑度和精确度。注意此函数的设置为立即生效，不支持队列执行。
//
// @Param params *JOGLParams true "JOGL运动参数：
//   - velocity: 关节角速度（单位：°/s）
//   - acceleration: 关节角加速度（单位：°/s²）
//     注意：参数设置不合理可能导致运动不稳定或不平滑"
//
// @Return uint64 "指令队列索引（此函数总是返回0，因为不支持队列执行）"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回0"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数为空
//   - 参数值超出范围
//   - 机械臂被锁定
//   - 机械臂处于报警状态
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置JOGL运动参数
//	params := &JOGLParams{
//	    Velocity:     20.0,  // 角速度20°/s
//	    Acceleration: 50.0,  // 角加速度50°/s²
//	}
//	_, err := dobot.SetJOGLParams(params)
//	if err != nil {
//	    log.Printf("设置JOGL参数失败: %v", err)
//	} else {
//	    log.Printf("JOGL参数设置成功")
//	}
func (dobot *Dobot) SetJOGLParams(params *JOGLParams) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolJOGLParams,
		RW:       true,
		IsQueued: false, // C++实现中强制为false
	}

	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()

	_, err = dobot.connector.SendMessage(context.Background(), message)
	return 0, err // 由于IsQueued为false，所以不返回queuedCmdIndex
}

// GetJOGLParams 获取JOGL参数
// @Summary 获取JOG模式下的连续关节角度增量参数
// @Description 获取机械臂在JOG（点动）模式下的连续关节角度增量运动参数
// @Description 设置。可用于确认当前的运动参数配置，或在修改参数前获取原始
// @Description 值作为参考。
//
// @Return *JOGLParams "JOGL运动参数：
//   - velocity: 关节角速度（单位：°/s）
//   - acceleration: 关节角加速度（单位：°/s²）"
//
// @Return error "错误信息"
// @Success 200 {object} *JOGLParams "返回JOGL运动参数结构体"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效
//   - 数据解析错误"
//
// @Example
//
//	// 获取当前JOGL运动参数
//	params, err := dobot.GetJOGLParams()
//	if err != nil {
//	    log.Printf("获取JOGL参数失败: %v", err)
//	} else {
//	    log.Printf("当前JOGL参数：")
//	    log.Printf("  角速度: %.2f °/s", params.Velocity)
//	    log.Printf("  角加速度: %.2f °/s²", params.Acceleration)
//	}
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

	params := &JOGLParams{}
	if err := binary.Read(bytes.NewReader(resp.Params), binary.LittleEndian, params); err != nil {
		return nil, fmt.Errorf("failed to read JOG L params: %v", err)
	}
	return params, nil
}

// SetJOGCommonParams 设置JOG通用参数
// @Summary 设置JOG模式下的通用运动参数
// @Description 设置机械臂在JOG（点动）模式下所有运动方式共用的基础参数。
// @Description 这些参数会影响所有JOG运动模式（包括关节运动、坐标运动等）的
// @Description 基本特性，如运动的平滑度和响应性。
//
// @Param params *JOGCommonParams true "JOG通用参数：
//   - velocityRatio: 速度比例，范围[0-100]
//   - accelerationRatio: 加速度比例，范围[0-100]
//     注意：参数设置过大可能导致运动不稳定，建议从小到大逐步调整"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 将指令加入队列，按顺序执行
//   - false: 立即执行该指令
//   - 建议使用队列模式以确保参数设置顺序"
//
// @Return uint64 "指令队列索引（当isQueued为true时有效）"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数为空
//   - 参数值超出范围
//   - 机械臂被锁定
//   - 机械臂处于报警状态
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置JOG通用运动参数
//	params := &JOGCommonParams{
//	    VelocityRatio:     50,  // 速度比例50%
//	    AccelerationRatio: 50,  // 加速度比例50%
//	}
//	index, err := dobot.SetJOGCommonParams(params, true)
//	if err != nil {
//	    log.Printf("设置JOG通用参数失败: %v", err)
//	} else {
//	    log.Printf("JOG通用参数设置成功，指令索引: %d", index)
//	}
func (dobot *Dobot) SetJOGCommonParams(params *JOGCommonParams, isQueued bool) (queuedCmdIndex uint64, err error) {
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
// @Summary 获取JOG模式下的通用运动参数
// @Description 获取机械臂在JOG（点动）模式下所有运动方式共用的基础参数
// @Description 设置。可用于确认当前的运动参数配置，或在修改参数前获取原始
// @Description 值作为参考。
//
// @Return *JOGCommonParams "JOG通用参数：
//   - velocityRatio: 速度比例，范围[0-100]
//   - accelerationRatio: 加速度比例，范围[0-100]"
//
// @Return error "错误信息"
// @Success 200 {object} *JOGCommonParams "返回JOG通用参数结构体"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效
//   - 数据解析错误"
//
// @Example
//
//	// 获取当前JOG通用运动参数
//	params, err := dobot.GetJOGCommonParams()
//	if err != nil {
//	    log.Printf("获取JOG通用参数失败: %v", err)
//	} else {
//	    log.Printf("当前JOG通用参数：")
//	    log.Printf("  速度比例: %d%%", params.VelocityRatio)
//	    log.Printf("  加速度比例: %d%%", params.AccelerationRatio)
//	}
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

	params := &JOGCommonParams{}
	if err := binary.Read(bytes.NewReader(resp.Params), binary.LittleEndian, params); err != nil {
		return nil, fmt.Errorf("failed to read JOG common params: %v", err)
	}
	return params, nil
}

// SetJOGCmd 设置JOG运动指令
// @Summary 发送JOG模式的单步运动指令
// @Description 控制机械臂在JOG（点动）模式下执行单步运动。可以选择关节运动
// @Description 或笛卡尔坐标运动模式，并指定运动方向。此指令用于手动控制机械
// @Description 臂的精确移动，常用于示教和位置微调。
//
// @Param cmd *JOGCmd true "JOG运动指令参数结构体，包含：
//   - isJoint: 运动模式选择（true为关节模式/false为坐标模式）
//   - index: 运动轴索引（关节模式：0-3为关节1-4；坐标模式：0-3为X/Y/Z/R轴）
//   - direction: 运动方向（1正向/0停止/-1负向）"
//
// @Param isQueued bool true "是否加入指令队列（true加入队列/false立即执行）"
// @Return uint64 "指令队列索引（当isQueued为true时有效）"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "执行失败，可能的错误：
//   - 参数为空
//   - 参数值无效
//   - 机械臂被锁定
//   - 机械臂处于报警状态
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 执行关节1正向运动
//	cmd := &JOGCmd{
//	    IsJoint:   true,     // 关节运动模式
//	    Index:     0,        // 关节1
//	    Direction: 1,        // 正向运动
//	}
//	index, err := dobot.SetJOGCmd(cmd, true)
//	if err != nil {
//	    log.Printf("执行JOG运动失败: %v", err)
//	} else {
//	    log.Printf("正在执行JOG运动，指令索引: %d", index)
//	    // 运动一段时间后停止
//	    time.Sleep(1 * time.Second)
//	    cmd.Direction = 0    // 停止运动
//	    _, err = dobot.SetJOGCmd(cmd, true)
//	}
func (dobot *Dobot) SetJOGCmd(cmd *JOGCmd, isQueued bool) (queuedCmdIndex uint64, err error) {
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
// @Summary 设置PTP模式下的关节运动参数
// @Description 设置机械臂在PTP（点到点）模式下各关节的运动参数。这些参数
// @Description 包括各关节的速度和加速度，直接影响机械臂在点到点运动时的运动
// @Description 特性。合理的参数设置可以确保运动平滑且高效。
//
// @Param params *PTPJointParams true "PTP关节运动参数：
//   - velocity: 各关节速度数组[4]float32（单位：°/s）
//   - acceleration: 各关节加速度数组[4]float32（单位：°/s²）
//     注意：参数设置不合理可能导致运动不稳定或耗时过长"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 将指令加入队列，按顺序执行
//   - false: 立即执行该指令
//   - 建议使用队列模式以确保参数设置顺序"
//
// @Return uint64 "指令队列索引（当isQueued为true时有效）"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数为空
//   - 参数值超出范围
//   - 机械臂被锁定
//   - 机械臂处于报警状态
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置PTP关节运动参数
//	params := &PTPJointParams{
//	    Velocity: [4]float32{
//	        60, 60, 60, 60,  // 各关节速度60°/s
//	    },
//	    Acceleration: [4]float32{
//	        100, 100, 100, 100,  // 各关节加速度100°/s²
//	    },
//	}
//	index, err := dobot.SetPTPJointParams(params, true)
//	if err != nil {
//	    log.Printf("设置PTP关节参数失败: %v", err)
//	} else {
//	    log.Printf("PTP关节参数设置成功，指令索引: %d", index)
//	}
func (dobot *Dobot) SetPTPJointParams(params *PTPJointParams, isQueued bool) (queuedCmdIndex uint64, err error) {
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
// @Summary 获取PTP模式下的关节运动参数
// @Description 获取机械臂在PTP（点到点）模式下各关节的当前运动参数设置
// @Return *PTPJointParams "PTP关节运动参数结构体"
// @Return error "错误信息"
// @Success 200 {object} *PTPJointParams "返回PTP关节运动参数结构体"
// @Failure 400 {error} "获取失败，可能的错误：通信错误、设备未连接、响应数据无效"
// @Example
//
//	params, err := dobot.GetPTPJointParams()
//	if err != nil {
//	    log.Printf("获取PTP关节参数失败: %v", err)
//	} else {
//	    log.Printf("当前关节速度: %.2f, %.2f, %.2f, %.2f °/s",
//	        params.Velocity[0], params.Velocity[1],
//	        params.Velocity[2], params.Velocity[3])
//	}
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
	if err := binary.Read(bytes.NewReader(resp.Params), binary.LittleEndian, params); err != nil {
		return nil, fmt.Errorf("failed to read PTP joint params: %v", err)
	}
	return params, nil
}

// SetPTPCoordinateParams 设置PTP坐标参数
// @Summary 设置PTP模式下的笛卡尔坐标运动参数
// @Description 设置机械臂在PTP（点到点）模式下笛卡尔坐标系中的运动参数。
// @Description 这些参数包括X、Y、Z轴的速度和加速度，以及姿态（R轴）的运动
// @Description 参数，直接影响机械臂在点到点运动时的运动特性。
//
// @Param params *PTPCoordinateParams true "PTP坐标运动参数：
//   - xyzVelocity: XYZ轴速度（单位：mm/s）
//   - rVelocity: R轴速度（单位：°/s）
//   - xyzAcceleration: XYZ轴加速度（单位：mm/s²）
//   - rAcceleration: R轴加速度（单位：°/s²）
//     注意：参数设置不合理可能导致运动不稳定或耗时过长"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 将指令加入队列，按顺序执行
//   - false: 立即执行该指令
//   - 建议使用队列模式以确保参数设置顺序"
//
// @Return uint64 "指令队列索引（当isQueued为true时有效）"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数为空
//   - 参数值超出范围
//   - 机械臂被锁定
//   - 机械臂处于报警状态
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置PTP坐标运动参数
//	params := &PTPCoordinateParams{
//	    XYZVelocity:      100,  // XYZ轴速度100mm/s
//	    RVelocity:        50,   // R轴速度50°/s
//	    XYZAcceleration:  200,  // XYZ轴加速度200mm/s²
//	    RAcceleration:    100,  // R轴加速度100°/s²
//	}
//	index, err := dobot.SetPTPCoordinateParams(params, true)
//	if err != nil {
//	    log.Printf("设置PTP坐标参数失败: %v", err)
//	} else {
//	    log.Printf("PTP坐标参数设置成功，指令索引: %d", index)
//	}
func (dobot *Dobot) SetPTPCoordinateParams(params *PTPCoordinateParams, isQueued bool) (queuedCmdIndex uint64, err error) {
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
// @Summary 获取PTP坐标参数
// @Description 获取当前PTP模式下笛卡尔坐标运动参数
// @Success 200 {object} *PTPCoordinateParams "返回PTP坐标参数结构体"
// @Failure 400 {object} error "获取PTP坐标参数失败时返回错误信息"
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

	params := &PTPCoordinateParams{}
	if err := binary.Read(bytes.NewReader(resp.Params), binary.LittleEndian, params); err != nil {
		return nil, fmt.Errorf("failed to read PTP coordinate params: %v", err)
	}
	return params, nil
}

// SetPTPLParams 设置PTPL参数
// @Summary 设置PTP模式下的线性插补运动参数
// @Description 设置机械臂在PTP（点到点）模式下的线性插补运动参数。这些参数
// @Description 用于控制机械臂在直线运动时的速度和加速度特性，确保运动轨迹的
// @Description 精确性和平滑性。
//
// @Param params *PTPLParams true "PTPL运动参数：
//   - velocity: 速度（单位：mm/s）
//   - acceleration: 加速度（单位：mm/s²）
//     注意：参数设置不合理可能导致运动不稳定或耗时过长"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 将指令加入队列，按顺序执行
//   - false: 立即执行该指令
//   - 建议使用队列模式以确保参数设置顺序"
//
// @Return uint64 "指令队列索引（当isQueued为true时有效）"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数为空
//   - 参数值超出范围
//   - 机械臂被锁定
//   - 机械臂处于报警状态
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置PTPL运动参数
//	params := &PTPLParams{
//	    Velocity:     100,  // 速度100mm/s
//	    Acceleration: 200,  // 加速度200mm/s²
//	}
//	index, err := dobot.SetPTPLParams(params, true)
//	if err != nil {
//	    log.Printf("设置PTPL参数失败: %v", err)
//	} else {
//	    log.Printf("PTPL参数设置成功，指令索引: %d", index)
//	}
func (dobot *Dobot) SetPTPLParams(params *PTPLParams, isQueued bool) (queuedCmdIndex uint64, err error) {
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
// @Summary 获取PTP模式下的线性插补运动参数
// @Description 获取机械臂在PTP（点到点）模式下的线性插补运动参数设置。
// @Description 可用于确认当前的运动参数配置，或在修改参数前获取原始值作为
// @Description 参考。返回的参数包括速度和加速度设置。
//
// @Return *PTPLParams "PTPL运动参数：
//   - velocity: 速度（单位：mm/s）
//   - acceleration: 加速度（单位：mm/s²）"
//
// @Return error "错误信息"
// @Success 200 {object} *PTPLParams "返回PTPL运动参数结构体"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效
//   - 数据解析错误"
//
// @Example
//
//	// 获取当前PTPL运动参数
//	params, err := dobot.GetPTPLParams()
//	if err != nil {
//	    log.Printf("获取PTPL参数失败: %v", err)
//	} else {
//	    log.Printf("当前PTPL参数：")
//	    log.Printf("  速度: %.2f mm/s", params.Velocity)
//	    log.Printf("  加速度: %.2f mm/s²", params.Acceleration)
//	}
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

	params := &PTPLParams{}
	if err := binary.Read(bytes.NewReader(resp.Params), binary.LittleEndian, params); err != nil {
		return nil, fmt.Errorf("failed to read PTP L params: %v", err)
	}
	return params, nil
}

// SetPTPJumpParams 设置PTP跳跃参数
// @Summary 设置PTP跳跃参数
// @Description 通过PTPJumpParams结构体设置跳跃运动参数
// @Param params body *PTPJumpParams true "PTP跳跃参数结构体"
// @Param isQueued query bool true "是否队列执行"
// @Success 200 {number} uint8 "返回命令队列索引"
// @Failure 400 {object} error "设置PTP跳跃参数失败时返回错误信息"
func (dobot *Dobot) SetPTPJumpParams(params *PTPJumpParams, isQueued bool) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolPTPJumpParams,
		RW:       true,
		IsQueued: isQueued,
	}

	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()

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
// @Summary 获取PTP跳跃参数
// @Description 获取当前PTP跳跃模式下的参数
// @Success 200 {object} *PTPJumpParams "返回PTP跳跃参数结构体"
// @Failure 400 {object} error "获取PTP跳跃参数失败时返回错误信息"
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

	params := &PTPJumpParams{}
	if err := binary.Read(bytes.NewReader(resp.Params), binary.LittleEndian, params); err != nil {
		return nil, fmt.Errorf("failed to read PTP jump params: %v", err)
	}
	return params, nil
}

// SetPTPJump2Params 设置PTP跳跃2参数
// @Summary 设置PTP跳跃2参数
// @Description 通过PTPJump2Params结构体设置跳跃2运动模式参数
// @Param params body *PTPJump2Params true "PTP跳跃2参数结构体"
// @Param isQueued query bool true "是否队列执行"
// @Success 200 {number} uint8 "返回命令队列索引"
// @Failure 400 {object} error "设置PTP跳跃2参数失败时返回错误信息"
func (dobot *Dobot) SetPTPJump2Params(params *PTPJump2Params, isQueued bool) (queuedCmdIndex uint64, err error) {
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
// @Summary 获取PTP跳跃2参数
// @Description 获取当前PTP跳跃2模式下的参数
// @Success 200 {object} *PTPJump2Params "返回PTP跳跃2参数结构体"
// @Failure 400 {object} error "获取PTP跳跃2参数失败时返回错误信息"
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

	params := &PTPJump2Params{}
	if err := binary.Read(bytes.NewReader(resp.Params), binary.LittleEndian, params); err != nil {
		return nil, fmt.Errorf("failed to read PTP jump2 params: %v", err)
	}
	return params, nil
}

// SetPTPCommonParams 设置PTP通用参数
// @Summary 设置PTP共通参数
// @Description 通过PTPCommonParams结构体设置所有PTP运动模式的共通参数
// @Param params body *PTPCommonParams true "PTP共通参数结构体"
// @Param isQueued query bool true "是否队列执行"
// @Success 200 {number} uint8 "返回命令队列索引"
// @Failure 400 {object} error "设置PTP共通参数失败时返回错误信息"
func (dobot *Dobot) SetPTPCommonParams(params *PTPCommonParams, isQueued bool) (queuedCmdIndex uint64, err error) {
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

	params := &PTPCommonParams{}
	if err := binary.Read(bytes.NewReader(resp.Params), binary.LittleEndian, params); err != nil {
		return nil, fmt.Errorf("failed to read PTP common params: %v", err)
	}
	return params, nil
}

// SetPTPCmd 设置PTP命令
// @Summary 设置PTP命令
// @Description 通过PTPCmd结构体发送PTP运动命令，控制机械臂快速移动到目标位姿
// @Param cmd body *PTPCmd true "PTP命令结构体"
// @Param isQueued query bool true "是否队列执行"
// @Success 200 {number} uint8 "返回命令队列索引"
// @Failure 400 {object} error "设置PTP命令失败时返回错误信息"
func (dobot *Dobot) SetPTPCmd(cmd *PTPCmd, isQueued bool) (queuedCmdIndex uint64, err error) {
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
// @Summary 设置带L参数的PTP命令
// @Description 通过PTPWithLCmd结构体发送带L参数的PTP命令，用于长臂补偿
// @Param cmd body *PTPWithLCmd true "带L参数的PTP命令结构体"
// @Param isQueued query bool true "是否队列执行"
// @Success 200 {number} uint8 "返回命令队列索引"
// @Failure 400 {object} error "设置PTP带L命令失败时返回错误信息"
func (dobot *Dobot) SetPTPWithLCmd(cmd *PTPWithLCmd, isQueued bool) (queuedCmdIndex uint64, err error) {
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
// @Summary 设置连续运动参数
// @Description 通过CPParams结构体设置连续运动模式下的运动参数
// @Param params body *CPParams true "CP参数结构体"
// @Param isQueued query bool true "是否队列执行"
// @Success 200 {number} uint8 "返回命令队列索引"
// @Failure 400 {object} error "设置CP参数失败时返回错误信息"
func (dobot *Dobot) SetCPParams(params *CPParams, isQueued bool) (queuedCmdIndex uint64, err error) {
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

// SetCPCmd 设置连续运动命令
// @Summary 设置CPCmd命令
// @Description 通过CPCmd结构体发送连续运动命令，定义运动路径
// @Param cmd body *CPCmd true "CPCmd命令结构体"
// @Param isQueued query bool true "是否队列执行"
// @Success 200 {number} uint8 "返回命令队列索引"
// @Failure 400 {object} error "设置CPCmd命令失败时返回错误信息"
func (dobot *Dobot) SetCPCmd(cmd *CPCmd, isQueued bool) (queuedCmdIndex uint64, err error) {
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

// SetCPLECmd 设置连续运动扩展命令
// @Summary 设置CPLE命令
// @Description 通过cpMode及坐标参数发送连续运动扩展命令
// @Param cpMode query uint8 true "连续运动模式"
// @Param x query float32 true "目标x坐标"
// @Param y query float32 true "目标y坐标"
// @Param z query float32 true "目标z坐标"
// @Param power query float32 true "功率参数"
// @Param isQueued query bool true "是否队列执行"
// @Success 200 {number} uint8 "返回命令队列索引"
// @Failure 400 {object} error "设置CPLE命令失败时返回错误信息"
func (dobot *Dobot) SetCPLECmd(cpMode uint8, x, y, z, power float32, isQueued bool) (queuedCmdIndex uint64, err error) {
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
// @Summary 设置CP运动保持使能
// @Description 设置是否启用CP模式下的保持功能
// @Param isEnable query bool true "启用为true，禁用为false"
// @Success 200 {string} "设置成功返回空字符串"
// @Failure 400 {object} error "设置CPR保持使能失败时返回错误信息"
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

// GetCPRHoldEnable 获取CP运动保持使能状态
// @Summary 获取CPR保持使能状态
// @Description 获取当前CP模式下是否启用保持功能
// @Success 200 {boolean} bool "返回使能状态"
// @Failure 400 {object} error "获取CPR保持使能状态失败时返回错误信息"
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
// @Summary 设置CP共通参数
// @Description 通过CPCommonParams结构体设置连续运动的共通参数
// @Param params body *CPCommonParams true "CP共通参数结构体"
// @Param isQueued query bool true "是否队列执行"
// @Success 200 {number} uint8 "返回命令队列索引"
// @Failure 400 {object} error "设置CP共通参数失败时返回错误信息"
func (dobot *Dobot) SetCPCommonParams(params *CPCommonParams, isQueued bool) (queuedCmdIndex uint64, err error) {
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
// @Summary 获取CP共通参数
// @Description 获取当前连续运动模式下的共通参数
// @Success 200 {object} *CPCommonParams "返回CP共通参数结构体"
// @Failure 400 {object} error "获取CP共通参数失败时返回错误信息"
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

	params := &CPCommonParams{}
	if err := binary.Read(bytes.NewReader(resp.Params), binary.LittleEndian, params); err != nil {
		return nil, fmt.Errorf("failed to read CP common params: %v", err)
	}
	return params, nil
}

// SetARCParams 设置ARC参数
// @Summary 设置ARC模式下的圆弧运动参数
// @Description 设置机械臂在ARC（圆弧）模式下的运动参数。这些参数用于控制
// @Description 机械臂在执行圆弧轨迹运动时的速度和加速度特性，直接影响运动
// @Description 的平滑度和精确度。
//
// @Param params *ARCParams true "ARC运动参数：
//   - xyzVelocity: XYZ轴速度（单位：mm/s）
//   - rVelocity: R轴速度（单位：°/s）
//   - xyzAcceleration: XYZ轴加速度（单位：mm/s²）
//   - rAcceleration: R轴加速度（单位：°/s²）
//     注意：参数设置不合理可能导致运动不稳定或轨迹偏离"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 将指令加入队列，按顺序执行
//   - false: 立即执行该指令
//   - 建议使用队列模式以确保参数设置顺序"
//
// @Return uint64 "指令队列索引（当isQueued为true时有效）"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数为空
//   - 参数值超出范围
//   - 机械臂被锁定
//   - 机械臂处于报警状态
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置ARC运动参数
//	params := &ARCParams{
//	    XYZVelocity:      100,  // XYZ轴速度100mm/s
//	    RVelocity:        50,   // R轴速度50°/s
//	    XYZAcceleration:  200,  // XYZ轴加速度200mm/s²
//	    RAcceleration:    100,  // R轴加速度100°/s²
//	}
//	index, err := dobot.SetARCParams(params, true)
//	if err != nil {
//	    log.Printf("设置ARC参数失败: %v", err)
//	} else {
//	    log.Printf("ARC参数设置成功，指令索引: %d", index)
//	}
func (dobot *Dobot) SetARCParams(params *ARCParams) (queuedCmdIndex uint64, err error) {
	if params == nil {
		return 0, errors.New("invalid params: params is nil")
	}

	message := &Message{
		Id:       ProtocolARCParams,
		RW:       true,
		IsQueued: true,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()
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

	params := &ARCParams{}
	if err := binary.Read(bytes.NewReader(resp.Params), binary.LittleEndian, params); err != nil {
		return nil, fmt.Errorf("failed to read ARC params: %v", err)
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
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, cmd)
	message.Params = writer.Bytes()

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
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, cmd)
	message.Params = writer.Bytes()

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
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, params)
	message.Params = writer.Bytes()

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
// @Summary 获取ARC模式下的通用运动参数
// @Description 获取机械臂在ARC（圆弧）模式下的通用运动参数设置。这些参数
// @Description 影响所有ARC运动的基本特性。可用于确认当前的运动参数配置，
// @Description 或在修改参数前获取原始值作为参考。
//
// @Return *ARCCommonParams "ARC通用参数：
//   - velocityRatio: 速度比例，范围[0-100]
//   - accelerationRatio: 加速度比例，范围[0-100]
//     注意：这些比例值会影响实际的运动速度和加速度"
//
// @Return error "错误信息"
// @Success 200 {object} *ARCCommonParams "返回ARC通用参数结构体"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效
//   - 数据解析错误"
//
// @Example
//
//	// 获取当前ARC通用运动参数
//	params, err := dobot.GetARCCommonParams()
//	if err != nil {
//	    log.Printf("获取ARC通用参数失败: %v", err)
//	} else {
//	    log.Printf("当前ARC通用参数：")
//	    log.Printf("  速度比例: %d%%", params.VelocityRatio)
//	    log.Printf("  加速度比例: %d%%", params.AccelerationRatio)
//	}
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

	params := &ARCCommonParams{}
	if err := binary.Read(bytes.NewReader(resp.Params), binary.LittleEndian, params); err != nil {
		return nil, fmt.Errorf("failed to read ARC common params: %v", err)
	}
	return params, nil
}

// SetWAITCmd 设置等待指令
// @Summary 设置机械臂等待指令
// @Description 控制机械臂执行等待操作。等待指令用于在运动序列中插入延时，
// @Description 可以用来协调不同动作之间的时序，或等待外部条件满足。此指令
// @Description 必须加入队列执行。
//
// @Param cmd *WAITCmd true "等待指令参数：
//   - timeout: 等待时间（单位：毫秒）
//     注意：等待时间不宜过长，以免影响整体执行效率"
//
// @Param isQueued bool true "是否加入指令队列：
//   - 此函数强制为true，必须加入队列执行
//   - 如果设为false也会被转为true"
//
// @Return uint64 "指令队列索引"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数为空
//   - 等待时间无效
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置等待1秒
//	cmd := &WAITCmd{
//	    Timeout: 1000,  // 等待1000毫秒
//	}
//	index, err := dobot.SetWAITCmd(cmd, true)
//	if err != nil {
//	    log.Printf("设置等待指令失败: %v", err)
//	} else {
//	    log.Printf("等待指令设置成功，指令索引: %d", index)
//	}
func (dobot *Dobot) SetWAITCmd(cmd *WAITCmd, isQueued bool) (queuedCmdIndex uint64, err error) {
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

// SetTRIGCmd 设置触发指令
// @Summary 设置机械臂触发指令
// @Description 控制机械臂执行触发操作。触发指令用于在运动过程中产生输出
// @Description 信号，可以用来控制外部设备或同步其他动作。触发条件可以是
// @Description 时间、位置等，具体由触发模式决定。
//
// @Param cmd *TRIGCmd true "触发指令参数：
//   - address: 触发地址（通常为IO端口号）
//   - mode: 触发模式（时间触发、位置触发等）
//   - condition: 触发条件（与模式相关的具体参数）
//   - threshold: 触发阈值
//     注意：参数设置需要与实际硬件配置相匹配"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 将指令加入队列，按顺序执行
//   - false: 立即执行该指令
//   - 建议使用队列模式以确保触发时序"
//
// @Return uint64 "指令队列索引（当isQueued为true时有效）"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数为空
//   - 触发参数无效
//   - 地址超出范围
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置时间触发指令
//	cmd := &TRIGCmd{
//	    Address:   1,        // IO端口1
//	    Mode:      0,        // 时间触发模式
//	    Condition: 500,      // 延时500ms后触发
//	    Threshold: 1,        // 输出高电平
//	}
//	index, err := dobot.SetTRIGCmd(cmd, true)
//	if err != nil {
//	    log.Printf("设置触发指令失败: %v", err)
//	} else {
//	    log.Printf("触发指令设置成功，指令索引: %d", index)
//	}
func (dobot *Dobot) SetTRIGCmd(cmd *TRIGCmd, isQueued bool) (queuedCmdIndex uint64, err error) {
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

// SetIOMultiplexing 设置IO复用功能
// @Summary 设置机械臂IO端口的复用功能
// @Description 配置机械臂IO端口的复用功能。通过复用功能，可以让一个IO端口
// @Description 实现多种不同的功能，如普通IO、PWM输出、编码器输入等。选择
// @Description 合适的复用功能可以更灵活地控制外部设备。
//
// @Param params *IOMultiplexing true "IO复用参数：
//   - address: IO端口地址
//   - multiplex: 复用功能代码，可选值：
//   - 0: 普通IO功能
//   - 1: PWM输出功能
//   - 2: 编码器输入功能
//   - 其他值参见具体型号说明
//     注意：不是所有IO端口都支持所有复用功能"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 将指令加入队列，按顺序执行
//   - false: 立即执行该指令
//   - 建议使用队列模式以确保设置顺序"
//
// @Return uint64 "指令队列索引（当isQueued为true时有效）"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数为空
//   - IO地址无效
//   - 复用功能不支持
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置IO端口1为PWM输出功能
//	params := &IOMultiplexing{
//	    Address:    1,    // IO端口1
//	    Multiplex:  1,    // PWM输出功能
//	}
//	index, err := dobot.SetIOMultiplexing(params, true)
//	if err != nil {
//	    log.Printf("设置IO复用功能失败: %v", err)
//	} else {
//	    log.Printf("IO复用功能设置成功，指令索引: %d", index)
//	}
func (dobot *Dobot) SetIOMultiplexing(params *IOMultiplexing, isQueued bool) (queuedCmdIndex uint64, err error) {
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
// @Summary 设置机械臂IO端口的数字输出状态
// @Description 控制机械臂IO端口的数字输出状态。数字输出可以用来控制外部
// @Description 设备，如电磁阀、继电器、指示灯等。在使用前需要确保IO端口
// @Description 已经配置为数字输出模式。
//
// @Param params *IODO true "IO数字输出参数：
//   - address: IO端口地址
//   - level: 输出电平
//   - 0: 低电平
//   - 1: 高电平
//     注意：输出前确保端口已正确配置为数字输出模式"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 将指令加入队列，按顺序执行
//   - false: 立即执行该指令
//   - 建议使用队列模式以确保输出时序"
//
// @Return uint64 "指令队列索引（当isQueued为true时有效）"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数为空
//   - IO地址无效
//   - 端口模式不正确
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置IO端口1输出高电平
//	params := &IODO{
//	    Address: 1,    // IO端口1
//	    Level:   1,    // 输出高电平
//	}
//	index, err := dobot.SetIODO(params, true)
//	if err != nil {
//	    log.Printf("设置IO数字输出失败: %v", err)
//	} else {
//	    log.Printf("IO数字输出设置成功，指令索引: %d", index)
//	}
func (dobot *Dobot) SetIODO(params *IODO, isQueued bool) (queuedCmdIndex uint64, err error) {
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
// @Summary 设置机械臂IO端口的PWM输出参数
// @Description 控制机械臂IO端口的PWM（脉冲宽度调制）输出。PWM输出可以用来
// @Description 控制电机速度、LED亮度、加热器功率等需要模拟量输出的场合。
// @Description 使用前需要确保IO端口已经配置为PWM输出模式。
//
// @Param params *IOPWM true "IO PWM输出参数：
//   - address: IO端口地址
//   - frequency: PWM频率（单位：Hz）
//   - dutyCycle: 占空比（范围：0-100）
//     注意：
//   - 频率范围受硬件限制，具体参见产品手册
//   - 占空比表示高电平时间占总周期的百分比"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 将指令加入队列，按顺序执行
//   - false: 立即执行该指令
//   - 建议使用队列模式以确保输出时序"
//
// @Return uint64 "指令队列索引（当isQueued为true时有效）"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数为空
//   - IO地址无效
//   - 端口模式不正确
//   - 频率超出范围
//   - 占空比超出范围
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置IO端口1输出1kHz、占空比50%的PWM信号
//	params := &IOPWM{
//	    Address:    1,     // IO端口1
//	    Frequency:  1000,  // 频率1000Hz
//	    DutyCycle:  50,    // 占空比50%
//	}
//	index, err := dobot.SetIOPWM(params, true)
//	if err != nil {
//	    log.Printf("设置IO PWM输出失败: %v", err)
//	} else {
//	    log.Printf("IO PWM输出设置成功，指令索引: %d", index)
//	}
func (dobot *Dobot) SetIOPWM(params *IOPWM, isQueued bool) (queuedCmdIndex uint64, err error) {
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
// @Summary 设置电机参数
// @Description 通过EMotor结构体设置单个电机的控制参数
// @Param params body *EMotor true "电机参数结构体"
// @Param isQueued query bool true "是否队列执行"
// @Success 200 {number} uint8 "返回命令队列索引"
// @Failure 400 {object} error "设置电机参数失败时返回错误信息"
func (dobot *Dobot) SetEMotor(params *EMotor, isQueued bool) (queuedCmdIndex uint64, err error) {
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
// @Summary 设置电机S参数
// @Description 通过EMotorS结构体同时设置多个电机的参数
// @Param params body *EMotorS true "电机S参数结构体"
// @Param isQueued query bool true "是否队列执行"
// @Success 200 {number} uint8 "返回命令队列索引"
// @Failure 400 {object} error "设置电机S参数失败时返回错误信息"
func (dobot *Dobot) SetEMotorS(params *EMotorS, isQueued bool) (queuedCmdIndex uint64, err error) {
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
// @Summary 设置颜色传感器
// @Description 根据传入的enable、colorPort和version设置颜色传感器状态
// @Param enable query bool true "是否启用颜色传感器"
// @Param colorPort query ColorPort true "颜色传感器端口"
// @Param version query uint8 true "传感器版本号"
// @Success 200 {string} "设置成功返回空字符串"
// @Failure 400 {object} error "设置颜色传感器失败时返回错误信息"
func (dobot *Dobot) SetColorSensor(enable bool, colorPort ColorPort, version uint8) error {
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

	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// GetColorSensor 获取颜色传感器数据
// @Summary 获取颜色传感器数据
// @Description 获取当前颜色传感器的RGB数值
// @Success 200 {number} uint8 "返回R值"
// @Success 200 {number} uint8 "返回G值"
// @Success 200 {number} uint8 "返回B值"
// @Failure 400 {object} error "获取颜色传感器数据失败时返回错误信息"
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
// @Summary 设置红外传感器
// @Description 根据enable、infraredPort和version设置红外传感器状态
// @Param enable query bool true "是否启用红外传感器"
// @Param infraredPort query InfraredPort true "红外传感器端口"
// @Param version query uint8 true "传感器版本号"
// @Success 200 {string} "设置成功返回空字符串"
// @Failure 400 {object} error "设置红外传感器失败时返回错误信息"
func (dobot *Dobot) SetInfraredSensor(enable bool, infraredPort InfraredPort, version uint8) error {
	message := &Message{
		Id:       ProtocolInfraredSensor,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, enable)
	binary.Write(writer, binary.LittleEndian, infraredPort)
	binary.Write(writer, binary.LittleEndian, version)
	message.Params = writer.Bytes()

	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// GetInfraredSensor 获取红外传感器数据
// @Summary 获取红外传感器数据
// @Description 获取指定端口红外传感器的数据
// @Param port query InfraredPort true "红外传感器端口"
// @Success 200 {number} uint8 "返回红外传感器数据值"
// @Failure 400 {object} error "获取红外传感器数据失败时返回错误信息"
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
// @Summary 设置角度传感器静态误差
// @Description 为后臂和前臂角度传感器设置静态误差补偿值，单位为度
// @Param rearArmAngleError query float32 true "后臂角度静态误差"
// @Param frontArmAngleError query float32 true "前臂角度静态误差"
// @Success 200 {string} "设置成功返回空字符串"
// @Failure 400 {object} error "设置角度传感器静态误差失败时返回错误信息"
func (dobot *Dobot) SetAngleSensorStaticError(rearArmAngleError, frontArmAngleError float32) error {
	message := &Message{
		Id:       ProtocolAngleSensorStaticError,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, rearArmAngleError)
	binary.Write(writer, binary.LittleEndian, frontArmAngleError)
	message.Params = writer.Bytes()

	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// GetAngleSensorStaticError 获取角度传感器静态误差
// @Summary 获取角度传感器静态误差
// @Description 获取后臂和前臂角度传感器的静态误差补偿值
// @Success 200 {number} float32 "返回后臂角度静态误差"
// @Success 200 {number} float32 "返回前臂角度静态误差"
// @Failure 400 {object} error "获取角度传感器静态误差失败时返回错误信息"
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
// @Summary 设置角度传感器系数
// @Description 为校正角度传感器误差，通过设置后臂和前臂的比例系数
// @Param rearArmAngleCoef query float32 true "后臂角度系数"
// @Param frontArmAngleCoef query float32 true "前臂角度系数"
// @Success 200 {string} "设置成功返回空字符串"
// @Failure 400 {object} error "设置角度传感器系数失败时返回错误信息"
func (dobot *Dobot) SetAngleSensorCoef(rearArmAngleCoef, frontArmAngleCoef float32) error {
	message := &Message{
		Id:       ProtocolAngleSensorCoef,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, rearArmAngleCoef)
	binary.Write(writer, binary.LittleEndian, frontArmAngleCoef)
	message.Params = writer.Bytes()

	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// GetAngleSensorCoef 获取角度传感器系数
// @Summary 获取角度传感器系数
// @Description 获取后臂和前臂角度传感器的校正比例系数
// @Success 200 {number} float32 "返回后臂角度系数"
// @Success 200 {number} float32 "返回前臂角度系数"
// @Failure 400 {object} error "获取角度传感器系数失败时返回错误信息"
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

// SetBaseDecoderStaticError 设置底座解码器静态误差
// @Summary 设置底座解码器静态误差
// @Description 为底座解码器设置静态补偿误差，单位可能为度或mm
// @Param baseDecoderError query float32 true "底座解码器静态误差"
// @Success 200 {string} "设置成功返回空字符串"
// @Failure 400 {object} error "设置底座解码器静态误差失败时返回错误信息"
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

// GetBaseDecoderStaticError 获取底座解码器静态误差
// @Summary 获取底座解码器静态误差
// @Description 获取当前底座解码器的静态误差补偿值
// @Success 200 {number} float32 "返回静态误差值"
// @Failure 400 {object} error "获取底座解码器静态误差失败时返回错误信息"
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
// @Summary 设置左右手校准值
// @Description 设置机械臂左右手的校准参数，用于精度调整
// @Param lrHandCalibrateValue query float32 true "左右手校准值"
// @Success 200 {string} "设置成功返回空字符串"
// @Failure 400 {object} error "设置左右手校准值失败时返回错误信息"
func (dobot *Dobot) SetLRHandCalibrateValue(lrHandCalibrateValue float32) error {
	message := &Message{
		Id:       ProtocolLRHandCalibrateValue,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, lrHandCalibrateValue)
	message.Params = writer.Bytes()

	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// GetLRHandCalibrateValue 获取左右手校准值
// @Summary 获取左右手校准值
// @Description 获取当前机械臂左右手的校准参数值
// @Success 200 {number} float32 "返回校准值"
// @Failure 400 {object} error "获取左右手校准值失败时返回错误信息"
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

// SetQueuedCmdStartExec 执行Dobot命令
// @Summary 开始执行指令队列
// @Description 开始执行之前加入队列的所有指令。指令队列特点：
// @Description 1. 队列中的指令按照加入顺序依次执行
// @Description 2. 每条指令执行完成后才会执行下一条指令
// @Description 3. 可以随时暂停、继续或停止队列执行
// @Description 注意：开始执行前确保队列中的指令是正确的，且机械臂处于使能状态
// @Success 200 {string} "成功返回空字符串"
// @Failure 400 {object} error "开始执行队列命令失败时返回错误信息"
func (dobot *Dobot) SetQueuedCmdStartExec() error {
	message := &Message{
		Id:       ProtocolQueuedCmdStartExec,
		RW:       true,
		IsQueued: false,
	}
	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// SetQueuedCmdStopExec 停止执行队列命令
// @Summary 停止执行队列中的指令
// @Description 停止执行队列中的指令，可以通过 SetQueuedCmdStartExec 重新开始执行
// @Success 200 {string} "成功返回空字符串"
// @Failure 400 {object} error "停止执行队列命令失败时返回错误信息"
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
// @Summary 强制停止执行队列中的指令
// @Description 立即停止执行队列中的指令，与普通停止相比，强制停止会立即中断当前动作
// @Success 200 {string} "成功返回空字符串"
// @Failure 400 {object} error "强制停止队列命令失败时返回错误信息"
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
// @Summary 开始下载队列命令
// @Description 通知机械臂开始下载命令队列，指定循环次数和每循环命令行数
// @Param totalLoop query uint32 true "总循环次数"
// @Param linePerLoop query uint32 true "每循环命令行数"
// @Success 200 {string} "成功返回空字符串"
// @Failure 400 {object} error "开始下载队列命令失败时返回错误信息"
func (dobot *Dobot) SetQueuedCmdStartDownload(totalLoop uint32, linePerLoop uint32) error {
	message := &Message{
		Id:       ProtocolQueuedCmdStartDownload,
		RW:       true,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, totalLoop)
	binary.Write(writer, binary.LittleEndian, linePerLoop)
	message.Params = writer.Bytes()

	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// SetQueuedCmdStopDownload 停止下载队列命令
// @Summary 停止下载队列命令
// @Description 通知机械臂停止下载命令队列
// @Success 200 {string} "成功返回空字符串"
// @Failure 400 {object} error "停止下载队列命令失败时返回错误信息"
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
// @Summary 清除命令队列
// @Description 清空机械臂命令队列中的所有指令
// @Success 200 {string} "成功返回空字符串"
// @Failure 400 {object} error "清除队列命令失败时返回错误信息"
func (dobot *Dobot) SetQueuedCmdClear() error {
	message := &Message{
		Id:       ProtocolQueuedCmdClear,
		RW:       true,
		IsQueued: false,
	}
	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

func (dobot *Dobot) GetQueuedCmdLeftSpace() (uint32, error) {
	message := &Message{
		Id:       ProtocolQueuedCmdLeftSpace,
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

// GetQueuedCmdCurrentIndex 获取当前队列命令索引
// @Summary 获取当前命令索引
// @Description 获取机械臂当前正在执行的命令队列索引
// @Success 200 {number} uint64 "返回当前命令索引"
// @Failure 400 {object} error "获取当前队列命令索引失败时返回错误信息"
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
// @Summary 获取命令执行状态
// @Description 检查当前队列命令是否已完成执行
// @Success 200 {boolean} bool "返回运动是否完成"
// @Failure 400 {object} error "获取命令执行状态失败时返回错误信息"
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
// @Summary 设置PTP PO命令
// @Description 通过PTPCmd及并行输出命令数组同时发送PTP运动及附加操作
// @Param ptpCmd body *PTPCmd true "PTP命令结构体"
// @Param parallelCmd body []ParallelOutputCmd true "并行输出命令数组"
// @Success 200 {number} uint64 "返回命令队列索引"
// @Failure 400 {object} error "设置PTP PO命令失败时返回错误信息"
func (dobot *Dobot) SetPTPPOCmd(ptpCmd *PTPCmd, parallelCmd []ParallelOutputCmd) (queuedCmdIndex uint64, err error) {
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
// @Summary 设置PTP PO with L命令
// @Description 通过PTPWithLCmd结构体和并行输出命令数组发送带L参数的PTP命令
// @Param ptpWithLCmd body *PTPWithLCmd true "带L参数的PTP命令结构体"
// @Param parallelCmd body []ParallelOutputCmd true "并行输出命令数组"
// @Success 200 {number} uint64 "返回命令队列索引"
// @Failure 400 {object} error "设置带L并行输出PTP命令失败时返回错误信息"
func (dobot *Dobot) SetPTPPOWithLCmd(ptpWithLCmd *PTPWithLCmd, parallelCmd []ParallelOutputCmd) (queuedCmdIndex uint64, err error) {
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
// @Summary 设置WIFI配置模式
// @Description 设置机械臂的WIFI配置模式，enable为true时启用配置
// @Param enable query bool true "是否启用WIFI配置模式"
// @Success 200 {string} "设置成功返回空字符串"
// @Failure 400 {object} error "设置WIFI配置模式失败时返回错误信息"
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

// GetWIFIConfigMode 获取WIFI配置模式状态
// @Summary 获取WIFI配置模式状态
// @Description 获取当前机械臂的WIFI配置模式状态
// @Success 200 {boolean} bool "返回配置模式状态"
// @Failure 400 {object} error "获取WIFI配置模式状态失败时返回错误信息"
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
// @Summary 设置WIFI SSID
// @Description 设置机械臂连接WIFI使用的SSID
// @Param ssid query string true "WIFI网络SSID"
// @Success 200 {string} "设置成功返回空字符串"
// @Failure 400 {object} error "设置WIFI SSID失败时返回错误信息"
func (dobot *Dobot) SetWIFISSID(ssid string) error {
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

	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// GetWIFISSID 获取WIFI SSID
// @Summary 获取WIFI SSID
// @Description 获取当前配置的WIFI网络SSID
// @Success 200 {string} string "返回WIFI SSID"
// @Failure 400 {object} error "获取WIFI SSID失败时返回错误信息"
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
// @Summary 设置WIFI密码
// @Description 设置机械臂连接WIFI使用的密码
// @Param password query string true "WIFI网络密码"
// @Success 200 {string} "设置成功返回空字符串"
// @Failure 400 {object} error "设置WIFI密码失败时返回错误信息"
func (dobot *Dobot) SetWIFIPassword(password string) error {
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
	writer.WriteByte(0) // 添加一个字节 0x00 作为校验字节
	message.Params = writer.Bytes()

	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// GetWIFIPassword 获取WIFI密码
// @Summary 获取WIFI密码
// @Description 获取当前配置的WIFI网络密码
// @Success 200 {string} string "返回WIFI密码"
// @Failure 400 {object} error "获取WIFI密码失败时返回错误信息"
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
// @Summary 设置丢步参数
// @Description 设置机械臂运动过程中丢步判断的阈值
// @Param threshold query float32 true "丢步阈值"
// @Param isQueued query bool true "是否队列执行"
// @Success 200 {number} uint64 "返回命令队列索引"
// @Failure 400 {object} error "设置丢步参数失败时返回错误信息"
func (dobot *Dobot) SetLostStepParams(threshold float32, isQueued bool) (queuedCmdIndex uint64, err error) {
	message := &Message{
		Id:       ProtocolLostStepSet,
		RW:       true,
		IsQueued: isQueued,
	}

	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, threshold)
	message.Params = writer.Bytes()

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
// @Summary 设置丢步命令
// @Description 发送丢步命令以校正运动过程中的丢步问题
// @Param isQueued query bool true "是否队列执行"
// @Success 200 {number} uint64 "返回命令队列索引"
// @Failure 400 {object} error "设置丢步命令失败时返回错误信息"
func (dobot *Dobot) SetLostStepCmd(isQueued bool) (queuedCmdIndex uint64, err error) {
	message := &Message{
		Id:       ProtocolLostStepDetect,
		RW:       true,
		IsQueued: isQueued,
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
