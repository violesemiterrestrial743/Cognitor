# Reports And Output Bundle

Cognitor can emit Markdown, JSON, SARIF, CSV, SQLite, and a bundle manifest.

## Full Bundle

```sh
cognitor compare old new --workdir out --all-formats
```

Outputs:

```text
out/findings.db
out/report.md
out/report.json
out/report.sarif
out/report.csv
out/cognitor-bundle.json
```

## Markdown

`report.md` is the primary analyst report.

Sections:

- Executive Summary
- Severity
- Likely Vulnerability Classes
- Analyst Guidance
- Priority Review Queue
- Automatic Change Inventory
- Top Changed Components
- Top Findings
- Semantic Change Clusters
- Sibling Bug Hypotheses
- Recommended Manual Review Plan

Recommended reading order:

1. Executive Summary
2. Priority Review Queue
3. Top Findings
4. Automatic Change Inventory
5. Researcher Checklist
6. Sibling Bug Hypotheses

## JSON

`report.json` preserves the full report model:

- metadata,
- summary,
- executive posture,
- changes,
- findings,
- graph.

Use JSON for custom dashboards, regression tests, and long-term archival.

## SARIF

`report.sarif` is intended for code scanning and security dashboards that understand SARIF 2.1.0.

Each Cognitor finding becomes a SARIF result with:

- rule ID,
- severity-mapped SARIF level,
- message,
- affected binary location.

## CSV

`report.csv` is a compact triage table for spreadsheet workflows.

Columns:

```text
type,target,category,severity,confidence,risk_score,reason,signals
```

CSV includes findings plus high-value binary and artifact changes.

## SQLite

`findings.db` stores:

- findings,
- graph nodes,
- graph edges,
- change summaries.

Use this when you want repeatable queries or staged processing.

## Bundle Manifest

`cognitor-bundle.json` records:

- generation time,
- tool version,
- old and new paths,
- risk level,
- priority,
- output files,
- SHA-256 hash for each output.

Example:

```json
{
  "generated_at": "2026-06-13T20:19:07Z",
  "tool_version": "0.1.0",
  "old_path": "./old",
  "new_path": "./new",
  "risk_level": "elevated",
  "priority": "same-day review",
  "outputs": [
    {
      "kind": "markdown",
      "path": "out/report.md",
      "sha256": "..."
    }
  ]
}
```

Use the manifest for handoff, CI retention, or verifying that reports were not modified after generation.

## CI Gates

Fail when a threshold is met:

```sh
cognitor compare old new --workdir out --all-formats --fail-on high
cognitor compare old new --workdir out --all-formats --fail-on medium
cognitor compare old new --workdir out --all-formats --fail-on low
```

Thresholds are inclusive. `--fail-on medium` fails on medium and high findings.

## Risk Posture

Executive risk posture is one of:

- `informational`
- `moderate`
- `elevated`
- `high`

It is based on finding severity, confidence, changed inventory, and prioritized review targets. It is a triage aid, not a vulnerability determination.
