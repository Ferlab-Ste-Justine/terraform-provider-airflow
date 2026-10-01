resource "airflow_fab_role" "test" {
  name       = "troubleshoot_role"
  permissions = [
    {
      action   = "can_read"
      resource = "Website"
    },
    {
      action   = "can_read"
      resource = "DAGs"
    }
  ] 
}

resource "airflow_fab_user" "test" {
  username   = "troubleshoot_user"
  password   = "TempPass123!!"
  email      = "troubleshoot@example.com"
  first_name = "Troubleshoot"
  last_name  = "User"
  roles      = ["User", airflow_fab_role.test.name]
}