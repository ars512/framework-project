package handlers

import (
	"errors"
	"strconv"
	"strings"
)

func parseUint(value string) (uint, error) {
	v, err := strconv.ParseUint(strings.TrimSpace(value), 10, 64)
	return uint(v), err
}

func parsePositiveInt(value string, fallback int) (int, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback, nil
	}

	v, err := strconv.Atoi(value)
	if err != nil || v <= 0 {
		return 0, errors.New("invalid positive integer")
	}

	return v, nil
}

func parseFloat(value string) (float64, error) {
	value = strings.TrimSpace(value)
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, errors.New("invalid number")
	}
	return v, nil
}
