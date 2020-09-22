package correlator

// LogCorrelatorInput
type LogCorrelatorInput struct {
	LogName    string
	OutputHash string
}

type LogCorrelatorOutput struct {
	// All lognames that matched this hash
	LogNames []string
	// The decided OutputHash
	OutputHash string
}

type LogCorrelator interface {
	// Decide - Decide on an output hash from aggregated logs
	Decide([]*LogCorrelatorInput) (*LogCorrelatorOutput, error)
}
