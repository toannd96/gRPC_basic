package providers

import (
	"encoding/json"
	"fmt"
)

const (
	magentoName    = "magento"
	magentoBaseURL = "https://magento23demo.connectpos.com"
)

type Magento struct{}

func (p *Magento) Name() string {
	return magentoName
}

func (p *Magento) Get(parameter string) (PlatformInfo, error) {
	pathURL := fmt.Sprintf("%s%s", magentoBaseURL, parameter)
	body, err := httpClient.getMagento(pathURL)
	if err != nil {
		return PlatformInfo{}, err
	}

	var result MagentoResult
	json.Unmarshal(body, &result)

	return result.asPlatformInfo(), nil
}

func (p *Magento) Query(keyword string) (PlatformInfo, error) {
	pathURL := fmt.Sprintf("%s/rest/V1/products?searchCriteria[filterGroups][0][filters][0][field]=name&searchCriteria[filterGroups][0][filters][0][value]=%s", magentoBaseURL, keyword)
	body, err := httpClient.getMagento(pathURL)
	if err != nil {
		return PlatformInfo{}, err
	}

	var result MagentoResult
	json.Unmarshal(body, &result)

	return result.asPlatformInfo(), nil
}

type MagentoResult struct {
	Item []Items `json:"items"`
}

type Items struct {
	Name               string              `json:"name"`
	Sku                string              `json:"sku"`
	Price              float32             `json:"price"`
	Type               string              `json:"type_id"`
	ExtensionAttribute ExtensionAttributes `json:"extension_attributes"`
}

type ExtensionAttributes struct {
	CategoryLink []CategoryLinks `json:"category_links"`
}

type CategoryLinks struct {
	CategoryID int32 `json:"category_id,string"`
}

func (r MagentoResult) asPlatformInfo() PlatformInfo {
	if len(r.Item) == 0 {
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

func (r MagentoResult) getName() string {
	return r.Item[0].Name
}

func (r MagentoResult) getSku() string {
	return r.Item[0].Sku
}

func (r MagentoResult) getPrice() float32 {
	return r.Item[0].Price
}

func (r MagentoResult) getType() string {
	return r.Item[0].Type
}

func (r MagentoResult) getCategories() []int32 {
	var listCategory []int32

	CategoryLinks := r.Item[0].ExtensionAttribute.CategoryLink
	for _, value := range CategoryLinks {
		listCategory = append(listCategory, value.CategoryID)
	}

	return listCategory
}
