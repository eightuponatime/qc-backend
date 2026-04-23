package impl

import "html/template"

const emailBodyTemplate = `
<!doctype html>
<html lang="ru">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Сводка по качеству питания</title>
</head>
<body style="margin:0;padding:0;background:#f3f5f7;font-family:Segoe UI,Arial,sans-serif;color:#1f2933;">
  <div style="max-width:760px;margin:0 auto;padding:24px 16px;">
    <div style="background:#11324d;color:#ffffff;border-radius:20px;padding:28px 32px;">
      <div style="font-size:13px;letter-spacing:0.08em;text-transform:uppercase;opacity:0.8;">Контроль качества питания</div>
      <h1 style="margin:10px 0 8px;font-size:30px;line-height:1.2;">Сводка за период {{.PeriodStartDisplay}} - {{.PeriodEndDisplay}}</h1>
      <p style="margin:0;font-size:15px;line-height:1.6;opacity:0.9;">Автоматически сформированная управленческая сводка по голосам и отзывам сотрудников.</p>
    </div>

    <div style="height:16px;"></div>

    <div style="background:#e0f2fe;border-radius:20px;padding:24px;border:1px solid #7dd3fc;">
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

    <div style="background:#ffffff;border-radius:20px;padding:24px;border:1px solid #dde4ea;">
      <h2 style="margin:0 0 18px;font-size:22px;color:#102a43;">Статистика по дням недели</h2>
      {{range .WeekdayStats}}
      <div style="padding:16px 0;border-bottom:1px solid #e6edf3;">
        <table role="presentation" width="100%" cellspacing="0" cellpadding="0" style="border-collapse:collapse;">
          <tr>
            <td style="vertical-align:top;width:34%;">
              <div style="font-size:16px;font-weight:700;color:#102a43;margin-bottom:6px;">{{weekdayRu .Weekday}}</div>
              <div style="font-size:13px;color:#486581;line-height:1.5;">Оценок: <strong>{{.TotalRatings}}</strong><br>Отзывов: <strong>{{.TextReviewsCount}}</strong></div>
            </td>
            <td style="vertical-align:top;width:66%;">
              <table role="presentation" width="100%" cellspacing="0" cellpadding="0" style="border-collapse:separate;border-spacing:6px 0;">
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
    </div>

    <div style="height:16px;"></div>

    <div style="background:#ffffff;border-radius:20px;padding:24px;border:1px solid #dde4ea;">
      <h2 style="margin:0 0 18px;font-size:22px;color:#102a43;">По приемам пищи</h2>
      <h3 style="margin:0 0 12px;font-size:16px;color:#334e68;">По дням недели</h3>
      <table role="presentation" width="100%" cellspacing="0" cellpadding="0" style="border-collapse:collapse;border:1px solid #dde4ea;border-radius:14px;overflow:hidden;">
        <tr>
          <th align="left" style="background:#f8fafc;padding:10px 12px;border-bottom:1px solid #dde4ea;color:#486581;font-size:12px;text-transform:uppercase;">День</th>
          <th align="left" style="background:#f8fafc;padding:10px 12px;border-bottom:1px solid #dde4ea;color:#486581;font-size:12px;text-transform:uppercase;">Завтрак</th>
          <th align="left" style="background:#f8fafc;padding:10px 12px;border-bottom:1px solid #dde4ea;color:#486581;font-size:12px;text-transform:uppercase;">Обед</th>
          <th align="left" style="background:#f8fafc;padding:10px 12px;border-bottom:1px solid #dde4ea;color:#486581;font-size:12px;text-transform:uppercase;">Ужин</th>
        </tr>
        {{range .WeekdayStats}}
        <tr>
          <td style="padding:12px;border-bottom:1px solid #e6edf3;color:#102a43;font-weight:700;">{{weekdayRu .Weekday}}</td>
          {{range .MealStats}}
          <td style="padding:12px;border-bottom:1px solid #e6edf3;color:#334e68;font-size:13px;line-height:1.5;">
            <div>Всего: <strong>{{.TotalRatings}}</strong></div>
            <div style="color:#9a3412;">Низких: <strong>{{.LowRatingsCount}}</strong></div>
          </td>
          {{end}}
        </tr>
        {{end}}
      </table>
    </div>

  </div>
</body>
</html>
`

func newHTMLTemplate(name string) *template.Template {
	return template.New(name).Funcs(template.FuncMap{
		"weekdayRu": weekdayToRussian,
		"mealRu":    mealToRussian,
	})
}
