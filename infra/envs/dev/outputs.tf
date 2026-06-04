output "resource_group_name" {
  value = azurerm_resource_group.main.name
}

output "acr_name" {
  value = azurerm_container_registry.main.name
}

output "acr_login_server" {
  value = azurerm_container_registry.main.login_server
}

output "aks_name" {
  value = azurerm_kubernetes_cluster.main.name
}

output "app_namespace" {
  value = kubernetes_namespace_v1.app.metadata[0].name
}

output "argocd_namespace" {
  value = kubernetes_namespace_v1.argocd.metadata[0].name
}

output "postgres_secret_name" {
  value = kubernetes_secret_v1.postgres.metadata[0].name
}
