package http

import (
	"context"
	"net/http"

	"github.com/faithol1024/bgp-hackaton/internal/entity/bid"
	ers "github.com/faithol1024/bgp-hackaton/lib/error"
	"github.com/tokopedia/tdk/go/httpt/response"
	"github.com/tokopedia/tdk/go/log"
	"github.com/tokopedia/tdk/go/tracer"
)

type bidUseCase interface {
	UpdateBidFRDB(ctx context.Context, bid bid.BidFirebaseRDB) error
}

type Handler struct {
	BidUC bidUseCase
}

func New(bidUC bidUseCase) *Handler {
	return &Handler{
		BidUC: bidUC,
	}
}

func (h *Handler) PublishBidFRDB(w http.ResponseWriter, r *http.Request) {
	span, ctx := tracer.StartFromRequest(r)
	defer span.Finish()

	//Allow CORS here By *
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	//// params checking
	//user_id, err := strconv.ParseInt(chi.URLParam(r, "user_id"), 10, 64)
	//if err != nil {
	//	log.Error("[bid.GetByID] error from Parse Param: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}

	// call the usecase
	// for test
	err := h.BidUC.UpdateBidFRDB(ctx, bid.BidFirebaseRDB{
		ProductID:    "af447cba-6220-4092-a46a-ceaf52df16d9",
		UserID:       "1cbe7199-646f-4d45-80b9-0b071f538791",
		CurrentPrice: 100000,
		BidderCount:  1,
	})
	if err != nil {
		log.Error("[bid.PublishBidFRDB] error from PublishBidFRDB: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
		response.WriteJSONAPIError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// send the response
	if _, err := response.WriteJSONAPIData(w, r, http.StatusOK, nil); err != nil {
		log.Error("[bid.PublishBidFRDB] error from WriteJSON: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	}
}
