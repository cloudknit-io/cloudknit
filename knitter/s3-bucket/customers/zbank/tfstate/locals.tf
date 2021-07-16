locals {
  aws = {
    region  = var.aws_region
    profile = var.aws_profile
  }
  s3  = {
    tfstate_bucket = "zlifecycle-tfstate-${var.customer_name}"
  }
  ddb = {
    tflock_table   = "zlifecycle-tflock-${var.customer_name}"
    read_capacity  = 20
    write_capacity = 20
  }
}
