#!/bin/bash

# SSM Agent config
yum install -y https://s3.us-east-2.amazonaws.com/amazon-ssm-us-east-2/latest/linux_amd64/amazon-ssm-agent.rpm

systemctl enable amazon-ssm-agent
systemctl start amazon-ssm-agent

# nginx config
amazon-linux-extras install -y nginx1

cat << EOF > /etc/nginx/conf.d/default.conf
server {
    listen 80 default_server;
    listen [::]:80 default_server;
    server_name _;
    return 301 ${redirect_url};
}
EOF

systemctl enable nginx
systemctl start nginx
