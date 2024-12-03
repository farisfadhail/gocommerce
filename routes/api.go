package routes

import (
	"github.com/go-playground/validator/v10"
	"gocommerce/config"
	"gocommerce/database"
	"gocommerce/handlers"
	"gocommerce/middleware"
	"gocommerce/repository"
	"gocommerce/utils"

	"github.com/gofiber/fiber/v2"
)

func RouteInit(app *fiber.App) {
	app.Static("/public", config.ProjectRootPath+"/public/asset")

	db := database.DatabaseInit()
	validate := validator.New()

	api := app.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Gocommerce")
	})

	// Auth API
	api.Post("/register", handlers.RegisterHandler).Name("register")
	api.Post("/login", handlers.LoginHandler).Name("login")

	// User API
	user := api.Group("/user", middleware.Authenticated)
	user.Get("/:userId", handlers.GetByIdUserHandler).Name("user.show")
	user.Put("/:userId", handlers.UpdateUserHandler).Name("user.update")
	user.Put("/:userId/update-email", handlers.UpdateEmailUserHandler).Name("user.emailUpdate")
	user.Put("/:userId/update-password", handlers.UpdatePasswordUserHandler).Name("user.passwordUpdate")

	// User API (Admin)
	userAdmin := user.Group("/admin", middleware.IsAdmin)
	userAdmin.Get("/", handlers.GetAllUserHandler).Name("user.index")
	userAdmin.Put("/:userId/update-role", handlers.UpdateRoleUserHandler, middleware.IsAdmin).Name("user.roleUpdate")
	userAdmin.Delete("/:userId", handlers.DeleteUserHandler, middleware.IsAdmin).Name("user.destroy")

	// Category API
	categoryRepo := repository.NewCategoryRepository(db)
	categoryHandler := handlers.NewCategoryHandler(categoryRepo, validate)
	category := api.Group("/category")
	category.Get("/", categoryHandler.GetAllCategoriesHandler).Name("category.index")
	category.Get("/:categorySlug", categoryHandler.GetBySlugCategoryHandler).Name("category.show")

	// Category API (Admin)
	categoryAdmin := category.Group("/admin", middleware.Authenticated, middleware.IsAdmin)
	categoryAdmin.Post("/", categoryHandler.StoreCategoryHandler).Name("category.store")
	categoryAdmin.Put("/:categoryId", categoryHandler.UpdateCategoryHandler).Name("category.update")
	categoryAdmin.Delete("/:categoryId", categoryHandler.DeleteCategoryHandler).Name("category.destroy")

	// Product API
	product := api.Group("/product")
	product.Get("/", handlers.GetAllProductsHandler).Name("product.index")
	product.Post("/", utils.HandleMultipleFile, handlers.StoreProductHandler).Name("product.store")
	product.Get("/:productSlug", handlers.GetBySlugProductHandler).Name("product.show")

	// Product Elasticsearch API
	product.Post("/search", handlers.SearchProductHandler).Name("product.search")

	// Product API (Admin)
	productAdmin := product.Group("/admin", middleware.Authenticated, middleware.IsAdmin)
	productAdmin.Put("/:productId", utils.HandleMultipleFile, handlers.UpdateProductHandler).Name("product.update")
	productAdmin.Delete("/:productId", handlers.DeleteProductHandler).Name("product.destroy")

	// Cart API
	cart := api.Group("/cart", middleware.Authenticated)
	cart.Get("/", handlers.GetAllCartsHandler).Name("cart.index")
	cart.Post("/", handlers.StoreCartHandler).Name("cart.store")
	cart.Get("/?cart=:cartId", handlers.ShowByIdCartHandler).Name("cart.show")
	cart.Get("/?user=:userId", handlers.ShowByUserIdCartHandler).Name("cart.showByUser")
	cart.Put("/:cartId", handlers.UpdateQuantityCartHandler).Name("cart.updateQuantity")
	cart.Delete("/delete", handlers.DeleteCartHandler).Name("cart.destroy")

	// Order API
	order := api.Group("/order", middleware.Authenticated)
	order.Get("/", handlers.GetAllOrderHandler).Name("order.index")
	order.Post("/create-order", handlers.CheckoutImmediatelyOrderHandler).Name("order.checkout-immediately")
	order.Post("/create-order-by-cart", handlers.CheckoutByCartOrderHandler).Name("order.checkout-by-cart")
	order.Put("/:orderNumber", handlers.UpdateOrderHandler).Name("order.update-status")

	// Payment API
	payment := api.Group("/payment", middleware.Authenticated)
	payment.Post("/midtrans-notification", handlers.NotificationMidtrans).Name("payment.midtrans-notification")
}
