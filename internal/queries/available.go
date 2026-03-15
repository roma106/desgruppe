package queries

import (
	"desgruppe/internal/entities"
	"desgruppe/internal/logger"
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

// func ListAvailable(c echo.Context, db *sqlx.DB) error {
// 	products := []entities.Product{}
// 	rows, err := db.Query(`SELECT "ID","Name","Type","Photo","Price","OnSale","Sale","Slug","Position" FROM "products" WHERE "Available"=true ORDER BY "Position"`)
// 	if err != nil {
// 		logger.Error("Failed to select availaable products from database: " + err.Error())
// 		return c.String(http.StatusInternalServerError, err.Error())
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		product := entities.Product{}
// 		photo := []byte{}
// 		err = rows.Scan(&product.ID,
// 			&product.Name,
// 			&product.Type,
// 			&photo,
// 			&product.Price,
// 			&product.OnSale,
// 			&product.Sale,
// 			&product.Slug,
// 			&product.Position,
// 		)
// 		if err != nil {
// 			logger.Error("Failed to scan products from database: " + err.Error())
// 			return c.String(http.StatusInternalServerError, err.Error())
// 		}
// 		product.Photo = make([]int, len(photo))
// 		for i, v := range photo {
// 			product.Photo[i] = int(v)
// 		}

// 		products = append(products, product)
// 	}
// 	return c.JSON(http.StatusOK, products)
// }

func EditAvailable(c echo.Context, db *sqlx.DB) error {
	product := entities.Product{}
	if err := json.NewDecoder(c.Request().Body).Decode(&product); err != nil {
		logger.Error("Error parsing response body to edit available product: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	_, err := db.Exec(`UPDATE "products" SET "Available"=$1 WHERE "ID"=$2`, product.Available, product.ID)
	if err != nil {
		logger.Error("Failed to edit available product: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "available product edited successfully")
}

func EditBest(c echo.Context, db *sqlx.DB) error {
	product := entities.Product{}
	if err := json.NewDecoder(c.Request().Body).Decode(&product); err != nil {
		logger.Error("Error parsing response body to edit available product: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	_, err := db.Exec(`UPDATE "products" SET "Best"=$1 WHERE "ID"=$2`, product.Best, product.ID)
	if err != nil {
		logger.Error("Failed to edit available product: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "available product edited successfully")
}
