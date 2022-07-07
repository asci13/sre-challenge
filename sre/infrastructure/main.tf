terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.27.0"
    }
  }
}


resource "google_service_account" "service_account" {
  account_id   = "asci13srechallengeciaccount"
  display_name = "Challenge CI Account"
  project = var.project_id
}

resource "google_container_registry" "registry" {
  project  = var.project_id
  location = "EU"
}

resource "google_project_iam_member" "registry_storage_admin" {
  project = var.project_id
  role    = "roles/storage.admin"
  member  = "serviceAccount:${google_service_account.service_account.email}"
}

resource "google_service_account_key" "registry_key" {
  service_account_id = google_service_account.service_account.name
  public_key_type    = "TYPE_X509_PEM_FILE"
}