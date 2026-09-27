package agency

import (
	"firefly-airline/internal/api"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

//	@Summary		Create a new agency
//	@Description	Creates a new travel agency and generates an API key for it.
//	@Tags			Agencies
//	@Accept			json
//	@Produce		json
//	@Security		AirlineAPIKey
//	@Param			request	body		CreateAgencyRequest	true	"Agency information"
//	@Success		201		{object}	api.APIResponse[CreateAgencyResponse]
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/internal/agencies [post]
func (h *Handler) Create(c *echo.Context) error {
	var req CreateAgencyRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	req.Name = strings.TrimSpace(req.Name)

	if req.Name == "" {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"name is required",
		)
	}

	walletBalance := int64(0)

	if req.WalletBalance != nil {
		walletBalance = *req.WalletBalance
	}

	agency, apiKey, err := h.service.Create(c.Request().Context(), req.Name, walletBalance)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			err.Error(),
		)
	}

	return c.JSON(http.StatusCreated, api.Success(CreateAgencyResponse{
		ID:            agency.ID,
		Name:          agency.Name,
		WalletBalance: agency.WalletBalance,
		APIKey:        apiKey,
	}))
}

// Add balance
//
// @Summary Add balance to agency wallet
// @Tags Agencies
// @Security AgencyAPIKey
// @Accept json
// @Produce json
// @Param request body AddBalanceRequest true "Deposit amount"
// @Success 200 {object} api.APIResponse[AddBalanceResponse]
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /wallet/deposit [post]
func (h *Handler) AddBalance(c *echo.Context) error {
	var req AddBalanceRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	agencyID, ok := c.Get("agency_id").(int64)
	if !ok || agencyID <= 0 {
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			"invalid agency context",
		)
	}

	balance, err := h.service.AddBalance(
		c.Request().Context(),
		agencyID,
		req.Amount,
	)
	if err != nil {
		return err
	}

	return c.JSON(
		http.StatusOK,
		api.Success(AddBalanceResponse{
			WalletBalance: balance,
		}),
	)
}