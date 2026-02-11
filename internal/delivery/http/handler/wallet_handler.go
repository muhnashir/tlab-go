package handler

import (
	"wallet-api/internal/domain"
	"wallet-api/internal/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type WalletHandler struct {
	Service domain.TransactionService
}

func NewWalletHandler(s domain.TransactionService) *WalletHandler {
	return &WalletHandler{Service: s}
}

type TransferRequest struct {
	ReceiverUserID int64   `json:"receiver_user_id" validate:"required"`
	Amount         float64 `json:"amount" validate:"required,gt=0"`
}

// TopUp godoc
// @Summary Top up wallet
// @Description Top up wallet balance (mock implementation)
// @Tags Wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {string} string "TopUp Not Implemented"
// @Router /api/wallets/topup [post]
func (h *WalletHandler) TopUp(c *fiber.Ctx) error {
	return c.SendString("TopUp Not Implemented")
}

// Transfer godoc
// @Summary Transfer funds
// @Description Transfer funds from logged-in user to another user
// @Tags Wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body TransferRequest true "Transfer Request"
// @Success 200 {object} utils.ApiResponse{data=domain.Transaction}
// @Failure 400 {object} utils.ApiResponse
// @Failure 401 {object} utils.ApiResponse
// @Router /transactions/transfer [post]
func (h *WalletHandler) Transfer(c *fiber.Ctx) error {
	// Parse user_id from middleware
	userID, ok := c.Locals("user_id").(int64)
	if !ok {
		return utils.Unauthorized(c, "Invalid user session")
	}

	var req TransferRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body", err.Error())
	}

	transaction, err := h.Service.Transfer(c.Context(), userID, req.ReceiverUserID, req.Amount)
	if err != nil {
		// Differentiate errors if possible
		return utils.BadRequest(c, err.Error(), nil)
	}

	return utils.Success(c, fiber.StatusOK, "Transfer successful", transaction)
}

// GetHistory godoc
// @Summary Get transaction history
// @Description Get transaction history for logged-in user with pagination
// @Tags Wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} utils.ApiResponse{data=[]domain.Transaction}
// @Failure 401 {object} utils.ApiResponse
// @Failure 500 {object} utils.ApiResponse
// @Router /transactions [get]
func (h *WalletHandler) GetHistory(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int64)
	if !ok {
		return utils.Unauthorized(c, "Invalid user session")
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	history, err := h.Service.GetHistory(c.Context(), userID, page, limit)
	if err != nil {
		return utils.InternalServerError(c, "Failed to retrieve history", err.Error())
	}

	return utils.Success(c, fiber.StatusOK, "History retrieved", history)
}

// GetBalance godoc
// @Summary Get user balance
// @Description Get current balance of logged-in user
// @Tags Wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.ApiResponse{data=domain.Wallet}
// @Failure 401 {object} utils.ApiResponse
// @Failure 404 {object} utils.ApiResponse
// @Failure 500 {object} utils.ApiResponse
// @Router /wallets/balance [get]
func (h *WalletHandler) GetBalance(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int64)
	if !ok {
		return utils.Unauthorized(c, "Invalid user session")
	}

	wallet, err := h.Service.GetBalance(c.Context(), userID)
	if err != nil {
		return utils.InternalServerError(c, "Failed to retrieve balance", err.Error())
	}
	if wallet == nil {
		return utils.NotFound(c, "Wallet not found")
	}

	return utils.Success(c, fiber.StatusOK, "Balance retrieved", wallet)
}
