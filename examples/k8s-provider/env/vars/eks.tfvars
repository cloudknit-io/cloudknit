cluster_name = "example-eks"
cluster_version = "1.24"
enable_irsa = true
create_aws_auth_configmap = true
manage_aws_auth_configmap = true

cluster_addons = {
  kube-proxy = {}
  vpc-cni    = {}
  # aws-ebs-csi-driver = {}
  # coredns = {
  #   configuration_values = jsonencode({
  #     computeType = "Fargate"
  #   })
  # }
}
