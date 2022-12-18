#!/bin/sh

# Import variables
source ./env.sh

function run_e2e() {
  # Creating accounts
  set -x
  tape --config ${TAPE_CONFIG_FILE_PUT} >${TAPE_STDOUT_FILE_PUT} 2>${TAPE_STDERR_FILE_PUT}
  cp ACCOUNTS.txt TRANSACTIONS.txt result

  # Transfering money
  tape --config ${TAPE_CONFIG_FILE_CONFLICT} >${TAPE_STDOUT_FILE_CONFLICT} 2>${TAPE_STDERR_FILE_CONFLICT}
  { set +x; } 2>/dev/null
}

run_e2e
