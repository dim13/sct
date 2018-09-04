package main

import "time"

const (
	inc     = 2
	midTerm = 4500
	maxMin  = 24 * 60 // minutes in a day
)

type Setter interface {
	Set(temp int) error
}

func minutes(t time.Time) int {
	return t.Hour()*60 + t.Minute()
}

func Run(s Setter) error {
	temp, now := midTerm, minutes(time.Now())

	if now > maxMin/2 {
		temp += inc * (maxMin - now)
	} else {
		temp += inc * now
	}

	if err := s.Set(temp); err != nil {
		return err
	}

	tick := time.NewTicker(time.Minute)
	defer tick.Stop()

	for t := range tick.C {
		now := minutes(t)

		if now > maxMin/2 {
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
