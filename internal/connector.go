package internal

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"strings"
	"time"

	"go.bug.st/serial"
)

const (
	SyncByte       = 0xAA
	MaxPayloadSize = SyncByte - 1 // 确保payload不大于SYNC_BYTE
)

var ErrLeftSpace = errors.New("left space is not enough")

// Packet 负载结构
type MessageAck struct {
	Message *Message
	Error   error
}

// Message 用于协议通信的消息结构
type Message struct {
	Id       ProtocolId // 原 id
	RW       bool       // 原 rw
	IsQueued bool       // 原 isQueued
	Params   []byte
	AckLen   uint8
}

func (message *Message) Ctrl() uint8 {
	ctrl := uint8(0)
	if message.RW {
		ctrl |= 0x01
	}
	if message.IsQueued {
		ctrl |= 0x02
	}
	return ctrl
}

func (message *Message) SetCtrl(ctrl uint8) {
	if ctrl&0x01 != 0 {
		message.RW = true
	}
	if (ctrl>>1)&0x01 != 0 {
		message.IsQueued = true
	}
}

func (message *Message) Data() []byte {
	return message.Params[:message.AckLen]
}

func (message *Message) Reader() io.Reader {
	return bytes.NewReader(message.Params)
}

func (message *Message) Read(data any) error {
	return binary.Read(message.Reader(), binary.LittleEndian, data)
}

func (message *Message) Bool() bool {
	return message.Params[0] != 0
}

func (message *Message) Uint16() uint16 {
	return binary.LittleEndian.Uint16(message.Params)
}

func (message *Message) Uint32() uint32 {
	return binary.LittleEndian.Uint32(message.Params)
}

func (message *Message) Uint64() uint64 {
	return binary.LittleEndian.Uint64(message.Params)
}

func (message *Message) Float32() float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(message.Params))
}

func (message *Message) Float64() float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(message.Params))
}

type outMessage struct {
	*Message
	done chan *MessageAck
}

func (outmsg *outMessage) Reply(message *Message) {
	if outmsg.done != nil {
		outmsg.done <- &MessageAck{Message: message, Error: nil}
	}
}

func (outmsg *outMessage) Error(err error) {
	if outmsg.done != nil {
		outmsg.done <- &MessageAck{Message: nil, Error: err}
	}
}

func (outmsg *outMessage) Close() {
	if outmsg.done != nil {
		close(outmsg.done)
		outmsg.done = nil
	}
}

type Connector struct {
	Alarms         []uint8
	Error          error
	port           io.ReadWriteCloser
	recevieError   chan error
	recevieMessage chan *Message
	sendingMessage chan *outMessage
	leftSpace      uint32
}

func (connector *Connector) Open(name string, baudrate uint32) error {
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
		if baudrate == 0 {
			baudrate = 115200
		}
		mode := &serial.Mode{
			BaudRate: int(baudrate),
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
	connector.recevieError = make(chan error)
	connector.recevieMessage = make(chan *Message)
	connector.sendingMessage = make(chan *outMessage)
	connector.leftSpace = 0
	go connector.receiveGoRoutine()
	go connector.processGoRoutine()
	return nil
}

func (connector *Connector) Close() error {
	return connector.port.Close()
}

func (connector *Connector) writeMessage(message *Message) error {
	var checksum uint8
	ctrl := message.Ctrl()
	checksum += uint8(message.Id)
	checksum += ctrl
	for _, v := range message.Params {
		checksum += v
	}
	checksum = uint8(0) - checksum
	var err error
	if _, err = connector.port.Write([]byte{byte(SyncByte), byte(SyncByte), byte(len(message.Params) + 2)}); err != nil {
		return err
	}
	if _, err = connector.port.Write([]byte{byte(message.Id), byte(ctrl)}); err != nil {
		return err
	}
	if _, err = connector.port.Write(message.Params); err != nil {
		return err
	}
	if _, err = connector.port.Write([]byte{byte(checksum)}); err != nil {
		return err
	}
	return nil
}

func (connector *Connector) sendMessage(message *Message) (*Message, error) {
	var err error
	const maxRetries = 3
	for retry := 0; retry < maxRetries; retry++ {
		if err = connector.writeMessage(message); err != nil {
			return nil, err
		}
		select {
		case ack := <-connector.recevieMessage:
			if ack.Id == message.Id {
				if ack.Id == ProtocolQueuedCmdLeftSpace {
					connector.leftSpace = binary.LittleEndian.Uint32(ack.Params)
				}
				return ack, nil
			}
			// 非预期消息，丢弃。TODO: 是否要判断丢弃几个？
		case <-time.After(3 * time.Second):
			// 超时
		}
	}
	return nil, nil
}

func (connector *Connector) receiveGoRoutine() {
	var err error
	reader := bufio.NewReader(connector.port)
	for {
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
		var idbyte, ctrlbyte uint8
		idbyte, _ = reader.ReadByte()
		ctrlbyte, _ = reader.ReadByte()
		message := &Message{}
		message.Id = ProtocolId(idbyte)
		message.SetCtrl(ctrlbyte)
		message.Params = make([]byte, MaxPayloadSize-2)
		message.AckLen = sbyte - 2
		if message.AckLen > 0 {
			reader.Read(message.Params[:message.AckLen])
		}
		reader.ReadByte()
		connector.recevieMessage <- message
	}
	connector.recevieError <- err
}

func (connector *Connector) processGoRoutine() {
	var err error
	breakfor := false
	alarmTicker := time.NewTicker(100 * time.Millisecond)
	defer alarmTicker.Stop()
	for !breakfor {
		select {
		case <-alarmTicker.C:
			alarmStateGet := &Message{Id: ProtocolAlarmsState, RW: false, IsQueued: false}
			var alarmState *Message
			if alarmState, err = connector.sendMessage(alarmStateGet); err != nil {
				breakfor = true
			} else if alarmState != nil {
				connector.Alarms = alarmState.Data()
			}
		case err = <-connector.recevieError:
			breakfor = true
		case outmsg := <-connector.sendingMessage:
			for i, alarm := range connector.Alarms {
				if alarm != 0 {
					outmsg.Error(fmt.Errorf("alarm: %d-%d", i, alarm))
					break
				}
			}
			var ack *Message
			for {
				if !outmsg.IsQueued || connector.leftSpace > 0 {
					// 非队列消息直接发送
					if ack, err = connector.sendMessage(outmsg.Message); err != nil {
						breakfor = true
						outmsg.Error(err)
					} else if ack != nil {
						if outmsg.IsQueued {
							connector.leftSpace--
						}
						outmsg.Reply(ack)
					} else {
						outmsg.Error(errors.New("send message timeout max retries"))
					}
					break
				}
				cmdGetLeftSpace := &Message{Id: ProtocolQueuedCmdLeftSpace, RW: false, IsQueued: false}
				if ack, err = connector.sendMessage(cmdGetLeftSpace); err != nil {
					breakfor = true
					outmsg.Error(err)
					break
				}
				if ack != nil {
					connector.leftSpace = binary.LittleEndian.Uint32(ack.Params)
					if connector.leftSpace == 0 {
						outmsg.Error(ErrLeftSpace)
						break
					}
				}
			}
		}
	}
	connector.Error = err
}

func (connector *Connector) SendMessage(message *Message) (*Message, error) {
	outmsg := &outMessage{Message: message, done: make(chan *MessageAck)}
	defer outmsg.Close()
	connector.sendingMessage <- outmsg
	ack := <-outmsg.done
	return ack.Message, ack.Error
}
