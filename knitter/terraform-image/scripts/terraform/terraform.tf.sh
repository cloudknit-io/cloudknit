# Copyright (C) 2020 CompuZest, Inc. - All Rights Reserved
#
# Unauthorized copying of this file, via any medium, is strictly prohibited
# Proprietary and confidential
#
# NOTICE: All information contained herein is, and remains the property of
# CompuZest, Inc. The intellectual and technical concepts contained herein are
# proprietary to CompuZest, Inc. and are protected by trade secret or copyright
# law. Dissemination of this information or reproduction of this material is
# strictly forbidden unless prior written permission is obtained from CompuZest, Inc.

set -eo pipefail

env_component_path=$1
team_name=$2
env_name=$3
config_name=$4

cat > $env_component_path/terraform.tf << EOL
terraform {
  required_version = "= 0.13.2"

  backend "s3" {
    profile                 = "compuzest-shared"

    bucket                  = "compuzest-zlifecycle-tfstate"
    key                     = "${team_name}/${env_name}/${config_name}/terraform.tfstate"
    region                  = "us-east-1"

    dynamodb_table          = "compuzest-zlifecycle-tflock"
    encrypt                 = true
  }
}
EOL

cat $env_component_path/terraform.tf
