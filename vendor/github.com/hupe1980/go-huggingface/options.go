package gohuggingface

type Options struct {
	// (Default: true). There is a cache layer on the inference API to speedup
	// requests we have already seen. Most models can use those results as is
	// as models are deterministic (meaning the results will be the same anyway).
	// However if you use a non deterministic model, you can set this parameter
	// to prevent the caching mechanism from being used resulting in a real new query.
	UseCache *bool `json:"use_cache,omitempty"`

	// (Default: false) If the model is not ready, wait for it instead of receiving 503.
	// It limits the number of requests required to get your inference done. It is advised
	// to only set this flag to true after receiving a 503 error as it will limit hanging
	// in your application to known places.
	WaitForModel *bool `json:"wait_for_model,omitempty"`
}
