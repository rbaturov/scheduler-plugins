/*
 * Copyright 2023 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package pfpstatus

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/k8stopologyawareschedwg/podfingerprint"
)

func TestDumpNodeStatus(t *testing.T) {
	testDir, err := os.MkdirTemp("", "dumpnodest")
	if err != nil {
		t.Errorf("mkdirtemp failed: %v", err)
	}
	defer os.RemoveAll(testDir)

	fmt.Fprintf(os.Stderr, "testDir=%q\n", testDir)

	// copypasted from example: https://pkg.go.dev/time#Parse - no spacial meaning
	timeStamp, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	if err != nil {
		t.Errorf("time parse failed: %v", err)
	}

	st := StatusInfo{
		NodeName: "testNode-1",
		Data: podfingerprint.Status{
			FingerprintExpected: "PFPExpectedOrig",
			FingerprintComputed: "PFPComputedOrig",
			Pods: []podfingerprint.NamespacedName{
				{
					Namespace: "NSOrig1",
					Name:      "Name1",
				},
				{
					Namespace: "NSOrig2",
					Name:      "Name2",
				},
			},
		},
		LastModified: timeStamp,
	}

	err = DumpNodeStatus(testDir, &st)
	if err != nil {
		t.Errorf("dumpnodestatus failed: %v", err)
	}

	st2, err := LoadNodeStatus(testDir, st.NodeName)
	if err != nil {
		t.Errorf("loadnodestatus failed: %v", err)
	}

	if diff := cmp.Diff(st, *st2); diff != "" {
		t.Errorf("RTT check failed:\n%v", diff)
	}
}

func TestTracingStatusBasicOperation(t *testing.T) {
	ch := make(chan *StatusInfo, 10) // we just need a buffer > 1
	ts := TracingStatus{
		nodeName: "test-node-2",
		updates:  ch,
	}

	ts.Start(3)
	ts.Add("ns-1", "pod-1")
	ts.Add("ns-2", "pod-2")
	ts.Sign("random-value-sign") // not a real PFP
	ts.Check("random-value-chk") // not a real PFP

	if len(ch) != 1 {
		t.Errorf("data not sent in the updates channel")
	}

	ret := <-ch

	// copypasted from example: https://pkg.go.dev/time#Parse - no spacial meaning
	timeStamp, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	if err != nil {
		t.Errorf("time parse failed: %v", err)
	}
	ret.LastModified = timeStamp // reset to known value

	// extra check: received copy doesn't alter the original
	ret.Data.Add("ns-3", "pod-3")

	data, err := json.Marshal(ret)
	if err != nil {
		t.Errorf("error marshalling value from chan: %v", err)
	}

	expTs := `{"nodeName":"test-node-2","data":{"fingerprintExpected":"random-value-chk","fingerprintComputed":"random-value-sign","pods":[{"Namespace":"ns-1","Name":"pod-1"},{"Namespace":"ns-2","Name":"pod-2"},{"Namespace":"ns-3","Name":"pod-3"}]},"lastModified":"2006-01-02T15:04:05Z"}`
	if string(data) != expTs {
		t.Errorf("unexpected data from chan:\n%s\nexpected:\n%s", string(data), expTs)
	}

	retRepr := ret.Data.Repr()
	expRepr := `> processing 3 pods
+ ns-1/pod-1
+ ns-2/pod-2
+ ns-3/pod-3
= random-value-sign
V random-value-chk
`
	if retRepr != expRepr {
		t.Errorf("unexpected change in received repr:\n%s\nexpected:\n%s", string(retRepr), expRepr)
	}

	// extra check reprise: did altering the received copy changed back the original?
	origRepr := ts.Repr()
	if err != nil {
		t.Errorf("error marshalling orig value: %v", err)
	}

	expOrigRepr := `> processing 2 pods
+ ns-1/pod-1
+ ns-2/pod-2
= random-value-sign
V random-value-chk
`
	if origRepr != expOrigRepr {
		t.Errorf("unexpected change in orig repr:\n%s\nexpected:\n%s", string(origRepr), expOrigRepr)
	}
}
