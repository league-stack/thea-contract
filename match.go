package contract

import (
	"bytes"
	"encoding/json"
	"iter"
)

const Version = 1

type Match struct {
	Version        int        `json:"version"`
	ID             string     `json:"id"`
	Teams          [2]Team    `json:"teams"`
	Spectators     []Player   `json:"spectators,omitempty"`
	LineupRevision int        `json:"lineup_revision"`
	Settings       Settings   `json:"settings"`
	Maps           []MapPlan  `json:"maps"`
	Veto           []VetoStep `json:"veto,omitempty"`
	FillWithBots   bool       `json:"fill_with_bots,omitempty"`
	OpenLineups    bool       `json:"open_lineups,omitempty"`
}

func (m Match) LineupPlayers() iter.Seq2[int, Player] {
	return func(yield func(int, Player) bool) {
		for team, lineup := range m.Teams {
			for _, player := range lineup.Players {
				if !yield(team, player) {
					return
				}
			}
		}
	}
}

type Team struct {
	Ref     string   `json:"ref"`
	Name    string   `json:"name"`
	Tag     string   `json:"tag"`
	Players []Player `json:"players"`
	Coaches []Player `json:"coaches,omitempty"`
	Captain string   `json:"captain,omitempty"`
}

type Player struct {
	Steam string `json:"steam"`
	Name  string `json:"name"`
}

type MapPlan struct {
	Map   string `json:"map"`
	Sides Sides  `json:"sides"`
}

type VetoAction string

const (
	VetoPick    VetoAction = "pick"
	VetoBan     VetoAction = "ban"
	VetoDecider VetoAction = "decider"
)

func (a VetoAction) known() bool {
	return a == VetoPick || a == VetoBan || a == VetoDecider
}

type VetoStep struct {
	Action VetoAction `json:"action"`
	Map    string     `json:"map"`
	Team   *int       `json:"team,omitempty"`
}

type Sides string

const (
	SidesKnife   Sides = "knife"
	SidesTeam1CT Sides = "team1_ct"
	SidesTeam2CT Sides = "team2_ct"
)

func (s Sides) known() bool {
	return s == SidesKnife || s == SidesTeam1CT || s == SidesTeam2CT
}

func (s Sides) StartingCT() (team int, decidedByKnife bool) {
	switch s {
	case SidesTeam1CT:
		return 0, false
	case SidesTeam2CT:
		return 1, false
	}
	return 0, true
}

type StartRule string

const (
	StartWhenReady     StartRule = "ready"
	StartWhenConnected StartRule = "connected"
)

func (r StartRule) known() bool {
	return r == StartWhenReady || r == StartWhenConnected
}

type DemoRule string

const (
	DemoRequired DemoRule = "required"
	DemoWanted   DemoRule = "wanted"
)

func (r DemoRule) known() bool {
	return r == DemoRequired || r == DemoWanted
}

type Settings struct {
	PlayersPerTeam    int          `json:"players_per_team"`
	Start             StartRule    `json:"start"`
	ReadySeconds      int          `json:"ready_seconds"`
	RoundsPerHalf     int          `json:"rounds_per_half"`
	Overtime          Overtime     `json:"overtime"`
	AllowDraw         bool         `json:"allow_draw"`
	Timeouts          Timeouts     `json:"timeouts"`
	TechPauses        TechPauses   `json:"tech_pauses"`
	Disconnects       Disconnects  `json:"disconnects"`
	NoShows           NoShows      `json:"no_shows,omitzero"`
	Surrender         Surrender    `json:"surrender,omitzero"`
	Remake            Remake       `json:"remake,omitzero"`
	VoteSeconds       int          `json:"vote_seconds,omitempty"`
	IdleRounds        int          `json:"idle_rounds,omitempty"`
	CaptainsCall      bool         `json:"captains_call,omitempty"`
	ClearSeconds      int          `json:"clear_seconds,omitempty"`
	FriendlyFire      FriendlyFire `json:"friendly_fire,omitempty"`
	SideChoiceSeconds int          `json:"side_choice_seconds"`
	BreakSeconds      int          `json:"break_seconds"`
	Demo              DemoRule     `json:"demo"`
	DamageReport      bool         `json:"damage_report"`
	TeamIntro         bool         `json:"team_intro"`
}

type Overtime struct {
	Enabled       bool `json:"enabled"`
	RoundsPerHalf int  `json:"rounds_per_half"`
	StartMoney    int  `json:"start_money"`
}

type Timeouts struct {
	PerTeam int `json:"per_team"`
	Seconds int `json:"seconds"`
}

type TechPauses struct {
	Allowed    bool `json:"allowed"`
	PerTeam    int  `json:"per_team"`
	MaxSeconds int  `json:"max_seconds"`
}

type Disconnects struct {
	Pause       bool `json:"pause"`
	WaitSeconds int  `json:"wait_seconds"`
	Forfeit     bool `json:"forfeit"`
	Vote        bool `json:"vote,omitempty"`
}

type FriendlyFire string

const (
	FriendlyFireOn      FriendlyFire = "on"
	FriendlyFireUtility FriendlyFire = "utility"
	FriendlyFireOff     FriendlyFire = "off"
)

func (f FriendlyFire) known() bool {
	return f == "" || f == FriendlyFireOn || f == FriendlyFireUtility || f == FriendlyFireOff
}

type NoShows struct {
	Forfeit            bool `json:"forfeit"`
	ReplacementSeconds int  `json:"replacement_seconds,omitempty"`
}

type Surrender struct {
	Allowed   bool `json:"allowed"`
	Needed    int  `json:"needed,omitempty"`
	FromRound int  `json:"from_round,omitempty"`
}

type Remake struct {
	Allowed    bool `json:"allowed"`
	UntilRound int  `json:"until_round,omitempty"`
}

const DefaultVoteSeconds = 45

func (s Settings) VoteLasts() int {
	if s.VoteSeconds > 0 {
		return s.VoteSeconds
	}
	return DefaultVoteSeconds
}

func DefaultSettings() Settings {
	return Settings{
		PlayersPerTeam:    5,
		Start:             StartWhenReady,
		RoundsPerHalf:     12,
		Overtime:          Overtime{Enabled: true, RoundsPerHalf: 3, StartMoney: 12500},
		Timeouts:          Timeouts{PerTeam: 3, Seconds: 30},
		TechPauses:        TechPauses{Allowed: true, PerTeam: 2, MaxSeconds: 300},
		Disconnects:       Disconnects{Pause: true, WaitSeconds: 300},
		SideChoiceSeconds: 60,
		BreakSeconds:      300,
		Demo:              DemoRequired,
		DamageReport:      true,
	}
}

func DecodeMatch(raw []byte) (Match, error) {
	return decodeStrictly[Match](raw)
}

func DecodeCommand(raw []byte) (Command, error) {
	return decodeStrictly[Command](raw)
}

func decodeStrictly[T any](raw []byte) (T, error) {
	var decoded T
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&decoded)
	return decoded, err
}
