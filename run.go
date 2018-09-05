package main

import "time"

const (
	inc      = 2
	initTemp = 4500
	dayMins  = 24 * 60 // minutes in a day
)

type Setter interface {
	Set(temp int) error
}

func minutes(t time.Time) int {
	return t.Hour()*60 + t.Minute()
}

func Run(s Setter) error {
	temp, nowMin := initTemp, minutes(time.Now())

	if nowMin > dayMins/2 {
		temp += inc * (dayMins - nowMin)
	} else {
		temp += inc * nowMin
	}

	if err := s.Set(temp); err != nil {
		return err
	}

	tick := time.NewTicker(time.Minute)
	defer tick.Stop()

	for t := range tick.C {
		nowMin := minutes(t)

		if nowMin > dayMins/2 {
			temp -= inc
		} else {
			temp += inc
		}

		if err := s.Set(temp); err != nil {
			return err
		}
	}

	return nil
}
