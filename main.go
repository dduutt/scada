package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dduutt/scada/meter"
)

func main() {
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	meter.InitLogger(logFile)
	p := "设备配置表.xlsx"
	groups, err := meter.MeterAddrGroupFromExcel(p)
	if err != nil {
		log.Fatalln("读取配置文件失败:", err)
	}
	c := make(chan *meter.MeterGroupResult)
	defer close(c)
	go meter.ReadMeterGroup(groups, c)

	for r := range c {
		if r.Error != nil {
			log.Printf("读取表:%s失败:%v\n", r.Meter.Code, r.Error)
		} else {
			fmt.Printf("读取表:%s成功:%v\n", r.Meter.Code, r.Meter.Value)
		}
	}
}
