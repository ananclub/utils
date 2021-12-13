package utils

import (
	"strconv"
)

type PDNS struct {
	APIUrl string
	APIKey string
}

func (pdns *PDNS) Add(hostname string, ip string, ttl uint) (err error) {

	_, err = pdns.pdnsApi("PATCH", `
		{"rrsets": [
		 	{"name": "`+hostname+`.cluster.melot.cn.",
			 "type": "A","ttl": `+strconv.Itoa(int(ttl))+`,
			 "changetype": "REPLACE",
			 "records": [ 
				{"content": "`+ip+`", "disabled": false }
				]
		 	}
		 ]
		}`)
	return

}
func (pdns *PDNS) Del(hostname string) (err error) {

	_, err = pdns.pdnsApi("PATCH", `{"rrsets": [
			 {"name": "`+hostname+`.cluster.melot.cn.",
				"type": "A",
				"changetype": "DELETE"
			 }
		 	]
		 }`)
	return

}
func (pdns *PDNS) GetAll() (rsp []byte, err error) {

	rsp, err = pdns.pdnsApi("GET", "")
	return

}
func (pdns *PDNS) pdnsApi(method string, data string) (rsp []byte, err error) {

	h := make(map[string]string)
	h["X-API-Key"] = pdns.APIKey

	if _, rsp, err = HttpDo(pdns.APIUrl, method, nil, nil, h, nil, []byte(data)); err != nil {
		return
	}

	return
}
