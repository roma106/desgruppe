package utils

import (
	"database/sql"
	"desgruppe/internal/entities"
	"desgruppe/internal/logger"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func GenerateSlug(s string) string {
	slug := strings.ToLower(s)

	slug = strings.ReplaceAll(slug, " ", "-")

	re := regexp.MustCompile(`[^a-z0-9-]`)
	slug = re.ReplaceAllString(slug, "")

	return slug
}

func GetProductById(db *sqlx.DB, id string) (entities.Product, error) {
	product := entities.Product{}
	rows, err := db.Query(`SELECT 
	"ID","Name","Type","SectionID","ProducerID","DesignerID","Size","Available", "Price","OnSale","Sale","Description","Colors","Position","Free Form","Slug"
	FROM "products" WHERE "ID"=$1`, id)
	if err != nil {
		logger.Error("Failed to select product from database: " + err.Error())
		return product, err
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
			&product.Description,
			&colors,
			&product.Position,
			&product.FreeForm,
			&product.Slug)
		if err != nil {
			logger.Error("Failed to scan product from database: " + err.Error())
			return product, err
		}
		product.Colors = ConvertPqIntArray(colors)
	}
	if err != nil {
		logger.Error("Failed to select product from database: " + err.Error())
		return product, err
	}
	return product, nil
}

func ListDesigners(db *sqlx.DB) ([]entities.Designer, error) {
	designers := []entities.Designer{}
	rows, err := db.Query(`SELECT "ID","Name","Description","Position","Slug" FROM "designers" ORDER BY "Position" ASC`)
	if err != nil {
		logger.Error("Failed to select designers from database: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		designer := entities.Designer{}
		err = rows.Scan(&designer.ID, &designer.Name, &designer.Description, &designer.Position, &designer.Slug)
		if err != nil {
			logger.Error("Failed to scan designers from database: " + err.Error())
			return nil, err
		}
		designers = append(designers, designer)
	}

	return designers, nil
}

func ListProducers(db *sqlx.DB) ([]entities.Producer, error) {
	producers := []entities.Producer{}
	rows, err := db.Query(`SELECT "ID","Name","Description","Position","Slug" FROM "producers" ORDER BY "Position" ASC`)
	if err != nil {
		logger.Error("Failed to select producers from database: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		producer := entities.Producer{}
		err = rows.Scan(&producer.ID, &producer.Name, &producer.Description, &producer.Position, &producer.Slug)
		if err != nil {
			logger.Error("Failed to scan producers from database: " + err.Error())
			return nil, err
		}
		producers = append(producers, producer)
	}

	return producers, nil
}

func ListSections(db *sqlx.DB) ([]entities.Section, error) {
	sections := []entities.Section{}
	rows, err := db.Query(`SELECT * FROM "sections"`)
	if err != nil {
		logger.Error("Failed to select sections from database: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		section := entities.Section{}
		err = rows.Scan(&section.ID, &section.Name, &section.Type)
		if err != nil {
			logger.Error("Failed to scan sections from database: " + err.Error())
			return nil, err
		}
		sections = append(sections, section)
	}

	return sections, nil
}

func ListColors(db *sqlx.DB) ([]entities.Color, error) {
	colors := []entities.Color{}

	rows, err := db.Query(`SELECT "ID","Name","Code","Position" FROM "colors" ORDER BY "Position" ASC`)
	if err != nil {
		logger.Error("Failed to select colors from database: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		color := entities.Color{}
		var code any
		err = rows.Scan(&color.ID, &color.Name, &code, &color.Position)
		if err != nil {
			logger.Error("Failed to scan colors from database: " + err.Error())
			return nil, err
		}

		if code != nil {
			color.Code = code.(string)
		}
		colors = append(colors, color)
	}

	return colors, nil
}

func ListProducts(db *sqlx.DB, ftype string, sort any) ([]entities.Product, error) {
	products := []entities.Product{}
	q := `SELECT "ID","Name","Type","SectionID","DesignerID","ProducerID","Size","Available","Price","OnSale","Sale","Free Form","Position","Slug","Best" FROM "products"`
	if ftype != "any" {
		q += fmt.Sprintf(` WHERE "Type"='%s'`, ftype)
	}
	switch sort {
	case nil:
		q += ` ORDER BY "Name" ASC`
	case "abc":
		q += ` ORDER BY "Name" ASC`
	case "random":
		q += ` ORDER BY "Position" ASC`
	case "increasing":
		q += ` ORDER BY "Price" ASC`
	case "decreasing":
		q += ` ORDER BY "Price" DESC`
	default:
		q += ` ORDER BY "Name" ASC`
	}

	stmt, err := db.Preparex(q)
	if err != nil {
		logger.Error("Failed to prepare query: " + err.Error())
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Queryx()
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
			&product.FreeForm,
			&product.Position,
			&product.Slug,
			&product.Best,
		)
		if err != nil {
			logger.Error("Failed to scan products from database: " + err.Error())
			return nil, err
		}
		product.Colors = ConvertPqIntArray(colors)
		products = append(products, product)
	}
	return products, nil
}
func ListAvailable(db *sqlx.DB) ([]entities.Product, error) {
	products := []entities.Product{}
	rows, err := db.Query(`SELECT "ID","Name","Type","Price","OnSale","Sale","Slug","Position" FROM "products" WHERE "Available"=true ORDER BY "Position"`)
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
			&product.Price,
			&product.OnSale,
			&product.Sale,
			&product.Slug,
			&product.Position,
		)
		if err != nil {
			logger.Error("Failed to scan products from database: " + err.Error())
			return nil, err
		}
		product.Colors = ConvertPqIntArray(colors)
		products = append(products, product)
	}
	return products, nil
}
func ListBest(db *sqlx.DB) ([]entities.Product, error) {
	products := []entities.Product{}
	rows, err := db.Query(`SELECT "ID","Name","Type","Price","OnSale","Sale","Slug","Position" FROM "products" WHERE "Best"=true ORDER BY "Position"`)
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
			&product.Price,
			&product.OnSale,
			&product.Sale,
			&product.Slug,
			&product.Position,
		)
		if err != nil {
			logger.Error("Failed to scan products from database: " + err.Error())
			return nil, err
		}
		product.Colors = ConvertPqIntArray(colors)
		products = append(products, product)
	}
	return products, nil
}

// func ListSales(db *sqlx.DB) ([]entities.Product, error) {
// 	products := []entities.Product{}
// 	rows, err := db.Query(`SELECT * FROM "sales" ORDER BY "Position"`)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		sale := entities.Sale{}
// 		err = rows.Scan(&sale.ID, &sale.ProductID, &sale.Position)
// 		if err != nil {
// 			logger.Error("Failed to list sale: " + err.Error())
// 			return nil, err
// 		}

// 		r, err := db.Query(`SELECT "ID","Name","Type","Price","OnSale","Sale","Slug","Position" FROM "products" WHERE "ID"=$1`, sale.ProductID)
// 		if err != nil {
// 			logger.Error("Failed to list sale: " + err.Error())
// 			return nil, err
// 		}
// 		defer r.Close()
// 		for r.Next() {
// 			product := entities.Product{}
// 			err = r.Scan(
// 				&product.ID,
// 				&product.Name,
// 				&product.Type,
// 				&product.Price,
// 				&product.OnSale,
// 				&product.Sale,
// 				&product.Slug,
// 				&product.Position,
// 			)
// 			if err != nil {
// 				logger.Error("Failed to add sale: " + err.Error())
// 				return nil, err
// 			}
// 			product.Position = sale.Position
// 			products = append(products, product)
// 		}
// 	}
// 	return products, nil
// }

func GetSettings(db *sqlx.DB) entities.Settings {
	settings := entities.Settings{}

	err := db.Get(&settings, `SELECT * FROM "settings" WHERE "ID"=$1`, 1)

	if err != nil {
		if err == sql.ErrNoRows {
			// Если нет строк, создаем новую запись
			_, err := db.Exec(`INSERT INTO "settings" ("ID", "ExchangeRate", "Email") VALUES ($1, $2, $3)`, 1, 100, ".")
			if err != nil {
				logger.Error("Failed to insert default settings into database: " + err.Error())
				return entities.Settings{}
			}
			settings.ID = 1
			settings.ExchangeRate = 1
			logger.Info("No settings found, inserting default settings with ExchangeRate = 1")
		} else {
			logger.Error("Failed to get settings from database: " + err.Error())
			return entities.Settings{}
		}
	}

	return settings
}

func GetOrderByID(db *sqlx.DB, id string) (entities.Order, error) {
	order := entities.Order{}
	rows, err := db.Query(`SELECT * FROM "orders" WHERE "ID"=$1`, id)
	if err != nil {
		logger.Error("Failed to select order from database: " + err.Error())
		return order, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&order.ID, &order.Date, &order.Name, &order.Phone, &order.Email, &order.Comment, &order.CartID, &order.Seen)
		if err != nil {
			logger.Error("Failed to scan order from database: " + err.Error())
			return order, err
		}
	}

	return order, nil
}

func GetImgBlob(c echo.Context, db *sqlx.DB, tableName string, columnName string, entityID string) error {
	bytes := []byte{}
	err := db.Get(&bytes, fmt.Sprintf(`SELECT "%s" FROM "%s" WHERE "ID"=%s`, columnName, tableName, entityID))
	if err != nil {
		logger.Error("Failed to select img from database: " + err.Error())
		return err
	}
	return c.Blob(200, "image/png", bytes)
}

func FilterProducts(products []entities.Product, desIDs []string, prodIDs []string, secIDs []string, available string) []entities.Product {
	var result []entities.Product
	filterFunc := func(filter []string, v int) bool {
		if len(filter) == 0 || (len(filter) == 1 && filter[0] == "") {
			return true
		}
		return slices.Contains(filter, strconv.Itoa(v))
	}
	for _, pr := range products {
		if !filterFunc(desIDs, pr.DesignerID) {
			continue
		}
		if !filterFunc(prodIDs, pr.ProducerID) {
			continue
		}
		if !filterFunc(secIDs, pr.SectionID) {
			continue
		}
		if available != "" {
			av, _ := strconv.ParseBool(available)
			if pr.Available != av {
				continue
			}
		}
		result = append(result, pr)
	}
	return result
}

func ConvertPqIntArray(a pq.Int64Array) []int {
	res := []int{}
	for _, v := range a {
		res = append(res, int(v))
	}
	return res
}
