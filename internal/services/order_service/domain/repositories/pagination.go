package repositories

type Pagination interface {
	SetCount(count int) Pagination
	SetOffset(offset int) Pagination
}
