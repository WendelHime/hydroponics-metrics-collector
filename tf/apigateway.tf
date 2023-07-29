locals {
  api_config_id_prefix     = "hydroponics-project"
  api_id                   = "hydroponics-project"
  gateway_id               = "metrics-collector"
  display_name             = "metrics-collector"
}

resource "google_api_gateway_api" "api_gw" {
  provider     = google-beta
  api_id       = local.api_id
  project      = var.gcp_project_id
  display_name = local.display_name
}

resource "google_api_gateway_api_config" "api_cfg" {
  provider             = google-beta
  api                  = google_api_gateway_api.api_gw.api_id
  api_config_id_prefix = local.api_config_id_prefix
  project              = var.gcp_project_id
  display_name         = local.display_name

  openapi_documents {
    document {
      path     = "openapi.yaml"
      contents = filebase64("../internal/api/openapi.spec.yaml")
    }
  }

  gateway_config {
    backend_config {
      google_service_account = "api-gateway@hydroponics-392400.iam.gserviceaccount.com"
    }
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "google_api_gateway_gateway" "gw" {
  provider = google-beta
  region   = var.region
  project  = var.gcp_project_id


  api_config   = google_api_gateway_api_config.api_cfg.id

  gateway_id   = local.gateway_id
  display_name = local.display_name

  depends_on   = [google_api_gateway_api_config.api_cfg]
}
