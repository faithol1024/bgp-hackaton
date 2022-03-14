package main

import (
	"fmt"
	"net/http"

	bidhandler "github.com/faithol1024/bgp-hackaton/internal/handler/http/bid"
	gopayhandler "github.com/faithol1024/bgp-hackaton/internal/handler/http/gopay"
	producthandler "github.com/faithol1024/bgp-hackaton/internal/handler/http/product"
	"github.com/go-chi/chi"
	"github.com/tokopedia/tdk/go/httpt/middleware"
	chiMW "github.com/tokopedia/tdk/go/httpt/middleware/chi"
	"github.com/tokopedia/tdk/go/panics"
)

type RouteHandlers struct {
	gopay   *gopayhandler.Handler
	bid     *bidhandler.Handler
	product *producthandler.Handler
}

func newRoutes(handler RouteHandlers) *chi.Mux {
	router := chi.NewRouter()

	mw := middleware.NewSet(
		middleware.Prometheus(repoName, chiMW.MetricLabels(router)),
		panics.CaptureHandlerFunc,
	)

	router.Method(http.MethodGet, "/gopay/get/{user_id}", mw.HandlerFunc(handler.gopay.GetByUserID))
	router.Method(http.MethodGet, "/products", mw.HandlerFunc(handler.product.GetAll))
	router.Method(http.MethodGet, "/products/buyer", mw.HandlerFunc(handler.product.GetAllByBuyer))
	router.Method(http.MethodGet, "/products/seller", mw.HandlerFunc(handler.product.GetAllBySeller))
	router.Method(http.MethodGet, "/product/{product_id}", mw.HandlerFunc(handler.product.GetByID))
	router.Method(http.MethodPost, "/product", mw.HandlerFunc(handler.product.Create))
	router.Method(http.MethodPost, "/bid", mw.HandlerFunc(handler.product.Bid))
	router.Method(http.MethodGet, "/ping", mw.HandlerFunc(Ping))

	//for test
	router.Method(http.MethodGet, "/bid", mw.HandlerFunc(handler.bid.PublishBidFRDB))

	return router
}

func Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Pong")
}
