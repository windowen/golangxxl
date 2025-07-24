package strhelper

import (
	"testing"
)

func TestGenerator_Encode(t *testing.T) {
	g := NewGenerator(6)

	t.Log("最大支持ID:", g.MaxSupportID())

	test := func(id uint64) bool {
		code := g.Encode(id)
		t.Logf("ID:%d code:%s", id, code)
		nid := g.Decode(code)

		if nid != id {
			t.Error(id, nid)
			return false
		}
		return true
	}

	var _min, _max uint64 = 0, 20
	for id := _min; id <= _max; id++ {
		if !test(id) {
			return
		}
	}
}
