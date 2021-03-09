package monitor

import (
	"bytes"
	"fmt"
	"log"
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
	aegisLoggedOut      = "<span><H1>Server Error"
)

func init() {
	RegisterCadMonitor("aegis", func() CadMonitor { return &AegisMonitor{} })
}

// AegisMonitor represents a CadMonitor for the AEGIS CAD system
type AegisMonitor struct {
	// Suffix is the required suffix for units. Leaving blank disables
	// qualification by this value.
	Suffix string

	// LoginURL (example: "http://cadview.qvec.org/NewWorld.CAD.ViewOnly/"
	// or "http://cadview.qvec.org/")
	BaseURL string

	// FDID represents the FDID string which will be right padded to select
	// only local FDID events.
	FDID string

	// Numeric protocol; 0 defaults to latest
	Protocol int64

	terminateMonitor bool
	browserObject    *browser.Browser
	initialized      bool
	debug            bool
	cacheUser        string
	cachePass        string
}

// SetDebug enables or disables debug
func (c *AegisMonitor) SetDebug(d bool) {
	c.debug = d
}

func (c *AegisMonitor) SetTerminateMonitor(t bool) {
	c.terminateMonitor = t
}

func (c AegisMonitor) TerminateMonitor() bool {
	return c.terminateMonitor
}

func (c *AegisMonitor) ConfigureFromValues(values map[string]string) error {
	var ok bool
	if c.Suffix, ok = values["suffix"]; !ok {
		//return errors.New("'suffix' not defined")
	}
	if c.BaseURL, ok = values["baseUrl"]; !ok {
		return errors.New("'baseUrl' not defined")
	}
	if c.FDID, ok = values["fdid"]; !ok {
		return errors.New("'fdid' not defined")
	}
	if x, ok := values["protocol"]; ok {
		c.Protocol, _ = strconv.ParseInt(x, 10, 64)
	}
	return nil
}

func (c *AegisMonitor) LoggedIn() bool {
	return strings.Index(c.browserObject.Body(), aegisLoggedOut) == -1
}

func (c *AegisMonitor) Login(user, pass string) error {
	b := surf.NewBrowser()
	c.browserObject = b

	// Just so we can reconnect if necessary
	c.cacheUser = user
	c.cachePass = pass

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

func (c *AegisMonitor) KeepAlive() error {
	if !c.initialized {
		return errors.New("Not initialized")
	}

	b := c.browserObject

	// Test main screen
	b.Open(c.BaseURL + aegisMainURL)

	if !c.LoggedIn() {
		err := c.Login(c.cacheUser, c.cachePass)
		return err
	}

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

	// Determine if we're logged out
	if !c.LoggedIn() {
		return calls, ErrCadMonitorLoggedOut
	}

	//log.Printf("%s", b.Body())

	b.Dom().Find("div#ctl00_content_uxCallGrid div.Body table tbody tr td.Key_CFSNumber a").Each(func(_ int, s *goquery.Selection) {
		x, exists := s.Attr("href")
		if exists {
			calls = append(calls, x)
		}
	})

	return calls, nil
}

func (c *AegisMonitor) GetActiveAndUnassignedCalls() (map[string]CallStatus, error) {
	calls := map[string]CallStatus{}

	if !c.initialized {
		return calls, errors.New("Not initialized")
	}

	b := c.browserObject

	// Return to main status screen
	b.Open(c.BaseURL + aegisMainURL)

	// Determine if we're logged out
	if !c.LoggedIn() {
		return calls, ErrCadMonitorLoggedOut
	}

	// Retain raw body

	//log.Printf("%s", b.Body())

	b.Dom().Find("div#ctl00_content_uxCallGrid div.Body table tbody tr").Each(func(_ int, s *goquery.Selection) {
		cs := CallStatus{}
		s.Find("td.Key_Quadrant").Each(func(_ int, s2 *goquery.Selection) {
			cs.District = s2.Text()
		})
		s.Find("td.Key_Address").Each(func(_ int, s2 *goquery.Selection) {
			cs.Location = s2.Text()
		})
		s.Find("td.Key_CallTime").Each(func(_ int, s2 *goquery.Selection) {
			if s2.Text() == "" {
				return
			}
			cs.CallTime, _ = time.Parse("01/02/2006 15:04:05", s2.Text())
		})
		s.Find("td.Key_Priority").Each(func(_ int, s2 *goquery.Selection) {
			if s2.Text() == "" {
				return
			}
			cs.Priority, _ = strconv.Atoi(s2.Text())
		})
		s.Find("td.Key_CFSNumber a").Each(func(_ int, s2 *goquery.Selection) {
			x, exists := s2.Attr("href")
			if exists {
				cs.ID = x
				cs.LastUpdated = time.Now()
				cs.RawHTML = b.Body()
				calls[x] = cs
			}
		})
	})

	return calls, nil
}

func (c *AegisMonitor) GetStatusFromURL(url string) (CallStatus, error) {
	b := c.browserObject
	err := b.Open(c.BaseURL + url)
	if err != nil {
		return CallStatus{}, err
	}
	body := b.Body()
	if strings.Index(body, "<h2> <i>Runtime Error</i> </h2></span>") != -1 {
		return CallStatus{}, errors.New("Error fetching call page")
	}

	// Determine if we're logged out
	if !c.LoggedIn() {
		return CallStatus{}, ErrCadMonitorLoggedOut
	}

	return c.GetStatus([]byte(body), url)
}

func (c *AegisMonitor) GetStatus(content []byte, id string) (CallStatus, error) {
	// Start this an adequate amount of time back
	latestTime := time.Now().Add(-10 * time.Hour)

	ret := CallStatus{}
	ret.Units = make([]UnitStatus, 0)
	ret.UnitStatusMap = map[string]UnitStatus{}
	ret.Narratives = make([]Narrative, 0)

	// Maintain the unique identifier
	ret.ID = id

	// Retain raw body
	ret.RawHTML = string(content)

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(content))
	if err != nil {
		return ret, err
	}

	var tableSelector string
	switch c.Protocol {
	case 1:
		tableSelector = "table#ctl00_content_uxCallDetail tr td table tr td"
		break
	case 2:
	default:
		tableSelector = "div#ctl00_content_uxDetailWrapper table#ctl00_content_uxCallDetail tr td div table.summary tr td"
	}
	doc.Find(tableSelector).Each(func(_ int, s *goquery.Selection) {
		var content string
		s.Find("span").Each(func(_ int, inner *goquery.Selection) {
			if inner.Text() == "" {
				// Skip false empties
				return
			}
			content = inner.Text()
		})
		s.Find("b").Each(func(_ int, inner *goquery.Selection) {
			switch inner.Text() {
			case "Call Date/Time: ":
				ret.CallTime, _ = time.Parse("01/02/2006 15:04:05", content) // Mon Jan 2 15:04:05 -0700 MST 2006
				if ret.CallTime.After(latestTime) {
					latestTime = ret.CallTime
				}
				break
			case "Dispatch Date/Time: ":
				ret.DispatchTime, _ = time.Parse("01/02/2006 15:04:05", content) // Mon Jan 2 15:04:05 -0700 MST 2006
				if ret.DispatchTime.After(latestTime) {
					latestTime = ret.DispatchTime
				}
				break
			case "Arrival Date/Time: ":
				ret.ArrivalTime, _ = time.Parse("01/02/2006 15:04:05", content) // Mon Jan 2 15:04:05 -0700 MST 2006
				if ret.ArrivalTime.After(latestTime) {
					latestTime = ret.ArrivalTime
				}
				break
			case "Caller Phone: ":
				ret.CallerPhone = content
				break
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

	// Determine call ID from incidents grid
	i := []Incident{}
	doc.Find("div#ctl00_content_uxIncidentsGrid div.Body table tbody tr").Each(func(_ int, s *goquery.Selection) {
		thisi := Incident{}
		s.Find("td").Each(func(_ int, inner *goquery.Selection) {
			cl, _ := inner.Attr("class")
			content := inner.Find("a").Text()
			switch cl {
			case "Key_ORI":
				thisi.FDID = content
				break
			case "Key_IncidentNumber":
				thisi.IncidentNumber = content
				break
			default:
			}
		})
		i = append(i, thisi)

		if thisi.FDID == c.FDID {
			ret.CallID = thisi.IncidentNumber
		}
	})

	doc.Find("div#ctl00_content_uxNarrativesGrid div.Body table tbody tr").Each(func(_ int, s *goquery.Selection) {
		var nRecordedTime time.Time
		nMessage := ""
		nUser := ""

		s.Find("td").Each(func(_ int, inner *goquery.Selection) {
			cl, _ := inner.Attr("class")
			content := inner.Find("a").Text()
			switch cl {
			case "Key_DateTime DateTime":
				nRecordedTime = dateTime(content)
				if nRecordedTime.After(latestTime) {
					latestTime = nRecordedTime
				}
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

	doc.Find("div#ctl00_content_uxUnitsGrid div.Body table tbody tr").Each(func(_ int, s *goquery.Selection) {
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
				if dateTime(dispatchTime).After(latestTime) {
					latestTime = dateTime(dispatchTime)
				}
				break
			case "Key_EnRouteTime DateTime":
				enrouteTime = content
				if dateTime(enrouteTime).After(latestTime) {
					latestTime = dateTime(enrouteTime)
				}
				break
			case "Key_ArrivedTime DateTime":
				arrivedTime = content
				if dateTime(arrivedTime).After(latestTime) {
					latestTime = dateTime(arrivedTime)
				}
				break
			case "Key_ClearedTime DateTime":
				clearedTime = content
				if dateTime(clearedTime).After(latestTime) {
					latestTime = dateTime(clearedTime)
				}
				break
			default:
			}
		})

		if c.Suffix != "" && !strings.HasSuffix(unit, c.Suffix) {
			return
		}

		st := UnitStatus{
			Unit:         unit,
			Status:       status,
			DispatchTime: dispatchTime,
			EnRouteTime:  enrouteTime,
			ArrivedTime:  arrivedTime,
			ClearedTime:  clearedTime,
		}
		ret.UnitStatusMap[unit] = st
		ret.Units = append(ret.Units, st)
	})

	// div#ctl00_content_uxUnitsGrid div.Body table tbody tr
	// td.Key_UnitNumber a == Unit Number (QVMEDIC)
	// td.Key_Status a == Unit Status (DISPATCHED)
	// td.DispatchTime a / td.Key_EnRouteTime a / td.Key_ArrivedTime a

	ret.Incidents = i

	// Only return the most recent time involved
	ret.LastUpdated = latestTime

	return ret, nil
}

// GetClearedCalls fetches all cleared calls for specified date in format
// MM/DD/YYYY.
func (c *AegisMonitor) GetClearedCalls(dt string) (map[string]string, error) {
	calls := make(map[string]string, 0)

	log.Printf("GetClearedCalls")

	if !c.initialized {
		return calls, errors.New("Not initialized")
	}

	b := c.browserObject

	// Return to main status screen
	b.Open(c.BaseURL + aegisMainURL)

	// Determine if we're logged out
	if !c.LoggedIn() {
		return calls, ErrCadMonitorLoggedOut
	}

	// Open cleared call search page
	//b.Click("a#ctl00_uxSearch")
	//b.Open(c.BaseURL + a)
	b.Open(c.BaseURL + aegisClearedCallURL)
	if c.debug {
		log.Printf("CLEARED CALL SEARCH BODY: %#v", b.Body())
	}

	if len(b.Forms()) < 1 {
		return calls, errors.New("Form does not exist")
	}
	f, err := b.Form("form#aspnetForm")
	if err != nil {
		return calls, err
	}
	switch c.Protocol {
	case 1:
		// Original protocol
		f.Input("ctl00$content$uxORI", fmt.Sprintf("%-9v", c.FDID))
		break
	case 2:
	default:
		// From protocol 2 on, no space padding
		f.Input("ctl00$content$uxORI", c.FDID)
	}
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

	if c.debug {
		log.Printf("CLEARED CALL BODY: %#v", b.Body())
	}
	b.Dom().Find("div#ctl00_content_uxClearedCallsGrid div.Body table tbody tr").Each(func(_ int, s *goquery.Selection) {
		//h, _ := s.Html()
		//log.Printf("OUTER : %#v", h)
		var url string
		var id string
		s.Find("td.Key_CFSNumber a").Each(func(_ int, s2 *goquery.Selection) {
			if url != "" {
				// Just accept the first one
				return
			}

			//h, _ := s2.Html()
			//log.Printf("INNER1: %#v", h)

			x, exists := s2.Attr("href")
			if exists {
				url = x
				//log.Printf("url = %s, x = %s, BaseURL = %s", url, x, c.BaseURL)
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

// Monitor runs a monitoring function with a callback function
func (c *AegisMonitor) Monitor(callback func(CallStatus) error, pollInterval int) error {
	currentCallMap := map[string]CallStatus{}
	for {
		// Poll for data
		calls, err := c.GetActiveAndUnassignedCalls()
		if err != nil {
			log.Printf("Monitor: %s", err.Error())
			goto continueOn
		}

		for _, call := range calls {
			if _, ok := currentCallMap[call.ID]; !ok {
				log.Printf("Found new call URL %s", call.ID)
				status, err := c.GetStatusFromURL(call.ID)
				if err != nil {
					log.Printf("Monitor: %s", err.Error())
					continue
				}
				// Record in current mapping so we don't do bad things
				currentCallMap[call.ID] = status
				// If there's a callback, send the data back
				if callback != nil {
					go func(status CallStatus) {
						err := callback(status)
						if err != nil {
							log.Printf("Monitor: Callback: %s", err.Error())
						}
					}(status)
				}
			} else {
				// Re-poll for call data
				log.Printf("Updating call URL %s", call.ID)
				status, err := c.GetStatusFromURL(call.ID)
				if err != nil {
					log.Printf("Monitor: %s", err.Error())
					continue
				}
				// Record in current mapping so we don't do bad things
				currentCallMap[call.ID] = status
				// If there's a callback, send the data back
				if callback != nil {
					go func(status CallStatus) {
						err := callback(status)
						if err != nil {
							log.Printf("Monitor: Callback: %s", err.Error())
						}
					}(status)
				}
			}
		}

		if c.TerminateMonitor() {
			log.Printf("Terminating monitor")
			break
		}

		// Sleep for pollInterval duration
	continueOn:
		for iter := 0; iter < pollInterval; iter++ {
			time.Sleep(time.Second)
			if c.TerminateMonitor() {
				log.Printf("Terminating monitor")
				break
			}
		}
	}
	return nil
}
