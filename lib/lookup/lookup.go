package lookup

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Table struct {
	mapping map[string]string
}

func (t *Table) GetTag(port int32, protocol string) string {
	return t.mapping[fmt.Sprintf("%d-%s", port, strings.ToLower(protocol))]
}

func (t *Table) LoadFile(fname string) error {
	var (
		file *os.File
		err  error
		rec  []string

		header = true
	)
	if file, err = os.Open(fname); err != nil {
		return err
	}
	t.mapping = make(map[string]string)
	rdr := csv.NewReader(file)
	for {
		if rec, err = rdr.Read(); err != nil {
			break
		}
		if _, err = strconv.Atoi(rec[0]); err != nil {
			// Support both headered and headerless input
			if header {
				header = false
				continue
			} else {
				break
			}
		}
		header = false
		t.mapping[fmt.Sprintf("%s-%s", rec[0], strings.ToLower(rec[1]))] = rec[2]
	}
	if err != io.EOF {
		return err
	}
	return nil
}
