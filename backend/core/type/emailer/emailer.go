package emailer

// Emailer -
type Emailer interface {
	SendEmail(from, to, subject, emailBody string) error
	SendHTMLEmail(from, to, subject, emailBody string) error
}
