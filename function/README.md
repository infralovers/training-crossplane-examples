# function/ — XBuckets composition function (deck 04)

These are the two files attendees paste into a scaffolded function project. They
are **fragments**, not a standalone module: deck 04 first runs
`crossplane xpkg init function-xbuckets function-template-go -d function-xbuckets`,
which creates the `go.mod`, `main.go`, `Dockerfile`, and SDK wiring. Attendees
then replace `fn.go` and `fn_test.go` with these:

```shell
$ cp ~/training-crossplane-examples/function/fn.go      ~/function-xbuckets/fn.go
$ cp ~/training-crossplane-examples/function/fn_test.go ~/function-xbuckets/fn_test.go
$ go mod tidy
$ go test -v -cover -timeout 120s .
```

For the local `crossplane render` preview, also copy the three files in `example/`
over the template's generated defaults (whose `xr.yaml` has no `spec.names`):

```shell
$ cp ~/training-crossplane-examples/function/example/*.yaml ~/function-xbuckets/example/
```

The function reads `spec.names` from an `XBuckets` composite and emits one
`minio.crossplane.io/v1` Bucket per name, against the local MinIO from deck 02.

> **Why unstructured, not the provider's Go types?** provider-minio's apis are
> built on `crossplane-runtime/v2`, while `function-sdk-go` is on
> `crossplane-runtime` v1. Importing the typed `Bucket` struct pulls both into one
> build, and they require incompatible `controller-runtime` versions (the build
> fails with `SubResourceWriter ... missing method Apply`). Building the Bucket as
> an unstructured object keeps the project on a single runtime, so `go mod tidy`
> resolves cleanly and there is no `go get` of the provider module. The typed-
> provider approach is a good Advanced-course topic once the runtimes align.
