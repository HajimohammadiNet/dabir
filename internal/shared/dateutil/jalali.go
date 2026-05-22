package dateutil

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var ErrInvalidJalaliDate = errors.New("invalid jalali date")

func ParseOfficialDate(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, errors.New("date is required")
	}

	value = NormalizeDigits(value)
	value = strings.ReplaceAll(value, "-", "/")

	parts := strings.Split(value, "/")
	if len(parts) != 3 {
		return time.Time{}, errors.New("invalid date format, expected YYYY/MM/DD")
	}

	y, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, errors.New("invalid date year")
	}

	m, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, errors.New("invalid date month")
	}

	d, err := strconv.Atoi(parts[2])
	if err != nil {
		return time.Time{}, errors.New("invalid date day")
	}

	// If year is lower than 1700, treat it as Jalali.
	if y < 1700 {
		gy, gm, gd, err := JalaliToGregorian(y, m, d)
		if err != nil {
			return time.Time{}, err
		}

		return time.Date(gy, time.Month(gm), gd, 0, 0, 0, 0, time.UTC), nil
	}

	// Backward compatibility: accept Gregorian too.
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC), nil
}

func ToJalaliString(t time.Time) string {
	jy, jm, jd := GregorianToJalali(t.Year(), int(t.Month()), t.Day())
	return fmt.Sprintf("%04d/%02d/%02d", jy, jm, jd)
}

func JalaliToGregorian(jy, jm, jd int) (int, int, int, error) {
	if jm < 1 || jm > 12 {
		return 0, 0, 0, ErrInvalidJalaliDate
	}

	if jd < 1 || jd > 31 {
		return 0, 0, 0, ErrInvalidJalaliDate
	}

	if jm > 6 && jd > 30 {
		return 0, 0, 0, ErrInvalidJalaliDate
	}

	jy += 1595

	days := -355668 + (365 * jy) + ((jy / 33) * 8) + (((jy % 33) + 3) / 4) + jd

	if jm < 7 {
		days += (jm - 1) * 31
	} else {
		days += ((jm - 7) * 30) + 186
	}

	gy := 400 * (days / 146097)
	days %= 146097

	if days > 36524 {
		gy += 100 * ((days - 1) / 36524)
		days = (days - 1) % 36524

		if days >= 365 {
			days++
		}
	}

	gy += 4 * (days / 1461)
	days %= 1461

	if days > 365 {
		gy += (days - 1) / 365
		days = (days - 1) % 365
	}

	gd := days + 1

	monthDays := []int{
		0,
		31,
		28,
		31,
		30,
		31,
		30,
		31,
		31,
		30,
		31,
		30,
		31,
	}

	if isGregorianLeap(gy) {
		monthDays[2] = 29
	}

	gm := 1
	for gm <= 12 && gd > monthDays[gm] {
		gd -= monthDays[gm]
		gm++
	}

	return gy, gm, gd, nil
}

func GregorianToJalali(gy, gm, gd int) (int, int, int) {
	gdm := []int{0, 31, 59, 90, 120, 151, 181, 212, 243, 273, 304, 334}

	var jy int

	if gy > 1600 {
		jy = 979
		gy -= 1600
	} else {
		jy = 0
		gy -= 621
	}

	gy2 := gy
	if gm > 2 {
		gy2++
	}

	days := (365 * gy) +
		((gy2 + 3) / 4) -
		((gy2 + 99) / 100) +
		((gy2 + 399) / 400) -
		80 +
		gd +
		gdm[gm-1]

	jy += 33 * (days / 12053)
	days %= 12053

	jy += 4 * (days / 1461)
	days %= 1461

	if days > 365 {
		jy += (days - 1) / 365
		days = (days - 1) % 365
	}

	var jm int
	var jd int

	if days < 186 {
		jm = 1 + (days / 31)
		jd = 1 + (days % 31)
	} else {
		jm = 7 + ((days - 186) / 30)
		jd = 1 + ((days - 186) % 30)
	}

	return jy, jm, jd
}

func isGregorianLeap(year int) bool {
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}

func NormalizeDigits(value string) string {
	replacer := strings.NewReplacer(
		"۰", "0",
		"۱", "1",
		"۲", "2",
		"۳", "3",
		"۴", "4",
		"۵", "5",
		"۶", "6",
		"۷", "7",
		"۸", "8",
		"۹", "9",

		"٠", "0",
		"١", "1",
		"٢", "2",
		"٣", "3",
		"٤", "4",
		"٥", "5",
		"٦", "6",
		"٧", "7",
		"٨", "8",
		"٩", "9",
	)

	return replacer.Replace(value)
}
