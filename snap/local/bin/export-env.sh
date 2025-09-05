#!/bin/bash -e

export CONFIG_FILE=$(snapctl get config-file)

VERBOSE=$(snapctl get verbose)
if [[ -n $VERBOSE ]]; then
  export VERBOSE=$VERBOSE
fi

exec "$@"
