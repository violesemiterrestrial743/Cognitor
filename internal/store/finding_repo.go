package store

import (
	"context"

	"github.com/kernelstub/cognitor/pkg/model"
)

func (s *Store) SaveFindings(ctx context.Context, findings []model.Finding) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `delete from findings`); err != nil {
		return err
	}
	for _, finding := range findings {
		evidence, err := encode(finding.Evidence)
		if err != nil {
			return err
		}
		oldEvidence, err := encode(finding.OldEvidence)
		if err != nil {
			return err
		}
		newEvidence, err := encode(finding.NewEvidence)
		if err != nil {
			return err
		}
		hints, err := encode(finding.SiblingBugSearchHints)
		if err != nil {
			return err
		}
		targets, err := encode(finding.RecommendedAuditTargets)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, `insert into findings(id,title,affected_binary,old_function,new_function,category,confidence,severity,risk_score,evidence_json,old_evidence_json,new_evidence_json,reasoning,sibling_hints_json,audit_targets_json,disclosure_note) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, finding.ID, finding.Title, finding.AffectedBinary, finding.OldFunction, finding.NewFunction, finding.Category, finding.Confidence, finding.Severity, finding.RiskScore, evidence, oldEvidence, newEvidence, finding.Reasoning, hints, targets, finding.ResponsibleDisclosureNote)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) LoadFindings(ctx context.Context) ([]model.Finding, error) {
	rows, err := s.db.QueryContext(ctx, `select id,title,affected_binary,old_function,new_function,category,confidence,severity,risk_score,evidence_json,old_evidence_json,new_evidence_json,reasoning,sibling_hints_json,audit_targets_json,disclosure_note from findings order by risk_score desc,id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var findings []model.Finding
	for rows.Next() {
		var finding model.Finding
		var evidence, oldEvidence, newEvidence, hints, targets string
		err := rows.Scan(&finding.ID, &finding.Title, &finding.AffectedBinary, &finding.OldFunction, &finding.NewFunction, &finding.Category, &finding.Confidence, &finding.Severity, &finding.RiskScore, &evidence, &oldEvidence, &newEvidence, &finding.Reasoning, &hints, &targets, &finding.ResponsibleDisclosureNote)
		if err != nil {
			return nil, err
		}
		if finding.Evidence, err = decode[[]string](evidence); err != nil {
			return nil, err
		}
		if finding.OldEvidence, err = decode[[]string](oldEvidence); err != nil {
			return nil, err
		}
		if finding.NewEvidence, err = decode[[]string](newEvidence); err != nil {
			return nil, err
		}
		if finding.SiblingBugSearchHints, err = decode[[]string](hints); err != nil {
			return nil, err
		}
		if finding.RecommendedAuditTargets, err = decode[[]string](targets); err != nil {
			return nil, err
		}
		findings = append(findings, finding)
	}
	return findings, rows.Err()
}
