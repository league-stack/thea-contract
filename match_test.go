package contract

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var sampleNames = []string{"match-bo3.json", "match-bo1-bots.json"}

func decodeSample(t *testing.T, name string) (Match, []byte) {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	match, err := DecodeMatch(raw)
	if err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	return match, raw
}

func compact(t *testing.T, raw []byte) string {
	t.Helper()
	var out bytes.Buffer
	if err := json.Compact(&out, raw); err != nil {
		t.Fatal(err)
	}
	return out.String()
}

func TestSamplesSurviveADecodeAndEncode(t *testing.T) {
	for _, name := range sampleNames {
		match, raw := decodeSample(t, name)
		encoded, err := json.Marshal(match)
		if err != nil {
			t.Fatal(err)
		}
		if string(encoded) != compact(t, raw) {
			t.Errorf("%s changed:\n got %s\nwant %s", name, encoded, compact(t, raw))
		}
	}
}

func TestBestOfThreeSampleHoldsItsSettings(t *testing.T) {
	match, _ := decodeSample(t, "match-bo3.json")
	if match.Settings != DefaultSettings() {
		t.Errorf("settings = %+v", match.Settings)
	}
	if len(match.Maps) != 3 || match.Maps[2].Sides != SidesKnife {
		t.Errorf("maps = %+v", match.Maps)
	}
}

func TestDecodeRefusesAnUnknownField(t *testing.T) {
	_, raw := decodeSample(t, "match-bo1-bots.json")
	withExtra := strings.Replace(string(raw), `"version": 1,`, `"version": 1, "cvars": {},`, 1)
	if _, err := DecodeMatch([]byte(withExtra)); err == nil {
		t.Fatal("an unknown field was accepted")
	}
}

func TestSidesSayWhichTeamStartsOnCT(t *testing.T) {
	cases := []struct {
		sides Sides
		team  int
		knife bool
	}{
		{SidesTeam1CT, 0, false},
		{SidesTeam2CT, 1, false},
		{SidesKnife, 0, true},
	}
	for _, c := range cases {
		team, knife := c.sides.StartingCT()
		if team != c.team || knife != c.knife {
			t.Errorf("%s = (%d, %v)", c.sides, team, knife)
		}
	}
}

func TestANewLineupNamesItsCaptains(t *testing.T) {
	before := Match{Teams: [2]Team{{Captain: "1", Players: []Player{{Steam: "1"}}}, {Players: []Player{{Steam: "2"}}}}}

	after := before.With(Lineup{Revision: 2, Teams: [2][]Player{{{Steam: "3"}}, {{Steam: "2"}}}, Captains: [2]string{"3", "2"}})

	if after.Teams[0].Captain != "3" || after.Teams[1].Captain != "2" || after.checkCaptains() != nil {
		t.Fatalf("teams = %+v", after.Teams)
	}
}
