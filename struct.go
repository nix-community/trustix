package main

// Input - Corresponds to Trillian Log LeafValue
type Input struct {
	inputHash string
}

// Create a new Input instance from inputHash string
func newInput(inputHash string) *Input {
	return &Input{
		inputHash: inputHash,
	}
}

// Marshal - A convenience method
// While we probably don't want any extra fields in the input
// this method exists in case we want to extend the data type
// at some point later on
func (t *Input) Marshal() ([]byte, error) {
	return []byte(t.inputHash), nil
}

// Output - Corresponds to Trillian Log ExtraData
type Output struct {
	outputHash string
}

// Create a new output struct from hash
// This is a model we may actually want to extend with extra metadata
// at some point
//
// Something that might be useful is to separate the hashing algo
// into a separate field, but for now this is just the NarHash
func newOutput(outputHash string) *Output {
	return &Output{
		outputHash: outputHash,
	}
}

// Marshal - Convenience method in case model is ever more complex
func (e *Output) Marshal() ([]byte, error) {
	return []byte(e.outputHash), nil
}
