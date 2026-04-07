package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/nekoimi/go-project-template/internal/dto"
	"github.com/nekoimi/go-project-template/internal/pkg/errcode"
	"github.com/nekoimi/go-project-template/internal/pkg/response"
	"github.com/nekoimi/go-project-template/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary      Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.RegisterRequest  true  "Register request"
// @Success      200   {object}  response.APIResponse{data=dto.AuthResponse}
// @Failure      400   {object}  response.APIResponse
// @Failure      409   {object}  response.APIResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		if appErr, ok := response.IsAppError(err); ok {
			response.AppErr(c, appErr)
			return
		}
		response.ErrorWithMsg(c, http.StatusInternalServerError, errcode.Internal, err.Error())
		return
	}

	response.Success(c, result)
}

// Login godoc
// @Summary      User login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.LoginRequest  true  "Login request"
// @Success      200   {object}  response.APIResponse{data=dto.AuthResponse}
// @Failure      400   {object}  response.APIResponse
// @Failure      401   {object}  response.APIResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		if appErr, ok := response.IsAppError(err); ok {
			response.AppErr(c, appErr)
			return
		}
		response.ErrorWithMsg(c, http.StatusInternalServerError, errcode.Internal, err.Error())
		return
	}

	response.Success(c, result)
}
