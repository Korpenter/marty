package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Mldlr/marty/internal/app/config"
	"github.com/Mldlr/marty/internal/app/constant"
	"github.com/Mldlr/marty/internal/app/logging"
	"github.com/Mldlr/marty/internal/app/models"
	"github.com/Mldlr/marty/internal/app/storage"
	"net/http"
	"strconv"
	"time"
)

type OrderService interface {
	GetAccrual(order *models.Order)
	AddOrder(ctx context.Context, order *models.Order) error
	GetOrdersByUser(ctx context.Context) ([]models.OrderItem, error)
	UpdateOrders(ctx context.Context)
	PollAccrual()
}

type OrderServiceImpl struct {
	repo         storage.Repository
	cfg          *config.Config
	updateQueue  chan *models.Order
	accrualQueue chan *models.Order
	accrual      string
	Queue        chan *models.Order
}

func NewOrderService(cfg *config.Config, repo storage.Repository) OrderService {
	return &OrderServiceImpl{
		repo:         repo,
		cfg:          cfg,
		updateQueue:  make(chan *models.Order, 1000),
		accrualQueue: make(chan *models.Order, 1000),
		accrual:      fmt.Sprintf("http://%s/api/orders/", cfg.AccrualAddress),
		Queue:        make(chan *models.Order, 1000),
	}
}

func (s *OrderServiceImpl) PollAccrual() {
	for {
		order := <-s.accrualQueue
		gotOrder, retryAfter, err := s.getAccrual(order)
		if retryAfter == -1 {
			continue
		}
		if err != nil {
			logging.Logger.Error("constant polling accrual service:" + err.Error())
			s.accrualQueue <- order
			time.Sleep(time.Duration(retryAfter) * time.Second)
			continue
		}
		if gotOrder.Status != order.Status {
			order.Status = gotOrder.Status
			s.updateQueue <- gotOrder
		}
		if order.Status == constant.StatusProcessing || order.Status == constant.StatusRegistered || order.Status == "" {
			s.accrualQueue <- order
		}
	}
}

func (s *OrderServiceImpl) getAccrual(order *models.Order) (*models.Order, int, error) {
	r, err := http.Get(s.accrual + order.OrderID)
	if err != nil {
		return nil, 0, err
	}
	switch {
	case r.StatusCode == http.StatusNoContent:
		return nil, -1, nil
	case r.StatusCode == http.StatusTooManyRequests:
		retryAfterHeader := r.Header.Get("Retry-After")
		retryAfter, _ := strconv.Atoi(retryAfterHeader)
		return nil, retryAfter, fmt.Errorf("too many requests")
	case r.StatusCode == http.StatusInternalServerError:
		return nil, 0, fmt.Errorf("accrual server constant")
	}
	var gotOrder models.Order
	defer r.Body.Close()
	if err = json.NewDecoder(r.Body).Decode(&gotOrder); err != nil {
		return nil, 0, fmt.Errorf("constant decoding json: %s", err)
	}
	return &gotOrder, 0, nil
}

func (s *OrderServiceImpl) AddOrder(ctx context.Context, order *models.Order) error {
	err := s.repo.AddOrder(ctx, order)
	if err != nil {
		logging.Logger.Error(err.Error())
		return err
	}
	return nil
}

func (s *OrderServiceImpl) GetAccrual(order *models.Order) {
	s.accrualQueue <- order
}

func (s *OrderServiceImpl) GetOrdersByUser(ctx context.Context) ([]models.OrderItem, error) {
	orders, err := s.repo.GetOrdersByUser(ctx)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *OrderServiceImpl) UpdateOrders(ctx context.Context) {
	for {
		order := <-s.updateQueue
		err := s.repo.UpdateOrder(ctx, order)
		if err != nil {
			s.updateQueue <- order
			logging.Logger.Error("constant updating order :" + err.Error())
		}
	}
}
