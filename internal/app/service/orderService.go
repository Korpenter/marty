package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Mldlr/marty/internal/app/config"
	"github.com/Mldlr/marty/internal/app/models"
	"github.com/Mldlr/marty/internal/app/storage"
	"log"
	"net/http"
	"strconv"
	"time"
)

type OrderService interface {
	GetAccrual(order *models.Order)
	AddOrder(ctx context.Context, order *models.Order) error
	GetOrdersByUser(ctx context.Context, login string) ([]models.OrderItem, error)
	PollAccrual()
}

type OrderServiceImpl struct {
	repo         storage.Repository
	cfg          *config.Config
	updateQueue  chan *models.Order
	accrualQueue chan *models.Order
	httpClient   http.Client
	accrual      string
	Queue        chan *models.Order
}

func NewOrderService(c *config.Config, repo storage.Repository) OrderService {
	return &OrderServiceImpl{
		repo:         repo,
		cfg:          c,
		updateQueue:  make(chan *models.Order, 1000),
		accrualQueue: make(chan *models.Order, 1000),
		httpClient:   http.Client{},
		accrual:      fmt.Sprintf("%s/api/orders/", c.AccrualAddress),
		Queue:        make(chan *models.Order, 1000),
	}
}

func (s *OrderServiceImpl) PollAccrual() {
	for {
		order := <-s.accrualQueue
		gotOrder, retryAfter, err := s.getAccrual(order)
		if err != nil {
			log.Println(err)
			s.accrualQueue <- order
			time.Sleep(time.Duration(retryAfter) * time.Second)
			continue
		}
		if gotOrder.Status != order.Status {
			order.Status = gotOrder.Status
			s.updateQueue <- order
		}
		if order.Status == "PROCESSING" || order.Status == "REGISTERED" || order.Status == "" {
			s.accrualQueue <- order
		}
	}
}

func (s *OrderServiceImpl) getAccrual(order *models.Order) (*models.Order, int, error) {
	r, err := s.httpClient.Get(s.accrual + order.OrderID)
	log.Println("getting accrual")
	if err != nil {
		return nil, 0, err
	}
	switch {
	case r.StatusCode == http.StatusTooManyRequests:
		retryAfterHeader := r.Header.Get("Retry-After")
		retryAfter, _ := strconv.Atoi(retryAfterHeader)
		return nil, retryAfter, fmt.Errorf("too many requests")
	case r.StatusCode == http.StatusInternalServerError:
		return nil, 0, fmt.Errorf("accrual server error")
	}
	var gotOrder models.Order
	if err = json.NewDecoder(r.Body).Decode(&gotOrder); err != nil {
		return nil, 0, fmt.Errorf("error decoding json: %s", err)
	}
	defer r.Body.Close()
	log.Println(gotOrder)
	return &gotOrder, 0, nil
}

func (s *OrderServiceImpl) AddOrder(ctx context.Context, order *models.Order) error {
	err := s.repo.AddOrder(ctx, order)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *OrderServiceImpl) GetAccrual(order *models.Order) {
	s.accrualQueue <- order
}

func (s *OrderServiceImpl) GetOrdersByUser(ctx context.Context, login string) ([]models.OrderItem, error) {
	orders, err := s.repo.GetOrdersByUser(ctx, login)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *OrderServiceImpl) updateOrders(ctx context.Context) {
	for {
		order := <-s.updateQueue
		err := s.repo.UpdateOrder(ctx, order)
		if err != nil {
			s.updateQueue <- order
			log.Println("error updating order ", err)
		}
	}
}
