package godobot

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.bug.st/serial"
)

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
	if connector.FirmwareType != "" || connector.FirmwareType != "DobotWifi" {
		if _, err := connector.port.Write([]byte("\nM10\nM10\nM10\nM10\nM10\nM13\nM13\nM13\nM13\nM13\n")); err != nil {
			return err
		}
		grblRegex := regexp.MustCompile(`GRBL\:\sV((\d+\.){2}\d+)`)
		marlinRegex := regexp.MustCompile(`MARLIN\:\sV((\d+\.){2}\d+)`)
		marlinTimeRegex := regexp.MustCompile(`Runtim1:(-?\d*\.\d*),`)
		reader := bufio.NewReader(connector.port)
		readbuf := make([]byte, 1024)
		var readText string
		startTime := time.Now()
		for {
			if time.Since(startTime) > 500*time.Millisecond {
				connector.FirmwareType = "DobotSerial"
				connector.FirmwareVersion = "0.0.0"
				break
			}
			n, err := reader.Read(readbuf)
			if err != nil {
				return err
			}
			if n > 0 {
				readText += string(readbuf[:n])
				// 检查是否是GRBL固件
				if matches := grblRegex.FindStringSubmatch(readText); matches != nil {
					connector.FirmwareType = "GRBL"
					connector.FirmwareVersion = matches[1]
				} else if matches := marlinRegex.FindStringSubmatch(readText); matches != nil {
					connector.FirmwareType = "MARLIN"
					connector.FirmwareVersion = matches[1]
					// 检查MARLIN的运行时间
					if timeMatches := marlinTimeRegex.FindStringSubmatch(readText); timeMatches != nil {
						runTimeFloat, _ := strconv.ParseFloat(timeMatches[1], 64)
						connector.FirmwareRunTime = runTimeFloat
					}
				}
			}
			time.Sleep(10 * time.Millisecond)
		}
		connector.port.Close()
		return errors.New("firmware type not supported")
	}
	connector.Running = true
	go connector.receiveGoRoutine(ctx)
	return nil
}

func (connector *Connector) receiveGoRoutine(ctx context.Context) {
	readbuf := make([]byte, 1024)
	reader := bufio.NewReader(connector.port)
	for connector.Running {
		if ctx.Err() != nil {
			break
		}
		readLen, err := reader.Read(readbuf)
		if err != nil {
			break
		}
		fmt.Printf("read text: %s\n", string(readbuf[:readLen]))
	}
}
