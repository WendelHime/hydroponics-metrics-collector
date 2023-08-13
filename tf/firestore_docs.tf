resource "google_firestore_document" "user_devices_document" {
  project = var.gcp_project_id
  collection  = google_project_service.firestore.name
  document_id = "my-doc-id"
  fields      = "{\"user_id\":{\"stringValue\":\"avalue\"},\"devices\":{\"arrayValue\":{\"stringValue\":\"avalue\"}}}"
  depends_on = [google_project_service.firestore]
}
