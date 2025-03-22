package meter

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"sync"

	"github.com/goburrow/modbus"
)

var ModbusTCPHandler = make(map[string]*modbus.TCPClientHandler)

func CloseModbusTCPHandler() {
	for _, v := range ModbusTCPHandler {
		err := v.Close()
		if err != nil {
			fmt.Println("关闭ModbusTCPHandler失败:", err)
		}
	}
}

func GetOrInitModbusTCPHandler(addr string, slaveID int) (*modbus.TCPClientHandler, error) {
	if h, ok := ModbusTCPHandler[addr]; ok {
		h.SlaveId = byte(slaveID)
		return h, nil
	}
	hanlder := modbus.NewTCPClientHandler(addr)
	hanlder.SlaveId = byte(slaveID)
	hanlder.Timeout = 1
	ModbusTCPHandler[addr] = hanlder
	return hanlder, nil
}

type EnergyMeter struct {
	Code        string
	WorkShop    string
	Room        string
	Name        string
	Protocol    string
	IP          string
	Port        int
	SlaveOrArea string
	Address     int
	Len         int
	Type        string
	Value       float64
	Bytes       []byte
}

type MeterGroupResult struct {
	Meter *EnergyMeter
	Error error
}

func ReadMeterGroup(e map[string][]*EnergyMeter, r chan *MeterGroupResult) {
	var wg sync.WaitGroup

	// 启动一个 goroutine 来关闭通道
	go func() {
		wg.Wait()
		close(r)
	}()

	for _, m := range e {
		wg.Add(1)
		go func(s []*EnergyMeter, c chan *MeterGroupResult) {
			defer wg.Done()
			for _, v := range s {
				b, err := v.Read()
				v.Bytes = b
				c <- &MeterGroupResult{
					Meter: v,
					Error: err,
				}
			}
		}(m, r)
	}
}

func (e *EnergyMeter) Read() ([]byte, error) {
	switch e.Protocol {
	case "S7":
		return e.ReadS7()
	case "modbus_tcp":
		return e.ReadModbusTCP()
	}
	return nil, fmt.Errorf("不支持的协议:%s", e.Protocol)
}

func (e *EnergyMeter) ReadS7() ([]byte, error) {
	return nil, nil
}

func (e *EnergyMeter) ReadModbusTCP() ([]byte, error) {
	slave, err := strconv.Atoi(e.SlaveOrArea)
	if err != nil {
		return nil, err
	}
	addr := fmt.Sprintf("%s:%d", e.IP, e.Port)
	h, err := GetOrInitModbusTCPHandler(addr, slave)
	if err != nil {
		return nil, err
	}
	client := modbus.NewClient(h)
	return client.ReadHoldingRegisters(uint16(e.Address), uint16(e.Len))
}

func ReadBytes(b []byte, bigEndian bool, v any) error {
	buf := bytes.NewReader(b)
	var order binary.ByteOrder

	if bigEndian {
		order = binary.BigEndian
	} else {
		order = binary.LittleEndian
	}
	return binary.Read(buf, order, v)
}

func (e *EnergyMeter) ParseBytes() error {
	if e.Bytes == nil {
		return fmt.Errorf("未读取到数据")
	}
	switch e.Type {
	case "uint16":
		var v uint16
		err := ReadBytes(e.Bytes, true, &v)
		if err != nil {
			return err
		}
		e.Value = float64(v)
	case "int16":
		var v int16
		err := ReadBytes(e.Bytes, true, &v)
		if err != nil {
			return err
		}
		e.Value = float64(v)
	case "uint32":
		var v uint32
		err := ReadBytes(e.Bytes, true, &v)
		if err != nil {
			return err
		}
		e.Value = float64(v)
	case "int32":
		var v int32
		err := ReadBytes(e.Bytes, true, &v)
		if err != nil {
			return err
		}
		e.Value = float64(v)
	case "float32":
		var v float32
		err := ReadBytes(e.Bytes, true, &v)
		if err != nil {
			return err
		}
		e.Value = float64(v)
	case "float64":
		var v float64
		err := ReadBytes(e.Bytes, true, &v)
		if err != nil {
			return err
		}
		e.Value = v
	}
	return fmt.Errorf("不支持的数据类型:%s", e.Type)
}
