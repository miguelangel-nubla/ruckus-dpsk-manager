package backup

import (
	"flag"
	"fmt"

	"github.com/miguelangel-nubla/ruckus-dpsk-manager/pkg/client"
)

var (
	outputFileFlag = flag.String("output", "file.bak", "Output file location")
)

func Handle(rc *client.Client, args []string) error {
	flag.Parse()

	err := rc.Backup(*outputFileFlag)
	if err != nil {
		return fmt.Errorf("error saving backup: %v", err)
	}

	return nil
}
