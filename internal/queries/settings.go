package queries

import (
	"database/sql"
	"desgruppe/internal/entities"
	"desgruppe/internal/logger"
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func GetSettings(c echo.Context, db *sqlx.DB) error {
	settings := entities.Settings{}

	err := db.Get(&settings, `SELECT * FROM "settings" WHERE "ID"=$1`, 1)

	if err != nil {
		if err == sql.ErrNoRows {
			// Если нет строк, создаем новую запись
			_, err := db.Exec(`INSERT INTO "settings" ("ID", "ExchangeRate", "Email") VALUES ($1, $2, $3)`, 1, 100, ".")
			if err != nil {
				logger.Error("Failed to insert default settings into database: " + err.Error())
				return c.String(http.StatusInternalServerError, err.Error())
			}
			settings.ID = 1
			settings.ExchangeRate = 1
			logger.Info("No settings found, inserting default settings with ExchangeRate = 1")
		} else {
			logger.Error("Failed to get settings from database: " + err.Error())
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, settings)
}

func AddSettings(c echo.Context, db *sqlx.DB) error {
	settings := entities.Settings{}

	if err := json.NewDecoder(c.Request().Body).Decode(&settings); err != nil {
		logger.Error("Error parsing response body to add settings: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	// DB
	var err error
	if settings.Email != "" {
		_, err = db.Exec(`
		UPDATE "settings" SET "Email"=$1 WHERE "ID"=1
		`, settings.Email)
	}
	if settings.ExchangeRate != 0 {
		_, err = db.Exec(`
			UPDATE "settings" SET "ExchangeRate"=$1 WHERE "ID"=1
			`, settings.ExchangeRate)
	}
	if err != nil {
		logger.Error("Failed to add settings to database: " + err.Error())
		return err
	}
	logger.Info("Actual currency rate: ", settings.ExchangeRate)
	return c.String(http.StatusCreated, "settings added successfully")
}
