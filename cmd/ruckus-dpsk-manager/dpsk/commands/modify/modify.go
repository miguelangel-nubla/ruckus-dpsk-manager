package modify

import (
	"encoding/xml"
	"flag"
	"fmt"
	"net"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/miguelangel-nubla/ruckus-dpsk-manager/internal/errors"
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/internal/helpers"
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/pkg/client"
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/pkg/data/dpsk"
)

type filter struct {
	value     *string
	validator func(string) (string, error)
}

func dpskFieldsToFlags(flagSet *flag.FlagSet, args []string) (map[string]string, error) {
	findArgs := make(map[string]filter)

	t := reflect.TypeOf(dpsk.Dpsk{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fullTag, ok := field.Tag.Lookup("xml")
		if ok {
			parts := strings.Split(fullTag, ",")
			tag := parts[0]

			switch tag {
			case "last-rekey", "next-rekey":
				findArgs[tag] = filter{
					value: flagSet.String(tag, "", fmt.Sprintf("filter by %s, valid formats: Unix timestamp, RFC3339 or YYYY-MM-DD HH:MM:SS", tag)),
					validator: func(v string) (string, error) {
						time, err := helpers.ParseTimestamp(v)
						if err != nil {
							return "", fmt.Errorf("invalid %s timestamp: %s", tag, v)
						}

						return strconv.FormatInt(time.Unix(), 10), nil
					},
				}
			case "mac":
				findArgs[tag] = filter{
					value: flagSet.String(tag, "", fmt.Sprintf("filter by %s, valid formats: case insensitive AA:BB:CC:DD:EE:FF or aa-bb-cc-dd-ee-ff", tag)),
					validator: func(v string) (string, error) {
						mac, ok := isValidMAC(v)
						if !ok {
							return "", fmt.Errorf("invalid %s address: %s", tag, v)
						}
						return mac.String(), nil
					},
				}
			default:
				findArgs[tag] = filter{
					value: flagSet.String(tag, "", fmt.Sprintf("filter by %s", tag)),
					validator: func(v string) (string, error) {
						// currently no additional validation
						return v, nil
					},
				}
			}
		}
	}

	flagSet.Parse(args)

	// iterate over the filters and validate them
	filters := make(map[string]string)
	for k, v := range findArgs {
		value := *v.value
		if value != "" {
			validator := v.validator
			validated, err := validator(value)
			if err != nil {
				return nil, &errors.CommandError{
					Msg:     err.Error(),
					FlagSet: flagSet,
				}
			}
			filters[k] = validated
		}
	}

	return filters, nil
}

func Handle(svc *client.DpskService, args []string) error {
	modifyArgs := args
	pos := FindString(args, "set")
	if pos != -1 {
		modifyArgs = args[:pos]
	}

	filterFlags := flag.NewFlagSet("filter flags", flag.ExitOnError)
	filters, err := dpskFieldsToFlags(filterFlags, modifyArgs)
	if err != nil {
		return err
	}

	if len(filters) == 0 {
		return &errors.CommandError{
			Msg:     "no filters specified",
			FlagSet: filterFlags,
		}
	}

	if pos == -1 {
		return fmt.Errorf("set directive not found")
	}

	setFlags := flag.NewFlagSet("set flags", flag.ExitOnError)
	fields, err := dpskFieldsToFlags(setFlags, args[pos+1:])
	if err != nil {
		return err
	}

	if len(fields) == 0 {
		return &errors.CommandError{
			Msg:     "no properties specified",
			FlagSet: filterFlags,
		}
	}

	if svc.Client.Debug {
		fmt.Println("Filtering by:")
		for k, v := range filters {
			fmt.Printf("  %s: %s\n", k, v)
		}

		fmt.Println("Setting fields:")
		for k, v := range fields {
			fmt.Printf("  %s: %s\n", k, v)
		}
	}

	if len(filters) == 0 {
		if svc.Client.Debug {
			fmt.Print("Filter is empty, nothing changed\n")
		}
		return nil
	}

	dpskList, err := svc.List()
	if err != nil {
		return fmt.Errorf("error getting DPSK list: %v", err)
	}

	matches := dpskList.FindByFields(filters)

	// Print before
	if svc.Client.Debug {
		fmt.Printf("Found %d matches:\n", len(matches))
		output, err := xml.MarshalIndent(matches, "", "  ")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	}

	// modify the matches
	for _, dpsk := range matches {
		if err := svc.Modify(dpsk.ID, fields); err != nil {
			return fmt.Errorf("error modifying DPSK %d: %v", dpsk.ID, err)
		}
	}

	// Print after
	if svc.Client.Debug {
		dpskList, err = svc.List()
		if err != nil {
			return fmt.Errorf("error getting DPSK list: %v", err)
		}

		matches := dpskList.FindByFields(filters)

		fmt.Printf("Modified %d records:\n", len(matches))
		output, err := xml.MarshalIndent(matches, "", "  ")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	}

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
