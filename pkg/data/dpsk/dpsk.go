package dpsk

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"strings"
)

// Define the struct to match the XML structure
type ajaxResponse struct {
	XMLName  xml.Name `xml:"ajax-response"`
	Response response `xml:"response"`
}

type response struct {
	Type     string       `xml:"type,attr"`
	ID       string       `xml:"id,attr"`
	Apstamgr apstamgrStat `xml:"apstamgr-stat"`
}

type apstamgrStat struct {
	DpskList dpskList `xml:"dpsk-list"`
}

type dpskList struct {
	Entries Entries `xml:"dpsk"`
}

type Entries []Dpsk

type Dpsk struct {
	ID           int    `xml:"id,attr"`
	RoleID       string `xml:"role-id,attr"`
	Mac          string `xml:"mac,attr"`
	WlansvcID    int    `xml:"wlansvc-id,attr"`
	DvlanID      int    `xml:"dvlan-id,attr"`
	User         string `xml:"user,attr"`
	LastRekey    string `xml:"last-rekey,attr"`
	NextRekey    string `xml:"next-rekey,attr"`
	Expire       string `xml:"expire,attr"`
	StartPoint   string `xml:"start-point,attr"`
	Passphrase   string `xml:"passphrase,attr"`
	IpAddr       string `xml:"ip-addr,attr"`
	CurSharedNum string `xml:"cur-shared-num,attr"`
	Usage        string `xml:"usage,attr"`
}

// Method to search for a specific DPSK user and return its passphrase
func (list *Entries) FindPassphrase(wlanID int, username string) (string, error) {
	for _, entry := range *list {
		if entry.User == username && entry.WlansvcID == wlanID {
			return entry.Passphrase, nil
		}
	}
	return "", fmt.Errorf("DPSK user not found for username: %s and wlanID: %d", username, wlanID)
}

func (list *Entries) FindByWlanUser(wlanID int, username string) (Dpsk, error) {
	for _, entry := range *list {
		if entry.User == username && entry.WlansvcID == wlanID {
			return entry, nil
		}
	}
	return Dpsk{}, fmt.Errorf("DPSK user not found for username: %s and wlanID: %d", username, wlanID)
}

func (list *Entries) FindByFields(fields map[string]string) Entries {
	var matches Entries

	for _, dpsk := range *list {
		t := reflect.TypeOf(dpsk)
		v := reflect.ValueOf(dpsk)

		match := true

		// iterate over struct fields
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fullTag, ok := field.Tag.Lookup("xml")
			if !ok {
				// no xml tag, ignore
				continue
			}

			parts := strings.Split(fullTag, ",")
			tag := parts[0]

			value, ok := fields[tag]
			if !ok {
				// field not in filter, ignore
				continue
			}

			if v.Field(i).String() != value {
				// field value does not match
				match = false
				break
			}
		}

		if match {
			// all checks passed, add to matches
			matches = append(matches, dpsk)
		}
	}

	return matches
}

func FromXml(xmlData []byte) (*Entries, error) {
	var response ajaxResponse
	if err := xml.Unmarshal(xmlData, &response); err != nil {
		return &Entries{}, fmt.Errorf("error unmarshalling XML: %v", err)
	}

	return &response.Response.Apstamgr.DpskList.Entries, nil
}
