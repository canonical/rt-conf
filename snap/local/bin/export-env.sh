#!/bin/bash -e

export CONFIG_FILE=$(snapctl get config-file)
export VERBOSE=$(snapctl get verbose)

exec "$@"
