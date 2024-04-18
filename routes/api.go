package routes

import (
	"github.com/gofiber/fiber/v2"
	"gocommerce/config"
	"gocommerce/handlers"
	"gocommerce/utils"
)

func RouteInit(app *fiber.App) {
	app.Static("/public", config.ProjectRootPath+"/public/asset")

	api := app.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Gocommerce")
	})

	// Auth API
	api.Post("/register", handlers.RegisterUserHandler).Name("register")

	// User API
	user := api.Group("/user")
	user.Get("/", handlers.GetAllUserHandler).Name("user.index")
	user.Get("/:userId", handlers.GetByIdUserHandler).Name("user.show")
	user.Put("/:userId", handlers.UpdateUserHandler).Name("user.update")
	user.Put("/:userId/update-email", handlers.UpdateEmailUserHandler).Name("user.emailUpdate")
	user.Delete("/:userId", handlers.DeleteUserHandler).Name("user.destroy")

	// Category API
	category := api.Group("/category")
	category.Get("/", handlers.GetAllCategoriesHandler).Name("category.index")
	category.Post("/", handlers.StoreCategoryHandler).Name("category.store")
	category.Get("/:categorySlug", handlers.GetBySlugCategoryHandler).Name("category.show")
	category.Put("/:categoryId", handlers.UpdateCategoryHandler).Name("category.update")
	category.Delete("/:categoryId", handlers.DeleteCategoryHandler).Name("category.destroy")

	// Product API
	product := api.Group("/product")
	product.Get("/", handlers.GetAllProductsHandler).Name("product.index")
	product.Post("/", utils.HandleMultipleFile, handlers.StoreProductHandler).Name("product.store")
	product.Get("/:productSlug", handlers.GetBySlugProductHandler).Name("product.show")
	product.Put("/:productId", utils.HandleMultipleFile, handlers.UpdateProductHandler).Name("product.update")
	product.Delete("/:productId", handlers.DeleteProductHandler).Name("product.destroy")

	// Cart API
	cart := api.Group("/cart")
	cart.Get("/", handlers.GetAllCartsHandler).Name("cart.index")
	cart.Post("/", handlers.StoreCartHandler).Name("cart.store")
	cart.Get("/?cart=:cartId", handlers.ShowByIdCartHandler).Name("cart.show")
	cart.Get("/?user=:userId", handlers.ShowByUserIdCartHandler).Name("cart.showByUser")
	cart.Put("/:cartId", handlers.UpdateQuantityCartHandler).Name("cart.updateQuantity")
	cart.Delete("/delete", handlers.DeleteCartHandler).Name("cart.destroy")

	// Order API
	order := api.Group("/order")
	order.Get("/", handlers.GetAllOrderHandler).Name("order.index")
	order.Post("/create-order", handlers.CheckoutImmediatelyOrderHandler).Name("order.checkout-immediately")
	order.Post("/create-order-by-cart", handlers.CheckoutByCartOrderHandler).Name("order.checkout-by-cart")
	order.Put("/:orderNumber", handlers.UpdateOrderHandler).Name("order.update-status")

}
