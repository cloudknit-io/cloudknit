resource "aws_s3_bucket" "compuzest_tfstate" {
  bucket = "compuzest-tfstate"
  acl    = "private"

  versioning {
    enabled = true
  }
  
  tags = {
    Terraform   = true
  }
}
