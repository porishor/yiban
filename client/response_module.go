package client

type YiBanResponse[T interface{}] struct {
	Code int64   `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}
