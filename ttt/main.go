package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// 获取原始数据
	initUrl := "http://query.sse.com.cn/infodisplay/querySpecialTipsInfoByPage.do?jsonCallBack=jsonpCallback34822&isPagination=true&searchDate=2020-12-02&bgFlag=1&searchDo=1&pageHelp.pageSize=25&pageHelp.pageNo=1&pageHelp.beginPage=1&pageHelp.cacheSize=1&pageHelp.endPage=5&_=1606875884329"
	header := make(map[string]string)
	header["Accept"] = "*/*"
	header["Accept-Encoding"] = "gzip, deflate"
	header["Accept-Language"] = "zh-CN,zh;q=0.9,en;q=0.8"
	header["Connection"] = "keep-alive"
	header["Cookie"] = "yfx_c_g_u_id_10000042=_ck20120210244318750116153672215; yfx_f_l_v_t_10000042=f_t_1606875883868__r_t_1606875883868__v_t_1606875883868__r_c_0; VISITED_MENU=%5B%228314%22%5D; JSESSIONID=3F5E79C5542A72B676476053F23B171E"
	header["Host"] = "query.sse.com.cn"
	header["Referer"] = "http://www.sse.com.cn/"
	header["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.67 Safari/537.36"

	client := &http.Client{}
	req, err := http.NewRequest("GET", initUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	for key, v := range header {
		req.Header.Add(key, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s", body)
}
