package queries

import (
	"desgruppe/internal/entities"
	"desgruppe/internal/logger"
	"desgruppe/internal/utils"
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func ListDesigners(c echo.Context, db *sqlx.DB) error {
	designers := []entities.Designer{}
	rows, err := db.Query(`SELECT * FROM "designers" ORDER BY "Position" ASC`)
	if err != nil {
		logger.Error("Failed to select designers from database: " + err.Error())
		return err
	}
	defer rows.Close()
	for rows.Next() {
		designer := entities.Designer{}
		bytes := []byte{}
		err = rows.Scan(&designer.ID, &designer.Name, &designer.Description, &bytes, &designer.Position, &designer.Slug)
		if err != nil {
			logger.Error("Failed to scan designers from database: " + err.Error())
			return err
		}
		for _, v := range bytes {
			designer.Picture = append(designer.Picture, int(v))
		}
		designers = append(designers, designer)
	}

	return c.JSON(http.StatusOK, designers)
}

func GetDesigner(c echo.Context, db *sqlx.DB) error {
	designer := entities.Designer{}
	rows, err := db.Query(`SELECT * FROM "designers" WHERE "ID"=$1`, c.QueryParam("id"))
	if err != nil {
		logger.Error("Failed to get designer from database: " + err.Error())
		return err
	}
	defer rows.Close()
	for rows.Next() {
		bytes := []byte{}
		err = rows.Scan(&designer.ID, &designer.Name, &designer.Description, &bytes, &designer.Position, &designer.Slug)
		if err != nil {
			logger.Error("Failed to scan designers from database: " + err.Error())
			return err
		}
		for _, v := range bytes {
			designer.Picture = append(designer.Picture, int(v))
		}
	}
	return c.JSON(http.StatusOK, designer)
}
func AddDesigner(c echo.Context, db *sqlx.DB) error {
	addDesigner := entities.Designer{}
	if err := json.NewDecoder(c.Request().Body).Decode(&addDesigner); err != nil {
		logger.Error("Error parsing response body to add color: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	bytes := make([]byte, len(addDesigner.Picture))

	for i, v := range addDesigner.Picture {
		bytes[i] = byte(v)
	}

	// DB
	_, err := db.Exec(`
		INSERT INTO "designers" ("Name", "Description", "Picture", "Position", "Slug")
		VALUES ($1, $2, $3, $4, $5)
		`, addDesigner.Name, addDesigner.Description, bytes, addDesigner.Position, utils.GenerateSlug(addDesigner.Name))
	if err != nil {
		logger.Error("Failed to add designer to database: " + err.Error())
		return err
	}
	logger.Info("designer added: ", addDesigner.Name)
	return c.String(http.StatusCreated, "Designer added successfully")
}
func EditDesigner(c echo.Context, db *sqlx.DB) error {
	addDesigner := entities.Designer{}
	if err := json.NewDecoder(c.Request().Body).Decode(&addDesigner); err != nil {
		logger.Error("Error parsing response body to edit color: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	bytes := make([]byte, len(addDesigner.Picture))

	for i, v := range addDesigner.Picture {
		bytes[i] = byte(v)
	}

	// DB
	var err error
	if len(addDesigner.Picture) == 0 {
		_, err = db.Exec(`
			UPDATE "designers" SET "Name"=$1, "Description"=$2, "Position"=$3, "Slug"=$4 WHERE "ID"=$5
			`, addDesigner.Name, addDesigner.Description, addDesigner.Position, utils.GenerateSlug(addDesigner.Name), addDesigner.ID)
	} else if len(addDesigner.Picture) > 0 {
		_, err = db.Exec(`
			UPDATE "designers" SET "Name"=$1, "Description"=$2, "Picture"=$3, "Position"=$4, "Slug"=$5 WHERE "ID"=$6
			`, addDesigner.Name, addDesigner.Description, bytes, addDesigner.Position, utils.GenerateSlug(addDesigner.Name), addDesigner.ID)
	}
	if err != nil {
		logger.Error("Failed to edit designer to database: " + err.Error())
		return err
	}
	logger.Info("designer edited: ", addDesigner.Name)
	return c.String(http.StatusOK, "Designer updated successfully")
}

func DeleteDesigner(c echo.Context, db *sqlx.DB) error {
	designer := entities.Designer{}
	if err := json.NewDecoder(c.Request().Body).Decode(&designer); err != nil {
		logger.Error("Error parsing response body to edit color: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	_, err := db.Exec(`DELETE FROM "designers" WHERE "ID"=$1`, designer.ID)
	if err != nil {
		logger.Error("Failed to delete designer from database: " + err.Error())
		return err
	}
	err = ClearDeletedValue(db, "DesignerID", designer.ID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	logger.Info("designer deleted: ", designer.ID)
	return c.String(http.StatusOK, "Designer deleted successfully")
}

func SearchDesigner(c echo.Context, db *sqlx.DB) error {
	results := []entities.Designer{}
	query := c.Param("q")
	rows, err := db.Query(`SELECT * FROM "designers" WHERE "Name" COLLATE "C" ILIKE $1`, `%`+query+`%`)
	if err != nil {
		logger.Error("Failed to search designers from database: " + err.Error())
		return err
	}
	defer rows.Close()
	for rows.Next() {
		designer := entities.Designer{}
		bytes := []byte{}
		err = rows.Scan(&designer.ID, &designer.Name, &designer.Description, &bytes, &designer.Position)
		if err != nil {
			logger.Error("Failed to scan designer from database: " + err.Error())
			return err
		}
		for _, v := range bytes {
			designer.Picture = append(designer.Picture, int(v))
		}
		results = append(results, designer)
	}
	return c.JSON(http.StatusOK, results)
}

func GetDesignerBySlug(c echo.Context, db *sqlx.DB) (entities.Designer, error) {
	designer := entities.Designer{}
	rows, err := db.Query(`SELECT * FROM "designers" WHERE "Slug"=$1`, c.Param("slug"))
	if err != nil {
		logger.Error("Failed to get designer from database: " + err.Error())
		return designer, err
	}
	defer rows.Close()
	for rows.Next() {
		bytes := []byte{}
		err = rows.Scan(&designer.ID, &designer.Name, &designer.Description, &bytes, &designer.Position, &designer.Slug)
		if err != nil {
			logger.Error("Failed to scan designers from database: " + err.Error())
			return designer, err
		}
		for _, v := range bytes {
			designer.Picture = append(designer.Picture, int(v))
		}
	}
	return designer, nil
}
func GetDesignerById(db *sqlx.DB, ID int) (entities.Designer, error) {
	designer := entities.Designer{}
	rows, err := db.Query(`SELECT "ID","Name","Description","Position","Slug" FROM "designers" WHERE "ID"=$1`, ID)
	if err != nil {
		logger.Error("Failed to get designer from database: " + err.Error())
		return designer, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&designer.ID, &designer.Name, &designer.Description, &designer.Position, &designer.Slug)
		if err != nil {
			logger.Error("Failed to scan designers from database: " + err.Error())
			return designer, err
		}
	}
	return designer, nil
}
