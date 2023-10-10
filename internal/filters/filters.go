package filters

import (
	"flag"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/miguelangel-nubla/ruckus-dpsk-manager/internal/helpers"
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/pkg/data/dpsk"
)

type ExtendedFilter interface {
	dpsk.Filter
	Validate() (bool, error)
	String() string
}

type FilterExact struct {
	value     *string
	validator func(string) (string, error)
	validated string
}

func NewFilterExact(value *string, validator func(string) (string, error)) FilterExact {
	return FilterExact{
		value:     value,
		validator: validator,
	}
}

func (filter *FilterExact) Validate() (bool, error) {
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

func (filter *FilterExact) Test(s string) bool {
	return filter.validated == s
}

func (filter *FilterExact) String() string {
	return filter.validated
}

type FilterRegexp struct {
	value *string
	r     *regexp.Regexp
}

func NewFilterRegexp(value *string) FilterRegexp {
	return FilterRegexp{
		value: value,
	}
}

func (filter *FilterRegexp) Validate() (bool, error) {
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

func (filter *FilterRegexp) Test(s string) bool {
	return filter.r.MatchString(s)
}

func (filter *FilterRegexp) String() string {
	return "regexp: " + *filter.value
}

func GenerateDpskFiltersExact(flagSet *flag.FlagSet) (map[string]ExtendedFilter, error) {
	flagList := make(map[string]ExtendedFilter)

	t := reflect.TypeOf(dpsk.Dpsk{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		tag, ok := field.Tag.Lookup("_dpsk_attr")
		if ok {
			var filter FilterExact
			switch tag {
			case "last-rekey", "next-rekey":
				filter = NewFilterExact(
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
				filter = NewFilterExact(
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
				filter = NewFilterExact(
					flagSet.String(tag, "", fmt.Sprintf("filter by %s", tag)),
					func(v string) (string, error) {
						// currently no additional validation
						return v, nil
					},
				)
			}

			flagList[tag] = &filter
		}
	}

	return flagList, nil
}

func GenerateDpskFiltersRegexp(flagSet *flag.FlagSet, args []string) (map[string]ExtendedFilter, error) {
	flagList := make(map[string]ExtendedFilter)

	t := reflect.TypeOf(dpsk.Dpsk{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		tag, ok := field.Tag.Lookup("_dpsk_attr")
		if ok {
			flagName := "regexp-" + tag

			var filter FilterRegexp
			switch tag {
			case "last-rekey", "next-rekey":
				filter = NewFilterRegexp(
					flagSet.String(flagName, "", fmt.Sprintf("filter by %s, format: unix timestamp", tag)),
				)

			case "mac":
				filter = NewFilterRegexp(
					flagSet.String(flagName, "", fmt.Sprintf("filter by %s, format: a6:b5:c4:d2:e2:f1 (lowercase)", tag)),
				)

			default:
				filter = NewFilterRegexp(
					flagSet.String(flagName, "", fmt.Sprintf("filter by %s", tag)),
				)

			}

			flagList[tag] = &filter
		}
	}

	return flagList, nil
}

func ValidateFilters(filtersMap map[string]ExtendedFilter) (map[string]ExtendedFilter, error) {
	filtersFinal := make(map[string]ExtendedFilter)
	for k, v := range filtersMap {
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

func ExtractValuesFromExactFilters(filters map[string]ExtendedFilter) (map[string]string, error) {
	values := make(map[string]string)
	for k, v := range filters {
		valid, err := v.Validate()
		if err != nil {
			return nil, err
		}

		if !valid {
			continue
		}
		values[k] = v.String()
	}

	return values, nil
}

func FlagSetUsageOrdered(flagSet *flag.FlagSet) func() {
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
