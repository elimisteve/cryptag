// Steve Phillips / elimisteve
// 2015.02.24

package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/elimisteve/cryptag"
	"github.com/elimisteve/fun"
	uuid "github.com/nu7hatch/gouuid"
)

type Row struct {
	// Populated by server
	Encrypted  []byte   `json:"data"`
	RandomTags []string `json:"tags"`

	// Populated locally
	decrypted []byte
	plainTags []string
	Nonce     *[24]byte `json:"nonce"`
}

var (
	ErrRowsNotFound = errors.New("No rows found")
)

func NewRow(decrypted []byte, plainTags []string) (*Row, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("Error generating new UUID for Row: %v", err)
	}
	// TODO(elimisteve): Document `id:`-prefix and related conventions
	uuidTag := "id:" + id.String()

	// For future queryability-related reasons, UUID must come first!
	plainTags = append([]string{uuidTag}, append(plainTags, "all")...)

	nonce, err := cryptag.RandomNonce()
	if err != nil {
		return nil, err
	}

	row := &Row{decrypted: decrypted, plainTags: plainTags, Nonce: nonce}

	return row, nil
}

func NewRowFromBytes(body []byte) (*Row, error) {
	row := &Row{}
	if err := json.Unmarshal(body, row); err != nil {
		return nil, fmt.Errorf("Error creating new row: `%v`. Input: `%s`", err,
			body)
	}
	if Debug {
		log.Printf("Created new Row `%#v` from bytes: `%s`\n", row, body)
	}
	return row, nil
}

func (row *Row) Decrypted() []byte {
	return row.decrypted
}

func (row *Row) PlainTags() []string {
	return row.plainTags
}

func (row *Row) HasRandomTag(randtag string) bool {
	return fun.SliceContains(row.RandomTags, randtag)
}

func (row *Row) HasPlainTag(plain string) bool {
	return fun.SliceContains(row.plainTags, plain)
}

func (row *Row) Format() string {
	return fmt.Sprintf("%s\t%s\n", row.decrypted, strings.Join(row.plainTags, "  "))
}

// Decrypt sets row.decrypted, row.nonce based upon row.Encrypted,
// nonce.  The passed-in `decrypt` function will typically be
// bkend.Decrypt, where `bkend` is the backend storing this Row.
func (row *Row) Decrypt(decrypt func(data []byte, nonce *[24]byte) ([]byte, error)) error {
	if len(row.Encrypted) == 0 {
		if Debug {
			log.Printf("row.Decrypt: no data to decrypt, returning nil (no error)\n")
		}
		return nil
	}

	dec, err := decrypt(row.Encrypted, row.Nonce)
	if err != nil {
		return fmt.Errorf("Error decrypting: %v", err)
	}

	row.decrypted = dec

	return nil
}

// SetPlainTags uses row.RandomTags and retrieved TagPairs to set
// row.plainTags
func (row *Row) SetPlainTags(getPairsFromRandom func(plain []string) (TagPairs, error)) error {
	pairs, err := getPairsFromRandom(row.RandomTags)
	if err != nil {
		return err
	}

	row.plainTags = pairs.AllPlain()
	if Debug {
		log.Printf("row.plainTags set to `%#v`\n", row.plainTags)
	}

	return nil
}
