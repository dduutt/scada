package meter

import (
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
)

var H = []string{"编号", "车间", "配电室", "名称", "协议", "IP", "PORT", "从站/区域", "地址", "长度", "类型"}

func FromExcel(file string) ([]*EnergyMeter, error) {
	f, err := excelize.OpenFile(file)

	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows("设备配置表")

	if err != nil {
		return nil, err
	}

	header := rows[0]
	if !compareHeader(header) {
		return nil, fmt.Errorf("表头应为:%v读取到:\n%v", H, header)
	}
	var meters []*EnergyMeter
	for i, row := range rows[1:] {
		l, err := strconv.Atoi(row[7])
		if err != nil {
			return nil, fmt.Errorf("第%d行长度错误:%v", i+2, row[7])
		}
		p, err := strconv.Atoi(row[6])
		if err != nil {
			return nil, fmt.Errorf("第%d行端口错误:%v", i+2, row[6])
		}
		address, err := strconv.Atoi(row[8])
		if err != nil {
			return nil, fmt.Errorf("第%d行地址错误:%v", i+2, row[8])
		}
		m := &EnergyMeter{
			Code:        row[0],
			WorkShop:    row[1],
			Room:        row[2],
			Name:        row[3],
			Protocol:    row[4],
			IP:          row[5],
			Port:        p,
			SlaveOrArea: row[7],
			Address:     address,
			Len:         l,
			Type:        row[10],
			Value:       0,
		}
		meters = append(meters, m)

	}
	return meters, nil

}

func MeterAddrGroupFromExcel(path string) (map[string][]*EnergyMeter, error) {

	meters, err := FromExcel(path)
	if err != nil {
		return nil, err
	}

	group := make(map[string][]*EnergyMeter)
	for _, m := range meters {
		if _, ok := group[m.WorkShop]; !ok {
			group[m.WorkShop] = make([]*EnergyMeter, 0)
		}
		group[m.WorkShop] = append(group[m.WorkShop], m)
	}
	return group, nil
}
func compareHeader(header []string) bool {
	if len(header) != len(H) {
		return false
	}
	for i, v := range header {
		if v != H[i] {
			return false
		}
	}
	return true
}
