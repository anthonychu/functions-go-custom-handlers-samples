package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type InvokeResponse struct {
	Outputs     map[string]interface{}
	Logs        []string
	ReturnValue interface{}
}

type InvokeRequest struct {
	Data     map[string]json.RawMessage
	Metadata map[string]interface{}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hello from go")
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	var invokeRequest InvokeRequest

	d := json.NewDecoder(r.Body)
	d.Decode(&invokeRequest)

	var reqData map[string]interface{}
	json.Unmarshal(invokeRequest.Data["req"], &reqData)

	outputs := make(map[string]interface{})
	outputs["message"] = reqData["Body"]

	resData := make(map[string]interface{})
	resData["body"] = "Message enqueued"
	outputs["res"] = resData
	invokeResponse := InvokeResponse{outputs, nil, ""}
	invokeResponse.Logs = append(invokeResponse.Logs, "hello from go http trigger!")

	responseJson, _ := json.Marshal(invokeResponse)
	fmt.Println(string(responseJson))
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJson)
}

func queueTriggerHandler(w http.ResponseWriter, r *http.Request) {
	var invokeRequest InvokeRequest

	d := json.NewDecoder(r.Body)
	d.Decode(&invokeRequest)

	order := invokeRequest.Data["order"]

	outputs := map[string]interface{}{"": ""}
	invokeResponse := InvokeResponse{outputs, nil, ""}
	invokeResponse.Logs = append(invokeResponse.Logs, "hello from queue trigger")
	invokeResponse.Logs = append(invokeResponse.Logs, string(order))

	responseJson, _ := json.Marshal(invokeResponse)
	fmt.Println(string(responseJson))
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJson)
}

func main() {
	customHandlerPort, exists := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if !exists {
		customHandlerPort = "8080"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/order", orderHandler)
	mux.HandleFunc("/queueTrigger", queueTriggerHandler)
	mux.HandleFunc("/api/hello", helloHandler)
	fmt.Println("Go server Listening on: ", customHandlerPort)
	log.Fatal(http.ListenAndServe(":"+customHandlerPort, mux))
}
