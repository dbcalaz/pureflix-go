package mail

import (
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func EnviarMailDeValidacion(emailUsuario string, tokenActivacion string) {
	m := gomail.NewMessage()

	from := os.Getenv("MAIL_FROM")
	mailUser := os.Getenv("MAIL_USER")
	mailPass := os.Getenv("MAIL_PASSWORD")
	mailHost := os.Getenv("MAIL_HOST")
	mailPort, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	frontURL := os.Getenv("FRONT_URL")

	m.SetHeader("From", from)
	m.SetHeader("To", emailUsuario)
	m.SetHeader("Subject", "Activá tu cuenta Pureflix")
	m.SetBody("text/html", `
		<h2>Bienvenido</h2>
		<p>Para activar tu cuenta:</p>
		<a href="`+frontURL+`/activar-cuenta?token=`+tokenActivacion+`">
			Activar cuenta
		</a>
	`)

	d := gomail.NewDialer(mailHost, mailPort, mailUser, mailPass)
	d.DialAndSend(m)
}

func EnviarMailDeRecuperacion(emailUsuario string, tokenRecuperacion string) {
	m := gomail.NewMessage()

	from := os.Getenv("MAIL_FROM")
	mailUser := os.Getenv("MAIL_USER")
	mailPass := os.Getenv("MAIL_PASSWORD")
	mailHost := os.Getenv("MAIL_HOST")
	mailPort, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	frontURL := os.Getenv("FRONT_URL")

	m.SetHeader("From", from)
	m.SetHeader("To", emailUsuario)
	m.SetHeader("Subject", "Actualizá tu contraseña")
	m.SetBody("text/html", `
		<h2>Recuperación de contraseña</h2>
		<p>Hacé clic para actualizar:</p>
		<a href="`+frontURL+`/new-pass?token=`+tokenRecuperacion+`">
			Actualizar contraseña
		</a>
	`)

	d := gomail.NewDialer(mailHost, mailPort, mailUser, mailPass)
	d.DialAndSend(m)
}
