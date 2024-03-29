package main

import (
	"fmt"
	"net/http"

	bidhandler "github.com/faithol1024/bgp-hackaton/internal/handler/http/bid"
	gopayhandler "github.com/faithol1024/bgp-hackaton/internal/handler/http/gopay"
	producthandler "github.com/faithol1024/bgp-hackaton/internal/handler/http/product"
	userhandler "github.com/faithol1024/bgp-hackaton/internal/handler/http/user"
	"github.com/go-chi/chi"
	"github.com/tokopedia/tdk/go/httpt/middleware"
	chiMW "github.com/tokopedia/tdk/go/httpt/middleware/chi"
	"github.com/tokopedia/tdk/go/panics"
)

type RouteHandlers struct {
	user    *userhandler.Handler
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
	router.Method(http.MethodGet, "/user/{user_id}", mw.HandlerFunc(handler.user.GetByUserID))
	router.Method(http.MethodPost, "/user", mw.HandlerFunc(handler.user.Create))

	router.Method(http.MethodGet, "/gopay/user/{user_id}", mw.HandlerFunc(handler.gopay.GetByUserID))
	router.Method(http.MethodGet, "/gopay/user/history/{user_id}", mw.HandlerFunc(handler.gopay.GetAllHistoryByUserID))
	router.Method(http.MethodGet, "/products", mw.HandlerFunc(handler.product.GetAll))
	router.Method(http.MethodGet, "/products/buyer/{user_id}", mw.HandlerFunc(handler.product.GetAllByBuyer))
	router.Method(http.MethodGet, "/products/seller/{user_id}", mw.HandlerFunc(handler.product.GetAllBySeller))
	router.Method(http.MethodGet, "/product/{product_id}", mw.HandlerFunc(handler.product.GetByID))
	router.Method(http.MethodPost, "/product", mw.HandlerFunc(handler.product.Create))
	router.Method(http.MethodPost, "/bid", mw.HandlerFunc(handler.product.Bid))
	router.Method(http.MethodGet, "/ping", mw.HandlerFunc(Ping))

	//for test
	router.Method(http.MethodGet, "/bid", mw.HandlerFunc(handler.bid.PublishBidFRDB))

	return router
}

func Ping(w http.ResponseWriter, r *http.Request) {

	//Allow CORS here By *
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	fmt.Println("Pong")
}
