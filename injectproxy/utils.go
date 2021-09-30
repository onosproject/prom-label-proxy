package injectproxy

func getEnterpriseName(groups []string) string {
	return groups[len(groups)-1]
}

func (r *routes) isAdminUser(groups []string) bool {
	if r.adminGroup == "" {
		return false
	}

	for _, gp := range groups {
		if gp == r.adminGroup {
			return true
		}
	}

	return false
}
