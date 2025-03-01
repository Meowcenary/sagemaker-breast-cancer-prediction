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
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sagemakerruntime"
	"github.com/gorilla/mux"
)


func orderFeatures(queryParams url.Values) []string {
	// TODO This should probably be updated to the median value for each feature,
	// but for now, using 0.0 will have to do
	featureOrder := map[string]string{
		"mean":                    "0.0",
		"radius":                  "0.0",
		"mean_texture":            "0.0",
		"mean_perimeter":          "0.0",
		"mean_area":               "0.0",
		"mean_smoothness":         "0.0",
		"mean_compactness":        "0.0",
		"mean_concavity":          "0.0",
		"mean_concave_points":     "0.0",
		"mean_symmetry":           "0.0",
		"mean_fractal_dimension":  "0.0",
		"radius_error":            "0.0",
		"texture_error":           "0.0",
		"perimeter_error":         "0.0",
		"area_error":              "0.0",
		"smoothness_error":        "0.0",
		"compactness_error":       "0.0",
		"concavity_error":         "0.0",
		"concave_points_error":    "0.0",
		"symmetry_error":          "0.0",
		"fractal_dimension_error": "0.0",
		"worst_radius":            "0.0",
		"worst_texture":           "0.0",
		"worst_perimeter":         "0.0",
		"worst_area":              "0.0",
		"worst_smoothness":        "0.0",
		"worst_compactness":       "0.0",
		"worst_concavity":         "0.0",
		"worst_concave_points":    "0.0",
		"worst_symmetry":          "0.0",
		"worst_fractal_dimension": "0.0",
		"target":                  "0.0",
	}

	var record []string
	for _, feature := range featureOrder {
		if val, exists := queryParams[feature]; exists {
			// Use provided value
			record = append(record, val[0])
		} else {
			 // Use default value
			record = append(record, featureOrder[feature])
		}
	}

	return record
}

func predictHandler(w http.ResponseWriter, r *http.Request) {
	// Get the query params
	queryParams := r.URL.Query()

	// should have at least one param
	if len(queryParams) == 0 {
		http.Error(w, "Missing query parameters", http.StatusBadRequest)
		return
	}

	// Parse query parameters to CSV row - this assumes feature values are
	// passed as `?feature1=val1&feature2=val2...) where feature is subbed out for
	// values in var featureOrder
	var record = orderFeatures(queryParams)

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

	// Call SageMaker endpoint - need to use text/csv because that is what it was trained with
	endpointName := "sagemaker-xgboost-2025-02-23-22-17-09-203"
	resp, err := sageMakerClient.InvokeEndpoint(context.Background(), &sagemakerruntime.InvokeEndpointInput{
		EndpointName: &endpointName,
		Body:         buf.Bytes(),
		ContentType:  aws.String("text/csv"),
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to call SageMaker: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert response body to a string
	predictionStr := string(resp.Body)
	if len(predictionStr) == 0 {
		http.Error(w, "Empty response from SageMaker", http.StatusInternalServerError)
		return
	}

	// Return prediction value as JSON
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
