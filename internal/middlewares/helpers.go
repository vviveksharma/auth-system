package middlewares

import "github.com/gofiber/fiber/v2"

// ApplyMiddlewareChain is a utility function to apply multiple middlewares to a route group
func ApplyMiddlewareChain(group fiber.Router, middlewares []fiber.Handler) {
	for _, middleware := range middlewares {
		group.Use(middleware)
	}
}

// WithMiddleware is a helper to combine middleware chain with handler
func WithMiddleware(middlewares []fiber.Handler, handler fiber.Handler) []fiber.Handler {
	return append(middlewares, handler)
}


