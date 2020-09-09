package packagecloud

import (
	"encoding/json"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Packagecloud struct {
	Token   string
	Repo	string
	User	string
	Value	int `json:"Value"`
}

var PackagecloudConfig = `
	token = ""
	user = ""
	repo = ""
`

func (s *Packagecloud) SampleConfig() string {
	return PackagecloudConfig
}

func (s *Packagecloud) Description() string {
	return "Fetching statistics from Packagecloud"
}

func (s *Packagecloud) Gather(acc telegraf.Accumulator) error {
	token := s.Token
	repo  := strings.ReplaceAll(s.Repo, "-", "_")
	user  := strings.ReplaceAll(s.User, "-", "_")

	url := "https://"+token+":@packagecloud.io/api/v1/repos/"+user+"/"+repo+"/stats/installs/count.json"

	packagecloudClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := packagecloudClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	downloadCount := Packagecloud{}

	jsonErr := json.Unmarshal(body, &downloadCount)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	fields := make(map[string]interface{})
	fields[user + "_" + repo + "_downloads"] = downloadCount.Value

	tags := make(map[string]string)

	acc.AddFields("packagecloud", fields, tags)

	return nil
}

func init() {
	inputs.Add("packagecloud", func() telegraf.Input { return &Packagecloud{Value: 0} })
}

