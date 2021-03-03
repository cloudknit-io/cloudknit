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

env_component_path=$1
config_name=$2
module_source=$3
module_source_path=$4
variables_file_path=$5

set -eo pipefail

if [ -n "${module_source_path}" ]; then
    full_module_source="${module_source}//${module_source_path}"
else
    full_module_source="${module_source}"
fi

cat > $env_component_path/module.tf << EOL
module "${config_name}" {
  source = "${full_module_source}"
EOL

cat $env_component_path/vars/$variables_file_path >> $env_component_path/module.tf

echo "}" >> $env_component_path/module.tf

cat $env_component_path/module.tf
