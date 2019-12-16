package blacklist

import (
	"github.com/fekle/simple-blacklist/pkg/domainfilter"
	"github.com/fekle/simple-blacklist/pkg/fetchlist"
)

// Blacklist holds all required information while fetching, filtering and merging blacklists
type Blacklist struct {
	Source  string
	Domains []string
	NTotal  int
	NUniq   int
}

// NewBlacklist initialises a new Blacklist instance
func NewBlacklist(source string) *Blacklist {
	return &Blacklist{
		Source:  source,
		Domains: []string{},
		NTotal:  0,
		NUniq:   0,
	}
}

// Process fetches all the domains in a blacklist source and filters them, keeping all unique results
func (bl *Blacklist) Process() error {
	lines, err := fetchlist.Fetch(bl.Source)
	if err != nil {
		return err
	}

	// apply filters
	if bl.Domains, err = domainfilter.Filter([]domainfilter.FilterFn{
		domainfilter.TrimWhiteSpaceFilter,
		domainfilter.DropCommentsFilter,
		domainfilter.DropIPAddressesFilter,
		domainfilter.ExtractNoSpaceGroupsFilter,
		domainfilter.HasValidTLDFilter,
		domainfilter.LeadingDotToWildcard,
	}, lines); err != nil {
		panic(err)
	}
	bl.NTotal = len(bl.Domains)

	// get unique domains
	bl.Domains = domainfilter.Uniq(bl.Domains)
	bl.NUniq = len(bl.Domains)

	return nil
}

// Merge merges multiple blacklists into a single big blacklist containing all unique domains
func Merge(lists []*Blacklist) *Blacklist {
	bl := NewBlacklist("final")

	for _, list := range lists {
		bl.Domains = append(bl.Domains, list.Domains...)
		bl.NTotal += list.NTotal
	}

	bl.Domains = domainfilter.Uniq(bl.Domains)
	bl.NUniq = len(bl.Domains)

	return bl
}
