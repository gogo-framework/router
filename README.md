# router

This package comes with a nicer to use API for the newly introduced routing enhancements in the Go 1.22 release. It adds the following features:

- It allows you to use HTTP methods as functions rather than it being part of the URL path.
- It allows you to group routes under a common prefix.
- It allows you to add middlewares to the router itself, single routes and route groups.

## Usage

### Creating a new router

To create a new router, use the `NewRouter` function.

```go
r := router.NewRouter()
```

### Set custom mux to the router

If you want to use a custom mux, you can set it using the `SetMux` method.

```go
r.SetMux(http.NewServeMux())
```

### Adding routes

You can add add routes in two ways; Single routes, or grouped routes.

#### Single routes

```go
r := router.NewRouter()

r.GET("/get-endpoint/", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
})

r.POST("/post-endpoint/", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You did a post request!"))
})

r.PUT("/put-endpoint/", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You did a put request!"))
})

r.PATCH("/patch-endpoint/", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You did a patch request!"))
})

r.DELETE("/delete-endpoint/", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You did a delete request!"))
})

r.OPTIONS("/options-endpoint/", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You did an options request!"))
})

r.HEAD("/head-endpoint/", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You did a head request!"))
})

r.CONNECT("/connect-endpoint/", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You did a connect request!"))
})

r.TRACE("/trace-endpoint/", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You did a trace request!"))
})
```

#### Grouped routes

You can group routes under a common prefix using the `Group` method.

```go
r := router.NewRouter()
r.Group("/group/", func(rg *router.Router) {
	rg.GET("get/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	rg.POST("post/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("You did a post request!"))
	})
})
```

The above will create two routes, `/group/get/` and `/group/post/`.

### Middlewares

You can add middlewares to router itself, single routes and route groups using the `Use` method.

```go
r := router.NewRouter()

// Define middleware functions
XTestGlobalMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Test-Global-Middleware", "This header is set from the global middleware!")
		next.ServeHTTP(w, r)
	}
}
XTestSingleHeaderMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Test-Single-Header", "This header is set from the single route middleware!")
		next.ServeHTTP(w, r)
	}
}
XTestGroupHeaderMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Test-Group-Header", "This header is set from the route group middleware!")
		next.ServeHTTP(w, r)
	}
}

// Adding middleware to the router itself
r.Use(XTestGlobalMiddleware)

// Adding middleware to a single route
r.GET("/get-endpoint/", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}).Use(XTestSingleHeaderMiddleware)

// Adding middleware to a route group
r.Group("/group/", func(rg *router.Router) {
	rg.GET("get/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
}).Use(XTestGroupHeaderMiddleware)
```

### Improved pattern matching

This is a todo, but the default pattern matching is quite strange and can result in unexpected behavior. See [this](https://pkg.go.dev/net/http#ServeMux) for more information. I'll probably have some settings added to the router to allow for more control over the pattern matching, but I'll need to do some research on how other routers do this.

## Things I'd like to add

- Host/domain matching
- Route naming (maybe?)
- ...More?
