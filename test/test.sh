#!/bin/bash

####################
# Import variables
####################
source ./env.sh

####################
# Import functions
####################
source ${SCRIPT_PATH}/util.sh

####################
# Functions
####################
function prepare_artifact() {
  ${SCRIPT_PATH}/prepare_artifact.sh
}

function deploy_service() {
  mkdir -p ${RESULT_PATH}

  if [[ ${DEPLOY_MODE} == 'swarm' ]]; then
    infoln "Initialize swarm with '${INIT_NODE_IP}' as manager node"
    sudo -E docker swarm init --advertise-addr ${INIT_NODE_IP} &>/dev/null

    infoln "Create overlay network '${NETWORK_NAME}'"
    sudo -E docker network create --driver overlay --attachable ${NETWORK_NAME}

    infoln "Deploy service to swarm"
    sudo -E docker stack deploy -c ${HLF_SWARM_FILE} --resolve-image never ${PROJECT_NAME}
  fi

  if [[ ${DEPLOY_MODE} == 'compose' ]]; then
    infoln "Deploy service to compose"
    sudo -E docker-compose -f ${HLF_COMPOSE_FILE} -p ${PROJECT_NAME} up -d 2>&1
  fi
}

function setup_network() {
  infoln "Creating channel and deploy network"
  local cli_id=''
  while :; do
    infoln "Waiting for cli to start up..."
    cli_id=$(sudo -E docker container ls -f 'name=cli' -q 2>/dev/null)
    if [[ -n ${cli_id} ]]; then
      break
    fi
    sleep ${DELAY}
  done

  sudo -E docker exec ${cli_id} bash -c "./scripts/setup_network.sh"
}

function benchmark() {
  infoln "Using tape to benchmark the network"
  local tape_id=''
  while :; do
    infoln "Waiting for tape to start up..."
    tape_id=$(sudo -E docker container ls -f 'name=tape' -q 2>/dev/null)
    if [[ -n ${tape_id} ]]; then
      break
    fi
    sleep ${DELAY}
  done

  sudo -E docker exec ${tape_id} sh "./scripts/benchmark.sh"
}

function leave() {
  if [[ ${DEPLOY_MODE} == 'swarm' ]]; then
    infoln "Remove services from swarm"
    sudo -E docker stack rm ${PROJECT_NAME}

    infoln "Removing chaincode containers"
    CONTAINER_IDS=$(sudo -E docker ps -a | awk '($2 ~ /dev-.*/) {print $1}')
    if [[ -z "${CONTAINER_IDS}" || "${CONTAINER_IDS}" == " " ]]; then
      infoln "No containers available for deletion"
    else
      sudo -E docker container rm -f ${CONTAINER_IDS}
    fi

    infoln "Removing chaincode images"
    DOCKER_IMAGE_IDS=$(sudo -E docker images | awk '($1 ~ /dev-.*/) {print $3}')
    if [[ -z "${DOCKER_IMAGE_IDS}" || "${DOCKER_IMAGE_IDS}" == " " ]]; then
      infoln "No images available for deletion"
    else
      sudo -E docker image remove -f ${DOCKER_IMAGE_IDS}
    fi

    infoln "Remove overlay network '${NETWORK_NAME}'"
    sudo -E docker network rm ${NETWORK_NAME}

    # infoln "Manager '${INIT_NODE_IP}' leaves swarm"
    # sudo -E docker swarm leave --force
  fi

  if [[ ${DEPLOY_MODE} == 'compose' ]]; then
    infoln "Remove services"
    sudo -E docker-compose -f ${HLF_COMPOSE_FILE} down --volumes --remove-orphans

    infoln "Remove chaincode containers"
    CONTAINER_IDS=$(sudo -E docker ps -a | awk '($2 ~ /dev-.*/) {print $1}')
    if [[ -z "${CONTAINER_IDS}" || "${CONTAINER_IDS}" == " " ]]; then
      infoln "No containers available for deletion"
    else
      sudo -E docker container rm -f ${CONTAINER_IDS}
    fi

    infoln "Remove chaincode images"
    DOCKER_IMAGE_IDS=$(sudo -E docker images | awk '($1 ~ /dev-.*/) {print $3}')
    if [[ -z "${DOCKER_IMAGE_IDS}" || "${DOCKER_IMAGE_IDS}" == " " ]]; then
      infoln "No images available for deletion"
    else
      sudo -E docker image rm -f ${DOCKER_IMAGE_IDS}
    fi
  fi
}


CMD=$1
if [[ $CMD == run ]]; then
  prepare_artifact
  deploy_service
  setup_network
  benchmark
  # sudo docker service logs hlf_org1peer0 &>hlf_org1peer0.log
  # sudo docker service logs hlf_org0orderer &>hlf_org0orderer.log
  leave
elif [[ $CMD == analyze ]]; then
  analyze
fi
