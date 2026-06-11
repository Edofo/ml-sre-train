IMAGE ?= ml-sre-train:dev
CLUSTER ?= ml-sre-train
NS ?= recommender-dev
KUSTOMIZE_DIR ?= infra/k8s/overlays/dev
TF_DIR ?= infra/terraform/envs/dev

.DEFAULT_GOAL := help
.PHONY: help test image cluster-up load deploy up port-forward down tf-validate fmt

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'

test: ## Run Go tests
	cd app && go test ./... -race -cover

image: ## Build the container image
	docker build -t $(IMAGE) -f app/deploy/Dockerfile app

cluster-up: ## Create the kind cluster if it does not exist
	@kind get clusters | grep -qx $(CLUSTER) || kind create cluster --name $(CLUSTER)

load: image ## Build then load the image into the kind cluster
	kind load docker-image $(IMAGE) --name $(CLUSTER)

deploy: ## Apply the Kustomize dev overlay
	kubectl apply -k $(KUSTOMIZE_DIR)

up: ## Zero-to-running: cluster + image + load + deploy
	$(MAKE) cluster-up
	$(MAKE) load
	$(MAKE) deploy

monitoring-up: ## Install the monitoring stack
	helmfile -f observability/helmfile.yaml sync

monitoring-down: ## Uninstall the monitoring stack
	helmfile -f observability/helmfile.yaml destroy

port-forward: ## Forward the service to localhost:8080
	kubectl port-forward -n $(NS) svc/recommender 8080:8080

down: ## Delete the kind cluster
	kind delete cluster --name $(CLUSTER)

tf-validate: ## Init (no backend) and validate Terraform
	terraform -chdir=$(TF_DIR) init -backend=false
	terraform -chdir=$(TF_DIR) validate

fmt: ## Format Go and Terraform
	cd app && gofmt -w .
	terraform -chdir=infra/terraform fmt -recursive
