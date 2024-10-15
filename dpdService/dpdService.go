package dpdService

import (
	drugProducts "RxHub/dpdService/drugProducts"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func getAllDrugProducts() []drugProducts.DrugProduct {
	const drugProductURL = "https://health-products.canada.ca/api/drug/drugproduct/?lang=en&type=json"
	res, err := http.Get(drugProductURL)

	if err != nil {
		fmt.Printf("There was an error trying to do GET %s\n", err)
		os.Exit(1)
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Printf("Could not read response body")
		os.Exit(1)
	}

	var products []drugProducts.DrugProduct
	err = json.Unmarshal(resBody, &products)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
	}

	return products
}

func GetDrugByBrandName(searchName string) []drugProducts.DrugProduct {
	//start := time.Now()
	var products = getAllDrugProducts()
	//apiEnd := time.Since(start)
	ingredMap := GetActiveIngredient()
	//ingredMapEnd := time.Since(start)

	var resultProducts []drugProducts.DrugProduct

	for _, product := range products {
		if strings.Contains(strings.ToLower(product.BrandName), strings.ToLower(searchName)) {
			product.ActiveIngredients = ingredMap[product.DrugCode]
			resultProducts = append(resultProducts, product)
		}
	}

	fmt.Println(resultProducts)
	//totalTimeTook := time.Since(start)

	//fmt.Printf("The time it took was %f\n The getting the drugs api took %f Getting the ingredients took %f\n", totalTimeTook.Seconds(), apiEnd.Seconds(), ingredMapEnd.Seconds())

	return resultProducts
}

func GetActiveIngredient() map[uint32][]drugProducts.ActiveIngredient {
	var activeIngredientUrl = "https://health-products.canada.ca/api/drug/activeingredient/?lang=en&type=json"

	var activeIngredients []drugProducts.ActiveIngredient

	res, err := http.Get(activeIngredientUrl)

	if err != nil {
		fmt.Printf("There was an error trying to do GET %s\n", err)
		os.Exit(1)
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Printf("Could not read response body")
		os.Exit(1)
	}

	err = json.Unmarshal(resBody, &activeIngredients)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
	}

	ingredMap := make(map[uint32][]drugProducts.ActiveIngredient)

	for _, ingredient := range activeIngredients {
		ingredMap[ingredient.DrugCode] = append(ingredMap[ingredient.DrugCode], ingredient)
	}

	return ingredMap
}
