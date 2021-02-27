package i18n

import (
	"strings"
	"time"

	"github.com/goodsign/monday"
)

const dateFormat = "Monday 02 January 2006"
const dateTimeFormat = dateFormat + " @ 15:04:05 GMT"

// FormatDate returns the provided formatted and localized date.
func FormatDate(dateTime time.Time, locale string) string {
	return format(dateTime, locale, dateFormat)
}

// FormatDateTime returns the provided formatted and localized date and time.
func FormatDateTime(dateTime time.Time, locale string) string {
	return format(dateTime, locale, dateTimeFormat)
}

func format(dateTime time.Time, locale, format string) string {
	return monday.Format(dateTime, format, monday.Locale(strings.Replace(locale, "-", "_", 1)))
}
