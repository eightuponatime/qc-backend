package impl

import (
	"context"
	"fmt"
	"qc/config"
	"qc/internal/dto"
	"qc/internal/service"
	"strconv"
	"strings"
	"time"

	"github.com/wneessen/go-mail"
)

type EmailService struct {
	reportService          service.ReportService
	analyticsAccessService service.AnalyticsAccessService
	cfg                    *config.Config
}

func NewEmailService(
	reportService service.ReportService,
	analyticsAccessService service.AnalyticsAccessService,
	cfg *config.Config,
) *EmailService {
	return &EmailService{
		reportService:          reportService,
		analyticsAccessService: analyticsAccessService,
		cfg:                    cfg,
	}
}

func (e *EmailService) SendEmail(ctx context.Context) error {
	summary, err := e.reportService.CreateSummary(ctx)
	if err != nil {
		return fmt.Errorf("couldn't get email summary: %w", err)
	}

	return e.SendPeriodReport(ctx, mustParseDate(summary.PeriodStart), mustParseDate(summary.PeriodEnd))
}

// main logic of configuring smtp server
func (e *EmailService) SendPeriodReport(ctx context.Context, periodStart, periodEnd time.Time) error {
	emailSummary, err := e.reportService.CreateSummaryForPeriod(ctx, periodStart, periodEnd)
	if err != nil {
		return fmt.Errorf("couldn't get email summary: %w", err)
	}

	accessValidFrom := periodEnd.AddDate(0, 0, 1)
	accessValidUntil := accessValidFrom.AddDate(0, 0, 14)
	accessCode, err := e.analyticsAccessService.CreateAccessCode(ctx, accessValidFrom, accessValidUntil)
	if err != nil {
		return fmt.Errorf("couldn't create analytics access code: %w", err)
	}

	m := mail.NewMsg()

	if err := m.From(e.cfg.SmtpFrom); err != nil {
		return fmt.Errorf("failed to set the email initiator: %w", err)
	}

	recipients := buildRecipients(e.cfg.ReportTo)
	if len(recipients) == 0 {
		return fmt.Errorf("no email recipients configured")
	}

	if err := m.To(recipients...); err != nil {
		return fmt.Errorf("failed to set target recipients for the report: %w", err)
	}

	subject := fmt.Sprintf(
		"[Вахта %s] Контроль качества еды",
		emailSummary.PeriodShortDisplay,
	)
	m.Subject(subject)

	emailData := reportEmailTemplateData{
		ReportSummaryDto: *emailSummary,
		AnalyticsURL:     e.cfg.AnalyticsURL,
		AccessCode:       accessCode,
		AccessValidUntil: formatRussianDate(accessValidUntil),
	}

	m, err = e.buildEmailBody(emailData, m)
	if err != nil {
		return fmt.Errorf("failed to build email body: %w", err)
	}

	port, err := strconv.Atoi(e.cfg.SmtpPort)
	if err != nil {
		return fmt.Errorf("couldn't parse smtp port: %w", err)
	}

	client, err := mail.NewClient(
		e.cfg.SmtpHost,
		mail.WithPort(port),
		mail.WithTLSPolicy(mail.TLSOpportunistic),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(e.cfg.SmtpUsername),
		mail.WithPassword(e.cfg.SmtpPassword),
	)
	if err != nil {
		return fmt.Errorf("failed to create smtp client: %w", err)
	}

	if err := client.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

type reportEmailTemplateData struct {
	dto.ReportSummaryDto
	AnalyticsURL     string
	AccessCode       string
	AccessValidUntil string
}

func (e *EmailService) buildEmailBody(data reportEmailTemplateData, m *mail.Msg) (*mail.Msg, error) {
	tpl, err := newHTMLTemplate("email_body").Parse(emailBodyTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html template: %w", err)
	}

	if err := m.SetBodyHTMLTemplate(tpl, data); err != nil {
		return nil, fmt.Errorf("failed to set html template into email body: %w", err)
	}

	return m, nil
}

func mustParseDate(date string) time.Time {
	parsed, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic(fmt.Sprintf("parse date %q: %v", date, err))
	}

	return parsed
}

func buildRecipients(reportTo string) []string {
	parts := strings.Split(reportTo, ",")
	recipients := make([]string, 0, len(parts))

	for _, part := range parts {
		recipient := strings.TrimSpace(part)
		if recipient == "" {
			continue
		}
		recipients = append(recipients, recipient)
	}

	return recipients
}
