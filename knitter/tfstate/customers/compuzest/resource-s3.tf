resource "aws_s3_bucket" "compuzest_tfstate" {
  bucket = "compuzest-zlifecycle-tfstate"
  acl    = "private"

  versioning {
    enabled = true
  }
  
  tags = {
    Terraform   = true
  }
}
