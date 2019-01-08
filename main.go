package main

import (
	"fmt"
	"net/url"

	"github.com/unix2dos/zj-business/pkg/util"
)

func main() {

	_url := "http://wsdj.zjaic.gov.cn/pda.do?method=mccxList&type=0"

	values := url.Values{}
	values.Set("zwzh", util.UTF82GBK("你"))
	//values.Set("currentPage", "3")
	//values.Set("ec_p", "2")
	body, err := util.Post(_url, values.Encode())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(util.GBK2UTF8(string(body)))
}
