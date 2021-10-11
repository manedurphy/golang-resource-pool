#! /bin/bash
STRATEGY=$1

if [[ -z "$STRATEGY" ]]; then
	STRATEGY="local"
fi

GODEBUG=gctrace=1 ./server/server --strategy=$STRATEGY
