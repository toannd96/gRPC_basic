package providers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	bigcommerceURL  = "https://api.bigcommerce.com/stores/jp6bmqxeb9/v3/catalog/products"
	bigcommerceName = "bigcommerce"
)

type BigCommerce struct{}

func (p *BigCommerce) Name() string {
	return bigcommerceName
}

func (p *BigCommerce) Query(q string) (PlatformInfo, error) {

	fmt.Println("name platform", p.Name())

	req, err := http.NewRequest("GET", bigcommerceURL, nil)
	req.Header.Set("X-Auth-Client", "hyr8hgnlkx54l16gn95v74kkmgnm518")
	req.Header.Set("X-Auth-Token", "d35jic4f6mxr59fsabmwuiguiss1yg2")

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

	var result BigCommerceResult

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return PlatformInfo{}, err
	}
	return result.asPlatformInfo(), nil
}

type BigCommerceResult struct {
	Data []struct {
		Name       string  `json:"name"`
		Sku        string  `json:"sku"`
		Price      float32 `json:"price"`
		Categories []int32 `json:"categories"`
		Type       string  `json:"type"`
	} `json:"data"`
}

func (r BigCommerceResult) asPlatformInfo() PlatformInfo {
	return PlatformInfo{
		Name:       r.getName(),
		Sku:        r.getSku(),
		Price:      r.getPrice(),
		Type:       r.getType(),
		Categories: r.getCategories(),
	}
}

func (r BigCommerceResult) getName() string {
	return r.Data[0].Name
}

func (r BigCommerceResult) getSku() string {
	return r.Data[0].Sku
}

func (r BigCommerceResult) getPrice() float32 {
	return r.Data[0].Price
}

func (r BigCommerceResult) getType() string {
	return r.Data[0].Type
}

func (r BigCommerceResult) getCategories() int32 {
	return r.Data[0].Categories[1]
}
