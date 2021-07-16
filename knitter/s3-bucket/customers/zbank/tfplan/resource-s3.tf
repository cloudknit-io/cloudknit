resource "aws_s3_bucket" "compuzest_tfplan" {
  bucket = local.s3.tfplan_bucket
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

resource "aws_s3_bucket_public_access_block" "compuzest_tfplan" {
  bucket = aws_s3_bucket.compuzest_tfplan.bucket

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}
