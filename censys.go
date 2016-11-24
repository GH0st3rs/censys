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

func getErrorString(StatusCode int) error {
	switch StatusCode {
	case 400:
		return fmt.Errorf("Error %d -> Query could not be parsed", StatusCode)
	case 404:
		return fmt.Errorf("Error %d -> Page not found", StatusCode)
	case 429:
		return fmt.Errorf("Error %d -> Rate limit exceeded", StatusCode)
	case 500:
		return fmt.Errorf("Error %d -> Internal server error", StatusCode)
	default:
		return fmt.Errorf("unknown error code %d", StatusCode)
	}
}

//Search Engine function
func Search(auth [2]string, index string, query string, page int) (*[]byte, error) {
	//Make request URL and post param
	url := fmt.Sprintf("%s/%s/%s", APIURL, SearchURL, index)
	//	fields := ``
	paramStr := fmt.Sprintf(`{"query": "%s", "page": %d, "fields":[]}`, query, page)
	data := bytes.NewBuffer([]byte(paramStr))
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		panic(err)
	}
	//Set Headers
	req.SetBasicAuth(auth[0], auth[1])
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		return nil, getErrorString(resp.StatusCode)
	}
	defer resp.Body.Close()
	//Read body request
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &body, nil
}

//Search IPv4
func SearchIPv4(auth [2]string, query string, page int) (*censysSearchIPv4, error) {
	const index = "ipv4"

	body, err := Search(auth, index, query, page)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(*body))

	cS := censysSearchIPv4{}
	if err = json.Unmarshal(*body, &cS); err != nil {
		return nil, err
	}

	return &cS, nil
}

//Search WebSites
func SearchWebSites(auth [2]string, query string, page int) (*censysSearchWebsites, error) {
	const index = "websites"

	body, err := Search(auth, index, query, page)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(*body))

	cS := censysSearchWebsites{}
	if err = json.Unmarshal(*body, &cS); err != nil {
		return nil, err
	}

	return &cS, nil
}

//Search Certificates
func SearchCertificates(auth [2]string, query string, page int) (*censysSearchCertificates, error) {
	const index = "certificates"

	body, err := Search(auth, index, query, page)
	if err != nil {
		return nil, err
	}

	cS := censysSearchCertificates{}
	if err = json.Unmarshal(*body, &cS); err != nil {
		return nil, err
	}

	return &cS, nil
}
