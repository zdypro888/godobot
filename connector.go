package godobot

import (
	"bufio"
	"context"
	"io"
	"net"
	"strings"

	"go.bug.st/serial"
)

const (
	SyncByte       = 0xAA
	MaxPayloadSize = SyncByte - 1 // 确保payload不大于SYNC_BYTE
)

// Packet 负载结构
type Packet struct {
	Id     uint8 // 命令ID
	Ctrl   uint8 // 控制字节
	Params []byte
}

// Message 用于协议通信的消息结构
type Message struct {
	Id       ProtocolId // 原 id
	RW       uint8      // 原 rw
	IsQueued uint8      // 原 isQueued
	Params   []byte
}

func (message *Message) ToPacket() *Packet {
	packet := &Packet{}
	packet.Id = uint8(message.Id)
	packet.Ctrl = 0
	packet.Ctrl |= message.RW & 0x01
	packet.Ctrl |= (message.IsQueued << 1) & 0x02
	packet.Params = make([]byte, len(message.Params))
	copy(packet.Params, message.Params)
	return packet
}

func (packet *Packet) ToMessage() *Message {
	message := &Message{}
	message.Id = ProtocolId(packet.Id)
	message.RW = packet.Ctrl & 0x01
	message.IsQueued = (packet.Ctrl >> 1) & 0x01
	message.Params = make([]byte, len(packet.Params))
	copy(message.Params, packet.Params)
	return message
}

type Connector struct {
	Running         bool
	FirmwareType    string
	FirmwareVersion string
	FirmwareRunTime float64
	port            io.ReadWriteCloser
}

func (connector *Connector) Open(ctx context.Context, name string) error {
	if !strings.HasPrefix(name, "/dev/") {
		udpAddr, err := net.ResolveUDPAddr("udp", name)
		if err != nil {
			return err
		}
		conn, err := net.DialUDP("udp", nil, udpAddr)
		if err != nil {
			return err
		}
		connector.port = conn
	} else {
		mode := &serial.Mode{
			BaudRate: 115200,
			DataBits: 8,
			Parity:   serial.NoParity,
			StopBits: serial.OneStopBit,
		}
		serialPort, err := serial.Open(name, mode)
		if err != nil {
			return err
		}
		connector.port = serialPort
	}
	connector.Running = true
	go connector.receiveGoRoutine(ctx)
	return nil
}

func (connector *Connector) WritePacket(packet *Packet) error {
	var checksum uint8
	checksum += packet.Id
	checksum += packet.Ctrl
	for _, v := range packet.Params {
		checksum += v
	}
	checksum = uint8(0) - checksum
	var err error
	if _, err = connector.port.Write([]byte{byte(SyncByte), byte(SyncByte), byte(len(packet.Params) + 2)}); err != nil {
		return err
	}
	if _, err = connector.port.Write([]byte{byte(packet.Id), byte(packet.Ctrl)}); err != nil {
		return err
	}
	if _, err = connector.port.Write(packet.Params); err != nil {
		return err
	}
	if _, err = connector.port.Write([]byte{byte(checksum)}); err != nil {
		return err
	}
	return nil
}

func (connector *Connector) receiveGoRoutine(ctx context.Context) {
	var err error
	reader := bufio.NewReader(connector.port)
	for connector.Running {
		if ctx.Err() != nil {
			break
		}
		var sbyte byte
		if sbyte, err = reader.ReadByte(); err != nil {
			break
		}
		if sbyte != SyncByte {
			continue
		}
		if sbyte, err = reader.ReadByte(); err != nil {
			break
		}
		if sbyte != SyncByte {
			continue
		}
		// PayloadLen
		if sbyte, err = reader.ReadByte(); err != nil {
			break
		}
		if sbyte >= SyncByte {
			continue
		}
		// id ctrl params + checksum
		var data []byte
		if data, err = reader.Peek(int(sbyte) + 1); err != nil {
			break
		}
		var checksum uint8
		for _, v := range data {
			checksum += v
		}
		if checksum != 0 {
			if _, err = reader.Discard(3); err != nil {
				break
			}
			continue
		}
		packet := &Packet{}
		packet.Id, _ = reader.ReadByte()
		packet.Ctrl, _ = reader.ReadByte()
		packet.Params = make([]byte, sbyte-2)
		reader.Read(packet.Params)
		reader.ReadByte()
		connector.onMessage(packet.ToMessage())
		packet = nil
	}
}

func (connector *Connector) onMessage(message *Message) {

}
