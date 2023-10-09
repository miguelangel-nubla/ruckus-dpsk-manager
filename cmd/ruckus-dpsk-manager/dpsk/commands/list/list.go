package list

import (
	"encoding/xml"
	"flag"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/miguelangel-nubla/ruckus-dpsk-manager/internal/errors"
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/internal/helpers"
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/pkg/client"
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/pkg/data/dpsk"
)

type extendedFilter interface {
	dpsk.Filter
	Validate() (bool, error)
	String() string
	Property() string
}

type filterExact struct {
	value     *string
	validator func(string) (string, error)
	validated string
	property  string
}

func newFilterExact(property string, value *string, validator func(string) (string, error)) filterExact {
	return filterExact{
		value:     value,
		validator: validator,
		property:  property,
	}
}

func (filter *filterExact) Validate() (bool, error) {
	value := *filter.value
	if value == "" {
		filter.validated = ""
		return false, nil
	}

	validated, err := filter.validator(value)
	if err != nil {
		return false, err
	}

	filter.validated = validated
	return true, nil
}

func (filter *filterExact) Test(s string) bool {
	return filter.validated == s
}

func (filter *filterExact) String() string {
	return filter.validated
}

func (filter *filterExact) Property() string {
	return filter.property
}

type filterRegexp struct {
	value    *string
	r        *regexp.Regexp
	property string
}

func newFilterRegexp(property string, value *string) filterRegexp {
	return filterRegexp{
		value:    value,
		property: property,
	}
}

func (filter *filterRegexp) Validate() (bool, error) {
	pattern := *filter.value
	if pattern == "" {
		return false, nil
	}

	r, err := regexp.Compile(pattern)
	if err != nil {
		return false, fmt.Errorf("Failed to compile regex pattern %q: %v", pattern, err)
	}

	filter.r = r
	return true, nil
}

func (filter *filterRegexp) Test(s string) bool {
	return filter.r.MatchString(s)
}

func (filter *filterRegexp) String() string {
	return "regexp: " + *filter.value
}

func (filter *filterRegexp) Property() string {
	return filter.property
}

func generateDpskFiltersExact(flagSet *flag.FlagSet) (map[string]extendedFilter, error) {
	flagList := make(map[string]extendedFilter)

	t := reflect.TypeOf(dpsk.Dpsk{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fullTag, ok := field.Tag.Lookup("xml")
		property := field.Name
		if ok {
			parts := strings.Split(fullTag, ",")
			tag := parts[0]

			var filter filterExact
			switch tag {
			case "last-rekey", "next-rekey":
				filter = newFilterExact(
					tag,
					flagSet.String(tag, "", fmt.Sprintf("filter by %s, valid formats: Unix timestamp, RFC3339 or YYYY-MM-DD HH:MM:SS", tag)),
					func(v string) (string, error) {
						time, err := helpers.ParseTimestamp(v)
						if err != nil {
							return "", fmt.Errorf("invalid %s timestamp '%s': %s", tag, v, err.Error())
						}

						return strconv.FormatInt(time.Unix(), 10), nil
					},
				)
			case "mac":
				filter = newFilterExact(
					tag,
					flagSet.String(tag, "", fmt.Sprintf("filter by %s, valid formats: case insensitive AA:BB:CC:DD:EE:FF or aa-bb-cc-dd-ee-ff", tag)),
					func(v string) (string, error) {
						mac, ok := isValidMAC(v)
						if !ok {
							return "", fmt.Errorf("invalid %s address: %s", tag, v)
						}
						return mac.String(), nil
					},
				)
			default:
				filter = newFilterExact(
					tag,
					flagSet.String(tag, "", fmt.Sprintf("filter by %s", tag)),
					func(v string) (string, error) {
						// currently no additional validation
						return v, nil
					},
				)
			}

			flagList[property] = &filter
		}
	}

	return flagList, nil
}

func generateDpskFiltersRegexp(flagSet *flag.FlagSet, args []string) (map[string]extendedFilter, error) {
	flagList := make(map[string]extendedFilter)

	t := reflect.TypeOf(dpsk.Dpsk{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fullTag, ok := field.Tag.Lookup("xml")
		property := field.Name
		if ok {
			parts := strings.Split(fullTag, ",")
			tag := parts[0]
			flagName := "regexp-" + tag

			var filter filterRegexp
			switch tag {
			case "last-rekey", "next-rekey":
				filter = newFilterRegexp(
					tag,
					flagSet.String(flagName, "", fmt.Sprintf("filter by %s, format: unix timestamp", tag)),
				)

			case "mac":
				filter = newFilterRegexp(
					tag,
					flagSet.String(flagName, "", fmt.Sprintf("filter by %s, format: a6:b5:c4:d2:e2:f1 (lowercase)", tag)),
				)

			default:
				filter = newFilterRegexp(
					tag,
					flagSet.String(flagName, "", fmt.Sprintf("filter by %s", tag)),
				)

			}

			flagList[property] = &filter
		}
	}

	return flagList, nil
}

func validateFilters(filters map[string]extendedFilter) (map[string]extendedFilter, error) {
	filtersFinal := make(map[string]extendedFilter)
	for k, v := range filters {
		valid, err := v.Validate()
		if err != nil {
			return nil, err
		}

		if !valid {
			continue
		}

		filtersFinal[k] = v
	}

	return filtersFinal, nil
}

func extractValuesFromExactFilters(filters map[string]extendedFilter) (map[string]string, error) {
	values := make(map[string]string)
	for _, v := range filters {
		valid, err := v.Validate()
		if err != nil {
			return nil, err
		}

		if !valid {
			continue
		}
		values[v.Property()] = v.String()
	}

	return values, nil
}

func flagSetUsageOrdered(flagSet *flag.FlagSet) func() {
	return func() {
		{
			// Lists to segregate flags
			var regexpFlags []*flag.Flag
			var otherFlags []*flag.Flag

			flagSet.VisitAll(func(f *flag.Flag) {
				if strings.HasPrefix(f.Name, "regexp-") {
					regexpFlags = append(regexpFlags, f)
				} else {
					otherFlags = append(otherFlags, f)
				}
			})

			// Print non-regexp flags
			for _, f := range otherFlags {
				fmt.Fprintf(flagSet.Output(), "  -%s: %s\n", f.Name, f.Usage)
			}

			// Print regexp flags
			for _, f := range regexpFlags {
				fmt.Fprintf(flagSet.Output(), "  -%s: %s\n", f.Name, f.Usage)
			}
		}
	}
}

func Handle(svc *client.DpskService, args []string) error {
	filterArgs := args

	// Generate flags for filtering

	filtersFlagSet := flag.NewFlagSet("filter flags", flag.ExitOnError)
	filtersFlagSet.Usage = flagSetUsageOrdered(filtersFlagSet)

	filtersExactAll, err := generateDpskFiltersExact(filtersFlagSet)
	if err != nil {
		return err
	}

	filtersRegexpAll, err := generateDpskFiltersRegexp(filtersFlagSet, filterArgs)
	if err != nil {
		return err
	}

	// Parse the filter flags here so we can validate them
	filtersFlagSet.Parse(filterArgs)

	filtersExact, err := validateFilters(filtersExactAll)
	if err != nil {
		return err
	}

	filtersRegexp, err := validateFilters(filtersRegexpAll)
	if err != nil {
		return err
	}

	filters := make(map[string]dpsk.Filter)
	for k, filter := range filtersExact {
		filters[k] = filter
	}

	for k, filter := range filtersRegexp {
		if _, ok := filters[k]; ok {
			return fmt.Errorf("duplicate property filter: %s", k)
		}

		filters[k] = filter
	}

	if len(filters) == 0 {
		return &errors.CommandError{
			Msg:     "no filters specified",
			FlagSet: filtersFlagSet,
		}
	}

	// Filter validation end

	if svc.Client.Debug {
		fmt.Println("Filtering by:")
		for _, v := range filtersExact {
			fmt.Printf("  %s: %s\n", v.Property(), v)
		}
		for _, v := range filtersRegexp {
			fmt.Printf("  %s: %s\n", v.Property(), v)
		}
	}

	dpskList, err := svc.List()
	if err != nil {
		return fmt.Errorf("error getting DPSK list: %v", err)
	}

	matches, err := dpskList.Filter(filters)
	if err != nil {
		return fmt.Errorf("error filtering DPSK list: %v", err)
	}

	output, err := xml.MarshalIndent(matches, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(output))

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

func isValidMAC(mac string) (net.HardwareAddr, bool) {
	// Check if it's a valid MAC format using net package
	addr, err := net.ParseMAC(mac)
	if err != nil {
		return addr, false
	}

	// Use regex to ensure it adheres strictly to the xx:xx:xx:xx:xx:xx format
	matched, _ := regexp.MatchString(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`, mac)
	return addr, matched
}
