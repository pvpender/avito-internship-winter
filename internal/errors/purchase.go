package errors

type PurchaseError struct {
}

func (e PurchaseError) Error() string {
	return "not enough coins"
}
