#!/usr/bin/env bash
VERSION_TAG=$(git tag -l | grep "v" | cut -c2-);
if [[ ${#VERSION_TAG} -ne 0 ]]; then
  echo "$VERSION_TAG";
else
  HASH=$(git rev-parse --short HEAD);
  printf "git-%s" "$HASH"
fi;
