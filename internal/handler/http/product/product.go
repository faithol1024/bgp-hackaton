package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/faithol1024/bgp-hackaton/internal/entity/bid"
	productEntity "github.com/faithol1024/bgp-hackaton/internal/entity/product"
	"github.com/faithol1024/bgp-hackaton/internal/entity/user"
	ers "github.com/faithol1024/bgp-hackaton/lib/error"
	"github.com/go-chi/chi"
	"github.com/tokopedia/tdk/go/httpt/response"
	"github.com/tokopedia/tdk/go/log"
	"github.com/tokopedia/tdk/go/tracer"
)

type productUseCase interface {
	Create(ctx context.Context, product productEntity.Product) error
	GetByID(ctx context.Context, id string) (productEntity.Product, error)
	GetAll(ctx context.Context, userID string, role string) ([]productEntity.Product, error)
	Bid(ctx context.Context, bid bid.Bid) (bid.Bid, error)
}

type Handler struct {
	ProductUC productUseCase
}

func New(productUC productUseCase) *Handler {
	return &Handler{
		ProductUC: productUC,
	}
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	span, ctx := tracer.StartFromRequest(r)
	defer span.Finish()

	// params checking
	productID := chi.URLParam(r, "product_id")
	if productID == "" {
		log.Error("[product.GetByID] error Invalid param")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// call the usecase
	product, err := h.ProductUC.GetByID(ctx, productID)
	if err != nil {
		log.Error("[product.GetByID] error from GetByID: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
		response.WriteJSONAPIError(w, r, http.StatusInternalServerError, `error get product by id`)
		return
	}

	if _, err := response.WriteJSONAPIData(w, r, http.StatusOK, product); err != nil {
		log.Error("[gopay.GetByID] error from WriteJSON: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	span, ctx := tracer.StartFromRequest(r)
	defer span.Finish()

	// call the usecase
	products, err := h.ProductUC.GetAll(ctx, "", "")
	if err != nil {
		log.Error("[product.GetAll] error from GetAll: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
		response.WriteJSONAPIError(w, r, http.StatusInternalServerError, `error get all product`)
		return
	}

	if _, err := response.WriteJSONAPIData(w, r, http.StatusOK, products); err != nil {
		log.Error("[product.GetAll] error from WriteJSON: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	}
}

func (h *Handler) GetAllBySeller(w http.ResponseWriter, r *http.Request) {
	span, ctx := tracer.StartFromRequest(r)
	defer span.Finish()

	// params checking
	userID := chi.URLParam(r, "user_id")
	if userID == "" {
		log.Error("[product.GetAllBySeller] Invalid Param")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// call the usecase
	products, err := h.ProductUC.GetAll(ctx, userID, user.RoleSeller)
	if err != nil {
		log.Error("[product.GetAllBySeller] error from GetAllBySeller: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
		response.WriteJSONAPIError(w, r, http.StatusInternalServerError, `error get seller product`)
		return
	}

	if _, err := response.WriteJSONAPIData(w, r, http.StatusOK, products); err != nil {
		log.Error("[product.GetAllBySeller] error from WriteJSON: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	}
}

func (h *Handler) GetAllByBuyer(w http.ResponseWriter, r *http.Request) {
	span, ctx := tracer.StartFromRequest(r)
	defer span.Finish()

	// params checking
	userID := chi.URLParam(r, "user_id")
	if userID == "" {
		log.Error("[product.GetAllByBuyer] Invalid Param")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// call the usecase
	products, err := h.ProductUC.GetAll(ctx, userID, user.RoleSeller)
	if err != nil {
		log.Error("[product.GetAllByBuyer] error from GetAllByBuyer: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
		response.WriteJSONAPIError(w, r, http.StatusInternalServerError, `error get buyer product`)
		return
	}

	if _, err := response.WriteJSONAPIData(w, r, http.StatusOK, products); err != nil {
		log.Error("[product.GetAllByBuyer] error from WriteJSON: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	span, ctx := tracer.StartFromRequest(r)
	defer span.Finish()

	var product productEntity.Product

	// params decode
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Errorf("[product.Create] failed to decode request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = product.Validate()
	if err != nil {
		log.Error("[product.Create] Invalid param")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// call the usecase
	product.Status = productEntity.StatusNew
	err = h.ProductUC.Create(ctx, product)
	if err != nil {
		log.Error("[product.Create] error from GetByID: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
		response.WriteJSONAPIError(w, r, http.StatusInternalServerError, `gaboleh bikin product yee`)
		return
	}

	// send the response
	if _, err := response.WriteJSONAPIData(w, r, http.StatusOK, product); err != nil {
		log.Error("[product.Create] error from WriteJSON: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	}
}

func (h *Handler) Bid(w http.ResponseWriter, r *http.Request) {
	span, ctx := tracer.StartFromRequest(r)
	defer span.Finish()

	var bid bid.Bid

	// params decode
	err := json.NewDecoder(r.Body).Decode(&bid)
	if err != nil {
		log.Errorf("[product.Bid] failed to decode request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = bid.Validate()
	if err != nil {
		log.Error("[product.Bid] Invalid param")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// call the usecase
	bidRes, err := h.ProductUC.Bid(ctx, bid)
	if err != nil {
		log.Error("[product.Bid] error from Bid: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
		response.WriteJSONAPIError(w, r, http.StatusInternalServerError, `error bidding`)
		return
	}

	// send the response
	if _, err := response.WriteJSONAPIData(w, r, http.StatusOK, bidRes); err != nil {
		log.Error("[product.Bid] error from WriteJSON: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	}
}
