package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (c *WordpressCollector) FetchJSONFromEndpoint(APIEndpoint string) []byte {
	APIBase := c.Wp.MonitoredWordpress
	HTTPClient := &http.Client{}
	fetchURL := fmt.Sprintf("%s%s", APIBase, APIEndpoint)
	fmt.Println(fetchURL)
	request, err := http.NewRequest("GET", fetchURL, nil)
	request.Header.Set("User-Agent", c.Wp.UserAgent)
	ErrCheck(err)
	if c.Wp.Auth.Use {
		request.SetBasicAuth(c.Wp.Auth.Username, c.Wp.Auth.Password)
	}
	response, err := HTTPClient.Do(request)
	ErrCheck(err)
	data, _ := io.ReadAll(response.Body)
	if err == nil && response.StatusCode != http.StatusOK {
		fmt.Printf("Error status: %d: %s", response.StatusCode, string(data))
		return nil
	}

	return data
}

// count items returned in JSON and return length
func CountJSONItems(JSONResponse []byte) (int, error) {
	var err error
	var JSONObject interface{}
	json.Unmarshal(JSONResponse, &JSONObject)

	JSONObjectSlice, isOK := JSONObject.([]interface{})
	if !isOK {
		// try as map
		JSONObjectMap, isOK2 := JSONObject.(map[string]interface{})
		if isOK2 {
			return len(JSONObjectMap), err
		}
		err = fmt.Errorf("cannot convert the JSON object")
		// return -1 if json cannot be parsed properly
		return -1, err
	}

	return len(JSONObjectSlice), err
}

func BasicAuth(username, password string) string {
	authString := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(authString))
}

func ErrCheck(e error) {
	if e != nil {
		log.Println(e)
	}
}
