package envs_test

import (
	"testing"

	"github.com/nyaruka/goflow/envs"

	"github.com/stretchr/testify/assert"
)

func TestGetTranslation(t *testing.T) {
	tests := []struct {
		locale  envs.Locale
		sun     string
		sunday  string
		jan     string
		january string
		am      string
	}{
		{envs.NewLocale(envs.Language("eng"), envs.Country("US")), "Sun", "Sunday", "Jan", "January", "AM"},
		{envs.NewLocale(envs.Language("eng"), envs.Country("GB")), "Sun", "Sunday", "Jan", "January", "am"},
		{envs.NewLocale(envs.Language("eng"), envs.NilCountry), "Sun", "Sunday", "Jan", "January", "am"},
		{envs.NewLocale(envs.Language("spa"), envs.Country("EC")), "dom", "domingo", "ene", "enero", ""},
		{envs.NewLocale(envs.Language("spa"), envs.NilCountry), "dom", "domingo", "ene", "enero", ""},
		{envs.NewLocale(envs.Language("por"), envs.Country("BR")), "dom", "domingo", "jan", "janeiro", ""},
		{envs.NewLocale(envs.Language("por"), envs.Country("PT")), "dom", "domingo", "jan", "janeiro", ""},
		{envs.NewLocale(envs.Language("por"), envs.NilCountry), "dom", "domingo", "jan", "janeiro", ""},
		{envs.NewLocale(envs.Language("kin"), envs.Country("RW")), "Mwe", "Ku cyumweru", "Mut", "Mutarama", ""},
		{envs.NewLocale(envs.Language("kin"), envs.NilCountry), "Mwe", "Ku cyumweru", "Mut", "Mutarama", ""},
		{envs.NewLocale(envs.Language("zho"), envs.Country("CN")), "日", "星期日", "1月", "一月", "上午"},
		{envs.NewLocale(envs.Language("zho"), envs.Country("HK")), "日", "星期日", "1月", "一月", "上午"},
		{envs.NewLocale(envs.Language("zho"), envs.Country("SG")), "日", "星期日", "一月", "一月", "上午"},
		{envs.NewLocale(envs.Language("zho"), envs.Country("TW")), "日", "週日", " 1月", "一月", "上午"},
		{envs.NewLocale(envs.Language("zho"), envs.NilCountry), "日", "星期日", "1月", "一月", "上午"}, // backs down to first zh translation
	}

	for _, tc := range tests {
		trans := envs.GetTranslation(tc.locale)
		assert.Equal(t, tc.sun, trans.ShortDays[0], "short day mismatch for locale %s-%s", tc.locale.Language, tc.locale.Country)
		assert.Equal(t, tc.sunday, trans.Days[0], "full day mismatch for locale %s-%s", tc.locale.Language, tc.locale.Country)
		assert.Equal(t, tc.jan, trans.ShortMonths[0], "short month mismatch for locale %s-%s", tc.locale.Language, tc.locale.Country)
		assert.Equal(t, tc.january, trans.Months[0], "full month mismatch for locale %s-%s", tc.locale.Language, tc.locale.Country)
		assert.Equal(t, tc.am, trans.AmPm[0], "AM mismatch for locale %s-%s", tc.locale.Language, tc.locale.Country)
	}
}
