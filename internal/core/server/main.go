package server

import "net/http"

type Router struct {
	Router *http.ServeMux
}

func NewRouter() *Router {
	return &Router{
		Router: http.NewServeMux(),
	}
}

func RemoveTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			if req.URL.Path[len(req.URL.Path)-1] == '/' {
				http.Redirect(w, req, req.URL.Path[:len(req.URL.Path)-1], http.StatusPermanentRedirect)
				return
			}
		}
		next.ServeHTTP(w, req)
	})
}

func allowMethod(
	method string, handler func(w http.ResponseWriter, req *http.Request),
) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case method:
			handler(w, req)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (r *Router) Get(pattern string, handler func(w http.ResponseWriter, req *http.Request)) {
	r.Router.HandleFunc(pattern, allowMethod(http.MethodGet, handler))
}

func (r *Router) Post(pattern string, handler func(w http.ResponseWriter, req *http.Request)) {
	r.Router.HandleFunc(pattern, allowMethod(http.MethodPost, handler))
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Router.ServeHTTP(w, req)
}
