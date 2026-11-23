package contract

type CommandKind string

const (
	CommandLoad        CommandKind = "load"
	CommandCancel      CommandKind = "cancel"
	CommandStart       CommandKind = "start"
	CommandPause       CommandKind = "pause"
	CommandUnpause     CommandKind = "unpause"
	CommandMoveOn      CommandKind = "move_on"
	CommandResume      CommandKind = "resume"
	CommandRestore     CommandKind = "restore"
	CommandReplayRound CommandKind = "replay_round"
	CommandNextMap     CommandKind = "next_map"
	CommandSetLineup   CommandKind = "set_lineup"
	CommandSwapSides   CommandKind = "swap_sides"
	CommandAnnounce    CommandKind = "announce"
	CommandForfeit     CommandKind = "forfeit"

	CommandReady         CommandKind = "ready"
	CommandUnready       CommandKind = "unready"
	CommandChooseSide    CommandKind = "choose_side"
	CommandTechPause     CommandKind = "tech_pause"
	CommandTimeout       CommandKind = "timeout"
	CommandAgreeToResume CommandKind = "agree_to_resume"
	CommandSurrender     CommandKind = "surrender"
	CommandRemake        CommandKind = "remake"
	CommandPlayOn        CommandKind = "play_on"
	CommandGiveUp        CommandKind = "give_up"

	CommandPoints CommandKind = "points"
)

type SideChoice string

const (
	ChoiceCT   SideChoice = "ct"
	ChoiceT    SideChoice = "t"
	ChoiceStay SideChoice = "stay"
	ChoiceSwap SideChoice = "swap"
)

type Command struct {
	ID     string      `json:"id"`
	Kind   CommandKind `json:"kind"`
	Match  *Match      `json:"match,omitempty"`
	Round  int         `json:"round,omitempty"`
	Winner *int        `json:"winner,omitempty"`
	Lineup *Lineup     `json:"lineup,omitempty"`
	Text   string      `json:"text,omitempty"`
	Team   *int        `json:"team,omitempty"`
	Steam  string      `json:"steam,omitempty"`
	Side   SideChoice  `json:"side,omitempty"`
	Points []Points    `json:"points,omitempty"`
}

type Points struct {
	Steam  string `json:"steam"`
	Change int    `json:"change"`
	Total  int    `json:"total"`
}

type Lineup struct {
	Revision   int         `json:"revision"`
	Teams      [2][]Player `json:"teams"`
	Coaches    [2][]Player `json:"coaches,omitzero"`
	Captains   [2]string   `json:"captains,omitzero"`
	Spectators []Player    `json:"spectators,omitempty"`
}

func (m Match) With(lineup Lineup) Match {
	for team := range m.Teams {
		m.Teams[team].Players, m.Teams[team].Coaches, m.Teams[team].Captain = lineup.Teams[team], lineup.Coaches[team], lineup.Captains[team]
	}
	m.Spectators, m.LineupRevision = lineup.Spectators, lineup.Revision
	return m
}

type Outcome string

const (
	OutcomeApplied Outcome = "applied"
	OutcomeRefused Outcome = "refused"
)

type Answer struct {
	CommandID string   `json:"command_id"`
	Outcome   Outcome  `json:"outcome"`
	Refusal   *Refusal `json:"refusal,omitempty"`
}
