#!/bin/bash

####################
# Configuration
####################

# for subcommand './consym.sh build'
export BUILD_HLF_PEER_IMAGE='false'
export BUILD_HLF_ORDERER_IMAGE='false'
export BUILD_HLF_TOOLS_IMAGE='false'
export BUILD_TAPE_IMAGE='false'

# for subcommand './consym.sh clean'
export CLEAN_HLF_PEER_IMAGE='false'
export CLEAN_HLF_ORDERER_IMAGE='false'
export CLEAN_HLF_TOOLS_IMAGE='false'
export CLEAN_TAPE_IMAGE='false'
export CLEAN_DANGLING_IMAGE='false'
export CLEAN_RESULT='true'

export HLF_TYPE='strawman'        # 'vanilla' 'strawman' 'consym'
export HLF_IMAGE_PREFIX='fabric'
export HLF_OLD_TAG='latest'
export HLF_NEW_TAG="${HLF_TYPE}"

export HOSTS=(
  202.45.128.165
  202.45.128.164
  202.45.128.163
)
export INIT_NODE_IP=${HOSTS[0]}

####################
# Path
####################
export CONSYM_PATH="${PWD}"
export TAPE_PATH="${CONSYM_PATH}/tape"
export SHARP_PATH="${CONSYM_PATH}/FabricSharp"
export TEST_PATH="${CONSYM_PATH}/test"
export IMAGE_PATH="${CONSYM_PATH}/image"
export PATH="${TEST_PATH}/bin:${PATH}"
export HLF_PATH="${CONSYM_PATH}/fabric"
