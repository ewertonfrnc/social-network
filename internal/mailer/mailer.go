package mailer

import "embed"

const (
	fromEmail                 = "Social Network"
	maxRetries                = 3
	UserWelcomeInviteTemplate = "send_invitation.tmpl"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile, username, email string, data any, isSandbox bool) error
}
