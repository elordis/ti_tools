package ti

import (
	"strconv"
	"strings"
)

const NameForMissing = "-"
const AbbrLetters = 2

type UnmarshallableFloat64 float64

func (u *UnmarshallableFloat64) UnmarshalJSON(data []byte) error {
	fl, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		*u = UnmarshallableFloat64(0.0)

		return nil //nolint:nilerr
	}

	*u = UnmarshallableFloat64(fl)

	return nil
}

func abbreviate(s string, abbrevLimit int, keepLast bool) string { //nolint:unparam
	abbrev := ""
	parts := strings.Split(s, " ")

	max := len(parts)
	if keepLast {
		max--
	}

	for i := 0; i < max; i++ {
		abbrev += parts[i][:min(abbrevLimit, len(parts[i]))]
	}

	if keepLast {
		abbrev += " " + parts[max]
	}

	return abbrev
}
