package server

import (
	"log/slog"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

func (s *Server) initAdminRoutes() {
	s.app.GET("/admin/dashboard/info", adminToken(), func(ctx *gin.Context) {
		products, err := s.usecases.Db.GetAllProducts()
		if err != nil {
			slog.Error("Failed to get all products", "error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error getting products"})
			return
		}
		transactions, err := s.usecases.Db.GetAllTransactions()
		if err != nil {
			slog.Error("Error getting all transactions", "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting transactions"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"products": products, "transactions": transactions})
	})

	s.app.GET("/admin/revenue", adminToken(), func(ctx *gin.Context) {
		transactions, err := s.usecases.Db.GetAllTransactions()
		if err != nil {
			slog.Error("Failed to build revenue report", "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting revenue"})
			return
		}

		byCurrency := make(map[string]*revenueCurrencySummary)
		for _, transaction := range transactions {
			paymentStatus := transaction.PaymentStatus
			if paymentStatus == "" {
				paymentStatus = transaction.Status
			}
			if paymentStatus != "paid" && paymentStatus != "refunded" && paymentStatus != "succeeded" {
				continue
			}
			currency := transaction.Currency
			if currency == "" {
				continue
			}
			summary := byCurrency[currency]
			if summary == nil {
				summary = &revenueCurrencySummary{Currency: currency}
				byCurrency[currency] = summary
			}
			summary.Gross += transaction.AmountTotal
			summary.Refunds += transaction.AmountRefunded
			summary.PaymentCount++
			if transaction.StripeFeeKnown {
				summary.StripeFees += transaction.StripeFee
			} else {
				summary.StripeFeePendingPayments++
			}
		}

		currencies := make([]string, 0, len(byCurrency))
		for currency := range byCurrency {
			currencies = append(currencies, currency)
		}
		sort.Strings(currencies)
		revenue := make([]revenueCurrencySummary, 0, len(currencies))
		for _, currency := range currencies {
			summary := *byCurrency[currency]
			summary.NetCompanyRevenue = summary.Gross - summary.Refunds - summary.StripeFees
			revenue = append(revenue, summary)
		}
		ctx.JSON(http.StatusOK, gin.H{"revenue": revenue})
	})
}

type revenueCurrencySummary struct {
	Currency                 string `json:"currency"`
	Gross                    int64  `json:"gross"`
	StripeFees               int64  `json:"stripe_fees"`
	Refunds                  int64  `json:"refunds"`
	NetCompanyRevenue        int64  `json:"net_company_revenue"`
	PaymentCount             int    `json:"payment_count"`
	StripeFeePendingPayments int    `json:"stripe_fee_pending_payments"`
}
