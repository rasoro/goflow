package envs

import (
	_ "embed"
	"encoding/json"
	"sort"
	"strings"
)

// file containing day and month translations, generated using https://github.com/nyaruka/go-locales
//
// ./localesdump --merge days=LC_TIME.day short_days=LC_TIME.abday months=LC_TIME.mon short_months=LC_TIME.abmon am_pm=LC_TIME.am_pm > dates.json
//
//go:embed i18n/i18n.json
var i18nJSON []byte

type Translation struct {
	Days        []string `json:"days"`
	ShortDays   []string `json:"short_days"`
	Months      []string `json:"months"`
	ShortMonths []string `json:"short_months"`
	AmPm        []string `json:"am_pm"`
}

var translations map[string]*Translation
var backdowns = map[string]*Translation{} // language only backdowns for locales that have countries
var defaultLang = "en"

func init() {
	err := json.Unmarshal(i18nJSON, &translations)
	if err != nil {
		panic(err)
	}

	// so that we can iterate translations deterministically (code a-z)
	codes := make([]string, len(translations))
	for c := range translations {
		codes = append(codes, c)
	}
	sort.Strings(codes)

	for _, code := range codes {
		if len(code) == 5 {
			lang := code[:2]
			if backdowns[lang] == nil {
				backdowns[lang] = translations[code] // using first is arbitary but best we can do
			}
		}
	}
}

// GetTranslation gets the best match translation for the given locale
func GetTranslation(locale Locale) *Translation {
	code := strings.Replace(locale.ToISO639_2(), "-", "_", 1)

	// try extract xx_YY match
	t := translations[code]
	if t != nil {
		return t
	}

	// try match by language xx only
	lang := code[:2]
	t = translations[lang]
	if t != nil {
		return t
	}

	// use backdown for this language
	t = backdowns[lang]
	if t != nil {
		return t
	}

	// use default
	return translations[defaultLang]
}
