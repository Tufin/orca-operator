#!/bin/bash

OPERATOR_IMAGE="registry\.connect\.redhat\.com\/tufin\/orca-operator"
TARGET_FILE="${1}"
NEW_TAG="${2}"

sed -i.bak "s/image: ${OPERATOR_IMAGE}.*/image: ${OPERATOR_IMAGE}:${NEW_TAG}/" ${TARGET_FILE}
