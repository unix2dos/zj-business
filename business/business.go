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
	mux      sync.Mutex
	currPage int
	maxPage  int
	ok       bool
	sum      int
}

func New() *Business {
	return &Business{}
}

func (b *Business) Run() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		wg.Done()
		b.start()
	}()

	wg.Wait()
}

func (b *Business) start() {

	for i := 0; i < 10; i++ {
		go func() {
			b.spider()
		}()
	}
}

func (b *Business) spider() {

	_url := "http://wsdj.zjaic.gov.cn/pda.do"
	values := url.Values{}
	values.Set("method", "mccxList")
	values.Set("type", "0")
	values.Set("zwzh", util.UTF82GBK("合作社"))

	var data []Context
	var word KeyWord

	for {
		values.Set("currentPage", strconv.Itoa(b.GetPage()))
		str := values.Encode()
		body, err := util.Post(_url, str)
		if err != nil {
			log.Errorf("[spiderType] get url err:%v", err)
			continue
		}

		data, word, err = GetNewContent(body)
		if err != nil {
			log.Errorf("[spiderType] parse content err:%v", err)
			continue
		}

		b.SetCount(str, word.Total, word.MaxPage, len(data))

		if b.ok {
			for _, v := range data {
				fmt.Println("-------------------new", v)
			}
		}

	}
}

// method=mccxList&type=0&zwzh=%B9%AB%CB%BE
// method=mccxList&type=0&zwzh=%B9%AB%CB%BE&currentPage=53234

func (b *Business) GetPage() int {
	b.mux.Lock()
	defer b.mux.Unlock()

	if b.currPage > b.maxPage {
		b.currPage = 0
	}
	b.currPage++

	return b.currPage
}

func (b *Business) SetCount(str string, total int, maxpage int, count int) {
	b.mux.Lock()
	defer b.mux.Unlock()

	b.maxPage = maxpage
	b.sum += count
	fmt.Printf("spiderType post: %s, maxpage: %d, total: %d,  sum: %d\n", str, maxpage, total, b.sum)

}
