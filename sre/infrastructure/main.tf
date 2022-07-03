terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.27.0"
    }
  }
}

provider "google" {

  project = var.gcp-project-name
  region  = "europe-north1"
  zone    = "europe-north1-c"
}

variable gcp-project-name { default="asci13-sre-challenge" }

module "project-services" {
  source  = "terraform-google-modules/project-factory/google//modules/project_services"
  enable_apis = true

  project_id                  = var.gcp-project-name

  activate_apis = [
    "containerregistry.googleapis.com",
    "iam.googleapis.com",
    "cloudresourcemanager.googleapis.com"
  ]
}

resource "google_service_account" "service_account" {
  account_id   = "asci13srechallengeciaccount"
  display_name = "Challenge CI Account"
  project = var.gcp-project-name
}

resource "google_container_registry" "registry" {
  project  = var.gcp-project-name
  location = "EU"
}

resource "google_project_iam_member" "registry_storage_admin" {
  project = var.gcp-project-name
  role    = "roles/storage.admin"
  member  = "serviceAccount:${google_service_account.service_account.email}"
}

resource "google_service_account_key" "registry_key" {
  service_account_id = google_service_account.service_account.name
  public_key_type    = "TYPE_X509_PEM_FILE"
}