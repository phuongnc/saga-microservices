package src

import (
	"context"
	"net/http"

	"infra/common/log"
	"infra/order"

	"github.com/labstack/echo/v4"
)

type KitchenHandler interface {
	RegisterEndpoints(echo *echo.Group)
	prepareOrderFailed(c echo.Context) error
	prepareOrderSuccess(c echo.Context) error
	getOrder(c echo.Context) error
}

func NewKitchenHandler(
	logger *log.Logger,
	kitchenService KitchenService,
) KitchenHandler {
	return &kitchenHandler{
		logger:         logger,
		kitchenService: kitchenService,
	}
}

type kitchenHandler struct {
	logger         *log.Logger
	kitchenService KitchenService
}

func (rc *kitchenHandler) RegisterEndpoints(echo *echo.Group) {
	echo.POST("/failed", rc.prepareOrderFailed)
	echo.POST("/success", rc.prepareOrderSuccess)
	echo.GET("/orders/:orderId", rc.getOrder)
}

func (rc *kitchenHandler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, "Ok")
}

func (rc *kitchenHandler) prepareOrderFailed(c echo.Context) error {
	req := &order.OrderDto{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Params")
	}

	model := &order.OrderModel{
		Id:            req.Id,
		Status:        order.ORDER_FAILED,
		SubStatus:     order.ORDER_KTCHENT_PREPARATION_FAILED,
		FailureReason: req.FailureReason,
	}
	ctx := context.WithValue(context.Background(), "db", c.Get("db"))
	_, err := rc.kitchenService.UpdateOrderKitchenStatus(ctx, model)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, nil)
}

func (rc *kitchenHandler) prepareOrderSuccess(c echo.Context) error {
	req := &order.OrderDto{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Params")
	}

	ctx := context.WithValue(context.Background(), "db", c.Get("db"))
	model := &order.OrderModel{
		Id:        req.Id,
		Status:    order.ORDER_DONE,
		SubStatus: order.ORDER_DELIVERED,
	}
	_, err := rc.kitchenService.UpdateOrderKitchenStatus(ctx, model)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, nil)
}

func (rc *kitchenHandler) getOrder(c echo.Context) error {
	orderId := c.Param("orderId")
	if orderId == "" {
		return c.JSON(http.StatusBadRequest, "Invalid Params")
	}
	ctx := context.WithValue(context.Background(), "db", c.Get("db"))
	existingOrder, err := rc.kitchenService.GetOrder(ctx, orderId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}
	if existingOrder == nil {
		return c.JSON(http.StatusNotFound, "Order is not exist")
	}
	return c.JSON(http.StatusOK, order.MapOrderFromModel(existingOrder))
}
