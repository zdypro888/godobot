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
type Dobot struct {
	connector *Connector
}

// NewDobot 创建新的Dobot实例
// @Summary 创建新的机械臂控制实例
// @Description 初始化并返回一个新的Dobot机械臂控制实例。该实例包含了
// @Description 所有控制机械臂所需的方法。在使用任何控制功能之前，必须
// @Description 先创建一个控制实例。
//
// @Return *Dobot "返回值说明：
//   - 返回Dobot实例指针
//   - 该实例用于后续所有的机械臂控制操作
//     注意：创建实例后需要调用Connect方法建立连接"
//
// @Success 200 {object} *Dobot "返回Dobot实例"
//
// @Example
//
//	// 创建新的机械臂控制实例
//	dobot := NewDobot()
//	if dobot != nil {
//	    // 创建成功，继续建立连接
//	    err := dobot.Connect(ctx, "/dev/ttyUSB0", 115200)
//	    if err != nil {
//	        log.Fatal("连接失败:", err)
//	    }
//	    log.Printf("机械臂控制实例创建并连接成功")
//	}
func NewDobot() *Dobot {
	return &Dobot{connector: &Connector{}}
}

// ConnectDobot 连接到Dobot设备
// @Summary 建立与机械臂的通信连接
// @Description 通过指定的串口和波特率连接到Dobot机械臂。此函数会尝试
// @Description 打开指定的串口并建立与机械臂的通信连接。连接成功后才能
// @Description 执行后续的控制命令。如果连接失败，需要检查串口设备是否
// @Description 正确以及波特率设置是否匹配。
//
// @Param ctx context.Context true "上下文对象：
//   - 用于控制连接超时和取消操作
//   - 可以使用 context.Background() 创建"
//
// @Param portName string true "串口设备名称：
//   - Windows系统：'COM1', 'COM2', ...
//   - Linux系统：'/dev/ttyUSB0', '/dev/ttyACM0', ...
//   - macOS系统：'/dev/cu.usbserial-*'
//     注意：确保有正确的串口访问权限"
//
// @Param baudrate uint32 true "串口通信波特率：
//   - 标准值：115200
//   - 必须与机械臂设置匹配
//     注意：不正确的波特率会导致通信错误"
//
// @Return error "错误信息"
// @Success 200 {string} "连接成功"
// @Failure 400 {error} "连接失败，可能的错误：
//   - 串口不存在
//   - 串口被占用
//   - 波特率不支持
//   - 权限不足
//   - 通信超时
//   - 设备未就绪"
//
// @Example
//
//	// 连接到机械臂
//	ctx := context.Background()
//	err := dobot.Connect(ctx, "/dev/ttyUSB0", 115200)
//	if err != nil {
//	    log.Printf("连接失败: %v", err)
//	    // 检查串口设备和权限
//	    return
//	}
//	log.Printf("成功连接到机械臂")
//
//	// 使用带超时的连接
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	err = dobot.Connect(ctx, "/dev/ttyUSB0", 115200)
//	if err != nil {
//	    if err == context.DeadlineExceeded {
//	        log.Printf("连接超时")
//	    } else {
//	        log.Printf("连接失败: %v", err)
//	    }
//	    return
//	}
func (dobot *Dobot) Connect(ctx context.Context, portName string, baudrate uint32) error {
	err := dobot.connector.Open(ctx, portName, baudrate)
	if err != nil {
		return err
	}
	return nil
}

// SetDeviceSN 设置设备序列号
// @Summary 设置机械臂设备序列号
// @Description 为机械臂设置唯一的序列号标识。序列号用于区分不同的设备，
// @Description 在多台机械臂同时使用的场景下特别有用。设置后的序列号将被
// @Description 永久保存在设备中，重启后仍然有效。
//
// @Param sn string true "设备序列号：
//   - 长度：必须小于64字符
//   - 格式：字母、数字和特殊字符的组合
//   - 建议：使用有意义的标识，如'DOBOT_LAB_001'
//     注意：序列号一旦设置通常不需要经常更改"
//
// @Return error "错误信息"
// @Success 200 {string} "设置成功"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 序列号为空
//   - 序列号格式无效
//   - 序列号长度超出范围
//   - 设备被锁定
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置机械臂序列号
//	err := dobot.SetDeviceSN("DOBOT_LAB_001")
//	if err != nil {
//	    log.Printf("设置序列号失败: %v", err)
//	    return
//	}
//	log.Printf("序列号设置成功")
//
//	// 验证序列号设置
//	sn, err := dobot.GetDeviceSN()
//	if err != nil {
//	    log.Printf("获取序列号失败: %v", err)
//	    return
//	}
//	if sn != "DOBOT_LAB_001" {
//	    log.Printf("序列号设置验证失败")
//	    return
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
// @Summary 获取机械臂设备序列号
// @Description 获取当前连接的机械臂设备的序列号。序列号是设备的唯一标识，
// @Description 可用于区分不同的机械臂设备，也可用于验证设备身份或追踪
// @Description 设备历史记录。
//
// @Return string "设备序列号：
//   - 返回当前设备的序列号字符串
//   - 如果设备未设置序列号，可能返回空字符串
//     注意：序列号格式取决于之前的设置"
//
// @Return error "错误信息"
// @Success 200 {string} "返回设备序列号"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 设备未连接
//   - 通信错误
//   - 响应数据无效
//   - 设备信息读取失败"
//
// @Example
//
//	// 获取机械臂序列号
//	sn, err := dobot.GetDeviceSN()
//	if err != nil {
//	    log.Printf("获取序列号失败: %v", err)
//	    return
//	}
//	if sn == "" {
//	    log.Printf("设备未设置序列号")
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
// @Summary 获取机械臂设备的版本信息
// @Description 获取当前连接的机械臂设备的固件版本信息。版本信息包括
// @Description 主版本号、次版本号、修订版本号和硬件版本号。这些信息
// @Description 对于确保软件兼容性和故障诊断非常重要。
//
// @Return uint8 "主版本号：
//   - 表示重大功能更新
//   - 不同主版本可能不兼容"
//
// @Return uint8 "次版本号：
//   - 表示功能改进或增强
//   - 通常向后兼容"
//
// @Return uint8 "修订版本号：
//   - 表示错误修复和小改动
//   - 完全兼容"
//
// @Return uint8 "硬件版本号：
//   - 表示硬件设计版本
//   - 用于确定硬件兼容性"
//
// @Return error "错误信息"
// @Success 200 {object} struct{majorVersion,minorVersion,revision,hwVersion uint8} "返回版本信息"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 设备未连接
//   - 通信错误
//   - 响应数据无效
//   - 设备信息读取失败"
//
// @Example
//
//	// 获取机械臂版本信息
//	major, minor, rev, hw, err := dobot.GetDeviceVersion()
//	if err != nil {
//	    log.Printf("获取版本信息失败: %v", err)
//	    return
//	}
//	log.Printf("设备版本信息：")
//	log.Printf("  主版本号: %d", major)
//	log.Printf("  次版本号: %d", minor)
//	log.Printf("  修订版本: %d", rev)
//	log.Printf("  硬件版本: %d", hw)
//
//	// 检查版本兼容性
//	if major != 1 {
//	    log.Printf("警告：当前软件可能与该版本硬件不兼容")
//	    return
//	}
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
// @Summary 设置机械臂L轴（第五轴）参数
// @Description 设置机械臂是否启用L轴（第五轴）功能。L轴是机械臂的扩展轴，
// @Description 可以提供额外的旋转自由度。启用L轴后，机械臂的运动规划和
// @Description 控制将考虑L轴的参数。
//
// @Param isWithL bool true "是否启用L轴：
//   - true: 启用L轴功能
//   - false: 禁用L轴功能
//     注意：启用L轴需要硬件支持"
//
// @Param version uint8 true "L轴版本号：
//   - 用于适配不同版本的L轴硬件
//   - 具体值参见产品手册
//     注意：版本号必须与实际硬件匹配"
//
// @Return uint64 "指令队列索引"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 硬件不支持L轴
//   - 版本号无效
//   - 机械臂被锁定
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 启用L轴功能
//	index, err := dobot.SetDeviceWithL(true, 1)
//	if err != nil {
//	    log.Printf("设置L轴失败: %v", err)
//	    return
//	}
//	log.Printf("L轴功能已启用，指令索引: %d", index)
//
//	// 验证L轴状态
//	enabled, err := dobot.GetDeviceWithL()
//	if err != nil {
//	    log.Printf("获取L轴状态失败: %v", err)
//	    return
//	}
//	if !enabled {
//	    log.Printf("L轴启用失败")
//	    return
//	}
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
// @Summary 获取机械臂L轴（第五轴）的启用状态
// @Description 获取机械臂L轴（第五轴）功能的当前启用状态。通过此函数
// @Description 可以确认L轴功能是否正确启用，以便进行后续的L轴相关操作。
//
// @Return bool "L轴状态：
//   - true: L轴功能已启用
//   - false: L轴功能未启用
//     注意：状态与硬件支持和配置有关"
//
// @Return error "错误信息"
// @Success 200 {boolean} bool "返回L轴启用状态"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	// 获取L轴状态
//	enabled, err := dobot.GetDeviceWithL()
//	if err != nil {
//	    log.Printf("获取L轴状态失败: %v", err)
//	    return
//	}
//	if enabled {
//	    log.Printf("L轴功能已启用")
//	} else {
//	    log.Printf("L轴功能未启用")
//	}
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
// @Summary 获取机械臂设备的运行时间
// @Description 获取机械臂设备自上次启动以来的运行时间。此时间可用于
// @Description 监控设备的工作时长，评估设备使用情况，以及进行维护计划。
// @Description 时间从设备每次启动时开始计数。
//
// @Return uint32 "运行时间：
//   - 单位：秒
//   - 从设备启动开始计数
//   - 关机后重置
//     注意：此时间仅反映当次启动后的运行时长"
//
// @Return error "错误信息"
// @Success 200 {number} uint32 "返回运行时间（秒）"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	// 获取设备运行时间
//	seconds, err := dobot.GetDeviceTime()
//	if err != nil {
//	    log.Printf("获取运行时间失败: %v", err)
//	    return
//	}
//	hours := float64(seconds) / 3600.0
//	log.Printf("设备运行时间: %.2f小时", hours)
//
//	// 检查是否需要维护
//	if hours > 100 {
//	    log.Printf("提示：设备运行超过100小时，建议进行例行维护")
//	}
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
// @Summary 获取机械臂设备的详细信息
// @Description 获取机械臂设备的综合信息，包括设备类型、累计运行时间、
// @Description 开机次数、错误次数和警告次数等统计信息。这些信息对于
// @Description 设备维护、故障诊断和使用状况分析非常有用。
//
// @Return *DeviceCountInfo "设备信息结构体：
//   - deviceType: 设备类型编号
//   - runTime: 累计运行时间（单位：秒）
//   - powerOnCount: 开机次数
//   - errorCount: 错误次数
//   - warningCount: 警告次数
//     注意：这些统计信息在设备断电后仍然保留"
//
// @Return error "错误信息"
// @Success 200 {object} *DeviceCountInfo "返回设备信息结构体"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效
//   - 数据解析错误"
//
// @Example
//
//	// 获取设备详细信息
//	info, err := dobot.GetDeviceInfo()
//	if err != nil {
//	    log.Printf("获取设备信息失败: %v", err)
//	    return
//	}
//	log.Printf("设备信息：")
//	log.Printf("  设备类型: %d", info.DeviceType)
//	log.Printf("  运行时间: %.2f小时", float64(info.RunTime)/3600.0)
//	log.Printf("  开机次数: %d", info.PowerOnCount)
//	log.Printf("  错误次数: %d", info.ErrorCount)
//	log.Printf("  警告次数: %d", info.WarningCount)
//
//	// 检查设备状态
//	if info.ErrorCount > 10 {
//	    log.Printf("警告：设备错误次数过多，建议进行检查")
//	}
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
// @Description 获取机械臂末端在笛卡尔坐标系下的位置和姿态信息，以及各关节的角度值。
// @Description 这些信息对于实时监控机械臂状态、轨迹规划和位置校准非常重要。
// @Description 建议在执行运动指令前后获取位姿进行验证。
//
// @Return *Pose "位姿结构体：
//   - x: X轴坐标（单位：mm）
//   - y: Y轴坐标（单位：mm）
//   - z: Z轴坐标（单位：mm）
//   - r: 末端旋转角度（单位：度）
//   - jointAngle: 各关节角度数组[4]（单位：度）
//     注意：坐标值基于机械臂基座坐标系"
//
// @Return error "错误信息"
// @Success 200 {object} *Pose "返回位姿信息结构体"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 传感器异常
//   - 数据解析错误"
//
// @Example
//
//	// 获取当前位姿信息
//	pose, err := dobot.GetPose()
//	if err != nil {
//	    log.Printf("获取位姿失败: %v", err)
//	    return
//	}
//	log.Printf("当前位姿信息：")
//	log.Printf("  位置: (%.2f, %.2f, %.2f) mm", pose.X, pose.Y, pose.Z)
//	log.Printf("  旋转角度: %.2f°", pose.R)
//	log.Printf("  关节角度:")
//	for i, angle := range pose.JointAngle {
//	    log.Printf("    关节%d: %.2f°", i+1, angle)
//	}
//
//	// 检查位置是否在工作空间内
//	if math.Sqrt(pose.X*pose.X + pose.Y*pose.Y) > 300 {
//	    log.Printf("警告：当前位置接近工作空间边界")
//	}
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
// @Summary 重置机械臂位姿到指定状态
// @Description 将机械臂重置到指定的位姿状态。可以选择手动模式或自动模式，
// @Description 并可以指定后臂和前臂的角度。此功能通常用于初始化机械臂位置
// @Description 或在异常情况后恢复到已知状态。
//
// @Param manual bool true "重置模式：
//   - true: 手动重置模式，使用指定的角度值
//   - false: 自动重置模式，使用默认角度值
//     注意：手动模式需要确保角度值安全可达"
//
// @Param rearArmAngle float32 true "后臂角度：
//   - 单位：度（°）
//   - 范围：通常在-90°到90°之间
//   - 在手动模式下生效
//     注意：角度值必须在机械臂可达范围内"
//
// @Param frontArmAngle float32 true "前臂角度：
//   - 单位：度（°）
//   - 范围：通常在-90°到90°之间
//   - 在手动模式下生效
//     注意：角度值必须在机械臂可达范围内"
//
// @Return error "错误信息"
// @Success 200 {string} "重置成功"
// @Failure 400 {error} "重置失败，可能的错误：
//   - 角度值超出范围
//   - 位置不可达
//   - 机械臂被锁定
//   - 机械臂处于报警状态
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 使用自动模式重置位姿
//	err := dobot.ResetPose(false, 0, 0)
//	if err != nil {
//	    log.Printf("自动重置位姿失败: %v", err)
//	    return
//	}
//	log.Printf("位姿已自动重置")
//
//	// 使用手动模式重置到特定角度
//	err = dobot.ResetPose(true, 45.0, 30.0)
//	if err != nil {
//	    log.Printf("手动重置位姿失败: %v", err)
//	    return
//	}
//	log.Printf("位姿已重置到指定角度")
//
//	// 验证重置后的位姿
//	pose, _ := dobot.GetPose()
//	if pose != nil {
//	    log.Printf("重置后的关节角度：")
//	    log.Printf("  后臂: %.2f°", pose.JointAngle[1])
//	    log.Printf("  前臂: %.2f°", pose.JointAngle[2])
//	}
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
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
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
	params := &HOMEParams{}
	if err := binary.Read(resp.Reader(), binary.LittleEndian, params); err != nil {
		return nil, fmt.Errorf("failed to read HOME params: %v", err)
	}
	return params, nil
}

// SetHOMECmd 执行回零操作
// @Summary 执行机械臂回零操作
// @Description 执行机械臂的回零操作。回零操作会将机械臂移动到预设的
// @Description 原点位置，这个位置通常是机械臂的安全位置。此功能用于
// @Description 初始化机械臂位置或在异常后恢复到已知状态。
//
// @Param cmd *HOMECmd true "回零命令参数：
//   - 包含回零操作的具体参数
//   - 具体参数含义参见HOMECmd结构体定义
//     注意：确保参数设置合理，避免碰撞"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 指令加入队列，按顺序执行
//   - false: 立即执行指令
//     注意：建议使用队列模式以确保操作顺序"
//
// @Return uint64 "指令索引：
//   - 返回指令队列索引
//   - 仅在isQueued为true时有效"
//
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令索引"
// @Failure 400 {error} "执行失败，可能的错误：
//   - 参数无效
//   - 机械臂被锁定
//   - 机械臂处于报警状态
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 执行回零操作
//	cmd := &HOMECmd{
//	    // 设置回零参数
//	}
//	index, err := dobot.SetHOMECmd(cmd, true)
//	if err != nil {
//	    log.Printf("执行回零操作失败: %v", err)
//	    return
//	}
//	log.Printf("回零操作开始执行，指令索引: %d", index)
//
//	// 等待回零操作完成
//	for {
//	    finished, _ := dobot.GetQueuedCmdMotionFinish()
//	    if finished {
//	        log.Printf("回零操作完成")
//	        break
//	    }
//	    time.Sleep(100 * time.Millisecond)
//	}
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
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
	return queuedCmdIndex, nil
}

// SetAutoLevelingCmd 执行自动调平
// @Summary 执行机械臂自动调平操作
// @Description 执行机械臂的自动调平操作。自动调平功能用于确保机械臂
// @Description 的基座处于水平状态，这对于保证运动精度和重复定位精度
// @Description 非常重要。调平过程会自动测量和补偿基座倾斜。
//
// @Param cmd *AutoLevelingCmd true "自动调平命令参数：
//   - 包含调平操作的具体参数
//   - 具体参数含义参见AutoLevelingCmd结构体定义
//     注意：确保调平区域内无障碍物"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 指令加入队列，按顺序执行
//   - false: 立即执行指令
//     注意：建议使用队列模式以确保操作顺序"
//
// @Return uint64 "指令索引：
//   - 返回指令队列索引
//   - 仅在isQueued为true时有效"
//
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令索引"
// @Failure 400 {error} "执行失败，可能的错误：
//   - 参数无效
//   - 机械臂被锁定
//   - 机械臂处于报警状态
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 执行自动调平
//	cmd := &AutoLevelingCmd{
//	    // 设置调平参数
//	}
//	index, err := dobot.SetAutoLevelingCmd(cmd, true)
//	if err != nil {
//	    log.Printf("执行自动调平失败: %v", err)
//	    return
//	}
//	log.Printf("自动调平开始执行，指令索引: %d", index)
//
//	// 等待调平完成
//	for {
//	    finished, _ := dobot.GetQueuedCmdMotionFinish()
//	    if finished {
//	        // 获取调平结果
//	        result, err := dobot.GetAutoLevelingResult()
//	        if err != nil {
//	            log.Printf("获取调平结果失败: %v", err)
//	        } else {
//	            log.Printf("调平完成，补偿角度: %.2f°", result)
//	        }
//	        break
//	    }
//	    time.Sleep(100 * time.Millisecond)
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
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
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
	return math.Float32frombits(binary.LittleEndian.Uint32(resp.Data())), nil
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
	return HHTTrigMode(resp.Data()[0]), nil
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
	return resp.Data()[0] != 0, nil
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
	return resp.Data()[0] != 0, nil
}

// SetEndEffectorParams 设置末端执行器参数
// @Summary 设置机械臂末端执行器参数
// @Description 设置机械臂末端执行器的坐标偏移量参数。这些参数用于定义
// @Description 末端执行器相对于机械臂末端法兰的位置偏移，确保机械臂能够
// @Description 准确定位到目标位置。
//
// @Param params *EndEffectorParams true "末端执行器参数：
//   - xBias: X轴偏移量（单位：mm）
//   - yBias: Y轴偏移量（单位：mm）
//   - zBias: Z轴偏移量（单位：mm）
//     注意：偏移量基于末端法兰坐标系"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 指令加入队列，按顺序执行
//   - false: 立即执行指令
//     注意：建议使用队列模式以确保参数设置顺序"
//
// @Return uint64 "指令索引：
//   - 返回指令队列索引
//   - 仅在isQueued为true时有效"
//
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数无效
//   - 偏移量超出范围
//   - 机械臂被锁定
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置末端执行器偏移量
//	params := &EndEffectorParams{
//	    XBias: 0,    // X轴无偏移
//	    YBias: 20,   // Y轴偏移20mm
//	    ZBias: 30,   // Z轴偏移30mm
//	}
//	index, err := dobot.SetEndEffectorParams(params, true)
//	if err != nil {
//	    log.Printf("设置末端执行器参数失败: %v", err)
//	    return
//	}
//	log.Printf("末端执行器参数设置成功，指令索引: %d", index)
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
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
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
	params := &EndEffectorParams{}
	if err := binary.Read(resp.Reader(), binary.LittleEndian, params); err != nil {
		return nil, fmt.Errorf("failed to read END_EFFECTOR_PARAMS: %v", err)
	}
	return params, nil
}

// SetEndEffectorLaser 设置末端激光状态
// @Summary 设置机械臂末端激光器的工作状态
// @Description 控制机械臂末端激光器的开关状态和使能状态。激光器可用于
// @Description 激光雕刻、切割、标记等应用。使用前需要确保激光器已正确
// @Description 安装并配置。
//
// @Param enableCtrl bool true "使能控制：
//   - true: 启用激光器控制
//   - false: 禁用激光器控制
//     注意：必须先使能才能控制激光器"
//
// @Param on bool true "激光器开关：
//   - true: 打开激光器
//   - false: 关闭激光器
//     注意：激光器打开时要注意安全"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 指令加入队列，按顺序执行
//   - false: 立即执行指令
//     注意：建议使用队列模式以确保控制顺序"
//
// @Return uint64 "指令索引：
//   - 返回指令队列索引
//   - 仅在isQueued为true时有效"
//
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 激光器未安装
//   - 机械臂被锁定
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 启用并打开激光器
//	index, err := dobot.SetEndEffectorLaser(true, true, true)
//	if err != nil {
//	    log.Printf("设置激光器状态失败: %v", err)
//	    return
//	}
//	log.Printf("激光器已启用并打开，指令索引: %d", index)
//
//	// 等待一段时间后关闭激光器
//	time.Sleep(5 * time.Second)
//	index, err = dobot.SetEndEffectorLaser(true, false, true)
//	if err != nil {
//	    log.Printf("关闭激光器失败: %v", err)
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

	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
	return queuedCmdIndex, nil
}

// GetEndEffectorLaser 获取末端激光状态
// @Summary 获取机械臂末端激光器的工作状态
// @Description 获取末端激光器的当前使能状态和开关状态。通过此功能可以
// @Description 确认激光器是否已正确启用和工作。在执行激光相关操作前，
// @Description 建议先检查激光器状态。
//
// @Return bool "使能状态：
//   - true: 激光器控制已启用
//   - false: 激光器控制已禁用"
//
// @Return bool "开关状态：
//   - true: 激光器已打开
//   - false: 激光器已关闭"
//
// @Return error "错误信息"
// @Success 200 {object} struct{isCtrlEnabled,isOn bool} "返回激光器状态"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 激光器未安装
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	// 获取激光器状态
//	enabled, on, err := dobot.GetEndEffectorLaser()
//	if err != nil {
//	    log.Printf("获取激光器状态失败: %v", err)
//	    return
//	}
//	log.Printf("激光器状态：")
//	if enabled {
//	    log.Printf("  控制已启用")
//	    if on {
//	        log.Printf("  当前已打开")
//	    } else {
//	        log.Printf("  当前已关闭")
//	    }
//	} else {
//	    log.Printf("  控制未启用")
//	}
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
	data := resp.Data()
	return data[0] != 0, data[1] != 0, nil
}

// SetEndEffectorSuctionCup 设置末端吸盘状态
// @Summary 设置机械臂末端吸盘的工作状态
// @Description 控制机械臂末端吸盘的使能和吸附状态。吸盘是一种常用的
// @Description 末端执行器，用于吸附和搬运物体。使用前需要确保吸盘已
// @Description 正确安装并连接气源。
//
// @Param enableCtrl bool true "使能控制：
//   - true: 启用吸盘控制
//   - false: 禁用吸盘控制
//     注意：必须先使能才能控制吸盘"
//
// @Param suck bool true "吸盘状态：
//   - true: 启动吸附
//   - false: 停止吸附
//     注意：确保物体重量在吸盘负载范围内"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 指令加入队列，按顺序执行
//   - false: 立即执行指令
//     注意：建议使用队列模式以确保控制顺序"
//
// @Return uint64 "指令索引：
//   - 返回指令队列索引
//   - 仅在isQueued为true时有效"
//
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 吸盘未安装
//   - 气源压力不足
//   - 机械臂被锁定
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 启用并打开吸盘
//	index, err := dobot.SetEndEffectorSuctionCup(true, true, true)
//	if err != nil {
//	    log.Printf("设置吸盘状态失败: %v", err)
//	    return
//	}
//	log.Printf("吸盘已启用并开始吸附，指令索引: %d", index)
//
//	// 等待一段时间后关闭吸盘
//	time.Sleep(2 * time.Second)
//	index, err = dobot.SetEndEffectorSuctionCup(true, false, true)
//	if err != nil {
//	    log.Printf("关闭吸盘失败: %v", err)
//	}
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

	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
	return queuedCmdIndex, nil
}

// GetEndEffectorSuctionCup 获取末端执行器吸盘状态
// @Summary 获取末端执行器吸盘的工作状态
// @Description 获取末端吸盘的当前使能状态和吸附状态。通过此功能可以
// @Description 确认吸盘是否已正确启用和工作。在执行吸附相关操作前，
// @Description 建议先检查吸盘状态。
//
// @Return bool "使能状态：
//   - true: 吸盘控制已启用
//   - false: 吸盘控制已禁用"
//
// @Return bool "吸附状态：
//   - true: 正在吸附
//   - false: 未吸附"
//
// @Return error "错误信息"
// @Success 200 {object} struct{isCtrlEnabled,isSucked bool} "返回吸盘状态"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 吸盘未安装
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	// 获取吸盘状态
//	enabled, sucking, err := dobot.GetEndEffectorSuctionCup()
//	if err != nil {
//	    log.Printf("获取吸盘状态失败: %v", err)
//	    return
//	}
//	log.Printf("吸盘状态：")
//	if enabled {
//	    log.Printf("  控制已启用")
//	    if sucking {
//	        log.Printf("  当前正在吸附")
//	    } else {
//	        log.Printf("  当前未吸附")
//	    }
//	} else {
//	    log.Printf("  控制未启用")
//	}
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
	data := resp.Data()
	return data[0] != 0, data[1] != 0, nil
}

// SetEndEffectorGripper 设置末端夹爪状态
// @Summary 设置机械臂末端夹爪的工作状态
// @Description 控制机械臂末端夹爪的使能和夹持状态。夹爪是一种常用的
// @Description 末端执行器，用于抓取和搬运物体。使用前需要确保夹爪已
// @Description 正确安装并配置。
//
// @Param enableCtrl bool true "使能控制：
//   - true: 启用夹爪控制
//   - false: 禁用夹爪控制
//     注意：必须先使能才能控制夹爪"
//
// @Param grip bool true "夹爪状态：
//   - true: 闭合夹爪
//   - false: 打开夹爪
//     注意：确保物体尺寸在夹爪范围内"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 指令加入队列，按顺序执行
//   - false: 立即执行指令
//     注意：建议使用队列模式以确保控制顺序"
//
// @Return uint64 "指令索引：
//   - 返回指令队列索引
//   - 仅在isQueued为true时有效"
//
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 夹爪未安装
//   - 机械臂被锁定
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 启用并闭合夹爪
//	index, err := dobot.SetEndEffectorGripper(true, true, true)
//	if err != nil {
//	    log.Printf("设置夹爪状态失败: %v", err)
//	    return
//	}
//	log.Printf("夹爪已启用并闭合，指令索引: %d", index)
//
//	// 等待一段时间后打开夹爪
//	time.Sleep(2 * time.Second)
//	index, err = dobot.SetEndEffectorGripper(true, false, true)
//	if err != nil {
//	    log.Printf("打开夹爪失败: %v", err)
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

	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
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
	data := resp.Data()
	return data[0] != 0, data[1] != 0, nil
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
//	    return
//	}
//	log.Printf("正在切换到左手方向，指令索引: %d", index)
//
//	// 等待切换完成
//	time.Sleep(2 * time.Second)
//	orientation, _ := dobot.GetArmOrientation()
//	if orientation == ArmOrientation_LeftyArmOrientation {
//	    log.Printf("机械臂方向切换成功")
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

	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
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
//	    return
//	}
//	if orientation == ArmOrientation_LeftyArmOrientation {
//	    log.Printf("当前为左手方向")
//	} else {
//	    log.Printf("当前为右手方向")
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
	return ArmOrientation(resp.Data()[0]), nil
}

// SetJOGJointParams 设置关节点动参数
// @Summary 设置机械臂关节点动运动的速度和加速度参数
// @Description 设置机械臂在JOG模式下各关节运动的速度和加速度参数。
// @Description 这些参数影响机械臂在手动点动模式下的运动特性。合适的
// @Description 参数设置可以确保运动平稳且可控。
//
// @Param params *JOGJointParams true "关节点动参数：
//   - velocity: 各关节速度数组（单位：°/s）
//   - acceleration: 各关节加速度数组（单位：°/s²）
//     注意：参数必须在机械臂规格范围内"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 指令加入队列，按顺序执行
//   - false: 立即执行指令
//     注意：建议使用队列模式以确保参数设置顺序"
//
// @Return uint64 "指令索引：
//   - 返回指令队列索引
//   - 仅在isQueued为true时有效"
//
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数无效
//   - 参数超出范围
//   - 机械臂被锁定
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置关节点动参数
//	params := &JOGJointParams{
//	    Velocity:     [4]float32{10, 10, 10, 10},  // 各关节速度10°/s
//	    Acceleration: [4]float32{50, 50, 50, 50},  // 各关节加速度50°/s²
//	}
//	index, err := dobot.SetJOGJointParams(params, true)
//	if err != nil {
//	    log.Printf("设置关节点动参数失败: %v", err)
//	    return
//	}
//	log.Printf("关节点动参数设置成功，指令索引: %d", index)
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
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
	return queuedCmdIndex, nil
}

// GetJOGJointParams 获取关节点动参数
// @Summary 获取机械臂关节点动运动的速度和加速度参数
// @Description 获取机械臂在JOG模式下各关节运动的当前速度和加速度参数。
// @Description 通过这些参数可以了解机械臂在手动点动模式下的运动特性
// @Description 设置。
//
// @Return *JOGJointParams "关节点动参数：
//   - velocity: 各关节速度数组（单位：°/s）
//   - acceleration: 各关节加速度数组（单位：°/s²）
//     注意：返回当前实际设置的参数值"
//
// @Return error "错误信息"
// @Success 200 {object} *JOGJointParams "返回关节点动参数结构体"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	// 获取关节点动参数
//	params, err := dobot.GetJOGJointParams()
//	if err != nil {
//	    log.Printf("获取关节点动参数失败: %v", err)
//	    return
//	}
//	log.Printf("当前关节点动参数：")
//	for i := 0; i < 4; i++ {
//	    log.Printf("  关节%d - 速度: %.1f°/s, 加速度: %.1f°/s²",
//	        i+1, params.Velocity[i], params.Acceleration[i])
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
	params := &JOGJointParams{}
	if err := binary.Read(resp.Reader(), binary.LittleEndian, params); err != nil {
		return nil, fmt.Errorf("failed to read JOG joint params: %v", err)
	}
	return params, nil
}

// SetJOGCoordinateParams 设置坐标点动参数
// @Summary 设置机械臂坐标点动运动的速度和加速度参数
// @Description 设置机械臂在JOG模式下笛卡尔坐标系运动的速度和加速度参数。
// @Description 这些参数影响机械臂在手动点动模式下的运动特性。合适的
// @Description 参数设置可以确保运动平稳且可控。
//
// @Param params *JOGCoordinateParams true "坐标点动参数：
//   - velocity: 各轴速度数组（X/Y/Z单位：mm/s，R单位：°/s）
//   - acceleration: 各轴加速度数组（X/Y/Z单位：mm/s²，R单位：°/s²）
//     注意：参数必须在机械臂规格范围内"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 指令加入队列，按顺序执行
//   - false: 立即执行指令
//     注意：建议使用队列模式以确保参数设置顺序"
//
// @Return uint64 "指令索引：
//   - 返回指令队列索引
//   - 仅在isQueued为true时有效"
//
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数无效
//   - 参数超出范围
//   - 机械臂被锁定
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置坐标点动参数
//	params := &JOGCoordinateParams{
//	    Velocity:     [4]float32{100, 100, 100, 50},  // XYZ速度100mm/s，R速度50°/s
//	    Acceleration: [4]float32{500, 500, 500, 200}, // XYZ加速度500mm/s²，R加速度200°/s²
//	}
//	index, err := dobot.SetJOGCoordinateParams(params, true)
//	if err != nil {
//	    log.Printf("设置坐标点动参数失败: %v", err)
//	    return
//	}
//	log.Printf("坐标点动参数设置成功，指令索引: %d", index)
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
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
	return queuedCmdIndex, nil
}

// GetJOGCoordinateParams 获取坐标点动参数
// @Summary 获取机械臂坐标点动运动的速度和加速度参数
// @Description 获取机械臂在JOG模式下笛卡尔坐标系运动的当前速度和加速度参数。
// @Description 通过这些参数可以了解机械臂在手动点动模式下的运动特性
// @Description 设置。
//
// @Return *JOGCoordinateParams "坐标点动参数：
//   - velocity: 各轴速度数组（X/Y/Z单位：mm/s，R单位：°/s）
//   - acceleration: 各轴加速度数组（X/Y/Z单位：mm/s²，R单位：°/s²）
//     注意：返回当前实际设置的参数值"
//
// @Return error "错误信息"
// @Success 200 {object} *JOGCoordinateParams "返回坐标点动参数结构体"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	// 获取坐标点动参数
//	params, err := dobot.GetJOGCoordinateParams()
//	if err != nil {
//	    log.Printf("获取坐标点动参数失败: %v", err)
//	    return
//	}
//	log.Printf("当前坐标点动参数：")
//	axes := []string{"X", "Y", "Z", "R"}
//	for i, axis := range axes {
//	    if i < 3 {
//	        log.Printf("  %s轴 - 速度: %.1fmm/s, 加速度: %.1fmm/s²",
//	            axis, params.Velocity[i], params.Acceleration[i])
//	    } else {
//	        log.Printf("  %s轴 - 速度: %.1f°/s, 加速度: %.1f°/s²",
//	            axis, params.Velocity[i], params.Acceleration[i])
//	    }
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

	params := &JOGCoordinateParams{}
	if err := binary.Read(resp.Reader(), binary.LittleEndian, params); err != nil {
		return nil, fmt.Errorf("failed to read JOG coordinate params: %v", err)
	}
	return params, nil
}

// SetJOGLParams 设置JOGL参数
// @Summary 设置JOG模式下的滑轨速度和加速度参数
// @Description 设置机械臂在JOG模式下滑轨运动的速度和加速度参数。
// @Description 这些参数影响机械臂在手动点动模式下的滑轨运动特性。合适的
// @Description 参数设置可以确保运动平稳且可控。
//
// @Param params *JOGLParams true "滑轨点动参数：
//   - velocity: 滑轨速度（单位：mm/s）
//   - acceleration: 滑轨加速度（单位：mm/s²）
//     注意：参数必须在机械臂规格范围内"
//
// @Return uint64 "指令索引：
//   - 返回指令队列索引
//   - 仅在isQueued为true时有效"
//
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数无效
//   - 参数超出范围
//   - 机械臂被锁定
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置滑轨点动参数
//	params := &JOGLParams{
//	    Velocity:     100,  // 滑轨速度100mm/s
//	    Acceleration: 500,  // 滑轨加速度500mm/s²
//	}
//	index, err := dobot.SetJOGLParams(params)
//	if err != nil {
//	    log.Printf("设置滑轨点动参数失败: %v", err)
//	    return
//	}
//	log.Printf("滑轨点动参数设置成功，指令索引: %d", index)
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
// @Summary 获取JOG模式下的滑轨速度和加速度参数
// @Description 获取机械臂在JOG模式下滑轨运动的当前速度和加速度参数。
// @Description 通过这些参数可以了解机械臂在手动点动模式下的滑轨运动
// @Description 特性设置。
//
// @Return *JOGLParams "滑轨点动参数：
//   - velocity: 滑轨速度（单位：mm/s）
//   - acceleration: 滑轨加速度（单位：mm/s²）
//     注意：返回当前实际设置的参数值"
//
// @Return error "错误信息"
// @Success 200 {object} *JOGLParams "返回滑轨点动参数结构体"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	// 获取滑轨点动参数
//	params, err := dobot.GetJOGLParams()
//	if err != nil {
//	    log.Printf("获取滑轨点动参数失败: %v", err)
//	    return
//	}
//	log.Printf("当前滑轨点动参数：")
//	log.Printf("  速度: %.1f mm/s", params.Velocity)
//	log.Printf("  加速度: %.1f mm/s²", params.Acceleration)
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

	params := &JOGLParams{}
	if err := binary.Read(resp.Reader(), binary.LittleEndian, params); err != nil {
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
//	    return
//	}
//	log.Printf("JOG通用参数设置成功，指令索引: %d", index)
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

	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
	return queuedCmdIndex, nil
}

// GetJOGCommonParams 获取JOG通用参数
// @Summary 获取JOG模式下的通用运动参数
// @Description 获取机械臂在JOG（点动）模式下所有运动方式共用的基础参数。
// @Description 通过这些参数可以了解当前JOG运动的速度和加速度设置，便于
// @Description 调整运动特性。
//
// @Return *JOGCommonParams "JOG通用参数：
//   - velocityRatio: 速度比例，范围[0-100]
//   - accelerationRatio: 加速度比例，范围[0-100]
//     注意：返回当前实际设置的参数值"
//
// @Return error "错误信息"
// @Success 200 {object} *JOGCommonParams "返回JOG通用参数结构体"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	// 获取JOG通用参数
//	params, err := dobot.GetJOGCommonParams()
//	if err != nil {
//	    log.Printf("获取JOG通用参数失败: %v", err)
//	    return
//	}
//	log.Printf("当前JOG通用参数：")
//	log.Printf("  速度比例: %d%%", params.VelocityRatio)
//	log.Printf("  加速度比例: %d%%", params.AccelerationRatio)
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

	params := &JOGCommonParams{}
	if err := binary.Read(resp.Reader(), binary.LittleEndian, params); err != nil {
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
// @Param cmd *JOGCmd true "JOG运动指令参数：
//   - isJoint: 运动模式选择（true为关节模式/false为坐标模式）
//   - index: 运动轴索引（关节模式：0-3为关节1-4；坐标模式：0-3为X/Y/Z/R轴）
//   - direction: 运动方向（1正向/0停止/-1负向）
//     注意：确保选择正确的运动模式和轴索引"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 指令加入队列，按顺序执行
//   - false: 立即执行指令
//     注意：建议使用队列模式以确保运动顺序"
//
// @Return uint64 "指令索引：
//   - 返回指令队列索引
//   - 仅在isQueued为true时有效"
//
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令索引"
// @Failure 400 {error} "执行失败，可能的错误：
//   - 参数无效
//   - 轴索引超出范围
//   - 机械臂被锁定
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
//	    return
//	}
//	log.Printf("JOG运动开始执行，指令索引: %d", index)
//
//	// 等待一段时间后停止运动
//	time.Sleep(1 * time.Second)
//	cmd.Direction = 0    // 停止运动
//	index, err = dobot.SetJOGCmd(cmd, true)
//	if err != nil {
//	    log.Printf("停止JOG运动失败: %v", err)
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

	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
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
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
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

	params := &PTPJointParams{}
	if err := binary.Read(resp.Reader(), binary.LittleEndian, params); err != nil {
		return nil, fmt.Errorf("failed to read PTP joint params: %v", err)
	}
	return params, nil
}

// SetPTPCoordinateParams 设置PTP坐标运动参数
// @Summary 设置机械臂PTP坐标运动的速度和加速度参数
// @Description 设置机械臂在PTP（点到点）坐标运动模式下的速度和加速度参数。
// @Description 这些参数直接影响机械臂在笛卡尔坐标系下运动的性能和轨迹特性。
// @Description 合理的参数设置可以在保证运动精度的同时优化运动时间。
//
// @Param params *PTPCoordinateParams true "坐标运动参数：
//   - xyzVelocity: XYZ轴速度（单位：mm/s）
//   - rVelocity: R轴旋转速度（单位：°/s）
//   - xyzAcceleration: XYZ轴加速度（单位：mm/s²）
//   - rAcceleration: R轴旋转加速度（单位：°/s²）
//     注意：参数必须在机械臂规格范围内"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 指令加入队列，按顺序执行
//   - false: 立即执行指令
//     注意：建议使用队列模式以确保参数设置顺序"
//
// @Return uint64 "指令索引：
//   - 返回指令队列索引
//   - 仅在isQueued为true时有效"
//
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数无效
//   - 参数超出范围
//   - 机械臂被锁定
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置PTP坐标运动参数
//	params := &PTPCoordinateParams{
//	    XYZVelocity:     200,  // XYZ轴速度200mm/s
//	    RVelocity:       100,  // R轴速度100°/s
//	    XYZAcceleration: 800,  // XYZ轴加速度800mm/s²
//	    RAcceleration:   400,  // R轴加速度400°/s²
//	}
//	index, err := dobot.SetPTPCoordinateParams(params, true)
//	if err != nil {
//	    log.Printf("设置PTP坐标运动参数失败: %v", err)
//	    return
//	}
//	log.Printf("PTP坐标运动参数设置成功，指令索引: %d", index)
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

	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
	return queuedCmdIndex, nil
}

// GetPTPCoordinateParams 获取PTP坐标运动参数
// @Summary 获取机械臂PTP坐标运动的速度和加速度参数
// @Description 获取机械臂在PTP（点到点）坐标运动模式下当前的速度和加速度
// @Description 参数。通过这些参数可以了解机械臂在笛卡尔坐标系下的运动性能
// @Description 配置，便于调整和优化运动控制。
//
// @Return *PTPCoordinateParams "坐标运动参数：
//   - xyzVelocity: XYZ轴速度（单位：mm/s）
//   - rVelocity: R轴旋转速度（单位：°/s）
//   - xyzAcceleration: XYZ轴加速度（单位：mm/s²）
//   - rAcceleration: R轴旋转加速度（单位：°/s²）
//     注意：返回当前实际设置的参数值"
//
// @Return error "错误信息"
// @Success 200 {object} *PTPCoordinateParams "返回坐标运动参数结构体"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	// 获取PTP坐标运动参数
//	params, err := dobot.GetPTPCoordinateParams()
//	if err != nil {
//	    log.Printf("获取PTP坐标运动参数失败: %v", err)
//	    return
//	}
//	log.Printf("当前PTP坐标运动参数：")
//	log.Printf("  XYZ轴速度: %.1f mm/s", params.XYZVelocity)
//	log.Printf("  R轴速度: %.1f °/s", params.RVelocity)
//	log.Printf("  XYZ轴加速度: %.1f mm/s²", params.XYZAcceleration)
//	log.Printf("  R轴加速度: %.1f °/s²", params.RAcceleration)
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

	params := &PTPCoordinateParams{}
	if err := binary.Read(resp.Reader(), binary.LittleEndian, params); err != nil {
		return nil, fmt.Errorf("failed to read PTP coordinate params: %v", err)
	}
	return params, nil
}

// SetPTPLParams 设置PTPL运动参数
// @Summary 设置机械臂PTPL运动的速度和加速度参数
// @Description 设置机械臂在PTPL（点到点线性）运动模式下的速度和加速度参数。
// @Description 这些参数直接影响机械臂在直线运动时的性能和轨迹特性。合理的
// @Description 参数设置可以在保证运动精度的同时优化运动时间。
//
// @Param params *PTPLParams true "PTPL运动参数：
//   - velocity: 速度（单位：mm/s）
//   - acceleration: 加速度（单位：mm/s²）
//     注意：参数必须在机械臂规格范围内"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 指令加入队列，按顺序执行
//   - false: 立即执行指令
//     注意：建议使用队列模式以确保参数设置顺序"
//
// @Return uint64 "指令索引：
//   - 返回指令队列索引
//   - 仅在isQueued为true时有效"
//
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数无效
//   - 参数超出范围
//   - 机械臂被锁定
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置PTPL运动参数
//	params := &PTPLParams{
//	    Velocity:     200,  // 速度200mm/s
//	    Acceleration: 800,  // 加速度800mm/s²
//	}
//	index, err := dobot.SetPTPLParams(params, true)
//	if err != nil {
//	    log.Printf("设置PTPL运动参数失败: %v", err)
//	    return
//	}
//	log.Printf("PTPL运动参数设置成功，指令索引: %d", index)
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

	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
	return queuedCmdIndex, nil
}

// GetPTPLParams 获取PTPL运动参数
// @Summary 获取机械臂PTPL运动的速度和加速度参数
// @Description 获取机械臂在PTPL（点到点线性）运动模式下当前的速度和加速度
// @Description 参数。通过这些参数可以了解机械臂在直线运动时的性能配置，
// @Description 便于调整和优化运动控制。
//
// @Return *PTPLParams "PTPL运动参数：
//   - velocity: 速度（单位：mm/s）
//   - acceleration: 加速度（单位：mm/s²）
//     注意：返回当前实际设置的参数值"
//
// @Return error "错误信息"
// @Success 200 {object} *PTPLParams "返回PTPL运动参数结构体"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	// 获取PTPL运动参数
//	params, err := dobot.GetPTPLParams()
//	if err != nil {
//	    log.Printf("获取PTPL运动参数失败: %v", err)
//	    return
//	}
//	log.Printf("当前PTPL运动参数：")
//	log.Printf("  速度: %.1f mm/s", params.Velocity)
//	log.Printf("  加速度: %.1f mm/s²", params.Acceleration)
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

	params := &PTPLParams{}
	if err := binary.Read(resp.Reader(), binary.LittleEndian, params); err != nil {
		return nil, fmt.Errorf("failed to read PTP L params: %v", err)
	}
	return params, nil
}

// SetPTPJumpParams 设置PTP跳跃参数
// @Summary 设置机械臂PTP跳跃运动的参数
// @Description 设置机械臂在PTP跳跃模式下的运动参数。跳跃模式是一种特殊的
// @Description 运动模式，机械臂会先抬升到指定高度，然后平移到目标位置上方，
// @Description 最后下降到目标位置。这种模式适用于需要避开中间障碍物的场景。
//
// @Param params *PTPJumpParams true "跳跃运动参数：
//   - jumpHeight: 跳跃高度（单位：mm）
//   - zLimit: Z轴限位（单位：mm）
//     注意：确保参数设置合理，避免碰撞"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 指令加入队列，按顺序执行
//   - false: 立即执行指令
//     注意：建议使用队列模式以确保参数设置顺序"
//
// @Return uint64 "指令索引：
//   - 返回指令队列索引
//   - 仅在isQueued为true时有效"
//
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数无效
//   - 参数超出范围
//   - 机械臂被锁定
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置PTP跳跃参数
//	params := &PTPJumpParams{
//	    JumpHeight: 50,   // 跳跃高度50mm
//	    ZLimit:     200,  // Z轴限位200mm
//	}
//	index, err := dobot.SetPTPJumpParams(params, true)
//	if err != nil {
//	    log.Printf("设置PTP跳跃参数失败: %v", err)
//	    return
//	}
//	log.Printf("PTP跳跃参数设置成功，指令索引: %d", index)
func (dobot *Dobot) SetPTPJumpParams(params *PTPJumpParams, isQueued bool) (queuedCmdIndex uint64, err error) {
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

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
	return queuedCmdIndex, nil
}

// GetPTPJumpParams 获取PTP跳跃参数
// @Summary 获取机械臂PTP跳跃运动的参数
// @Description 获取机械臂在PTP跳跃模式下当前的运动参数。通过这些参数可以
// @Description 了解机械臂在跳跃运动时的高度和限位设置，便于调整和优化运动
// @Description 轨迹。
//
// @Return *PTPJumpParams "跳跃运动参数：
//   - jumpHeight: 跳跃高度（单位：mm）
//   - zLimit: Z轴限位（单位：mm）
//     注意：返回当前实际设置的参数值"
//
// @Return error "错误信息"
// @Success 200 {object} *PTPJumpParams "返回跳跃运动参数结构体"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	// 获取PTP跳跃参数
//	params, err := dobot.GetPTPJumpParams()
//	if err != nil {
//	    log.Printf("获取PTP跳跃参数失败: %v", err)
//	    return
//	}
//	log.Printf("当前PTP跳跃参数：")
//	log.Printf("  跳跃高度: %.1f mm", params.JumpHeight)
//	log.Printf("  Z轴限位: %.1f mm", params.ZLimit)
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

	params := &PTPJumpParams{}
	if err := binary.Read(resp.Reader(), binary.LittleEndian, params); err != nil {
		return nil, fmt.Errorf("failed to read PTP jump params: %v", err)
	}
	return params, nil
}

// SetPTPJump2Params 设置PTP跳跃2参数
// @Summary 设置机械臂PTP跳跃2运动的参数
// @Description 设置机械臂在PTP跳跃2模式下的运动参数。跳跃2模式是跳跃模式
// @Description 的扩展，允许设置更多的参数来精确控制跳跃轨迹。这种模式提供
// @Description 了更灵活的障碍物避让能力。
//
// @Param params *PTPJump2Params true "跳跃2运动参数：
//   - startJumpHeight: 起点跳跃高度（单位：mm）
//   - endJumpHeight: 终点跳跃高度（单位：mm）
//   - zLimit: Z轴限位（单位：mm）
//     注意：确保参数设置合理，避免碰撞"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 指令加入队列，按顺序执行
//   - false: 立即执行指令
//     注意：建议使用队列模式以确保参数设置顺序"
//
// @Return uint64 "指令索引：
//   - 返回指令队列索引
//   - 仅在isQueued为true时有效"
//
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数无效
//   - 参数超出范围
//   - 机械臂被锁定
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置PTP跳跃2参数
//	params := &PTPJump2Params{
//	    StartJumpHeight: 30,   // 起点跳跃高度30mm
//	    EndJumpHeight:   50,   // 终点跳跃高度50mm
//	    ZLimit:          200,  // Z轴限位200mm
//	}
//	index, err := dobot.SetPTPJump2Params(params, true)
//	if err != nil {
//	    log.Printf("设置PTP跳跃2参数失败: %v", err)
//	    return
//	}
//	log.Printf("PTP跳跃2参数设置成功，指令索引: %d", index)
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

	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
	return queuedCmdIndex, nil
}

// GetPTPJump2Params 获取PTP跳跃2参数
// @Summary 获取机械臂PTP跳跃2运动的参数
// @Description 获取机械臂在PTP跳跃2模式下当前的运动参数。通过这些参数可以
// @Description 了解机械臂在跳跃运动时的起点和终点高度设置，以及限位设置，
// @Description 便于调整和优化运动轨迹。
//
// @Return *PTPJump2Params "跳跃2运动参数：
//   - startJumpHeight: 起点跳跃高度（单位：mm）
//   - endJumpHeight: 终点跳跃高度（单位：mm）
//   - zLimit: Z轴限位（单位：mm）
//     注意：返回当前实际设置的参数值"
//
// @Return error "错误信息"
// @Success 200 {object} *PTPJump2Params "返回跳跃2运动参数结构体"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	// 获取PTP跳跃2参数
//	params, err := dobot.GetPTPJump2Params()
//	if err != nil {
//	    log.Printf("获取PTP跳跃2参数失败: %v", err)
//	    return
//	}
//	log.Printf("当前PTP跳跃2参数：")
//	log.Printf("  起点跳跃高度: %.1f mm", params.StartJumpHeight)
//	log.Printf("  终点跳跃高度: %.1f mm", params.EndJumpHeight)
//	log.Printf("  Z轴限位: %.1f mm", params.ZLimit)
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

	params := &PTPJump2Params{}
	if err := binary.Read(resp.Reader(), binary.LittleEndian, params); err != nil {
		return nil, fmt.Errorf("failed to read PTP jump2 params: %v", err)
	}
	return params, nil
}

// SetPTPCommonParams 设置PTP通用参数
// @Summary 设置机械臂PTP运动的通用参数
// @Description 设置机械臂在PTP（点到点）运动模式下的通用参数。这些参数
// @Description 会影响所有PTP运动模式（包括关节运动、直线运动等）的基本
// @Description 特性，如运动的平滑度和响应性。
//
// @Param params *PTPCommonParams true "PTP通用参数：
//   - velocityRatio: 速度比例，范围[0-100]
//   - accelerationRatio: 加速度比例，范围[0-100]
//     注意：参数设置过大可能导致运动不稳定，建议从小到大逐步调整"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 指令加入队列，按顺序执行
//   - false: 立即执行指令
//     注意：建议使用队列模式以确保参数设置顺序"
//
// @Return uint64 "指令索引：
//   - 返回指令队列索引
//   - 仅在isQueued为true时有效"
//
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数无效
//   - 参数超出范围
//   - 机械臂被锁定
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置PTP通用运动参数
//	params := &PTPCommonParams{
//	    VelocityRatio:     50,  // 速度比例50%
//	    AccelerationRatio: 50,  // 加速度比例50%
//	}
//	index, err := dobot.SetPTPCommonParams(params, true)
//	if err != nil {
//	    log.Printf("设置PTP通用参数失败: %v", err)
//	    return
//	}
//	log.Printf("PTP通用参数设置成功，指令索引: %d", index)
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

	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
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

	params := &PTPCommonParams{}
	if err := binary.Read(resp.Reader(), binary.LittleEndian, params); err != nil {
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
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
	return queuedCmdIndex, nil
}

// SetPTPWithLCmd 设置带L轴的PTP运动指令
// @Summary 设置带L轴的机械臂PTP运动指令
// @Description 控制带有L轴（第五轴）的机械臂执行PTP（点到点）运动。除了
// @Description 标准的XYZ和R轴运动外，还可以控制L轴的位置。这适用于需要
// @Description 额外自由度的应用场景。
//
// @Param cmd *PTPWithLCmd true "PTP运动指令参数：
//   - ptpMode: 运动模式（MOVJ/MOVL/JUMP/JUMPXY）
//   - x: X轴目标位置（单位：mm）
//   - y: Y轴目标位置（单位：mm）
//   - z: Z轴目标位置（单位：mm）
//   - r: R轴目标角度（单位：度）
//   - l: L轴目标位置（单位：mm）
//     注意：确保目标位置在工作空间内"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 指令加入队列，按顺序执行
//   - false: 立即执行指令
//     注意：建议使用队列模式以确保运动顺序"
//
// @Return uint64 "指令索引：
//   - 返回指令队列索引
//   - 仅在isQueued为true时有效"
//
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令索引"
// @Failure 400 {error} "执行失败，可能的错误：
//   - 参数无效
//   - 目标位置不可达
//   - L轴未启用
//   - 机械臂被锁定
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 执行带L轴的点到点运动
//	cmd := &PTPWithLCmd{
//	    PTPMode: PTPMode_MOVJ,           // 关节运动模式
//	    X: 200, Y: 0, Z: 100, R: 0,      // 目标位置和姿态
//	    L: 50,                           // L轴目标位置
//	}
//	index, err := dobot.SetPTPWithLCmd(cmd, true)
//	if err != nil {
//	    log.Printf("执行带L轴的PTP运动失败: %v", err)
//	    return
//	}
//	log.Printf("带L轴的PTP运动开始执行，指令索引: %d", index)
//
//	// 等待运动完成
//	for {
//	    finished, _ := dobot.GetQueuedCmdMotionFinish()
//	    if finished {
//	        log.Printf("带L轴的PTP运动完成")
//	        break
//	    }
//	    time.Sleep(100 * time.Millisecond)
//	}
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
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
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

	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
	return queuedCmdIndex, nil
}

// SetCPCmd 设置连续运动命令
// @Summary 设置机械臂连续路径运动指令
// @Description 通过CP（连续路径）运动命令控制机械臂执行连续轨迹运动。
// @Description CP模式下，机械臂会按照设定的路径点平滑运动，适用于需要
// @Description 连续轨迹的应用场景，如画线、涂胶等。
//
// @Param cmd *CPCmd true "CP运动命令参数：
//   - cpMode: CP运动模式
//   - x,y,z,r: 目标位置坐标
//   - velocity: 运动速度
//   - acceleration: 加速度
//     注意：确保运动参数在安全范围内"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 将指令加入队列，按顺序执行
//   - false: 立即执行该指令
//   - 建议使用队列模式以确保运动连续性"
//
// @Return uint64 "指令队列索引（当isQueued为true时有效）"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数为空
//   - 运动参数无效
//   - 机械臂被锁定
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置CP运动指令
//	cmd := &CPCmd{
//	    CPMode: 0,
//	    X: 200, Y: 0, Z: 50, R: 0,
//	    Velocity: 100,    // 速度100mm/s
//	    Acceleration: 200, // 加速度200mm/s²
//	}
//	index, err := dobot.SetCPCmd(cmd, true)
//	if err != nil {
//	    log.Printf("设置CP运动指令失败: %v", err)
//	} else {
//	    log.Printf("CP运动指令设置成功，指令索引: %d", index)
//	}
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

	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
	return queuedCmdIndex, nil
}

// SetCPLECmd 设置连续运动扩展命令
// @Summary 设置机械臂连续路径运动扩展指令
// @Description 通过CPLE（连续路径扩展）运动命令控制机械臂执行连续轨迹
// @Description 运动。与基本CP命令相比，扩展命令增加了功率参数控制，可以
// @Description 更精细地控制运动过程中的功率输出。
//
// @Param cpMode uint8 true "CP运动模式：
//   - 用于指定连续路径运动的具体模式
//   - 具体值含义参见产品手册"
//
// @Param x float32 true "目标X坐标（单位：mm）"
// @Param y float32 true "目标Y坐标（单位：mm）"
// @Param z float32 true "目标Z坐标（单位：mm）"
// @Param power float32 true "功率参数：
//   - 范围：0-100
//   - 用于控制运动过程中的功率输出
//     注意：功率设置过大可能影响运动精度"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 将指令加入队列，按顺序执行
//   - false: 立即执行该指令
//   - 建议使用队列模式以确保运动连续性"
//
// @Return uint64 "指令队列索引（当isQueued为true时有效）"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数无效
//   - 运动模式不支持
//   - 坐标超出范围
//   - 功率参数超出范围
//   - 机械臂被锁定
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置CPLE运动指令
//	index, err := dobot.SetCPLECmd(0, 200, 0, 50, 80, true)
//	if err != nil {
//	    log.Printf("设置CPLE运动指令失败: %v", err)
//	} else {
//	    log.Printf("CPLE运动指令设置成功，指令索引: %d", index)
//	}
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

	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
	return queuedCmdIndex, nil
}

// SetCPRHoldEnable 设置CPR保持使能
// @Summary 设置机械臂CP运动保持功能的使能状态
// @Description 设置机械臂在CP（连续路径）运动模式下保持功能的使能状态。
// @Description CP运动保持功能可以让机械臂在CP运动过程中保持特定的状态或
// @Description 参数，提高运动的稳定性和精确度。
//
// @Param isEnable bool true "是否启用CP运动保持功能：
//   - true: 启用保持功能
//   - false: 禁用保持功能
//     注意：启用后可能会影响运动速度"
//
// @Return error "错误信息"
// @Success 200 {string} "设置成功"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 机械臂被锁定
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 启用CP运动保持功能
//	err := dobot.SetCPRHoldEnable(true)
//	if err != nil {
//	    log.Printf("设置CP运动保持功能失败: %v", err)
//	} else {
//	    log.Printf("CP运动保持功能已启用")
//	}
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
// @Summary 获取机械臂CP运动保持功能的使能状态
// @Description 获取机械臂在CP（连续路径）运动模式下保持功能的使能状态。
// @Description CP运动保持功能可以让机械臂在CP运动过程中保持特定的状态或
// @Description 参数，提高运动的稳定性和精确度。
//
// @Return bool "使能状态：
//   - true: CP运动保持功能已启用
//   - false: CP运动保持功能已禁用"
//
// @Return error "错误信息"
// @Success 200 {boolean} bool "返回使能状态"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	// 获取CP运动保持功能状态
//	enabled, err := dobot.GetCPRHoldEnable()
//	if err != nil {
//	    log.Printf("获取CP运动保持状态失败: %v", err)
//	} else {
//	    if enabled {
//	        log.Printf("CP运动保持功能已启用")
//	    } else {
//	        log.Printf("CP运动保持功能已禁用")
//	    }
//	}
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
	return resp.Data()[0] != 0, nil
}

// SetCPCommonParams 设置CP通用参数
// @Summary 设置机械臂CP运动的通用参数
// @Description 设置机械臂在CP（连续路径）运动模式下的通用参数。这些参数
// @Description 影响所有CP运动的基本特性，如加加速度、加速度等。合理的参数
// @Description 设置可以优化运动的平滑度和精确度。
//
// @Param params *CPCommonParams true "CP通用参数：
//   - planAcc: 规划加速度
//   - junctionVel: 拐点速度
//   - acc: 加速度
//   - realTimeTrack: 实时轨迹开关
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
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置CP运动通用参数
//	params := &CPCommonParams{
//	    PlanAcc:       100,   // 规划加速度
//	    JunctionVel:   50,    // 拐点速度
//	    Acc:           200,   // 加速度
//	    RealTimeTrack: false, // 关闭实时轨迹
//	}
//	index, err := dobot.SetCPCommonParams(params, true)
//	if err != nil {
//	    log.Printf("设置CP通用参数失败: %v", err)
//	} else {
//	    log.Printf("CP通用参数设置成功，指令索引: %d", index)
//	}
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

	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
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

	params := &CPCommonParams{}
	if err := binary.Read(resp.Reader(), binary.LittleEndian, params); err != nil {
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
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
	return queuedCmdIndex, nil
}

// GetARCParams 获取ARC参数
// @Summary 获取机械臂ARC运动参数
// @Description 获取机械臂在ARC（圆弧）模式下的运动参数。这些参数控制
// @Description 机械臂在执行圆弧轨迹运动时的速度和加速度特性。可用于
// @Description 确认当前的运动参数配置。
//
// @Return *ARCParams "ARC运动参数：
//   - xyzVelocity: XYZ轴速度（单位：mm/s）
//   - rVelocity: R轴速度（单位：°/s）
//   - xyzAcceleration: XYZ轴加速度（单位：mm/s²）
//   - rAcceleration: R轴加速度（单位：°/s²）"
//
// @Return error "错误信息"
// @Success 200 {object} *ARCParams "返回ARC参数结构体"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效
//   - 数据解析错误"
//
// @Example
//
//	// 获取当前ARC运动参数
//	params, err := dobot.GetARCParams()
//	if err != nil {
//	    log.Printf("获取ARC参数失败: %v", err)
//	} else {
//	    log.Printf("当前ARC参数：")
//	    log.Printf("  XYZ轴速度: %.2f mm/s", params.XYZVelocity)
//	    log.Printf("  R轴速度: %.2f °/s", params.RVelocity)
//	    log.Printf("  XYZ轴加速度: %.2f mm/s²", params.XYZAcceleration)
//	    log.Printf("  R轴加速度: %.2f °/s²", params.RAcceleration)
//	}
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

	params := &ARCParams{}
	if err := binary.Read(resp.Reader(), binary.LittleEndian, params); err != nil {
		return nil, fmt.Errorf("failed to read ARC params: %v", err)
	}
	return params, nil
}

// SetARCCmd 设置ARC命令
// @Summary 设置机械臂ARC运动指令
// @Description 发送ARC（圆弧）运动指令给机械臂。通过指定圆弧运动的
// @Description 关键点，控制机械臂执行圆弧轨迹运动。圆弧由起点（当前
// @Description 位置）和两个路径点确定。
//
// @Param cmd *ARCCmd true "ARC运动指令参数：
//   - point1: 第一个路径点坐标（x,y,z,r）
//   - point2: 第二个路径点坐标（x,y,z,r）
//     注意：
//   - 起点为机械臂当前位置
//   - 三点不能共线，否则无法确定圆弧"
//
// @Return uint64 "指令队列索引"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数为空
//   - 点位坐标无效
//   - 三点共线
//   - 圆弧半径过大
//   - 机械臂被锁定
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置ARC运动指令
//	cmd := &ARCCmd{
//	    Point1: Point{X: 200, Y: 0, Z: 0, R: 0},    // 第一个路径点
//	    Point2: Point{X: 200, Y: 200, Z: 0, R: 0},  // 第二个路径点
//	}
//	index, err := dobot.SetARCCmd(cmd)
//	if err != nil {
//	    log.Printf("设置ARC运动指令失败: %v", err)
//	} else {
//	    log.Printf("ARC运动指令设置成功，指令索引: %d", index)
//	}
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

	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
	return queuedCmdIndex, nil
}

// SetCircleCmd 设置圆周运动命令
// @Summary 设置机械臂圆周运动指令
// @Description 发送圆周运动指令给机械臂。通过指定圆周运动的关键点和圈数，
// @Description 控制机械臂执行完整的圆周运动。圆周由起点（当前位置）和两个
// @Description 路径点确定，机械臂将按指定圈数重复运动。
//
// @Param cmd *CircleCmd true "圆周运动指令参数：
//   - point1: 第一个路径点坐标（x,y,z,r）
//   - point2: 第二个路径点坐标（x,y,z,r）
//   - count: 圆周运动圈数
//     注意：
//   - 起点为机械臂当前位置
//   - 三点不能共线，否则无法确定圆周
//   - 圈数必须为正整数"
//
// @Return uint64 "指令队列索引"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数为空
//   - 点位坐标无效
//   - 三点共线
//   - 圈数无效
//   - 圆周半径过大
//   - 机械臂被锁定
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置圆周运动指令，运动3圈
//	cmd := &CircleCmd{
//	    Point1: Point{X: 200, Y: 0, Z: 0, R: 0},    // 第一个路径点
//	    Point2: Point{X: 200, Y: 200, Z: 0, R: 0},  // 第二个路径点
//	    Count:  3,                                   // 运动3圈
//	}
//	index, err := dobot.SetCircleCmd(cmd)
//	if err != nil {
//	    log.Printf("设置圆周运动指令失败: %v", err)
//	} else {
//	    log.Printf("圆周运动指令设置成功，指令索引: %d", index)
//	}
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

	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
	return queuedCmdIndex, nil
}

// SetARCCommonParams 设置ARC通用参数
// @Summary 设置机械臂ARC运动的通用参数
// @Description 设置机械臂在ARC（圆弧）运动模式下的通用参数。这些参数用于
// @Description 调整所有ARC运动的基本特性，如速度比例和加速度比例。合理的
// @Description 参数设置可以优化运动性能和精度。
//
// @Param params *ARCCommonParams true "ARC通用参数：
//   - velocityRatio: 速度比例（范围：0-100）
//   - accelerationRatio: 加速度比例（范围：0-100）
//     注意：
//   - 比例值越大，实际运动速度和加速度越大
//   - 设置过大可能影响运动精度"
//
// @Return uint64 "指令队列索引"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数为空
//   - 比例值超出范围
//   - 机械臂被锁定
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置ARC通用参数
//	params := &ARCCommonParams{
//	    VelocityRatio:     50,  // 速度比例50%
//	    AccelerationRatio: 50,  // 加速度比例50%
//	}
//	index, err := dobot.SetARCCommonParams(params)
//	if err != nil {
//	    log.Printf("设置ARC通用参数失败: %v", err)
//	} else {
//	    log.Printf("ARC通用参数设置成功，指令索引: %d", index)
//	}
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

	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
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

	params := &ARCCommonParams{}
	if err := binary.Read(resp.Reader(), binary.LittleEndian, params); err != nil {
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

	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
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

	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
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

	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
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

	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
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

	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
	return queuedCmdIndex, nil
}

// SetEMotor 设置扩展电机参数
// @Summary 设置扩展电机的控制参数
// @Description 控制机械臂扩展接口上连接的电机。可以设置电机的运行模式、
// @Description 速度、加速度等参数。扩展电机可用于传送带、转盘等外部设备
// @Description 的驱动控制。
//
// @Param params *EMotor true "扩展电机参数：
//   - index: 电机索引（从0开始）
//   - enabled: 是否使能电机
//   - speed: 电机速度（单位取决于电机类型）
//   - direction: 旋转方向（0:正向，1:反向）
//     注意：参数设置需要与实际电机规格相匹配"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 将指令加入队列，按顺序执行
//   - false: 立即执行该指令
//   - 建议使用队列模式以确保控制顺序"
//
// @Return uint64 "指令队列索引（当isQueued为true时有效）"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数为空
//   - 电机索引无效
//   - 速度超出范围
//   - 电机未连接
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置电机0正向运行
//	params := &EMotor{
//	    Index:     0,     // 电机0
//	    Enabled:   true,  // 使能电机
//	    Speed:     1000,  // 设置速度
//	    Direction: 0,     // 正向旋转
//	}
//	index, err := dobot.SetEMotor(params, true)
//	if err != nil {
//	    log.Printf("设置扩展电机失败: %v", err)
//	} else {
//	    log.Printf("扩展电机设置成功，指令索引: %d", index)
//	}
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

	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
	return queuedCmdIndex, nil
}

// SetEMotorS 设置扩展步进电机参数
// @Summary 设置多个扩展步进电机的控制参数
// @Description 同时控制多个连接在机械臂扩展接口上的步进电机。可以批量设置
// @Description 电机的使能状态、速度、方向等参数。这种批量设置方式适用于需要
// @Description 多个电机协同工作的场合。
//
// @Param params *EMotorS true "扩展步进电机参数：
//   - count: 要控制的电机数量
//   - index[8]: 电机索引数组（从0开始）
//   - enabled[8]: 电机使能状态数组
//   - speed[8]: 电机速度数组
//   - direction[8]: 电机方向数组（0:正向，1:反向）
//     注意：
//   - 数组长度固定为8，未使用的位置填0
//   - count必须小于等于8"
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 将指令加入队列，按顺序执行
//   - false: 立即执行该指令
//   - 建议使用队列模式以确保控制顺序"
//
// @Return uint64 "指令队列索引（当isQueued为true时有效）"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数为空
//   - 电机数量无效
//   - 电机索引无效
//   - 速度超出范围
//   - 电机未连接
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 同时设置两个步进电机
//	params := &EMotorS{
//	    Count: 2,                    // 控制2个电机
//	    Index:     [8]int{0, 1},     // 电机0和1
//	    Enabled:   [8]bool{true, true},
//	    Speed:     [8]int{1000, 1000},
//	    Direction: [8]int{0, 0},      // 都设为正向
//	}
//	index, err := dobot.SetEMotorS(params, true)
//	if err != nil {
//	    log.Printf("设置步进电机失败: %v", err)
//	} else {
//	    log.Printf("步进电机设置成功，指令索引: %d", index)
//	}
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
	queuedCmdIndex = binary.LittleEndian.Uint64(resp.Data())
	return queuedCmdIndex, nil
}

// SetColorSensor 设置颜色传感器
// @Summary 设置机械臂颜色传感器的工作状态
// @Description 配置机械臂上的颜色传感器。颜色传感器可用于识别和分类不同
// @Description 颜色的物体，在分拣、质检等应用中非常有用。使用前需要正确
// @Description 连接传感器并设置其工作参数。
//
// @Param enable bool true "传感器使能状态：
//   - true: 启用传感器
//   - false: 禁用传感器"
//
// @Param colorPort ColorPort true "传感器端口：
//   - ColorPort_PORT1: 端口1
//   - ColorPort_PORT2: 端口2
//     注意：具体可用端口取决于机械臂型号"
//
// @Param version uint8 true "传感器版本号：
//   - 用于适配不同版本的传感器
//   - 具体值参见传感器说明"
//
// @Return error "错误信息"
// @Success 200 {string} "设置成功"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 传感器未连接
//   - 端口不支持
//   - 版本不匹配
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 在端口1启用颜色传感器
//	err := dobot.SetColorSensor(true, ColorPort_PORT1, 1)
//	if err != nil {
//	    log.Printf("设置颜色传感器失败: %v", err)
//	} else {
//	    log.Printf("颜色传感器设置成功")
//	}
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
// @Summary 获取机械臂颜色传感器的RGB数据
// @Description 获取当前颜色传感器检测到的RGB颜色值。这些数据可用于
// @Description 颜色识别、物体分类等应用。使用前需要确保颜色传感器已
// @Description 正确配置并启用。
//
// @Return uint8 "R值：
//   - 红色分量值
//   - 范围：0-255"
//
// @Return uint8 "G值：
//   - 绿色分量值
//   - 范围：0-255"
//
// @Return uint8 "B值：
//   - 蓝色分量值
//   - 范围：0-255"
//
// @Return error "错误信息"
// @Success 200 {number} uint8 "返回RGB三个分量值"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 传感器未启用
//   - 传感器未连接
//   - 通信错误
//   - 响应数据无效"
//
// @Example
//
//	// 获取颜色传感器数据
//	r, g, b, err := dobot.GetColorSensor()
//	if err != nil {
//	    log.Printf("获取颜色传感器数据失败: %v", err)
//	} else {
//	    log.Printf("当前检测到的颜色：")
//	    log.Printf("  R: %d", r)
//	    log.Printf("  G: %d", g)
//	    log.Printf("  B: %d", b)
//	}
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
	data := resp.Data()
	return data[0], data[1], data[2], nil
}

// SetInfraredSensor 设置红外传感器
// @Summary 设置机械臂红外传感器的工作状态
// @Description 配置机械臂上的红外传感器。红外传感器可用于检测物体的存在
// @Description 和距离，在避障、物体检测等应用中非常有用。使用前需要正确
// @Description 连接传感器并设置其工作参数。
//
// @Param enable bool true "传感器使能状态：
//   - true: 启用传感器
//   - false: 禁用传感器"
//
// @Param infraredPort InfraredPort true "传感器端口：
//   - InfraredPort_PORT1: 端口1
//   - InfraredPort_PORT2: 端口2
//     注意：具体可用端口取决于机械臂型号"
//
// @Param version uint8 true "传感器版本号：
//   - 用于适配不同版本的传感器
//   - 具体值参见传感器说明"
//
// @Return error "错误信息"
// @Success 200 {string} "设置成功"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 传感器未连接
//   - 端口不支持
//   - 版本不匹配
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 在端口1启用红外传感器
//	err := dobot.SetInfraredSensor(true, InfraredPort_PORT1, 1)
//	if err != nil {
//	    log.Printf("设置红外传感器失败: %v", err)
//	} else {
//	    log.Printf("红外传感器设置成功")
//	}
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
// @Summary 获取机械臂红外传感器的数据
// @Description 获取指定端口红外传感器的检测数据。红外传感器可用于检测
// @Description 物体的存在和距离，数据值通常与物体的距离或反射特性相关。
// @Description 使用前需要确保红外传感器已正确配置并启用。
//
// @Param port InfraredPort true "红外传感器端口：
//   - InfraredPort_PORT1: 端口1
//   - InfraredPort_PORT2: 端口2
//     注意：具体可用端口取决于机械臂型号"
//
// @Return uint8 "传感器数据：
//   - 数值范围：0-255
//   - 具体含义取决于传感器型号和应用场景
//   - 通常值越大表示检测到的物体越近"
//
// @Return error "错误信息"
// @Success 200 {number} uint8 "返回传感器数据值"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 传感器未启用
//   - 传感器未连接
//   - 端口无效
//   - 通信错误
//   - 响应数据无效"
//
// @Example
//
//	// 获取端口1的红外传感器数据
//	value, err := dobot.GetInfraredSensor(InfraredPort_PORT1)
//	if err != nil {
//	    log.Printf("获取红外传感器数据失败: %v", err)
//	} else {
//	    log.Printf("当前红外传感器数据: %d", value)
//	}
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
	data := resp.Data()
	return data[0], nil
}

// SetAngleSensorStaticError 设置角度传感器静态误差
// @Summary 设置机械臂关节角度传感器的静态误差补偿
// @Description 设置机械臂后臂和前臂关节角度传感器的静态误差补偿值。
// @Description 这些补偿值用于修正传感器的系统误差，提高角度测量的准确性。
// @Description 通常在标定过程中设置，不建议随意修改。
//
// @Param rearArmAngleError float32 true "后臂角度静态误差：
//   - 单位：度（°）
//   - 范围：通常在±1°以内
//     注意：补偿值过大可能表示传感器异常"
//
// @Param frontArmAngleError float32 true "前臂角度静态误差：
//   - 单位：度（°）
//   - 范围：通常在±1°以内
//     注意：补偿值过大可能表示传感器异常"
//
// @Return error "错误信息"
// @Success 200 {string} "设置成功"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 补偿值超出范围
//   - 机械臂被锁定
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置角度传感器静态误差补偿
//	err := dobot.SetAngleSensorStaticError(0.1, 0.2)
//	if err != nil {
//	    log.Printf("设置角度传感器静态误差失败: %v", err)
//	} else {
//	    log.Printf("角度传感器静态误差设置成功")
//	}
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
// @Summary 获取机械臂关节角度传感器的静态误差补偿值
// @Description 获取机械臂后臂和前臂关节角度传感器当前的静态误差补偿值。
// @Description 通过查看当前的补偿值，可以了解传感器的校准状态。这些值
// @Description 通常在标定时设置。
//
// @Return float32 "后臂角度静态误差：
//   - 单位：度（°）
//   - 范围：通常在±1°以内
//     注意：补偿值过大可能表示传感器异常"
//
// @Return float32 "前臂角度静态误差：
//   - 单位：度（°）
//   - 范围：通常在±1°以内
//     注意：补偿值过大可能表示传感器异常"
//
// @Return error "错误信息"
// @Success 200 {number} float32 "返回后臂和前臂的角度静态误差补偿值"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	// 获取角度传感器静态误差补偿值
//	rearError, frontError, err := dobot.GetAngleSensorStaticError()
//	if err != nil {
//	    log.Printf("获取角度传感器静态误差失败: %v", err)
//	} else {
//	    log.Printf("当前静态误差补偿值：")
//	    log.Printf("  后臂: %.2f°", rearError)
//	    log.Printf("  前臂: %.2f°", frontError)
//	}
func (dobot *Dobot) GetAngleSensorStaticError() (float32, float32, error) {
	message := &Message{
		Id:       ProtocolAngleSensorStaticError,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
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
// @Summary 设置机械臂关节角度传感器的校准系数
// @Description 设置机械臂后臂和前臂关节角度传感器的校准系数。这些系数用于
// @Description 校正传感器的线性误差，通过调整比例系数来提高角度测量的准确性。
// @Description 通常在出厂标定时设置，不建议用户随意修改。
//
// @Param rearArmAngleCoef float32 true "后臂角度系数：
//   - 用于校正后臂角度传感器的线性误差
//   - 范围：通常在0.9-1.1之间
//     注意：系数偏离1过大可能表示传感器异常"
//
// @Param frontArmAngleCoef float32 true "前臂角度系数：
//   - 用于校正前臂角度传感器的线性误差
//   - 范围：通常在0.9-1.1之间
//     注意：系数偏离1过大可能表示传感器异常"
//
// @Return error "错误信息"
// @Success 200 {string} "设置成功"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 系数值超出范围
//   - 机械臂被锁定
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置角度传感器校准系数
//	err := dobot.SetAngleSensorCoef(1.02, 0.98)
//	if err != nil {
//	    log.Printf("设置角度传感器系数失败: %v", err)
//	} else {
//	    log.Printf("角度传感器系数设置成功")
//	}
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
// @Description 获取机械臂后臂和前臂关节角度传感器当前的校准系数。这些系数
// @Description 用于校正传感器的线性误差，通过查看当前的系数值可以了解传感器
// @Description 的校准状态。
//
// @Return float32 "后臂角度系数：
//   - 用于校正后臂角度传感器的线性误差
//   - 范围：通常在0.9-1.1之间
//     注意：系数偏离1过大可能表示传感器异常"
//
// @Return float32 "前臂角度系数：
//   - 用于校正前臂角度传感器的线性误差
//   - 范围：通常在0.9-1.1之间
//     注意：系数偏离1过大可能表示传感器异常"
//
// @Return error "错误信息"
// @Success 200 {number} float32 "返回后臂和前臂的角度校准系数"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	// 获取角度传感器校准系数
//	rearCoef, frontCoef, err := dobot.GetAngleSensorCoef()
//	if err != nil {
//	    log.Printf("获取角度传感器系数失败: %v", err)
//	} else {
//	    log.Printf("当前校准系数：")
//	    log.Printf("  后臂: %.3f", rearCoef)
//	    log.Printf("  前臂: %.3f", frontCoef)
//	}
func (dobot *Dobot) GetAngleSensorCoef() (float32, float32, error) {
	message := &Message{
		Id:       ProtocolAngleSensorCoef,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
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
// @Summary 设置机械臂底座编码器的静态误差补偿
// @Description 设置机械臂底座旋转编码器的静态误差补偿值。这个补偿值用于
// @Description 修正底座旋转时的系统误差，提高旋转定位的准确性。通常在
// @Description 标定过程中设置，不建议随意修改。
//
// @Param baseDecoderError float32 true "底座编码器静态误差：
//   - 单位：度（°）
//   - 范围：通常在±1°以内
//     注意：补偿值过大可能导致旋转定位不准"
//
// @Return error "错误信息"
// @Success 200 {string} "设置成功"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 补偿值超出范围
//   - 机械臂被锁定
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置底座编码器静态误差补偿
//	err := dobot.SetBaseDecoderStaticError(0.1)
//	if err != nil {
//	    log.Printf("设置底座编码器静态误差失败: %v", err)
//	} else {
//	    log.Printf("底座编码器静态误差设置成功")
//	}
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
// @Summary 获取机械臂底座编码器的静态误差补偿值
// @Description 获取机械臂底座旋转编码器当前设置的静态误差补偿值。这个值
// @Description 用于修正底座旋转时的系统误差，通过查看当前的补偿值可以了解
// @Description 底座旋转精度的校准状态。
//
// @Return float32 "底座编码器静态误差：
//   - 单位：度（°）
//   - 范围：通常在±1°以内
//     注意：补偿值过大可能表示校准异常"
//
// @Return error "错误信息"
// @Success 200 {number} float32 "返回底座编码器的静态误差补偿值"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	// 获取底座编码器静态误差补偿值
//	error, err := dobot.GetBaseDecoderStaticError()
//	if err != nil {
//	    log.Printf("获取底座编码器静态误差失败: %v", err)
//	} else {
//	    log.Printf("当前底座编码器静态误差: %.2f°", error)
//	}
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
	reader := resp.Reader()
	var baseDecoderError float32
	binary.Read(reader, binary.LittleEndian, &baseDecoderError)
	return baseDecoderError, nil
}

// SetLRHandCalibrateValue 设置左右手校准值
// @Summary 设置机械臂左右手模式的校准参数
// @Description 设置机械臂在左右手模式切换时的校准参数。这些参数用于确保
// @Description 机械臂在左右手模式切换后的位置精度。通常在出厂标定时设置，
// @Description 不建议用户随意修改。
//
// @Param lrHandCalibrateValue float32 true "左右手校准值：
//   - 用于补偿左右手切换时的位置偏差
//   - 单位：毫米（mm）
//   - 范围：通常在±5mm以内
//     注意：校准值过大可能表示机械结构异常"
//
// @Return error "错误信息"
// @Success 200 {string} "设置成功"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 校准值超出范围
//   - 机械臂被锁定
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置左右手校准值
//	err := dobot.SetLRHandCalibrateValue(1.5)
//	if err != nil {
//	    log.Printf("设置左右手校准值失败: %v", err)
//	} else {
//	    log.Printf("左右手校准值设置成功")
//	}
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
// @Summary 获取机械臂左右手模式的校准参数
// @Description 获取机械臂在左右手模式切换时的校准参数。通过查看当前的校准值，
// @Description 可以了解机械臂左右手切换的精度补偿状态。这个值通常在出厂时
// @Description 标定设置。
//
// @Return float32 "左右手校准值：
//   - 用于补偿左右手切换时的位置偏差
//   - 单位：毫米（mm）
//   - 范围：通常在±5mm以内
//     注意：校准值过大可能表示机械结构异常"
//
// @Return error "错误信息"
// @Success 200 {number} float32 "返回左右手校准值"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	// 获取左右手校准值
//	value, err := dobot.GetLRHandCalibrateValue()
//	if err != nil {
//	    log.Printf("获取左右手校准值失败: %v", err)
//	} else {
//	    log.Printf("当前左右手校准值: %.2fmm", value)
//	}
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
	reader := resp.Reader()
	var lrHandCalibrateValue float32
	binary.Read(reader, binary.LittleEndian, &lrHandCalibrateValue)
	return lrHandCalibrateValue, nil
}

// SetQueuedCmdStartExec 开始执行指令队列
// @Summary 开始执行机械臂指令队列
// @Description 开始执行之前加入队列的所有指令。指令队列具有以下特点：
// @Description 1. 队列中的指令按照加入顺序依次执行
// @Description 2. 每条指令执行完成后才会执行下一条指令
// @Description 3. 可以随时暂停、继续或停止队列执行
// @Description 4. 队列执行过程中可以继续添加新指令
// @Description 注意：开始执行前确保队列中的指令是正确的，且机械臂处于使能状态
//
// @Return error "错误信息"
// @Success 200 {string} "执行成功"
// @Failure 400 {error} "执行失败，可能的错误：
//   - 队列为空
//   - 机械臂未使能
//   - 机械臂处于报警状态
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 开始执行指令队列
//	err := dobot.SetQueuedCmdStartExec()
//	if err != nil {
//	    log.Printf("开始执行指令队列失败: %v", err)
//	} else {
//	    log.Printf("指令队列开始执行")
//	}
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
// @Summary 停止执行机械臂指令队列
// @Description 停止执行当前正在执行的指令队列。停止后，机械臂将完成当前
// @Description 正在执行的指令，然后暂停执行队列中的后续指令。可以通过
// @Description SetQueuedCmdStartExec 重新开始执行队列。
//
// @Return error "错误信息"
// @Success 200 {string} "停止成功"
// @Failure 400 {error} "停止失败，可能的错误：
//   - 队列未在执行
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 停止执行指令队列
//	err := dobot.SetQueuedCmdStopExec()
//	if err != nil {
//	    log.Printf("停止执行指令队列失败: %v", err)
//	} else {
//	    log.Printf("指令队列已停止执行")
//	}
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
// @Summary 强制停止执行机械臂指令队列
// @Description 立即强制停止执行指令队列。与普通停止不同，强制停止会立即
// @Description 中断当前正在执行的指令，机械臂可能会在当前位置突然停止。
// @Description 仅在紧急情况下使用此功能。
//
// @Return error "错误信息"
// @Success 200 {string} "强制停止成功"
// @Failure 400 {error} "强制停止失败，可能的错误：
//   - 队列未在执行
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 强制停止执行指令队列
//	err := dobot.SetQueuedCmdForceStopExec()
//	if err != nil {
//	    log.Printf("强制停止指令队列失败: %v", err)
//	} else {
//	    log.Printf("指令队列已强制停止")
//	}
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
// @Summary 开始下载机械臂指令队列
// @Description 通知机械臂开始下载指令队列。通过指定循环次数和每循环的
// @Description 指令数量，可以实现指令队列的循环执行功能。这对于需要重复
// @Description 执行的动作序列非常有用。
//
// @Param totalLoop uint32 true "总循环次数：
//   - 指定指令队列需要循环执行的次数
//   - 0表示无限循环
//     注意：设置合理的循环次数以避免过度运行"
//
// @Param linePerLoop uint32 true "每循环指令行数：
//   - 指定每次循环中要执行的指令数量
//   - 必须小于等于队列中的总指令数
//     注意：确保数量与实际队列指令数匹配"
//
// @Return error "错误信息"
// @Success 200 {string} "下载开始成功"
// @Failure 400 {error} "下载失败，可能的错误：
//   - 参数无效
//   - 队列为空
//   - 机械臂未就绪
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置指令队列循环执行3次，每次执行5条指令
//	err := dobot.SetQueuedCmdStartDownload(3, 5)
//	if err != nil {
//	    log.Printf("开始下载指令队列失败: %v", err)
//	} else {
//	    log.Printf("指令队列开始下载")
//	}
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
// @Summary 停止下载机械臂指令队列
// @Description 通知机械臂停止下载指令队列。当需要中断正在进行的队列
// @Description 下载过程时使用此功能。停止下载后，已下载的指令仍然保留
// @Description 在队列中。
//
// @Return error "错误信息"
// @Success 200 {string} "停止下载成功"
// @Failure 400 {error} "停止失败，可能的错误：
//   - 当前没有下载任务
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 停止下载指令队列
//	err := dobot.SetQueuedCmdStopDownload()
//	if err != nil {
//	    log.Printf("停止下载指令队列失败: %v", err)
//	} else {
//	    log.Printf("指令队列下载已停止")
//	}
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
// @Summary 清除机械臂指令队列
// @Description 清空机械臂指令队列中的所有指令。此操作会删除队列中所有
// @Description 待执行的指令，包括已下载但尚未执行的指令。清除后，需要
// @Description 重新添加指令才能执行新的动作序列。
//
// @Return error "错误信息"
// @Success 200 {string} "清除成功"
// @Failure 400 {error} "清除失败，可能的错误：
//   - 队列正在执行
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 清除指令队列
//	err := dobot.SetQueuedCmdClear()
//	if err != nil {
//	    log.Printf("清除指令队列失败: %v", err)
//	} else {
//	    log.Printf("指令队列已清除")
//	}
func (dobot *Dobot) SetQueuedCmdClear() error {
	message := &Message{
		Id:       ProtocolQueuedCmdClear,
		RW:       true,
		IsQueued: false,
	}
	_, err := dobot.connector.SendMessage(context.Background(), message)
	return err
}

// GetQueuedCmdLeftSpace 获取队列剩余空间
// @Summary 获取机械臂指令队列的剩余空间
// @Description 获取机械臂指令队列中还可以添加的指令数量。通过检查剩余
// @Description 空间，可以避免队列溢出，确保所有指令都能被正确添加到队列中。
//
// @Return uint32 "剩余空间：
//   - 返回队列中还可以添加的指令数量
//   - 0表示队列已满，无法添加新指令"
//
// @Return error "错误信息"
// @Success 200 {number} uint32 "返回剩余空间数量"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	// 获取队列剩余空间
//	space, err := dobot.GetQueuedCmdLeftSpace()
//	if err != nil {
//	    log.Printf("获取队列剩余空间失败: %v", err)
//	} else {
//	    log.Printf("当前队列剩余空间: %d", space)
//	}
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
	var leftSpace uint32
	binary.Read(resp.Reader(), binary.LittleEndian, &leftSpace)
	return leftSpace, nil
}

// GetQueuedCmdCurrentIndex 获取当前队列命令索引
// @Summary 获取机械臂当前执行的指令索引
// @Description 获取机械臂当前正在执行的指令队列索引。通过此索引可以
// @Description 追踪指令的执行进度，了解当前执行到队列中的哪条指令。
// @Description 对于调试和同步控制非常有用。
//
// @Return uint64 "当前指令索引：
//   - 返回当前正在执行的指令在队列中的索引
//   - 0表示队列尚未开始执行或已执行完毕"
//
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回当前指令索引"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 队列未在执行
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	// 获取当前执行的指令索引
//	index, err := dobot.GetQueuedCmdCurrentIndex()
//	if err != nil {
//	    log.Printf("获取当前指令索引失败: %v", err)
//	} else {
//	    log.Printf("当前执行的指令索引: %d", index)
//	}
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
	var currentIndex uint64
	binary.Read(resp.Reader(), binary.LittleEndian, &currentIndex)
	return currentIndex, nil
}

// GetQueuedCmdMotionFinish 获取队列命令运动是否完成
// @Summary 获取机械臂指令队列的执行状态
// @Description 检查当前指令队列是否已完成执行。此函数可用于同步控制，
// @Description 在需要等待当前动作完成后再执行下一步操作时特别有用。
// @Description 建议定期轮询此状态以监控执行进度。
//
// @Return bool "执行完成状态：
//   - true: 队列中的所有指令已执行完成
//   - false: 队列仍在执行中"
//
// @Return error "错误信息"
// @Success 200 {boolean} bool "返回执行完成状态"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 队列未启动
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	// 检查指令队列执行状态
//	finished, err := dobot.GetQueuedCmdMotionFinish()
//	if err != nil {
//	    log.Printf("获取执行状态失败: %v", err)
//	} else {
//	    if finished {
//	        log.Printf("指令队列已执行完成")
//	    } else {
//	        log.Printf("指令队列仍在执行中")
//	    }
//	}
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
	var motionFinish bool
	binary.Read(resp.Reader(), binary.LittleEndian, &motionFinish)
	return motionFinish, nil
}

// SetPTPPOCmd 设置PTP并行输出命令
// @Summary 设置机械臂PTP运动和并行输出命令
// @Description 通过PTP运动命令和并行输出命令数组同时控制机械臂运动和IO输出。
// @Description 这种组合命令可以实现运动过程中的同步IO控制，适用于需要在
// @Description 特定位置触发外部设备的应用场景。
//
// @Param ptpCmd *PTPCmd true "PTP运动命令参数：
//   - 包含目标位置、速度等运动参数
//   - PTPMode: PTP运动模式（0-4）
//   - X,Y,Z: 目标位置坐标（单位：mm）
//   - R: 末端旋转角度（单位：度）
//     注意：确保运动参数在机械臂工作空间范围内"
//
// @Param parallelCmd []ParallelOutputCmd true "并行输出命令数组：
//   - Address: IO地址（0-255）
//   - Level: 输出电平（0:低，1:高）
//   - Time: 延时时间（单位：ms）
//     注意：时间设置会影响运动执行效率"
//
// @Return uint64 "指令队列索引：
//   - 用于跟踪命令执行状态
//   - 可通过此索引查询命令是否完成"
//
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数为空或无效
//   - 目标位置超出范围
//   - IO地址无效
//   - 机械臂被锁定
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置PTP运动并在运动过程中控制IO输出
//	ptpCmd := &PTPCmd{
//	    PTPMode: PTPMode_MovJ,  // 关节运动模式
//	    X: 200, Y: 0, Z: 50, R: 0,
//	}
//	parallelCmd := []ParallelOutputCmd{
//	    {Address: 1, Level: 1, Time: 500},  // 延时500ms后输出高电平
//	}
//	index, err := dobot.SetPTPPOCmd(ptpCmd, parallelCmd)
//	if err != nil {
//	    log.Printf("设置PTP并行输出命令失败: %v", err)
//	} else {
//	    log.Printf("PTP并行输出命令设置成功，指令索引: %d", index)
//	}
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

	reader := resp.Reader()
	binary.Read(reader, binary.LittleEndian, &queuedCmdIndex)
	return queuedCmdIndex, nil
}

// SetPTPPOWithLCmd 设置带L轴的PTP并行输出命令
// @Summary 设置带L轴的机械臂PTP运动和并行输出命令
// @Description 通过带L轴的PTP运动命令和并行输出命令数组同时控制机械臂运动
// @Description 和IO输出。这种组合命令除了基本的PTP运动和IO控制外，还包含
// @Description L轴（第五轴）的控制参数，适用于需要协调L轴运动的应用场景。
//
// @Param ptpWithLCmd *PTPWithLCmd true "带L轴的PTP运动命令参数：
//   - 包含目标位置、速度等运动参数
//   - 包含L轴的位置参数
//   - 详细参数说明参见PTPWithLCmd结构体定义
//     注意：确保所有轴的运动参数在安全范围内"
//
// @Param parallelCmd []ParallelOutputCmd true "并行输出命令数组：
//   - 可以包含多个IO输出控制命令
//   - 这些命令将在运动过程中按设定时序执行
//     注意：命令数量不要过多，以免影响运动性能"
//
// @Return uint64 "指令队列索引"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 参数为空
//   - 运动参数无效
//   - L轴参数无效
//   - 输出命令无效
//   - 机械臂被锁定
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置带L轴的PTP运动并在运动过程中控制IO输出
//	ptpCmd := &PTPWithLCmd{
//	    PTPMode: 0,
//	    X: 200, Y: 0, Z: 50, R: 0,
//	    L: 90,  // L轴旋转到90度位置
//	}
//	parallelCmd := []ParallelOutputCmd{
//	    {Address: 1, Level: 1, Time: 500},  // 延时500ms后输出高电平
//	}
//	index, err := dobot.SetPTPPOWithLCmd(ptpCmd, parallelCmd)
//	if err != nil {
//	    log.Printf("设置带L轴的PTP并行输出命令失败: %v", err)
//	} else {
//	    log.Printf("带L轴的PTP并行输出命令设置成功，指令索引: %d", index)
//	}
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

	reader := resp.Reader()
	binary.Read(reader, binary.LittleEndian, &queuedCmdIndex)
	return queuedCmdIndex, nil
}

// SetWIFIConfigMode 设置WIFI配置模式
// @Summary 设置机械臂WIFI配置模式
// @Description 设置机械臂的WIFI配置模式。启用配置模式后，可以对机械臂的
// @Description WIFI连接参数进行设置，如SSID、密码等。完成配置后应及时
// @Description 关闭配置模式以确保安全。
//
// @Param enable bool true "是否启用WIFI配置模式：
//   - true: 启用配置模式
//   - false: 关闭配置模式
//     注意：配置完成后记得关闭配置模式"
//
// @Return error "错误信息"
// @Success 200 {string} "设置成功"
// @Failure 400 {error} "设置失败，可能的错误：
//   - WIFI模块未就绪
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 启用WIFI配置模式
//	err := dobot.SetWIFIConfigMode(true)
//	if err != nil {
//	    log.Printf("设置WIFI配置模式失败: %v", err)
//	} else {
//	    log.Printf("WIFI配置模式已启用")
//	}
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
// @Summary 获取机械臂WIFI配置模式的状态
// @Description 获取当前机械臂WIFI配置模式的启用状态。通过此函数可以
// @Description 确认当前是否处于WIFI配置模式，以便进行相应的配置操作
// @Description 或检查配置模式是否已正确关闭。
//
// @Return bool "配置模式状态：
//   - true: 当前处于配置模式
//   - false: 当前不在配置模式"
//
// @Return error "错误信息"
// @Success 200 {boolean} bool "返回配置模式状态"
// @Failure 400 {error} "获取失败，可能的错误：
//   - WIFI模块未就绪
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	// 获取WIFI配置模式状态
//	enabled, err := dobot.GetWIFIConfigMode()
//	if err != nil {
//	    log.Printf("获取WIFI配置模式状态失败: %v", err)
//	} else {
//	    if enabled {
//	        log.Printf("当前处于WIFI配置模式")
//	    } else {
//	        log.Printf("当前不在WIFI配置模式")
//	    }
//	}
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
	reader := resp.Reader()
	var mode bool
	binary.Read(reader, binary.LittleEndian, &mode)
	return mode, nil
}

// SetWIFISSID 设置WIFI SSID
// @Summary 设置机械臂WIFI网络的SSID
// @Description 设置机械臂要连接的WIFI网络的SSID。此设置需要在WIFI配置
// @Description 模式下进行。设置SSID后，还需要设置对应的密码才能完成
// @Description WIFI网络的配置。
//
// @Param ssid string true "WIFI网络SSID：
//   - 不能为空
//   - 长度通常不超过32字符
//   - 支持ASCII字符
//     注意：确保SSID正确，否则无法连接网络"
//
// @Return error "错误信息"
// @Success 200 {string} "设置成功"
// @Failure 400 {error} "设置失败，可能的错误：
//   - SSID为空
//   - SSID格式无效
//   - 未处于配置模式
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置WIFI网络SSID
//	err := dobot.SetWIFISSID("MyNetwork")
//	if err != nil {
//	    log.Printf("设置WIFI SSID失败: %v", err)
//	} else {
//	    log.Printf("WIFI SSID设置成功")
//	}
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
// @Failure 400 {error} "获取WIFI SSID失败时返回错误信息"
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
// @Summary 设置机械臂WIFI网络的密码
// @Description 设置机械臂要连接的WIFI网络的密码。此设置需要在WIFI配置
// @Description 模式下进行，且应该在设置SSID之后进行。设置完成后，机械臂
// @Description 将尝试使用这些凭据连接到指定的WIFI网络。
//
// @Param password string true "WIFI网络密码：
//   - 不能为空
//   - 长度通常在8-63字符之间
//   - 支持ASCII字符
//     注意：密码将以加密方式存储"
//
// @Return error "错误信息"
// @Success 200 {string} "设置成功"
// @Failure 400 {error} "设置失败，可能的错误：
//   - 密码为空
//   - 密码格式无效
//   - 未处于配置模式
//   - 未先设置SSID
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置WIFI网络密码
//	err := dobot.SetWIFIPassword("MyPassword123")
//	if err != nil {
//	    log.Printf("设置WIFI密码失败: %v", err)
//	} else {
//	    log.Printf("WIFI密码设置成功")
//	}
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

// GetPTPTime 获取PTP运动时间
// @Summary 获取PTP运动的预计执行时间
// @Description 获取PTP（点到点）运动指令的预计执行时间。此功能可用于
// @Description 预估运动时间，帮助进行运动规划和时序控制。时间计算基于
// @Description 当前的速度和加速度参数设置。
//
// @Return float32 "预计执行时间（单位：秒）"
// @Return error "错误信息"
// @Success 200 {number} float32 "返回预计执行时间"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	// 获取PTP运动时间
//	time, err := dobot.GetPTPTime()
//	if err != nil {
//	    log.Printf("获取PTP运动时间失败: %v", err)
//	} else {
//	    log.Printf("预计执行时间: %.2f秒", time)
//	}
func (dobot *Dobot) GetPTPTime() (float32, error) {
	message := &Message{
		Id:       ProtocolPTPTime,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}
	reader := resp.Reader()
	var time float32
	if err := binary.Read(reader, binary.LittleEndian, &time); err != nil {
		return 0, fmt.Errorf("failed to read PTP time: %v", err)
	}
	return time, nil
}

// FirmwareMode 固件模式
type FirmwareMode uint8

const (
	FirmwareModeNormal    FirmwareMode = 0 // 正常模式
	FirmwareModeUpgrade   FirmwareMode = 1 // 升级模式
	FirmwareModeCalibrate FirmwareMode = 2 // 校准模式
)

// GetFirmwareMode 获取固件模式
// @Summary 获取机械臂固件的当前工作模式
// @Description 获取机械臂固件的当前工作模式。固件可能处于正常工作模式、
// @Description 固件升级模式或校准模式。不同模式下支持的功能不同。
//
// @Return FirmwareMode "固件模式：
//   - 0: 正常模式
//   - 1: 升级模式
//   - 2: 校准模式"
//
// @Return error "错误信息"
// @Success 200 {number} FirmwareMode "返回固件模式"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	// 获取固件模式
//	mode, err := dobot.GetFirmwareMode()
//	if err != nil {
//	    log.Printf("获取固件模式失败: %v", err)
//	    return
//	}
//	switch mode {
//	case FirmwareModeNormal:
//	    log.Printf("固件处于正常工作模式")
//	case FirmwareModeUpgrade:
//	    log.Printf("固件处于升级模式")
//	case FirmwareModeCalibrate:
//	    log.Printf("固件处于校准模式")
//	default:
//	    log.Printf("未知的固件模式: %d", mode)
//	}
func (dobot *Dobot) GetFirmwareMode() (FirmwareMode, error) {
	message := &Message{
		Id:       ProtocolFirmwareMode,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}
	return FirmwareMode(resp.Data()[0]), nil
}

// SetLostStepParams 设置丢步参数
// @Summary 设置机械臂丢步检测参数
// @Description 设置机械臂运动过程中丢步检测的阈值参数。丢步检测用于
// @Description 监控机械臂运动精度，当检测到丢步超过阈值时，可以及时
// @Description 进行校正，确保运动精度。
//
// @Param threshold float32 true "丢步检测阈值：
//   - 单位：度（°）
//   - 范围：通常在0.1-1.0之间
//     注意：阈值过小可能导致误报，过大可能无法及时发现问题"
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
//   - 阈值参数无效
//   - 机械臂被锁定
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 设置丢步检测阈值为0.5度
//	index, err := dobot.SetLostStepParams(0.5, true)
//	if err != nil {
//	    log.Printf("设置丢步参数失败: %v", err)
//	} else {
//	    log.Printf("丢步参数设置成功，指令索引: %d", index)
//	}
func (dobot *Dobot) SetLostStepParams(threshold float32, isQueued bool) (uint64, error) {
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

	reader := resp.Reader()
	var queuedCmdIndex uint64
	if err := binary.Read(reader, binary.LittleEndian, &queuedCmdIndex); err != nil {
		return 0, fmt.Errorf("failed to read queued cmd index: %v", err)
	}
	return queuedCmdIndex, nil
}

// SetLostStepCmd 设置丢步命令
// @Summary 发送机械臂丢步检测命令
// @Description 发送丢步检测命令，触发机械臂执行丢步检测和校正操作。
// @Description 当检测到丢步时，机械臂会自动进行位置校正，以确保运动
// @Description 精度。建议在精度要求高的场合定期执行此命令。
//
// @Param isQueued bool true "是否加入指令队列：
//   - true: 将指令加入队列，按顺序执行
//   - false: 立即执行该指令
//   - 建议使用队列模式以确保检测时序"
//
// @Return uint64 "指令队列索引（当isQueued为true时有效）"
// @Return error "错误信息"
// @Success 200 {number} uint64 "返回指令队列索引"
// @Failure 400 {error} "执行失败，可能的错误：
//   - 机械臂被锁定
//   - 机械臂处于报警状态
//   - 通信错误
//   - 设备未连接"
//
// @Example
//
//	// 执行丢步检测
//	index, err := dobot.SetLostStepCmd(true)
//	if err != nil {
//	    log.Printf("执行丢步检测失败: %v", err)
//	} else {
//	    log.Printf("丢步检测命令已发送，指令索引: %d", index)
//	}
func (dobot *Dobot) SetLostStepCmd(isQueued bool) (uint64, error) {
	message := &Message{
		Id:       ProtocolLostStepDetect,
		RW:       true,
		IsQueued: isQueued,
	}

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return 0, err
	}

	reader := resp.Reader()
	var queuedCmdIndex uint64
	if err := binary.Read(reader, binary.LittleEndian, &queuedCmdIndex); err != nil {
		return 0, fmt.Errorf("failed to read queued cmd index: %v", err)
	}
	return queuedCmdIndex, nil
}

// GetJOGCmd 获取JOG运动指令状态
// @Summary 获取当前JOG运动指令的状态
// @Description 获取机械臂当前JOG（点动）模式下的运动状态信息。通过此函数
// @Description 可以了解当前的运动模式、运动轴和运动方向等信息。这对于监控
// @Description 和同步控制非常有用。
//
// @Return *JOGCmd "JOG运动状态：
//   - isJoint: 当前运动模式（true为关节模式/false为坐标模式）
//   - index: 当前运动轴索引
//   - direction: 当前运动方向
//     注意：返回的是最后一次设置的状态"
//
// @Return error "错误信息"
// @Success 200 {object} *JOGCmd "返回JOG运动状态结构体"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效"
//
// @Example
//
//	// 获取当前JOG运动状态
//	cmd, err := dobot.GetJOGCmd()
//	if err != nil {
//	    log.Printf("获取JOG运动状态失败: %v", err)
//	    return
//	}
//
//	// 输出状态信息
//	log.Printf("当前JOG运动状态：")
//	if cmd.IsJoint {
//	    log.Printf("  模式: 关节运动")
//	    log.Printf("  关节索引: %d", cmd.Index)
//	} else {
//	    log.Printf("  模式: 坐标运动")
//	    axes := []string{"X", "Y", "Z", "R"}
//	    log.Printf("  运动轴: %s", axes[cmd.Index])
//	}
//
//	switch cmd.Direction {
//	case 1:
//	    log.Printf("  方向: 正向")
//	case -1:
//	    log.Printf("  方向: 负向")
//	default:
//	    log.Printf("  状态: 停止")
//	}
func (dobot *Dobot) GetJOGCmd() (*JOGCmd, error) {
	message := &Message{
		Id:       ProtocolJOGCmd,
		RW:       false,
		IsQueued: false,
	}
	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return nil, err
	}

	cmd := &JOGCmd{}
	reader := resp.Reader()
	if err := binary.Read(reader, binary.LittleEndian, cmd); err != nil {
		return nil, fmt.Errorf("failed to read JOG cmd: %v", err)
	}
	return cmd, nil
}

// GetUART4PeripheralsType 获取UART4外设类型
// @Summary 获取UART4外设类型
// @Description 获取当前连接到UART4端口的外设类型
// @Return (uint8, error) "返回外设类型和可能的错误"
func (dobot *Dobot) GetUART4PeripheralsType(ctx context.Context) (uint8, error) {
	message := &Message{
		Id:       ProtocolCheckUART4PeripheralsModel,
		RW:       false,
		IsQueued: false,
		Params:   []byte{},
	}

	response, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	return response.Params[0], nil
}

// SetUART4PeripheralsEnable 设置UART4外设使能状态
// @Summary 设置UART4外设使能状态
// @Description 启用或禁用UART4外设
// @Param isEnable bool "是否启用UART4外设"
// @Return error "可能的错误信息"
func (dobot *Dobot) SetUART4PeripheralsEnable(ctx context.Context, isEnable bool) error {
	message := &Message{
		Id:       ProtocolUART4PeripheralsEnabled,
		RW:       true,
		IsQueued: false,
		Params:   []byte{boolToByte(isEnable)},
	}

	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// GetUART4PeripheralsEnable 获取UART4外设使能状态
// @Summary 获取UART4外设使能状态
// @Description 获取UART4外设当前是否启用
// @Return (bool, error) "返回是否启用和可能的错误"
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

	return response.Params[0] != 0, nil
}

// PluseCmd 脉冲控制命令
type PluseCmd struct {
	ControlMode uint8   // 控制模式
	Port        uint8   // 端口号
	Speed       float32 // 速度（脉冲个数每秒）
	Distance    float32 // 移动距离（脉冲个数）
}

// SendPluse 发送脉冲控制命令
// @Summary 发送脉冲控制命令
// @Description 发送脉冲控制命令，控制电机运动
// @Param pluseCmd PluseCmd "脉冲控制命令参数"
// @Param isQueued bool "是否加入队列"
// @Return (uint64, error) "返回队列索引（如果加入队列）和可能的错误"
func (dobot *Dobot) SendPluse(ctx context.Context, pluseCmd *PluseCmd, isQueued bool) (uint64, error) {
	params := make([]byte, 10) // 2个uint8 + 2个float32
	params[0] = pluseCmd.ControlMode
	params[1] = pluseCmd.Port
	binary.LittleEndian.PutUint32(params[2:6], math.Float32bits(pluseCmd.Speed))
	binary.LittleEndian.PutUint32(params[6:10], math.Float32bits(pluseCmd.Distance))

	message := &Message{
		Id:       ProtocolFunctionPulseMode,
		RW:       true,
		IsQueued: isQueued,
		Params:   params,
	}

	response, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}

	if isQueued {
		reader := response.Reader()
		var queuedCmdIndex uint64
		if err := binary.Read(reader, binary.LittleEndian, &queuedCmdIndex); err != nil {
			return 0, fmt.Errorf("failed to read queued cmd index: %v", err)
		}
		return queuedCmdIndex, nil
	}
	return 0, nil
}

// SendPluseEx 发送脉冲控制命令并等待执行完成
// @Summary 发送脉冲控制命令并等待执行完成
// @Description 发送脉冲控制命令，并等待命令执行完成
// @Param pluseCmd PluseCmd "脉冲控制命令参数"
// @Return error "可能的错误信息"
func (dobot *Dobot) SendPluseEx(ctx context.Context, pluseCmd *PluseCmd) error {
	index, err := dobot.SendPluse(ctx, pluseCmd, true)
	if err != nil {
		return err
	}

	for {
		curIndex, err := dobot.GetQueuedCmdCurrentIndex()
		if err != nil {
			return err
		}
		if curIndex >= index {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
	return nil
}

func boolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
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
// @Summary 设置WIFI IP地址
// @Description 设置机械臂的WIFI IP地址配置
// @Param wifiIPAddress WIFIIPAddress "IP地址配置"
// @Return error "可能的错误信息"
func (dobot *Dobot) SetWIFIIPAddress(ctx context.Context, wifiIPAddress *WIFIIPAddress) error {
	params := make([]byte, 5) // 1字节DHCP + 4字节IP地址
	if wifiIPAddress.DHCP {
		params[0] = 1
	}
	copy(params[1:], wifiIPAddress.Addr[:])

	message := &Message{
		Id:       ProtocolWIFIIPAddress,
		RW:       true,
		IsQueued: false,
		Params:   params,
	}

	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// GetWIFIIPAddress 获取WIFI IP地址
// @Summary 获取WIFI IP地址
// @Description 获取机械臂当前的WIFI IP地址配置
// @Return (*WIFIIPAddress, error) "返回IP地址配置和可能的错误"
func (dobot *Dobot) GetWIFIIPAddress(ctx context.Context) (*WIFIIPAddress, error) {
	message := &Message{
		Id:       ProtocolWIFIIPAddress,
		RW:       false,
		IsQueued: false,
		Params:   []byte{},
	}

	response, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	result := &WIFIIPAddress{}
	reader := response.Reader()
	if err := binary.Read(reader, binary.LittleEndian, result); err != nil {
		return nil, fmt.Errorf("failed to read WIFI IP address: %v", err)
	}
	return result, nil
}

// SetWIFINetmask 设置WIFI子网掩码
// @Summary 设置WIFI子网掩码
// @Description 设置机械臂的WIFI子网掩码
// @Param wifiNetmask WIFINetmask "子网掩码配置"
// @Return error "可能的错误信息"
func (dobot *Dobot) SetWIFINetmask(ctx context.Context, wifiNetmask *WIFINetmask) error {
	message := &Message{
		Id:       ProtocolWIFINetmask,
		RW:       true,
		IsQueued: false,
		Params:   wifiNetmask.Addr[:],
	}

	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// GetWIFINetmask 获取WIFI子网掩码
// @Summary 获取WIFI子网掩码
// @Description 获取机械臂当前的WIFI子网掩码
// @Return (*WIFINetmask, error) "返回子网掩码配置和可能的错误"
func (dobot *Dobot) GetWIFINetmask(ctx context.Context) (*WIFINetmask, error) {
	message := &Message{
		Id:       ProtocolWIFINetmask,
		RW:       false,
		IsQueued: false,
		Params:   []byte{},
	}

	response, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	result := &WIFINetmask{}
	reader := response.Reader()
	if err := binary.Read(reader, binary.LittleEndian, result); err != nil {
		return nil, fmt.Errorf("failed to read WIFI netmask: %v", err)
	}
	return result, nil
}

// SetWIFIGateway 设置WIFI网关
// @Summary 设置WIFI网关
// @Description 设置机械臂的WIFI网关地址
// @Param wifiGateway WIFIGateway "网关配置"
// @Return error "可能的错误信息"
func (dobot *Dobot) SetWIFIGateway(ctx context.Context, wifiGateway *WIFIGateway) error {
	message := &Message{
		Id:       ProtocolWIFIGateway,
		RW:       true,
		IsQueued: false,
		Params:   wifiGateway.Addr[:],
	}

	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// GetWIFIGateway 获取WIFI网关
// @Summary 获取WIFI网关
// @Description 获取机械臂当前的WIFI网关地址
// @Return (*WIFIGateway, error) "返回网关配置和可能的错误"
func (dobot *Dobot) GetWIFIGateway(ctx context.Context) (*WIFIGateway, error) {
	message := &Message{
		Id:       ProtocolWIFIGateway,
		RW:       false,
		IsQueued: false,
		Params:   []byte{},
	}

	response, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	result := &WIFIGateway{}
	reader := response.Reader()
	if err := binary.Read(reader, binary.LittleEndian, result); err != nil {
		return nil, fmt.Errorf("failed to read WIFI gateway: %v", err)
	}
	return result, nil
}

// SetWIFIDNS 设置WIFI DNS
// @Summary 设置WIFI DNS
// @Description 设置机械臂的WIFI DNS服务器地址
// @Param wifiDNS WIFIDNS "DNS配置"
// @Return error "可能的错误信息"
func (dobot *Dobot) SetWIFIDNS(ctx context.Context, wifiDNS *WIFIDNS) error {
	message := &Message{
		Id:       ProtocolWIFIDNS,
		RW:       true,
		IsQueued: false,
		Params:   wifiDNS.Addr[:],
	}

	_, err := dobot.connector.SendMessage(ctx, message)
	return err
}

// GetWIFIDNS 获取WIFI DNS
// @Summary 获取WIFI DNS
// @Description 获取机械臂当前的WIFI DNS服务器地址
// @Return (*WIFIDNS, error) "返回DNS配置和可能的错误"
func (dobot *Dobot) GetWIFIDNS(ctx context.Context) (*WIFIDNS, error) {
	message := &Message{
		Id:       ProtocolWIFIDNS,
		RW:       false,
		IsQueued: false,
		Params:   []byte{},
	}

	response, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	result := &WIFIDNS{}
	reader := response.Reader()
	if err := binary.Read(reader, binary.LittleEndian, result); err != nil {
		return nil, fmt.Errorf("failed to read WIFI DNS: %v", err)
	}
	return result, nil
}

// GetWIFIConnectStatus 获取WIFI连接状态
// @Summary 获取WIFI连接状态
// @Description 获取机械臂当前的WIFI连接状态
// @Return (bool, error) "返回是否已连接和可能的错误"
func (dobot *Dobot) GetWIFIConnectStatus(ctx context.Context) (bool, error) {
	message := &Message{
		Id:       ProtocolWIFIConnectStatus,
		RW:       false,
		IsQueued: false,
		Params:   []byte{},
	}

	response, err := dobot.connector.SendMessage(ctx, message)
	if err != nil {
		return false, err
	}
	reader := response.Reader()
	var status bool
	if err := binary.Read(reader, binary.LittleEndian, &status); err != nil {
		return false, fmt.Errorf("failed to read WIFI connect status: %v", err)
	}
	return status, nil
}

// GetIOMultiplexing 获取IO复用状态
// @Summary 获取IO引脚的复用功能配置
// @Description 获取指定IO引脚的复用功能配置。IO引脚可以配置为不同的功能，
// @Description 如普通IO、PWM输出、ADC输入等。
//
// @Param address "IO引脚地址"
// @Return IOFunction "IO复用功能：
//   - IOFunctionDummy: 未配置
//   - IOFunctionDO: 数字输出
//   - IOFunctionPWM: PWM输出
//   - IOFunctionDI: 数字输入
//   - IOFunctionADC: 模拟输入"
//
// @Return error "错误信息"
// @Success 200 {number} IOFunction "返回IO复用功能"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效
//   - 参数无效"
//
// @Example
//
//	// 获取IO引脚的复用功能
//	function, err := dobot.GetIOMultiplexing(1)
//	if err != nil {
//	    log.Printf("获取IO复用功能失败: %v", err)
//	    return
//	}
//	switch function {
//	case IOFunctionDO:
//	    log.Printf("IO引脚配置为数字输出")
//	case IOFunctionPWM:
//	    log.Printf("IO引脚配置为PWM输出")
//	case IOFunctionDI:
//	    log.Printf("IO引脚配置为数字输入")
//	case IOFunctionADC:
//	    log.Printf("IO引脚配置为模拟输入")
//	default:
//	    log.Printf("IO引脚未配置")
//	}
func (dobot *Dobot) GetIOMultiplexing(ioMultiplexing *IOMultiplexing) (*IOMultiplexing, error) {
	message := &Message{
		Id:       ProtocolIOMultiplexing,
		RW:       false,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, ioMultiplexing)
	message.Params = writer.Bytes()

	resp, err := dobot.connector.SendMessage(context.Background(), message)
	if err != nil {
		return nil, err
	}
	reader := resp.Reader()
	var multiplexing IOMultiplexing
	if err := binary.Read(reader, binary.LittleEndian, &multiplexing); err != nil {
		return nil, fmt.Errorf("failed to read IOMultiplexing: %v", err)
	}
	return &multiplexing, nil
}

// GetIODO 获取IO数字输出状态
// @Summary 获取IO引脚的数字输出状态
// @Description 获取指定IO引脚的数字输出状态。在使用此函数前，
// @Description 需要先将IO引脚配置为数字输出模式。
//
// @Param address "IO引脚地址"
// @Return bool "输出状态：true表示高电平，false表示低电平"
// @Return error "错误信息"
// @Success 200 {bool} "返回输出状态"
// @Failure 400 {error} "获取失败，可能的错误：
//   - 通信错误
//   - 设备未连接
//   - 响应数据无效
//   - 参数无效
//   - IO引脚未配置为数字输出"
//
// @Example
//
//	// 获取IO引脚的数字输出状态
//	level, err := dobot.GetIODO(1)
//	if err != nil {
//	    log.Printf("获取IO数字输出状态失败: %v", err)
//	    return
//	}
//	if level {
//	    log.Printf("IO引脚输出高电平")
//	} else {
//	    log.Printf("IO引脚输出低电平")
//	}
func (dobot *Dobot) GetIODO(ioDO *IODO) (*IODO, error) {
	message := &Message{
		Id:       ProtocolIODO,
		RW:       false,
		IsQueued: false,
	}
	writer := &bytes.Buffer{}
	binary.Write(writer, binary.LittleEndian, ioDO)
	message.Params = writer.Bytes()

	resp, err := dobot.connector.SendMessage(context.Background(), message)
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
