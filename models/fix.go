package models

type Fix struct {
    Link
    Severity   string `json:"severity"`
    Reason     string `json:"reason"`
    Suggestion string `json:"suggestion"`
}
