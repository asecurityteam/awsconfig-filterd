package filter

import "github.com/asecurityteam/awsconfig-filterd/pkg/domain"

// AnyMatch tries all filters until one returns true.
type AnyMatch []domain.ConfigFilterer

// FilterConfig calls the underlying filters until one returns true.
func (f AnyMatch) FilterConfig(c domain.ConfigurationItem) bool {
	for _, filter := range f {
		if filter.FilterConfig(c) {
			return true
		}
	}
	return false
}
