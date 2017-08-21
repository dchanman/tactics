package game

type Notification struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type Subscribable interface {
	Subscribe(id uint64) chan Notification
	Unsubscribe(id uint64)
}
