package app

import (
	"desgruppe/internal/admin"
	"desgruppe/internal/cart"
	"desgruppe/internal/config"
	"desgruppe/internal/databases"
	"desgruppe/internal/html"
	"desgruppe/internal/logger"
	"desgruppe/internal/notification"
	"desgruppe/internal/queries"
	"desgruppe/internal/templ"
	"desgruppe/internal/templ/api"
	pagesAdmin "desgruppe/internal/templ/views/pages/admin"
	"desgruppe/internal/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type App struct {
	Server   *echo.Echo
	Database *sqlx.DB
	Config   *config.Config
}

func New() (*App, error) {
	app := new(App)

	// Config

	cfg, err := config.New()
	if err != nil {
		return nil, err
	}
	app.Config = cfg

	e := echo.New()
	renderer := &templ.Renderer{}
	e.Renderer = renderer

	// ENDPOINTS
	e.GET("/robots.txt", func(c echo.Context) error {
		file := "User-agent: * \nAllow: / \nSitemap: https://desgruppe.shop/sitemap.xml"
		return c.String(200, file)
	})
	e.GET("/terms", func(c echo.Context) error { return html.LoadHtmlPage(c, "terms.html") })
	e.GET("/about", func(c echo.Context) error { return html.LoadHtmlPage(c, "about.html") })
	e.GET("/contacts", func(c echo.Context) error { return html.LoadHtmlPage(c, "contacts.html") })
	e.GET("/personal-data", func(c echo.Context) error { return html.LoadHtmlPage(c, "personal-data.html") })
	e.GET("/logistics", func(c echo.Context) error { return html.LoadHtmlPage(c, "logistics.html") })
	e.GET("/designers", func(c echo.Context) error { return html.LoadHtmlPage(c, "designers.html") })
	e.GET("/producers", func(c echo.Context) error { return html.LoadHtmlPage(c, "producers.html") })
	e.GET("/cart", func(c echo.Context) error {
		return api.RenderCartPage(c, app.Database)
	})
	e.GET("/", func(c echo.Context) error {
		cart.CartMiddleware(c, app.Database)
		return api.RenderIndexPage(c, app.Database)
	})
	e.GET("/available", func(c echo.Context) error {
		cart.CartMiddleware(c, app.Database)
		return api.RenderAvailablePage(c, app.Database)
	})
	e.GET("/furniture", func(c echo.Context) error {
		cart.CartMiddleware(c, app.Database)
		return api.RenderFurniturePage(c, app.Database, "furniture", "Мебель")
	})
	e.GET("/light", func(c echo.Context) error {
		cart.CartMiddleware(c, app.Database)
		return api.RenderFurniturePage(c, app.Database, "light", "Свет")
	})
	e.GET("/interior", func(c echo.Context) error {
		cart.CartMiddleware(c, app.Database)
		return api.RenderFurniturePage(c, app.Database, "interior", "Аксессуары")
	})
	e.GET("/furniturefilter", func(c echo.Context) error {
		cart.CartMiddleware(c, app.Database)
		return api.RenderFilterFurniturePage(c, app.Database, "furniture", "Мебель")
	})
	e.GET("/lightfilter", func(c echo.Context) error {
		cart.CartMiddleware(c, app.Database)
		return api.RenderFilterFurniturePage(c, app.Database, "light", "Свет")
	})
	e.GET("/interiorfilter", func(c echo.Context) error {
		cart.CartMiddleware(c, app.Database)
		return api.RenderFilterFurniturePage(c, app.Database, "interior", "Аксессуары")
	})

	e.Static("/static/", "frontend/static/")
	e.Static("/script/", "internal/templ/views/static/")
	e.Static("/admin/static/admin/", "frontend/static/admin/")
	e.Static("/imgs/", "frontend/imgs")
	e.Static("/admin/imgs/", "frontend/imgs")
	e.Static("/product/static/", "frontend/static")
	e.Static("/designer/static/", "frontend/static")
	e.Static("/producer/static/", "frontend/static")
	e.Static("/product/imgs/", "frontend/imgs")
	e.Static("/designer/imgs/", "frontend/imgs")
	e.Static("/producer/imgs/", "frontend/imgs")

	e.GET("/admin/", admin.AdminMiddleware(func(c echo.Context) error { return html.LoadAdminPage(c, "sections") }))
	e.GET("/admin", admin.AdminMiddleware(func(c echo.Context) error { return html.LoadAdminPage(c, "sections") }))

	e.POST("/auth-des-admin", func(c echo.Context) error { return admin.LoginAdmin(c, app.Config) })

	e.GET("/admin-login", func(c echo.Context) error {
		authcookie, err := c.Cookie("admin-logged")

		if err != nil || authcookie.Value != "true" {
			if err := html.LoadAdminPage(c, "login"); err != nil {
				return c.String(http.StatusInternalServerError, "Ошибка при загрузке страницы: "+err.Error())
			}
			return nil
		}

		return c.Redirect(http.StatusSeeOther, "/admin/sections")
	})
	admin := e.Group("/admin", admin.AdminMiddleware)

	admin.GET("/:page", func(c echo.Context) error { return html.LoadAdminPage(c, c.Param("page")) })

	e.GET("/:tableName/getImgBlob/:columnName/:ID", func(c echo.Context) error {
		return utils.GetImgBlob(c, app.Database, c.Param("tableName"), c.Param("columnName"), c.Param("ID"))
	})
	admin.GET("/products", func(c echo.Context) error { return api.RenderProductsPage(c, app.Database) })
	admin.GET("/available", func(c echo.Context) error { return api.RenderAdminAvailablePage(c, app.Database) })
	admin.GET("/best", func(c echo.Context) error { return api.RenderAdminBestPage(c, app.Database) })
	admin.GET("/products/editpage", func(c echo.Context) error {
		product, err := utils.GetProductById(app.Database, c.QueryParam("id"))
		if err != nil {
			return err
		}
		designers, err := utils.ListDesigners(app.Database)
		if err != nil {
			return err
		}
		producers, err := utils.ListProducers(app.Database)
		if err != nil {
			return err
		}
		sections, err := utils.ListSections(app.Database)
		if err != nil {
			return err
		}
		colors, err := utils.ListColors(app.Database)
		if err != nil {
			return err
		}
		return c.Render(200, "editpage", pagesAdmin.ProductsEditPage(pagesAdmin.EditPageProps{
			Product:   product,
			Designers: designers,
			Producers: producers,
			Sections:  sections,
			Colors:    colors,
			Type:      c.QueryParam("type"),
		}))
	})
	admin.GET("/orders/editpage", func(c echo.Context) error { return api.RenderOrderEditPage(c, app.Database) })
	admin.GET("/:page/editpage", func(c echo.Context) error { return html.LoadAdminEditPage(c, c.Param("page")) })

	// colors
	admin.POST("/colors/add", func(c echo.Context) error { return queries.AddColor(c, app.Database) })
	admin.PUT("/colors/edit", func(c echo.Context) error { return queries.EditColor(c, app.Database) })
	admin.GET("/colors/list", func(c echo.Context) error { return queries.ListColors(c, app.Database) })
	e.GET("/colors/get", func(c echo.Context) error { return queries.GetColor(c, app.Database) })

	admin.GET("/colors/get", func(c echo.Context) error { return queries.GetColor(c, app.Database) })
	admin.DELETE("/colors/delete", func(c echo.Context) error { return queries.DeleteColor(c, app.Database) })
	admin.GET("/colors/search", func(c echo.Context) error { return queries.SearchColor(c, app.Database) })

	// designers
	admin.POST("/designers/add", func(c echo.Context) error { return queries.AddDesigner(c, app.Database) })
	admin.PUT("/designers/edit", func(c echo.Context) error { return queries.EditDesigner(c, app.Database) })
	admin.GET("/designers/list", func(c echo.Context) error { return queries.ListDesigners(c, app.Database) })
	e.GET("/designers", func(c echo.Context) error { return api.RenderDesignersPage(c, app.Database) })
	e.GET("/designers/list", func(c echo.Context) error { return queries.ListDesigners(c, app.Database) })
	e.GET("/designers/get", func(c echo.Context) error { return queries.GetDesigner(c, app.Database) })
	admin.GET("/designers/get", func(c echo.Context) error { return queries.GetDesigner(c, app.Database) })
	admin.DELETE("/designers/delete", func(c echo.Context) error { return queries.DeleteDesigner(c, app.Database) })
	admin.GET("/designers/search", func(c echo.Context) error { return queries.SearchDesigner(c, app.Database) })
	e.GET("/designer/:slug", func(c echo.Context) error { return html.LoadDesignerPage(c, app.Database) })
	e.GET("/designer/get/:slug", func(c echo.Context) error {
		designer, err := queries.GetDesignerBySlug(c, app.Database)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, designer)
	})

	// producers
	admin.POST("/producers/add", func(c echo.Context) error { return queries.AddProducer(c, app.Database) })
	admin.PUT("/producers/edit", func(c echo.Context) error { return queries.EditProducer(c, app.Database) })
	admin.GET("/producers/list", func(c echo.Context) error { return queries.ListProducers(c, app.Database) })
	e.GET("/producers/list", func(c echo.Context) error { return queries.ListProducers(c, app.Database) })
	e.GET("/producers/get", func(c echo.Context) error { return queries.GetProducer(c, app.Database) })
	admin.GET("/producers/get", func(c echo.Context) error { return queries.GetProducer(c, app.Database) })
	admin.DELETE("/producers/delete", func(c echo.Context) error { return queries.DeleteProducer(c, app.Database) })
	admin.GET("/producers/search", func(c echo.Context) error { return queries.SearchProducer(c, app.Database) })
	e.GET("/producer/:slug", func(c echo.Context) error { return html.LoadProducerPage(c, app.Database) })
	e.GET("/producer/get/:slug", func(c echo.Context) error {
		designer, err := queries.GetProducerBySlug(c, app.Database)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, designer)
	})

	// sections
	admin.POST("/sections/add", func(c echo.Context) error { return queries.AddSection(c, app.Database) })
	admin.GET("/sections/list", func(c echo.Context) error { return queries.ListSections(c, app.Database) })
	e.GET("/sections/list", func(c echo.Context) error { return queries.ListSections(c, app.Database) })
	admin.DELETE("/sections/delete", func(c echo.Context) error { return queries.DeleteSection(c, app.Database) })

	// products
	admin.POST("/products/add", func(c echo.Context) error { return queries.AddProduct(c, app.Database) })
	admin.PUT("/products/edit", func(c echo.Context) error { return queries.EditProduct(c, app.Database) })
	// admin.GET("/products/list", func(c echo.Context) error { return queries.AdminListProducts(c, app.Database) })
	// e.GET("/products/list1", func(c echo.Context) error { return queries.FirstListProducts(c, app.Database) })
	// e.GET("/products/list2", func(c echo.Context) error { return queries.SecondListProducts(c, app.Database) })
	admin.GET("/products/get", func(c echo.Context) error { return queries.GetProduct(c, app.Database) })
	admin.DELETE("/products/delete", func(c echo.Context) error { return queries.DeleteProduct(c, app.Database) })
	admin.GET("/products/searchcolor", func(c echo.Context) error { return queries.SearchColor(c, app.Database) })
	// admin.GET("/products/search", func(c echo.Context) error { return queries.SearchProduct(c, app.Database) })
	// photos
	admin.GET("/products/photos", func(c echo.Context) error { return html.LoadAdminEditPage(c, "products-photos") })
	admin.GET("/products/photos/list", func(c echo.Context) error { return queries.ListPhotos(c, app.Database) })
	admin.GET("/products/photos/getmain", func(c echo.Context) error { return queries.GetMainPhoto(c, app.Database) })
	admin.POST("/products/photos/editmain", func(c echo.Context) error { return queries.EditMainPhoto(c, app.Database) })
	admin.POST("/products/photos/edit", func(c echo.Context) error { return queries.EditPhoto(c, app.Database) })
	admin.POST("/products/photos/add", func(c echo.Context) error { return queries.AddPhoto(c, app.Database) })
	admin.DELETE("/products/photos/delete", func(c echo.Context) error { return queries.DeletePhoto(c, app.Database) })
	admin.DELETE("/products/photos/deleteall", func(c echo.Context) error { return queries.DeletePhotosForProduct(c, app.Database) })
	e.GET("/product/photos/get", func(c echo.Context) error {
		products, err := queries.ListPhotosByID(app.Database, c.QueryParam("id"))
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, products)
	})

	e.GET("/product/:slug", func(c echo.Context) error { return api.RenderProductPage(c, app.Database) })

	e.GET("/products/listbyfilter", func(c echo.Context) error { return queries.ListProductsByFilter(c, app.Database) })

	// settings
	e.GET("/ex-rate", func(c echo.Context) error { return queries.GetSettings(c, app.Database) })
	admin.GET("/exchange-rate", func(c echo.Context) error { return api.RenderSettingsPage(c, app.Database) })
	admin.PUT("/ex-rate", func(c echo.Context) error { return queries.AddSettings(c, app.Database) })

	// carts
	e.GET("/cart/new-id", func(c echo.Context) error { return queries.CreateNewCart(c, app.Database) })
	e.POST("/cart/product-add", func(c echo.Context) error { return queries.AddProductToCart(c, app.Database) })
	// e.GET("/cart/products-list", func(c echo.Context) error {
	// 	return queries.ListProductsFromCart(c, app.Database)
	// })
	e.PUT("/cart/product-edit", func(c echo.Context) error { return queries.EditProductQuantity(c, app.Database) })
	e.DELETE("/cart/product-delete", func(c echo.Context) error { return queries.DeleteProductFromCart(c, app.Database) })

	// orders
	e.POST("/orders/add", func(c echo.Context) error {
		err := queries.AddOrder(c, app.Database)
		if err == nil {
			go notification.SendOrderNotification(c, app.Database)
		}
		return err
	})
	admin.GET("/orders", func(c echo.Context) error { return api.RenderOrdersPage(c, app.Database) })
	admin.PUT("/orders/seen", func(c echo.Context) error { return queries.OrderSeen(c, app.Database) })
	// admin.GET("/orders/list", func(c echo.Context) error { return queries.ListOrders(c, app.Database) })
	admin.GET("/orders/get", func(c echo.Context) error { return queries.GetOrder(c, app.Database) })
	admin.DELETE("/orders/delete", func(c echo.Context) error { return queries.DeleteOrder(c, app.Database) })
	// admin.GET("/orders/search", func(c echo.Context) error { return queries.SearchOrder(c, app.Database) })

	// questions
	e.POST("/questions/add", func(c echo.Context) error { return queries.AddQuestion(c, app.Database) })
	e.PUT("/questions/seen", func(c echo.Context) error { return queries.QuestionSeen(c, app.Database) })
	admin.GET("/questions/list", func(c echo.Context) error { return queries.ListQuestions(c, app.Database) })
	admin.GET("/questions/get", func(c echo.Context) error { return queries.GetQuestion(c, app.Database) })
	admin.DELETE("/questions/delete", func(c echo.Context) error { return queries.DeleteQuestion(c, app.Database) })
	// admin.GET("/orders/search", func(c echo.Context) error { return queries.SearchOrder(c, app.Database) })

	// available
	// e.GET("/available/list", func(c echo.Context) error { return queries.ListAvailable(c, app.Database) })
	admin.PUT("/available/edit", func(c echo.Context) error { return queries.EditAvailable(c, app.Database) })

	// best
	admin.PUT("/best/edit", func(c echo.Context) error { return queries.EditBest(c, app.Database) })

	// Sitemap

	e.GET("/sitemap.xml", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/xml")
		sm, err := utils.GetSitemap(app.Database)
		if err != nil {
			return c.String(echo.ErrInternalServerError.Code, err.Error())
		}
		return c.String(200, sm)
	})

	// БД

	db := databases.ConnectToDB(app.Config)
	if db == nil {
		return nil, fmt.Errorf("Failed to connect to database")
	}
	db.SetMaxOpenConns(7)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(time.Minute * 60)
	err = databases.InitProducts(db)
	err = databases.InitColors(db)
	err = databases.InitDesigners(db)
	err = databases.InitProducers(db)
	err = databases.InitSections(db)
	err = databases.InitProductsGallery(db)
	err = databases.InitSettings(db)
	err = databases.InitCarts(db)
	err = databases.InitOrders(db)
	if err != nil {
		return nil, err
	}
	// AddBestColumn(db)
	app.Server = e
	app.Database = db
	return app, nil
}

func (app *App) Run() error {
	logger.Info("Starting app...")
	err := app.Server.Start(":8080")
	if err != nil {
		logger.Error("Error starting server: ", err)
		return err
	}
	return nil
}

func AddBestColumn(db *sqlx.DB) {
	_, err := db.Query(`ALTER TABLE "products" ADD "Best" boolean NOT NULL DEFAULT false`)
	if err != nil {
		logger.Error(err.Error())
	}
	_, err = db.Query(`DROP TABLE "sales"`)
	if err != nil {
		logger.Error(err.Error())
	}
	_, err = db.Query(`DROP TABLE "best"`)
	if err != nil {
		logger.Error(err.Error())
	}
}
