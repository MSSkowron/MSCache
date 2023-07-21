package protocol

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytes(t *testing.T) {
	cmdValid := &CommandSet{
		Key:   []byte("Foo"),
		Value: []byte("Bar"),
		TTL:   2,
	}

	_, err := cmdValid.Bytes()
	assert.Nil(t, err)
}

func TestParseSetCommand(t *testing.T) {
	cmd := &CommandSet{
		Key:   []byte("Foo"),
		Value: []byte("Bar"),
		TTL:   2,
	}

	b, err := cmd.Bytes()
	assert.Nil(t, err)

	r := bytes.NewReader(b)

	pcmd, err := ParseCommand(r)
	assert.Nil(t, err)

	assert.Equal(t, cmd, pcmd)
}

func TestParseGetCommand(t *testing.T) {
	cmd := &CommandGet{
		Key: []byte("Foo"),
	}

	b, err := cmd.Bytes()
	assert.Nil(t, err)

	r := bytes.NewReader(b)

	pcmd, err := ParseCommand(r)
	assert.Nil(t, err)

	assert.Equal(t, cmd, pcmd)
}

func TestParseDeleteCommand(t *testing.T) {
	cmd := &CommandDelete{
		Key: []byte("Foo"),
	}

	b, err := cmd.Bytes()
	assert.Nil(t, err)

	r := bytes.NewReader(b)

	pcmd, err := ParseCommand(r)
	assert.Nil(t, err)

	assert.Equal(t, cmd, pcmd)
}
