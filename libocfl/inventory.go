package libocfl

import (
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"time"
)

type InventoryConfig struct {
	DigestAlgorithm	string
	Id				string
}

type VersionConfig struct {
	UserName		string
	UserEmail		string
}

type InventoryBlock struct {
	DigestAlgorithm	string			`json:"digestAlgorithm"`
	Head			string			`json:"head"`
	Id				string			`json:"id"`
	Manifest		ManifestBlock	`json:"manifest"`
	Versions		map[string]VersionBlock	`json:"versions"`
	Fixity			FixityBlock		`json:"fixity,omitempty"`
}

type VersionBlock struct {
	Created 		string			`json:"created"`
	Message			string			`json:"message,omitempty"`
	State			StateBlock		`json:"state"`
	User			UserBlock		`json:"user"`
}

type UserBlock struct {
	Name			string			`json:"name"`
	Address			string			`json:"address"`
}

type StateBlock 	map[string][]string
type ManifestBlock 	map[string][]string
type FixityBlock	map[string]map[string][]string

func CreateBlankInventory(i *InventoryConfig, v *VersionConfig) (*InventoryBlock, error) {

	uid, _ := uuid.NewV4()

	vers := VersionBlock{
		Created: time.Now().Format(time.RFC3339),
		User: UserBlock{
			Name: v.UserName,
			Address: v.UserEmail,
		},
	}

	versMap := map[string]VersionBlock{
		"v1": vers,
	}

	inv := &InventoryBlock{
		DigestAlgorithm: i.DigestAlgorithm,
		Id:				fmt.Sprintf("urn:uuid:%s", uid),
		Head:			"v1",
		Versions:		versMap,
	}

	var err error

	return inv, err
}

func (i *InventoryBlock) ToJSON() ([]byte, error) {
	b, err := json.Marshal(i)
	return b, err
}

func (i *InventoryBlock) ToIndentedJSON() ([]byte, error) {
	b, err := json.MarshalIndent(i, "", "    ")
	return b, err
}