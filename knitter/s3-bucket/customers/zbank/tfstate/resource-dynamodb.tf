resource "aws_dynamodb_table" "tflock" {
  name           = local.ddb.tflock_table
  hash_key       = "LockID"
  read_capacity  = local.ddb.read_capacity
  write_capacity = local.ddb.write_capacity

  tags = {
    Terraform   = true
  }

  attribute {
    name = "LockID"
    type = "S"
  }
}

