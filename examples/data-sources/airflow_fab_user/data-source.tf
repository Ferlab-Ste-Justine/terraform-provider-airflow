data "airflow_fab_user" "admin" {
  username       = "admin"
}

output "admin_user" {
  value = data.airflow_fab_user.admin
}