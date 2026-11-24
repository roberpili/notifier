package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-co-op/gocron/v2"
	email "github.com/rober0xf/notifier/internal/infraestructure/email"
)

type Alerter struct {
	sender email.EmailSender
	logger *slog.Logger
}

func NewAlerter(sender email.EmailSender, logger *slog.Logger) *Alerter {
	if logger == nil {
		logger = slog.Default()
	}

	return &Alerter{sender: sender, logger: logger}
}

func (a *Alerter) InitCron() (gocron.Scheduler, error) {
	s, err := gocron.NewScheduler(gocron.WithLocation(time.Local))
	if err != nil {
		return nil, fmt.Errorf("error creating scheduler: %w", err)
	}

	at7 := gocron.NewAtTimes(gocron.NewAtTime(7, 0, 0))

	if _, err = s.NewJob(gocron.DailyJob(1, at7), gocron.NewTask(a.SendDailyAlert)); err != nil {
		return nil, fmt.Errorf("daily job: %w", err)
	}
	if _, err = s.NewJob(
		gocron.WeeklyJob(1, gocron.NewWeekdays(time.Monday), at7),
		gocron.NewTask(a.SendWeeklyAlert),
	); err != nil {
		return nil, fmt.Errorf("weekly job: %w", err)
	}
	if _, err = s.NewJob(
		gocron.MonthlyJob(1, gocron.NewDaysOfTheMonth(1), at7),
		gocron.NewTask(a.SendMonthlyAlert),
	); err != nil {
		return nil, fmt.Errorf("monthly job: %w", err)
	}

	s.Start() // non-blocking, no goroutine or select{} needed

	a.logger.Info("cron scheduler started",
		slog.String("timezone", time.Local.String()),
		slog.String("daily", "07:00"),
		slog.String("weekly", "Monday 07:00"),
		slog.String("monthly", "1st day 07:00"),
	)

	return s, nil
}

func (a *Alerter) run(ctx context.Context, job, title string, target time.Time) {
	a.logger.InfoContext(ctx, "cron job triggered",
		slog.String("job", job),
		slog.Time("target_date", target),
	)

	if err := a.SendPaymentAlert(ctx, title, target); err != nil {
		a.logger.ErrorContext(ctx, "cron job failed",
			slog.String("job", job),
			slog.Any("error", err),
		)
		return
	}

	a.logger.InfoContext(ctx, "cron job completed",
		slog.String("job", job),
	)
}
