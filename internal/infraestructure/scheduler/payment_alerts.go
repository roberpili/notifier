package scheduler

import (
	"context"
	"fmt"
	"html"
	"log"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rober0xf/notifier/pkg/database"
)

type payment struct {
	Email   string
	Name    string
	Amount  float64
	PayType string
}

func (a *Alerter) SendPaymentAlert(title string, targetDate time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	dateStr := targetDate.Format("2006-01-02")
	log.Printf("starting SendPaymentAlert for date: %s", dateStr)

	query := `
				SELECT users.email, payments.name, payments.amount, payments.type
				FROM users
				LEFT JOIN payments ON payments.user_id = users.id
				WHERE DATE(payments.due_date) = DATE($1)
				`
	rows, err := database.DB.Query(ctx, query, targetDate)
	if err != nil {
		log.Printf("database error: %v", err)
		return err
	}
	defer rows.Close()

	byUser := make(map[string][]payment)
	for rows.Next() {
		var userEmail string
		var p payment
		if err := rows.Scan(&userEmail, &p.Name, &p.Amount, &p.PayType); err != nil {
			log.Printf("error scanning row: %v", err)
			continue
		}
		byUser[userEmail] = append(byUser[userEmail], p)
	}
	if err := rows.Err(); err != nil {
		log.Printf("error iterating rows: %v", err)
		return err
	}

	if len(byUser) == 0 {
		log.Printf("there are no payments due on %s", dateStr)
		return nil
	}

	sent, failed := 0, 0
	for addr, payments := range byUser {
		body := buildBody(dateStr, payments)

		sendCtx, sendCancel := context.WithTimeout(ctx, 15*time.Second)
		err := a.sender.Send(sendCtx, []string{addr}, title, body)
		sendCancel()

		if err != nil {
			log.Printf("error sending email to %s: %v", addr, err)
			failed++
			continue
		}
		log.Printf("email sent to: %s", addr)
		sent++
	}

	log.Printf("email stats - sent: %d, failed: %d", sent, failed)
	return nil
}

func buildBody(dateStr string, payments []payment) string {
	var b strings.Builder
	fmt.Fprintf(&b, "<p>These are the payments that will be executed on %s:</p><ul>", html.EscapeString(dateStr))
	for _, p := range payments {
		fmt.Fprintf(&b, "<li><strong>%s</strong> &mdash; %.2f (%s)</li>",
			html.EscapeString(p.Name), p.Amount, html.EscapeString(p.PayType))
	}
	b.WriteString("</ul>")
	return b.String()
}

func (a *Alerter) SendDailyAlert() {
	a.run("SendDailyAlert", "daily payment alert", time.Now().AddDate(0, 0, 1))
}

func (a *Alerter) SendWeeklyAlert() {
	a.run("SendWeeklyAlert", "weekly payment alert", time.Now().AddDate(0, 0, 7))
}

func (a *Alerter) SendMonthlyAlert() {
	a.run("SendMonthlyAlert", "monthly payment alert", time.Now().AddDate(0, 1, 0))
}
