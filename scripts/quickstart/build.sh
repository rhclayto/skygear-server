#!/bin/sh

: ${DOCKER_COMPOSE_VERSION:=1.6.2}
: ${SKYGEAR_VERSION:=latest}

DOCKER_ENGINE_PACKAGE=docker-engine=1.10.2-0~precise

# Update apt packages
echo "deb https://apt.dockerproject.org/repo ubuntu-precise main" > /etc/apt/sources.list.d/docker.list
apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D
apt-get update
apt-get upgrade -y

locale-gen UTF-8

# Install linux-image-extra so that Docker will use aufs
apt-get -y install linux-image-extra-$(uname -r)

# Install Docker Engine
apt-get install -y --no-install-recommends $DOCKER_ENGINE_PACKAGE

# Pull images used by Skygear
docker pull nginx:1.9
docker pull mdillon/postgis:9.4
docker pull redis:3.0
docker pull skygeario/skygear-server:$SKYGEAR_VERSION

# Install Docker Compose
curl -L https://github.com/docker/compose/releases/download/$DOCKER_COMPOSE_VERSION/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# Install CloudFormation Helper Scripts
apt-get install -y python-pip
pip install pyopenssl ndg-httpsclient pyasn1
pip install https://s3.amazonaws.com/cloudformation-examples/aws-cfn-bootstrap-latest.tar.gz

# Install jinja2-cli to modify config files
pip install jinja2-cli

# Move other files into place
if [ -f /tmp/kickstart.sh]; then
  mv /tmp/kickstart.sh /usr/local/bin/kickstart.sh
  chmod +x /usr/local/bin/kickstart.sh
fi
