package golog

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type (
	WriterMock struct {
		mock.Mock
	}

	FormaterMock struct {
		mock.Mock
	}
)

func (f *FormaterMock) Format(level int, format string, a ...interface{}) string {
	args := f.Called(level, format, a)
	return args.Get(0).(string)
}

func (w *WriterMock) Write(b []byte) (n int, err error) {
	args := w.Called(b)
	return args.Get(0).(int), args.Error(1)

}

func TestNew(t *testing.T) {
	w := new(WriterMock)
	f := new(FormaterMock)

	golog := New(DEBUG, w, f)

	assert.Equal(t, DEBUG, golog.Level)
	assert.Equal(t, w, golog.Writer)
	assert.Equal(t, f, golog.Formater)
}

func TestMethods(t *testing.T) {
	msg := "hello"
	w := new(WriterMock)
	w.On("Write", []byte(msg)).Return(len(msg), nil)

	f := new(FormaterMock)
	levels := []int{DEBUG, INFO, WARN, ERROR, FATAL}
	for _, level := range levels {
		f.On("Format", level, "foo", []interface{}(nil)).Return(msg)
	}

	golog := New(DEBUG, w, f)
	golog.Debug("foo")
	golog.Info("foo")
	golog.Warn("foo")
	golog.Error("foo")
	golog.Fatal("foo")
}

func TestComponentMethodsCalled(t *testing.T) {
	f := new(FormaterMock)
	f.On("Format", DEBUG, "foo %s", []interface{}{"bar"}).Return("hello")

	w := new(WriterMock)
	w.On("Write", []byte("hello")).Return(len("hello"), nil)

	golog := New(DEBUG, w, f)
	golog.Printf(DEBUG, "foo %s", "bar")

	f.AssertNumberOfCalls(t, "Format", 1)
	w.AssertNumberOfCalls(t, "Write", 1)
}

func TestCorrectLevel(t *testing.T) {
	f := new(FormaterMock)
	f.On("Format", DEBUG, "foo", []interface{}(nil)).Return("hello")
	f.On("Format", ERROR, "bar", []interface{}(nil)).Return("hello")
	f.On("Format", FATAL, "baz", []interface{}(nil)).Return("hello")

	w := new(WriterMock)
	w.On("Write", []byte("hello")).Return(len("hello"), nil)

	golog := New(WARN, w, f)
	golog.Debug("foo")
	golog.Error("bar")
	golog.Fatal("baz")

	f.AssertNumberOfCalls(t, "Format", 2)
	w.AssertNumberOfCalls(t, "Write", 2)
}
