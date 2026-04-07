package orderdto

import "time"

type GetInput struct {
	User_id  string `json:"user_id"`
	Order_id string `json:"order_id"`
}

type GetOutput struct {
	Order_id     string
	Order_status string
}

type StreamOutput struct {
	Order_status string
	Update_time  time.Time
}

type Output struct {
	Order_id     string
	Order_status string
}

func NewOutput(order_id, order_status string) Output {
	return Output{order_id, order_status}
}
