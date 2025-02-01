package godobot

import (
	"net"
	"syscall"
	"time"
)

const (
	BroadcastPort    = 48899
	BroadcastKeyword = "Who is Dobot?"
	LocalPort        = 2046
)

func SearchDobot() ([]*net.UDPAddr, error) {
	// 创建未绑定接口的 UDP 套接字（用于接收响应）
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: LocalPort})
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	// 启用广播权限
	if file, err := conn.File(); file == nil {
		return nil, err
	} else {
		fd := int(file.Fd())
		syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
	}
	// 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	// 遍历每个接口，发送子网广播
	for _, iface := range interfaces {
		// 跳过无效接口（无 IP 或未启用）
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagBroadcast == 0 {
			continue
		}
		// 获取接口的地址列表
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		// 遍历接口地址，找到 IPv4 地址并计算子网广播地址
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok || ipNet.IP.To4() == nil || ipNet.IP.IsLoopback() {
				continue
			}
			// 计算子网广播地址
			broadcastIP := calculateBroadcastIP(ipNet)
			if broadcastIP == nil {
				continue
			}
			// 发送到子网广播地址
			targetAddr := &net.UDPAddr{
				IP:   broadcastIP,
				Port: BroadcastPort,
			}
			message := []byte(BroadcastKeyword)
			conn.WriteToUDP(message, targetAddr)
		}
	}
	// 接收响应（超时 5 秒）
	buffer := make([]byte, 1024)
	var dobotAddrs []*net.UDPAddr
	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			break
		}
		if n == 0 {
			continue
		}
		if string(buffer[:n]) == addr.String() {
			dobotAddrs = append(dobotAddrs, addr)
		}
	}
	return dobotAddrs, nil
}

// 计算子网广播地址（IP | ^NetMask）
func calculateBroadcastIP(ipNet *net.IPNet) net.IP {
	ip := ipNet.IP.To4()
	if ip == nil {
		return nil
	}
	mask := ipNet.Mask
	broadcast := make(net.IP, len(ip))
	for i := range ip {
		broadcast[i] = ip[i] | ^mask[i]
	}
	return broadcast
}
