// example for calling api: curl "http://localhost:8080/predict?feature1=5.1&feature2=3.5&feature3=1.4&feature4=0.2"
package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sagemakerruntime"
	"github.com/gorilla/mux"
)

func predictHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters - this assumes feature values are passed as `?feature1=val1&feature2=val2...)
	queryParams := r.URL.Query()
	// should have at least one param
	if len(queryParams) == 0 {
		http.Error(w, "Missing query parameters", http.StatusBadRequest)
		return
	}

	// Convert query parameters to a CSV row
	var record []string
	for _, values := range queryParams {
    	// Get the first value for each parameter
		record = append(record, values[0])
	}

	// Convert to CSV format
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	err := writer.Write(record)
	if err != nil {
		http.Error(w, "Failed to write CSV data", http.StatusInternalServerError)
		return
	}
	writer.Flush()

	// Load AWS credentials and config - this neeeds to be set in ~/.aws/credentials
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	if err != nil {
		http.Error(w, "Failed to load AWS config", http.StatusInternalServerError)
		return
	}

	// Create SageMaker runtime client
	sageMakerClient := sagemakerruntime.NewFromConfig(cfg)

	endpointName := "sagemaker-xgboost-2025-02-23-22-17-09-203"
	// Call SageMaker endpoint - need to use text/csv because that is what it was trained with
	resp, err := sageMakerClient.InvokeEndpoint(context.Background(), &sagemakerruntime.InvokeEndpointInput{
		EndpointName: &endpointName,
		Body:         buf.Bytes(),
		ContentType:  aws.String("text/csv"),
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to call SageMaker: %v", err), http.StatusInternalServerError)
		return
	}

	// Read and decode SageMaker response
	var result map[string]interface{}
	err = json.NewDecoder(bytes.NewReader(resp.Body)).Decode(&result)
	if err != nil {
	    http.Error(w, "Failed to parse response from SageMaker", http.StatusInternalServerError)
	    return
	}

	// Extract the prediction value (assuming it is inside a key like "predictions")
	prediction, ok := result["predictions"]
	if !ok {
	    http.Error(w, "Prediction key not found in response", http.StatusInternalServerError)
	    return
	}

	// Convert prediction to string
	predictionStr := fmt.Sprintf("%v", prediction)

	// Send response back to client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"prediction": predictionStr})
}

func main() {
  // router is similar to demo video, but different handlers
	r := mux.NewRouter()
	r.HandleFunc("/predict", predictHandler).Methods("GET")

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
