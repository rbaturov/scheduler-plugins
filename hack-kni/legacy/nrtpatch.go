package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-logr/logr"

	topologyv1alpha2 "github.com/k8stopologyawareschedwg/noderesourcetopology-api/pkg/apis/topology/v1alpha2"
	kubeletconfig "k8s.io/kubernetes/pkg/kubelet/apis/config"
	"sigs.k8s.io/scheduler-plugins/pkg/noderesourcetopology/nodeconfig"
)

func toTopologyManagerPolicy(conf nodeconfig.TopologyManager) string {
	if conf.Policy == kubeletconfig.SingleNumaNodeTopologyManagerPolicy && conf.Scope == kubeletconfig.PodTopologyManagerScope {
		return string(topologyv1alpha2.SingleNUMANodePodLevel)
	}
	if conf.Policy == kubeletconfig.SingleNumaNodeTopologyManagerPolicy && conf.Scope == kubeletconfig.ContainerTopologyManagerScope {
		return string(topologyv1alpha2.SingleNUMANodeContainerLevel)
	}
	if conf.Policy == kubeletconfig.BestEffortTopologyManagerPolicy && conf.Scope == kubeletconfig.PodTopologyManagerScope {
		return string(topologyv1alpha2.BestEffortPodLevel)
	}
	if conf.Policy == kubeletconfig.BestEffortTopologyManagerPolicy && conf.Scope == kubeletconfig.ContainerTopologyManagerScope {
		return string(topologyv1alpha2.BestEffortContainerLevel)
	}
	if conf.Policy == kubeletconfig.RestrictedTopologyManagerPolicy && conf.Scope == kubeletconfig.PodTopologyManagerScope {
		return string(topologyv1alpha2.RestrictedPodLevel)
	}
	if conf.Policy == kubeletconfig.RestrictedTopologyManagerPolicy && conf.Scope == kubeletconfig.ContainerTopologyManagerScope {
		return string(topologyv1alpha2.RestrictedContainerLevel)
	}
	return ""
}

func main() {
	var nrt topologyv1alpha2.NodeResourceTopology
	err := json.NewDecoder(os.Stdin).Decode(&nrt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot decode object: %v\n", err)
	}
	tm := nodeconfig.TopologyManagerFromNodeResourceTopology(logr.Discard(), &nrt)
	fmt.Println(fmt.Sprintf(`kubectl patch noderesourcetopologies.topology.node.k8s.io %s --type=merge -p '{"topologyPolicies":["%s"]}'`, nrt.Name, toTopologyManagerPolicy(tm)))
}
