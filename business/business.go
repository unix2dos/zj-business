package business

import (
	"fmt"
	"net/url"
	"strconv"
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
	var word KeyWord
	var maxPage = 1

	for {
		values.Set("currentPage", strconv.Itoa(maxPage))
		str := values.Encode()
		fmt.Println("post :", str)

		body, err := util.Post(_url, str)
		if err != nil {
			log.Errorf("[Spider] get url err:%v", err)
			goto Sleep
		}

		data, word, err = GetNewContent(body)
		if err != nil {
			log.Errorf("[Spider] parse content err:%v", err)
			goto Sleep
		}

		maxPage = word.MaxPage
		fmt.Println("word: ", word)
		for _, v := range data {
			fmt.Println(v)
		}

	Sleep:
		time.Sleep(time.Minute * 10)
	}

}

// method=mccxList
// method=mccxList&currentPage=2712

// method=mccxList&type=0&zwzh=%B9%AB%CB%BE
// method=mccxList&type=0&zwzh=%B9%AB%CB%BE&currentPage=53234
