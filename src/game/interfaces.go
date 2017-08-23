package game

// Notification is a JSONRPC notification
type Notification struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

// Subscribable describes a type which can send notifications
type Subscribable interface {
	Subscribe(id uint64) chan Notification
	Unsubscribe(id uint64)
}
