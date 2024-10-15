package main

import (
	dpdService "RxHub/dpdService"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func getByBrandName(response http.ResponseWriter, request *http.Request) {
	brandName := request.URL.Query().Get("brandName")

	if brandName == "" {
		http.Error(response, "Brand name is required", http.StatusBadRequest)
		return
	}

	products := dpdService.GetDrugByBrandName(brandName)

	response.Header().Set("Content-Type", "application/json")

	jsonResponse, err := json.Marshal(products)

	if err != nil {
		http.Error(response, "Error processing data", http.StatusInternalServerError)
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Write(jsonResponse)
}

func main() {
	http.HandleFunc("/getByBrandName", getByBrandName)

	err := http.ListenAndServe(":3333", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server closed \n")
	} else if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
