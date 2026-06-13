package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kernelstub/cognitor/internal/util"
	"github.com/kernelstub/cognitor/pkg/model"
)

func (s *Store) SaveSnapshot(ctx context.Context, snapshot model.Snapshot) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `delete from snapshots`); err != nil {
		return err
	}
	for _, table := range []string{"binaries", "artifacts", "services", "registry_keys"} {
		if _, err := tx.ExecContext(ctx, `delete from `+table); err != nil {
			return err
		}
	}
	if _, err := tx.ExecContext(ctx, `insert into snapshots(id,name,path,created_at) values(?,?,?,?)`, snapshot.ID, snapshot.Name, snapshot.Path, snapshot.CreatedAt.Format("2006-01-02T15:04:05Z07:00")); err != nil {
		return err
	}
	for _, binary := range snapshot.Binaries {
		if err := insertBinary(ctx, tx, binary); err != nil {
			return err
		}
	}
	for _, artifact := range snapshot.Artifacts {
		if err := insertArtifact(ctx, tx, artifact); err != nil {
			return err
		}
	}
	for _, service := range snapshot.Services {
		id := util.StableID(snapshot.ID, "service", service.Name)
		if _, err := tx.ExecContext(ctx, `insert into services(id,name,binary_path,permissions,start_type) values(?,?,?,?,?)`, id, service.Name, service.BinaryPath, service.Permissions, service.StartType); err != nil {
			return err
		}
	}
	for _, key := range snapshot.Registry {
		id := util.StableID(snapshot.ID, "registry", key.Path)
		if _, err := tx.ExecContext(ctx, `insert into registry_keys(id,path,acl,description) values(?,?,?,?)`, id, key.Path, key.ACL, key.Description); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func insertBinary(ctx context.Context, tx *sql.Tx, binary model.Binary) error {
	imports, err := encode(binary.Imports)
	if err != nil {
		return err
	}
	exports, err := encode(binary.Exports)
	if err != nil {
		return err
	}
	sections, err := encode(binary.Sections)
	if err != nil {
		return err
	}
	strings, err := encode(binary.Strings)
	if err != nil {
		return err
	}
	functions, err := encode(binary.Functions)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `insert into binaries(id,snapshot_id,path,name,kind,sha256,size,version,signer,imports_json,exports_json,sections_json,strings_json,functions_json,manifest) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, binary.ID, binary.SnapshotID, binary.Path, binary.Name, binary.Kind, binary.SHA256, binary.Size, binary.Version, binary.Signer, imports, exports, sections, strings, functions, binary.Manifest)
	return err
}

func insertArtifact(ctx context.Context, tx *sql.Tx, artifact model.Artifact) error {
	stringsRaw, err := encode(artifact.Strings)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `insert into artifacts(id,snapshot_id,path,name,kind,sha256,size,strings_json) values(?,?,?,?,?,?,?,?)`, artifact.ID, artifact.SnapshotID, artifact.Path, artifact.Name, artifact.Kind, artifact.SHA256, artifact.Size, stringsRaw)
	return err
}

func (s *Store) LoadSnapshot(ctx context.Context) (model.Snapshot, error) {
	var snapshot model.Snapshot
	var created string
	err := s.db.QueryRowContext(ctx, `select id,name,path,created_at from snapshots order by created_at desc limit 1`).Scan(&snapshot.ID, &snapshot.Name, &snapshot.Path, &created)
	if errors.Is(err, sql.ErrNoRows) {
		return snapshot, sql.ErrNoRows
	}
	if err != nil {
		return snapshot, err
	}
	rows, err := s.db.QueryContext(ctx, `select id,snapshot_id,path,name,kind,sha256,size,version,signer,imports_json,exports_json,sections_json,strings_json,functions_json,manifest from binaries order by name,path`)
	if err != nil {
		return snapshot, err
	}
	defer rows.Close()
	for rows.Next() {
		binary, err := scanBinary(rows)
		if err != nil {
			return snapshot, err
		}
		snapshot.Binaries = append(snapshot.Binaries, binary)
	}
	if err := rows.Err(); err != nil {
		return snapshot, err
	}
	artifactRows, err := s.db.QueryContext(ctx, `select id,snapshot_id,path,name,kind,sha256,size,strings_json from artifacts order by path`)
	if err != nil {
		return snapshot, err
	}
	defer artifactRows.Close()
	for artifactRows.Next() {
		artifact, err := scanArtifact(artifactRows)
		if err != nil {
			return snapshot, err
		}
		snapshot.Artifacts = append(snapshot.Artifacts, artifact)
	}
	if err := artifactRows.Err(); err != nil {
		return snapshot, err
	}
	serviceRows, err := s.db.QueryContext(ctx, `select name,binary_path,permissions,start_type from services order by name`)
	if err != nil {
		return snapshot, err
	}
	defer serviceRows.Close()
	for serviceRows.Next() {
		var service model.Service
		if err := serviceRows.Scan(&service.Name, &service.BinaryPath, &service.Permissions, &service.StartType); err != nil {
			return snapshot, err
		}
		snapshot.Services = append(snapshot.Services, service)
	}
	if err := serviceRows.Err(); err != nil {
		return snapshot, err
	}
	registryRows, err := s.db.QueryContext(ctx, `select path,acl,description from registry_keys order by path`)
	if err != nil {
		return snapshot, err
	}
	defer registryRows.Close()
	for registryRows.Next() {
		var key model.RegistryKey
		if err := registryRows.Scan(&key.Path, &key.ACL, &key.Description); err != nil {
			return snapshot, err
		}
		snapshot.Registry = append(snapshot.Registry, key)
	}
	return snapshot, registryRows.Err()
}

func scanBinary(rows *sql.Rows) (model.Binary, error) {
	var binary model.Binary
	var imports, exports, sections, stringsRaw, functions string
	err := rows.Scan(&binary.ID, &binary.SnapshotID, &binary.Path, &binary.Name, &binary.Kind, &binary.SHA256, &binary.Size, &binary.Version, &binary.Signer, &imports, &exports, &sections, &stringsRaw, &functions, &binary.Manifest)
	if err != nil {
		return binary, err
	}
	if binary.Imports, err = decode[[]string](imports); err != nil {
		return binary, err
	}
	if binary.Exports, err = decode[[]string](exports); err != nil {
		return binary, err
	}
	if binary.Sections, err = decode[[]model.Section](sections); err != nil {
		return binary, err
	}
	if binary.Strings, err = decode[[]string](stringsRaw); err != nil {
		return binary, err
	}
	if binary.Functions, err = decode[[]model.Function](functions); err != nil {
		return binary, err
	}
	return binary, nil
}

func scanArtifact(rows *sql.Rows) (model.Artifact, error) {
	var artifact model.Artifact
	var stringsRaw string
	err := rows.Scan(&artifact.ID, &artifact.SnapshotID, &artifact.Path, &artifact.Name, &artifact.Kind, &artifact.SHA256, &artifact.Size, &stringsRaw)
	if err != nil {
		return artifact, err
	}
	if artifact.Strings, err = decode[[]string](stringsRaw); err != nil {
		return artifact, err
	}
	return artifact, nil
}
