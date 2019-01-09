package business

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/unix2dos/zj-business/pkg/log"
	"github.com/unix2dos/zj-business/pkg/util"
)

type Business struct {
}

func New() *Business {
	return &Business{}
}

func (b *Business) Run() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		wg.Done()
		b.spider()
	}()

	wg.Wait()
}

func (b *Business) spider() {

	_url := "http://wsdj.zjaic.gov.cn/pda.do"
	values := url.Values{}
	values.Set("method", "mccxList")
	var data []Context

	for {
		body, err := util.Post(_url, values.Encode())
		if err != nil {
			log.Errorf("[Spider] get url err:%v", err)
			goto Sleep
		}

		data, err = GetNewContent(body)
		if err != nil {
			log.Errorf("[Spider] parse content err:%v", err)
			goto Sleep
		}

		for _, v := range data {
			fmt.Println(v)
		}

	Sleep:
		time.Sleep(time.Minute * 5)
	}

}

// _url := "http://wsdj.zjaic.gov.cn/pda.do?method=mccxList&type=0"
// values := url.Values{}
// values.Set("zwzh", util.UTF82GBK("公司"))
// values.Set("currentPage", "2")
//values.Set("ec_p", "2")
