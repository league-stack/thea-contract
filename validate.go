package contract

import (
	"fmt"
	"slices"
)

type RefusalCode string

const (
	RefusedVersion         RefusalCode = "version"
	RefusedNoID            RefusalCode = "no_id"
	RefusedNoMaps          RefusalCode = "no_maps"
	RefusedBadSides        RefusalCode = "bad_sides"
	RefusedBadVeto         RefusalCode = "bad_veto"
	RefusedBadCaptain      RefusalCode = "bad_captain"
	RefusedLineupSize      RefusalCode = "lineup_size"
	RefusedNoSteam         RefusalCode = "no_steam"
	RefusedDuplicatePlayer RefusalCode = "duplicate_player"
	RefusedBadSetting      RefusalCode = "bad_setting"
	RefusedMatchLoaded     RefusalCode = "match_loaded"
	RefusedNoMatch         RefusalCode = "no_match"
	RefusedWrongPhase      RefusalCode = "wrong_phase"
	RefusedUnknownCommand  RefusalCode = "unknown_command"
	RefusedBadCommand      RefusalCode = "bad_command"
	RefusedNoCheckpoint    RefusalCode = "no_checkpoint"
	RefusedActionFailed    RefusalCode = "action_failed"
	RefusedBotsNotAllowed  RefusalCode = "bots_not_allowed"
	RefusedStaleLineup     RefusalCode = "stale_lineup"
	RefusedNotAllowed      RefusalCode = "not_allowed"
)

type Refusal struct {
	Code   RefusalCode `json:"code"`
	Detail string      `json:"detail,omitempty"`
}

func (r *Refusal) Error() string {
	return fmt.Sprintf("%s: %s", r.Code, r.Detail)
}

func Refuse(code RefusalCode, format string, args ...any) *Refusal {
	return &Refusal{Code: code, Detail: fmt.Sprintf(format, args...)}
}

func (m Match) Validate() *Refusal {
	checks := []func() *Refusal{
		m.checkIdentity,
		m.Settings.check,
		m.checkMaps,
		m.checkVeto,
		m.checkCaptains,
		m.checkSteamIDs,
		m.checkLineupSizes,
		m.checkDuplicates,
	}
	for _, check := range checks {
		if refusal := check(); refusal != nil {
			return refusal
		}
	}
	return nil
}

func (m Match) checkIdentity() *Refusal {
	if m.Version != Version {
		return Refuse(RefusedVersion, "version %d, want %d", m.Version, Version)
	}
	if m.ID == "" {
		return Refuse(RefusedNoID, "the match has no id")
	}
	return nil
}

func (m Match) checkMaps() *Refusal {
	if len(m.Maps) == 0 {
		return Refuse(RefusedNoMaps, "the match has no maps")
	}
	for i, plan := range m.Maps {
		if plan.Map == "" {
			return Refuse(RefusedNoMaps, "map %d has no name", i+1)
		}
		if !plan.Sides.known() {
			return Refuse(RefusedBadSides, "map %d sides %q", i+1, plan.Sides)
		}
	}
	return nil
}

func (m Match) checkVeto() *Refusal {
	for i, step := range m.Veto {
		switch {
		case !step.Action.known():
			return Refuse(RefusedBadVeto, "veto step %d action %q", i+1, step.Action)
		case step.Map == "":
			return Refuse(RefusedBadVeto, "veto step %d has no map", i+1)
		case step.Team != nil && *step.Team != 0 && *step.Team != 1:
			return Refuse(RefusedBadVeto, "veto step %d team must be 0 or 1", i+1)
		case step.Team == nil && step.Action != VetoDecider:
			return Refuse(RefusedBadVeto, "veto step %d needs a team", i+1)
		}
	}
	return nil
}

func (m Match) checkCaptains() *Refusal {
	for _, team := range m.Teams {
		listed := slices.ContainsFunc(team.Players, func(player Player) bool { return player.Steam == team.Captain })
		if team.Captain != "" && !listed {
			return Refuse(RefusedBadCaptain, "%s's captain is not in its lineup", team.Name)
		}
	}
	return nil
}

func (m Match) checkSteamIDs() *Refusal {
	for _, player := range m.everyone() {
		if player.Steam == "" {
			return Refuse(RefusedNoSteam, "%q has no steam id", player.Name)
		}
	}
	return nil
}

func (m Match) checkLineupSizes() *Refusal {
	if m.FillWithBots || m.OpenLineups {
		return nil
	}
	for i, team := range m.Teams {
		if len(team.Players) < m.Settings.PlayersPerTeam {
			return Refuse(RefusedLineupSize, "team %d has %d players, needs %d", i+1, len(team.Players), m.Settings.PlayersPerTeam)
		}
	}
	return nil
}

func (m Match) checkDuplicates() *Refusal {
	seen := map[string]bool{}
	for _, player := range m.everyone() {
		if seen[player.Steam] {
			return Refuse(RefusedDuplicatePlayer, "%s is listed twice", player.Steam)
		}
		seen[player.Steam] = true
	}
	return nil
}

func (m Match) everyone() []Player {
	return slices.Concat(m.Teams[0].Players, m.Teams[1].Players, m.Teams[0].Coaches, m.Teams[1].Coaches, m.Spectators)
}

func (s Settings) check() *Refusal {
	problems := []struct {
		wrong  bool
		detail string
	}{
		{s.PlayersPerTeam < 1, "players_per_team must be at least 1"},
		{s.RoundsPerHalf < 1, "rounds_per_half must be at least 1"},
		{!s.Start.known(), "start must be ready or connected"},
		{!s.Demo.known(), "demo must be required or wanted"},
		{!s.FriendlyFire.known(), "friendly_fire must be on, utility or off"},
		{s.Overtime.Enabled && s.Overtime.RoundsPerHalf < 1, "overtime rounds_per_half must be at least 1"},
		{s.Overtime.StartMoney < 0, "overtime start_money cannot be negative"},
		{s.Timeouts.PerTeam < 0 || s.Timeouts.Seconds < 0, "timeouts cannot be negative"},
		{s.TechPauses.PerTeam < 0 || s.TechPauses.MaxSeconds < 0, "tech_pauses cannot be negative"},
		{s.AllowDraw == s.Overtime.Enabled, "allow_draw must be on exactly when overtime is off"},
		{s.Disconnects.WaitSeconds < 0, "disconnects wait_seconds cannot be negative"},
		{s.Disconnects.SubstituteSeconds < 0, "disconnects substitute_seconds cannot be negative"},
		{s.ReadySeconds < 0, "ready_seconds cannot be negative"},
		{s.StartSeconds < 0, "start_seconds cannot be negative"},
		{s.SideChoiceSeconds < 0, "side_choice_seconds cannot be negative"},
		{s.BreakSeconds < 0, "break_seconds cannot be negative"},
		{s.NoShows.ReplacementSeconds < 0, "no_shows replacement_seconds cannot be negative"},
		{s.NoShows.Forfeit && s.NoShows.ReplacementSeconds > 0, "no_shows takes forfeit or replacement_seconds, not both"},
		{s.Surrender.Needed < 0 || s.Surrender.FromRound < 0, "surrender cannot be negative"},
		{s.Remake.UntilRound < 0, "remake until_round cannot be negative"},
		{s.VoteSeconds < 0 || s.IdleRounds < 0 || s.ClearSeconds < 0, "vote_seconds, idle_rounds and clear_seconds cannot be negative"},
	}
	for _, problem := range problems {
		if problem.wrong {
			return Refuse(RefusedBadSetting, "%s", problem.detail)
		}
	}
	return nil
}
