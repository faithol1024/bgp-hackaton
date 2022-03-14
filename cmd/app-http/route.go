package main

import (
	"fmt"
	"net/http"

	gopayhandler "github.com/faithol1024/bgp-hackaton/internal/handler/http/gopay"
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

	router.Method(http.MethodGet, "/gopay/get/{user_id}", mw.HandlerFunc(handler.gopay.GetByUserID))
	router.Method(http.MethodGet, "/ping", mw.HandlerFunc(Ping))

	return router
}

func Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Pong")
}
