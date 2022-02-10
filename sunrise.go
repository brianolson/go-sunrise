// transliteration of noaa.gov javascript code from
// http://www.srrb.noaa.gov/highlights/sunrise/sunrise.html
//
// Go translation by Brian Olson github.com/brianolson/go-sunrise
// This code, like its js version, is public domain.
package sunrise

import (
	"math"
	"time"
)

func isLeapYear(year int) bool {
	if year%4 != 0 {
		return false
	}
	if year%400 == 0 {
		return true
	}
	if year%100 == 0 {
		return false
	}
	return true
}

func radToDeg(ang float64) float64 {
	return 180.0 * ang / math.Pi
}
func degToRad(ang float64) float64 {
	return math.Pi * ang / 180.0
}

func calcDayOfYear(mn, dy int, lpyr bool) int {
	k := 2
	if lpyr {
		k = 1
	}
	return int(math.Floor(float64(275*mn)/9)) - k*int(math.Floor(float64(mn+9)/12)) + dy - 30
}

var weekDayNames []string = []string{
	"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday",
}

func calcDayOfWeek(juld float64) string {
	a := int(math.Floor(juld+1.5)) % 7
	return weekDayNames[a]
}

func CalcJD(when time.Time) float64 {
	when = when.UTC()
	return calcJDYMD(when.Year(), int(when.Month()), when.Day())
}

func calcJDYMD(year, month, day int) float64 {
	if month <= 2 {
		year -= 1
		month += 12
	}
	a := math.Floor(float64(year) / 100.0)
	b := 2 - a + math.Floor(a/4)
	return math.Floor(365.25*float64(year+4716)) + math.Floor(30.6001*float64(month+1)) + float64(day) + b - 1524.5
}

func calcTimeJulianCent(jd float64) float64 {
	return (jd - 2451545.0) / 36525.0
}
func calcJDFromJulianCent(t float64) float64 {
	return t*36525.0 + 2451545.0
}

func calcGeomMeanLongSun(t float64) float64 {
	L0 := 280.46646 + t*(36000.76983+0.0003032*t)
	for L0 > 360.0 {
		L0 -= 360.0
	}
	for L0 < 0.0 {
		L0 += 360.0
	}

	return L0 // in degrees
}

func calcGeomMeanAnomalySun(t float64) float64 {
	return 357.52911 + t*(35999.05029-0.0001537*t)
}

func calcEccentricityEarthOrbit(t float64) float64 {
	return 0.016708634 - t*(0.000042037+0.0000001267*t)
}

func calcSunEqOfCenter(t float64) float64 {
	m := calcGeomMeanAnomalySun(t)

	mrad := degToRad(m)
	sinm := math.Sin(mrad)
	var sin2m = math.Sin(mrad + mrad)
	var sin3m = math.Sin(mrad + mrad + mrad)

	return sinm*(1.914602-t*(0.004817+0.000014*t)) + sin2m*(0.019993-0.000101*t) + sin3m*0.000289
}

func calcSunTrueLong(t float64) float64 {
	l0 := calcGeomMeanLongSun(t)
	c := calcSunEqOfCenter(t)

	return l0 + c
}

func calcSunTrueAnomaly(t float64) float64 {
	m := calcGeomMeanAnomalySun(t)
	c := calcSunEqOfCenter(t)

	return m + c
}

func calcSunRadVector(t float64) float64 {
	v := calcSunTrueAnomaly(t)
	e := calcEccentricityEarthOrbit(t)

	return (1.000001018 * (1 - e*e)) / (1 + e*math.Cos(degToRad(v)))
}

func calcSunApparentLong(t float64) float64 {
	o := calcSunTrueLong(t)

	omega := 125.04 - 1934.136*t
	return o - 0.00569 - 0.00478*math.Sin(degToRad(omega)) // in degrees
}

func calcMeanObliquityOfEcliptic(t float64) float64 {
	seconds := 21.448 - t*(46.8150+t*(0.00059-t*(0.001813)))
	return 23.0 + (26.0+(seconds/60.0))/60.0 // in degrees
}

func calcObliquityCorrection(t float64) float64 {
	e0 := calcMeanObliquityOfEcliptic(t)

	omega := 125.04 - 1934.136*t
	return e0 + 0.00256*math.Cos(degToRad(omega)) // in degrees
}

func calcSunRtAscension(t float64) float64 {
	e := calcObliquityCorrection(t)
	lambda := calcSunApparentLong(t)

	tananum := (math.Cos(degToRad(e)) * math.Sin(degToRad(lambda)))
	tanadenom := (math.Cos(degToRad(lambda)))
	return radToDeg(math.Atan2(tananum, tanadenom)) // in degrees
}

func calcSunDeclination(t float64) float64 {
	e := calcObliquityCorrection(t)
	lambda := calcSunApparentLong(t)

	sint := math.Sin(degToRad(e)) * math.Sin(degToRad(lambda))
	return radToDeg(math.Asin(sint)) // in degrees
}

func calcEquationOfTime(t float64) float64 {
	var epsilon = calcObliquityCorrection(t)
	l0 := calcGeomMeanLongSun(t)
	e := calcEccentricityEarthOrbit(t)
	m := calcGeomMeanAnomalySun(t)

	y := math.Tan(degToRad(epsilon) / 2.0)
	y *= y

	sin2l0 := math.Sin(2.0 * degToRad(l0))
	sinm := math.Sin(degToRad(m))
	cos2l0 := math.Cos(2.0 * degToRad(l0))
	sin4l0 := math.Sin(4.0 * degToRad(l0))
	sin2m := math.Sin(2.0 * degToRad(m))

	Etime := y*sin2l0 - 2.0*e*sinm + 4.0*e*y*sinm*cos2l0 - 0.5*y*y*sin4l0 - 1.25*e*e*sin2m

	return radToDeg(Etime) * 4.0 // in minutes of time
}
func calcHourAngleSunrise(lat, solarDec float64) float64 {
	latRad := degToRad(lat)
	sdRad := degToRad(solarDec)

	//HAarg := (math.Cos(degToRad(90.833))/(math.Cos(latRad)*math.Cos(sdRad))-math.Tan(latRad) * math.Tan(sdRad))

	return (math.Acos(math.Cos(degToRad(90.833))/(math.Cos(latRad)*math.Cos(sdRad)) - math.Tan(latRad)*math.Tan(sdRad))) // in radians
}
func calcHourAngleSunset(lat, solarDec float64) float64 {
	latRad := degToRad(lat)
	sdRad := degToRad(solarDec)

	//HAarg := (math.Cos(degToRad(90.833))/(math.Cos(latRad)*math.Cos(sdRad))-math.Tan(latRad) * math.Tan(sdRad))

	HA := (math.Acos(math.Cos(degToRad(90.833))/(math.Cos(latRad)*math.Cos(sdRad)) - math.Tan(latRad)*math.Tan(sdRad)))

	return -HA // in radians
}

func utcMinutesToTimeUTC(utcMinutes float64, basisUtc time.Time) time.Time {
	utcHours := math.Floor(utcMinutes / 60)
	utcMinutes -= utcHours * 60
	tm := math.Floor(utcMinutes)
	utcSeconds := (utcMinutes - tm) * 60
	utcMinutes = tm
	ts := math.Floor(utcSeconds)
	utcNsec := (utcSeconds - ts) * 1000000000
	utcSeconds = ts
	return time.Date(basisUtc.Year(), basisUtc.Month(), basisUtc.Day(), int(utcHours), int(utcMinutes), int(utcSeconds), int(utcNsec), time.UTC)
}

// CalcSunrise returns the sunrise before the solar noon nearest to `when`
// latitude degrees North
// longitude degrees East
func CalcSunrise(when time.Time, latitude, longitude float64) time.Time {
	// library internally coded in degrees WEST
	longitude *= -1
	solnoon := CalcSolNoon(when, longitude*-1)
	JD := CalcJD(solnoon)
	t := calcTimeJulianCent(JD)
	snutc := solnoon.UTC()
	noonmin := float64(snutc.Hour()*60) + float64(snutc.Minute()) + (float64(snutc.Second()) / 60.0) + (float64(snutc.Nanosecond()) / 60000000000.0)
	tnoon := calcTimeJulianCent(JD + noonmin/1440.0)
	utcMinutes := calcSunriseUTCInner(t, tnoon, latitude, longitude)
	oututc := utcMinutesToTimeUTC(utcMinutes, when.UTC())
	return oututc.In(when.Location())
}

// latitude degrees North
// longitude degrees WEST
func calcSunriseUTC(JD, latitude, longitude float64) float64 {
	t := calcTimeJulianCent(JD)

	// *** Find the time of solar noon at the location, and use
	//     that declination. This is better than start of the
	//     Julian day

	noonmin := calcSolNoonUTC(t, longitude)
	tnoon := calcTimeJulianCent(JD + noonmin/1440.0)
	return calcSunriseUTCInner(t, tnoon, latitude, longitude)
}

// latitude degrees North
// longitude degrees WEST
func calcSunriseUTCInner(t, tnoon, latitude, longitude float64) float64 {
	// *** First pass to approximate sunrise (using solar noon)

	eqTime := calcEquationOfTime(tnoon)
	solarDec := calcSunDeclination(tnoon)
	hourAngle := calcHourAngleSunrise(latitude, solarDec)

	delta := longitude - radToDeg(hourAngle)
	timeDiff := 4 * delta              // in minutes of time
	timeUTC := 720 + timeDiff - eqTime // in minutes

	// alert("eqTime = " + eqTime + "\nsolarDec = " + solarDec + "\ntimeUTC = " + timeUTC)

	// *** Second pass includes fractional jday in gamma calc

	newt := calcTimeJulianCent(calcJDFromJulianCent(t) + timeUTC/1440.0)
	eqTime = calcEquationOfTime(newt)
	solarDec = calcSunDeclination(newt)
	hourAngle = calcHourAngleSunrise(latitude, solarDec)
	delta = longitude - radToDeg(hourAngle)
	timeDiff = 4 * delta
	timeUTC = 720 + timeDiff - eqTime // in minutes

	// alert("eqTime = " + eqTime + "\nsolarDec = " + solarDec + "\ntimeUTC = " + timeUTC)

	return timeUTC
}

// CalcSolNoon returns the solar noon nearest the parameter when at the longitude.
// Internally checks yesterday, today, tomorrow.
// This works around edges in "today" converting to UTC and back.
// longitude degrees East
func CalcSolNoon(when time.Time, longitude float64) time.Time {
	// this library internally coded in degrees WEST
	longitude *= -1
	utct := when.UTC()
	jd := CalcJD(utct)
	_, _, out := calcSolNoonUTCInner(utct, jd, longitude)
	return out.In(when.Location())
}

func calcSolNoonUTCInner(utct time.Time, jd, longitude float64) (noonmin, bestjd float64, noonUtc time.Time) {
	bestd := 3600.0 * 24.0 // 1 day of seconds
	for dj := -1; dj <= 1; dj++ {
		t := calcTimeJulianCent(jd + float64(dj))
		utcMinutes := calcSolNoonUTC(t, longitude)
		xu := time.Date(utct.Year(), utct.Month(), utct.Day()+dj, utct.Hour(), utct.Minute(), utct.Second(), utct.Nanosecond(), utct.Location())
		ottu := utcMinutesToTimeUTC(utcMinutes, xu)
		ds := math.Abs(ottu.Sub(utct).Seconds())
		if ds < bestd {
			bestd = ds
			noonUtc = ottu
			noonmin = utcMinutes
			bestjd = jd + float64(dj)
		}
	}
	return
}

// t is from calcTimeJulianCent(jd)
// longitude is degrees west
func calcSolNoonUTC(t, longitude float64) float64 {
	// First pass uses approximate solar noon to calculate eqtime
	tnoon := calcTimeJulianCent(calcJDFromJulianCent(t) + longitude/360.0)
	eqTime := calcEquationOfTime(tnoon)
	solNoonUTC := 720 + (longitude * 4) - eqTime // min

	newt := calcTimeJulianCent(calcJDFromJulianCent(t) - 0.5 + solNoonUTC/1440.0)

	eqTime = calcEquationOfTime(newt)
	// solarNoonDec := calcSunDeclination(newt)
	solNoonUTC = 720 + (longitude * 4) - eqTime // min

	return solNoonUTC
}

// CalcSunset finds the sunset after the solar noon nearest to `when`
// latitude degrees North
// longitude degrees East
func CalcSunset(when time.Time, latitude, longitude float64) time.Time {
	// this library internally coded in degrees WEST
	longitude *= -1
	JD := CalcJD(when)
	t := calcTimeJulianCent(JD)

	// *** Find the time of solar noon at the location, and use
	//     that declination. This is better than start of the
	//     Julian day

	noonmin := calcSolNoonUTC(t, longitude)
	tnoon := calcTimeJulianCent(JD + noonmin/1440.0)
	utcMinutes := calcSunsetUTCInner(t, tnoon, latitude, longitude)
	return utcMinutesToTimeUTC(utcMinutes, when.UTC()).In(when.Location())
}

// latitude degrees North
// longitude degrees WEST
func calcSunsetUTCInner(t, tnoon, latitude, longitude float64) float64 {
	// First calculates sunrise and approx length of day

	eqTime := calcEquationOfTime(tnoon)
	solarDec := calcSunDeclination(tnoon)
	hourAngle := calcHourAngleSunset(latitude, solarDec)

	delta := longitude - radToDeg(hourAngle)
	timeDiff := 4 * delta
	timeUTC := 720 + timeDiff - eqTime

	// first pass used to include fractional day in gamma calc

	newt := calcTimeJulianCent(calcJDFromJulianCent(t) + timeUTC/1440.0)
	eqTime = calcEquationOfTime(newt)
	solarDec = calcSunDeclination(newt)
	hourAngle = calcHourAngleSunset(latitude, solarDec)

	delta = longitude - radToDeg(hourAngle)
	timeDiff = 4 * delta
	timeUTC = 720 + timeDiff - eqTime // in minutes

	return timeUTC
}

// SunsForDay calculates rise, noon, and set for day around solar noon closest to `when`.
// Slightly more efficient than calling 3 functions separately.
// latitude degrees North
// longitude degrees East
func SunsForDay(when time.Time, latitude, longitude float64) (rise, noon, set time.Time) {
	// this library internally coded in degrees WEST
	longitude *= -1
	utct := when.UTC()
	jd := CalcJD(utct)

	noonmin, bestjd, noonUtc := calcSolNoonUTCInner(utct, jd, longitude)
	jd = bestjd
	t := calcTimeJulianCent(jd)
	noon = noonUtc.In(when.Location())
	tnoon := calcTimeJulianCent(jd + noonmin/1440.0)
	utcMinutes := calcSunriseUTCInner(t, tnoon, latitude, longitude)
	rise = utcMinutesToTimeUTC(utcMinutes, when.UTC()).In(when.Location())
	utcMinutes = calcSunsetUTCInner(t, tnoon, latitude, longitude)
	set = utcMinutesToTimeUTC(utcMinutes, when.UTC()).In(when.Location())
	return
}
