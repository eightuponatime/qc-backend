package impl

import "time"

func resolveBusinessDateByCutoff(localNow time.Time, cutoffHour int) time.Time {
	if localNow.Hour() < cutoffHour {
		localNow = localNow.AddDate(0, 0, -1)
	}

	return localNow
}

func resolveCurrentPeriodBounds(now time.Time, location *time.Location, cutoffHour int) (time.Time, time.Time) {
	businessDate := resolveBusinessDateByCutoff(now.In(location), cutoffHour)
	year := businessDate.Year()
	month := businessDate.Month()

	if businessDate.Day() <= 15 {
		return time.Date(year, month, 1, 0, 0, 0, 0, location),
			time.Date(year, month, 15, 0, 0, 0, 0, location)
	}

	return time.Date(year, month, 16, 0, 0, 0, 0, location),
		lastDayOfMonth(year, month, location)
}

func resolveLatestReportPeriod(currentTime time.Time, location *time.Location, reportSendHour int) (time.Time, time.Time, bool) {
	localNow := currentTime.In(location)
	year := localNow.Year()
	month := localNow.Month()
	day := localNow.Day()
	hour := localNow.Hour()

	switch {
	case day > 16 || (day == 16 && hour >= reportSendHour):
		return time.Date(year, month, 1, 0, 0, 0, 0, location),
			time.Date(year, month, 15, 0, 0, 0, 0, location),
			true
	case day > 1 || (day == 1 && hour >= reportSendHour):
		previousMonthDate := localNow.AddDate(0, -1, 0)
		prevYear := previousMonthDate.Year()
		prevMonth := previousMonthDate.Month()
		return time.Date(prevYear, prevMonth, 16, 0, 0, 0, 0, location),
			lastDayOfMonth(prevYear, prevMonth, location),
			true
	default:
		previousMonthDate := localNow.AddDate(0, -1, 0)
		prevYear := previousMonthDate.Year()
		prevMonth := previousMonthDate.Month()
		return time.Date(prevYear, prevMonth, 1, 0, 0, 0, 0, location),
			time.Date(prevYear, prevMonth, 15, 0, 0, 0, 0, location),
			true
	}
}

func resolveNextPeriodBounds(periodStart time.Time, location *time.Location) (time.Time, time.Time) {
	periodStart = periodStart.In(location)

	if periodStart.Day() == 1 {
		return time.Date(periodStart.Year(), periodStart.Month(), 16, 0, 0, 0, 0, location),
			lastDayOfMonth(periodStart.Year(), periodStart.Month(), location)
	}

	nextMonthDate := periodStart.AddDate(0, 1, 0)
	return time.Date(nextMonthDate.Year(), nextMonthDate.Month(), 1, 0, 0, 0, 0, location),
		time.Date(nextMonthDate.Year(), nextMonthDate.Month(), 15, 0, 0, 0, 0, location)
}

func lastDayOfMonth(year int, month time.Month, location *time.Location) time.Time {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, location)
}
