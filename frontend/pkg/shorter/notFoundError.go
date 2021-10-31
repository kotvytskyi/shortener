package shorter

type NotFoundError struct{}

func (err *NotFoundError) Error() string { return "Not found" }
