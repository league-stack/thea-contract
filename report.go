package contract

import "time"

type Phase string

const (
	PhaseIdle        Phase = "idle"
	PhaseLoading     Phase = "loading"
	PhaseWarmup      Phase = "warmup"
	PhaseKnife       Phase = "knife"
	PhaseSideChoice  Phase = "side_choice"
	PhaseLive        Phase = "live"
	PhaseHalftime    Phase = "halftime"
	PhaseBetweenMaps Phase = "between_maps"
	PhaseFinished    Phase = "finished"
	PhaseCancelled   Phase = "cancelled"
)

type ReportKind string

const (
	ReportState            ReportKind = "state"
	ReportMapLive          ReportKind = "map_live"
	ReportRoundEnded       ReportKind = "round_ended"
	ReportMapEnded         ReportKind = "map_ended"
	ReportSeriesEnded      ReportKind = "series_ended"
	ReportCheckpoint       ReportKind = "checkpoint"
	ReportNoShow           ReportKind = "no_show"
	ReportSubstituteNeeded ReportKind = "substitute_needed"
)

type Report struct {
	Seq        int64       `json:"seq"`
	Kind       ReportKind  `json:"kind"`
	State      State       `json:"state"`
	Winner     *int        `json:"winner,omitempty"`
	Checkpoint *Checkpoint `json:"checkpoint,omitempty"`
}

type State struct {
	MatchID        string        `json:"match_id"`
	Phase          Phase         `json:"phase"`
	Pause          Pause         `json:"pause,omitzero"`
	Hold           Hold          `json:"hold,omitzero"`
	Deadline       time.Time     `json:"deadline,omitzero"`
	StartsAt       time.Time     `json:"starts_at,omitzero"`
	Cancelled      CancelReason  `json:"cancelled,omitempty"`
	Forfeited      *int          `json:"forfeited,omitempty"`
	ForfeitReason  ForfeitReason `json:"forfeit_reason,omitempty"`
	Replacing      bool          `json:"awaiting_replacement,omitempty"`
	Vote           *Vote         `json:"vote,omitempty"`
	MapIndex       int           `json:"map_index"`
	Round          int           `json:"round"`
	Score          [2]int        `json:"score"`
	MapsWon        [2]int        `json:"maps_won"`
	MapScores      [][2]int      `json:"map_scores,omitempty"`
	CTTeam         int           `json:"ct_team"`
	KnifeWinner    *int          `json:"knife_winner,omitempty"`
	TechPausesUsed [2]int        `json:"tech_pauses_used"`
	TimeoutsUsed   [2]int        `json:"timeouts_used"`
	Players        []PlayerState `json:"players"`
	Commands       []CommandKind `json:"commands"`
}

type PlayerState struct {
	Steam       string    `json:"steam"`
	Team        int       `json:"team"`
	Connected   bool      `json:"connected"`
	Ready       bool      `json:"ready"`
	Left        int       `json:"left,omitempty"`
	LeftAt      time.Time `json:"left_at,omitzero"`
	AwaySeconds int       `json:"away_seconds,omitempty"`
	Idle        bool      `json:"idle,omitempty"`
	TeamDamage  int       `json:"team_damage,omitempty"`
	TeamKills   int       `json:"team_kills,omitempty"`
}

type ForfeitReason string

const (
	ForfeitByStaff    ForfeitReason = "staff"
	ForfeitDisconnect ForfeitReason = "disconnect"
	ForfeitNoShow     ForfeitReason = "no_show"
	ForfeitSurrender  ForfeitReason = "surrender"
	ForfeitGaveUp     ForfeitReason = "gave_up"
)

type VoteKind string

const (
	VoteSurrender VoteKind = "surrender"
	VotePlayOn    VoteKind = "play_on"
)

type Vote struct {
	Kind   VoteKind  `json:"kind"`
	Team   int       `json:"team"`
	Yes    []string  `json:"yes"`
	No     []string  `json:"no,omitempty"`
	Needed int       `json:"needed"`
	EndsAt time.Time `json:"ends_at"`
}

type PauseKind string

const (
	PauseStaff        PauseKind = "staff"
	PauseRestore      PauseKind = "restore"
	PauseRecovery     PauseKind = "recovery"
	PauseTech         PauseKind = "tech"
	PauseSubstitution PauseKind = "substitution"
	PauseDisconnect   PauseKind = "disconnect"
)

type Pause struct {
	Kind      PauseKind `json:"kind"`
	Reason    string    `json:"reason,omitempty"`
	ResumeAt  time.Time `json:"resume_at,omitzero"`
	GivesUpAt time.Time `json:"gives_up_at,omitzero"`
	CalledBy  *int      `json:"called_by,omitempty"`
	Agreed    [2]bool   `json:"agreed_to_resume,omitzero"`
}

func (p Pause) Active() bool { return p.Kind != "" }

type HoldReason string

const (
	HoldServerKeepsFailing HoldReason = "server_keeps_failing"
	HoldActionFailed       HoldReason = "action_failed"
	HoldRestoreFailed      HoldReason = "restore_failed"
	HoldPlayersMissing     HoldReason = "players_missing"
)

type Hold struct {
	Reason HoldReason `json:"reason"`
	Detail string     `json:"detail,omitempty"`
	Tried  []string   `json:"tried,omitempty"`
}

func (h Hold) Active() bool { return h.Reason != "" }

type Checkpoint struct {
	State          State          `json:"state"`
	LineupRevision int            `json:"lineup_revision"`
	DemoPart       int            `json:"demo_part"`
	Money          map[string]int `json:"money,omitempty"`
}

type CancelReason string

const (
	CancelledByStaff  CancelReason = "staff"
	CancelledNotReady CancelReason = "not_ready"
	CancelledNoShow   CancelReason = "no_show"
	CancelledRemake   CancelReason = "remake"
)
