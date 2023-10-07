package modify

import (
	"flag"
	"fmt"

	"github.com/miguelangel-nubla/ruckus-dpsk-manager/internal/helpers"
	"github.com/miguelangel-nubla/ruckus-dpsk-manager/pkg/client"
)

func Handle(svc *client.DpskService, args []string) error {
	findCmd := flag.NewFlagSet("modify", flag.ExitOnError)
	expirationString := findCmd.String("expiration", "", "Expiration time, valid formats: Unix timestamp, RFC3339 or YYYY-MM-DD HH:MM:SS")
	findCmd.Parse(args)

	if *expirationString != "" {
		expiration, err := helpers.ParseTimestamp(*expirationString)
		if err != nil {
			err := fmt.Errorf("Failed to parse '%s': %s\n", *expirationString, err)
			return err
		}

		err = svc.Modify(12, expiration)
		if err != nil {
			err := fmt.Errorf("Failed to set expiration: %v", err)
			return err
		}
	}

	fmt.Println("Success")

	return nil
}
