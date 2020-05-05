package providers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	magentoURL  = "https://magento23demo.connectpos.com/rest/V1/products?searchCriteria[pageSize]=1"
	magentoName = "magento"
)

type Magento struct{}

func (p *Magento) Name() string {
	return magentoName
}

func (p *Magento) Query(q string) (PlatformInfo, error) {

	fmt.Println("name platform", p.Name())

	var bearer = "Bearer " + "ca3fokgt3anygnrlxx3ouk3ngguwagnj"
	req, err := http.NewRequest("GET", magentoURL, nil)
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != 200 {
		respErr := fmt.Errorf("Unexpected response: %s", resp.Status)
		log.Fatalf("request failed: %v", respErr)
	}

	defer resp.Body.Close()

	var result MagentoResult

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return PlatformInfo{}, err
	}

	return result.asPlatformInfo(), nil
}

type MagentoResult struct {
	Items []struct {
		Name       string  `json:"name"`
		Sku        string  `json:"sku"`
		Price      float32 `json:"price"`
		Categories []int32 `json:"category_id"`
		Type       string  `json:"type_id"`
	} `json:"items"`
}

func (r MagentoResult) asPlatformInfo() PlatformInfo {
	return PlatformInfo{
		Name:       r.getName(),
		Sku:        r.getSku(),
		Price:      r.getPrice(),
		Type:       r.getType(),
		Categories: r.getCategories(),
	}
}

func (r MagentoResult) getName() string {
	return r.Items[0].Name
}

func (r MagentoResult) getSku() string {
	return r.Items[0].Sku
}

func (r MagentoResult) getPrice() float32 {
	return r.Items[0].Price
}

func (r MagentoResult) getType() string {
	return r.Items[0].Type
}

func (r MagentoResult) getCategories() int32 {
	return r.Items[0].Categories[1]
}
