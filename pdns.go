package utils

import (
	"strconv"
)

type PDNS struct {
	APIUrl string
	APIKey string
}

const (
	PATHAPI   = "/api/v1/servers/localhost"
	PATHZones = PATHAPI + "/zones"
)

func (pdns *PDNS) Add(zone, hostname, ip string, ttl uint) (err error) {

	_, err = pdns.pdnsApi(PATHZones+"/"+zone, "PATCH", `
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

	_, err = pdns.pdnsApi(PATHZones+"/"+zone, "PATCH", `{"rrsets": [
			 {"name": "`+hostname+"."+zone+`.",
				"type": "A",
				"changetype": "DELETE"
			 }
		 	]
		 }`)
	return

}
func (pdns *PDNS) GetAllRR(zone string) (rsp []byte, err error) {

	rsp, err = pdns.pdnsApi(PATHZones+"/"+zone, "GET", "")
	return

}
func (pdns *PDNS) GetAllZones() (rsp []byte, err error) {

	rsp, err = pdns.pdnsApi(PATHZones, "GET", "")
	return

}
func (pdns *PDNS) pdnsApi(path, method, data string) (rsp []byte, err error) {

	h := make(map[string]string)
	h["X-API-Key"] = pdns.APIKey

	if _, rsp, err = HttpDo(pdns.APIUrl+path, method, nil, nil, h, nil, []byte(data)); err != nil {
		return
	}

	return
}
