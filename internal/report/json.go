package report

import (
	"encoding/json"

	"github.com/kernelstub/cognitor/pkg/model"
)

func JSON(report model.Report) ([]byte, error) {
	return json.MarshalIndent(report, "", "  ")
}
