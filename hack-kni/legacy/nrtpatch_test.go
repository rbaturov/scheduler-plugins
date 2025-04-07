package main

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"reflect"
	"strings"
	"testing"

	"github.com/go-logr/logr"
	topologyv1alpha2 "github.com/k8stopologyawareschedwg/noderesourcetopology-api/pkg/apis/topology/v1alpha2"
	"sigs.k8s.io/scheduler-plugins/pkg/noderesourcetopology/nodeconfig"
)

type tcase struct {
	attrs topologyv1alpha2.AttributeList
	exp   nodeconfig.TopologyManager
}

func TestNRTPatch(t *testing.T) {

	argvs := [][]string{
		[]string{"go", "run", "./nrtpatch.go", "-R"},
		[]string{"bash", "-x", "./nrtpatch.sh", "-R"},
	}
	tcases := []tcase{
		{
			attrs: topologyv1alpha2.AttributeList{
				{
					Name:  "topologyManagerScope",
					Value: "container",
				},
				{
					Name:  "topologyManagerPolicy",
					Value: "single-numa-node",
				},
			},
			exp: nodeconfig.TopologyManager{
				Scope:  "container",
				Policy: "single-numa-node",
			},
		},
		{
			attrs: topologyv1alpha2.AttributeList{
				{
					Name:  "topologyManagerScope",
					Value: "container",
				},
				{
					Name:  "topologyManagerPolicy",
					Value: "restricted",
				},
			},
			exp: nodeconfig.TopologyManager{
				Scope:  "container",
				Policy: "restricted",
			},
		},
		{
			attrs: topologyv1alpha2.AttributeList{
				{
					Name:  "topologyManagerScope",
					Value: "container",
				},
				{
					Name:  "topologyManagerPolicy",
					Value: "best-effort",
				},
			},
			exp: nodeconfig.TopologyManager{
				Scope:  "container",
				Policy: "best-effort",
			},
		},
		{
			attrs: topologyv1alpha2.AttributeList{
				{
					Name:  "topologyManagerScope",
					Value: "pod",
				},
				{
					Name:  "topologyManagerPolicy",
					Value: "single-numa-node",
				},
			},
			exp: nodeconfig.TopologyManager{
				Scope:  "pod",
				Policy: "single-numa-node",
			},
		},
		{
			attrs: topologyv1alpha2.AttributeList{
				{
					Name:  "topologyManagerScope",
					Value: "pod",
				},
				{
					Name:  "topologyManagerPolicy",
					Value: "restricted",
				},
			},
			exp: nodeconfig.TopologyManager{
				Scope:  "pod",
				Policy: "restricted",
			},
		},
		{
			attrs: topologyv1alpha2.AttributeList{
				{
					Name:  "topologyManagerScope",
					Value: "pod",
				},
				{
					Name:  "topologyManagerPolicy",
					Value: "best-effort",
				},
			},
			exp: nodeconfig.TopologyManager{
				Scope:  "pod",
				Policy: "best-effort",
			},
		},
	}

	for _, argv := range argvs {
		for _, tcase := range tcases {
			t.Run(argv[0]+"_"+tcase.exp.String(), func(t *testing.T) {
				var nrt topologyv1alpha2.NodeResourceTopology
				nrt.Name = "test"
				nrt.Attributes = append(nrt.Attributes, tcase.attrs...)
				var inBuf bytes.Buffer
				if err := json.NewEncoder(&inBuf).Encode(&nrt); err != nil {
					t.Fatalf("encode failed: %v", err)
				}
				jsonText := inBuf.String()

				var outBuf bytes.Buffer
				var errBuf bytes.Buffer
				cmd := exec.Command(argv[0], argv[1:]...)
				cmd.Stdin = strings.NewReader(jsonText)
				cmd.Stdout = &outBuf
				cmd.Stderr = &errBuf
				if err := cmd.Run(); err != nil {
					t.Fatalf("run failed: %v stderr=%s (text=%s)", err, errBuf.String(), jsonText)
				}

				nrt2 := nrt.DeepCopy()
				nrt2.TopologyPolicies = []string{outBuf.String()}
				got := nodeconfig.TopologyManagerFromNodeResourceTopology(logr.Discard(), nrt2)
				if !reflect.DeepEqual(tcase.exp, got) {
					t.Fatalf("expected=%q got=%q", tcase.exp.String(), got.String())
				}
			})
		}
	}
}
