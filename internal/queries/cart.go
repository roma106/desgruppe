package queries

import (
	"desgruppe/internal/entities"
	"desgruppe/internal/logger"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func CreateNewCart(c echo.Context, db *sqlx.DB) error {
	var cartID int
	err := db.QueryRow(`INSERT INTO "carts" 
	("ProductIDs") VALUES 
	($1) RETURNING "ID"`, nil).Scan(&cartID)

	if err != nil {
		logger.Error("Failed to find ID for new user cart: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusCreated, fmt.Sprint(cartID))
}

func AddProductToCart(c echo.Context, db *sqlx.DB) error {
	newProduct := entities.CartProduct{}
	if err := json.NewDecoder(c.Request().Body).Decode(&newProduct); err != nil {
		logger.Error("Error parsing response body to add product to cart: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	cart := entities.Cart{}
	rows, err := db.Query(`SELECT * FROM "carts" WHERE "ID"=$1`, newProduct.ID)
	if err != nil {
		logger.Error("Error getting cart with ID: ", newProduct.ID, err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	for rows.Next() {
		productIDs := pq.Int64Array{}
		colorIDs := pq.Int64Array{}
		quantitiesIDs := pq.Int64Array{}
		err = rows.Scan(&cart.ID, &productIDs, &colorIDs, &quantitiesIDs)
		if err != nil {
			logger.Error("Error scanning cart with ID: ", newProduct.ID, err)
			return c.String(http.StatusBadRequest, err.Error())
		}
		for _, v := range productIDs {
			cart.ProductIDs = append(cart.ProductIDs, int(v))
		}
		for _, v := range colorIDs {
			cart.ColorIDs = append(cart.ColorIDs, int(v))
		}
		for _, v := range quantitiesIDs {
			cart.QuantitiesIDs = append(cart.QuantitiesIDs, int(v))
		}
	}
	appending := true
	for i, prID := range cart.ProductIDs {
		if prID == newProduct.ProductID && cart.ColorIDs[i] == newProduct.ColorID {
			appending = false
			cart.QuantitiesIDs[i] += newProduct.QuantitiesID
		}
	}
	if appending {
		cart.ProductIDs = append(cart.ProductIDs, newProduct.ProductID)
		cart.ColorIDs = append(cart.ColorIDs, newProduct.ColorID)
		cart.QuantitiesIDs = append(cart.QuantitiesIDs, newProduct.QuantitiesID)
	}
	// fmt.Println(cart)
	_, err = db.Exec(`UPDATE "carts" SET "ProductIDs"=$1, "ColorIDs"=$2, "QuantitiesIDs"=$3 WHERE "ID"=$4`, pq.Array(cart.ProductIDs), pq.Array(cart.ColorIDs), pq.Array((cart.QuantitiesIDs)), cart.ID)
	if err != nil {
		logger.Error("Error inserting values into cart with ID: ", newProduct.ID, err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(200, "ok")
}

// EDITING  CART

func EditProductQuantity(c echo.Context, db *sqlx.DB) error {
	newProduct := entities.CartProduct{}
	if err := json.NewDecoder(c.Request().Body).Decode(&newProduct); err != nil {
		logger.Error("Error parsing response body to add product to cart: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	_, err := db.Exec(fmt.Sprintf(`UPDATE "carts" SET "QuantitiesIDs"[%d]=%d WHERE "ID"=%d`, newProduct.Index+1, newProduct.QuantitiesID, newProduct.ID))
	if err != nil {
		logger.Error("Error inserting values into cart with ID: ", newProduct.ID, err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(200, "ok")
}

func DeleteProductFromCart(c echo.Context, db *sqlx.DB) error {
	newProduct := entities.CartProduct{}
	if err := json.NewDecoder(c.Request().Body).Decode(&newProduct); err != nil {
		logger.Error("Error parsing response body to add product to cart: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	_, err := db.Exec(`UPDATE "carts" SET "QuantitiesIDs" = "QuantitiesIDs"[1:$1] || "QuantitiesIDs"[$2:array_length("QuantitiesIDs", 1)] WHERE "ID"=$3`, newProduct.Index, newProduct.Index+2, newProduct.ID)
	if err != nil {
		logger.Error("Error deleting values into cart with ID: ", newProduct.ID, err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	_, err = db.Exec(`UPDATE "carts" SET "ColorIDs" = "ColorIDs"[1:$1] || "ColorIDs"[$2:array_length("ColorIDs", 1)] WHERE "ID"=$3`, newProduct.Index, newProduct.Index+2, newProduct.ID)
	if err != nil {
		logger.Error("Error deleting values into cart with ID: ", newProduct.ID, err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	_, err = db.Exec(`UPDATE "carts" SET "ProductIDs" = "ProductIDs"[1:$1] || "ProductIDs"[$2:array_length("ProductIDs", 1)] WHERE "ID"=$3`, newProduct.Index, newProduct.Index+2, newProduct.ID)
	if err != nil {
		logger.Error("Error deleting values into cart with ID: ", newProduct.ID, err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(200, "ok")
}
