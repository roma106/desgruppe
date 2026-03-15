package queries

import (
	"desgruppe/internal/entities"
	"desgruppe/internal/logger"
	"desgruppe/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

// func FirstListProducts(c echo.Context, db *sqlx.DB) error {
// 	products := []entities.Product{}

// 	stmt, err := db.Preparex(`SELECT "ID","Name","Type","Photo","Price","OnSale","Sale","Slug" FROM "products" ORDER BY "Name" ASC LIMIT 9`)
// 	if err != nil {
// 		logger.Error("Failed to prepare query: " + err.Error())
// 		return c.String(http.StatusInternalServerError, err.Error())
// 	}
// 	defer stmt.Close()

// 	rows, err := stmt.Queryx()
// 	if err != nil {
// 		logger.Error("Failed to select products from database: " + err.Error())
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

// func SecondListProducts(c echo.Context, db *sqlx.DB) error {
// 	products := []entities.Product{}

// 	stmt, err := db.Preparex(`SELECT "ID","Name","Type","Photo","Price","OnSale","Sale","Slug" FROM "products" ORDER BY "Name" ASC OFFSET 9`)
// 	if err != nil {
// 		logger.Error("Failed to prepare query: " + err.Error())
// 		return c.String(http.StatusInternalServerError, err.Error())
// 	}
// 	defer stmt.Close()

// 	rows, err := stmt.Queryx()
// 	if err != nil {
// 		logger.Error("Failed to select products from database: " + err.Error())
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

// func AdminListProducts(c echo.Context, db *sqlx.DB) error {
// 	products := []entities.Product{}

// 	stmt, err := db.Preparex(`SELECT * FROM "products" ORDER BY "Position" ASC`)
// 	if err != nil {
// 		logger.Error("Failed to prepare query: " + err.Error())
// 		return c.String(http.StatusInternalServerError, err.Error())
// 	}
// 	defer stmt.Close()

// 	rows, err := stmt.Queryx()
// 	if err != nil {
// 		logger.Error("Failed to select products from database: " + err.Error())
// 		return c.String(http.StatusInternalServerError, err.Error())
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		product := entities.Product{}
// 		photo := []byte{}
// 		colors := pq.Int64Array{}
// 		err = rows.Scan(&product.ID,
// 			&product.Name,
// 			&product.Type,
// 			&product.SectionID,
// 			&product.ProducerID,
// 			&product.DesignerID,
// 			&product.Size,
// 			&product.Available,
// 			&product.Price,
// 			&product.OnSale,
// 			&product.Sale,
// 			&photo,
// 			&product.Description,
// 			&colors,
// 			&product.Position,
// 			&product.Slug,
// 			&product.FreeForm,
// 		)
// 		if err != nil {
// 			logger.Error("Failed to scan products from database: " + err.Error())
// 			return c.String(http.StatusInternalServerError, err.Error())
// 		}
// 		product.Colors = make([]int, len(colors))
// 		product.Photo = make([]int, len(photo))
// 		go func() {
// 			for i, v := range colors {
// 				product.Colors[i] = int(v)
// 			}
// 		}()
// 		go func() {
// 			for i, v := range photo {
// 				product.Photo[i] = int(v)
// 			}
// 		}()

// 		products = append(products, product)
// 	}
// 	return c.JSON(http.StatusOK, products)
// }

func ListProductsByFilter(c echo.Context, db *sqlx.DB) error {
	products := []entities.Product{}
	sectionIds := c.QueryParam("sections")
	producerIds := c.QueryParam("producers")
	designerIds := c.QueryParam("designers")
	sale := c.QueryParam("sale")
	available := c.QueryParam("available")
	minprice := c.QueryParam("minprice")
	maxprice := c.QueryParam("maxprice")

	query := `SELECT * FROM "products" WHERE 1=1`
	if sectionIds != "" {
		query += fmt.Sprintf(` AND "SectionID" IN (%s)`, sectionIds)
	}
	if producerIds != "" {
		query += fmt.Sprintf(` AND "ProducerID" IN (%s)`, producerIds)
	}
	if designerIds != "" {
		query += fmt.Sprintf(` AND "DesignerID" IN (%s)`, designerIds)
	}
	if sale != "" {
		query += fmt.Sprintf(` AND "OnSale"=True`)
	}
	if minprice != "" {
		query += fmt.Sprintf(` AND "Price" >= %s`, minprice)
	}
	if maxprice != "" {
		query += fmt.Sprintf(` AND "Price" <= %s`, maxprice)
	}
	if available == "true" {
		query += fmt.Sprintf(` AND "Available"=True`)
	}

	stmt, err := db.Preparex(query)
	if err != nil {
		logger.Error("Failed to prepare query: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.Queryx()
	// rows, err := db.Query(query)
	if err != nil {
		logger.Error("Failed to select products from database: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		product := entities.Product{}
		photo := []byte{}
		colors := pq.Int64Array{}
		err = rows.Scan(&product.ID,
			&product.Name,
			&product.Type,
			&product.SectionID,
			&product.ProducerID,
			&product.DesignerID,
			&product.Size,
			&product.Available,
			&product.Price,
			&product.OnSale,
			&product.Sale,
			&photo,
			&product.Description,
			&colors,
			&product.Position,
			&product.Slug,
			&product.FreeForm,
			&product.Best,
		)
		if err != nil {
			logger.Error("Failed to scan products from database: " + err.Error())
			return c.String(http.StatusInternalServerError, err.Error())
		}
		product.Colors = make([]int, len(colors))
		product.Photo = make([]int, len(photo))
		go func() {
			for i, v := range colors {
				product.Colors[i] = int(v)
			}
		}()
		go func() {
			for i, v := range photo {
				product.Photo[i] = int(v)
			}
		}()

		products = append(products, product)
	}
	return c.JSON(http.StatusOK, products)
}

func GetProduct(c echo.Context, db *sqlx.DB) error {
	product := entities.Product{}
	id := c.QueryParam("id")

	rows, err := db.Query(`SELECT * FROM "products" WHERE "ID"=$1`, id)
	if err != nil {
		logger.Error("Failed to select product from database: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		colors := pq.Int64Array{}
		err = rows.Scan(&product.ID,
			&product.Name,
			&product.Type,
			&product.SectionID,
			&product.ProducerID,
			&product.DesignerID,
			&product.Size,
			&product.Available,
			&product.Price,
			&product.OnSale,
			&product.Sale,
			&[]byte{},
			&product.Description,
			&colors,
			&product.Position,
			&product.Slug,
			&product.FreeForm,
			&product.Best)
		if err != nil {
			logger.Error("Failed to scan product from database: " + err.Error())
			return c.String(http.StatusInternalServerError, err.Error())
		}
		for _, v := range colors {
			product.Colors = append(product.Colors, int(v))
		}
	}

	return c.JSON(http.StatusOK, product)
}

func GetProductBySlug(db *sqlx.DB, slug string) (entities.Product, error) {
	product := entities.Product{}
	rows, err := db.Query(`SELECT * FROM "products" WHERE "Slug"=$1`, slug)
	if err != nil {
		logger.Error("Failed to select product from database: " + err.Error())
		return product, err
	}
	defer rows.Close()

	for rows.Next() {
		colors := pq.Int64Array{}
		photo := []byte{}
		err = rows.Scan(&product.ID,
			&product.Name,
			&product.Type,
			&product.SectionID,
			&product.ProducerID,
			&product.DesignerID,
			&product.Size,
			&product.Available,
			&product.Price,
			&product.OnSale,
			&product.Sale,
			&photo,
			&product.Description,
			&colors,
			&product.Position,
			&product.Slug,
			&product.FreeForm,
			&product.Best)
		if err != nil {
			logger.Error("Failed to scan product from database: " + err.Error())
			return product, err
		}
		for _, v := range photo {
			product.Photo = append(product.Photo, int(v))
		}
		for _, v := range colors {
			product.Colors = append(product.Colors, int(v))
		}
	}
	if err != nil {
		logger.Error("Failed to select product from database: " + err.Error())
		return product, err
	}
	return product, nil
}

func AddProduct(c echo.Context, db *sqlx.DB) error {
	addProduct := entities.Product{}

	if err := json.NewDecoder(c.Request().Body).Decode(&addProduct); err != nil {
		logger.Error("Error parsing response body to add product: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	if addProduct.Position == 0 {
		prs, err := utils.ListProducts(db, "any", nil)
		if err != nil {
			return err
		}
		addProduct.Position = len(prs) + 1
	}

	addProduct.Slug = utils.GenerateSlug(addProduct.Name)

	photoBytes := make([]byte, len(addProduct.Photo))
	for i, v := range addProduct.Photo {
		photoBytes[i] = byte(v)
	}

	var err error
	var salePercent any
	if addProduct.OnSale {
		salePercent = addProduct.Sale
	} else {
		salePercent = 0
	}
	var productID int
	err = db.QueryRow(`INSERT INTO "products" 
	("Name", "Type", "SectionID", "ProducerID", "DesignerID",
	"Size", "Available", "Price", "OnSale", "Sale", 
	"Photo", "Description", "Colors", "Position", "Slug", "Free Form") VALUES 
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) RETURNING "ID"`,
		addProduct.Name, addProduct.Type, addProduct.SectionID,
		addProduct.ProducerID, addProduct.DesignerID, addProduct.Size,
		addProduct.Available, addProduct.Price, addProduct.OnSale, salePercent,
		photoBytes, addProduct.Description, pq.Array(addProduct.Colors), addProduct.Position,
		addProduct.Slug, addProduct.FreeForm,
	).Scan(&productID)

	if err != nil {
		logger.Error("Failed to add product to database: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	logger.Info("product added: ", addProduct.Name)
	return c.String(http.StatusCreated, fmt.Sprint(productID))
}

// РЕДАКТИРОВАНИЕ ПРОДУКТА БЕЗ ФОТОГРАФИЙ
func EditProduct(c echo.Context, db *sqlx.DB) error {
	addproduct := entities.Product{}

	if err := json.NewDecoder(c.Request().Body).Decode(&addproduct); err != nil {
		logger.Error("Error parsing response body to edit product: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	var err error
	_, err = db.Exec(`UPDATE "products" SET 
		"Name"=$1, "Type"=$2, "SectionID"=$3, "ProducerID"=$4, "DesignerID"=$5,
	"Size"=$6, "Available"=$7, "Price"=$8, "OnSale"=$9, "Sale"=$10, 
	"Description"=$11, "Colors"=$12, "Position"=$13, "Free Form"=$14, "Slug"=$15 WHERE "ID"=$16`,
		addproduct.Name, addproduct.Type, addproduct.SectionID, addproduct.ProducerID, addproduct.DesignerID,
		addproduct.Size, addproduct.Available, addproduct.Price, addproduct.OnSale, addproduct.Sale, addproduct.Description,
		pq.Array(addproduct.Colors), addproduct.Position, addproduct.FreeForm, addproduct.Slug, addproduct.ID,
	)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			if pqErr.Code == "23505" {
				return c.String(501, "slug repeating")
			}
		}
		logger.Error("Failed to edit product: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	logger.Info("product edited: ", addproduct.Name)
	return c.String(http.StatusOK, "product updated successfully")
}

func DeleteProduct(c echo.Context, db *sqlx.DB) error {
	product := entities.Product{}

	if err := json.NewDecoder(c.Request().Body).Decode(&product); err != nil {
		logger.Error("Error parsing response body to delete product: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	_, err := db.Exec(`DELETE FROM "products" WHERE "ID"=$1`, product.ID)
	if err != nil {
		logger.Error("Failed to delete product from database: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	logger.Info("product deleted: ", product.ID)
	return c.String(http.StatusOK, "product deleted successfully")
}

// func Searchproduct(c echo.Context, db *sqlx.DB) error {
// 	results := []entities.product{}
// 	query := c.QueryParam("q")

// 	rows, err := db.Query(`SELECT * FROM "products" WHERE "Name" COLLATE "C" ILIKE $1`, `%`+query+`%`)
// 	if err != nil {
// 		logger.Error("Failed to search products from database: " + err.Error())
// 		return c.String(http.StatusInternalServerError, err.Error())
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		product := entities.product{}
// 		bytes := []byte{}
// 		var code any
// 		err = rows.Scan(&product.ID, &product.Name, &code, &bytes, &product.Position)
// 		if err != nil {
// 			logger.Error("Failed to scan product from database: " + err.Error())
// 			return c.String(http.StatusInternalServerError, err.Error())
// 		}

// 		if code == nil {
// 			for _, v := range bytes {
// 				product.Picture = append(product.Picture, int(v))
// 			}
// 		} else {
// 			product.Code = code.(string)
// 		}
// 		results = append(results, product)
// 	}

// 	return c.JSON(http.StatusOK, results)
// }

func ClearDeletedValue(db *sqlx.DB, name string, id int) error {
	_, err := db.Exec(fmt.Sprintf(`UPDATE "products" SET "%s"=NULL WHERE "%s"=$1`, name, name), id)
	if err != nil {
		logger.Error("Failed to clear deleted value from database: " + err.Error())
		return err
	}
	return nil
}

func ListRecommendations(db *sqlx.DB, sectionID int, withoutID int) ([]entities.Product, error) {
	products := []entities.Product{}
	q := `SELECT "ID","Name","Type","SectionID","DesignerID","ProducerID","Size","Available","Price","OnSale","Sale", "Position","Slug" FROM "products" WHERE "SectionID"=$1 AND "ID"!=$2`

	stmt, err := db.Preparex(q)
	if err != nil {
		logger.Error("Failed to prepare query: " + err.Error())
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Queryx(sectionID, withoutID)
	if err != nil {
		logger.Error("Failed to select products from database: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		product := entities.Product{}
		colors := pq.Int64Array{}
		err = rows.Scan(&product.ID,
			&product.Name,
			&product.Type,
			&product.SectionID,
			&product.DesignerID,
			&product.ProducerID,
			&product.Size,
			&product.Available,
			&product.Price,
			&product.OnSale,
			&product.Sale,
			&product.Position,
			&product.Slug,
		)
		if err != nil {
			logger.Error("Failed to scan products from database: " + err.Error())
			return nil, err
		}
		product.Colors = make([]int, len(colors))
		for i, v := range colors {
			product.Colors[i] = int(v)
		}
		products = append(products, product)
	}
	return products, nil
}
