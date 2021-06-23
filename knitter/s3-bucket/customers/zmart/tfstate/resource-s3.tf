resource "aws_s3_bucket" "compuzest_tfstate" {
  bucket = "zlifecycle-tfstate-zmart"
  acl    = "private"

  versioning {
    enabled = true
  }
  
  tags = {
    Terraform   = true
  }
}
