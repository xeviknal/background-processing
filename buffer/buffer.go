package buffer

import "time"

type Buffer interface {
	Add(t *Descriptor)
	Next() *Descriptor
	Size() int
	Close()
	IsClosed() bool
}

type TaskBuffer struct {
	// FIFO to append the background tasks
	Channel chan *Descriptor

	// Tx is the transmission rate
	// It is the number of tasks to process per second
	// If Tx is 10, then every task will be executed every 100ms
	Tx int

	// Whether the channel is closed
	Closed bool
}

var tb *Buffer
var defaultBufferCap = 100 // 100 maxim tasks, otherwise wait the execution
var defaultBufferTx = 10   // 10 tasks per second

func NewTaskBuffer(cap, tx int) Buffer {
	tb := &TaskBuffer{
		Channel: make(chan *Descriptor, cap),
		Tx:      tx,
	}
	return Buffer(tb)
}

func GetBuffer() Buffer {
	if tb != nil && !(*tb).IsClosed() {
		return *tb
	}

	ntb := NewTaskBuffer(defaultBufferCap, defaultBufferTx)
	tb = &ntb
	return *tb
}

func SetBufferSettings(cap, tx int) {
	defaultBufferCap = cap
	defaultBufferTx = 1000 / tx // 1 second split into TX tasks
}

func (tb *TaskBuffer) Add(t *Descriptor) {
	tb.Channel <- t
}

func (tb *TaskBuffer) Next() *Descriptor {
	next := <-tb.Channel
	time.Sleep(time.Duration(tb.Tx) * time.Millisecond)
	return next
}

func (tb *TaskBuffer) Size() int {
	return len(tb.Channel)
}

func (tb *TaskBuffer) Close() {
	close(tb.Channel)
	tb.Closed = true
}

func (tb *TaskBuffer) IsClosed() bool {
	return tb.Closed
}

type Descriptor interface {
	GetConsumer() string
	GetDescriptor() map[string]interface{}
}

type TaskDescriptor struct {
	Consumer   string
	Descriptor map[string]interface{}
}

func (td TaskDescriptor) GetConsumer() string {
	return td.Consumer
}

func (td TaskDescriptor) GetDescriptor() map[string]interface{} {
	return td.Descriptor
}
