package assert

import (
	"fmt"
	"testing"
)

type Equals int64

func (e Equals) CheckUint8(t *testing.T, desc string, val uint8) {
	e.CheckInt64(t, desc, int64(val))
}

func (e Equals) CheckUint16(t *testing.T, desc string, val uint16) {
	e.CheckInt64(t, desc, int64(val))
}

func (e Equals) CheckUint32(t *testing.T, desc string, val uint32) {
	e.CheckInt64(t, desc, int64(val))
}

func (e Equals) CheckUint64(t *testing.T, desc string, val uint64) {
	e.CheckInt64(t, desc, int64(val))
}

func (e Equals) CheckInt(t *testing.T, desc string, val int) {
	e.CheckInt64(t, desc, int64(val))
}

func (e Equals) CheckInt8(t *testing.T, desc string, val int8) {
	e.CheckInt64(t, desc, int64(val))
}

func (e Equals) CheckInt16(t *testing.T, desc string, val int16) {
	e.CheckInt64(t, desc, int64(val))
}

func (e Equals) CheckInt32(t *testing.T, desc string, val int32) {
	e.CheckInt64(t, desc, int64(val))
}

func (e Equals) CheckInt64(t *testing.T, desc string, val int64) {
	t.Run(fmt.Sprintf("%s equals %d", desc, e), func(t *testing.T) {
		expected := int64(e)
		if val != expected {
			t.Errorf("expected %d but got %d", expected, val)
		}
	})
}
