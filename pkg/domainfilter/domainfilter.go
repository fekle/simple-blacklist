package domainfilter

import (
	"golang.org/x/sync/errgroup"
	"sort"
	"sync"
)

// FilterFn implements a filter function
type FilterFn func(in string) (out []string)

// Filter applies a list of filters to a string slice
func Filter(filters []FilterFn, input []string) (output []string, err error) {
	output = input

	for _, f := range filters {
		var tmp []string
		var tmpMutex sync.Mutex

		var eg errgroup.Group

		for i, _ := range output {

			l := output[i]
			eg.Go(func() error {
				result := f(l)

				tmpMutex.Lock()
				tmp = append(tmp, result...)
				tmpMutex.Unlock()

				return nil
			})
		}

		if err := eg.Wait(); err != nil {
			return nil, err
		}

		output = tmp
	}

	return output, nil
}

// Uniq returns a slice of unique entries
func Uniq(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	sort.Strings(u)

	return u
}
