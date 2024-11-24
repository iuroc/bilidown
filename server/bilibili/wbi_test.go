package bilibili

import (
	"bilidown/util"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestGetWbiKey(t *testing.T) {
	db := util.MustGetDB("../data.db")
	defer db.Close()
	sessdata, err := GetSessdata(db)
	if err != nil {
		t.Error(err)
	}
	client := BiliClient{SESSDATA: sessdata}
	mixinKey, err := client.GetMixinKey(db)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(mixinKey)
}

func TestTime(t *testing.T) {
	fmt.Println(time.Now().Unix())
}

func TestURLEncode(t *testing.T) {
	str := url.Values{
		"foo": {"one one four"},
		"bar": {"五一四"},
		"baz": {"1919810"},
	}.Encode()
	str = strings.ReplaceAll(str, "+", "%20")
	fmt.Println(str)
}

func TestWbiSign(t *testing.T) {
	db := util.MustGetDB("../data.db")
	defer db.Close()
	sessdata, err := GetSessdata(db)
	if err != nil {
		t.Error(err)
	}
	client := BiliClient{SESSDATA: sessdata}
	mixinKey, err := client.GetMixinKey(db)
	if err != nil {
		t.Error(err)
	}
	newParams := WbiSign(map[string]string{
		"dm_cover_img_str": "QU5HTEUgKEludGVsLCBJbnRlbChSKSBJcmlzKFIpIFhlIEdyYXBoaWNzICgweDAwMDA5QTQ5KSBEaXJlY3QzRDExIHZzXzVfMCBwc181XzAsIEQzRDExKUdvb2dsZSBJbmMuIChJbnRlbC",
		"dm_img_inter":     `{"ds":[{"t":10,"c":"YmUtcGFnZXItaXRlbSBiZS1wYWdlci1pdGVtLWFjdGl2ZQ","p":[330,110,110],"s":[179,179,358]}],"wh":[5762,6704,10],"of":[318,636,318]}`,
		"dm_img_list":      `[{"x":1847,"y":-1616,"z":0,"timestamp":2486463,"k":70,"type":0},{"x":1972,"y":-1005,"z":24,"timestamp":2486563,"k":120,"type":0},{"x":2564,"y":-638,"z":72,"timestamp":2486665,"k":107,"type":0},{"x":4200,"y":-753,"z":199,"timestamp":2486765,"k":91,"type":0},{"x":4922,"y":-73,"z":426,"timestamp":2486865,"k":106,"type":0},{"x":4656,"y":-512,"z":35,"timestamp":2486965,"k":94,"type":0},{"x":4942,"y":-590,"z":79,"timestamp":2487066,"k":126,"type":0},{"x":5200,"y":-414,"z":123,"timestamp":2487167,"k":60,"type":0},{"x":4916,"y":-820,"z":228,"timestamp":2563922,"k":111,"type":0},{"x":5004,"y":-50,"z":570,"timestamp":2564022,"k":110,"type":0},{"x":4784,"y":330,"z":666,"timestamp":2564122,"k":98,"type":0},{"x":5029,"y":767,"z":956,"timestamp":2564223,"k":115,"type":0},{"x":5327,"y":1400,"z":1307,"timestamp":2564324,"k":107,"type":0},{"x":4598,"y":671,"z":578,"timestamp":2564424,"k":125,"type":0},{"x":5228,"y":1283,"z":1193,"timestamp":2564525,"k":83,"type":0},{"x":5210,"y":932,"z":978,"timestamp":2564625,"k":66,"type":0},{"x":5034,"y":195,"z":484,"timestamp":2564725,"k":112,"type":0},{"x":6496,"y":1467,"z":1849,"timestamp":2564825,"k":83,"type":0},{"x":6592,"y":1387,"z":1852,"timestamp":2564929,"k":122,"type":0},{"x":6001,"y":728,"z":1235,"timestamp":2565030,"k":95,"type":0},{"x":6666,"y":1393,"z":1900,"timestamp":2627309,"k":121,"type":1},{"x":5486,"y":5603,"z":2283,"timestamp":2627842,"k":102,"type":0},{"x":5269,"y":4186,"z":1733,"timestamp":2627947,"k":117,"type":0},{"x":4439,"y":3174,"z":897,"timestamp":2628056,"k":117,"type":0},{"x":5231,"y":3959,"z":1687,"timestamp":2628268,"k":79,"type":0},{"x":6197,"y":4786,"z":2587,"timestamp":2628368,"k":96,"type":0},{"x":4667,"y":3298,"z":724,"timestamp":2628470,"k":91,"type":0},{"x":4658,"y":3436,"z":205,"timestamp":2629433,"k":88,"type":0},{"x":5380,"y":3443,"z":1117,"timestamp":2630390,"k":121,"type":0},{"x":5559,"y":2678,"z":1920,"timestamp":2630492,"k":116,"type":0},{"x":5517,"y":2465,"z":2023,"timestamp":2630634,"k":86,"type":0},{"x":6584,"y":3533,"z":3087,"timestamp":2630782,"k":96,"type":0},{"x":6183,"y":3242,"z":2839,"timestamp":2630892,"k":103,"type":0},{"x":6277,"y":580,"z":1449,"timestamp":2663911,"k":68,"type":0},{"x":7556,"y":3077,"z":3053,"timestamp":2664013,"k":114,"type":0},{"x":5321,"y":1642,"z":1201,"timestamp":2664114,"k":104,"type":0},{"x":5839,"y":2261,"z":1761,"timestamp":2664215,"k":72,"type":0},{"x":6303,"y":2726,"z":2222,"timestamp":2664317,"k":93,"type":0},{"x":7553,"y":4255,"z":3532,"timestamp":2664420,"k":104,"type":0},{"x":7677,"y":4381,"z":3650,"timestamp":2664521,"k":119,"type":0},{"x":8184,"y":4762,"z":4098,"timestamp":2664622,"k":122,"type":0},{"x":4555,"y":740,"z":406,"timestamp":2664722,"k":118,"type":0},{"x":4305,"y":241,"z":75,"timestamp":2664822,"k":105,"type":0},{"x":6373,"y":2218,"z":2117,"timestamp":2664934,"k":110,"type":0},{"x":6079,"y":1757,"z":1680,"timestamp":2665036,"k":90,"type":0},{"x":8158,"y":3462,"z":3524,"timestamp":2665136,"k":100,"type":0},{"x":9312,"y":4205,"z":4531,"timestamp":2665236,"k":124,"type":0},{"x":5306,"y":89,"z":487,"timestamp":2665337,"k":89,"type":0},{"x":9721,"y":4399,"z":4941,"timestamp":2665437,"k":100,"type":0},{"x":9058,"y":3727,"z":4282,"timestamp":2665540,"k":107,"type":0}]`,
		"dm_img_str":       "V2ViR0wgMS4wIChPcGVuR0wgRVMgMi4wIENocm9taXVtKQ",
		"keyword":          "",
		"mid":              "596866446",
		"order":            "pubdate",
		"order_avoided":    "true",
		"platform":         "web",
		"pn":               "3",
		"ps":               "30",
		"tid":              "0",
		"w_webid":          "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzcG1faWQiOiIwLjAiLCJidXZpZCI6IjBGRTlGRDMxLTZGQTQtRjU2RC1GQjEwLUI4OTg2ODUyRUZDNzUwMDUyaW5mb2MiLCJ1c2VyX2FnZW50IjoiTW96aWxsYS81LjAgKFdpbmRvd3MgTlQgMTAuMDsgV2luNjQ7IHg2NCkgQXBwbGVXZWJLaXQvNTM3LjM2IChLSFRNTCwgbGlrZSBHZWNrbykgQ2hyb21lLzEzMS4wLjAuMCBTYWZhcmkvNTM3LjM2IEVkZy8xMzEuMC4wLjAiLCJidXZpZF9mcCI6IjY3YzZkNzg2NTgyYzQ2OGYwNjRiZTE5MTk2OWU4ZWE4IiwiYmlsaV90aWNrZXQiOiI4ODA0MDgzYmYzZGQ1NGNkNzI4OWZjZDE0MmIwN2U4ZCIsImNyZWF0ZWRfYXQiOjE3MzI0NzE3MTIsInR0bCI6ODY0MDAsInVybCI6Ii81OTY4NjY0NDYvdmlkZW8_dGlkPTBcdTAwMjZwbj0yXHUwMDI2a2V5d29yZD1cdTAwMjZvcmRlcj1wdWJkYXRlIiwicmVzdWx0Ijoibm9ybWFsIiwiaXNzIjoiZ2FpYSIsImlhdCI6MTczMjQ3MTcxMn0.pTVyhdRcB2VFcXqeYoi3TjnlEFLAXMpNNJK-dW9Gq3vqBL4YldGgZBuGBXXF7Ldwg_vW6TQg6pQCQ7vz357ws2Z9g2-kLIuNmR3j8oMg2zAXAND1q5oJNw0jNnevhLlB8_vOcip0eIJSHRjqbPbNShKLOcSnSfLaiI64EHwjRFAOEPYVw3evLXKB4TFnxSRi1WDSI684TfNrXp0_2yTJvuPheHQmQC1NcUP_P9tqTuRiDy3YfkuR8PlRcQxmKHVs-byObL5WEPqMMQa8b8zidRtzEkbGV7ra8gTgv1HwgpSDi_Y1VNsX2WtNcBPTSwURWc7zREf-RJqruzgcf1ErpA",
		"web_location":     "1550101",
	}, mixinKey)

	fmt.Println(newParams.Encode())
}
