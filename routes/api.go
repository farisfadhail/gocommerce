package routes

import (
	"github.com/gofiber/fiber/v2"
	"gocommerce/config"
	"gocommerce/handlers"
	"gocommerce/middleware"
	"gocommerce/utils"
)

func RouteInit(app *fiber.App) {
	app.Static("/public", config.ProjectRootPath+"/public/asset")

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
	category := api.Group("/category")
	category.Get("/", handlers.GetAllCategoriesHandler).Name("category.index")
	category.Get("/:categorySlug", handlers.GetBySlugCategoryHandler).Name("category.show")

	// Category API (Admin)
	categoryAdmin := category.Group("/admin", middleware.Authenticated, middleware.IsAdmin)
	categoryAdmin.Post("/", handlers.StoreCategoryHandler).Name("category.store")
	categoryAdmin.Put("/:categoryId", handlers.UpdateCategoryHandler).Name("category.update")
	categoryAdmin.Delete("/:categoryId", handlers.DeleteCategoryHandler).Name("category.destroy")

	// Product API
	product := api.Group("/product")
	product.Get("/", handlers.GetAllProductsHandler).Name("product.index")
	product.Post("/", utils.HandleMultipleFile, handlers.StoreProductHandler).Name("product.store")
	product.Get("/:productSlug", handlers.GetBySlugProductHandler).Name("product.show")

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
