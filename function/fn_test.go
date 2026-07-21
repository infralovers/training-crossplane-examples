// Unit test for the XBuckets function (deck 04). Feeds a RunFunctionRequest with
// an observed XBuckets composite and asserts the function adds one MinIO Bucket
// per name to the desired state.
//
// VALIDATE: assertion details against the function-sdk-go version in go.mod;
// the structure mirrors the official function-template-go test.
package main

import (
	"context"
	"testing"

	"github.com/crossplane/function-sdk-go/logging"
	fnv1 "github.com/crossplane/function-sdk-go/proto/v1"
	"github.com/crossplane/function-sdk-go/resource"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestRunFunction(t *testing.T) {
	cases := map[string]struct {
		req       *fnv1.RunFunctionRequest
		wantCount int
	}{
		"AddThreeBuckets": {
			req: &fnv1.RunFunctionRequest{
				Observed: &fnv1.State{
					Composite: &fnv1.Resource{
						Resource: resource.MustStructJSON(`{
							"apiVersion": "example.crossplane.io/v1",
							"kind": "XBuckets",
							"metadata": {"name": "example-buckets"},
							"spec": {"names": ["site-a", "site-b", "site-c"]}
						}`),
					},
				},
			},
			wantCount: 3,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			f := &Function{log: logging.NewNopLogger()}
			rsp, err := f.RunFunction(context.Background(), tc.req)
			if err != nil {
				t.Fatalf("RunFunction returned error: %v", err)
			}
			got := len(rsp.GetDesired().GetResources())
			if diff := cmp.Diff(tc.wantCount, got); diff != "" {
				t.Errorf("desired bucket count mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// ensure structpb stays referenced if the template wiring needs it
var _ = structpb.NewStruct
