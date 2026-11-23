package contract

import "testing"

func validMatch(t *testing.T) Match {
	t.Helper()
	match, _ := decodeSample(t, "match-bo3.json")
	return match
}

func TestSamplesAreValid(t *testing.T) {
	for _, name := range sampleNames {
		match, _ := decodeSample(t, name)
		if refusal := match.Validate(); refusal != nil {
			t.Errorf("%s refused: %v", name, refusal)
		}
	}
}

func TestValidateRefusals(t *testing.T) {
	cases := []struct {
		name  string
		spoil func(*Match)
		want  RefusalCode
	}{
		{"another version", func(m *Match) { m.Version = Version + 1 }, RefusedVersion},
		{"no id", func(m *Match) { m.ID = "" }, RefusedNoID},
		{"no maps", func(m *Match) { m.Maps = nil }, RefusedNoMaps},
		{"a map with no name", func(m *Match) { m.Maps[0].Map = "" }, RefusedNoMaps},
		{"unknown sides", func(m *Match) { m.Maps[1].Sides = "team3_ct" }, RefusedBadSides},
		{"a short lineup", func(m *Match) { m.Teams[0].Players = m.Teams[0].Players[:4] }, RefusedLineupSize},
		{"a player on both teams", func(m *Match) { m.Teams[1].Players[0] = m.Teams[0].Players[0] }, RefusedDuplicatePlayer},
		{"a player who also spectates", func(m *Match) { m.Spectators[0] = m.Teams[0].Players[2] }, RefusedDuplicatePlayer},
		{"a player with no steam id", func(m *Match) { m.Teams[0].Players[1].Steam = "" }, RefusedNoSteam},
		{"no players per team", func(m *Match) { m.Settings.PlayersPerTeam = 0 }, RefusedBadSetting},
		{"no rounds", func(m *Match) { m.Settings.RoundsPerHalf = 0 }, RefusedBadSetting},
		{"unknown start rule", func(m *Match) { m.Settings.Start = "whenever" }, RefusedBadSetting},
		{"unknown demo rule", func(m *Match) { m.Settings.Demo = "maybe" }, RefusedBadSetting},
		{"overtime with no rounds", func(m *Match) { m.Settings.Overtime.RoundsPerHalf = 0 }, RefusedBadSetting},
		{"negative timeouts", func(m *Match) { m.Settings.Timeouts.PerTeam = -1 }, RefusedBadSetting},
		{"negative break", func(m *Match) { m.Settings.BreakSeconds = -1 }, RefusedBadSetting},
		{"no draw and no overtime", func(m *Match) { m.Settings.Overtime.Enabled = false }, RefusedBadSetting},
		{"a draw with overtime", func(m *Match) { m.Settings.AllowDraw = true }, RefusedBadSetting},
		{"negative ready limit", func(m *Match) { m.Settings.ReadySeconds = -1 }, RefusedBadSetting},
		{"negative technical pauses", func(m *Match) { m.Settings.TechPauses.PerTeam = -1 }, RefusedBadSetting},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			match := validMatch(t)
			c.spoil(&match)
			refusal := match.Validate()
			if refusal == nil || refusal.Code != c.want {
				t.Fatalf("refusal = %v, want %s", refusal, c.want)
			}
		})
	}
}

func TestOvertimeRoundsAreNotCheckedWhenOvertimeIsOff(t *testing.T) {
	match := validMatch(t)
	match.Settings.Overtime, match.Settings.AllowDraw = Overtime{}, true
	if refusal := match.Validate(); refusal != nil {
		t.Fatalf("refused: %v", refusal)
	}
}

func TestAShortLineupLoadsWhileLineupsAreOpen(t *testing.T) {
	match := validMatch(t)
	match.Teams[0].Players, match.OpenLineups = match.Teams[0].Players[:4], true

	if refusal := match.Validate(); refusal != nil {
		t.Fatalf("refused: %v", refusal)
	}
}

func TestAVetoStepNeedsAKnownActionAMapAndATeam(t *testing.T) {
	team, nobody := 1, 2
	for name, step := range map[string]VetoStep{
		"unknown action": {Action: "skip", Map: "de_mirage", Team: &team},
		"no map":         {Action: VetoBan, Team: &team},
		"no team":        {Action: VetoPick, Map: "de_mirage"},
		"a third team":   {Action: VetoBan, Map: "de_mirage", Team: &nobody},
	} {
		match := Match{Veto: []VetoStep{step}}
		if refusal := match.checkVeto(); refusal == nil || refusal.Code != RefusedBadVeto {
			t.Errorf("%s: refusal = %v", name, refusal)
		}
	}
	decided := Match{Veto: []VetoStep{{Action: VetoBan, Map: "de_nuke", Team: &team}, {Action: VetoDecider, Map: "de_mirage"}}}
	if refusal := decided.checkVeto(); refusal != nil {
		t.Errorf("a ban and a decider were refused: %v", refusal)
	}
}
