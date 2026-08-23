# Third-party notices

## Grype executable image and vulnerability database

`adapters/grype-v0.116.1.yaml` references the official immutable
[Grype v0.116.1](https://github.com/anchore/grype/releases/tag/v0.116.1) OCI
image for offline repository dependency-vulnerability analysis. The image and
the separately acquired vulnerability database are not vendored in this
repository. Grype is copyright Anchore, Inc. and contributors and is licensed
under the [Apache License 2.0](https://github.com/anchore/grype/blob/v0.116.1/LICENSE).

Grype databases are distributed by Anchore from `grype.anchore.io`. A scan
accepts only a locally mounted schema-v6 database whose official archive URL
contains an immutable SHA-256 checksum and whose recorded build time is no more
than 120 hours old. Database content and provider feeds retain their respective
upstream terms; consult the database source metadata before redistribution.

## Syft executable image

`adapters/syft-v1.51.0.yaml` references the official immutable non-root
[Syft v1.51.0](https://github.com/anchore/syft/releases/tag/v1.51.0) OCI image
for offline software bill-of-materials generation. The image is not vendored in
this repository. Syft is copyright Anchore, Inc. and contributors and is
licensed under the [Apache License 2.0](https://github.com/anchore/syft/blob/v1.51.0/LICENSE).

## Gitleaks default configuration

`scanner/adapter/gitleaks-v8.30.0.toml.gz` is a deterministic gzip container for
the unmodified default configuration from
[Gitleaks v8.30.0](https://github.com/gitleaks/gitleaks/tree/v8.30.0), embedded
for deterministic offline secret scanning. The archive SHA-256 is
`da1838bd7cedb1bbd297163435c1675e777babb98ddf641a653f309733bd1fd1`;
the exact decompressed configuration SHA-256 is
`e163e53b9e7e8a8511e77271e2b323ed057759542a6d988258afe3a1fa329caf`.

The repository's `.gitleaksignore` contains one full fingerprint for the
historical uncompressed vendored configuration, whose own detector patterns
triggered one of its rules. It does not suppress a rule or current path.

MIT License

Copyright (c) 2019 Zachary Rice

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
