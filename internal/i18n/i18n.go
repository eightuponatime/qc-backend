package i18n

import (
	"encoding/json"
	"os"
)

type Translations map[string]string

// possible langs: ru, kk
func Load(lang string) (Translations, error) {
	data, err := os.ReadFile("i18n/" + lang + ".json")

	if err != nil {
		return nil, err
	}

	var t Translations
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, err
	}

	return t, nil
}
