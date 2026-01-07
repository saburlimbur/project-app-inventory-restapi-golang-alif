package service

import (
	"alfdwirhmn/inventory/database"
	"alfdwirhmn/inventory/dto"
	"alfdwirhmn/inventory/model"
	"alfdwirhmn/inventory/repository"
	"alfdwirhmn/inventory/utils"
	"context"
	"errors"

	"go.uber.org/zap"
)

type SaleService interface {
	Create(ctx context.Context, usr *model.User, req dto.CreateSaleRequest) (*model.Sale, error)
	FindAll(ctx context.Context, page, limit int) (*[]model.Sale, *dto.Pagination, error)
	Detail(ctx context.Context, usr *model.User, id int) (*model.Sale, error)
	Update(ctx context.Context, usr *model.User, sale *model.Sale) error
	Delete(ctx context.Context, usr *model.User, id int) error
}

type saleService struct {
	repo     repository.SaleRepository
	itemRepo repository.ItemsRepository
	txMgr    database.TxManager // transaction db
	permSvc  PermissionService
	log      *zap.Logger
}

func NewSaleService(repo repository.SaleRepository, itemRepo repository.ItemsRepository, permSvc PermissionService, tx database.TxManager, log *zap.Logger) SaleService {
	return &saleService{
		repo:     repo,
		itemRepo: itemRepo,
		permSvc:  permSvc,
		txMgr:    tx,
		log:      log,
	}
}

func (s *saleService) Create(ctx context.Context, usr *model.User, req dto.CreateSaleRequest) (*model.Sale, error) {
	if !s.permSvc.CanCreateMasterData(usr.Role) {
		return nil, errors.New("forbidden")
	}

	// start transtaction
	tx, err := s.txMgr.Begin(ctx)
	if err != nil {
		return nil, err
	}

	// rollback if error
	defer tx.Rollback(ctx)

	// new repo with transaction context
	saleRepo := repository.NewSaleRepository(tx, s.log)
	itemRepo := repository.NewItemsRepository(tx, s.log)

	var total float64

	// loop items
	for _, it := range req.Items {
		// find from db
		item, err := itemRepo.FindByID(ctx, it.ItemID)
		if err != nil {
			return nil, err
		}

		if item.Stock < it.Quantity {
			return nil, errors.New("stock not enough")
		}

		// subtotal per item = harga * qty - discount
		sub := (item.Price * float64(it.Quantity)) - it.Discount
		total += sub
	}

	grandTotal := total - req.Discount + req.Tax

	sale := &model.Sale{
		InvoiceNumber: req.InvoiceNumber,
		CustomerName:  &req.CustomerName,
		CustomerPhone: &req.CustomerPhone,
		CustomerEmail: &req.CustomerEmail,
		TotalAmount:   total,
		Discount:      req.Discount,
		Tax:           req.Tax,
		GrandTotal:    grandTotal,
		PaymentMethod: &req.PaymentMethod,
		Notes:         &req.Notes,
		CreatedBy:     &usr.ID,
	}

	// save sale header to DB
	sale, err = saleRepo.Create(ctx, sale)
	if err != nil {
		return nil, err
	}

	// loop back, for save sale_item (detaill)
	for _, it := range req.Items {
		item, _ := itemRepo.FindByID(ctx, it.ItemID)

		sub := (item.Price * float64(it.Quantity)) - it.Discount

		err = saleRepo.CreateItem(ctx, &model.SaleItem{
			SaleID:    sale.ID,
			ItemID:    it.ItemID,
			Quantity:  it.Quantity,
			UnitPrice: item.Price,
			Subtotal:  sub,
			Discount:  it.Discount,
		})
		if err != nil {
			return nil, err
		}

		// kurangin stok items berdasarkan qty
		itemRepo.ReduceStock(ctx, it.ItemID, it.Quantity)
	}

	// commit transaksi
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return sale, nil
}

func (s *saleService) FindAll(ctx context.Context, page, limit int) (*[]model.Sale, *dto.Pagination, error) {
	sale, total, err := s.repo.Lists(ctx, page, limit)
	if err != nil {
		return nil, nil, err
	}

	pagination := dto.Pagination{
		Page:       page,
		Limit:      limit,
		TotalPages: utils.TotalPage(limit, int64(total)),
		TotalRows:  total,
	}

	return &sale, &pagination, nil
}

func (s *saleService) Detail(ctx context.Context, usr *model.User, id int) (*model.Sale, error) {
	if !s.permSvc.CanViewSale(usr.Role) {
		return nil, errors.New("forbidden")
	}
	return s.repo.FindDetailByID(ctx, id)
}

func (s *saleService) Update(ctx context.Context, usr *model.User, sale *model.Sale) error {
	if !s.permSvc.CanUpdateSale(usr.Role) {
		return errors.New("forbidden")
	}
	return s.repo.Update(ctx, sale)
}

func (s *saleService) Delete(ctx context.Context, usr *model.User, id int) error {
	if !s.permSvc.CanDeleteSale(usr.Role) {
		return errors.New("forbidden")
	}
	return s.repo.Delete(ctx, id)
}
