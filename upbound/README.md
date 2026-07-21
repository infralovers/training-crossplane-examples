# upbound/ — up-CLI project snippets (bonus deck)

Prepared copy-paste snippets for the bonus hands-on
`slidedecks/crossplane/bonus_crossplane_upbound_handson.md` (the **StorageBucket
on AWS S3** walkthrough). Attendees scaffold the project themselves with
`up project init` / `up ... generate`; these files are pasted into the generated
scaffolds so nobody debugs KCL live.

| File | Used on slide | Paste into |
|---|---|---|
| `main.k` | "Step 4: composition and function" | the generated `functions/test-function/main.k` |
| `test-main.k` | "Step 6 (optional): test the composition" | the generated test `main.k` |

- `main.k` reads `spec.parameters` (`region`, `versioning`, `acl`) from the
  `StorageBucket` XR and composes the S3 `Bucket`, `BucketOwnershipControls`,
  `BucketPublicAccessBlock`, `BucketACL`, `BucketServerSideEncryptionConfiguration`,
  and — when `versioning: true` — a `BucketVersioning`.
- There is intentionally **no `composition.yaml`** here: attendees generate the
  Composition with `up composition generate`. (The upstream repo still carries a
  stale GitLab `XTeamEnvironment` `composition.yaml` from a pre-rewrite bonus; it
  is unrelated to this example and is not shipped with the current course.)

> **VALIDATE:** `test-main.k` references `apis/xstoragebuckets/...` and
> `examples/storagebucket/xr.yaml`, whereas the deck's generated layout uses
> `apis/storagebuckets/definition.yaml` and `examples/storagebucket/example.yaml`.
> Confirm the exact directory/file names `up xrd generate` produces (the composite
> kind is likely `XStorageBucket` → `xstoragebuckets/`) and align the
> `compositionPath` / `xrdPath` / `xrPath` fields in the test before the optional
> `up test run` step.
