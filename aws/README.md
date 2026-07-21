# AWS S3 finale (deck 08)

The optional grand finale: the **same** `StaticSite` API from `minio/`, running on
**real AWS S3**. Only the Composition changes. The XRD, the request a developer
writes, and the go-templating loop are the ones built in decks 03 and 05.

> **Real cloud, real cost.** This creates billable S3 resources and needs an AWS
> account. It is a trainer demo by default; attendees with their own AWS keys can
> follow along. Clean up afterwards.

> **Status: validated-on-cluster pending.** Authored as the spec for the finale.
> Provider tag, S3 resource field names, `bucketSelector.matchControllerRef`, the
> external-name-as-bucket-name pattern, and the website endpoint host are flagged
> `# VALIDATE:` inline and must be checked on a real cluster before delivery.

## Prerequisites

- The StaticSite XRD already applied: `minio/staticsite-xrd.yaml`
- Both functions installed: `minio/patch-and-transform.yaml`, `minio/go-templating.yaml`
- An AWS access key and secret allowed to manage S3

## Apply

```shell
# 1. Install the AWS S3 provider and wait for it.
kubectl apply -f aws/01-provider-aws-s3.yaml
kubectl wait --for=condition=Healthy provider/provider-aws-s3 --timeout=300s

# 2. Add your AWS credentials, then the ProviderConfig.
#    Edit aws/02-providerconfig-aws.yaml first (replace the placeholder keys).
kubectl apply -f aws/02-providerconfig-aws.yaml

# 3. The AWS Composition (same XRD, new backend) and one request.
#    Edit aws/staticsite-aws.yaml: siteName must be a globally-unique bucket name.
kubectl apply -f aws/staticsite-composition-aws.yaml
kubectl apply -f aws/staticsite-aws.yaml
```

## Reach it

```shell
kubectl get staticsite hello-aws -o jsonpath='{.status.url}'
# open the printed http://<bucket>.s3-website-us-east-1.amazonaws.com URL
```

## Cleanup (avoid charges)

```shell
kubectl delete -f aws/staticsite-aws.yaml
kubectl delete -f aws/staticsite-composition-aws.yaml
# then remove the provider + config once no StaticSites remain
```
