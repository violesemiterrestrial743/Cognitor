package store

import "database/sql"

func Migrate(db *sql.DB) error {
	statements := []string{
		`create table if not exists snapshots(id text primary key, name text not null, path text not null, created_at text not null)`,
		`create table if not exists binaries(id text primary key, snapshot_id text not null, path text not null, name text not null, kind text not null, sha256 text not null, size integer not null, version text not null, signer text not null, imports_json text not null, exports_json text not null, sections_json text not null, strings_json text not null, functions_json text not null, manifest text not null)`,
		`create table if not exists artifacts(id text primary key, snapshot_id text not null, path text not null, name text not null, kind text not null, sha256 text not null, size integer not null, strings_json text not null)`,
		`create table if not exists services(id text primary key, name text not null, binary_path text not null, permissions text not null, start_type text not null)`,
		`create table if not exists registry_keys(id text primary key, path text not null, acl text not null, description text not null)`,
		`create table if not exists findings(id text primary key, title text not null, affected_binary text not null, old_function text not null, new_function text not null, category text not null, confidence real not null, severity text not null, risk_score real not null, evidence_json text not null, old_evidence_json text not null, new_evidence_json text not null, reasoning text not null, sibling_hints_json text not null, audit_targets_json text not null, disclosure_note text not null)`,
		`create table if not exists change_summaries(id text primary key, summary_json text not null)`,
		`create table if not exists graph_nodes(id text primary key, type text not null, label text not null, attrs_json text not null)`,
		`create table if not exists graph_edges(id text primary key, from_id text not null, to_id text not null, type text not null, attrs_json text not null)`,
	}
	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return err
		}
	}
	return nil
}
