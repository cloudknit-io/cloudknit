resource "aws_s3_bucket" "compuzest_tfstate" {
  bucket = local.s3.tfstate_bucket
  acl    = "private"

  versioning {
    enabled = true
  }

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }

  tags = {
    Terraform   = true
  }
}

resource "aws_s3_bucket_public_access_block" "compuzest_tfstate" {
  bucket = aws_s3_bucket.compuzest_tfstate.bucket

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}
