package disklayout

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// A FlagSet is a slice of bit-flags and their name.
type FlagSet []struct {
	Flag uint64
	Name string
}

// Parse returns a pretty version of val, using the flag names for known flags.
// Unknown flags remain numeric.
func (s FlagSet) Parse(val uint64) string {
	var flags []string

	for _, f := range s {
		if val&f.Flag == f.Flag {
			flags = append(flags, f.Name)
			val &^= f.Flag
		}
	}

	if val != 0 {
		flags = append(flags, "0x"+strconv.FormatUint(val, 16))
	}

	if len(flags) == 0 {
		// Prefer 0 to an empty string.
		return "0x0"
	}

	return strings.Join(flags, "|")
}

// ValueSet is a map of syscall values to their name. Parse will use the name
// or the value if unknown.
type ValueSet map[uint64]string

// Parse returns the name of the value associated with `val`. Unknown values
// are converted to hex.
func (s ValueSet) Parse(val uint64) string {
	if v, ok := s[val]; ok {
		return v
	}
	return fmt.Sprintf("%#x", val)
}

// ParseDecimal returns the name of the value associated with `val`. Unknown
// values are converted to decimal.
func (s ValueSet) ParseDecimal(val uint64) string {
	if v, ok := s[val]; ok {
		return v
	}
	return fmt.Sprintf("%d", val)
}

// ParseName returns the flag value associated with 'name'. Returns false
// if no value is found.
func (s ValueSet) ParseName(name string) (uint64, bool) {
	for k, v := range s {
		if v == name {
			return k, true
		}
	}
	return math.MaxUint64, false
}
