resource "aws_dynamodb_table" "tflock" {
  name           = "zlifecycle-tflock-zmart"
  hash_key       = "LockID"
  read_capacity  = 20
  write_capacity = 20

  tags = {
    Terraform   = true
  }

  attribute {
    name = "LockID"
    type = "S"
  }
}

