provider "google" {
  project = var.gcp_project_id
}

resource "google_cloud_run_v2_service" "metrics_collector_api" {
  name     = "metrics-collector"
  location = var.region
  ingress = "INGRESS_TRAFFIC_ALL"

  template {
    service_account = google_service_account.hydroponics_metrics_collector_sa.email
    containers {
      name = "api"
      image = "${var.region}-docker.pkg.dev/${var.gcp_project_id}/hydroponics-repository/hydroponics-metrics-collector:${var.image_tag}"
      env {
        name = "DATABASE"
        value = "hydroponics"
      }
      env {
        name  = "INFLUXDB_HOST"
        value = "https://us-east-1-1.aws.cloud2.influxdata.com"
      }
      env {
        name  = "INFLUXDB_TOKEN"
        value_source {
          secret_key_ref {
            secret = google_secret_manager_secret.influxdb_token.secret_id
            version = "latest"
          }
        }
      }
    }
  }

  traffic {
    percent = 100
    type = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
  }

  depends_on = [google_secret_manager_secret_version.influxdb_token_version, google_secret_manager_secret_iam_binding.binding]
}
