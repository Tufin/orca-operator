#!/bin/bash

PREFIX="${1}"
TAG="${2}"
COMMIT_MSG=CHECK_COMMIT_MSG="$(git log -1 --pretty=%B)"

if [[ "${COMMIT_MSG}" == *"[${PREFIX} ${TAG}]"* ]]; then
    exit 0
fi

exit 1
