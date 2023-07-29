resource "google_secret_manager_secret" "influxdb_token" {
  secret_id = "influxdb_token"

  labels = {
    label = "influxdb"
  }

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "influxdb_token_version" {
  secret = google_secret_manager_secret.influxdb_token.id

  secret_data = var.influxdb_token
}
