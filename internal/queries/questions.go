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

func ListQuestions(c echo.Context, db *sqlx.DB) error {
	questions := []entities.Question{}

	rows, err := db.Query(`SELECT * FROM "questions"`)
	if err != nil {
		logger.Error("Failed to select questions from database: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		question := entities.Question{}
		err = rows.Scan(&question.ID, &question.Date, &question.Name, &question.Phone, &question.Email, &question.Message, &question.Seen)
		if err != nil {
			logger.Error("Failed to scan questions from database: " + err.Error())
			return c.String(http.StatusInternalServerError, err.Error())
		}
		questions = append(questions, question)
	}

	return c.JSON(http.StatusOK, questions)
}

func GetQuestion(c echo.Context, db *sqlx.DB) error {
	question := entities.Question{}
	id := c.QueryParam("id")

	rows, err := db.Query(`SELECT * FROM "questions" WHERE "ID"=$1`, id)
	if err != nil {
		logger.Error("Failed to select question from database: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&question.ID, &question.Date, &question.Name, &question.Phone, &question.Email, &question.Message, &question.Seen)
		if err != nil {
			logger.Error("Failed to scan question from database: " + err.Error())
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, question)
}

func AddQuestion(c echo.Context, db *sqlx.DB) error {
	addquestion := entities.Question{}

	if err := json.NewDecoder(c.Request().Body).Decode(&addquestion); err != nil {
		logger.Error("Error parsing response body to add question: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	addquestion.Date = time.Now()
	addquestion.Seen = false

	_, err := db.Exec(`INSERT INTO "questions" ("Date", "Name", "Phone", "Email", "Message", "Seen") VALUES ($1, $2, $3, $4, $5)`,
		addquestion.Date, addquestion.Name, addquestion.Phone, addquestion.Email, addquestion.Message, addquestion.Seen)
	if err != nil {
		logger.Error("Failed to add question to database: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	logger.Info("question added: ", addquestion.Name)
	return c.String(http.StatusCreated, "question added successfully")
}

func DeleteQuestion(c echo.Context, db *sqlx.DB) error {
	question := entities.Question{}

	if err := json.NewDecoder(c.Request().Body).Decode(&question); err != nil {
		logger.Error("Error parsing response body to delete question: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	_, err := db.Exec(`DELETE FROM "questions" WHERE "ID"=$1`, question.ID)
	if err != nil {
		logger.Error("Failed to delete question from database: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	logger.Info("question deleted: ", question.ID)
	return c.String(http.StatusOK, "question deleted successfully")
}

func QuestionSeen(c echo.Context, db *sqlx.DB) error {
	question := entities.Question{}

	if err := json.NewDecoder(c.Request().Body).Decode(&question); err != nil {
		logger.Error("Error parsing response body to edit question: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	_, err := db.Exec(`UPDATE "questions" SET "Seen"=$1 WHERE "ID"=$2`, question.Seen, question.ID)
	if err != nil {
		logger.Error("Failed to edit question to database: " + err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "question updated successfully")
}

// func Searchquestion(c echo.Context, db *sqlx.DB) error {
// 	results := []entities.question{}
// 	query := c.QueryParam("q")

// 	rows, err := db.Query(`SELECT * FROM "questions" WHERE "Name" COLLATE "C" ILIKE $1`, `%`+query+`%`)
// 	if err != nil {
// 		logger.Error("Failed to search questions from database: " + err.Error())
// 		return c.String(http.StatusInternalServerError, err.Error())
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		question := entities.question{}
// 		bytes := []byte{}
// 		var code any
// 		err = rows.Scan(&question.ID, &question.Name, &code, &bytes, &question.Position)
// 		if err != nil {
// 			logger.Error("Failed to scan question from database: " + err.Error())
// 			return c.String(http.StatusInternalServerError, err.Error())
// 		}

// 		if code == nil {
// 			for _, v := range bytes {
// 				question.Picture = append(question.Picture, int(v))
// 			}
// 		} else {
// 			question.Code = code.(string)
// 		}
// 		results = append(results, question)
// 	}

// 	return c.JSON(http.StatusOK, results)
// }
