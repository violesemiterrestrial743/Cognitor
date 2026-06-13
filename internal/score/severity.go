package score

func severity(score float64) string {
	switch {
	case score >= 8:
		return "high"
	case score >= 5:
		return "medium"
	default:
		return "low"
	}
}
