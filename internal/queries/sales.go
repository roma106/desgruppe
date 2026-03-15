package queries

import (
	"desgruppe/internal/entities"
	"desgruppe/internal/logger"
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func ListSales(c echo.Context, db *sqlx.DB) error {
	products := []entities.Product{}
	rows, err := db.Query(`SELECT * FROM "sales" ORDER BY "Position"`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		sale := entities.Sale{}
		err = rows.Scan(&sale.ID, &sale.ProductID, &sale.Position)
		if err != nil {
			logger.Error("Failed to add sale: " + err.Error())
			return c.String(http.StatusInternalServerError, err.Error())
		}

		r, err := db.Query(`SELECT "ID","Name","Photo","Price","OnSale","Sale","Slug" FROM "products" WHERE "ID"=$1`, sale.ProductID)
		if err != nil {
			logger.Error("Failed to add sale: " + err.Error())
			return c.String(http.StatusInternalServerError, err.Error())
		}
		defer r.Close()
		for r.Next() {
			product := entities.Product{}
			photo := []byte{}
			err = r.Scan(&product.ID, &product.Name, &photo, &product.Price, &product.OnSale, &product.Sale, &product.Slug)
			if err != nil {
				logger.Error("Failed to add sale: " + err.Error())
				return c.String(http.StatusInternalServerError, err.Error())
			}
			for _, v := range photo {
				product.Photo = append(product.Photo, int(v))
			}
			product.Position = sale.Position
			products = append(products, product)
		}
	}
	return c.JSON(200, products)
}

func AddSale(c echo.Context, db *sqlx.DB) error {
	addsale := entities.Sale{}
	if err := json.NewDecoder(c.Request().Body).Decode(&addsale); err != nil {
		logger.Error("Error parsing response body to add sale: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	_, err := db.Exec(`INSERT INTO "sales" ("ProductID", "Position") VALUES ($1, $2)`, addsale.ProductID, addsale.Position)
	if err != nil {
		logger.Error("Failed to add sale: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "sale added successfully")
}

func EditSale(c echo.Context, db *sqlx.DB) error {
	addsale := entities.Sale{}
	if err := json.NewDecoder(c.Request().Body).Decode(&addsale); err != nil {
		logger.Error("Error parsing response body to edit sale: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	_, err := db.Exec(`UPDATE "sales" SET "Position"=$1 WHERE "ProductID"=$2`, addsale.Position, addsale.ProductID)
	if err != nil {
		logger.Error("Failed to edit sale: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "sale edited successfully")
}

func DeleteSale(c echo.Context, db *sqlx.DB) error {
	id := c.QueryParam("id")
	_, err := db.Exec(`DELETE FROM "sales" WHERE "ProductID"=$1`, id)
	if err != nil {
		logger.Error("Failed to delete sale: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "sale deleted successfully")
}
