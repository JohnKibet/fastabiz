package order

type Repository interface {
	Create(order *Order) error
	GetByID(id int64) (*Order, error)
	ListByCustomer(customerID int64) ([]*Order, error)
	UpdateStatus(orderID int64, status OrderStatus) error
}
