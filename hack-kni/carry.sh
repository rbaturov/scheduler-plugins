#!/bin/bash

set -e

UPSTREAM="upstream"

PR="$1"
BRANCH="$2"

if [ -z "$PR" ]; then
	echo "usage: $0 PR [branch]" 1>&2
	exit 1
fi
if [ -z "$BRANCH" ]; then
	BRANCH="pr$PR"
fi

echo "fetching PR $PR into branch $BRANCH"
git fetch upstream pull/$PR/head:tmp-pr-$PR

git co master
git co -B $BRANCH
git merge -m "[CARRY] merge $PR" tmp-pr-$PR
git branch -D tmp-pr-$PR
