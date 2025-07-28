package api

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

// GetTokenBalance returns the BOGO token balance for a specific address
// @Summary Get BOGO token balance
// @Description Returns the balance of BOGO tokens for a given address
// @Tags Tokens
// @Param address path string true "Wallet address"
// @Success 200 {object} sdk.TokenBalance
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /token/balance/{address} [get]
func (h *Handler) GetTokenBalance(c *gin.Context) {
	address := c.Param("address")

	// Validate Ethereum address
	if !common.IsHexAddress(address) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid Ethereum address"})
		return
	}

	balance, err := h.SDK.GetTokenBalance(address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, balance)
}

