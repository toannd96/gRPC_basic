package providers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	magentoName = "magento"
)

type Magento struct{}

func (p *Magento) Name() string {
	return magentoName
}

func (p *Magento) Get(api string) (PlatformInfo, error) {

	var bearer = "Bearer " + os.Getenv("ACCESS_TOKEN_MAGENTO")
	req, err := http.NewRequest("GET", api, nil)
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
		Name                string  `json:"name"`
		Sku                 string  `json:"sku"`
		Price               float32 `json:"price"`
		ExtensionAttributes struct {
			CategoryLinks []struct {
				CategoryID int32 `json:"category_id,string"`
			} `json:"category_links"`
		} `json:"extension_attributes"`
		Type string `json:"type_id"`
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

func (r MagentoResult) getCategories() []int32 {
	var listCategory []int32
	CategoryLinks := r.Items[0].ExtensionAttributes.CategoryLinks
	for _, value := range CategoryLinks {
		listCategory = append(listCategory, value.CategoryID)
	}
	return listCategory
}
