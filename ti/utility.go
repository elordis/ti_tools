package ti

import "strconv"

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
