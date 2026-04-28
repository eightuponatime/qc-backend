package impl

import (
	"fmt"
	"html/template"
	"math"
)

const emailBodyTemplate = `
<!doctype html>
<html lang="ru">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Сводка по качеству питания</title>
  <style>
    @media only screen and (max-width: 640px) {
      .email-shell {
        padding: 16px 10px !important;
      }
      .email-card {
        padding: 18px !important;
        border-radius: 16px !important;
      }
      .email-hero {
        padding: 22px 20px !important;
      }
      .email-title {
        font-size: 24px !important;
      }
      .email-shift-header td,
      .email-weekday-layout td,
      .email-date-layout td {
        display: block !important;
        width: 100% !important;
        box-sizing: border-box !important;
      }
      .email-date-rating {
        padding-top: 10px !important;
        text-align: left !important;
      }
      .email-average-pill {
        min-width: 0 !important;
        padding: 6px 10px !important;
        border-radius: 999px !important;
        text-align: left !important;
      }
      .email-average-label,
      .email-average-value {
        display: inline-block !important;
        vertical-align: middle !important;
      }
      .email-average-label {
        margin-right: 8px !important;
        font-size: 10px !important;
      }
      .email-average-value {
        margin-top: 0 !important;
        font-size: 15px !important;
      }
    }
  </style>
</head>
<body style="margin:0;padding:0;background:#f3f5f7;font-family:Segoe UI,Arial,sans-serif;color:#1f2933;">
  <div class="email-shell" style="max-width:760px;margin:0 auto;padding:24px 16px;">
    <div class="email-card email-hero" style="background:#11324d;color:#ffffff;border-radius:20px;padding:28px 32px;">
      <div style="font-size:13px;letter-spacing:0.08em;text-transform:uppercase;opacity:0.8;">Контроль качества питания</div>
      <h1 class="email-title" style="margin:10px 0 8px;font-size:30px;line-height:1.2;">Сводка за период {{.PeriodStartDisplay}} - {{.PeriodEndDisplay}}</h1>
      <p style="margin:0;font-size:15px;line-height:1.6;opacity:0.9;">Автоматически сформированная управленческая сводка по голосам и отзывам сотрудников.</p>
    </div>

    <div style="height:16px;"></div>

    <div class="email-card" style="background:#e0f2fe;border-radius:20px;padding:24px;border:1px solid #7dd3fc;">
      <h2 style="margin:0 0 12px;font-size:22px;color:#075985;">Подробная аналитика</h2>
      <p style="margin:0 0 14px;color:#0c4a6e;line-height:1.6;">Для просмотра полной аналитики по оценкам и отзывам перейдите на сайт.</p>
      <table role="presentation" width="100%" cellspacing="0" cellpadding="0" style="border-collapse:separate;border-spacing:0 10px;">
        <tr>
          <td style="font-size:13px;color:#0369a1;width:120px;">Ссылка</td>
          <td style="font-size:15px;font-weight:700;"><a href="{{.AnalyticsURL}}" style="color:#075985;text-decoration:none;">{{.AnalyticsURL}}</a></td>
        </tr>
        <tr>
          <td style="font-size:13px;color:#0369a1;width:120px;">Код доступа</td>
          <td>
            <span style="display:inline-block;background:#ffffff;border:1px solid #7dd3fc;border-radius:12px;padding:10px 14px;font-size:20px;font-weight:800;letter-spacing:0.08em;color:#0f172a;">{{.AccessCode}}</span>
          </td>
        </tr>
        <tr>
          <td style="font-size:13px;color:#0369a1;width:120px;">Действует до</td>
          <td style="font-size:15px;font-weight:700;color:#0f172a;">{{.AccessValidUntil}}</td>
        </tr>
      </table>
    </div>

    <div style="height:16px;"></div>

    {{range $shift := .ShiftSummaries}}
    <div class="email-card" style="background:#ffffff;border-radius:20px;padding:24px;border:1px solid #dde4ea;">
      <table class="email-shift-header" role="presentation" width="100%" cellspacing="0" cellpadding="0" style="border-collapse:collapse;margin-bottom:18px;">
        <tr>
          <td style="vertical-align:top;">
            <h2 style="margin:0;font-size:22px;color:#102a43;">{{shiftRu $shift.ShiftType}}</h2>
          </td>
          <td style="vertical-align:top;text-align:right;">
            <div style="display:inline-block;background:#f8fafc;border:1px solid #dde4ea;border-radius:14px;padding:10px 14px;font-size:13px;color:#486581;line-height:1.6;">
              Оценок: <strong>{{$shift.TotalRatings}}</strong><br>
              Отзывов: <strong>{{$shift.TextReviewsCount}}</strong>
            </div>
          </td>
        </tr>
      </table>

      <h3 style="margin:0 0 12px;font-size:18px;color:#102a43;">Статистика по датам</h3>
      {{range $shift.DateStats}}
      <div style="padding:14px 0;border-bottom:1px solid #e6edf3;">
        <table class="email-date-layout" role="presentation" width="100%" cellspacing="0" cellpadding="0" style="border-collapse:collapse;">
          <tr>
            <td style="vertical-align:top;padding-right:10px;">
              <div style="font-size:16px;font-weight:700;color:#102a43;margin-bottom:6px;">{{.BusinessDateDisplay}} · {{.BusinessWeekday}}</div>
              <div style="font-size:13px;color:#486581;line-height:1.5;">Оценок: <strong>{{.TotalRatings}}</strong> · Отзывов: <strong>{{.TextReviewsCount}}</strong></div>
            </td>
            <td class="email-date-rating" style="vertical-align:top;text-align:right;">
              <div class="email-average-pill" style="display:inline-block;min-width:74px;background:{{avgRatingBackground .AverageRating}};border:1px solid {{avgRatingBorder .AverageRating}};border-radius:12px;padding:8px 10px;color:{{avgRatingText .AverageRating}};text-align:center;">
                <span class="email-average-label" style="display:block;font-size:10px;font-weight:800;letter-spacing:0.03em;text-transform:uppercase;">Средняя</span>
                <span class="email-average-value" style="display:block;margin-top:5px;font-size:17px;font-weight:900;line-height:1;">{{formatAverageRating .AverageRating}}</span>
              </div>
            </td>
          </tr>
        </table>
      </div>
      {{end}}

      <div style="height:18px;"></div>
    </div>

    <div style="height:16px;"></div>
    {{end}}

  </div>
</body>
</html>
`

func newHTMLTemplate(name string) *template.Template {
	return template.New(name).Funcs(template.FuncMap{
		"weekdayRu":           weekdayToRussian,
		"mealRu":              mealToRussian,
		"shiftRu":             shiftToRussian,
		"shiftMealRu":         shiftMealToRussian,
		"formatAverageRating": formatAverageRating,
		"avgRatingBackground": avgRatingBackground,
		"avgRatingBorder":     avgRatingBorder,
		"avgRatingText":       avgRatingText,
	})
}

func shiftToRussian(shiftType string) string {
	switch shiftType {
	case "day":
		return "Дневная смена"
	case "night":
		return "Ночная смена"
	default:
		return shiftType
	}
}

func shiftMealToRussian(shiftType string, mealType string) string {
	if shiftType == "night" {
		switch mealType {
		case "breakfast":
			return "Первый прием пищи"
		case "lunch":
			return "Второй прием пищи"
		case "dinner":
			return "Третий прием пищи"
		}
	}

	return mealToRussian(mealType)
}

func formatAverageRating(value float64) string {
	if value == 0 {
		return "—"
	}

	return fmt.Sprintf("%.1f", value)
}

func avgRatingBackground(value float64) string {
	switch averageRatingBucket(value) {
	case 5:
		return "#ecfdf5"
	case 4:
		return "#f0fdf4"
	case 3:
		return "#fffbeb"
	case 2:
		return "#fff7ed"
	default:
		return "#fef2f2"
	}
}

func avgRatingBorder(value float64) string {
	switch averageRatingBucket(value) {
	case 5:
		return "#bbf7d0"
	case 4:
		return "#dcfce7"
	case 3:
		return "#fde68a"
	case 2:
		return "#fed7aa"
	default:
		return "#fecaca"
	}
}

func avgRatingText(value float64) string {
	switch averageRatingBucket(value) {
	case 5:
		return "#166534"
	case 4:
		return "#15803d"
	case 3:
		return "#92400e"
	case 2:
		return "#9a3412"
	default:
		return "#991b1b"
	}
}

func averageRatingBucket(value float64) int {
	if value == 0 {
		return 3
	}

	if value >= 5 {
		return 5
	}

	bucket := int(math.Floor(value))
	if bucket < 1 {
		return 1
	}
	if bucket > 5 {
		return 5
	}

	return bucket
}
