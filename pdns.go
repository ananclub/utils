package utils

import (
	"strconv"
)

type PDNS struct {
	APIUrl string
	APIKey string
}

func (pdns *PDNS) Add(zone, hostname, ip string, ttl uint) (err error) {

	_, err = pdns.pdnsApi("PATCH", zone, `
		{"rrsets": [
		 	{"name": "`+hostname+"."+zone+`.",
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
func (pdns *PDNS) Del(zone, hostname string) (err error) {

	_, err = pdns.pdnsApi("PATCH", zone, `{"rrsets": [
			 {"name": "`+hostname+"."+zone+`.",
				"type": "A",
				"changetype": "DELETE"
			 }
		 	]
		 }`)
	return

}
func (pdns *PDNS) GetAllRR(zone string) (rsp []byte, err error) {

	rsp, err = pdns.pdnsApi("GET", zone, "")
	return

}
func (pdns *PDNS) pdnsApi(method, zone, data string) (rsp []byte, err error) {

	h := make(map[string]string)
	h["X-API-Key"] = pdns.APIKey

	if _, rsp, err = HttpDo(pdns.APIUrl+"/"+zone, method, nil, nil, h, nil, []byte(data)); err != nil {
		return
	}

	return
}
