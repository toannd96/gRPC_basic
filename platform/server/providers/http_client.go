package providers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	errUnexpectedRespnse = "unexpected response: %s"
)

type HTTPClient struct{}

var (
	httpClient = HTTPClient{}
)

func (c HTTPClient) getBigCommerce(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Auth-Client", os.Getenv("CLIENT_ID"))
	req.Header.Set("X-Auth-Token", os.Getenv("ACCESS_TOKEN_BIGCOMMERCE"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	c.info(fmt.Sprintf("GET %s -> %d", url, resp.StatusCode))

	if resp.StatusCode != 200 {
		respErr := fmt.Errorf(errUnexpectedRespnse, resp.Status)
		c.info(fmt.Sprintf("request failed: %v", respErr))
		return nil, respErr
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (c HTTPClient) getMagento(url string) ([]byte, error) {
	var bearer = "Bearer " + os.Getenv("ACCESS_TOKEN_MAGENTO")
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	c.info(fmt.Sprintf("GET %s -> %d", url, resp.StatusCode))

	if resp.StatusCode != 200 {
		respErr := fmt.Errorf(errUnexpectedRespnse, resp.Status)
		c.info(fmt.Sprintf("request failed: %v", respErr))
		return nil, respErr
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (c HTTPClient) info(msg string) {
	log.Printf("[JSONClient] %s\n", msg)
}
