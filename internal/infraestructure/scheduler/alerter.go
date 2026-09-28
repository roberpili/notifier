package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/rober0xf/notifier/pkg/email"
)

type Alerter struct {
	sender email.EmailSender
}

func NewAlerter(sender email.EmailSender) *Alerter {
	return &Alerter{sender: sender}
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
	log.Println("cronjobs set")
	return s, nil
}

func (a *Alerter) run(job, title string, target time.Time) {
	log.Printf("========== %s triggered ==========", job)
	if err := a.SendPaymentAlert(title, target); err != nil {
		log.Printf("%s error: %v", job, err)
		return
	}
	log.Printf("========== %s finished ==========", job)
}
