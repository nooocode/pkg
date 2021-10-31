package modbus

import (
	"fmt"
	"sync"

	modbus "github.com/thinkgos/gomodbus/v2"
)

func NewModbusClient(addr string) (*ModbusClient, error) {
	c := &ModbusClient{
		addr: addr,
	}
	err := c.Connect()
	return c, err
}

type ModbusClient struct {
	addr string
	c    modbus.Client
	lock sync.Mutex
}

//Connect 连接modbus
func (client *ModbusClient) Connect() error {
	client.lock.Lock()
	defer client.lock.Unlock()

	if client.c != nil && client.c.IsConnected() {
		return nil
	}

	if client.c != nil {
		client.c.Close()
		client.c = nil
	}

	c := modbus.NewClient(modbus.NewTCPClientProvider(client.addr, modbus.WithEnableLogger(), modbus.WithAutoReconnect(0x06)))
	err := c.Connect()
	if err != nil {
		return err
	}
	client.c = c
	return nil
}

//Close 关闭
func (client *ModbusClient) Close() {
	client.lock.Lock()
	defer client.lock.Unlock()
	if client.c != nil {
		client.c.Close()
		client.c = nil
	}
}

//ReadHoldingRegistersBytes 读取寄存器值
func (client *ModbusClient) ReadHoldingRegistersBytes(slaveID byte, address, quantity uint16) ([]byte, error) {
	// 重连
	err := client.Connect()
	if err != nil {
		client.Close()
		return nil, err
	}

	result, err := client.c.ReadHoldingRegistersBytes(slaveID, address, quantity)
	if err == modbus.ErrClosedConnection {
		client.Close()
		return nil, err
	}

	if err != nil {
		fmt.Println("ReadHoldingRegistersBytes error:", err)
		return nil, err
	}
	return result, nil
}
