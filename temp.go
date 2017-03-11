package main

import (
	"fmt"
	"strconv"
)

type Temp struct {
	Value int
	Name  string
}

var (
	DefTemp = Temp{6500, "default"}
	minTemp = Temp{1000, "minimal"}
	maxTemp = Temp{10000, "maximal"}
	presets = []Temp{
		{2300, "candle"},
		{2700, "tungsten"},
		{3400, "halogen"},
		{4200, "fluorescent"},
		{5000, "daylight"},
	}
)

func (m *Temp) Set(s string) error {
	for _, v := range presets {
		if v.Name == s {
			*m = v
			return nil
		}
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	if i < minTemp.Value || i > maxTemp.Value {
		return fmt.Errorf("out of range")
	}
	*m = Temp{i, "manual"}
	return nil
}

func (m Temp) String() string {
	return fmt.Sprint(m.Value)
}

func (_ Temp) Usage() string {
	s := fmt.Sprintf("color temperature in range %v..%v or", minTemp, maxTemp)
	for _, v := range presets {
		s += fmt.Sprintf(" %v", v.Name)
	}
	return s
}
