package payment_statuses

type PaymentStatuses string

const (
	Pending           PaymentStatuses = "pending"
	WaitingForCapture PaymentStatuses = "waiting_for_capture"
	Succeeded         PaymentStatuses = "succeeded"
	Canceled          PaymentStatuses = "canceled"
	Refunded          PaymentStatuses = "refunded"
)
