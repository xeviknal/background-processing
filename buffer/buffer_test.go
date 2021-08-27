package buffer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Creates a buffer as expected
func TestNewTaskBuffer(t *testing.T) {
	asserter := assert.New(t)
	buffer := NewTaskBuffer(10, 1)

	buf := buffer.(*TaskBuffer)
	asserter.Equal(cap(buf.Channel), 10)
	asserter.Equal(buf.Tx, 1)
}

// Always gets the same instance of the buffer
func TestGetSingletonBuffer(t *testing.T) {
	asserter := assert.New(t)
	buffer := GetBuffer()
	buf := buffer.(*TaskBuffer)
	asserter.Equal(100, cap(buf.Channel))
	asserter.Equal(10, buf.Tx)
}

// Adding task method works
func TestAddTask(t *testing.T) {
	asserter := assert.New(t)

	SetBufferSettings(3, 1)
	b := GetBuffer()

	td := Descriptor(TaskDescriptor{
		Consumer:   "consumer",
		Descriptor: "descriptor",
	})

	b.Add(&td)

	asserter.Equal(b.Size(), 1)

	tb := b.(*TaskBuffer)
	buffTask := <-tb.Channel
	asserter.Equal("consumer", (*buffTask).GetConsumer())
	asserter.Equal("descriptor", (*buffTask).GetDescriptor())

	td = Descriptor(TaskDescriptor{
		Consumer:   "consumer",
		Descriptor: "descriptor",
	})

	b.Add(&td)
	b.Add(&td)

	asserter.Equal(b.Size(), 2)
	tb.Close()
}

// Taking one task works
func TestNextTask(t *testing.T) {
	asserter := assert.New(t)

	SetBufferSettings(3, 1)
	b := GetBuffer()

	td := Descriptor(TaskDescriptor{
		Consumer:   "consumer",
		Descriptor: "descriptor",
	})

	b.Add(&td)
	b.Add(&td)
	b.Add(&td)
	asserter.Equal(b.Size(), 3)

	b.Next()
	asserter.Equal(b.Size(), 2)
	b.Next()
	asserter.Equal(b.Size(), 1)
	b.Add(&td)
	asserter.Equal(b.Size(), 2)

	b.Close()
}
