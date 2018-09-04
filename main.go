package main

import (
	"flag"
	"log"
)

func main() {
	temp := defTemp
	var run bool
	flag.Var(&temp, "temp", temp.Usage())
	flag.BoolVar(&run, "run", false, "run in background")
	flag.Parse()

	x, err := NewX(Points)
	if err != nil {
		log.Fatal(err)
	}
	defer x.Close()

	if run {
		Run(x)
	} else {
		if err := x.Set(temp.Value); err != nil {
			log.Fatal(err)
		}
	}
}
