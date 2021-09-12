package main

func distinct(duplicates []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, item := range duplicates {
		if _, value := keys[item]; !value {
			keys[item] = true
			list = append(list, item)
		}
	}
	return list
}
