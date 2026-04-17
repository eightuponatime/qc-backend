package i18n

import (
	"net/http"
	"strings"
)

var supportedLangs = map[string]struct{}{
	"ru": {},
	"kk": {},
	"en": {},
}

func DetectLanguage(r *http.Request) string {
	lang := r.URL.Query().Get("lang")
	if isSupported(lang) {
		return lang
	}

	if c, err := r.Cookie("lang"); err == nil {
		if isSupported(c.Value) {
			return c.Value
		}
	}

	header := r.Header.Get("Accept-Language")
	if header != "" {
		langs := strings.Split(header, ",")

		for _, l := range langs {
			code := strings.TrimSpace(strings.Split(l, ";")[0])

			if len(code) >= 2 {
				code = strings.ToLower(code[:2])
				if isSupported(code) {
					return code
				}
			}
		}
	}

	return "ru"
}

func isSupported(lang string) bool {
	_, ok := supportedLangs[lang]
	return ok
}
