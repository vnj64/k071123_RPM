package repositories

import (
	"gorm.io/gorm"
	"k071123/internal/services/order_service/domain/repositories"
)

type sqlPagination struct {
	count  *int
	offset *int
}

func (p *sqlPagination) SetCount(count int) repositories.Pagination {
	p.count = &count
	return p
}

func (p *sqlPagination) SetOffset(offset int) repositories.Pagination {
	p.offset = &offset
	return p
}

func (p *sqlPagination) query(tx *gorm.DB) *gorm.DB {
	if p.count != nil {
		tx = tx.Limit(*p.count)
	}

	if p.offset != nil {
		tx = tx.Offset(*p.offset)
	}

	return tx
}
