package business

import (
	"fmt"
	"net/url"
	"strconv"
	"sync"

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
	values.Set("type", "0")
	values.Set("zwzh", util.UTF82GBK("食用菌"))

	var data []Context
	var word KeyWord
	var page = 1
	var sum = 0
	var first = true

	for {
		values.Set("currentPage", strconv.Itoa(page))
		str := values.Encode()

		body, err := util.Post(_url, str)
		if err != nil {
			log.Errorf("[spiderType] get url err:%v", err)
			goto Sleep
		}

		data, word, err = GetNewContent(body)
		if err != nil {
			log.Errorf("[spiderType] parse content err:%v", err)
			goto Sleep
		}

		if !first {
			for _, v := range data {
				fmt.Println("-------------------new", v)
			}
		}

		page++
		sum += len(data)
		fmt.Printf("spiderType post: %s, total: %d,  sum: %d\n", str, word.Total, sum)

		if page > word.MaxPage {
			first = false
			page = 1
		}

	Sleep:
		// time.Sleep(time.Second)
	}
}

// method=mccxList&type=0&zwzh=%B9%AB%CB%BE
// method=mccxList&type=0&zwzh=%B9%AB%CB%BE&currentPage=53234
