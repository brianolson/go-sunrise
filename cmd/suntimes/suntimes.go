package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/brianolson/go-sunrise"
)

func main() {
	var latitude float64
	var longitude float64
	flag.Float64Var(&latitude, "lat", 42.4, "latitude, degrees North")
	flag.Float64Var(&longitude, "lon", -71.1, "longitude, degrees East")
	flag.Parse()

	now := time.Now()

	// functions one at a time
	fmt.Printf("%d-%d-%d\n", now.Year(), int(now.Month()), now.Day())
	jd := sunrise.CalcJD(now)
	fmt.Printf("jd %f\n", jd)
	rise := sunrise.CalcSunrise(now, latitude, longitude)
	fmt.Printf("sunrise:    %s\n", rise)
	solNoon := sunrise.CalcSolNoon(now, longitude)
	fmt.Printf("solar noon: %s\n", solNoon)
	set := sunrise.CalcSunset(now, latitude, longitude)
	fmt.Printf("sunset:     %s\n", set)

	fmt.Println()

	// rise,noon,set in one shot
	rise, solNoon, set = sunrise.SunsForDay(now, latitude, longitude)
	fmt.Printf("sunrise:    %s\n", rise)
	fmt.Printf("solar noon: %s\n", solNoon)
	fmt.Printf("sunset:     %s\n", set)
}
