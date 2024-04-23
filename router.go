package router

import (
	"fmt"
	"log"
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func applyMiddlewares(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

type Route struct {
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	Middlewares []Middleware
}

func (r *Route) Use(middleware ...Middleware) *Route {
	r.Middlewares = append(r.Middlewares, middleware...)
	return r
}

type RouteGroup struct {
	Prefix      string
	Middlewares []Middleware
	Routes      []*Route
}

func (rg *RouteGroup) Use(middleware ...Middleware) *RouteGroup {
	rg.Middlewares = append(rg.Middlewares, middleware...)
	return rg
}

type Router struct {
	mux            *http.ServeMux
	routes         []*Route
	routeGroups    []*RouteGroup
	middlewares    []Middleware
	hasSetupRoutes bool
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) SetMux(mux *http.ServeMux) {
	r.mux = mux
}

func (r *Router) RegisterRoute(method string, pattern string, handler http.HandlerFunc) *Route {
	route := &Route{
		Method:      method,
		Pattern:     pattern,
		HandlerFunc: handler,
		Middlewares: nil,
	}
	r.routes = append(r.routes, route)
	return route
}

func (r *Router) GET(pattern string, handler http.HandlerFunc) *Route {
	return r.RegisterRoute(http.MethodGet, pattern, handler)
}

func (r *Router) POST(pattern string, handler http.HandlerFunc) *Route {
	return r.RegisterRoute(http.MethodPost, pattern, handler)
}

func (r *Router) PUT(pattern string, handler http.HandlerFunc) *Route {
	return r.RegisterRoute(http.MethodPut, pattern, handler)
}

func (r *Router) DELETE(pattern string, handler http.HandlerFunc) *Route {
	return r.RegisterRoute(http.MethodDelete, pattern, handler)
}

func (r *Router) PATCH(pattern string, handler http.HandlerFunc) *Route {
	return r.RegisterRoute(http.MethodPatch, pattern, handler)
}

func (r *Router) OPTIONS(pattern string, handler http.HandlerFunc) *Route {
	return r.RegisterRoute(http.MethodOptions, pattern, handler)
}

func (r *Router) HEAD(pattern string, handler http.HandlerFunc) *Route {
	return r.RegisterRoute(http.MethodHead, pattern, handler)
}

func (r *Router) CONNECT(pattern string, handler http.HandlerFunc) *Route {
	return r.RegisterRoute(http.MethodConnect, pattern, handler)
}

func (r *Router) TRACE(pattern string, handler http.HandlerFunc) *Route {
	return r.RegisterRoute(http.MethodTrace, pattern, handler)
}

func (r *Router) Group(prefix string, group func(r *Router)) *RouteGroup {
	tmpRouter := &Router{middlewares: make([]Middleware, len(r.middlewares))}
	copy(tmpRouter.middlewares, r.middlewares)
	group(tmpRouter)
	rg := &RouteGroup{
		Prefix:      prefix,
		Routes:      tmpRouter.routes,
		Middlewares: tmpRouter.middlewares,
	}
	r.routeGroups = append(r.routeGroups, rg)
	return rg
}

func (r *Router) Use(middleware ...Middleware) {
	r.middlewares = append(r.middlewares, middleware...)
}

func (r *Router) SetupRoutes() {
	if r.mux == nil {
		log.Println("Warning: ServeMux is nil, creating a default one")
		r.mux = http.NewServeMux()
	}

	// This function combines the global middlewares with the route middlewares and the route group middlewares
	combineMiddlewares := func(routeMiddlewares []Middleware, globalMiddlewares []Middleware) []Middleware {
		allMiddlewares := make([]Middleware, 0, len(globalMiddlewares)+len(routeMiddlewares))
		allMiddlewares = append(allMiddlewares, globalMiddlewares...)
		allMiddlewares = append(allMiddlewares, routeMiddlewares...)
		return allMiddlewares
	}

	// Set up single routes
	for _, route := range r.routes {
		handler := applyMiddlewares(
			route.HandlerFunc,
			combineMiddlewares(route.Middlewares, r.middlewares)...,
		)
		r.mux.HandleFunc(fmt.Sprintf("%s %s", route.Method, route.Pattern), func(w http.ResponseWriter, req *http.Request) {
			handler(w, req)
		})
	}

	// Set up route groups routes
	for _, routeGroup := range r.routeGroups {
		for _, route := range routeGroup.Routes {
			handler := applyMiddlewares(
				route.HandlerFunc,
				combineMiddlewares(append(routeGroup.Middlewares, route.Middlewares...), r.middlewares)...,
			)
			r.mux.HandleFunc(fmt.Sprintf("%s %s%s", route.Method, routeGroup.Prefix, route.Pattern), func(w http.ResponseWriter, req *http.Request) {
				handler(w, req)
			})
		}
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !r.hasSetupRoutes {
		r.SetupRoutes()
		r.hasSetupRoutes = true
	}
	r.mux.ServeHTTP(w, req)
}
