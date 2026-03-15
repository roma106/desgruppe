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

func ListProducers(c echo.Context, db *sqlx.DB) error {
	producers := []entities.Producer{}
	rows, err := db.Query(`SELECT * FROM "producers" ORDER BY "Position" ASC`)
	if err != nil {
		logger.Error("Failed to select producers from database: " + err.Error())
		return err
	}
	defer rows.Close()
	for rows.Next() {
		producer := entities.Producer{}
		bytes := []byte{}
		err = rows.Scan(&producer.ID, &producer.Name, &producer.Description, &bytes, &producer.Position, &producer.Slug)
		if err != nil {
			logger.Error("Failed to scan producers from database: " + err.Error())
			return err
		}
		for _, v := range bytes {
			producer.Picture = append(producer.Picture, int(v))
		}
		producers = append(producers, producer)
	}

	return c.JSON(http.StatusOK, producers)
}

func GetProducer(c echo.Context, db *sqlx.DB) error {
	producer := entities.Producer{}
	rows, err := db.Query(`SELECT * FROM "producers" WHERE "ID"=$1`, c.QueryParam("id"))
	if err != nil {
		logger.Error("Failed to get producer from database: " + err.Error())
		return err
	}
	defer rows.Close()
	for rows.Next() {
		bytes := []byte{}
		err = rows.Scan(&producer.ID, &producer.Name, &producer.Description, &bytes, &producer.Position, &producer.Slug)
		if err != nil {
			logger.Error("Failed to scan producers from database: " + err.Error())
			return err
		}
		for _, v := range bytes {
			producer.Picture = append(producer.Picture, int(v))
		}
	}
	return c.JSON(http.StatusOK, producer)
}
func AddProducer(c echo.Context, db *sqlx.DB) error {
	addproducer := entities.Producer{}
	if err := json.NewDecoder(c.Request().Body).Decode(&addproducer); err != nil {
		logger.Error("Error parsing response body to add color: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	bytes := make([]byte, len(addproducer.Picture))

	for i, v := range addproducer.Picture {
		bytes[i] = byte(v)
	}

	// DB
	_, err := db.Exec(`
		INSERT INTO "producers" ("Name", "Description", "Picture", "Position", "Slug")
		VALUES ($1, $2, $3, $4, $5)
		`, addproducer.Name, addproducer.Description, bytes, addproducer.Position, utils.GenerateSlug(addproducer.Name))
	if err != nil {
		logger.Error("Failed to add producer to database: " + err.Error())
		return err
	}
	logger.Info("producer added: ", addproducer.Name)
	return c.String(http.StatusCreated, "producer added successfully")
}
func EditProducer(c echo.Context, db *sqlx.DB) error {
	addproducer := entities.Producer{}
	if err := json.NewDecoder(c.Request().Body).Decode(&addproducer); err != nil {
		logger.Error("Error parsing response body to edit color: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	bytes := make([]byte, len(addproducer.Picture))

	for i, v := range addproducer.Picture {
		bytes[i] = byte(v)
	}

	// DB
	var err error
	if len(addproducer.Picture) == 0 {
		_, err = db.Exec(`
			UPDATE "producers" SET "Name"=$1, "Description"=$2, "Position"=$3, "Slug"=$4 WHERE "ID"=$5
			`, addproducer.Name, addproducer.Description, addproducer.Position, utils.GenerateSlug(addproducer.Name), addproducer.ID)
	} else if len(addproducer.Picture) > 0 {
		_, err = db.Exec(`
			UPDATE "producers" SET "Name"=$1, "Description"=$2, "Picture"=$3, "Position"=$4, "Slug"=$5 WHERE "ID"=$6
			`, addproducer.Name, addproducer.Description, bytes, addproducer.Position, utils.GenerateSlug(addproducer.Name), addproducer.ID)
	}
	if err != nil {
		logger.Error("Failed to edit producer to database: " + err.Error())
		return err
	}
	logger.Info("producer edited: ", addproducer.Name)
	return c.String(http.StatusOK, "producer updated successfully")
}

func DeleteProducer(c echo.Context, db *sqlx.DB) error {
	producer := entities.Producer{}
	if err := json.NewDecoder(c.Request().Body).Decode(&producer); err != nil {
		logger.Error("Error parsing response body to edit color: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	_, err := db.Exec(`DELETE FROM "producers" WHERE "ID"=$1`, producer.ID)
	if err != nil {
		logger.Error("Failed to delete producer from database: " + err.Error())
		return err
	}
	err = ClearDeletedValue(db, "ProducerID", producer.ID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	logger.Info("producer deleted: ", producer.ID)
	return c.String(http.StatusOK, "producer deleted successfully")
}

func SearchProducer(c echo.Context, db *sqlx.DB) error {
	results := []entities.Producer{}
	query := c.Param("q")
	rows, err := db.Query(`SELECT * FROM "producers" WHERE "Name" COLLATE "C" ILIKE $1`, `%`+query+`%`)
	if err != nil {
		logger.Error("Failed to search producers from database: " + err.Error())
		return err
	}
	defer rows.Close()
	for rows.Next() {
		producer := entities.Producer{}
		bytes := []byte{}
		err = rows.Scan(&producer.ID, &producer.Name, &producer.Description, &bytes, &producer.Position, &producer.Slug)
		if err != nil {
			logger.Error("Failed to scan producer from database: " + err.Error())
			return err
		}
		for _, v := range bytes {
			producer.Picture = append(producer.Picture, int(v))
		}
		results = append(results, producer)
	}
	return c.JSON(http.StatusOK, results)
}

func GetProducerBySlug(c echo.Context, db *sqlx.DB) (entities.Producer, error) {
	producer := entities.Producer{}
	rows, err := db.Query(`SELECT * FROM "producers" WHERE "Slug"=$1`, c.Param("slug"))
	if err != nil {
		logger.Error("Failed to get producer from database: " + err.Error())
		return producer, err
	}
	defer rows.Close()
	for rows.Next() {
		bytes := []byte{}
		err = rows.Scan(&producer.ID, &producer.Name, &producer.Description, &bytes, &producer.Position, &producer.Slug)
		if err != nil {
			logger.Error("Failed to scan producers from database: " + err.Error())
			return producer, err
		}
		for _, v := range bytes {
			producer.Picture = append(producer.Picture, int(v))
		}
	}
	return producer, nil
}

func GetProducerById(db *sqlx.DB, ID int) (entities.Producer, error) {
	producer := entities.Producer{}
	rows, err := db.Query(`SELECT "ID","Name","Description","Position","Slug" FROM "producers" WHERE "ID"=$1`, ID)
	if err != nil {
		logger.Error("Failed to get producer from database: " + err.Error())
		return producer, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&producer.ID, &producer.Name, &producer.Description, &producer.Position, &producer.Slug)
		if err != nil {
			logger.Error("Failed to scan producers from database: " + err.Error())
			return producer, err
		}
	}
	return producer, nil
}
