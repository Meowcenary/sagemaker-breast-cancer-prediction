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
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func featureOrder() []string {
	// Go maps are not ordered, so to preserve the order, must create a separate array
	return []string{
		"mean_radius",
		"mean_texture",
		"mean_perimeter",
		"mean_area",
		"mean_smoothness",
		"mean_compactness",
		"mean_concavity",
		"mean_concave_points",
		"mean_symmetry",
		"mean_fractal_dimension",
		"radius_error",
		"texture_error",
		"perimeter_error",
		"area_error",
		"smoothness_error",
		"compactness_error",
		"concavity_error",
		"concave_points_error",
		"symmetry_error",
		"fractal_dimension_error",
		"worst_radius",
		"worst_texture",
		"worst_perimeter",
		"worst_area",
		"worst_smoothness",
		"worst_compactness",
		"worst_concavity",
		"worst_concave_points",
		"worst_symmetry",
		"worst_fractal_dimension",
	}
}

// These functions are used to set default values
func medianValues() map[string]string {
	return map[string]string{
		"mean_radius":             "13.3",
		"mean_texture":            "18.68",
		"mean_perimeter":          "85.98",
		"mean_area":               "551.7",
		"mean_smoothness":         "0.09462",
		"mean_compactness":        "0.09097",
		"mean_concavity":          "0.06154",
		"mean_concave_points":     "0.03341",
		"mean_symmetry":           "0.1792",
		"mean_fractal_dimension":  "0.06148",
		"radius_error":            "0.3237",
		"texture_error":           "1.095",
		"perimeter_error":         "2.287",
		"area_error":              "24.72",
		"smoothness_error":        "0.00638",
		"compactness_error":       "0.02042",
		"concavity_error":         "0.02615",
		"concave_points_error":    "0.0111",
		"symmetry_error":          "0.01872",
		"fractal_dimension_error": "0.003211",
		"worst_radius":            "14.97",
		"worst_texture":           "25.22",
		"worst_perimeter":         "97.67",
		"worst_area":              "686.6",
		"worst_smoothness":        "0.1309",
		"worst_compactness":       "0.2101",
		"worst_concavity":         "0.2264",
		"worst_concave_points":    "0.09861",
		"worst_symmetry":          "0.2827",
		"worst_fractal_dimension": "0.08006",
		"target":                  "1.0",
	}
}

func meanValues() map[string]string {
	return map[string]string{
		"mean_radius":             "14.117635164835166",
		"mean_texture":            "19.185032967032967",
		"mean_perimeter":          "91.88224175824176",
		"mean_area":               "654.3775824175823",
		"mean_smoothness":         "0.09574402197802198",
		"mean_compactness":        "0.10361931868131868",
		"mean_concavity":          "0.08889814505494506",
		"mean_concave_points":     "0.04827987032967032",
		"mean_symmetry":           "0.18109868131868131",
		"mean_fractal_dimension":  "0.06275676923076923",
		"radius_error":            "0.40201582417582415",
		"texture_error":           "1.2026868131868131",
		"perimeter_error":         "2.858253406593406",
		"area_error":              "40.0712989010989",
		"smoothness_error":        "0.006989074725274725",
		"compactness_error":       "0.025635448351648354",
		"concavity_error":         "0.03282367230769231",
		"concave_points_error":    "0.01189394065934066",
		"symmetry_error":          "0.02057351208791209",
		"fractal_dimension_error": "0.003820455604395604",
		"worst_radius":            "16.235103296703297",
		"worst_texture":           "25.53569230769231",
		"worst_perimeter":         "107.10312087912088",
		"worst_area":              "876.9870329670329",
		"worst_smoothness":        "0.13153213186813187",
		"worst_compactness":       "0.25274180219780223",
		"worst_concavity":         "0.27459456923076925",
		"worst_concave_points":    "0.11418222197802198",
		"worst_symmetry":          "0.29050219780219777",
		"worst_fractal_dimension": "0.08386784615384615",
		"target":                  "0.6285714285714286",
	}
}

func orderFeatures(queryParams url.Values) []string {
	// for now assigning one of the default maps directly will be fine
	defaultValues := medianValues()

	var record []string
	for _, feature := range featureOrder() {
		if val, exists := queryParams[feature]; exists {
			// Use provided value
			record = append(record, val[0])
		} else {
			 // Use default value
			record = append(record, defaultValues[feature])
		}
	}

	return record
}

// This is not a great approach, but it makes it so we can reuse existing code easily
func decodeJsonToValues(r *http.Request) (url.Values, error) {
	var requestData map[string]string
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		return nil, err
	}

	// Convert map[string]string to url.Values (map[string][]string)
	values := url.Values{}
	for key, val := range requestData {
		// `Set` ensures each value is stored as a single-element slice
		values.Set(key, val)
	}

	return values, nil
}

// record is csv row which is necessary for prediction because sagemaker was
// trained on csv data
func performPrediction(w http.ResponseWriter, buf bytes.Buffer) string {
	// Load AWS credentials and config - this neeeds to be set in ~/.aws/credentials
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	if err != nil {
		http.Error(w, "Failed to load AWS config", http.StatusInternalServerError)
		return ""
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
		return ""
	}

	fmt.Println("Response:", &resp.ContentType)
	// Convert response body to a string
	predictionStr := string(resp.Body)
	if len(predictionStr) == 0 {
		http.Error(w, "Empty response from SageMaker", http.StatusInternalServerError)
		return ""
	}

	return string(resp.Body)
}

func predictHandler(w http.ResponseWriter, r *http.Request) {
	incrementRequestCount(r)
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
	fmt.Printf("Query Params: %+v\n", queryParams)
	record := orderFeatures(queryParams)
	fmt.Println("Record: ", record)

	// Convert to CSV format
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	err := writer.Write(record)
	if err != nil {
		http.Error(w, "Failed to write CSV data", http.StatusInternalServerError)
		return
	}
	writer.Flush()

	predictionStr := performPrediction(w, buf)

	// Return prediction value as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"prediction": predictionStr})
}

func predictHandlerJson(w http.ResponseWriter, r *http.Request) {
	incrementRequestCount(r)
	requestData, err := decodeJsonToValues(r)

	// Ensure we have at least one feature
	if len(requestData) == 0 {
		http.Error(w, "Missing JSON fields", http.StatusBadRequest)
		return
	}

	fmt.Printf("Request Data: %+v\n", requestData)

	// Convert JSON input to a CSV row
	record := orderFeatures(requestData)
	fmt.Println("Record: ", record)

	// Convert to CSV format
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	err = writer.Write(record)
	if err != nil {
		http.Error(w, "Failed to write CSV data", http.StatusInternalServerError)
		return
	}
	writer.Flush()

	// Load AWS credentials and config
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	if err != nil {
		http.Error(w, "Failed to load AWS config", http.StatusInternalServerError)
		return
	}

	// Create SageMaker runtime client
	sageMakerClient := sagemakerruntime.NewFromConfig(cfg)

	// Call SageMaker endpoint
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

	fmt.Println("Response:", &resp.ContentType)

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

func apiStatusHandler(w http.ResponseWriter, r *http.Request) {
	incrementRequestCount(r)
	// Return prediction value as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"code": "200"})
}

// Use this to increment request counts for each endpoint
func incrementRequestCount(r *http.Request) {
	label := fmt.Sprintf("%s - %s", r.URL.Path, r.Method)
	requestCount.WithLabelValues(label).Inc()
}

var requestCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "requests_total",
		Help: "Total number of requests",
	},
	[]string{"method"},
)

func init() {
	// having this call in init() ensures it registers properly
	prometheus.MustRegister(requestCount)
	// Set to 0
	requestCount.WithLabelValues("GET - ").Add(0)
}

func main() {
	// router is similar to demo video, but different handlers
	router := mux.NewRouter()
	router.HandleFunc("/status", apiStatusHandler).Methods("GET")
	router.HandleFunc("/predict", predictHandler).Methods("GET")
	router.HandleFunc("/predict/json", predictHandlerJson).Methods("GET")
	router.Handle("/metrics", promhttp.Handler())

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
