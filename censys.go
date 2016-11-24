// censys project censys.go
package censys

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const APIURL = "https://www.censys.io/api/v1"
const SearchURL = "search"
const ExportURL = "export"

type Metadata struct {
	Count int    `json:"count"`
	Query string `json:"query"`
	Page  int    `jsno:"page"`
	Pages int    `json:"pages"`
}

type censysSearchIPv4 struct {
	Status   string   `json:"status"`
	MetaData Metadata `json:"metadata"`
	Results  []struct {
		IP        string   `json:"ip"`
		Protocols []string `json:"protocols"`
	} `json:"results"`
}

type censysSearchWebsites struct {
	Status   string   `json:"status"`
	MetaData Metadata `json:"metadata"`
	Results  []struct {
		Domain string `json:"domain"`
		Rank   []int  `json:"alexa_rank"`
	} `json:"results"`
}

type censysSearchCertificates struct {
	Status   string   `json:"status"`
	MetaData Metadata `json:"metadata"`
	Results  []struct {
		FingerprintSHA256 []string `json:"parsed.fingerprint_sha256"`
		SubjectDN         []string `json:"parsed.subject_dn"`
		IssuerDN          []string `json:"parsed.issuer_dn"`
	} `json:"results"`
}

type censysExport struct {
	Status string `json:"status"`
	Config struct {
		Format   string `json:"format"`
		Compress bool   `json:"compress"`
		Headers  bool   `json:"headers"`
		Flatten  bool   `json:"flatten"`
		Query    string `json:"query"`
	} `json:"configuration"`
	JobID string `json:"job_id"`
}

func getErrorString(StatusCode int, url, paramStr string) error {
	switch StatusCode {
	case 400:
		return fmt.Errorf("Error %d -> Query could not be parsed. Url: %s. Param: %s", StatusCode, url, paramStr)
	case 404:
		return fmt.Errorf("Error %d -> Page not found. Url: %s. Param: %s", StatusCode, url, paramStr)
	case 429:
		return fmt.Errorf("Error %d -> Rate limit exceeded. Url: %s. Param: %s", StatusCode, url, paramStr)
	case 500:
		return fmt.Errorf("Error %d -> Internal server error. Url: %s. Param: %s", StatusCode, url, paramStr)
	default:
		return fmt.Errorf("unknown error code %d. Url: %s. Param: %s", StatusCode, url, paramStr)
	}
}

func request(auth [2]string, url, paramStr *string) (*[]byte, error) {
	data := bytes.NewBuffer([]byte(*paramStr))
	req, err := http.NewRequest("POST", *url, data)
	if err != nil {
		panic(err)
	}
	//Set Headers
	req.SetBasicAuth(auth[0], auth[1])
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		return nil, getErrorString(resp.StatusCode, *url, *paramStr)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

//Search Engine function
func Search(auth [2]string, index string, query string, page int) (*[]byte, error) {
	//Make request URL and post param
	url := fmt.Sprintf("%s/%s/%s", APIURL, SearchURL, index)
	//	fields := ``
	paramStr := fmt.Sprintf(`{"query": "%s", "page": %d, "fields":[]}`, query, page)

	//Read body request
	body, err := request(auth, &url, &paramStr)
	if err != nil {
		return nil, err
	}

	return body, nil
}

//Search IPv4
func SearchIPv4(auth [2]string, query string, page int) (*censysSearchIPv4, error) {

	body, err := Search(auth, "ipv4", query, page)
	if err != nil {
		return nil, err
	}

	cS := censysSearchIPv4{}
	if err = json.Unmarshal(*body, &cS); err != nil {
		return nil, err
	}

	return &cS, nil
}

//Search WebSites
func SearchWebSites(auth [2]string, query string, page int) (*censysSearchWebsites, error) {

	body, err := Search(auth, "websites", query, page)
	if err != nil {
		return nil, err
	}

	cS := censysSearchWebsites{}
	if err = json.Unmarshal(*body, &cS); err != nil {
		return nil, err
	}

	return &cS, nil
}

//Search Certificates
func SearchCertificates(auth [2]string, query string, page int) (*censysSearchCertificates, error) {

	body, err := Search(auth, "certificates", query, page)
	if err != nil {
		return nil, err
	}

	cS := censysSearchCertificates{}
	if err = json.Unmarshal(*body, &cS); err != nil {
		return nil, err
	}

	return &cS, nil
}

//Export Engine function
func Export(auth [2]string, query string) (*censysExport, error) {
	//Make request URL and post param
	url := fmt.Sprintf("%s/%s", APIURL, ExportURL)
	paramStr := fmt.Sprintf(`{"query": "%s", "format":"json"}`, query)

	//Read body request
	body, err := request(auth, &url, &paramStr)
	if err != nil {
		return nil, err
	}

	cE := censysExport{}
	if err = json.Unmarshal(*body, &cE); err != nil {
		return nil, err
	}

	return &cE, nil
}

func GetExportStatus(auth [2]string, job_id string) (*[]byte, error) {
	//Make request URL and post param
	url := fmt.Sprintf("%s/%s", APIURL, ExportURL)

	//Read body request
	body, err := request(auth, &url, &job_id)
	if err != nil {
		return nil, err
	}

	return body, nil
}
