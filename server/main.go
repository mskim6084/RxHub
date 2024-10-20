package main

import (
	"RxHub/client/user"
	users "RxHub/client/user"
	dpdService "RxHub/dpdService"
	middleware "RxHub/server/middleware"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

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
	var addUserBody users.User

	if err := json.NewDecoder(request.Body).Decode(&addUserBody); err != nil {
		fmt.Printf("Invalid JSON %s\n", err)
	}

	fmt.Printf("The request body was %s\n", addUserBody)

	users.CreatedUsers[users.UserCount] = addUserBody

	users.UserCount += 1

	fmt.Printf("New user added\n")
	fmt.Println(users.CreatedUsers)

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

func checkHealth(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(response, "Wrong Method", http.StatusMethodNotAllowed)
		return
	}

	resToClient := map[string]interface{}{
		"status":          200,
		"responseMessage": "Health is good",
	}

	jsonData, err := json.Marshal(resToClient)

	if err != nil {
		fmt.Printf("Could not convert to JSON %s\n", err)
	}

	response.Write(jsonData)
}

func deleteUser(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodDelete {
		http.Error(response, "/deleteUser got Wrong method\n", http.StatusMethodNotAllowed)
		return
	}

	userIdStr := request.URL.Query().Get("userId")

	userId, err := strconv.ParseInt(userIdStr, 10, 64)

	if err != nil {
		fmt.Printf("Couldn't delete user %s\n", err)
		http.Error(response, "Could not delete user, please try again", http.StatusInternalServerError)
		return
	}

	delete(users.CreatedUsers, int(userId))

	fmt.Println(user.CreatedUsers)

	resToClient := map[string]interface{}{
		"status":          200,
		"responseMessage": "Successfully deleted a user",
	}

	jsonData, err := json.Marshal(resToClient)

	if err != nil {
		fmt.Printf("Could not convert to JSON %s\n", err)
		http.Error(response, "Could not delete user, please try again", http.StatusInternalServerError)
	}

	response.Write(jsonData)
}

func main() {
	brandNameHandler := http.HandlerFunc(getByBrandName)
	brandNameWithWorkerHandler := http.HandlerFunc(getByBrandNameWithWorkers)

	http.Handle("/getByBrandName", middlewareChain(brandNameHandler, middleware.TimeBasedRateLimiter))
	http.Handle("/getByBrandNameWithWorkers", middlewareChain(brandNameWithWorkerHandler, middleware.TokenBucketRateLimiter))
	http.HandleFunc("/addUser", addUser)
	http.HandleFunc("/health", checkHealth)
	http.HandleFunc("/deleteUser", deleteUser)

	err := http.ListenAndServe(":3333", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server closed \n")
	} else if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
