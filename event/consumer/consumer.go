package consumer

type Consumer interface {
	Listen(url string)
}
