package server

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// RegisterStaticRoutes serves the built frontend SPA from ./webui/dist
// with automatic fallback to index.html for client-side routing.
func (s *Server) RegisterStaticRoutes(mux *http.ServeMux) {
	distDir := "./webui/dist"

	serveSPA := func(w http.ResponseWriter, r *http.Request) {
		// If webui/dist does not exist, return a descriptive message
		if _, err := os.Stat(distDir); os.IsNotExist(err) {
			http.Error(w, "Carrpigeon Web UI is not built. Run 'cd webui && npm run build'", http.StatusNotFound)
			return
		}

		reqPath := strings.TrimPrefix(r.URL.Path, "/")
		if reqPath == "" {
			http.ServeFile(w, r, filepath.Join(distDir, "index.html"))
			return
		}

		// Check if the requested file exists on disk
		targetPath := filepath.Join(distDir, filepath.Clean(reqPath))
		info, err := os.Stat(targetPath)
		if err == nil && !info.IsDir() {
			http.ServeFile(w, r, targetPath)
			return
		}

		// Fallback to index.html for SPA routing (e.g. /dashboard, /send, /receivers)
		http.ServeFile(w, r, filepath.Join(distDir, "index.html"))
	}

	// Serve root and any client-side routes not claimed by API handlers
	mux.HandleFunc("GET /{$}", serveSPA)
	mux.HandleFunc("GET /{path...}", func(w http.ResponseWriter, r *http.Request) {
		path := r.PathValue("path")
		// Do not intercept swagger or known prefixes if not handled
		if strings.HasPrefix(path, "swagger") {
			http.NotFound(w, r)
			return
		}
		serveSPA(w, r)
	})
}
