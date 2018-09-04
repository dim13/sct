package main

import "time"

const inc = 2

type Setter interface {
	Set(temp int) error
}

func minutes(t time.Time) int {
	return t.Hour()*60 + t.Minute()
}

func Run(s Setter) error {
	temp := 4500
	now := minutes(time.Now())

	if now > 720 {
		temp += inc * (1440 - now)
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

		if now > 720 {
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
