package envs_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/goflow/envs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var laTZ, _ = time.LoadLocation("America/Los_Angeles")

func TestDateTimeFromString(t *testing.T) {

	testCases := []struct {
		DateFormat envs.DateFormat
		TimeFormat envs.TimeFormat
		Timezone   string
		FillTime   bool
		Value      string
		Expected   string
		Error      bool
	}{
		// valid cases, varying formats
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "01-02-2001", "01-02-2001 00:00:00 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "date is 01.02.2001 yes", "01-02-2001 00:00:00 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "date is 1-2-99 yes", "01-02-1999 00:00:00 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "01/02/2001", "01-02-2001 00:00:00 +0000 UTC", false},

		// must be strict iso to match despite format
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "2001-01-02T10:34:56Z", "02-01-2001 10:34:56 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "2001-01-02T10:34:56+02:00", "02-01-2001 10:34:56 +0200 MST", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "2001-01-02", "02-01-2001 00:00:00 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", true, "2001-01-02", "02-01-2001 13:36:30.123456789 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "America/Los_Angeles", true, "2001-01-02", "02-01-2001 06:36:30.123456789 -0800 PST", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, " 2001-01-02 ", "02-01-2001 00:00:00 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "on 2001-01-02 ", "", true},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "2001_01_02", "", true},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "2001-1-2", "", true},

		// month first
		{envs.DateFormatMonthDayYear, envs.TimeFormatHourMinute, "UTC", false, "01-02-2001", "02-01-2001 00:00:00 +0000 UTC", false},
		{envs.DateFormatMonthDayYear, envs.TimeFormatHourMinute, "UTC", false, "2001-01-02", "02-01-2001 00:00:00 +0000 UTC", false},
		{envs.DateFormatMonthDayYear, envs.TimeFormatHourMinute, "UTC", false, "2001-1-2", "", true},

		// year first
		{envs.DateFormatYearMonthDay, envs.TimeFormatHourMinute, "UTC", false, "2001-02-01", "01-02-2001 00:00:00 +0000 UTC", false},
		{envs.DateFormatYearMonthDay, envs.TimeFormatHourMinute, "UTC", false, "99-02-01", "01-02-1999 00:00:00 +0000 UTC", false},

		// specific timezone
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "America/Los_Angeles", false, "01\\02\\2001", "01-02-2001 00:00:00 -0800 PST", false},

		// with time filling
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", true, "01-02-2001", "01-02-2001 13:36:30.123456789 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", true, "01-02-2001 04:23", "01-02-2001 04:23:00 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "America/Los_Angeles", true, "01-02-2001", "01-02-2001 06:36:30.123456789 -0800 PST", false},

		// illegal day
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "33-01-2001", "01-01-0001 00:00:00 +0000 UTC", true},

		// illegal month
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "01-13-2001", "01-01-0001 00:00:00 +0000 UTC", true},

		// valid two digit cases
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "01-01-99", "01-01-1999 00:00:00 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "01-01-16", "01-01-2016 00:00:00 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "01-01-16a", "", true},

		// iso dates
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "2016-05-01T18:30:15-08:00", "01-05-2016 18:30:15 -0800 PST", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "2016-05-01T18:30:15Z", "01-05-2016 18:30:15 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "2016-05-01T18:30:15.250Z", "01-05-2016 18:30:15.250 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "1977-06-23T08:34:00.000-07:00", "23-06-1977 15:34:00.000 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "1977-06-23T08:34:00.000250-07:00", "23-06-1977 15:34:00.000250 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "1977-06-23T08:34:00.000250500-07:00", "23-06-1977 15:34:00.000250500 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "2017-06-10T17:34-06:00", "10-06-2017 23:34:00.000000 +0000 UTC", false},

		// with time
		{envs.DateFormatYearMonthDay, envs.TimeFormatHourMinute, "UTC", false, "2001-02-01 03:15", "01-02-2001 03:15:00 +0000 UTC", false},
		{envs.DateFormatYearMonthDay, envs.TimeFormatHourMinute, "UTC", false, "2001-02-01 03:15pm", "01-02-2001 15:15:00 +0000 UTC", false},
		{envs.DateFormatYearMonthDay, envs.TimeFormatHourMinute, "UTC", false, "2001-02-01 03:15 AM", "01-02-2001 03:15:00 +0000 UTC", false},
		{envs.DateFormatYearMonthDay, envs.TimeFormatHourMinute, "UTC", false, "2001-02-01 03:15:34", "01-02-2001 03:15:34 +0000 UTC", false},
		{envs.DateFormatYearMonthDay, envs.TimeFormatHourMinute, "UTC", false, "2001-02-01 03:15:34.123", "01-02-2001 03:15:34.123 +0000 UTC", false},
		{envs.DateFormatYearMonthDay, envs.TimeFormatHourMinute, "UTC", false, "2001-02-01 03:15:34.123456", "01-02-2001 03:15:34.123456 +0000 UTC", false},
	}

	dates.SetNowSource(dates.NewFixedNowSource(time.Date(2018, 9, 13, 13, 36, 30, 123456789, time.UTC)))
	defer dates.SetNowSource(dates.DefaultNowSource)

	for _, tc := range testCases {
		timezone, err := time.LoadLocation(tc.Timezone)
		require.NoError(t, err)

		env := envs.NewBuilder().WithDateFormat(tc.DateFormat).WithTimeFormat(tc.TimeFormat).WithTimezone(timezone).Build()

		value, err := envs.DateTimeFromString(env, tc.Value, tc.FillTime)

		if tc.Error {
			assert.Error(t, err)
		} else {
			require.NoError(t, err, "error parsing date %s", tc.Value)

			expected, err := time.Parse("02-01-2006 15:04:05.999999999 -0700 MST", tc.Expected)
			require.NoError(t, err, "error parsing expected date %s", tc.Expected)

			if !expected.Equal(value) {
				assert.Fail(t, "", "mismatch for date input %s, expected %s, got %s", tc.Value, expected, value)
			}
		}
	}
}

func TestDateFromString(t *testing.T) {
	testCases := []struct {
		dateFormat envs.DateFormat
		value      string
		expected   dates.Date
		hasError   bool
	}{
		{envs.DateFormatDayMonthYear, "it's 31-12-2018 ok", dates.NewDate(2018, 12, 31), false},
		{envs.DateFormatDayMonthYear, "it's 31-12-18 ok", dates.NewDate(2018, 12, 31), false},
		{envs.DateFormatMonthDayYear, "it's 12-31-2018 ok", dates.NewDate(2018, 12, 31), false},
		{envs.DateFormatMonthDayYear, "it's 12-31-18 ok", dates.NewDate(2018, 12, 31), false},
		{envs.DateFormatYearMonthDay, "it's 2018-12-31 ok", dates.NewDate(2018, 12, 31), false},
		{envs.DateFormatYearMonthDay, "it's 18-12-31 ok", dates.NewDate(2018, 12, 31), false},

		// valid ISO always accepted
		{envs.DateFormatDayMonthYear, "2018-12-31", dates.NewDate(2018, 12, 31), false},
		{envs.DateFormatDayMonthYear, "2018-12-31T18:30:15-08:00", dates.NewDate(2018, 12, 31), false},

		{envs.DateFormatDayMonthYear, "it's ok", dates.ZeroDate, true},
		{envs.DateFormatDayMonthYear, "it's 2018-13-01 ok", dates.ZeroDate, true},
		{envs.DateFormatDayMonthYear, "it's 2018-12-32 ok", dates.ZeroDate, true},
	}

	for _, tc := range testCases {
		env := envs.NewBuilder().WithDateFormat(tc.dateFormat).Build()
		parsed, err := envs.DateFromString(env, tc.value)

		if tc.hasError {
			assert.Error(t, err)
		} else {
			require.NoError(t, err, "error parsing time %s", tc.value)

			if !tc.expected.Equal(parsed) {
				assert.Fail(t, "", "mismatch for time input %s, expected %s, got %s", tc.value, tc.expected, parsed)
			}
		}
	}
}

func TestTimeFromString(t *testing.T) {
	testCases := []struct {
		value    string
		expected dates.TimeOfDay
		hasError bool
	}{
		{"it's 10 ok", dates.NewTimeOfDay(10, 0, 0, 0), false},
		{"it's 10 PM ok", dates.NewTimeOfDay(22, 0, 0, 0), false},
		{"it's 10:30 ok", dates.NewTimeOfDay(10, 30, 0, 0), false},
		{"it's 10:30pm ok", dates.NewTimeOfDay(22, 30, 0, 0), false},
		{"it's 10:30 pm ok", dates.NewTimeOfDay(22, 30, 0, 0), false},
		{"it's 1030 ok", dates.NewTimeOfDay(10, 30, 0, 0), false},
		{"it's 1030 PM ok", dates.NewTimeOfDay(22, 30, 0, 0), false},
		{"it's 10:30:45 ok", dates.NewTimeOfDay(10, 30, 45, 0), false},
		{"it's 10:30:45 pm ok", dates.NewTimeOfDay(22, 30, 45, 0), false},
		{"it's 10:30:45.123 ok", dates.NewTimeOfDay(10, 30, 45, 123000000), false},
		{"it's 10:30:45.123 pm ok", dates.NewTimeOfDay(22, 30, 45, 123000000), false},
		{"it's 10:30:45.123456 ok", dates.NewTimeOfDay(10, 30, 45, 123456000), false},
		{"it's 10:30:45.123456 pm ok", dates.NewTimeOfDay(22, 30, 45, 123456000), false},
		{"it's 10:30:45.123456789 ok", dates.NewTimeOfDay(10, 30, 45, 123456789), false},
		{"it's 10:30:45.123456789 pm ok", dates.NewTimeOfDay(22, 30, 45, 123456789), false},

		// 12 am and 12 pm
		{"it's 12:00 AM", dates.NewTimeOfDay(0, 0, 0, 0), false},
		{"it's 12pm", dates.NewTimeOfDay(12, 0, 0, 0), false},

		// fractional component can be any length
		{"it's 10:30:45.123456789123456789 ok", dates.NewTimeOfDay(10, 30, 45, 123456789), false},
		{"it's 10:30:45.1 ok", dates.NewTimeOfDay(10, 30, 45, 100000000), false},

		// 24 can be used to mean midnight
		{"it's 24:00 ok", dates.NewTimeOfDay(0, 0, 0, 0), false},
		{"it's 24:00:00 ok", dates.NewTimeOfDay(0, 0, 0, 0), false},

		{"it's ok", dates.ZeroTimeOfDay, true},
		{"it's 25:30", dates.ZeroTimeOfDay, true},
		{"it's 10:61", dates.ZeroTimeOfDay, true},
		{"it's 10:30:61", dates.ZeroTimeOfDay, true},
	}

	for _, tc := range testCases {
		parsed, err := envs.TimeFromString(tc.value)

		if tc.hasError {
			assert.Error(t, err)
		} else {
			require.NoError(t, err, "error parsing time %s", tc.value)

			if !tc.expected.Equal(parsed) {
				assert.Fail(t, "", "mismatch for time input %s, expected %s, got %s", tc.value, tc.expected, parsed)
			}
		}
	}
}

func TestFormatDatetime(t *testing.T) {
	mst, err := time.LoadLocation("MST")
	require.NoError(t, err)

	d1 := time.Date(2006, 1, 2, 15, 4, 5, 123456789, mst)
	d2 := time.Date(1998, 4, 18, 9, 45, 30, 123456789, time.UTC)
	en_US := envs.NewLocale("eng", "US")
	// es_EC = envs.NewLocale("spa", "EC")

	tests := []struct {
		value    time.Time
		layout   string
		locale   envs.Locale
		expected string
		err      string
	}{
		{d1, "MM-DD-YYYY", en_US, "01-02-2006", ""},
		{d1, "M-D-YY", en_US, "1-2-06", ""},
		{d1, "h:m", en_US, "3:4", ""},
		{d1, "h:m:s aa", en_US, "3:4:5 pm", ""},
		{d1, "h:m:s AA", en_US, "3:4:5 PM", ""},
		{d1, "tt:mm:ss", en_US, "15:04:05", ""},
		{d2, "tt:mm:ss", en_US, "09:45:30", ""},
		{d2, "t:mm:ss", en_US, "9:45:30", ""},
		{d1, "YYYY-MM-DDTtt:mm:ssZZZ", en_US, "2006-01-02T15:04:05-07:00", ""},
		{d1, "YYYY-MM-DDTtt:mm:ssZZZ", en_US, "2006-01-02T15:04:05-07:00", ""},
		{d1, "YYYY-MM-DDThh:mm:ss.fffZZZ", en_US, "2006-01-02T03:04:05.123-07:00", ""},
		{d1, "YYYY-MM-DDThh:mm:ss.fffZ", en_US, "2006-01-02T03:04:05.123-07:00", ""},
		{d1, "YY-M-D", en_US, "06-1-2", ""},
		{d1, "YYYY-MM-DD", en_US, "2006-01-02", ""},
		{d1, "YYYY-MMM-DD", en_US, "2006-Jan-02", ""},
		{d1, "YYYY-MMMM-DD", en_US, "2006-January-02", ""},
		{d1, "EEE", en_US, "Mon", ""},
		{d1, "EEEE", en_US, "Monday", ""},
		{d1, "//YY--MM::DD..", en_US, "//06--01::02..", ""},

		{d1, "tt:mm:ss.fff", en_US, "15:04:05.123", ""},
		{d1, "tt:mm:ss.ffffff", en_US, "15:04:05.123456", ""},
		{d1, "tt:mm:ss.fffffffff", en_US, "15:04:05.123456789", ""},

		{d1, "YYY-MM-DD", en_US, "", "'YYY' is not a valid format sequence"},
		{d1, "YYYY-MMMMM-DD", en_US, "", "'MMMMM' is not a valid format sequence"},
		{d1, "EE", en_US, "", "'EE' is not a valid format sequence"},
		{d1, "tt:mm:ss.ffff", en_US, "", "'ffff' is not a valid format sequence"},
		{d1, "tt:mmm:ss.ffff", en_US, "", "'mmm' is not a valid format sequence"},
		{d1, "tt:mm:sss", en_US, "", "'sss' is not a valid format sequence"},
		{d1, "tt:mm:ss a", en_US, "", "'a' is not a valid format sequence"},
		{d1, "tt:mm:ss A", en_US, "", "'A' is not a valid format sequence"},
		{d1, "tt:mm:ssZZZZ", en_US, "", "'ZZZZ' is not a valid format sequence"},
		{d1, "2006-01-02", en_US, "", "'2' is not a valid format sequence"},
	}

	for _, tc := range tests {
		desc := fmt.Sprintf("%s as '%s' in '%s'", tc.value.String(), tc.layout, tc.locale.ToISO639_2())

		actual, err := envs.FormatDateTime(tc.value, tc.layout, tc.locale)
		if tc.err == "" {
			assert.NoError(t, err, "unexpected error for %s", desc)
			assert.Equal(t, tc.expected, actual, "format mismatch for %s", desc)
		} else {
			assert.EqualError(t, err, tc.err, "error mismatch for %s", desc)
			assert.Equal(t, actual, "")
		}
	}
}
