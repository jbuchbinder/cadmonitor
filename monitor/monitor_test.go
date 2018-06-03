package monitor

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

	/*
		ret, err := m.GetStatus("http://cadview.qvec.org/NewWorld.CAD.ViewOnly/CFSDetail.aspx?cfsID=-153484")
		if err != nil {
			t.Error(err)
			t.Fail()
		}
		t.Logf("%#v", ret)
	*/

	calls, err := m.GetClearedCalls("05/13/2018")
	//calls, err := m.GetActiveCalls()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Logf("#%v", calls)
	for k, v := range calls {
		cs, err := m.GetStatus(v)
		if err != nil {
			t.Logf("%s : %s", k, err.Error())
			continue
		}
		t.Logf("%s : %#v", k, cs)
	}
}
