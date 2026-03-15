package notification

import (
	"desgruppe/internal/logger"
	"desgruppe/internal/templ/api"
	"desgruppe/internal/utils"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"gopkg.in/gomail.v2"
)

func SendOrderNotification(c echo.Context, db *sqlx.DB) {
	settings := utils.GetSettings(db)
	mailHtml := api.RenderMailPage(c, db)
	m := gomail.NewMessage()
	m.SetHeader("From", "desgruppemanager@yandex.ru")
	m.SetHeader("To", settings.Email)
	m.SetHeader("Subject", "Заказ от "+time.Now().Format("02.01.2006"))
	m.SetBody("text/html", mailHtml)
	d := gomail.NewDialer("smtp.yandex.ru", 587, "desgruppemanager@yandex.ru", "kdnqybhhrneaogpa")
	if err := d.DialAndSend(m); err != nil {
		logger.Error("Не удалось отправить email:", err)
	}
}
