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
	user_id, err := strconv.ParseInt(chi.URLParam(r, "user_id"), 10, 64)
	if err != nil {
		log.Error("[gopay.GetByUserID] error from Parse Param: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// call the usecase
	gopay, err := h.GopayUC.GetByUserID(ctx, user_id)
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
