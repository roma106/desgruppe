package html

import (
	"desgruppe/internal/logger"
	"desgruppe/internal/queries"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func LoadHtmlPage(c echo.Context, name string) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("frontend/%s", name))
	if err != nil {
		logger.Error("Error loading template: ", err)
		return c.String(http.StatusInternalServerError, "Ошибка загрузки шаблона: "+err.Error())
	}
	return tmpl.Execute(c.Response().Writer, nil)
}

func LoadImg(c echo.Context) error {
	filename := strings.TrimPrefix(c.Request().URL.Path, "/imgs/")
	filePath := fmt.Sprintf("frontend/imgs/%s", filename)
	return c.File(filePath)
}

func LoadStaticFile(c echo.Context) error {
	filename := strings.TrimPrefix(c.Request().URL.Path, "/static/")
	filePath := fmt.Sprintf("frontend/static/%s", filename)
	logger.Info(filePath)
	return c.File(filePath)
}

func LoadAdminPage(c echo.Context, name string) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("frontend/admin/admin-%s.html", name))
	if err != nil {
		logger.Error("Error loading template: ", err)
		return c.String(http.StatusInternalServerError, "Ошибка загрузки шаблона: "+err.Error())
	}
	return tmpl.Execute(c.Response().Writer, nil)
}

func LoadAdminEditPage(c echo.Context, name string) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("frontend/admin/edit/admin-edit-%s.html", name))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Ошибка загрузки шаблона: "+err.Error())
	}
	return tmpl.Execute(c.Response().Writer, nil)
}

func LoadProductPage(c echo.Context, db *sqlx.DB) error {
	tmpl, err := template.ParseFiles("frontend/product.html")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Ошибка загрузки шаблона: "+err.Error())
	}
	product, err := queries.GetProductBySlug(db, c.Param("slug"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Ошибка загрузки шаблона: "+err.Error())
	}
	return tmpl.Execute(c.Response().Writer, product)
}

func LoadDesignerPage(c echo.Context, db *sqlx.DB) error {
	tmpl, err := template.ParseFiles("frontend/designer.html")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Ошибка загрузки шаблона: "+err.Error())
	}
	product, err := queries.GetDesignerBySlug(c, db)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Ошибка загрузки шаблона: "+err.Error())
	}
	return tmpl.Execute(c.Response().Writer, product)
}

func LoadProducerPage(c echo.Context, db *sqlx.DB) error {
	tmpl, err := template.ParseFiles("frontend/producer.html")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Ошибка загрузки шаблона: "+err.Error())
	}
	product, err := queries.GetProducerBySlug(c, db)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Ошибка загрузки шаблона: "+err.Error())
	}
	return tmpl.Execute(c.Response().Writer, product)
}
