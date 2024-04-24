package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gogo-framework/router"
)

func TestReadmeSingleRoutes(t *testing.T) {
	// Create a new router instance
	r := router.NewRouter()

	// Register routes for each HTTP method
	r.GET("/get-endpoint", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	r.POST("/post-endpoint", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("You did a post request!"))
	})

	r.PUT("/put-endpoint", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("You did a put request!"))
	})

	r.PATCH("/patch-endpoint", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("You did a patch request!"))
	})

	r.DELETE("/delete-endpoint", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("You did a delete request!"))
	})

	r.OPTIONS("/options-endpoint", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("You did an options request!"))
	})

	r.HEAD("/head-endpoint", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("You did a head request!"))
	})

	r.CONNECT("/connect-endpoint", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("You did a connect request!"))
	})

	r.TRACE("/trace-endpoint", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("You did a trace request!"))
	})

	// Define test cases for each route
	tests := []struct {
		method     string
		path       string
		statusCode int
		response   string
	}{
		{http.MethodGet, "/get-endpoint/", http.StatusOK, "Hello, World!"},
		{http.MethodPost, "/post-endpoint/", http.StatusOK, "You did a post request!"},
		{http.MethodPut, "/put-endpoint/", http.StatusOK, "You did a put request!"},
		{http.MethodPatch, "/patch-endpoint/", http.StatusOK, "You did a patch request!"},
		{http.MethodDelete, "/delete-endpoint/", http.StatusOK, "You did a delete request!"},
		{http.MethodOptions, "/options-endpoint/", http.StatusOK, "You did an options request!"},
		{http.MethodHead, "/head-endpoint/", http.StatusOK, "You did a head request!"},
		{http.MethodConnect, "/connect-endpoint/", http.StatusOK, "You did a connect request!"},
		{http.MethodTrace, "/trace-endpoint/", http.StatusOK, "You did a trace request!"},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			// Create a request for the specific route
			req, err := http.NewRequest(tt.method, tt.path, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Call the router's ServeHTTP method directly with the test request and ResponseRecorder
			r.ServeHTTP(rr, req)

			// Check the status code
			if status := rr.Code; status != tt.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.statusCode)
			}

			// Check the response body
			if body := rr.Body.String(); body != tt.response {
				t.Errorf("handler returned unexpected body: got %v want %v", body, tt.response)
			}
		})
	}
}

func TestReadmeRouteGroups(t *testing.T) {
	// Create a new router instance
	r := router.NewRouter()

	// Define the route group
	r.Group("grouped", func(rg *router.Router) {
		rg.GET("/get", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, World!"))
		})
		rg.POST("/post", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("You did a post request!"))
		})
	})

	r.Group("multi/level/group", func(rg *router.Router) {
		rg.GET("/get", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, World!"))
		})
		rg.POST("/post", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("You did a post request!"))
		})
	})

	// Define test cases for each route within the group
	tests := []struct {
		method     string
		path       string
		statusCode int
		response   string
	}{
		{http.MethodGet, "/grouped/get/", http.StatusOK, "Hello, World!"},
		{http.MethodPost, "/grouped/post/", http.StatusOK, "You did a post request!"},
		{http.MethodGet, "/multi/level/group/get/", http.StatusOK, "Hello, World!"},
		{http.MethodPost, "/multi/level/group/post/", http.StatusOK, "You did a post request!"},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			// Create a request for the specific route
			req, err := http.NewRequest(tt.method, tt.path, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Call the router's ServeHTTP method directly with the test request and ResponseRecorder
			r.ServeHTTP(rr, req)

			// Check the status code
			if status := rr.Code; status != tt.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.statusCode)
			}

			// Check the response body
			if body := rr.Body.String(); body != tt.response {
				t.Errorf("handler returned unexpected body: got %v want %v", body, tt.response)
			}
		})
	}
}

func TestReadmeMiddleware(t *testing.T) {
	// Create a new router instance
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
	r.GET("/get-endpoint", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}).Use(XTestSingleHeaderMiddleware)

	// Adding middleware to a route group
	r.Group("group", func(rg *router.Router) {
		rg.GET("/get", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, World!"))
		})
	}).Use(XTestGroupHeaderMiddleware)

	// Define test cases for each route
	tests := []struct {
		method     string
		path       string
		statusCode int
		headers    map[string]string
	}{
		// Test case for a route with global middleware
		{http.MethodGet, "/get-endpoint/", http.StatusOK, map[string]string{"X-Test-Global-Middleware": "This header is set from the global middleware!"}},
		// Test case for a route with single route middleware
		{http.MethodGet, "/get-endpoint/", http.StatusOK, map[string]string{"X-Test-Single-Header": "This header is set from the single route middleware!"}},
		// Test case for a route within a route group with group middleware
		{http.MethodGet, "/group/get/", http.StatusOK, map[string]string{"X-Test-Group-Header": "This header is set from the route group middleware!"}},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			// Create a request for the specific route
			req, err := http.NewRequest(tt.method, tt.path, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Call the router's ServeHTTP method directly with the test request and ResponseRecorder
			r.ServeHTTP(rr, req)

			// Check the status code
			if status := rr.Code; status != tt.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.statusCode)
			}

			// Check the response headers
			for key, value := range tt.headers {
				if headerValue := rr.Header().Get(key); headerValue != value {
					t.Errorf("handler returned wrong header value for %s: got %v want %v", key, headerValue, value)
				}
			}
		})
	}
}
