# Countries API

A minimal Go backend with no web framework. It exposes one endpoint using the standard library and stores country data in PostgreSQL.

The app reads these database environment variables:

```sh
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
DB_SSLMODE=disable
PORT=8080
```

## Run

```sh
DB_HOST=localhost DB_USER=admin DB_PASSWORD=password DB_NAME=mydb go run .
```

## Docker

```sh
docker build -t lets-go:v1 .
docker network create lets-go-net
docker network connect lets-go-net postgres
docker run --rm \
  --name lets-go-api \
  --network lets-go-net \
  -p 8080:8080 \
  -e DB_HOST=postgres \
  -e DB_USER=admin \
  -e DB_PASSWORD=password \
  -e DB_NAME=mydb \
  lets-go:v1
```

## Azure AKS Deployment

This repo includes Terraform for Azure infrastructure and a Helm chart for Argo CD GitOps deployment.

Terraform creates:

- Azure Container Registry
- AKS with Azure RBAC enabled
- Log Analytics
- Argo CD installed by Helm
- Argo CD application pointing to `deploy/helm/lets-go`
- Kubernetes secret for Postgres credentials

The Helm chart deploys:

- Go API as a Kubernetes Deployment and LoadBalancer Service
- Postgres 17 as a StatefulSet with a persistent volume

### 1. Bootstrap Terraform Remote State

Create a unique storage account name, then run:

```sh
cd infra/bootstrap
cp terraform.tfvars.example terraform.tfvars
terraform init
terraform apply
```

Use the outputs to configure the dev backend:

```sh
cd ../envs/dev
cp backend.hcl.example backend.hcl
cp terraform.tfvars.example terraform.tfvars
```

Update:

- `backend.hcl`
- `terraform.tfvars`

The `acr_name` value must be globally unique and must match the GitHub variable `AZURE_ACR_NAME`.

### 2. Create Azure OIDC Credentials for GitHub Actions

Create an Azure app registration/service principal with federated credentials for:

```text
repo:amarmanhala/lets-go:ref:refs/heads/main
repo:amarmanhala/lets-go:pull_request
```

Add these GitHub Actions secrets:

```text
AZURE_CLIENT_ID
AZURE_TENANT_ID
AZURE_SUBSCRIPTION_ID
```

Add these GitHub Actions variables:

```text
AZURE_ACR_NAME
AZURE_TFSTATE_RESOURCE_GROUP
AZURE_TFSTATE_STORAGE_ACCOUNT
AZURE_TFSTATE_CONTAINER
```

### 3. Apply Infrastructure

```sh
cd infra/envs/dev
terraform init -backend-config=backend.hcl
terraform apply
```

### 4. CI/CD Flow

On pull requests, GitHub Actions runs:

- Go formatting, vet, race tests, and vulnerability scan
- Terraform validate and plan
- Docker build and Trivy image scan

On push to `main`, after all checks pass:

- The Docker image is pushed to Azure Container Registry with the Git SHA tag
- `deploy/helm/lets-go/values.yaml` is updated with that image tag
- Argo CD syncs the Helm chart into AKS

## Endpoint

```sh
curl http://localhost:8080/api/countries
```

Returns 50 mocked country records:

```json
[
  {
    "name": "Canada",
    "population": 40097761,
    "flag_colors": ["red", "white"],
    "language": "English"
  }
]
```
