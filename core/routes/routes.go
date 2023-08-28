package routes

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	ServerConfig "github.com/hope-ag/go-dynamo/config"
	HealthHandler "github.com/hope-ag/go-dynamo/core/handlers/health"
	ProductHandler "github.com/hope-ag/go-dynamo/core/handlers/product"
	"github.com/hope-ag/go-dynamo/core/repository/adapter"
)

type Router struct {
	config *Config
	router *chi.Mux
}

func NewRouter() *Router {
	return &Router{
		config: NewConfig(ServerConfig.GetConfig().Timeout),
		router: chi.NewRouter(),
	}
}

func (r *Router) SetRouters(repository adapter.Interface) *chi.Mux {
	r.setConfigRouters()
	r.RouterHealth(repository)
	r.RouterProduct(repository)
	return r.router
}

func (r *Router) setConfigRouters() *chi.Mux {
	r.EnableCors()
	r.EnableLogger()
	r.EnableTimeout()
	r.EnableRecovery()
	r.EnableRequestId()
	r.EnableRequestIP()
	r.EnableRequestId()
	return r.router
}

func (r *Router) RouterHealth(repository adapter.Interface) {
	handler := HealthHandler.NewHandler(repository)
	r.router.Route("/health", func (route chi.Router) {
		route.Post("/", handler.Post)
		route.Get("/", handler.Get)
		route.Put("/", handler.Put)
		route.Delete("/", handler.Delete)
		route.Options("/", handler.Options)
	})
}
func (r *Router) RouterProduct(repository adapter.Interface) {
	handler := ProductHandler.NewHandler(repository)
	r.router.Route("/products", func (route chi.Router) {
		route.Get("/", handler.Get)
		route.Get("/{ID}", handler.Get)
		route.Post("/", handler.Post)
		route.Put("/{ID}", handler.Put)
		route.Delete("/{ID}", handler.Delete)
		route.Options("/", handler.Options)
	})
}
func (r *Router) EnableCors() * Router  {
	r.router.Use(r.config.Cors)
	return r
}

func (r *Router) EnableLogger() *Router {
	r.router.Use(middleware.Logger)
	return r
}

func (r *Router) EnableTimeout() * Router {
	r.router.Use(middleware.Timeout(r.config.GetTimeout()))
	return r
}

func (r *Router) EnableRecovery() * Router  {
	r.router.Use(middleware.Recoverer)
	return r
}

func (r *Router) EnableRequestId() * Router  {
	r.router.Use(middleware.RequestID)
	return r
}

func (r *Router) EnableRequestIP() * Router  {
	r.router.Use(middleware.RealIP)
	return r
}