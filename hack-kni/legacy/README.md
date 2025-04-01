# nrtpatch

use this tool to automatically set the deprecated `topologyPolicies` field with a value consistent with the content of the NRT object attributes.
This is effectively a wasteful no-operation which can be needed in corner cases scenarios.

NOTE: you can replace `go run ./nrtpatch.go` with `nrtpatch.sh` in the following lines if you don't have the go toolchain available.

assuming a valid KUBECONFIG
check the computed value:
```
kubectl get noderesourcetopologies.topology.node.k8s.io node19.lab.example.com -o json | go run ./nrtpatch.go
```

automatically patch the relevant object:
```
eval $( kubectl get noderesourcetopologies.topology.node.k8s.io node19.lab.example.com -o json | go run ./nrtpatch.go )
```

undo the changes:
```
kubectl patch noderesourcetopologies.topology.node.k8s.io node19.lab.example.com --type=merge -p {"topologyPolicies":[]}
```
