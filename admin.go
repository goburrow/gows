// Copyright 2015 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license. See the LICENSE file for details.
package gows

import (
	"fmt"
	"github.com/goburrow/health"
	"net/http"
)

const (
	metricsUri     = "/metrics"
	pingUri        = "/ping"
	runtimeUri     = "/runtime"
	healthCheckUri = "/healthcheck"

	adminHTML = `<!DOCTYPE html>
<html>
<head>
	<title>Operational Menu</title>
</head>
<body>
	<h1>Operational Menu</h1>
	<ul>
		<li><a href="%[1]s%[2]s">Metrics</a></li>
		<li><a href="%[1]s%[3]s">Ping</a></li>
		<li><a href="%[1]s%[4]s">Runtime</a></li>
		<li><a href="%[1]s%[5]s">Healthcheck</a></li>
	</ul>
</body>
</html>
`
)

type AdminEnvironment struct {
	ServerHandler       ServerHandler
	HealthCheckRegistry health.Registry
}

func NewAdminEnvironment() *AdminEnvironment {
	return &AdminEnvironment{
		HealthCheckRegistry: health.NewRegistry(),
	}
}

// Initialize registers all required HTTP handlers
func (env *AdminEnvironment) Initialize(contextPath string) {
	env.ServerHandler.Handle(pingUri, http.HandlerFunc(handleAdminPing))
	env.ServerHandler.Handle(runtimeUri, http.HandlerFunc(handleAdminRuntime))
	env.ServerHandler.Handle(healthCheckUri, NewHealthCheckHandler(env.HealthCheckRegistry))
	env.ServerHandler.Handle("/", NewAdminHandler(contextPath))
}

// AddTask adds a new task to admin environment
func (env *AdminEnvironment) AddTask(name string, task Task) {
	path := "/tasks/" + name
	env.ServerHandler.Handle(path, task)
}

// DefaultAdminHandler implement http.Handler
type DefaultAdminHandler struct {
	contextPath string
}

// NewAdminHTTPHandler allocates and returns a new adminHTTPHandler
func NewAdminHandler(contextPath string) *DefaultAdminHandler {
	return &DefaultAdminHandler{
		contextPath: contextPath,
	}
}

// ServeHTTP handles request to the root of Admin page
func (handler *DefaultAdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// The "/" pattern matches everything, so we need to check
	// that we're at the root here.
	rootUri := handler.contextPath + "/"
	if r.URL.Path != rootUri {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Cache-Control", "must-revalidate,no-cache,no-store")
	w.Header().Set("Content-Type", "text/html")

	fmt.Fprintf(w, adminHTML, handler.contextPath, metricsUri, pingUri, runtimeUri, healthCheckUri)
	// TODO: handle error
}

// handleAdminPing handles ping request to admin /ping
func handleAdminPing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "must-revalidate,no-cache,no-store")
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("pong\n"))
}
