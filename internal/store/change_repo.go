package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kernelstub/cognitor/pkg/model"
)

func (s *Store) SaveChangeSummary(ctx context.Context, summary model.ChangeSummary) error {
	raw, err := encode(summary)
	if err != nil {
		return err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `delete from change_summaries`); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `insert into change_summaries(id,summary_json) values(?,?)`, "latest", raw); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) LoadChangeSummary(ctx context.Context) (model.ChangeSummary, error) {
	var raw string
	err := s.db.QueryRowContext(ctx, `select summary_json from change_summaries where id = ?`, "latest").Scan(&raw)
	if errors.Is(err, sql.ErrNoRows) {
		return model.ChangeSummary{}, nil
	}
	if err != nil {
		return model.ChangeSummary{}, err
	}
	return decode[model.ChangeSummary](raw)
}
