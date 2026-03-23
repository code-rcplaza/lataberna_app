package email

import (
	"context"
	"fmt"

	resend "github.com/resend/resend-go/v2"
)

type ResendMailer struct {
	client *resend.Client
	from   string
}

func NewResendMailer(apiKey, from string) *ResendMailer {
	return &ResendMailer{
		client: resend.NewClient(apiKey),
		from:   from,
	}
}

func (m *ResendMailer) SendMagicLink(ctx context.Context, email, magicLink string) error {
	params := &resend.SendEmailRequest{
		From:    m.from,
		To:      []string{email},
		Subject: "Tu enlace mágico — La Taberna",
		Html:    buildMagicLinkEmail(magicLink),
	}
	_, err := m.client.Emails.Send(params)
	return err
}

func buildMagicLinkEmail(link string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="es">
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"></head>
<body style="margin:0;padding:0;background:#fdf6ee;font-family:Georgia,serif;">
  <div style="max-width:480px;margin:48px auto;background:#fff;border:1px solid #e8ddd0;padding:40px 32px;">
    <h1 style="font-size:24px;color:#3b1f0e;margin:0 0 8px;font-style:italic;">La Taberna</h1>
    <p style="color:#7a5c3e;font-size:12px;letter-spacing:2px;text-transform:uppercase;margin:0 0 32px;">Tu puerta al multiverso</p>
    <p style="color:#3b1f0e;font-size:16px;line-height:1.6;margin:0 0 24px;">
      Alguien — esperamos que seas vos — solicitó acceso a La Taberna.<br>
      Hacé clic en el botón para ingresar. El enlace expira en <strong>15 minutos</strong> y solo puede usarse una vez.
    </p>
    <a href="%s"
       style="display:inline-block;background:#3b1f0e;color:#fdf6ee;text-decoration:none;padding:14px 28px;font-size:12px;letter-spacing:2px;text-transform:uppercase;font-family:Georgia,serif;font-weight:bold;">
      Entrar a La Taberna
    </a>
    <p style="color:#a08060;font-size:12px;margin:32px 0 0;line-height:1.5;">
      Si no pediste este enlace, ignorá este mensaje. No pasa nada.
    </p>
    <hr style="border:none;border-top:1px solid #e8ddd0;margin:24px 0;">
    <p style="color:#c0a882;font-size:11px;margin:0;">lataberna.app</p>
  </div>
</body>
</html>`, link)
}
