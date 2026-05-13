package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/dto"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/response"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/service"
)

type SubscriptionHandler struct {
	subscriptions *service.SubscriptionService
}

func NewSubscriptionHandler(subscriptions *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{subscriptions: subscriptions}
}

func (h *SubscriptionHandler) Create(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}

	var req dto.SubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid_request", "invalid subscription request")
		return
	}

	subscription, err := h.subscriptions.Create(c.Request.Context(), userID, req)
	if err != nil {
		writeSubscriptionError(c, err, "create_subscription_failed")
		return
	}

	response.Success(c, subscription)
}

func (h *SubscriptionHandler) List(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}

	query, ok := subscriptionQuery(c)
	if !ok {
		return
	}
	subscriptions, err := h.subscriptions.List(c.Request.Context(), userID, query)
	if err != nil {
		writeSubscriptionError(c, err, "list_subscriptions_failed")
		return
	}

	response.Success(c, subscriptions)
}

func (h *SubscriptionHandler) Detail(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}
	id, ok := subscriptionID(c)
	if !ok {
		return
	}

	subscription, err := h.subscriptions.Detail(c.Request.Context(), userID, id)
	if err != nil {
		writeSubscriptionError(c, err, "subscription_detail_failed")
		return
	}

	response.Success(c, subscription)
}

func (h *SubscriptionHandler) Update(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}
	id, ok := subscriptionID(c)
	if !ok {
		return
	}

	var req dto.SubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid_request", "invalid subscription request")
		return
	}

	subscription, err := h.subscriptions.Update(c.Request.Context(), userID, id, req)
	if err != nil {
		writeSubscriptionError(c, err, "update_subscription_failed")
		return
	}

	response.Success(c, subscription)
}

func (h *SubscriptionHandler) Delete(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}
	id, ok := subscriptionID(c)
	if !ok {
		return
	}

	subscription, err := h.subscriptions.Delete(c.Request.Context(), userID, id)
	if err != nil {
		writeSubscriptionError(c, err, "delete_subscription_failed")
		return
	}

	response.Success(c, subscription)
}

func subscriptionID(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid_subscription_id", "invalid subscription id")
		return uuid.Nil, false
	}

	return id, true
}

func subscriptionQuery(c *gin.Context) (service.SubscriptionQuery, bool) {
	page, ok := optionalPositiveInt(c, "page")
	if !ok {
		return service.SubscriptionQuery{}, false
	}
	pageSize, ok := optionalPositiveInt(c, "page_size")
	if !ok {
		return service.SubscriptionQuery{}, false
	}

	var categoryID *uuid.UUID
	if rawCategoryID := c.Query("category_id"); rawCategoryID != "" {
		parsed, err := uuid.Parse(rawCategoryID)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "invalid_request", "invalid subscription query")
			return service.SubscriptionQuery{}, false
		}
		categoryID = &parsed
	}

	return service.SubscriptionQuery{
		Page:         page,
		PageSize:     pageSize,
		Status:       c.Query("status"),
		CategoryID:   categoryID,
		BillingCycle: c.Query("billing_cycle"),
	}, true
}

func optionalPositiveInt(c *gin.Context, key string) (int, bool) {
	value := c.Query(key)
	if value == "" {
		return 0, true
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		response.Error(c, http.StatusBadRequest, "invalid_request", "invalid subscription query")
		return 0, false
	}

	return parsed, true
}

func writeSubscriptionError(c *gin.Context, err error, fallbackCode string) {
	switch {
	case errors.Is(err, service.ErrInvalidSubscription), errors.Is(err, service.ErrInvalidSubscriptionQry):
		response.Error(c, http.StatusBadRequest, "invalid_request", "invalid subscription request")
	case errors.Is(err, service.ErrSubscriptionCategory):
		response.Error(c, http.StatusNotFound, "category_not_found", "category not found")
	case errors.Is(err, service.ErrSubscriptionNotFound):
		response.Error(c, http.StatusNotFound, "subscription_not_found", "subscription not found")
	default:
		response.Error(c, http.StatusInternalServerError, fallbackCode, "subscription operation failed")
	}
}
