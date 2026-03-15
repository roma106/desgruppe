package queries

import (
	"desgruppe/internal/entities"
	"desgruppe/internal/logger"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func ListColors(c echo.Context, db *sqlx.DB) error {
	colors := []entities.Color{}

	rows, err := db.Query(`SELECT * FROM "colors" ORDER BY "Position" ASC`)
	if err != nil {
		logger.Error("Failed to select colors from database: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		color := entities.Color{}
		var code any
		bytes := []byte{}
		err = rows.Scan(&color.ID, &color.Name, &code, &bytes, &color.Position)
		if err != nil {
			logger.Error("Failed to scan colors from database: " + err.Error())
			return c.String(http.StatusInternalServerError, err.Error())
		}

		if code == nil {
			for _, v := range bytes {
				color.Picture = append(color.Picture, int(v))
			}
		} else {
			color.Code = code.(string)
		}

		colors = append(colors, color)
	}

	return c.JSON(http.StatusOK, colors)
}

func GetColor(c echo.Context, db *sqlx.DB) error {
	color := entities.Color{}
	id := c.QueryParam("id")

	rows, err := db.Query(`SELECT * FROM "colors" WHERE "ID"=$1`, id)
	if err != nil {
		logger.Error("Failed to select color from database: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		bytes := []byte{}
		var code any
		err = rows.Scan(&color.ID, &color.Name, &code, &bytes, &color.Position)
		if err != nil {
			logger.Error("Failed to scan color from database: " + err.Error())
			return c.String(http.StatusInternalServerError, err.Error())
		}

		if code == nil {
			for _, v := range bytes {
				color.Picture = append(color.Picture, int(v))
			}
		} else {
			color.Code = code.(string)
		}
	}

	return c.JSON(http.StatusOK, color)
}

func AddColor(c echo.Context, db *sqlx.DB) error {
	addcolor := entities.Color{}

	if err := json.NewDecoder(c.Request().Body).Decode(&addcolor); err != nil {
		logger.Error("Error parsing response body to add color: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	bytes := make([]byte, len(addcolor.Picture))
	for i, v := range addcolor.Picture {
		bytes[i] = byte(v)
	}

	var err error
	if addcolor.IsCode {
		_, err = db.Exec(`INSERT INTO "colors" ("Name", "Code", "Position") VALUES ($1, $2, $3)`, addcolor.Name, addcolor.Code, addcolor.Position)
	} else if addcolor.IsPicture {
		_, err = db.Exec(`INSERT INTO "colors" ("Name", "Picture", "Position") VALUES ($1, $2, $3)`, addcolor.Name, bytes, addcolor.Position)
	}

	if err != nil {
		logger.Error("Failed to add color to database: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	logger.Info("Color added: ", addcolor.Name)
	return c.String(http.StatusCreated, "Color added successfully")
}

func EditColor(c echo.Context, db *sqlx.DB) error {
	addcolor := entities.Color{}

	if err := json.NewDecoder(c.Request().Body).Decode(&addcolor); err != nil {
		logger.Error("Error parsing response body to edit color: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	bytes := make([]byte, len(addcolor.Picture))
	for i, v := range addcolor.Picture {
		bytes[i] = byte(v)
	}

	var err error
	if addcolor.IsCode {
		_, err = db.Exec(`UPDATE "colors" SET "Name"=$1, "Code"=$2, "Picture"=$3, "Position"=$4 WHERE "ID"=$5`, addcolor.Name, addcolor.Code, nil, addcolor.Position, addcolor.ID)
	} else if addcolor.IsPicture {
		if len(addcolor.Picture) == 0 {
			_, err = db.Exec(`UPDATE "colors" SET "Name"=$1, "Position"=$2 WHERE "ID"=$3`, addcolor.Name, addcolor.Position, addcolor.ID)
		} else if len(addcolor.Picture) > 0 {
			_, err = db.Exec(`UPDATE "colors" SET "Name"=$1, "Code"=$2, "Picture"=$3, "Position"=$4 WHERE "ID"=$5`, addcolor.Name, nil, bytes, addcolor.Position, addcolor.ID)
		}
	}

	if err != nil {
		logger.Error("Failed to edit color to database: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	logger.Info("Color edited: ", addcolor.Name)
	return c.String(http.StatusOK, "Color updated successfully")
}

func DeleteColor(c echo.Context, db *sqlx.DB) error {
	color := entities.Color{}

	if err := json.NewDecoder(c.Request().Body).Decode(&color); err != nil {
		logger.Error("Error parsing response body to delete color: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	_, err := db.Exec(`DELETE FROM "colors" WHERE "ID"=$1`, color.ID)
	if err != nil {
		logger.Error("Failed to delete color from database: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	logger.Info("Color deleted: ", color.ID)
	return c.String(http.StatusOK, "Color deleted successfully")
}

func SearchColor(c echo.Context, db *sqlx.DB) error {
	results := []entities.Color{}
	query := c.QueryParam("q")

	rows, err := db.Query(`SELECT * FROM "colors" WHERE "Name" COLLATE "C" ILIKE $1`, `%`+query+`%`)
	if err != nil {
		logger.Error("Failed to search colors from database: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		color := entities.Color{}
		bytes := []byte{}
		var code any
		err = rows.Scan(&color.ID, &color.Name, &code, &bytes, &color.Position)
		if err != nil {
			logger.Error("Failed to scan color from database: " + err.Error())
			return c.String(http.StatusInternalServerError, err.Error())
		}

		if code == nil {
			for _, v := range bytes {
				color.Picture = append(color.Picture, int(v))
			}
		} else {
			color.Code = code.(string)
		}
		results = append(results, color)
	}

	return c.JSON(http.StatusOK, results)
}

func ListColorsByIds(db *sqlx.DB, ids []int) ([]entities.Color, error) {
	colors := []entities.Color{}
	for _, id := range ids {
		if id == 0 {
			colors = append(colors, entities.Color{ID: 0, Name: ""})
			continue
		}
		rows, err := db.Query(`SELECT "ID","Name","Code","Position" FROM "colors" WHERE "ID"=$1`, strconv.Itoa(id))
		if err != nil {
			logger.Error("Failed to select colors from database: " + err.Error())
			return colors, err
		}
		defer rows.Close()

		for rows.Next() {
			color := entities.Color{}
			var code any
			err = rows.Scan(&color.ID, &color.Name, &code, &color.Position)
			if err != nil {
				logger.Error("Failed to scan colors from database: " + err.Error())
				return colors, err
			}
			if code != nil {
				color.Code = code.(string)
			}

			colors = append(colors, color)
		}
	}

	return colors, nil
}
