package controllers

// Controller handle the first layer of the application including routes,
// middlewares, and handlers.
type Controller struct {
}

func New() *Controller {
	return &Controller{}
}
