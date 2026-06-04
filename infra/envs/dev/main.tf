locals {
  app_name        = "lets-go"
  postgres_secret = "lets-go-postgres"

  tags = {
    app         = local.app_name
    environment = var.environment
    managed_by  = "terraform"
  }
}

resource "azurerm_resource_group" "main" {
  name     = var.resource_group_name
  location = var.location
  tags     = local.tags
}

resource "azurerm_log_analytics_workspace" "main" {
  name                = "law-${local.app_name}-${var.environment}"
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
  tags                = local.tags
}

resource "azurerm_container_registry" "main" {
  name                = var.acr_name
  resource_group_name = azurerm_resource_group.main.name
  location            = azurerm_resource_group.main.location
  sku                 = "Basic"
  admin_enabled       = false
  tags                = local.tags
}

resource "azurerm_kubernetes_cluster" "main" {
  name                = var.aks_name
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
  dns_prefix          = "lets-go-${var.environment}"
  sku_tier            = "Free"
  tags                = local.tags

  default_node_pool {
    name       = "system"
    node_count = var.node_count
    vm_size    = var.node_vm_size
  }

  identity {
    type = "SystemAssigned"
  }

  azure_active_directory_role_based_access_control {
    tenant_id          = data.azurerm_client_config.current.tenant_id
    azure_rbac_enabled = true
  }

  oms_agent {
    log_analytics_workspace_id = azurerm_log_analytics_workspace.main.id
  }

  network_profile {
    network_plugin    = "kubenet"
    load_balancer_sku = "standard"
    outbound_type     = "loadBalancer"
  }
}

resource "azurerm_role_assignment" "aks_acr_pull" {
  scope                = azurerm_container_registry.main.id
  role_definition_name = "AcrPull"
  principal_id         = azurerm_kubernetes_cluster.main.kubelet_identity[0].object_id
}

resource "kubernetes_namespace_v1" "app" {
  metadata {
    name = var.app_namespace
    labels = {
      "app.kubernetes.io/name"       = local.app_name
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

resource "kubernetes_namespace_v1" "argocd" {
  metadata {
    name = var.argocd_namespace
    labels = {
      "app.kubernetes.io/name"       = "argocd"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

resource "random_password" "postgres" {
  length  = 24
  special = true
}

resource "kubernetes_secret_v1" "postgres" {
  metadata {
    name      = local.postgres_secret
    namespace = kubernetes_namespace_v1.app.metadata[0].name
  }

  data = {
    username = var.postgres_user
    password = coalesce(var.postgres_password, random_password.postgres.result)
    database = var.postgres_database
  }

  type = "Opaque"
}

resource "helm_release" "argocd" {
  name       = "argocd"
  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-cd"
  version    = "7.7.16"
  namespace  = kubernetes_namespace_v1.argocd.metadata[0].name

  values = [
    yamlencode({
      server = {
        service = {
          type = "LoadBalancer"
        }
      }
    })
  ]
}

resource "helm_release" "argocd_application" {
  depends_on = [helm_release.argocd]

  name      = "lets-go-argocd-application"
  chart     = "${path.module}/charts/argocd-application"
  namespace = kubernetes_namespace_v1.argocd.metadata[0].name

  values = [
    yamlencode({
      application = {
        name            = local.app_name
        namespace       = kubernetes_namespace_v1.argocd.metadata[0].name
        project         = "default"
        repositoryUrl   = var.git_repository_url
        targetRevision  = var.git_target_revision
        path            = "deploy/helm/lets-go"
        targetNamespace = kubernetes_namespace_v1.app.metadata[0].name
      }
    })
  ]
}
