package types

import "testing"

func TestWei(t *testing.T) {
	var w *Wei
	num := 1.11
	w = NewWeiFromGwei(num)

	t.Log(w.String())
	if w.String() != "1110000000" {
		t.Fail()
	}

	w = NewWeiFromEther(num)

	t.Log(w.String())
	if w.String() != "1110000000000000000" {
		t.Fail()
	}

	if w.ToEther() != num {
		t.Log(w.ToEther())
		t.Fail()
	}

}
