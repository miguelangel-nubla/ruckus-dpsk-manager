package list

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/miguelangel-nubla/ruckus-dpsk-manager/internal/errors"
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/internal/filters"
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/pkg/client"
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/pkg/data/dpsk"
)

func Handle(svc *client.DpskService, args []string) error {
	filterArgs := args

	// Generate flags for filtering
	filtersFlagSet := flag.NewFlagSet("filter flags", flag.ExitOnError)
	filtersFlagSet.Usage = filters.FlagSetUsageOrdered(filtersFlagSet)

	filtersExactAll, err := filters.GenerateDpskFiltersExact(filtersFlagSet)
	if err != nil {
		return err
	}

	filtersRegexpAll, err := filters.GenerateDpskFiltersRegexp(filtersFlagSet, filterArgs)
	if err != nil {
		return err
	}

	// Parse the filter flags here so we can validate them
	filtersFlagSet.Parse(filterArgs)

	filtersExact, err := filters.ValidateFilters(filtersExactAll)
	if err != nil {
		return err
	}

	filtersRegexp, err := filters.ValidateFilters(filtersRegexpAll)
	if err != nil {
		return err
	}

	filterMap := make(map[string]dpsk.Filter)
	for k, filter := range filtersExact {
		filterMap[k] = filter
	}

	for k, filter := range filtersRegexp {
		if _, ok := filterMap[k]; ok {
			return fmt.Errorf("duplicate property filter: %s", k)
		}

		filterMap[k] = filter
	}

	if len(filterMap) == 0 {
		return &errors.CommandError{
			Msg:     "no filters specified",
			FlagSet: filtersFlagSet,
		}
	}

	// Filter validation end
	dpskList, err := svc.List()
	if err != nil {
		return fmt.Errorf("error getting DPSK list: %v", err)
	}

	matches, err := dpskList.Filter(filterMap)
	if err != nil {
		return fmt.Errorf("error filtering DPSK list: %v", err)
	}

	output, err := json.Marshal(matches)
	if err != nil {
		return err
	}

	fmt.Println(string(output))

	return nil
}
