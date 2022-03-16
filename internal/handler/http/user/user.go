package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/faithol1024/bgp-hackaton/internal/entity/user"
	ers "github.com/faithol1024/bgp-hackaton/lib/error"
	"github.com/go-chi/chi"
	"github.com/tokopedia/tdk/go/httpt/response"
	"github.com/tokopedia/tdk/go/log"
	"github.com/tokopedia/tdk/go/tracer"
)

type userUseCase interface {
	GetByID(ctx context.Context, userID string) (user.User, error)
	Create(ctx context.Context, user user.User) (user.User, error)
}

type Handler struct {
	UserUC userUseCase
}

func New(userUC userUseCase) *Handler {
	return &Handler{
		UserUC: userUC,
	}
}

func (h *Handler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	span, ctx := tracer.StartFromRequest(r)
	defer span.Finish()

	// params checking
	userID := chi.URLParam(r, "user_id")
	if userID == "" {
		log.Error("[user.GetByID] error Invalid param")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// call the usecase
	user, err := h.UserUC.GetByID(ctx, userID)
	if err != nil {
		log.Error("[user.GetByUserID] error from GetByID: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
		response.WriteJSONAPIError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// send the response
	if _, err := response.WriteJSONAPIData(w, r, http.StatusOK, user); err != nil {
		log.Error("[user.GetByUserID] error from WriteJSON: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	span, ctx := tracer.StartFromRequest(r)
	defer span.Finish()

	var user user.User

	// params decode
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Errorf("[user.Create] failed to decode request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = user.Validate()
	if err != nil {
		log.Error("[user.Create] Invalid param")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// call the usecase
	user, err = h.UserUC.Create(ctx, user)
	if err != nil {
		log.Error("[user.Create] error from Create: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
		response.WriteJSONAPIError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// send the response
	if _, err := response.WriteJSONAPIData(w, r, http.StatusOK, user); err != nil {
		log.Error("[gopay.Create] error from WriteJSON: ", ers.ErrorAddTrace(err), ers.ErrorGetTrace(err))
	}
}
