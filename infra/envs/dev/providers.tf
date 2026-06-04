terraform {
  required_version = ">= 1.7.0"

  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 4.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.16"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.33"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.6"
    }
  }
}

provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

provider "kubernetes" {
  host                   = azurerm_kubernetes_cluster.main.kube_admin_config[0].host
  client_certificate     = base64decode(azurerm_kubernetes_cluster.main.kube_admin_config[0].client_certificate)
  client_key             = base64decode(azurerm_kubernetes_cluster.main.kube_admin_config[0].client_key)
  cluster_ca_certificate = base64decode(azurerm_kubernetes_cluster.main.kube_admin_config[0].cluster_ca_certificate)
}

provider "helm" {
  kubernetes {
    host                   = azurerm_kubernetes_cluster.main.kube_admin_config[0].host
    client_certificate     = base64decode(azurerm_kubernetes_cluster.main.kube_admin_config[0].client_certificate)
    client_key             = base64decode(azurerm_kubernetes_cluster.main.kube_admin_config[0].client_key)
    cluster_ca_certificate = base64decode(azurerm_kubernetes_cluster.main.kube_admin_config[0].cluster_ca_certificate)
  }
}
