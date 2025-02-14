package errors

type PurchaseError struct {
	Msg string
}

func (e PurchaseError) Error() string {
	return e.Msg
}
