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
	UpdateBidFRDB(ctx context.Context, bid bid.Bid) error
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

	//// params checking
	//user_id, err := strconv.ParseInt(chi.URLParam(r, "user_id"), 10, 64)
	//if err != nil {
	//	log.Error("[bid.GetByID] error from Parse Param: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}

	// call the usecase
	// for test
	err := h.BidUC.UpdateBidFRDB(ctx, bid.Bid{BidID: "1", Amount: 100})
	if err != nil {
		log.Error("[bid.PublishBidFRDB] error from PublishBidFRDB: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
		response.WriteJSONAPIError(w, r, http.StatusInternalServerError, `error get bid`)
		return
	}

	// send the response
	if _, err := response.WriteJSONAPIData(w, r, http.StatusOK, nil); err != nil {
		log.Error("[bid.PublishBidFRDB] error from WriteJSON: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	}
}
