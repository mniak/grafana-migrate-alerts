package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDuration_Unmarshal(t *testing.T) {
	testdata := []struct {
		expression string
		duration   Duration
	}{
		{
			expression: "1d2h3m4s",
			duration:   Duration(24*time.Hour + 2*time.Hour + 3*time.Minute + 4*time.Second),
		},
		{
			expression: "2h",
			duration:   Duration(2 * time.Hour),
		},
		{
			expression: "3m",
			duration:   Duration(3 * time.Minute),
		},
		{
			expression: "4s",
			duration:   Duration(4 * time.Second),
		},
		{
			expression: "5d",
			duration:   Duration(5 * 24 * time.Hour),
		},
		{
			expression: "6ms",
			duration:   Duration(24 * time.Millisecond),
		},
	}
	for _, td := range testdata {
		t.Run(td.expression, func(t *testing.T) {
			j, err := json.Marshal(td.expression)
			require.NoError(t, err)

			var dur Duration
			err = json.Unmarshal(j, &dur)
			require.NoError(t, err)

			assert.Equal(t, td.duration, dur)
		})
	}
}

func TestDuration_UnmarshalFail(t *testing.T) {
	testdata := []struct {
		expression string
	}{
		{
			expression: "",
		},
		{
			expression: "-1s-",
		},
		{
			expression: "1x",
		},
		{
			expression: "1us",
		},
	}
	for _, td := range testdata {
		t.Run(td.expression, func(t *testing.T) {
			j, err := json.Marshal(td.expression)
			require.NoError(t, err)

			var dur Duration
			err = json.Unmarshal(j, &dur)
			assert.EqualError(t, err, "invalid duration")
		})
	}
}
