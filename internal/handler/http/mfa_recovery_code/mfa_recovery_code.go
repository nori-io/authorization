package mfa_recovery_code

import (
	"net/http"

	"github.com/nori-io/common/v4/pkg/domain/logger"

	"github.com/nori-plugins/authentication/internal/domain/entity"
	"github.com/nori-plugins/authentication/internal/domain/service"
)

type MfaRecoveryCodeHandler struct {
	mfaRecoveryCodeService service.MfaRecoveryCodeService
	logger                 logger.FieldLogger
}

type Params struct {
	mfaRecoveryCodeService service.MfaRecoveryCodeService
	logger                 logger.FieldLogger
}

func New(params Params) *MfaRecoveryCodeHandler {
	return &MfaRecoveryCodeHandler{
		mfaRecoveryCodeService: params.mfaRecoveryCodeService,
		logger:                 params.logger,
	}
}

func (h *MfaRecoveryCodeHandler) GetMfaRecoveryCodes(w http.ResponseWriter, r *http.Request) {
	sessionId, err := r.Cookie("ssid")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	if _, err := h.mfaRecoveryCodeService.GetMfaRecoveryCodes(r.Context(), &entity.Session{SessionKey: []byte(sessionId.Value)}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//@todo path
	http.Redirect(w, r, "/", 0)
}
