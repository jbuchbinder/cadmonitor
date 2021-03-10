package monitor

import (
	"os"
	"testing"

	// Autoload .env files
	_ "github.com/joho/godotenv/autoload"
)

func Test_AegisMonitor(t *testing.T) {
	m := AegisMonitor{
		BaseURL: os.Getenv("BASEURL"),
		FDID:    os.Getenv("FDID"),
	}
	err := m.Login(os.Getenv("CADUSER"), os.Getenv("CADPASS"))
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

	calls, err := m.GetClearedCalls("01/21/2019")
	//calls, err := m.GetActiveCalls()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Logf("#%v", calls)
	for k, v := range calls {
		cs, err := m.GetStatusFromURL(v)
		if err != nil {
			t.Logf("%s : %s", k, err.Error())
			continue
		}
		t.Logf("%s : %#v", k, cs)
	}
}
