package errors

type InvalidAmountError struct {
}

func (e InvalidAmountError) Error() string {
	return "amount must be positive"
}

type InvalidCoinOperationError struct {
}

func (e InvalidCoinOperationError) Error() string {
	return "coin operation cannot be processed, check user and amount"
}
