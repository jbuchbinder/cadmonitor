package main

import (
	"testing"
)

func Test_Monitor(t *testing.T) {
	m := CadBrowser{}
	err := m.Login(USER, PASS)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	ret, err := m.GetStatus("http://cadview.qvec.org/NewWorld.CAD.ViewOnly/CFSDetail.aspx?cfsID=-153484")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Logf("%#v", ret)

}
