#!/bin/bash

# Import variables and functions
# because these scripts are run in a separate docker container
source ./env.sh
source ./scripts/util.sh

function create_channel() {
  if [[ ! -d "./organization" ]]; then
    fatalln "Organization crypto materials not found"
  fi

  ./scripts/create_channel.sh
}

function deploy_chaincode() {
  ./scripts/deploy_chaincode.sh
}

export FABRIC_CFG_PATH="${CONFIG_PATH}"

create_channel

deploy_chaincode
