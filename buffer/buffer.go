package buffer

type Buffer interface {
	Add(t *Descriptor)
	Next() *Descriptor
	Size() int
	Close()
}

type TaskBuffer struct {
	Channel chan *Descriptor
	Tx      int
}

var tb *Buffer
var defaultBufferCap = 10
var defaultBufferTx = 1

func NewTaskBuffer(cap, tx int) Buffer {
	tb := &TaskBuffer{
		Channel: make(chan *Descriptor, cap),
		Tx:      tx,
	}
	return Buffer(tb)
}

func GetBuffer() Buffer {
	if tb != nil {
		return *tb
	}

	ntb := NewTaskBuffer(defaultBufferCap, defaultBufferTx)
	tb = &ntb
	return *tb
}

func SetBufferSettings(cap, tx int) {
	defaultBufferCap = cap
	defaultBufferTx = tx
}

func (tb *TaskBuffer) Add(t *Descriptor) {
	tb.Channel <- t
}

func (tb *TaskBuffer) Next() *Descriptor {
	return <-tb.Channel
}

func (tb *TaskBuffer) Size() int {
	return len(tb.Channel)
}

func (tb *TaskBuffer) Close() {
	close(tb.Channel)
}

type Descriptor interface {
	GetConsumer() interface{}
	GetDescriptor() interface{}
}

type TaskDescriptor struct {
	Consumer   interface{}
	Descriptor interface{}
}

func (td TaskDescriptor) GetConsumer() interface{} {
	return td.Consumer
}

func (td TaskDescriptor) GetDescriptor() interface{} {
	return td.Descriptor
}
