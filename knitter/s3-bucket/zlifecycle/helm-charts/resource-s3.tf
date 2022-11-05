resource "aws_s3_bucket" "zlifecycle_helm_charts" {
  bucket = "zlifecycle-helm-charts"
  acl    = "private"

  versioning {
    enabled = true
  }
  
  tags = {
    Terraform = true
  }
}
