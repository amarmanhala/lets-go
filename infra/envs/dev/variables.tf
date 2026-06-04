variable "location" {
  description = "Azure region for dev resources."
  type        = string
  default     = "canadacentral"
}

variable "environment" {
  description = "Environment name used in resource names and tags."
  type        = string
  default     = "dev"
}

variable "resource_group_name" {
  description = "Resource group for the dev AKS stack."
  type        = string
  default     = "rg-lets-go-dev"
}

variable "acr_name" {
  description = "Globally unique Azure Container Registry name. Use lowercase letters and numbers only."
  type        = string
}

variable "aks_name" {
  description = "AKS cluster name."
  type        = string
  default     = "aks-lets-go-dev"
}

variable "node_vm_size" {
  description = "VM size for the low-cost dev AKS node."
  type        = string
  default     = "Standard_B2s"
}

variable "node_count" {
  description = "Node count for the dev AKS system node pool."
  type        = number
  default     = 1
}

variable "app_namespace" {
  description = "Kubernetes namespace for the API and Postgres."
  type        = string
  default     = "lets-go"
}

variable "argocd_namespace" {
  description = "Kubernetes namespace for Argo CD."
  type        = string
  default     = "argocd"
}

variable "git_repository_url" {
  description = "Git repository URL Argo CD watches."
  type        = string
  default     = "https://github.com/amarmanhala/lets-go.git"
}

variable "git_target_revision" {
  description = "Git revision Argo CD watches."
  type        = string
  default     = "main"
}

variable "postgres_user" {
  description = "Postgres username used by the app and Postgres container."
  type        = string
  default     = "admin"
}

variable "postgres_database" {
  description = "Postgres database name used by the app."
  type        = string
  default     = "mydb"
}

variable "postgres_password" {
  description = "Optional Postgres password. If null, Terraform generates one."
  type        = string
  default     = null
  sensitive   = true
}
