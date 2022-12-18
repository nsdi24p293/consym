#!/bin/bash

function generate_create_channel_tx() {
  which configtxgen &>/dev/null || fatalln "configtxgen tool not found"

  mkdir ${APP_CHANNEL_PATH}

  infoln "Generating channel creation transaction"
  set -x
  configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ${TX_FILE} -channelID ${APP_CHANNEL_NAME}
  res=$?
  { set +x; } 2>/dev/null

  verify_result $res "Generate channel creation transaction" "Failed to generate channel creation transaction"
}

function create_channel() {
  use_peer 1 0

  # Poll in case the raft leader is not set yet
  infoln "Creating channel '${APP_CHANNEL_NAME}'"
  local rc=1
  local counter=0
  while ((rc != 0 && counter < MAX_RETRY)); do
    sleep ${DELAY}

    set -x
    peer channel create -o ${ORDERER_ENDPOINT} -c ${APP_CHANNEL_NAME} --ordererTLSHostnameOverride ${ORDERER_HOSTNAME} -f ${TX_FILE} --outputBlock ${BLOCK_FILE} --tls --cafile ${ORG0_CA}
    res=$?
    { set +x; } 2>/dev/null

    ((rc = $res))
    ((counter = counter + 1))
  done

  verify_result $rc "Create channel '${APP_CHANNEL_NAME}'" "Failed to create channel '${APP_CHANNEL_NAME}' after ${MAX_RETRY} retries"
}

# join_channel org peer
function join_channel() {
  local org=$1
  local peer=$2
  local peer_fullname="org${org}peer${peer}"

  use_peer ${org} ${peer}

  infoln "Joining peer '${peer_fullname}' to channel '${APP_CHANNEL_NAME}'"
  local rc=1
  local counter=0
  ## Sometimes Join takes time, hence retry
  while ((rc != 0 && counter < MAX_RETRY)); do
    sleep ${DELAY}

    set -x
    peer channel join -b ${BLOCK_FILE}
    rc=$?
    { set +x; } 2>/dev/null

    ((rc = $res))
    ((counter = counter + 1))
  done

  verify_result $rc "Join peer '${peer_fullname}' to channel '${APP_CHANNEL_NAME}'" "Failed to join peer '${peer_fullname}' to channel '${APP_CHANNEL_NAME}'"
}

# set_anchor_peer org peer
function set_anchor_peer() {
  local org=$1
  local peer=$2

  ${SCRIPT_PATH}/set_anchor_peer.sh ${org} ${peer} ${APP_CHANNEL_NAME}
}

generate_create_channel_tx

## Create channel
infoln "Creating channel ${APP_CHANNEL_NAME}"
create_channel

if [[ ${DEPLOY_SCHEME} == 1p1o ]]; then
  ## Join all the peers to the channel
  infoln "Joining org1 peer to the channel..."
  join_channel 1 0

  ## Set the anchor peers for each org in the channel
  infoln "Setting anchor peer for org1..."
  set_anchor_peer 1 0

elif [[ ${DEPLOY_SCHEME} == 4p1o ]]; then
  ## Join all the peers to the channel
  infoln "Joining org1 peer to the channel..."
  join_channel 1 0
  join_channel 1 1
  infoln "Joining org2 peer to the channel..."
  join_channel 2 0
  join_channel 2 1

  ## Set the anchor peers for each org in the channel
  infoln "Setting anchor peer for org1..."
  set_anchor_peer 1 0
  infoln "Setting anchor peer for org2..."
  set_anchor_peer 2 0
fi

successln "Channel '${APP_CHANNEL_NAME}' joined"
