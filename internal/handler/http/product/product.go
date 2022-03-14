package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/faithol1024/bgp-hackaton/internal/entity/bid"
	"github.com/faithol1024/bgp-hackaton/internal/entity/product"
	ers "github.com/faithol1024/bgp-hackaton/lib/error"
	"github.com/go-chi/chi"
	"github.com/tokopedia/tdk/go/httpt/response"
	"github.com/tokopedia/tdk/go/log"
)

type productUseCase interface {
	Create(ctx context.Context, product product.Product) error
}

type Handler struct {
	ProductUC productUseCase
}

func New() *Handler {
	return &Handler{
		// GopayUC: gopayUC,
	}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	// span, ctx := tracer.StartFromRequest(r)
	// defer span.Finish()

	// // params checking
	// user_id, err := strconv.ParseInt(chi.URLParam(r, "user_id"), 10, 64)
	// if err != nil {
	// 	log.Error("[gopay.GetByUserID] error from Parse Param: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// // call the usecase
	// gopay, err := h.GopayUC.GetByUserID(ctx, user_id)
	// if err != nil {
	// 	log.Error("[gopay.GetByUserID] error from GetByUserID: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	// 	response.WriteJSONAPIError(w, r, http.StatusInternalServerError, `error get gopay`)
	// 	return
	// }

	// send the response
	products := []product.Product{
		{
			ProductID:    "1asdasd",
			UserID:       "1asssasd",
			Name:         "Waifu elit",
			ImageURL:     "https://www.seekpng.com/ipng/u2q8o0i1u2o0y3w7_mai-only-waifu-lucky-star-kagami-png/",
			Description:  "Image waifu",
			StartBid:     100000,
			MultipleBid:  3000,
			StartTime:    1647253367,
			EndTime:      1647771767,
			HighestBidID: "asda",
			TotalBidder:  33,
		},
		{
			ProductID:    "w23",
			UserID:       "sss",
			Name:         "Waifu elit",
			ImageURL:     "https://www.seekpng.com/ipng/u2q8o0i1u2o0y3w7_mai-only-waifu-lucky-star-kagami-png/",
			Description:  "Image waifu",
			StartBid:     99999,
			MultipleBid:  333,
			StartTime:    1647253367,
			EndTime:      1647771799,
			HighestBidID: "9sd",
			TotalBidder:  343,
		},
		{
			ProductID:    "1asdasd",
			UserID:       "1asssdasd",
			Name:         "Waifu elit",
			ImageURL:     "https://www.seekpng.com/ipng/u2q8o0i1u2o0y3w7_mai-only-waifu-lucky-star-kagami-png/",
			Description:  "Image waifu",
			StartBid:     12345,
			MultipleBid:  900,
			StartTime:    1647253367,
			EndTime:      1647771000,
			HighestBidID: "7",
			TotalBidder:  323,
		},
	}
	if _, err := response.WriteJSONAPIData(w, r, http.StatusOK, products); err != nil {
		log.Error("[gopay.GetByUserID] error from WriteJSON: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	}
}

func (h *Handler) GetAllBySeller(w http.ResponseWriter, r *http.Request) {
	// span, ctx := tracer.StartFromRequest(r)
	// defer span.Finish()

	// // params checking
	// user_id, err := strconv.ParseInt(chi.URLParam(r, "user_id"), 10, 64)
	// if err != nil {
	// 	log.Error("[gopay.GetByUserID] error from Parse Param: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// // call the usecase
	// gopay, err := h.GopayUC.GetByUserID(ctx, user_id)
	// if err != nil {
	// 	log.Error("[gopay.GetByUserID] error from GetByUserID: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	// 	response.WriteJSONAPIError(w, r, http.StatusInternalServerError, `error get gopay`)
	// 	return
	// }

	// send the response
	products := []product.Product{
		{
			ProductID:    "1asdasd",
			UserID:       "1asssasd",
			Name:         "Waifu elit",
			ImageURL:     "https://www.seekpng.com/ipng/u2q8o0i1u2o0y3w7_mai-only-waifu-lucky-star-kagami-png/",
			Description:  "Image waifu",
			StartBid:     100000,
			MultipleBid:  3000,
			StartTime:    1647253367,
			EndTime:      1647771767,
			HighestBidID: "asda",
			TotalBidder:  33,
		},
		{
			ProductID:    "1asdasd",
			UserID:       "1asssdasd",
			Name:         "Waifu elit",
			ImageURL:     "https://www.seekpng.com/ipng/u2q8o0i1u2o0y3w7_mai-only-waifu-lucky-star-kagami-png/",
			Description:  "Image waifu",
			StartBid:     12345,
			MultipleBid:  900,
			StartTime:    1647253367,
			EndTime:      1647771000,
			HighestBidID: "7",
			TotalBidder:  323,
		},
	}
	if _, err := response.WriteJSONAPIData(w, r, http.StatusOK, products); err != nil {
		log.Error("[gopay.GetByUserID] error from WriteJSON: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	}
}

func (h *Handler) GetAllByBuyer(w http.ResponseWriter, r *http.Request) {
	// span, ctx := tracer.StartFromRequest(r)
	// defer span.Finish()

	// // params checking
	// user_id, err := strconv.ParseInt(chi.URLParam(r, "user_id"), 10, 64)
	// if err != nil {
	// 	log.Error("[gopay.GetByUserID] error from Parse Param: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// // call the usecase
	// gopay, err := h.GopayUC.GetByUserID(ctx, user_id)
	// if err != nil {
	// 	log.Error("[gopay.GetByUserID] error from GetByUserID: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	// 	response.WriteJSONAPIError(w, r, http.StatusInternalServerError, `error get gopay`)
	// 	return
	// }

	// send the response
	products := []product.Product{
		{
			ProductID:    "w23",
			UserID:       "sss",
			Name:         "Waifu elit",
			ImageURL:     "https://www.seekpng.com/ipng/u2q8o0i1u2o0y3w7_mai-only-waifu-lucky-star-kagami-png/",
			Description:  "Image waifu",
			StartBid:     99999,
			MultipleBid:  333,
			StartTime:    1647253367,
			EndTime:      1647771799,
			HighestBidID: "9sd",
			TotalBidder:  343,
		},
	}
	if _, err := response.WriteJSONAPIData(w, r, http.StatusOK, products); err != nil {
		log.Error("[gopay.GetByUserID] error from WriteJSON: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	}
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	// span, ctx := tracer.StartFromRequest(r)
	// defer span.Finish()

	// params checking
	productID := chi.URLParam(r, "product_id")
	if productID == "" {
		log.Error("[product.GetByID] error Invalid param")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// // call the usecase
	// gopay, err := h.GopayUC.GetByUserID(ctx, user_id)
	// if err != nil {
	// 	log.Error("[gopay.GetByUserID] error from GetByUserID: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	// 	response.WriteJSONAPIError(w, r, http.StatusInternalServerError, `error get gopay`)
	// 	return
	// }

	// send the response
	product := product.Product{
		ProductID:    productID,
		UserID:       "sdskknf",
		Name:         "Waifu elit",
		ImageURL:     "https://www.seekpng.com/ipng/u2q8o0i1u2o0y3w7_mai-only-waifu-lucky-star-kagami-png/",
		Description:  "Image waifu",
		StartBid:     100000,
		MultipleBid:  3000,
		StartTime:    1647253367,
		EndTime:      1647771767,
		HighestBidID: "7",
		TotalBidder:  33,
	}

	if _, err := response.WriteJSONAPIData(w, r, http.StatusOK, product); err != nil {
		log.Error("[gopay.GetByUserID] error from WriteJSON: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	// span, ctx := tracer.StartFromRequest(r)
	// defer span.Finish()

	var product product.Product

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

	// // call the usecase
	// gopay, err := h.GopayUC.GetByUserID(ctx, user_id)
	// if err != nil {
	// 	log.Error("[gopay.GetByUserID] error from GetByUserID: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	// 	response.WriteJSONAPIError(w, r, http.StatusInternalServerError, `error get gopay`)
	// 	return
	// }

	// send the response
	if _, err := response.WriteJSONAPIData(w, r, http.StatusOK, product); err != nil {
		log.Error("[gopay.GetByUserID] error from WriteJSON: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	}
}

func (h *Handler) Bid(w http.ResponseWriter, r *http.Request) {
	// span, ctx := tracer.StartFromRequest(r)
	// defer span.Finish()

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

	// // call the usecase
	// gopay, err := h.GopayUC.GetByUserID(ctx, user_id)
	// if err != nil {
	// 	log.Error("[gopay.GetByUserID] error from GetByUserID: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	// 	response.WriteJSONAPIError(w, r, http.StatusInternalServerError, `error get gopay`)
	// 	return
	// }

	bid.BidID = "adasds"

	// send the response
	if _, err := response.WriteJSONAPIData(w, r, http.StatusOK, bid); err != nil {
		log.Error("[gopay.GetByUserID] error from WriteJSON: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	}
}
