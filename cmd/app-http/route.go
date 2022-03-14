package main

import (
	"net/http"

	gopayhandler "github.com/faithol1024/bgp-hackhaton/internal/handler/http/gopay"
	"github.com/go-chi/chi"
	"github.com/tokopedia/tdk/go/httpt/middleware"
	chiMW "github.com/tokopedia/tdk/go/httpt/middleware/chi"
	"github.com/tokopedia/tdk/go/panics"
)

type RouteHandlers struct {
	gopay *gopayhandler.Handler
}

func newRoutes(handler RouteHandlers) *chi.Mux {
	router := chi.NewRouter()

	mw := middleware.NewSet(
		middleware.Prometheus(repoName, chiMW.MetricLabels(router)),
		panics.CaptureHandlerFunc,
	)

	router.Method(http.MethodGet, "/gopay/get", mw.HandlerFunc(handler.gopay.GetByUserID))

	return router
}
