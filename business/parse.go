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
	ID    string
	Type  string
	State string
}

var (
	TOTAL_SEND = make(map[string]int, 100)
	Mux        sync.Mutex
)

func ParseContent(data []byte) (res []Context, total int, err error) {

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}

	total, _ = strconv.Atoi(doc.Find(".wrap .cxjg span").Text())

	doc.Find(".search_result tr").Each(func(i int, s *goquery.Selection) {

		_company := s.Find(".company").Text()
		_type := s.Find(".search_txt").Text()
		_state := s.Find(".state").Text()
		if _company == "" {
			return
		}

		context := Context{}
		context.ID = strings.TrimSpace(util.GBK2UTF8(_company))
		context.Type = strings.TrimSpace(util.GBK2UTF8(_type))
		context.State = strings.TrimSpace(util.GBK2UTF8(_state))

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

func GetNewContent(data []byte) (newS []Context, err error) {
	Mux.Lock()
	defer func() { Mux.Unlock() }()

	res, _, err := ParseContent(data)
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

	if len(newS) >= 20 {
		// newS = []Context{}
	}

	return
}
