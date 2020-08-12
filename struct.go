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

// Corresponds to Trillian Log ExtraData
type Output struct {
	name string
}

func newOutput(name string) *Output {
	return &Output{
		name: name,
	}
}

func (e *Output) Marshal() ([]byte, error) {
	return []byte(e.name), nil
}
