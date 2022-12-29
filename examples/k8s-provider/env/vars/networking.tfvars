name            = "example-vpc"
cidr            = "10.16.0.0/16"
azs             = ["us-east-1a", "us-east-1b", "us-east-1c"]
private_subnets = ["10.16.0.0/19", "10.16.64.0/19", "10.16.128.0/19"]
public_subnets  = ["10.16.48.0/20", "10.16.112.0/20", "10.16.176.0/20"]
enable_nat_gateway = false
single_nat_gateway = false
enable_dns_hostnames = true
