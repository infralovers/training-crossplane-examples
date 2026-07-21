// Package main implements the XBuckets composition function for the MinIO
// edition of the course (deck 04). It reads a list of names from an XBuckets
// composite resource and emits one MinIO Bucket managed resource per name.
//
// The Bucket is built as an UNSTRUCTURED composed resource, so this function has
// no dependency on the provider-minio Go module (its typed structs pull in
// crossplane-runtime/v2 and clash with the function SDK). It also imports no
// crossplane-runtime package directly: it uses the function-sdk-go re-exports and
// apimachinery's unstructured helpers, so it does not matter which
// crossplane-runtime version the SDK happens to resolve.
package main

import (
	"context"
	"fmt"

	"github.com/crossplane/function-sdk-go/logging"
	fnv1 "github.com/crossplane/function-sdk-go/proto/v1"
	"github.com/crossplane/function-sdk-go/request"
	"github.com/crossplane/function-sdk-go/resource"
	"github.com/crossplane/function-sdk-go/resource/composed"
	"github.com/crossplane/function-sdk-go/response"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Function returns one MinIO Bucket per name in the XBuckets spec.
type Function struct {
	fnv1.UnimplementedFunctionRunnerServiceServer
	log logging.Logger
}

// RunFunction reads spec.names from the observed composite and adds a Bucket to
// the desired state for each name.
func (f *Function) RunFunction(_ context.Context, req *fnv1.RunFunctionRequest) (*fnv1.RunFunctionResponse, error) {
	rsp := response.To(req, response.DefaultTTL)

	oxr, err := request.GetObservedCompositeResource(req)
	if err != nil {
		response.Fatal(rsp, fmt.Errorf("cannot get observed composite resource: %w", err))
		return rsp, nil
	}

	names, err := oxr.Resource.GetStringArray("spec.names")
	if err != nil {
		response.Fatal(rsp, fmt.Errorf("cannot read spec.names from the XBuckets resource: %w", err))
		return rsp, nil
	}

	desired, err := request.GetDesiredComposedResources(req)
	if err != nil {
		response.Fatal(rsp, fmt.Errorf("cannot get desired composed resources: %w", err))
		return rsp, nil
	}

	for _, name := range names {
		// Build the MinIO Bucket as an unstructured object: set only the fields
		// the provider understands, no Go types from provider-minio required.
		b := composed.New()
		b.SetAPIVersion("minio.crossplane.io/v1")
		b.SetKind("Bucket")
		b.SetName(name)

		if err := unstructured.SetNestedField(b.Object, name, "spec", "forProvider", "bucketName"); err != nil {
			response.Fatal(rsp, fmt.Errorf("cannot set bucketName for %q: %w", name, err))
			return rsp, nil
		}
		if err := unstructured.SetNestedField(b.Object, "us-east-1", "spec", "forProvider", "region"); err != nil {
			response.Fatal(rsp, fmt.Errorf("cannot set region for %q: %w", name, err))
			return rsp, nil
		}
		if err := unstructured.SetNestedField(b.Object, "default", "spec", "providerConfigRef", "name"); err != nil {
			response.Fatal(rsp, fmt.Errorf("cannot set providerConfigRef for %q: %w", name, err))
			return rsp, nil
		}

		desired[resource.Name(fmt.Sprintf("bucket-%s", name))] = &resource.DesiredComposed{Resource: b}
	}

	if err := response.SetDesiredComposedResources(rsp, desired); err != nil {
		response.Fatal(rsp, fmt.Errorf("cannot set desired composed resources: %w", err))
		return rsp, nil
	}

	f.log.Info("Added buckets to desired state", "count", len(names))
	return rsp, nil
}
