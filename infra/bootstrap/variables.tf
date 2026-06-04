variable "location" {
  description = "Azure region for Terraform state resources."
  type        = string
  default     = "canadacentral"
}

variable "resource_group_name" {
  description = "Resource group that stores Terraform remote state."
  type        = string
  default     = "rg-lets-go-tfstate"
}

variable "storage_account_name" {
  description = "Globally unique Azure Storage account name for Terraform state. Use lowercase letters and numbers only."
  type        = string
}

variable "container_name" {
  description = "Blob container name for Terraform state."
  type        = string
  default     = "tfstate"
}
