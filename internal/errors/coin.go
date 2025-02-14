package errors

type InvalidAmount struct {
}

func (e InvalidAmount) Error() string {
	return "amount must be positive"
}

type InvalidCoinOperation struct {
}

func (e InvalidCoinOperation) Error() string {
	return "coin operation cannot be processed, check user and amount"
}
