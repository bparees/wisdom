BINARY_NAME = wisdom
IMAGE_NAME = quay.io/bparees/wisdom
IMAGE_TAG = latest

# Targets
build:
	go build -o $(BINARY_NAME) ./cmd/wisdom

image:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

clean:
	rm -f $(BINARY_NAME)

.PHONY: build image clean