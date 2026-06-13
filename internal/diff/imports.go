package diff

import "sort"

func AddedStrings(oldValues []string, newValues []string) []string {
	oldSet := map[string]struct{}{}
	for _, value := range oldValues {
		oldSet[value] = struct{}{}
	}
	var added []string
	for _, value := range newValues {
		if _, ok := oldSet[value]; !ok {
			added = append(added, value)
		}
	}
	sort.Strings(added)
	return added
}
