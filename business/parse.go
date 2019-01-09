package business

import (
	"bytes"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/unix2dos/zj-business/pkg/util"
)

type Context struct {
	Index   int
	ID      string
	Company string //名字
	Type    string //企业
	State   string //状态
}

type KeyWord struct {
	Word    string
	Total   int
	MaxPage int
}

var (
	TOTAL_SEND = make(map[string]int, 100)
	Mux        sync.Mutex
)

func ParseContent(data []byte) (res []Context, word KeyWord, err error) {

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}

	maxPage := strings.TrimSpace(util.GBK2UTF8(doc.Find(".maxPage").Text()))
	maxPage = strings.TrimPrefix(maxPage, "共")
	maxPage = strings.TrimSuffix(maxPage, "页")
	word.Total, _ = strconv.Atoi(strings.TrimSpace(doc.Find(".wrap .cxjg span").Text()))
	word.MaxPage, _ = strconv.Atoi(maxPage)

	doc.Find(".search_result tr").Each(func(i int, s *goquery.Selection) {

		_index := s.Children().First().Text()
		_company := s.Find(".company").Text()
		_type := s.Find(".search_txt").Text()
		_state := s.Find(".state").Text()
		if _company == "" {
			return
		}

		context := Context{}
		context.Index, _ = strconv.Atoi(strings.TrimSpace(util.GBK2UTF8(_index)))
		context.Company = strings.TrimSpace(util.GBK2UTF8(_company))
		context.Type = strings.TrimSpace(util.GBK2UTF8(_type))
		context.State = strings.TrimSpace(util.GBK2UTF8(_state))
		context.ID = context.Company + "_" + context.State

		if context.State == "选择" {
			str, ok := s.Find(".state a").Attr("onclick")
			if ok {
				str = strings.TrimPrefix(str, "javascript:regi_dddl('")
				str = strings.TrimSuffix(str, "')")
				context.State = str
			}
		}
		res = append(res, context)
	})

	return
}

func GetNewContent(data []byte) (newS []Context, word KeyWord, err error) {
	Mux.Lock()
	defer func() { Mux.Unlock() }()

	res, word, err := ParseContent(data)
	if err != nil {
		return
	}
	for _, v := range res {
		if _, ok := TOTAL_SEND[v.ID]; ok {
			TOTAL_SEND[v.ID]++
		} else {
			TOTAL_SEND[v.ID] = 0
			newS = append(newS, v)
		}
	}
	return
}
