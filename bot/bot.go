package bots

// Bot interface
type Bot interface {
	Receive() <-chan string
	Emit(Message) error
}
