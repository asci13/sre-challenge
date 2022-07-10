resource "google_project_service" "iam" {
  project = var.project_id
  service = "iam.googleapis.com"

  timeouts {
    create = "30m"
    update = "40m"
  }

  disable_dependent_services = true
  
  lifecycle {
    prevent_destroy = true
  }
}
resource "google_project_service" "containerregistry" {
  project = var.project_id
  service = "containerregistry.googleapis.com"

  timeouts {
    create = "30m"
    update = "40m"
  }

  disable_dependent_services = true
  
  lifecycle {
    prevent_destroy = true
  }
}

resource "google_project_service" "resourcemanager" {
  project = var.project_id
  service = "cloudresourcemanager.googleapis.com"

  timeouts {
    create = "30m"
    update = "40m"
  }

  disable_dependent_services = true
  
  lifecycle {
    prevent_destroy = true
  }
}

resource "google_project_service" "container" {
  project = var.project_id
  service = "container.googleapis.com"

  timeouts {
    create = "30m"
    update = "40m"
  }

  disable_dependent_services = true
  
  lifecycle {
    prevent_destroy = true
  }
}

resource "google_project_service" "compute" {
  project = var.project_id
  service = "compute.googleapis.com"

  timeouts {
    create = "30m"
    update = "40m"
  }

  disable_dependent_services = true
  
  lifecycle {
    prevent_destroy = true
  }
}

resource "google_project_service" "logging" {
  project = var.project_id
  service = "logging.googleapis.com"

  timeouts {
    create = "30m"
    update = "40m"
  }

  disable_dependent_services = true
  
  lifecycle {
    prevent_destroy = true
  }
}

resource "google_project_service" "monitoring" {
  project = var.project_id
  service = "monitoring.googleapis.com"

  timeouts {
    create = "30m"
    update = "40m"
  }

  disable_dependent_services = true
  
  lifecycle {
    prevent_destroy = true
  }
}

resource "google_project_service" "serviceusage" {
  project = var.project_id
  service = "serviceusage.googleapis.com"

  timeouts {
    create = "30m"
    update = "40m"
  }

  disable_dependent_services = true
  
  lifecycle {
    prevent_destroy = true
  }
}