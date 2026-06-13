# Rule Guide

Cognitor rules detect newly added defensive behavior in matched functions. Rules are intentionally evidence-oriented and defensive. They produce review leads, not exploit conclusions.

## How Rules Work

For each matched function pair, Cognitor computes:

- added calls,
- added strings,
- added operations,
- match similarity,
- match reason.

Rules inspect those additions for hardening signals. Findings are then scored by confidence, severity, and risk.

## Rule Families

### Access Control

Category: `access-control`

Signals:

- `SeAccessCheck`
- `NtAccessCheck`
- access mask validation
- object type validation

Researcher checks:

- Does the new check dominate the privileged operation?
- Is the checked access mask the same one used by the sensitive operation?
- Are sibling call sites missing the new guard?

### Memory And Bounds Safety

Category: `memory-safety`

Signals:

- `ProbeForRead`
- `ProbeForWrite`
- length checks
- bounds checks
- integer overflow checks
- null checks
- size validation

Researcher checks:

- Is user-controlled size captured before use?
- Are integer conversions or multiplications protected?
- Is the validated buffer the same buffer later copied or parsed?

### Native API / Syscall Boundary

Category: `native-api-boundary`

Signals:

- syscall boundary validation
- previous mode validation
- user/kernel buffer assumptions
- probe/capture behavior
- native `Nt*` or `Zw*` function hardening

Researcher checks:

- Is `PreviousMode` or equivalent caller mode now consulted?
- Are user-mode buffers probed or captured before trust?
- Does the change affect an exported native API such as `ntdll.dll` stubs or kernel syscall handlers?

### Handle And Object Validation

Category: `handle-validation`

Signals:

- `ObReferenceObjectByHandle`
- `ObOpenObjectByPointer`
- `ObGetObjectType`
- `GrantedAccess`
- `DesiredAccess`
- `OBJ_KERNEL_HANDLE`
- handle type validation
- object type validation

Researcher checks:

- Is the object type verified before use?
- Is granted access checked against the sensitive operation?
- Are kernel handles and user handles separated correctly?

### Object Lifetime

Category: `object-lifetime`

Signals:

- `ObReferenceObject`
- `ObDereferenceObject`
- `ExAcquireRundownProtection`
- `ExReleaseRundownProtection`
- reference count changes
- rundown protection
- use-after-free guard language

Researcher checks:

- Does the new reference cover all later uses?
- Are error paths balanced?
- Is rundown protection released exactly once?

### Token And Impersonation

Category: `privilege-boundary`

Signals:

- `SeSinglePrivilegeCheck`
- privilege checks
- token privilege validation
- impersonation validation

Researcher checks:

- Is the effective token or client token being checked?
- Is impersonation reverted on all paths?
- Are low-privilege callers still able to reach sibling paths?

### IOCTL Hardening

Category: `ioctl-hardening`

Signals:

- IOCTL strings
- `METHOD_NEITHER`
- input buffer validation
- `FILE_READ_DATA`
- `FILE_WRITE_DATA`

Researcher checks:

- Does the dispatch routine validate method, access, and buffer length?
- Are all IOCTL codes covered?
- Are output buffer lengths validated separately?

### RPC Hardening

Category: `rpc-hardening`

Signals:

- `RpcBindingInqAuthClient`
- `RpcBindingInqAuthClientEx`
- `RpcImpersonateClient`
- `RpcRevertToSelf`
- `RpcServerRegisterIf3`
- `RpcServerRegisterIfEx`
- `RPC_IF_ALLOW_SECURE_ONLY`
- authn/authz levels
- interface security callbacks

Researcher checks:

- Is the interface registered with stronger flags?
- Are authentication and authorization levels checked before privileged methods?
- Are impersonation windows balanced with revert calls?

### RPC / Structured Input Marshalling

Category: `marshalling-validation`

Signals:

- NDR/MIDL additions
- conformant array range checks
- wire length checks
- max count checks
- deserialization validation
- string binding validation

Researcher checks:

- Are array counts and buffer sizes validated before allocation or copy?
- Are nested structures validated recursively?
- Are client-controlled string bindings constrained?

### COM Hardening

Category: `com-hardening`

Signals:

- `CoInitializeSecurity`
- `CoImpersonateClient`
- `CoRevertToSelf`
- `LaunchPermission`
- `AccessPermission`
- `AppID`
- `CLSID`
- impersonation level changes
- DCOM hardening

Researcher checks:

- Did launch/access permissions become stricter?
- Did impersonation level change?
- Do AppID/CLSID registry entries match the changed server binary?

### ALPC Hardening

Category: `alpc-hardening`

Signals:

- ALPC port validation
- message attribute validation
- security context checks

Researcher checks:

- Are message lengths and attributes validated?
- Is peer identity checked before privileged handling?

### Registry And Service Hardening

Categories:

- `registry-hardening`
- `service-hardening`

Signals:

- registry ACL changes
- policy gate additions
- service permission changes
- start type changes
- service binary path changes

Researcher checks:

- Does a binary consume the changed policy value?
- Did default permissions become stricter?
- Are service control manager permissions still consistent with the threat model?

## Scoring

Findings receive:

- confidence: rule certainty multiplied by match similarity,
- severity: derived from risk score,
- risk score: weighted from severity, confidence, category, and evidence.

The report also assigns an executive risk posture:

- `informational`
- `moderate`
- `elevated`
- `high`

## Adding A New Rule

1. Create a file in `internal/analysis`.
2. Implement:

```go
type ExampleRule struct{}

func (ExampleRule) ID() string { return "example-rule" }

func (ExampleRule) Evaluate(ctx context.Context, change model.SemanticChange) []model.Finding {
    hits := hasAny(append(change.AddedCalls, change.AddedOps...), "ExampleSignal")
    if len(hits) == 0 {
        return nil
    }
    return []model.Finding{finding(change, "example-category", "Example defensive change", hits, 0.8*change.Similarity)}
}
```

3. Add it to `DefaultEngine()` in `internal/analysis/rules.go`.
4. Add tests in `internal/analysis/rules_test.go`.
5. Update this document.

Keep new rules defensive. Do not add exploit guidance, payload generation, or bypass instructions.
