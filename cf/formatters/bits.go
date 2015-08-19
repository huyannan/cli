package formatters

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	. "github.com/cloudfoundry/cli/cf/i18n"
)

const (
	BIT     = 1.0
	KILOBIT = 1024 * BIT
	MEGABIT = 1024 * KILOBIT
	GIGABIT = 1024 * MEGABIT
	TERABIT = 1024 * GIGABIT
)

func BitSize(bits int64) string {
	unit := ""
	value := float32(bits)

	switch {
	case bits >= TERABIT:
		unit = "Tb"
		value = value / TERABIT
	case bits >= GIGABIT:
		unit = "Gb"
		value = value / GIGABIT
	case bits >= MEGABIT:
		unit = "Mb"
		value = value / MEGABIT
	case bits >= KILOBIT:
		unit = "Kb"
		value = value / KILOBIT
	case bits == 0:
		return "0"
	}

	stringValue := fmt.Sprintf("%.1f", value)
	stringValue = strings.TrimSuffix(stringValue, ".0")
	return fmt.Sprintf("%s%s", stringValue, unit)
}

func ToKilobits(s string) (int64, error) {
	parts := bitsPattern.FindStringSubmatch(strings.TrimSpace(s))
	if len(parts) < 3 {
		return 0, invalidBitQuantityError()
	}

	value, err := strconv.ParseInt(parts[1], 10, 0)
	if err != nil {
		return 0, invalidBitQuantityError()
	}

	var bits int64
	unit := strings.ToUpper(parts[2])
	switch unit {
	case "T":
		bits = value * TERABIT
	case "G":
		bits = value * GIGABIT
	case "M":
		bits = value * MEGABIT
	case "K":
		bits = value * KILOBIT
	}

	return bits / KILOBIT, nil
}

var (
	bitsPattern *regexp.Regexp = regexp.MustCompile(`(?i)^(-?\d+)([KMGT])b?$`)
)

func invalidBitQuantityError() error {
	return errors.New(T("Bit quantity must be an integer with a unit of measurement like K, Kb, M, Mb, G, or Gb"))
}
