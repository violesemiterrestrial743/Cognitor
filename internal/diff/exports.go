package diff

func RemovedStrings(oldValues []string, newValues []string) []string {
	newSet := map[string]struct{}{}
	for _, value := range newValues {
		newSet[value] = struct{}{}
	}
	var removed []string
	for _, value := range oldValues {
		if _, ok := newSet[value]; !ok {
			removed = append(removed, value)
		}
	}
	return removed
}
