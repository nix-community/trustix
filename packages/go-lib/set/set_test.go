package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	set := NewSet[string]()
	set.Add("123")

	assert.NotNil(t, set.values["123"])
}

func TestHas(t *testing.T) {
	set := NewSet[string]()
	set.Add("123")
	
	assert.Equal(t, set.Has("123"), true)
	assert.Equal(t, set.Has("1234"), false)
}

func TestRemove(t *testing.T) {
	set := NewSet[string]()

	set.Add("123")
	assert.Equal(t, set.Has("123"), true)

	set.Remove("123")
	assert.Equal(t, set.Has("123"), false)
}

func TestValues(t *testing.T) {
	set := NewSet[string]()

	// Add set members out of order to also test stable output order of Values()
	set.Add("321")
	set.Add("123")

	values := set.Values()
	
	assert.Equal(t, len(values), 2)
	assert.Equal(t, values[0], "123")
	assert.Equal(t, values[1], "321")
}

func TestUnion(t *testing.T) {
	a := NewSet[string]()
	a.Add("A")
	
	b := NewSet[string]()
	b.Add("B")

	u := a.Union(b)

	uv := u.Values()
	assert.Equal(t, len(uv), 2)

	assert.Equal(t, uv[0], "A")
	assert.Equal(t, uv[1], "B")
}

func TestDiff(t *testing.T) {
	a := NewSet[string]()
	a.Add("A")
	
	b := NewSet[string]()
	b.Add("A")
	b.Add("B")

	d := b.Diff(a)

	dv := d.Values()
	assert.Equal(t, len(dv), 1)

	assert.Equal(t, dv[0], "B")
}

func TestCopy(t *testing.T) {
	a := NewSet[string]()
	a.Add("A")
	a.Add("B")
	a.Add("C")

	b := a.Copy()

	assert.True(t, a != b) // Ensure different addr

	// Ensure no difference from copied set
	assert.Equal(t, len(a.Diff(b).Values()), 0)
	assert.Equal(t, len(b.Diff(a).Values()), 0)
}

func TestUpdate(t *testing.T) {
	a := NewSet[string]()
	a.Add("A")

	b := NewSet[string]()
	b.Add("B")
	b.Add("C")

	a.Update(b)

	values := a.Values()
	assert.Equal(t, len(values), 3)

	assert.Equal(t, values[0], "A")
	assert.Equal(t, values[1], "B")
	assert.Equal(t, values[2], "C")
}
