# go-sunrise
Sunrise/Noon/Sunset math in Go

transliteration of noaa.gov javascript code from
http://www.srrb.noaa.gov/highlights/sunrise/sunrise.html

Go translation by Brian Olson github.com/brianolson/go-sunrise

This code, like its js version, is public domain.

```
func CalcJD(when time.Time) float64
func CalcSolNoon(when time.Time, longitude float64) time.Time
func CalcSunrise(when time.Time, latitude, longitude float64) time.Time
func CalcSunset(when time.Time, latitude, longitude float64) time.Time
func SunsForDay(when time.Time, latitude, longitude float64) (rise, noon, set time.Time)
```
