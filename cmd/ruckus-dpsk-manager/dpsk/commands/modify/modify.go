package modify

import (
	"flag"
	"fmt"

	"github.com/miguelangel-nubla/ruckus-dpsk-manager/internal/errors"
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/internal/filters"
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/pkg/client"
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/pkg/data/dpsk"
)

func Handle(svc *client.DpskService, args []string) error {
	filterArgs := args
	valueArgs := []string{}
	pos := FindString(args, "set")
	if pos != -1 {
		filterArgs = args[:pos]
		valueArgs = args[pos+1:]
	}

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

	fmt.Println("Filtering by:")
	for k, v := range filtersExact {
		fmt.Printf("  %s: %s\n", k, v)
	}
	for k, v := range filtersRegexp {
		fmt.Printf("  %s: %s\n", k, v)
	}

	if pos == -1 {
		return fmt.Errorf("set directive not found")
	}

	// Generate flags for setting values

	valuesFlagSet := flag.NewFlagSet("value flags", flag.ExitOnError)

	filtersExactAsValues, err := filters.GenerateDpskFiltersExact(valuesFlagSet)
	if err != nil {
		return err
	}

	// Parse the value flags here so we can validate them
	valuesFlagSet.Parse(valueArgs)

	valuesToSet, err := filters.ExtractValuesFromExactFilters(filtersExactAsValues)
	if err != nil {
		return err
	}

	if len(valuesToSet) == 0 {
		return &errors.CommandError{
			Msg:     "no properties specified to modify",
			FlagSet: filtersFlagSet,
		}
	}

	// Values validation end

	fmt.Println("Setting attributes:")
	for k, v := range valuesToSet {
		fmt.Printf("  %s: %s\n", k, v)
	}

	dpskListOriginal, err := svc.List()
	if err != nil {
		return fmt.Errorf("error getting original DPSK list: %v", err)
	}

	matches, err := dpskListOriginal.Filter(filterMap)
	if err != nil {
		return fmt.Errorf("error filtering original DPSK list: %v", err)
	}

	// modify the matches
	for _, dpsk := range matches {
		if err := svc.Modify(dpsk.ID, valuesToSet); err != nil {
			return fmt.Errorf("error modifying DPSK %d: %v", dpsk.ID, err)
		}
	}

	// iterate over dpskList to print the modified records
	fmt.Printf("Modified %d records successfully\n", len(matches))

	return nil
}

func FindString(slice []string, target string) int {
	for i, v := range slice {
		if v == target {
			return i
		}
	}
	return -1
}
