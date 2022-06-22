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
type RREntry struct {
	Content  string `json:"content"`
	Disabled bool   `json:"disabled"`
	SetPtr   bool   `json:"set-ptr"`
}
type RRSet struct {
	Comments   []Comment `json:"comments"`
	Name       string    `json:"name"`
	Records    []RREntry `json:"records"`
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

type RRSets struct {
	RRSets []RRSet `json:"rrsets"`
}

/*
{
  "type": "Server",
  "id": "localhost",
  "url": "/api/v1/servers/localhost",
  "daemon_type": "recursor",
  "version": "4.1.0",
  "config_url": "/api/v1/servers/localhost/config{/config_setting}",
  "zones_url": "/api/v1/servers/localhost/zones{/zone}",
}
*/
type Server struct {
	Type       string `json:"type"`
	Id         string `json:"id"`
	Url        string `json:"url"`
	DaemonType string `json:"daemon_type"`
	Version    string `json:"version"`
	ConfigUrl  string `json:"config_url"`
	ZoneUrl    string `json:"zones_url"`
}

const (
	PATHBASE            = "/api/v1"
	PATHServers         = PATHBASE + "/servers"
	PATHServerLocalhost = PATHServers + "/localhost"
	PATHZones           = "/zones"
	PATHLocalZones      = PATHServerLocalhost + PATHZones
)

func (pdns *PDNS) Add(zone, hostname, typ, content string, ttl uint) (err error) {

	rrs := RRSets{RRSets: []RRSet{
		RRSet{
			Name:       hostname + "." + zone + ".",
			Type:       typ,
			ChangeType: "REPLACE",
			TTL:        ttl,
			Records:    []RREntry{RREntry{Content: content}},
		}}}
	b, err := json.Marshal(rrs)
	if err != nil {
		return
	}
	_, err = pdns.Api(PATHLocalZones+"/"+zone, "PATCH", string(b))
	return

}
func (pdns *PDNS) Del(zone, hostname, typ string) (err error) {
	rrs := RRSets{RRSets: []RRSet{
		RRSet{
			Name:       hostname + "." + zone + ".",
			Type:       typ,
			ChangeType: "DELETE",
		}}}

	b, err := json.Marshal(rrs)
	if err != nil {
		return
	}
	_, err = pdns.Api(PATHLocalZones+"/"+zone, "PATCH", string(b))
	return

}
func (pdns *PDNS) GetAllRR(zone string) (rrsets []RRSet, err error) {

	rsp, err := pdns.Api(PATHLocalZones+"/"+zone, "GET", "")
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

	rsp, err = pdns.Api(PATHLocalZones, "GET", "")
	return
}
func (pdns *PDNS) ListServers() (rsp []byte, err error) {

	rsp, err = pdns.Api(PATHServers, "GET", "")
	return
}
func (pdns *PDNS) Api(path, method, data string) (rsp []byte, err error) {

	h := make(map[string]string)
	h["X-API-Key"] = pdns.APIKey

	if _, rsp, err = HttpDo(pdns.APIUrl+path, method, nil, nil, h, nil, []byte(data)); err != nil {
		return
	}

	return
}
