package main

type SaveMessageDTO struct {
	Text string `json:"message"`
}

type MessageStatisticsDTO struct {
	Processing uint64 `json:"processing"`
	Completed  uint64 `json:"completed"`
}
