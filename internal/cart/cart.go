package cart

import (
	"desgruppe/internal/entities"
	"desgruppe/internal/logger"
	"desgruppe/internal/queries"
	"desgruppe/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func CartMiddleware(c echo.Context, db *sqlx.DB) {
	co, _ := c.Cookie("cartID")
	if co == nil || co.Value == "" {
		c.SetCookie(&http.Cookie{
			Name:    "cartID",
			Value:   CreateNewCart(c, db),
			Path:    "/",
			Expires: time.Now().Add(time.Hour * 24 * 30),
		})
		return
	}
}

func CreateNewCart(c echo.Context, db *sqlx.DB) string {
	var cartID int
	err := db.QueryRow(`INSERT INTO "carts" 
	("ProductIDs") VALUES 
	($1) RETURNING "ID"`, nil).Scan(&cartID)

	if err != nil {
		logger.Error("Failed to find ID for new user cart: " + err.Error())
		return ""
	}

	return strconv.Itoa(cartID)
}

func GetUserCart(c echo.Context, db *sqlx.DB) ([]entities.Product, []string, []int) {
	products := []entities.Product{}
	colors := []string{}
	qtys := []int{}

	cart := entities.Cart{}
	cartIdC, err := c.Cookie("cartID")
	if err != nil {
		if err != http.ErrNoCookie {
			logger.Error(err.Error())
		}
		return nil, nil, nil
	}
	cartId := cartIdC.Value
	productIDs := pq.Int64Array{}
	colorIDs := pq.Int64Array{}
	quantitiesIDs := pq.Int64Array{}
	err = db.QueryRowx(`SELECT * FROM "carts" WHERE "ID"=$1`, cartId).Scan(&cart.ID, &productIDs, &colorIDs, &quantitiesIDs)
	if err != nil {
		logger.Error("Error getting cart with ID: ", cartId, err)
		return nil, nil, nil
	}
	qtys = utils.ConvertPqIntArray(quantitiesIDs)

	for _, id := range productIDs {
		product, err := utils.GetProductById(db, strconv.Itoa(int(id)))
		if err != nil {
			logger.Error("Error getting product for cart with ID: ", id, err)
			continue
		}
		products = append(products, product)
	}
	colorNodes, err := queries.ListColorsByIds(db, utils.ConvertPqIntArray(colorIDs))
	if err != nil {
		logger.Error("Error getting colors for cart with ID: ", cartId, err)
		return nil, nil, nil
	}
	for _, c := range colorNodes {
		colors = append(colors, c.Name)
	}
	return products, colors, qtys
}

func GetUserCartByID(db *sqlx.DB, id string) ([]entities.Product, []string, []int) {
	products := []entities.Product{}
	colors := []string{}
	qtys := []int{}

	cart := entities.Cart{}
	cartId := id
	rows, err := db.Query(`SELECT * FROM "carts" WHERE "ID"=$1`, cartId)
	// err := db.QueryRowx(`SELECT * FROM "carts" WHERE "ID"=$1`, cartId).Scan(&cart.ID, &productIDs, &colorIDs, &quantitiesIDs)
	if err != nil {
		logger.Error("Error getting cart with ID: ", cartId, err)
		return nil, nil, nil
	}
	for rows.Next() {
		productIDs := pq.Int64Array{}
		colorIDs := pq.Int64Array{}
		quantitiesIDs := pq.Int64Array{}
		err = rows.Scan(&cart.ID, &productIDs, &colorIDs, &quantitiesIDs)
		if err != nil {
			logger.Error("Error getting cart with ID: ", cartId, err)
			return nil, nil, nil
		}
		qtys = utils.ConvertPqIntArray(quantitiesIDs)

		for _, id := range productIDs {
			product, err := utils.GetProductById(db, strconv.Itoa(int(id)))
			if err != nil {
				logger.Error("Error getting product for cart with ID: ", id, err)
				continue
			}
			products = append(products, product)
		}
		colorNodes, err := queries.ListColorsByIds(db, utils.ConvertPqIntArray(colorIDs))
		if err != nil {
			logger.Error("Error getting colors for cart with ID: ", cartId, err)
			return nil, nil, nil
		}
		for _, c := range colorNodes {
			colors = append(colors, c.Name)
		}
	}
	return products, colors, qtys
}
