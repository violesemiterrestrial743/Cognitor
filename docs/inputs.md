# Snapshot Inputs

Cognitor compares two snapshot directories: an older state and a newer patched state.

## Directory Layout

Recommended layout:

```text
snapshots/
  old/
    system32/
      ntdll.dll
      kernel32.dll
    drivers/
      example.sys
    services.json
    registry.json
  new/
    system32/
      ntdll.dll
      kernel32.dll
    drivers/
      example.sys
    services.json
    registry.json
```

Preserve relative paths when possible. Cognitor matches binaries by normalized relative path using Windows-style case-insensitive comparison.

## Supported Binary Inputs

Cognitor scans:

- `.exe`
- `.dll`
- `.sys`
- `.drv`
- `.ocx`
- `.cpl`

DLLs are first-class targets. Use `--focus ntdll.dll`, `--focus kernel32.dll`, or `--focus "*.dll"` for library-focused reviews.

## Supported Evidence Artifacts

Cognitor also tracks:

- `.edb`
- `.dat`
- `.log`
- `.evtx`
- `.etl`
- `.reg`
- `.json`
- `.xml`
- `.ini`
- `.inf`
- `.cfg`
- `.conf`

Artifacts are hashed and string-scanned. Reports show added, removed, and changed artifacts with risk signals when possible.

## Sidecar Files

Sidecars add richer context without requiring Cognitor to directly integrate with a disassembler.

### Analysis Sidecar

Name:

```text
binary.dll.analysis.json
```

Shape:

```json
{
  "functions": [
    {
      "name": "NtCreateFile",
      "basic_block_count": 15,
      "calls": ["NtCreateFile", "RtlValidSecurityDescriptor", "ObReferenceObjectByHandle"],
      "strings": ["object manager", "access mask validation", "syscall previous mode validation"],
      "imports": ["ntoskrnl.exe!ObReferenceObjectByHandle"],
      "operations": ["syscall dispatch", "access mask validation", "length check", "handle type validation"]
    }
  ]
}
```

Fields:

- `name`: function symbol or synthetic name.
- `basic_block_count`: rough structural size.
- `calls`: called APIs or functions.
- `strings`: strings referenced by the function.
- `imports`: imports associated with the function.
- `operations`: normalized semantic notes from an exporter.

### Symbols Sidecar

Name:

```text
binary.dll.symbols.json
```

Shape:

```json
["NtCreateFile", "NtOpenProcess", "RtlAllocateHeap"]
```

### Version Sidecar

Name:

```text
binary.dll.version
```

Shape:

```text
10.0.22621.3593
```

### Manifest Sidecar

Names:

```text
binary.exe.manifest
binary.exe.manifest.json
```

Manifest content is carried into the binary model for reporting and later analysis.

## Service Context

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

## Registry Context

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

## Snapshot Creation

Initialize empty folders:

```sh
cognitor snapshot create --name old --path snapshots/old
cognitor snapshot create --name new --path snapshots/new
```

Copy supported inputs from a source tree:

```sh
cognitor snapshot create --name new --path snapshots/new --source /mnt/windows-build
```

## Practical Tips

- Keep old and new relative paths consistent.
- Include sidecars for high-value targets such as `ntdll.dll`, drivers, RPC services, COM brokers, and parsers.
- Include service and registry context when reviewing privilege boundaries.
- Use `--focus` to keep reports small during deep dives.
