package core

import (
	"log"
	"strconv"
)

type LRSEvent struct {
	ElapsedTime uint   `json:"elapsedTime"`
	UUID        string `json:"uuid"`
	OrderType   string `json:"orderType"`
	TableName   string `json:"locationName"`
	State       string `json:"state"`
	PagerNumber string `json:"name"`
	Paged       bool   `json:"paged"`
}

func (app AppContext) HandleLRSEvent(siteID KountaID, lrsEvent LRSEvent) {
	pagerNumber, err := strconv.ParseInt(lrsEvent.PagerNumber, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}

	switch lrsEvent.State {
	case "started":
		_, err := app.CreateOrderForPager(siteID, pagerNumber)
		if err != nil {
			log.Println(err)
			return
		}

	case "located":
		err := app.LinkOrderWithTable(siteID, pagerNumber, lrsEvent.TableName)
		if err != nil {
			log.Println(err)
			return
		}

		//no longer handling "cleared" since LRS will automatically send us a cleared after a while and screw things up
	}
}
