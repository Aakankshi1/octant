/*
 *  Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 *  SPDX-License-Identifier: Apache-2.0
 *
 */

package octant

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/vmware-tanzu/octant/internal/api"
	"github.com/vmware-tanzu/octant/internal/log"
)

// HandlerFactory is a factory that generate's Octant's HTTP handler. Octant has
// a frontend and backend handler and they will both be populated in the
// generated handler.
type HandlerFactory struct {
	frontendHandler func(ctx context.Context) (http.Handler, error)
	backendHandler  func(ctx context.Context) (http.Handler, error)
}

// NewHandlerFactory creates an instance of HandlerFactory.
func NewHandlerFactory(optionList ...Option) *HandlerFactory {
	opts := buildOptions(optionList...)

	hf := HandlerFactory{
		frontendHandler: opts.frontendHandler,
		backendHandler:  opts.backendHandler,
	}

	return &hf
}

// Handler returns an HTTP handler.
func (hf *HandlerFactory) Handler(ctx context.Context) (http.Handler, error) {
	router := mux.NewRouter()

	backendHandler, err := hf.backendHandler(ctx)
	if err != nil {
		return nil, err
	}

	frontendHandler, err := hf.frontendHandler(ctx)
	if err != nil {
		return nil, err
	}

	router.PathPrefix(api.PathPrefix).Handler(backendHandler)

	router.PathPrefix("/").Handler(frontendHandler)

	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Language", "Origin", "Content-Type"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	return handlers.CORS(allowedOrigins, allowedHeaders, allowedMethods)(router), nil
}

// ProxiedFrontendHandler creates an HTTP handler that proxies to a target URL.
type ProxiedFrontendHandler struct {
	handler *httputil.ReverseProxy
}

var _ http.Handler = &ProxiedFrontendHandler{}

// NewProxiedFrontend creates an instance of ProxiedFrontendHandler. It will return an
// error if the proxy can not be created.
func NewProxiedFrontend(ctx context.Context, targetURL string) (*ProxiedFrontendHandler, error) {
	logger := log.From(ctx)
	logger.With("proxy-url", targetURL).
		Infof("Creating reverse proxy for Octant's frontend")

	handler, err := createProxyTarget(targetURL)
	if err != nil {
		return nil, err
	}

	pf := ProxiedFrontendHandler{
		handler: handler,
	}

	return &pf, nil
}

// ServeHTTP allows ProxiedFrontendHandler to satisfy the http.Handler interface.
func (p ProxiedFrontendHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.handler.ServeHTTP(w, r)
}

func createProxyTarget(targetURL string) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	if target.Scheme == "" {
		target.Scheme = "http"
	}

	handler := httputil.NewSingleHostReverseProxy(target)
	return handler, nil
}
