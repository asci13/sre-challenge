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

resource "google_project_service" "project" {
  project = var.gcp-project-name
  service = "containerregistry.googleapis.com"

  timeouts {
    create = "30m"
    update = "40m"
  }

  disable_dependent_services = true
}

resource "google_container_registry" "registry" {
  project  = var.gcp-project-name
  location = "EU"
}