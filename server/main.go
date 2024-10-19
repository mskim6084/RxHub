package main

import (
	dpdService "RxHub/dpdService"
	middleware "RxHub/server/middleware"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type RequestBody struct {
	Key   string `json:"name"`
	Value string `json:"Value"`
}

func middlewareChain(handler http.Handler, middleswares ...func(http.Handler) http.Handler) http.Handler {
	for _, m := range middleswares {
		handler = m(handler)
	}

	return handler
}

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

func getByBrandNameWithWorkers(response http.ResponseWriter, request *http.Request) {
	brandName := request.URL.Query().Get("brandName")

	if brandName == "" {
		http.Error(response, "Brand name is required", http.StatusBadRequest)
		return
	}

	products := dpdService.GetDrugByBrandNameWithWorkers(brandName)

	//fmt.Printf("%f\n", products)

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

func addUser(response http.ResponseWriter, request *http.Request) {
	var reqBody RequestBody

	if err := json.NewDecoder(request.Body).Decode(&reqBody); err != nil {
		fmt.Printf("Invalid JSON %s\n", err)
	}

	fmt.Printf("The request body was %s\n", reqBody)

	defer request.Body.Close()

	resToClient := map[string]interface{}{
		"status":          200,
		"responseMessage": "Successfully added!",
	}

	jsonData, err := json.Marshal(resToClient)

	if err != nil {
		fmt.Printf("Could not convert to JSON %s\n", err)
	}

	response.Write(jsonData)
}

func main() {
	brandNameHandler := http.HandlerFunc(getByBrandName)

	http.Handle("/getByBrandName", middlewareChain(brandNameHandler, middleware.RateLimiter))
	http.HandleFunc("/getByBrandNameWithWorkers", getByBrandNameWithWorkers)
	http.HandleFunc("/addUser", addUser)

	err := http.ListenAndServe(":3333", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server closed \n")
	} else if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
