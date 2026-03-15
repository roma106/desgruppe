package api

import (
	"bytes"
	"desgruppe/internal/cart"
	"desgruppe/internal/entities"
	"desgruppe/internal/queries"
	"desgruppe/internal/templ/views/layout"
	pagesAdmin "desgruppe/internal/templ/views/pages/admin"

	// pages "desgruppe/internal/templ/views/pages/user"
	pagesUser "desgruppe/internal/templ/views/pages/user"
	"desgruppe/internal/utils"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

// USER PAGES

func RenderIndexPage(c echo.Context, db *sqlx.DB) error {
	available, err := utils.ListAvailable(db)
	if err != nil {
		return err
	}
	best, err := utils.ListBest(db)
	if err != nil {
		return err
	}
	settings := utils.GetSettings(db)
	cartProducts, cartColors, cartQtys := cart.GetUserCart(c, db)
	return c.Render(200, "index", pagesUser.IndexPage(pagesUser.IndexPageProps{
		Available: available,
		Best:      best,
		ExRate:    settings.ExchangeRate,
		Cart:      layout.UserCart{cartProducts, cartColors, cartQtys, settings.ExchangeRate},
	}))
}
func RenderAvailablePage(c echo.Context, db *sqlx.DB) error {
	available, err := utils.ListAvailable(db)
	if err != nil {
		return err
	}
	cartProducts, cartColors, cartQtys := cart.GetUserCart(c, db)
	settings := utils.GetSettings(db)
	return c.Render(200, "available", pagesUser.AvailablePage(pagesUser.AvailablePageProps{
		Products: available,
		ExRate:   settings.ExchangeRate,
		Cart:     layout.UserCart{cartProducts, cartColors, cartQtys, settings.ExchangeRate},
	}))
}
func RenderFurniturePage(c echo.Context, db *sqlx.DB, ftype string, pageTitle string) error {
	productsCh := make(chan []entities.Product)
	go func() {
		products, err := utils.ListProducts(db, ftype, nil)
		if err != nil {
			return
		}
		productsCh <- products
	}()
	designers, err := utils.ListDesigners(db)
	if err != nil {
		return err
	}
	producers, err := utils.ListProducers(db)
	if err != nil {
		return err
	}
	sections, err := utils.ListSections(db)
	if err != nil {
		return err
	}
	products := <-productsCh

	minPrice := 999999999999
	maxPrice := 0
	for _, pr := range products {
		if minPrice > int(pr.Price) && pr.Type == ftype {
			minPrice = int(pr.Price)
		}
		if maxPrice < int(pr.Price) && pr.Type == ftype {
			maxPrice = int(pr.Price)
		}
	}
	settings := utils.GetSettings(db)
	cartProducts, cartColors, cartQtys := cart.GetUserCart(c, db)
	return c.Render(200, ftype, pagesUser.FurniturePage(pagesUser.FurniturePageProps{
		Products:    products,
		Designers:   designers,
		Producers:   producers,
		Sections:    sections,
		MaxPrice:    maxPrice,
		MinPrice:    minPrice,
		ExRate:      settings.ExchangeRate,
		Type:        ftype,
		PageTitle:   pageTitle,
		QueryParams: pagesUser.QueryParams{},
		Cart:        layout.UserCart{cartProducts, cartColors, cartQtys, settings.ExchangeRate},
	}))
}

func RenderFilterFurniturePage(c echo.Context, db *sqlx.DB, ftype string, pageTitle string) error {
	designerIds := strings.Split(c.QueryParam("designers"), ",")
	producerIds := strings.Split(c.QueryParam("producers"), ",")
	sectionIds := strings.Split(c.QueryParam("sections"), ",")
	available := c.QueryParam("available")
	sort := c.QueryParam("sort")

	productsCh := make(chan []entities.Product)
	go func() {
		products, err := utils.ListProducts(db, ftype, sort)
		if err != nil {
			return
		}
		productsCh <- products
	}()
	designers, err := utils.ListDesigners(db)
	if err != nil {
		return err
	}
	producers, err := utils.ListProducers(db)
	if err != nil {
		return err
	}
	sections, err := utils.ListSections(db)
	if err != nil {
		return err
	}
	products := <-productsCh
	products = utils.FilterProducts(products, designerIds, producerIds, sectionIds, available)

	minPrice := 999999999999
	maxPrice := 0
	for _, pr := range products {
		if minPrice > int(pr.Price) && pr.Type == ftype {
			minPrice = int(pr.Price)
		}
		if maxPrice < int(pr.Price) && pr.Type == ftype {
			maxPrice = int(pr.Price)
		}
	}
	settings := utils.GetSettings(db)
	cartProducts, cartColors, cartQtys := cart.GetUserCart(c, db)
	return c.Render(200, ftype+"filter", pagesUser.FurniturePage(pagesUser.FurniturePageProps{
		Products:  products,
		Designers: designers,
		Producers: producers,
		Sections:  sections,
		MaxPrice:  maxPrice,
		MinPrice:  minPrice,
		ExRate:    settings.ExchangeRate,
		Type:      ftype,
		PageTitle: pageTitle,
		QueryParams: pagesUser.QueryParams{
			DesignerIds: designerIds,
			ProducerIds: producerIds,
			SectionIds:  sectionIds,
			Available:   available,
			Sort:        sort,
		},
		Cart: layout.UserCart{cartProducts, cartColors, cartQtys, settings.ExchangeRate},
	}))
}

func RenderProductPage(c echo.Context, db *sqlx.DB) error {
	slug := c.Param("slug")
	product, err := queries.GetProductBySlug(db, slug)
	designer, err := queries.GetDesignerById(db, product.DesignerID)
	if err != nil {
		return err
	}
	producers, err := queries.GetProducerById(db, product.ProducerID)
	if err != nil {
		return err
	}
	colors, err := queries.ListColorsByIds(db, product.Colors)
	if err != nil {
		return err
	}
	recs, err := queries.ListRecommendations(db, product.SectionID, product.ID)
	if err != nil {
		return err
	}
	photos, err := queries.ListPhotosByID(db, strconv.Itoa(product.ID))
	if err != nil {
		return err
	}
	settings := utils.GetSettings(db)
	cartProducts, cartColors, cartQtys := cart.GetUserCart(c, db)
	return c.Render(200, slug, pagesUser.ProductPage(pagesUser.ProductPageProps{
		Product:         product,
		Designer:        designer,
		Producer:        producers,
		Colors:          colors,
		ExRate:          settings.ExchangeRate,
		Photos:          photos,
		Recommendations: recs,
		Cart:            layout.UserCart{cartProducts, cartColors, cartQtys, settings.ExchangeRate},
	}))
}

func RenderDesignersPage(c echo.Context, db *sqlx.DB) error {
	designers, err := utils.ListDesigners(db)
	if err != nil {
		return err
	}
	settings := utils.GetSettings(db)
	cartProducts, cartColors, cartQtys := cart.GetUserCart(c, db)
	return c.Render(200, "designers", pagesUser.DesignersPage(pagesUser.DesignersPageProps{
		Designers: designers,
		Cart:      layout.UserCart{cartProducts, cartColors, cartQtys, settings.ExchangeRate},
	}))
}

func RenderCartPage(c echo.Context, db *sqlx.DB) error {
	settings := utils.GetSettings(db)
	cartProducts, cartColors, cartQtys := cart.GetUserCart(c, db)
	return c.Render(200, "cart", pagesUser.CartPage(pagesUser.CartPageProps{
		ExRate: settings.ExchangeRate,
		Cart:   layout.UserCart{cartProducts, cartColors, cartQtys, settings.ExchangeRate},
	}))
}

// ADMIN PAGES

func RenderProductsPage(c echo.Context, db *sqlx.DB) error {
	products, err := utils.ListProducts(db, "any", "random")
	if err != nil {
		return err
	}
	designers, err := utils.ListDesigners(db)
	if err != nil {
		return err
	}
	producers, err := utils.ListProducers(db)
	if err != nil {
		return err
	}
	sections, err := utils.ListSections(db)
	if err != nil {
		return err
	}
	colors, err := utils.ListColors(db)
	if err != nil {
		return err
	}
	return c.Render(200, "productsPage", pagesAdmin.ProductsPage(pagesAdmin.ProductsPageProps{
		Products:  products,
		Designers: designers,
		Producers: producers,
		Sections:  sections,
		Colors:    colors,
	}))
}

func RenderOrdersPage(c echo.Context, db *sqlx.DB) error {
	orders, err := queries.ListOrders(db)
	if err != nil {
		return err
	}
	return c.Render(200, "orders", pagesAdmin.OrdersPage(pagesAdmin.OrdersPageProps{
		Orders: orders,
	}))
}

func RenderOrderEditPage(c echo.Context, db *sqlx.DB) error {
	order, err := utils.GetOrderByID(db, c.QueryParam("id"))
	if err != nil {
		return err
	}
	cartProducts, cartColors, cartQtys := cart.GetUserCartByID(db, order.CartID)

	settings := utils.GetSettings(db)
	return c.Render(200, "order", pagesAdmin.OrderEditPage(pagesAdmin.OrderEditPageProps{
		Order:    order,
		Products: cartProducts,
		Colors:   cartColors,
		Qtys:     cartQtys,
		ExRate:   settings.ExchangeRate,
	}))
}

func RenderAdminAvailablePage(c echo.Context, db *sqlx.DB) error {
	products, err := utils.ListProducts(db, "any", "random")
	if err != nil {
		return err
	}
	return c.Render(200, "available", pagesAdmin.AvailablePage(pagesAdmin.AvailablePageProps{
		Products: products,
	}))
}

func RenderAdminBestPage(c echo.Context, db *sqlx.DB) error {
	products, err := utils.ListProducts(db, "any", "random")
	if err != nil {
		return err
	}
	return c.Render(200, "best", pagesAdmin.BestPage(pagesAdmin.BestPageProps{
		Products: products,
	}))
}

func RenderSettingsPage(c echo.Context, db *sqlx.DB) error {
	settings := utils.GetSettings(db)
	return c.Render(200, "settings", pagesAdmin.SettingsPage(pagesAdmin.SettingsPageProps{
		Settings: settings,
	}))
}

func RenderMailPage(c echo.Context, db *sqlx.DB) string {
	buf := bytes.Buffer{}
	order, _ := utils.GetOrderByID(db, strconv.Itoa(c.Get("orderID").(int)))

	cartProducts, cartColors, cartQtys := cart.GetUserCartByID(db, order.CartID)
	settings := utils.GetSettings(db)
	t := pagesAdmin.NotificationMail(pagesAdmin.NotificationMailProps{
		Order: pagesAdmin.Order{
			Name:    order.Name,
			Date:    order.Date.Format("02.01.2006"),
			Phone:   order.Phone,
			Email:   order.Email,
			Comment: order.Comment,
		},
		Products: cartProducts,
		Colors:   cartColors,
		Qtys:     cartQtys,
		ExRate:   settings.ExchangeRate,
	})
	err := t.Render(c.Request().Context(), &buf)
	if err != nil {
		return ""
	}
	return buf.String()
}
