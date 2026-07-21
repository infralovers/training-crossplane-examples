# training-crossplane-examples

Lab files for **Crossplane Essentials**. The course spine is a self-service
**StaticSite** platform API: a developer asks for a site, and Crossplane provisions
a MinIO bucket, a public-read policy, uploads the page content, and serves it
**inside the cluster**. No external GitLab or public cloud is required.

> **Status: validated-on-cluster pending.** These files were authored as the
> source-of-truth spec for the course rewrite. Provider package tags, the
> `provider-minio` ProviderConfig schema, the `provider-kubernetes` RBAC, and the
> `mc` upload Job must be validated on a real single-node K3s + MinIO before
> delivery. Flagged spots are marked `# VALIDATE:` inline.

## Layout

```
minio/      # the StaticSite spine (decks 02, 03, 05)
function/   # a Go composition function, XBuckets -> MinIO buckets (deck 04)
upbound/    # up-CLI control-plane project, MinIO edition (bonus)
aws/        # the same StaticSite API on real AWS S3 (deck 08 finale)
```

## Prerequisites in the training environment

- A single-node **K3s** cluster reachable from the workstation VM
- **Crossplane v2** installed (deck 02 does this)
- Helm and kubectl on the workstation

## Quick path (what the decks walk through)

```shell
# 1. Object store + providers (deck 02)
kubectl apply -f minio/00-minio.yaml
kubectl apply -f minio/01-provider-minio.yaml
kubectl apply -f minio/02-provider-kubernetes.yaml

# WAIT for the providers to install their CRDs before applying any
# ProviderConfig, otherwise you get "no matches for kind ProviderConfig".
kubectl wait --for=condition=Healthy provider/provider-minio provider/provider-kubernetes --timeout=180s

kubectl apply -f minio/03-providerconfig-minio.yaml
kubectl apply -f minio/04-providerconfig-kubernetes.yaml
kubectl apply -f minio/bucket.yaml            # warm-up: one bucket

# 2. The StaticSite platform API (deck 03)
kubectl apply -f minio/staticsite-xrd.yaml
kubectl apply -f minio/patch-and-transform.yaml
kubectl apply -f minio/staticsite-composition.yaml
kubectl apply -f minio/staticsite.yaml        # one XR -> a live in-cluster site

# 3. Dynamic, multi-page sites (deck 05)
kubectl apply -f minio/go-templating.yaml
kubectl apply -f minio/staticsite-dynamic-composition.yaml
kubectl apply -f minio/staticsite-dynamic.yaml
```

Reaching the site from the workstation VM:

```shell
kubectl -n minio port-forward svc/minio 9000:9000
# then open http://localhost:9000/<bucket>/index.html
```

Default MinIO credentials for the training are `minioadmin` / `minioadmin`.

## Grand finale: the same API on real AWS S3 (deck 08)

Optional, real-cloud, trainer demo by default. Reuses the StaticSite XRD and both
functions from `minio/`; only the Composition changes to target AWS S3. Needs an
AWS account. See `aws/README.md` for the apply order and cleanup.
