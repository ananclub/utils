package utils

type PDNS struct {
	APIUrl string
	APIKey string
}
type Comment struct {
	Content    string `json:"content"`     //– The actual comment
	Account    string `json:"account"`     //– Name of an account that added the comment
	ModifiedAt int    `json:"modified_at"` //– Timestamp of the last change to the comment
}
type Record struct {
	Content  string `json:"content"`
	Disabled bool   `json:"disabled"`
}
type RRSet struct {
	Comments   []Comment `json:"comments"`
	Name       string    `json:"name"`
	Records    []Record  `json:"records"`
	TTL        uint      `json:"ttl"`
	Type       string    `json:"type"`
	ChangeType string    `json:"changetype"`
}
type Zone struct {
	Account          string   `json:"account"`
	ApiRectify       bool     `json:"api_rectify"`
	Dnssec           bool     `json:"dnssec"`
	EditedSerial     int      `json:"edited_serial"`
	Id               string   `json:"id"`
	Kind             string   `json:"kind"`
	LastCheck        int      `json:"last_check"`
	MasterTsigKeyIds []string `json:"master_tsig_key_ids"`
	Masters          []string `json:"masters"`
	Name             string   `json:"name"`
	NotifiedSerial   int      `json:"notified_serial"`
	Nsec3narrow      bool     `json:"nsec3narrow"`
	Nsec3param       string   `json:"nsec3param"`
	Serial           int      `json:"serial"`
	SlaveTsigKeyIds  []string `json:"slave_tsig_key_ids"`
	SoaEdit          string   `json:"soa_edit"`
	SoaEditApi       string   `json:"soa_edit_api"`
	Url              string   `json:"url"`
	RRSets           []RRSet  `json:"rrsets"`
}

const (
	PATHBASE            = "/api/v1"
	PATHServers         = PATHBASE + "/servers"
	PATHServerLocalhost = PATHServers + "/localhost"
	PATHZones           = "/zones"
	PATHLocalZones      = PATHServerLocalhost + PATHZones
)

func (pdns *PDNS) Add(zone, hostname, typ, content string, ttl uint) (err error) {

	rr := RRSet{
		Name:       hostname + "." + zone + ".",
		Type:       typ,
		ChangeType: "REPLACE",
		TTL:        uint(ttl),
		Records:    []Record{Record{Content: content}},
	}
	rrs := map[string]interface{}{"rrsets": []RRSet{rr}}
	b, err := json.Marshal(rrs)
	if err != nil {
		return
	}
	_, err = pdns.pdnsApi(PATHLocalZones+"/"+zone, "PATCH", string(b))
	return

}
func (pdns *PDNS) Del(zone, hostname, typ string) (err error) {
	rr := RRSet{
		Name:       hostname + "." + zone + ".",
		Type:       typ,
		ChangeType: "DELETE",
	}

	rrs := map[string]interface{}{"rrsets": []RRSet{rr}}
	b, err := json.Marshal(rrs)
	if err != nil {
		return
	}
	_, err = pdns.pdnsApi(PATHLocalZones+"/"+zone, "PATCH", string(b))
	return

}
func (pdns *PDNS) GetAllRR(zone string) (rrsets []RRSet, err error) {

	rsp, err := pdns.pdnsApi(PATHLocalZones+"/"+zone, "GET", "")
	if err != nil {
		return
	}
	//println(string(rsp))
	var z = &Zone{}
	err = json.Unmarshal(rsp, z)
	if err != nil {
		return
	}
	rrsets = z.RRSets
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
