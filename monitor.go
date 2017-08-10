package main

import (
	"errors"
	//"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/agent"
	"github.com/headzoo/surf/browser"
)

type CadBrowser struct {
	browserObject *browser.Browser
	initialized   bool
}

func (c *CadBrowser) Login(user, pass string) error {
	b := surf.NewBrowser()
	c.browserObject = b

	b.SetUserAgent(agent.Chrome())

	// Required to not have ASP.NET garbage yak all over me
	b.AddRequestHeader("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")

	err := b.Open("http://cadview.qvec.org/NewWorld.CAD.ViewOnly/Login.aspx")
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
		return err
	}

	c.initialized = true

	return nil
}

func (c *CadBrowser) GetActiveCalls() error {
	if !c.initialized {
		return errors.New("Not initialized")
	}

	// div.ctl00_content_uxCallGrid div.Body a

	return nil
}

func (c *CadBrowser) GetStatus(url string) (map[string]string, error) {
	b := c.browserObject

	ret := map[string]string{}

	err := b.Open(url)
	if err != nil {
		return ret, err
	}

	b.Dom().Find("div#ctl00_content_uxUnitsGrid div.Body table tbody tr").Each(func(_ int, s *goquery.Selection) {
		//fmt.Println("Found unit row")

		unit := ""
		status := ""

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
			default:
			}
		})

		ret[unit] = status
	})

	// div#ctl00_content_uxUnitsGrid div.Body table tbody tr
	// td.Key_UnitNumber a == Unit Number (QVMEDIC)
	// td.Key_Status a == Unit Status (DISPATCHED)
	// td.DispatchTime a / td.Key_EnRouteTime a / td.Key_ArrivedTime a

	return ret, nil
}
