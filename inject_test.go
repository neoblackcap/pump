package pump_test

import (
	"github.com/stretchr/testify/assert"
	"pump"
	"testing"
)

type Foo struct {
	Bar string `inject:"haha"`
}

type Bar struct {
	baz string `inject:"good"`
}

type F struct {
	f map[string]string `inject:"channel"`
}

type A struct {
	b []string `inject:"slice"`
}

type B struct {
	b chan<- string `inject:"slice"`
}

func TestWire(t *testing.T) {
	container := pump.NewContainer()
	foo := new(Foo)
	container.Register("haha", "string")
	container.Wire(&foo)
	assert.Equal(t, foo.Bar, "string")
}

func TestWireWithPrivateField(t *testing.T) {
	container := pump.NewContainer()
	bar := new(Bar)
	container.Register("good", "string")
	container.Wire(&bar)
	assert.Equal(t, bar.baz, "string")
}

func TestWireWithMapField(t *testing.T) {
	container := pump.NewContainer()
	f := new(F)
	c := make(map[string]string)
	c["foo"] = "bar"
	container.Register("channel", c)
	container.Wire(&f)
	assert.IsType(t, f.f, c)
	assert.Equal(t, f.f["foo"], "bar")
}

func TestWireWithArrayField(t *testing.T) {
	container := pump.NewContainer()
	f := new(A)
	c := []string{"foo", "bar"}

	container.Register("slice", c)
	container.Wire(&f)
	assert.IsType(t, f.b, c)
	assert.Equal(t, f.b[0], "foo")
}

func TestWireWithChanField(t *testing.T) {
	container := pump.NewContainer()
	f := new(B)
	c := make(chan<- string)

	container.Register("slice", c)
	container.Wire(&f)
	assert.Equal(t, f.b, c)
}
