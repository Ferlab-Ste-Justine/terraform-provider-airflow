data "airflow_fab_role" "viewer" {
  name       = "Viewer"
}

output "viewer_role" {
  value = data.airflow_fab_role.viewer
}