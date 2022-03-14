package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/faithol1024/bgp-hackaton/internal/entity/gopay"
	ers "github.com/faithol1024/bgp-hackaton/lib/error"
	"github.com/go-chi/chi"
	"github.com/tokopedia/tdk/go/httpt/response"
	"github.com/tokopedia/tdk/go/log"
	"github.com/tokopedia/tdk/go/tracer"
)

type gopayUseCase interface {
	GetByUserID(ctx context.Context, userID int64) (gopay.GopaySaldo, error)
	GetHistoryByUserID(ctx context.Context, userID int64) ([]gopay.GopayHistory, error)
}

type Handler struct {
	GopayUC gopayUseCase
}

func New(gopayUC gopayUseCase) *Handler {
	return &Handler{
		GopayUC: gopayUC,
	}
}

func (h *Handler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	span, ctx := tracer.StartFromRequest(r)
	defer span.Finish()

	// params checking
	userID := chi.URLParam(r, "user_id")
	if userID == "" {
		log.Error("[gopay.GetByUserID] Invalid Param")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// call the usecase
	gopay, err := h.GopayUC.GetByUserID(ctx, userID)
	if err != nil {
		log.Error("[gopay.GetByUserID] error from GetByUserID: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
		response.WriteJSONAPIError(w, r, http.StatusInternalServerError, `error get gopay`)
		return
	}

	// send the response
	if _, err := response.WriteJSONAPIData(w, r, http.StatusOK, gopay); err != nil {
		log.Error("[gopay.GetByUserID] error from WriteJSON: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	}
}

func (h *Handler) GetHistoryByUserID(w http.ResponseWriter, r *http.Request) {
	span, ctx := tracer.StartFromRequest(r)
	defer span.Finish()

	// params checking
	userID, err := strconv.ParseInt(chi.URLParam(r, "user_id"), 10, 64)
	if err != nil {
		log.Error("[gopay.GetHistoryByUserID] error from Parse Param: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// call the usecase
	gopayHistories, err := h.GopayUC.GetHistoryByUserID(ctx, userID)
	if err != nil {
		log.Error("[gopay.GetHistoryByUserID] error from GetHistoryByUserID: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
		response.WriteJSONAPIError(w, r, http.StatusInternalServerError, `error get gopay history`)
		return
	}

	// send the response
	if _, err := response.WriteJSONAPIData(w, r, http.StatusOK, gopayHistories); err != nil {
		log.Error("[gopay.GetHistoryByUserID] error from WriteJSON: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	}
}
