package providers

import (
	"encoding/json"
	"fmt"
)

const (
	bigcommerceName    = "bigcommerce"
	bigcommerceBaseURL = "https://api.bigcommerce.com/stores/jp6bmqxeb9"
)

type BigCommerce struct{}

func (p *BigCommerce) Name() string {
	return bigcommerceName
}

func (p *BigCommerce) Get(parameter string) (PlatformInfo, error) {
	pathURL := fmt.Sprintf("%s%s", bigcommerceBaseURL, parameter)
	body, err := httpClient.getBigCommerce(pathURL)
	if err != nil {
		return PlatformInfo{}, err
	}

	var result BigCommerceResult
	json.Unmarshal(body, &result)

	return result.asPlatformInfo(), nil
}

func (p *BigCommerce) Query(keyword string) (PlatformInfo, error) {
	queryURL := fmt.Sprintf("%s/v3/catalog/products?keyword=%s", bigcommerceBaseURL, keyword)
	body, err := httpClient.getBigCommerce(queryURL)
	if err != nil {
		return PlatformInfo{}, err
	}

	var result BigCommerceResult
	json.Unmarshal(body, &result)

	return result.asPlatformInfo(), nil
}

type BigCommerceResult struct {
	Data []Datas `json:"data"`
}

type Datas struct {
	Name       string  `json:"name"`
	Sku        string  `json:"sku"`
	Price      float32 `json:"price"`
	Categories []int32 `json:"categories,string"`
	Type       string  `json:"type"`
}

func (r BigCommerceResult) asPlatformInfo() PlatformInfo {
	if len(r.Data) == 0 {
		return PlatformInfo{}
	}

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

func (r BigCommerceResult) getCategories() []int32 {
	return r.Data[0].Categories
}
