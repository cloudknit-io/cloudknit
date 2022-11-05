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

mkdir ~/.aws
cat <<EOT >> ~/.aws/credentials
[default]
aws_access_key_id = ${CUSTOMER_AWS_ACCESS_KEY_ID}
aws_secret_access_key = ${CUSTOMER_AWS_SECRET_ACCESS_KEY}
[compuzest-shared]
aws_access_key_id = ${SHARED_AWS_ACCESS_KEY_ID}
aws_secret_access_key = ${SHARED_AWS_SECRET_ACCESS_KEY}
EOT
