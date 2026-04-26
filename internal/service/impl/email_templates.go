package impl

import "html/template"

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
      .email-meal-table td,
      .email-meal-table th {
        display: block !important;
        width: 100% !important;
        box-sizing: border-box !important;
      }
      .email-weekday-stars td {
        display: inline-block !important;
        width: calc(20% - 8px) !important;
        margin-bottom: 8px !important;
      }
      .email-meal-table tr {
        display: block !important;
        border-bottom: 1px solid #e6edf3 !important;
      }
      .email-meal-table__head {
        display: none !important;
      }
      .email-meal-table td {
        border-bottom: none !important;
        padding: 8px 12px !important;
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

      <h3 style="margin:0 0 12px;font-size:18px;color:#102a43;">Статистика по дням недели</h3>
      {{range $shift.WeekdayStats}}
      <div style="padding:14px 0;border-bottom:1px solid #e6edf3;">
        <table class="email-weekday-layout" role="presentation" width="100%" cellspacing="0" cellpadding="0" style="border-collapse:collapse;">
          <tr>
            <td style="vertical-align:top;padding-right:10px;">
              <div style="font-size:16px;font-weight:700;color:#102a43;margin-bottom:6px;">{{weekdayRu .Weekday}}</div>
              <div style="font-size:13px;color:#486581;line-height:1.5;">Оценок: <strong>{{.TotalRatings}}</strong> · Отзывов: <strong>{{.TextReviewsCount}}</strong></div>
            </td>
            <td style="vertical-align:top;">
              <table class="email-weekday-stars" role="presentation" width="100%" cellspacing="0" cellpadding="0" style="border-collapse:separate;border-spacing:6px 0;">
                <tr>
                  <td align="center" style="background:#ecfdf5;border:1px solid #bbf7d0;border-radius:12px;padding:9px 4px;color:#166534;font-size:13px;"><strong>5★</strong><br>{{.RatingDistribution.Five}}</td>
                  <td align="center" style="background:#f0fdf4;border:1px solid #dcfce7;border-radius:12px;padding:9px 4px;color:#15803d;font-size:13px;"><strong>4★</strong><br>{{.RatingDistribution.Four}}</td>
                  <td align="center" style="background:#fffbeb;border:1px solid #fde68a;border-radius:12px;padding:9px 4px;color:#92400e;font-size:13px;"><strong>3★</strong><br>{{.RatingDistribution.Three}}</td>
                  <td align="center" style="background:#fff7ed;border:1px solid #fed7aa;border-radius:12px;padding:9px 4px;color:#9a3412;font-size:13px;"><strong>2★</strong><br>{{.RatingDistribution.Two}}</td>
                  <td align="center" style="background:#fef2f2;border:1px solid #fecaca;border-radius:12px;padding:9px 4px;color:#991b1b;font-size:13px;"><strong>1★</strong><br>{{.RatingDistribution.One}}</td>
                </tr>
              </table>
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
		"weekdayRu":   weekdayToRussian,
		"mealRu":      mealToRussian,
		"shiftRu":     shiftToRussian,
		"shiftMealRu": shiftMealToRussian,
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
