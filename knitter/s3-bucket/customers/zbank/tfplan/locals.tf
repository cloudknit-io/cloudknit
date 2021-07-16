locals {
  aws = {
    region  = var.aws_region
    profile = var.aws_profile
  }
  s3  = {
    tfplan_bucket = "zlifecycle-tfplan-${var.customer_name}"
  }
}
