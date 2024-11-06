package config

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

type DurationList struct {
	vals []time.Duration
}

func (d *DurationList) Vals() []time.Duration {
	return d.vals
}

func (d *DurationList) Decode(val string) error {
	if d == nil {
		*d = DurationList{} //nolint:exhaustruct
	}

	for _, v := range strings.Split(val, ",") {
		dur, err := time.ParseDuration(v)
		if err != nil {
			return errors.Wrapf(err, "failed to parse duration: %s", v)
		}

		d.vals = append(d.vals, dur)
	}

	return nil
}
