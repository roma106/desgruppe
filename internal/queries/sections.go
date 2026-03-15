package queries

import (
	"desgruppe/internal/entities"
	"desgruppe/internal/logger"
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func ListSections(c echo.Context, db *sqlx.DB) error {
	sections := []entities.Section{}
	rows, err := db.Query(`SELECT * FROM "sections"`)
	if err != nil {
		logger.Error("Failed to select sections from database: " + err.Error())
		return err
	}
	defer rows.Close()
	for rows.Next() {
		section := entities.Section{}
		err = rows.Scan(&section.ID, &section.Name, &section.Type)
		if err != nil {
			logger.Error("Failed to scan sections from database: " + err.Error())
			return err
		}
		sections = append(sections, section)
	}

	return c.JSON(http.StatusOK, sections)
}

func AddSection(c echo.Context, db *sqlx.DB) error {
	addsection := entities.Section{}
	if err := json.NewDecoder(c.Request().Body).Decode(&addsection); err != nil {
		logger.Error("Error parsing response body to add section: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	// DB
	_, err := db.Exec(`
		INSERT INTO "sections" ("Name", "Type")
		VALUES ($1, $2)
		`, addsection.Name, addsection.Type)
	if err != nil {
		logger.Error("Failed to add section to database: " + err.Error())
		return err
	}
	logger.Info("section added: ", addsection.Name)
	return c.String(http.StatusCreated, "section added successfully")
}

func DeleteSection(c echo.Context, db *sqlx.DB) error {
	section := entities.Section{}
	if err := json.NewDecoder(c.Request().Body).Decode(&section); err != nil {
		logger.Error("Error parsing response body to edit color: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	_, err := db.Exec(`DELETE FROM "sections" WHERE "ID"=$1`, section.ID)
	if err != nil {
		logger.Error("Failed to delete section from database: " + err.Error())
		return err
	}
	logger.Info("section deleted: ", section.ID)
	return c.String(http.StatusOK, "section deleted successfully")
}
