#!/bin/bash

set -eu

if [ "$#" != "3" ]; then
	echo "usage: $0 commit branch-name version"
	exit 1
fi
COMMITHASH="$1"
TARGETBRANCH="$2"
VERSION="$3"
CURRENTBRANCH=$( git rev-parse --abbrev-ref HEAD )

git checkout -q release-$VERSION
git rebase -q origin/release-$VERSION
git checkout -q -b $TARGETBRANCH-$VERSION
git cherry-pick -x $COMMITHASH > /dev/null
git log -n 1 --format=%B | sed s/KNI/KNI\]\[release-$VERSION/ | git commit --quiet --amend -F -
git checkout -q $CURRENTBRANCH
echo "$TARGETBRANCH-$VERSION"
