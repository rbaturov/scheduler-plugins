#!/bin/bash

set -eu

# check jq is available - requires set -e
jq --version > /dev/null

NAME=$( jq -r '.metadata.name' /dev/stdin )
SCOPE=$( jq -r '.attributes[] | select(.name=="topologyManagerScope").value' /dev/stdin )
POLICY=$( jq -r '.attributes[] | select(.name=="topologyManagerPolicy").value' /dev/stdin )

FIX=""
if [[ $POLICY == "single-numa-node" && $SCOPE == "pod" ]]; then
	FIX="SingleNUMANodePodLevel"
elif [[ $POLICY == "single-numa-node" && $SCOPE == "container" ]]; then
	FIX="SingleNUMANodeContainerLevel"
elif [[ $POLICY == "best-effort" && $SCOPE == "pod" ]]; then
	FIX="BestEffortPodLevel"
elif [[ $POLICY == "best-effort" && $SCOPE == "container" ]]; then
	FIX="BestEffortContainerLevel"
elif [[ $POLICY == "restricted" && $SCOPE == "pod" ]]; then
	FIX="RestrictedPodLevel"
elif [[ $POLICY == "restricted" && $SCOPE == "container" ]]; then
	FIX="RestrictedContainerLevel"
elif [[ $POLICY == "none" ]]; then
	# nothing to do
	exit 0
else
	echo "cannot decode JSON input"
	exit 1
fi

echo "kubectl patch noderesourcetopologies.topology.node.k8s.io $NAME --type=merge -p '{\"topologyPolicies\":[\"$FIX\"]}'"
