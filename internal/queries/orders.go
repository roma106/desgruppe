package queries

import (
	"desgruppe/internal/entities"
	"desgruppe/internal/logger"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

// func ListOrders(c echo.Context, db *sqlx.DB) error {
// 	orders := []entities.Order{}

// 	rows, err := db.Query(`SELECT * FROM "orders" ORDER BY "ID" DESC`)
// 	if err != nil {
// 		logger.Error("Failed to select orders from database: " + err.Error())
// 		return c.String(http.StatusInternalServerError, err.Error())
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		order := entities.Order{}
// 		err = rows.Scan(&order.ID, &order.Date, &order.Name, &order.Phone, &order.Email, &order.Comment, &order.CartID, &order.Seen)
// 		if err != nil {
// 			logger.Error("Failed to scan orders from database: " + err.Error())
// 			return c.String(http.StatusInternalServerError, err.Error())
// 		}
// 		orders = append(orders, order)
// 	}

// 	return c.JSON(http.StatusOK, orders)
// }

func GetOrder(c echo.Context, db *sqlx.DB) error {
	order := entities.Order{}
	id := c.QueryParam("id")

	rows, err := db.Query(`SELECT * FROM "orders" WHERE "ID"=$1`, id)
	if err != nil {
		logger.Error("Failed to select order from database: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&order.ID, &order.Date, &order.Name, &order.Phone, &order.Email, &order.Comment, &order.CartID, &order.Seen)
		if err != nil {
			logger.Error("Failed to scan order from database: " + err.Error())
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, order)
}

func AddOrder(c echo.Context, db *sqlx.DB) error {
	addorder := entities.Order{}
	if err := json.NewDecoder(c.Request().Body).Decode(&addorder); err != nil {
		logger.Error("Error parsing response body to add order: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	addorder.Date = time.Now()
	addorder.Seen = false

	err := db.QueryRowx(`INSERT INTO "orders" ("Date", "Name", "Phone", "Email", "Comment", "CartID", "Seen") VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING "ID"`,
		addorder.Date, addorder.Name, addorder.Phone, addorder.Email, addorder.Comment, addorder.CartID, addorder.Seen).Scan(&addorder.ID)
	if err != nil {
		logger.Error("Failed to add order to database: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	c.SetCookie(&http.Cookie{
		Name:    "cartID",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	})
	logger.Info("order added: ", addorder.Name)
	c.Set("orderID", addorder.ID)
	return c.String(http.StatusCreated, "order added successfully")
}

func DeleteOrder(c echo.Context, db *sqlx.DB) error {
	order := entities.Order{}

	if err := json.NewDecoder(c.Request().Body).Decode(&order); err != nil {
		logger.Error("Error parsing response body to delete order: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	var cartID string
	err := db.QueryRow(`DELETE FROM "orders" WHERE "ID"=$1 RETURNING "CartID"`, order.ID).Scan(&cartID)
	if err != nil {
		logger.Error("Failed to delete order from database: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	_, err = db.Exec(`DELETE FROM "carts" WHERE "ID"=$1`, cartID)
	if err != nil {
		logger.Error("Failed to delete cart from database: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	logger.Info("order deleted: ", order.ID)
	return c.String(http.StatusOK, "order deleted successfully")
}

func OrderSeen(c echo.Context, db *sqlx.DB) error {
	order := entities.Order{}

	if err := json.NewDecoder(c.Request().Body).Decode(&order); err != nil {
		logger.Error("Error parsing response body to edit order: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	_, err := db.Exec(`UPDATE "orders" SET "Seen"=$1 WHERE "ID"=$2`, true, order.ID)
	if err != nil {
		logger.Error("Failed to edit order to database: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "Order updated successfully")
}

// func Searchorder(c echo.Context, db *sqlx.DB) error {
// 	results := []entities.order{}
// 	query := c.QueryParam("q")

// 	rows, err := db.Query(`SELECT * FROM "orders" WHERE "Name" COLLATE "C" ILIKE $1`, `%`+query+`%`)
// 	if err != nil {
// 		logger.Error("Failed to search orders from database: " + err.Error())
// 		return c.String(http.StatusInternalServerError, err.Error())
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		order := entities.order{}
// 		bytes := []byte{}
// 		var code any
// 		err = rows.Scan(&order.ID, &order.Name, &code, &bytes, &order.Position)
// 		if err != nil {
// 			logger.Error("Failed to scan order from database: " + err.Error())
// 			return c.String(http.StatusInternalServerError, err.Error())
// 		}

// 		if code == nil {
// 			for _, v := range bytes {
// 				order.Picture = append(order.Picture, int(v))
// 			}
// 		} else {
// 			order.Code = code.(string)
// 		}
// 		results = append(results, order)
// 	}

// 	return c.JSON(http.StatusOK, results)
// }

func ListOrders(db *sqlx.DB) ([]entities.Order, error) {
	orders := []entities.Order{}
	// err := db.Get(&orders, `SELECT * FROM "orders" ORDER BY "ID" DESC`)
	// if err != nil {
	// 	logger.Error("Failed to select orders from database: " + err.Error())
	// 	return nil, err
	// }

	rows, err := db.Query(`SELECT * FROM "orders" ORDER BY "ID" DESC`)
	if err != nil {
		logger.Error("Failed to select orders from database: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		order := entities.Order{}
		err = rows.Scan(&order.ID, &order.Date, &order.Name, &order.Phone, &order.Email, &order.Comment, &order.CartID, &order.Seen)
		if err != nil {
			logger.Error("Failed to scan orders from database: " + err.Error())
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}
