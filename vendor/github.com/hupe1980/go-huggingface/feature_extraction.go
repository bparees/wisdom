package gohuggingface

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// Request structure for the feature extraction endpoint
type FeatureExtractionRequest struct {
	// String to get the features from
	Inputs  []string `json:"inputs"`
	Options Options  `json:"options,omitempty"`
	Model   string   `json:"-"`
}

// Response structure for the feature extraction endpoint
type FeatureExtractionResponse [][][][]float64

// FeatureExtraction performs feature extraction using the specified model.
// It sends a POST request to the Hugging Face inference endpoint with the provided input data.
// The response contains the extracted features or an error if the request fails.
func (ic *InferenceClient) FeatureExtraction(ctx context.Context, req *FeatureExtractionRequest) (FeatureExtractionResponse, error) {
	if len(req.Inputs) == 0 {
		return nil, errors.New("inputs are required")
	}

	body, err := ic.post(ctx, req.Model, "feature-extraction", req)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))

	featureExtractionResponse := FeatureExtractionResponse{}
	if err := json.Unmarshal(body, &featureExtractionResponse); err != nil {
		return nil, err
	}

	return featureExtractionResponse, nil
}
