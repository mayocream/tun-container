#!/bin/sh

DLVDAP=$(which dlv-dap)

exec sudo dlv-dap --only-same-user=false "$@"
