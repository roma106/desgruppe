package queries

import (
	"desgruppe/internal/entities"
	"desgruppe/internal/logger"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

// PHOTOS

func ListPhotos(c echo.Context, db *sqlx.DB) error {
	productID := c.QueryParam("id")
	photos, err := ListPhotosByID(db, productID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, photos)
}

func ListPhotosByID(db *sqlx.DB, id string) ([]entities.ProductPhoto, error) {
	photos := []entities.ProductPhoto{}
	rows, err := db.Query(fmt.Sprintf(`SELECT * FROM "products_gallery" WHERE "ProductID"=%s ORDER BY "Position"`, id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		photo := entities.ProductPhoto{}
		photoBytes := []byte{}
		err = rows.Scan(&photo.ID, &photo.ProductID, &photoBytes, &photo.Position)
		if err != nil {
			return nil, err
		}

		for _, v := range photoBytes {
			photo.Photo = append(photo.Photo, int(v))
		}

		photos = append(photos, photo)
	}
	return photos, nil
}

func AddPhoto(c echo.Context, db *sqlx.DB) error {
	addphoto := entities.ProductPhoto{}
	if err := json.NewDecoder(c.Request().Body).Decode(&addphoto); err != nil {
		logger.Error("Error parsing response body to add photo: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	photoBytes := make([]byte, len(addphoto.Photo))
	for i, v := range addphoto.Photo {
		photoBytes[i] = byte(v)
	}
	_, err := db.Exec(`INSERT INTO "products_gallery" ("ProductID", "Photo", "Position") VALUES ($1, $2, $3)`, addphoto.ProductID, photoBytes, addphoto.Position)
	if err != nil {
		logger.Error("Failed to add photo: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	logger.Info("Photo added to gallery. ID: ", addphoto.ID)
	return c.String(http.StatusOK, "photo added successfully")
}

func EditPhoto(c echo.Context, db *sqlx.DB) error {
	addphoto := entities.ProductPhoto{}
	if err := json.NewDecoder(c.Request().Body).Decode(&addphoto); err != nil {
		logger.Error("Error parsing response body to edit photo: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	_, err := db.Exec(`UPDATE "products_gallery" SET "Position"=$1 WHERE "ID"=$2`, addphoto.Position, addphoto.ID)
	if err != nil {
		logger.Error("Failed to edit photo: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	logger.Info("Photo edited. ID: ", addphoto.ID)
	return c.String(http.StatusOK, "photo edited successfully")
}

func DeletePhoto(c echo.Context, db *sqlx.DB) error {
	productID := c.QueryParam("id")
	_, err := db.Exec(`DELETE FROM "products_gallery" WHERE "ID"=$1`, productID)
	if err != nil {
		logger.Error("Failed to delete photo: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	logger.Info("Photo for productdeleted. ID: ", productID)
	return c.String(http.StatusOK, "photo deleted successfully")
}

func DeletePhotosForProduct(c echo.Context, db *sqlx.DB) error {
	productID := c.QueryParam("id")
	_, err := db.Exec(`DELETE FROM "products_gallery" WHERE "ProductID"=$1`, productID)
	if err != nil {
		logger.Error("Failed to delete photos for product: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "photos deleted successfully")
}

// MAIN PHOTO

func GetMainPhoto(c echo.Context, db *sqlx.DB) error {
	productID := c.QueryParam("id")
	product := entities.Product{}
	photoBytes := []byte{}
	err := db.Get(&photoBytes, `SELECT "Photo" FROM "products" WHERE "ID"=$1`, productID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	for _, v := range photoBytes {
		product.Photo = append(product.Photo, int(v))
	}
	return c.JSON(http.StatusOK, product.Photo)
}

func EditMainPhoto(c echo.Context, db *sqlx.DB) error {
	addproduct := entities.Product{}

	if err := json.NewDecoder(c.Request().Body).Decode(&addproduct); err != nil {
		logger.Error("Error parsing response body to edit product: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	photoBytes := make([]byte, len(addproduct.Photo))
	for i, v := range addproduct.Photo {
		photoBytes[i] = byte(v)
	}

	var err error
	_, err = db.Exec(`UPDATE "products" SET "Photo"=$1 WHERE "ID"=$2`,
		photoBytes, addproduct.ID,
	)

	if err != nil {
		logger.Error("Failed to edit product: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	logger.Info("Main photo of the product edited. ID: ", addproduct.ID)
	return c.String(http.StatusOK, "product updated successfully")
}
