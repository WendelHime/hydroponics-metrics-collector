resource "google_service_account" "hydroponics_metrics_collector_sa" {
  account_id   = "cloudrun-sa"
  display_name = "A service account for hydroponics-metrics-collector cloud run instance"
}

resource "google_secret_manager_secret_iam_binding" "binding" {
  project = var.gcp_project_id
  secret_id = google_secret_manager_secret.influxdb_token.secret_id
  role = "roles/secretmanager.secretAccessor"

  members = [
    "serviceAccount:${google_service_account.hydroponics_metrics_collector_sa.email}",
  ]

  depends_on = [google_service_account.hydroponics_metrics_collector_sa]
}
