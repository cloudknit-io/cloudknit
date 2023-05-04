mkdir ~/.aws
cat <<EOT >> ~/.aws/credentials
[default]
aws_access_key_id = ${CUSTOMER_AWS_ACCESS_KEY_ID}
aws_secret_access_key = ${CUSTOMER_AWS_SECRET_ACCESS_KEY}
EOT

if [ ! -z "${CUSTOMER_AWS_SESSION_TOKEN}" ]
cat <<EOT >> ~/.aws/credentials
aws_session_token = ${CUSTOMER_AWS_SESSION_TOKEN}
EOT
fi

cat <<EOT >> ~/.aws/credentials
[compuzest-shared]
aws_access_key_id = ${SHARED_AWS_ACCESS_KEY_ID}
aws_secret_access_key = ${SHARED_AWS_SECRET_ACCESS_KEY}
EOT

if [ ! -z "${SHARED_AWS_SESSION_TOKEN}" ]
cat <<EOT >> ~/.aws/credentials
aws_session_token = ${SHARED_AWS_SESSION_TOKEN}
EOT
fi
