package utils

import (
	"strconv"
)

type PDNS struct {
	APIUrl string
	APIKey string
}

const (
	PATHBASE            = "/api/v1"
	PATHServers         = PATHBASE + "/servers"
	PATHServerLocalhost = PATHServers + "/localhost"
	PATHZones           = "/zones"
	PATHLocalZones      = PATHServerLocalhost + PATHZones
)

func (pdns *PDNS) Add(hostname, zone, ip string, ttl uint) (err error) {

	_, err = pdns.pdnsApi(PATHLocalZones+"/"+zone, "PATCH", `
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
func (pdns *PDNS) Del(hostname, zone string) (err error) {

	_, err = pdns.pdnsApi(PATHLocalZones+"/"+zone, "PATCH", `{"rrsets": [
			 {"name": "`+hostname+"."+zone+`.",
				"type": "A",
				"changetype": "DELETE"
			 }
		 	]
		 }`)
	return

}
func (pdns *PDNS) GetAllRR(zone string) (rsp []byte, err error) {

	rsp, err = pdns.pdnsApi(PATHLocalZones+"/"+zone, "GET", "")
	return
}
func (pdns *PDNS) ListLocalZones() (rsp []byte, err error) {

	rsp, err = pdns.pdnsApi(PATHLocalZones, "GET", "")
	return
}
func (pdns *PDNS) ListServers() (rsp []byte, err error) {

	rsp, err = pdns.pdnsApi(PATHServers, "GET", "")
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
