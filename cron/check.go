package cron

import (
	"log"
	"sync"
	"time"

	"github.com/urlooker/web/api"
	webg "github.com/urlooker/web/g"

	"github.com/urlooker/agent/backend"
	"github.com/urlooker/agent/g"
	"github.com/urlooker/agent/utils"
)

var wg *sync.WaitGroup

func StartCheck() {
	t1 := time.NewTicker(time.Duration(g.Config.Web.Interval) * time.Second)
	for {
		items, _ := GetItem()
		for _, item := range items {
			g.WorkerChan <- 1
			go utils.CheckTargetStatus(item, wg)
			time.Sleep(500 * time.Millisecond)
		}
		<-t1.C
	}
}

func GetItem() ([]*webg.DetectedItem, error) {
	hostname, _ := g.Hostname()

	var resp api.GetItemResponse
	err := backend.CallRpc("Web.GetItem", hostname, &resp)
	if err != nil {
		log.Println(err)
	}
	if resp.Message != "" {
		log.Println(resp.Message)
	}

	return resp.Data, err
}
