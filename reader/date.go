package reader

import "fmt"

type Date struct {
	year   int // since 1900
	month  int // 1 -> 12
	day    int // 1 -> 31
	hour   int // 0 -> 23
	minute int // 0 -> 59
}

var (
	months = [12]string{"Jan", "Feb", "Mar", "Apr", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
)

// Returns the date formatted as "15 Jan 1980"
func (d *Date) FormatDate() string {
	return fmt.Sprintf("%d %s %d", d.day, months[d.month-1], d.year+1900)
}

// Returns the time formatted as "16:30"
func (d *Date) FormatTime() string {
	hour, minute := fmt.Sprint(d.hour), fmt.Sprint(d.minute)
	if d.hour < 10 {
		hour = "0" + hour
	}

	if d.minute < 10 {
		minute = "0" + minute
	}

	return fmt.Sprintf("%s:%s", hour, minute)
}
