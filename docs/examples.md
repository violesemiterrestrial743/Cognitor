# Examples

This page shows common Cognitor workflows. Use binaries and snapshots you are authorized to analyze.

## 1. Fast Local Demo

Run the included fixture:

```sh
go run ./cmd/cognitor compare ./testdata/snapshots/old ./testdata/snapshots/new --workdir ./out --all-formats
```

Open the Markdown report:

```sh
less ./out/report.md
```

Generated files:

```text
out/findings.db
out/report.md
out/report.json
out/report.sarif
out/report.csv
out/cognitor-bundle.json
```

## 2. Patch-Diff Two Real Folders

Prepare two directories:

```text
snapshots/old
snapshots/new
```

Place older files in `old` and newer files in `new`, preserving relative paths where possible. Then run:

```sh
cognitor compare snapshots/old snapshots/new --workdir out --all-formats
```

## 3. Focus On One DLL

For `ntdll.dll`:

```sh
cognitor compare snapshots/old snapshots/new --focus ntdll.dll --workdir out-ntdll --all-formats
```

For all DLLs:

```sh
cognitor compare snapshots/old snapshots/new --focus "*.dll" --workdir out-dlls --all-formats
```

For a path-specific target:

```sh
cognitor compare snapshots/old snapshots/new --focus "system32/ntdll.dll" --workdir out-ntdll
```

Focus matching is case-insensitive and supports names, relative paths, and globs.

## 4. Driver Patch Review

```sh
cognitor compare snapshots/old snapshots/new --focus "*.sys" --workdir out-drivers --all-formats
```

Look first at:

- `Priority Review Queue`,
- `Top Findings`,
- `Automatic Change Inventory`,
- `Researcher Checklist`.

## 5. CI Gate

Fail a pipeline when a high-severity finding appears:

```sh
cognitor compare snapshots/old snapshots/new --workdir out --all-formats --fail-on high
```

Fail on medium or higher:

```sh
cognitor compare snapshots/old snapshots/new --workdir out --all-formats --fail-on medium
```

Archive these files from CI:

```text
out/report.md
out/report.json
out/report.sarif
out/report.csv
out/cognitor-bundle.json
```

## 6. Staged Pipeline

Use this when you want intermediate SQLite databases:

```sh
cognitor scan --snapshot old --path snapshots/old --out old.db
cognitor scan --snapshot new --path snapshots/new --out new.db
cognitor diff --old old.db --new new.db --out findings.db
cognitor report --db findings.db --format markdown --out report.md
cognitor report --db findings.db --format json --out report.json
cognitor report --db findings.db --format sarif --out report.sarif
cognitor report --db findings.db --format csv --out report.csv
```

## 7. Initialize Snapshot Folders

```sh
cognitor snapshot create --name old --path snapshots/old
cognitor snapshot create --name new --path snapshots/new
```

Copy supported files from a source tree:

```sh
cognitor snapshot create --name new --path snapshots/new --source /mnt/windows-build
```

## 8. Add Disassembler Sidecars

For a binary:

```text
snapshots/new/ntdll.dll
snapshots/new/ntdll.dll.analysis.json
```

Example sidecar:

```json
{
  "functions": [
    {
      "name": "NtCreateFile",
      "basic_block_count": 15,
      "calls": ["NtCreateFile", "RtlValidSecurityDescriptor", "ObReferenceObjectByHandle"],
      "strings": ["object manager", "access mask validation", "syscall previous mode validation"],
      "operations": ["syscall dispatch", "access mask validation", "length check", "handle type validation"]
    }
  ]
}
```

Sidecars improve function matching and give rules stronger evidence.

## 9. Add Service And Registry Context

`services.json`:

```json
[
  {
    "name": "SampleSvc",
    "binary_path": "driver.sys",
    "permissions": "restricted",
    "start_type": "manual"
  }
]
```

`registry.json`:

```json
[
  {
    "path": "HKLM\\Software\\Sample",
    "acl": "administrators-only",
    "description": "sample defensive fixture"
  }
]
```

These inputs help reports call out service and policy drift.

## 10. Read The Report

Recommended review order:

1. Executive Summary
2. Priority Review Queue
3. Analyst Guidance
4. Top Findings
5. Automatic Change Inventory
6. Sibling Bug Hypotheses
7. Recommended Manual Review Plan

Treat findings as review leads. Cognitor highlights defensive hardening signals; it does not prove exploitability.
