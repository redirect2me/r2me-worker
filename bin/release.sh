#!/usr/bin/env bash
#
# run locally
#

set -o errexit
set -o pipefail
set -o nounset

echo "INFO: release starting at $(date -u +%Y-%m-%dT%H:%M:%SZ)"

TAG=${1:-AUTO}

if [ "$TAG" == "AUTO" ]; then
  # get the latest tag
  echo "INFO: Calculating next tag"
  TAG=$(git tag | sort --numeric-sort --reverse | head -n 1)
  if [ -z "$TAG" ]; then
	echo "INFO: No tags found, using 0.1.0"
	TAG="0.1.0"
  else
	# increment the patch version
	echo "INFO: Latest tag is $TAG, incrementing patch version"
	IFS='.' read -r -a PARTS <<< "$TAG"
	PARTS[2]=$((PARTS[2] + 1))
	TAG="${PARTS[0]}.${PARTS[1]}.${PARTS[2]}"
  fi
fi

echo "INFO: Releasing version $TAG"

git tag ${TAG}
git push --tags

echo "INFO: release complete at $(date -u +%Y-%m-%dT%H:%M:%SZ)"
