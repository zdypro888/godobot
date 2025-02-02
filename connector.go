package godobot

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"go.bug.st/serial"
)

const (
	SyncByte       = 0xAA
	MaxPayloadSize = SyncByte - 1 // 确保payload不大于SYNC_BYTE
)

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
	}
}

type Connector struct {
	port         io.ReadWriteCloser
	stopChan     chan error
	errChan      chan error
	messageAck   chan *Message
	messageQueue chan *outMessage

	leftSpace uint32
}

func (connector *Connector) Open(ctx context.Context, name string, baudrate uint32) error {
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
	connector.stopChan = make(chan error, 0x01)
	connector.errChan = make(chan error)
	connector.messageAck = make(chan *Message)
	connector.messageQueue = make(chan *outMessage)
	go connector.receiveGoRoutine(ctx)
	go connector.processGoRoutine(ctx)
	connector.messageQueue <- &outMessage{Message: &Message{Id: ProtocolQueuedCmdLeftSpace}}
	return nil
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

func (connector *Connector) receiveGoRoutine(ctx context.Context) {
	var err error
	reader := bufio.NewReader(connector.port)
	for {
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
		var idbyte, ctrlbyte uint8
		idbyte, _ = reader.ReadByte()
		ctrlbyte, _ = reader.ReadByte()
		message := &Message{}
		message.Id = ProtocolId(idbyte)
		message.SetCtrl(ctrlbyte)
		message.Params = make([]byte, sbyte-2)
		reader.Read(message.Params)
		reader.ReadByte()
		connector.messageAck <- message
	}
	connector.errChan <- err
}

func (connector *Connector) sendMessage(ctx context.Context, message *Message) (*Message, error) {
	var err error
	const maxRetries = 3
	for retry := 0; retry < maxRetries; retry++ {
		log.Printf("发送信息: %d", message.Id)
		if err = connector.writeMessage(message); err != nil {
			return nil, err
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case ack := <-connector.messageAck:
			if ack.Id == message.Id {
				log.Printf("收到信息: %d", ack.Id)
				return ack, nil
			}
			// 非预期消息，丢弃。TODO: 是否要判断丢弃几个？
			log.Printf("丢弃信息: %d", ack.Id)
		case <-time.After(3 * time.Second):
			// 超时
			log.Printf("发送信息超时: %d", message.Id)
		}
	}
	return nil, nil
}

func (connector *Connector) processGoRoutine(ctx context.Context) {
	var err error
	breakfor := false
	for !breakfor {
		select {
		case <-ctx.Done():
			breakfor = false
		case err = <-connector.errChan:
			breakfor = true
		case message := <-connector.messageQueue:
			var ack *Message
			if message.Id != ProtocolQueuedCmdLeftSpace && (!message.IsQueued || connector.leftSpace > 0) {
				// 非队列消息直接发送
				if ack, err = connector.sendMessage(ctx, message.Message); err != nil {
					breakfor = true
					message.Error(err)
				} else if ack != nil {
					message.Reply(ack)
					connector.leftSpace--
				} else {
					message.Error(errors.New("send message timeout max retries"))
				}
			} else {
				var cmdGetLeftSpace *Message
				if message.Id == ProtocolQueuedCmdLeftSpace {
					cmdGetLeftSpace = message.Message
				} else {
					cmdGetLeftSpace = &Message{Id: ProtocolQueuedCmdLeftSpace, RW: false, IsQueued: false}
				}
				if ack, err = connector.sendMessage(ctx, cmdGetLeftSpace); err != nil {
					breakfor = true
					message.Error(err)
				} else if ack != nil {
					numbuf := make([]byte, 4)
					copy(numbuf, ack.Params)
					connector.leftSpace = binary.LittleEndian.Uint32(numbuf)
					if connector.leftSpace == 0 {
						message.Error(errors.New("left space is 0"))
					} else if message.Id == ProtocolQueuedCmdLeftSpace {
						message.Reply(ack)
					} else if ack, err = connector.sendMessage(ctx, message.Message); err != nil {
						breakfor = true
						message.Error(err)
					} else if ack != nil {
						message.Reply(ack)
						connector.leftSpace--
					} else {
						message.Error(errors.New("send message timeout max retries"))
					}
				}
			}
		}
	}
	connector.stopChan <- err
}

func (connector *Connector) SendMessage(ctx context.Context, message *Message) (*Message, error) {
	outmsg := &outMessage{Message: message, done: make(chan *MessageAck)}
	defer outmsg.Close()
	connector.messageQueue <- outmsg
	ack := <-outmsg.done
	return ack.Message, ack.Error
}
