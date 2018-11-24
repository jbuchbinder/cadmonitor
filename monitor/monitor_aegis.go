package monitor

import (
	//"fmt"

	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/agent"
	"github.com/headzoo/surf/browser"
	"github.com/pkg/errors"
)

const (
	aegisClearedCallURL = "ClearedCallSearch.aspx"
	aegisLoginURL       = "Login.aspx"
	aegisMainURL        = "default.aspx"
)

func init() {
	RegisterCadMonitor("aegis", func() CadMonitor { return &AegisMonitor{} })
}

type AegisMonitor struct {
	// Suffix is the required suffix for units. Leaving blank disables
	// qualification by this value.
	Suffix string

	// LoginURL (example: "http://cadview.qvec.org/NewWorld.CAD.ViewOnly/")
	BaseURL string

	browserObject *browser.Browser
	initialized   bool
}

func (c *AegisMonitor) ConfigureFromValues(values map[string]string) error {
	var ok bool
	if c.Suffix, ok = values["suffix"]; !ok {
		return errors.New("'suffix' not defined")
	}
	if c.BaseURL, ok = values["baseUrl"]; !ok {
		return errors.New("'baseUrl' not defined")
	}
	return nil
}

func (c *AegisMonitor) Login(user, pass string) error {
	b := surf.NewBrowser()
	c.browserObject = b

	b.SetUserAgent(agent.Chrome())

	// Required to not have ASP.NET garbage yak all over me
	b.AddRequestHeader("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")

	err := b.Open(c.BaseURL + aegisLoginURL)
	if err != nil {
		return err
	}

	if len(b.Forms()) < 1 {
		return errors.New("Form does not exist")
	}
	f, err := b.Form("form#aspnetForm")
	if err != nil {
		return err
	}
	f.Input("ctl00$content$Login1$UserName", user)
	f.Input("ctl00$content$Login1$Password", pass)

	if f.Submit() != nil {
		return err
	}

	err = b.Bookmark("main")
	if err != nil {
		return errors.Wrap(err, "Unable to bookmark main")
	}

	c.initialized = true

	return nil
}

func (c *AegisMonitor) GetActiveCalls() ([]string, error) {
	calls := make([]string, 0)

	if !c.initialized {
		return calls, errors.New("Not initialized")
	}

	b := c.browserObject

	// Return to main status screen
	b.Open(c.BaseURL + aegisMainURL)

	b.Dom().Find("div.ctl00_content_uxCallGrid div.Body a").Each(func(_ int, s *goquery.Selection) {
		x, exists := s.Attr("href")
		if exists {
			calls = append(calls, x)
		}
	})

	return calls, nil
}

func (c *AegisMonitor) GetStatus(url string) (CallStatus, error) {
	b := c.browserObject

	ret := CallStatus{}
	ret.Units = map[string]UnitStatus{}
	ret.Narratives = make([]Narrative, 0)

	err := b.Open(url)
	if err != nil {
		return ret, err
	}

	b.Dom().Find("table#ctl00_content_uxCallDetail tr td table tr td").Each(func(_ int, s *goquery.Selection) {
		var content string
		s.Find("span").Each(func(_ int, inner *goquery.Selection) {
			content = inner.Text()
		})
		s.Find("b").Each(func(_ int, inner *goquery.Selection) {
			switch inner.Text() {
			case "Priority: ":
				ret.Priority, _ = strconv.Atoi(content)
				break
			case "Call Type: ":
				ret.CallType = content
				break
			case "Nature of Call: ":
				ret.NatureOfCall = content
				break
			case "Location: ":
				ret.Location = content
				break
			case "Cross Streets: ":
				ret.CrossStreets = content
				break
			default:
			}
		})
	})

	b.Dom().Find("div#ctl00_content_uxNarrativesGrid div.Body table tbody tr").Each(func(_ int, s *goquery.Selection) {
		var nRecordedTime time.Time
		nMessage := ""
		nUser := ""

		s.Find("td").Each(func(_ int, inner *goquery.Selection) {
			cl, _ := inner.Attr("class")
			content := inner.Find("a").Text()
			switch cl {
			case "Key_DateTime DateTime":
				nRecordedTime = dateTime(content)
				break
			case "Key_Narrative":
				nMessage = content
				break
			case "Key_UserName":
				nUser = content
				break
			default:
			}
		})

		ret.Narratives = append(ret.Narratives, Narrative{
			RecordedTime: nRecordedTime,
			Message:      nMessage,
			User:         nUser,
		})
	})

	b.Dom().Find("div#ctl00_content_uxUnitsGrid div.Body table tbody tr").Each(func(_ int, s *goquery.Selection) {
		//fmt.Println("Found unit row")

		unit := ""
		status := ""
		dispatchTime := ""
		enrouteTime := ""
		arrivedTime := ""
		clearedTime := ""

		s.Find("td").Each(func(_ int, inner *goquery.Selection) {
			cl, _ := inner.Attr("class")
			content := inner.Find("a").Text()
			//fmt.Printf("--> %s : %s\n", cl, content)
			switch cl {
			case "Key_UnitNumber":
				unit = content
				break
			case "Key_Status":
				status = content
				break
			case "Key_DispatchTime DateTime":
				dispatchTime = content
				break
			case "Key_EnRouteTime DateTime":
				enrouteTime = content
				break
			case "Key_ArrivedTime DateTime":
				arrivedTime = content
				break
			case "Key_ClearedTime DateTime":
				clearedTime = content
				break
			default:
			}
		})

		if c.Suffix != "" && strings.HasSuffix(unit, c.Suffix) {
			return
		}

		ret.Units[unit] = UnitStatus{
			Unit:         unit,
			Status:       status,
			DispatchTime: dispatchTime,
			EnRouteTime:  enrouteTime,
			ArrivedTime:  arrivedTime,
			ClearedTime:  clearedTime,
		}
	})

	// div#ctl00_content_uxUnitsGrid div.Body table tbody tr
	// td.Key_UnitNumber a == Unit Number (QVMEDIC)
	// td.Key_Status a == Unit Status (DISPATCHED)
	// td.DispatchTime a / td.Key_EnRouteTime a / td.Key_ArrivedTime a

	return ret, nil
}

// GetClearedCalls fetches all cleared calls for specified date in format
// MM/DD/YYYY.
func (c *AegisMonitor) GetClearedCalls(dt string) (map[string]string, error) {
	calls := make(map[string]string, 0)

	if !c.initialized {
		return calls, errors.New("Not initialized")
	}

	b := c.browserObject

	// Return to main status screen
	b.Open(c.BaseURL + aegisMainURL)

	// Open cleared call search page
	//b.Click("a#ctl00_uxSearch")
	//b.Open(c.BaseURL + a)
	b.Open(c.BaseURL + aegisClearedCallURL)
	//log.Printf("CLEARED CALL SEARCH BODY: %#v", b.Body())

	if len(b.Forms()) < 1 {
		return calls, errors.New("Form does not exist")
	}
	f, err := b.Form("form#aspnetForm")
	if err != nil {
		return calls, err
	}
	f.Input("ctl00$content$uxORI", "04042    ")
	f.Input("ctl00$content$uxFromDate", dt)
	f.Input("ctl00$content$uxThruDate", dt)
	f.Input("ctl00$content$uxFromTime", "")
	f.Input("ctl00$content$uxThruTime", "")
	f.Input("ctl00$content$uxUnitNum", "")
	f.Input("ctl00$content$uxStreet", "")
	f.Input("ctl00$content$uxSearch", "Search")
	f.Input("ctl00$content$uxClearedCallsGrid$ctl19", "(All)")
	f.Input("ctl00$content$uxClearedCallsGrid$ctl31", "(All)")
	f.Input("ctl00$content$uxClearedCallsGrid$ctl43", "(All)")
	f.Input("ctl00$content$uxClearedCallsGrid$ctl49", "(All)")
	f.Input("ctl00$content$uxClearedCallsGrid$ctl55", "(All)")
	f.Input("ctl00$content$uxClearedCallsGrid$ColumnOrder", "CFSNumber|100px,CallType|171px,Address|481px,CallORI|100px,CallTime|134px,DispatchTime|131px,ArriveTime|130px,IncidentNumber|100px")
	f.Input("ctl00$content$uxClearedCallsGrid$uxDropDownPages", "1")
	f.Input("ctl00$content$uxClearedCallsGrid$uxGridSettingsId", "ASP.clearedcallsearch_aspx|uxClearedCallsGrid")

	//ccsf, _ := f.Dom().Html()
	//log.Printf("CLEARED CALL SEARCH FORM: %s", ccsf)

	if f.Click("ctl00$content$uxSearch") != nil {
		return calls, errors.Wrap(err, "Unable to click")
	}

	//log.Printf("CLEARED CALL BODY: %#v", b.Body())
	b.Dom().Find("div#ctl00_content_uxClearedCallsGrid div.Body table tbody tr").Each(func(_ int, s *goquery.Selection) {
		//h, _ := s.Html()
		//log.Printf("OUTER : %#v", h)
		var url string
		var id string
		s.Find("td.Key_CFSNumber a").Each(func(_ int, s2 *goquery.Selection) {
			//h, _ := s2.Html()
			//log.Printf("INNER1: %#v", h)

			x, exists := s2.Attr("href")
			if exists {
				url = URLPREFIX + x
			}
		})
		s.Find("td.Key_IncidentNumber a").Each(func(_ int, s2 *goquery.Selection) {
			//h, _ := s2.Html()
			//log.Printf("INNER2: %#v", h)

			id = s2.Text()
		})
		if id != "" && url != "" {
			calls[id] = url
		}
	})

	return calls, nil
}
