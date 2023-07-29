variable "gcp_project_id" {
    description = "the gcp project ID being used"
    type = string
}

variable "region" {
    description = "the gcp region being used"
    type = string
}

variable "image_tag" {
    description = "the docker image tag"
    type = string
}

variable "influxdb_token" {
    description = "influxdb token"
    type = string
    sensitive = true
}
