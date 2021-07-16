variable "aws_region" {
  description = "Region for AWS access"
  type        = string
  default     = "us-east-1"
}

variable "aws_profile" {
  description = "AWS profile"
  type        = string
  default     = null
}

variable "customer_name" {
  description = "Name of the customer (company)"
  type        = string
}
