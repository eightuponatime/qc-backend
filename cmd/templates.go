package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"qc/internal/dto"
	"qc/internal/i18n"
)

var tmpl *template.Template

func initTemplates() {
	tmpl = template.New("").Funcs(template.FuncMap{
		"asset":    asset,
		"assetCSS": assetCSS,
		"toJSON": func(v any) template.JS {
			data, err := json.Marshal(v)
			if err != nil {
				return template.JS("{}")
			}
			return template.JS(data)
		},
		"findMeal": func(vote *dto.VoteResponseDto, mealType string) *dto.VoteMealItemResponseDto {
			if vote == nil {
				return nil
			}

			for _, item := range vote.Items {
				if item.MealType == mealType {
					copied := item
					return &copied
				}
			}

			return nil
		},
		"mealLabel": func(mealType string, t i18n.Translations) string {
			return t["meals."+mealType]
		},
		"shiftMealLabel": func(mealType string, shiftType string, t i18n.Translations) string {
			if shiftType == "night" {
				return t["night_meals."+mealType]
			}
			return t["meals."+mealType]
		},
		"stars": func(rating *int16) string {
			if rating == nil || *rating <= 0 {
				return "—"
			}

			result := ""
			for i := 0; i < int(*rating); i++ {
				result += "★"
			}
			return result
		},
		"dict": func(values ...any) (map[string]any, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("dict requires even number of args")
			}

			result := make(map[string]any, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				result[key] = values[i+1]
			}

			return result, nil
		},
		"ratingEq": func(rating *int16, value int) bool {
			if rating == nil {
				return false
			}
			return int(*rating) == value
		},
		"mealRu": func(mealType string) string {
			switch mealType {
			case "breakfast":
				return "завтрак"
			case "lunch":
				return "обед"
			case "dinner":
				return "ужин"
			default:
				return mealType
			}
		},
		"shiftMealRu": func(mealType string, shiftType string) string {
			if shiftType == "night" {
				switch mealType {
				case "breakfast":
					return "первый прием пищи"
				case "lunch":
					return "второй прием пищи"
				case "dinner":
					return "третий прием пищи"
				default:
					return mealType
				}
			}

			switch mealType {
			case "breakfast":
				return "завтрак"
			case "lunch":
				return "обед"
			case "dinner":
				return "ужин"
			default:
				return mealType
			}
		},
		"shiftRu": func(shiftType string) string {
			switch shiftType {
			case "night":
				return "Ночная смена"
			default:
				return "Дневная смена"
			}
		},
		"weekdayRu": func(weekday string) string {
			switch weekday {
			case "Monday":
				return "понедельник"
			case "Tuesday":
				return "вторник"
			case "Wednesday":
				return "среда"
			case "Thursday":
				return "четверг"
			case "Friday":
				return "пятница"
			case "Saturday":
				return "суббота"
			case "Sunday":
				return "воскресенье"
			default:
				return weekday
			}
		},
		"hasText": func(value string) bool {
			return value != ""
		},
	})

	tmpl = template.Must(tmpl.ParseGlob("templates/*.html"))
}
