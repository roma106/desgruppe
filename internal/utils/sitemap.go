package utils

import (
	"bufio"
	"desgruppe/internal/entities"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

func GetSitemap(db *sqlx.DB) (string, error) {
	file, err := os.Open("internal/utils/sitemap-template.xml")
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return "", err
	}
	defer file.Close()

	sm := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sm += scanner.Text() + "\n"
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка чтения файла:", err)
	}

	productsCh := make(chan []entities.Product)
	go func() {
		products, err := ListProducts(db, "any", nil)
		if err != nil {
			return
		}
		productsCh <- products
	}()
	designers, err := ListDesigners(db)
	if err != nil {
		return "", err
	}
	for _, des := range designers {
		sm += fmt.Sprintf("<url><loc>https://desgruppe.shop/designer/%s</loc><priority>0.50</priority></url>", des.Slug)
	}
	producers, err := ListProducers(db)
	if err != nil {
		return "", err
	}
	for _, pr := range producers {
		sm += fmt.Sprintf("<url><loc>https://desgruppe.shop/producer/%s</loc><priority>0.50</priority></url>", pr.Slug)
	}
	products := <-productsCh
	for _, pr := range products {
		sm += fmt.Sprintf("<url><loc>https://desgruppe.shop/product/%s</loc><priority>0.58</priority></url>", pr.Slug)
	}

	sm += fmt.Sprint(`</urlset>` + "\n")
	return sm, nil
}
