// Package locales holds the translations of the removed site UI. Nothing
// serves them anymore. They are kept as a source for
// wallet/appserver/newui/locales/cmd/importlocales.
package locales

import (
	"fmt"

	"github.com/bisoncraft/caribou/wallet/intl"
	"golang.org/x/text/language"
)

var (
	Locales map[string]map[string]*intl.Translation
)

func init() {
	Locales = map[string]map[string]*intl.Translation{
		"en-US": EnUS,
		"pt-BR": PtBr,
		"zh-CN": ZhCN,
		"pl-PL": PlPL,
		"de-DE": DeDE,
		"ar":    Ar,
	}

	for localeName := range Locales {
		_, err := language.Parse(localeName)
		if err != nil {
			panic(fmt.Sprintf("failed to parse locale name %v: %v", localeName, err))
		}
	}
}
