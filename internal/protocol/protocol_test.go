package protocol

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseSetBytes(t *testing.T) {
	respValid := &ResponseSet{
		Status: StatusOK,
	}

	b, err := respValid.Bytes()
	assert.NoError(t, err)

	expected := []byte{0x1} // StatusOK = 1
	assert.Equal(t, expected, b)
}

func TestResponseSetParse(t *testing.T) {
	resp := &ResponseSet{
		Status: StatusOK,
	}

	b, err := resp.Bytes()
	assert.NoError(t, err)

	r := bytes.NewReader(b)

	presp, err := ParseSetResponse(r)
	assert.NoError(t, err)

	assert.Equal(t, resp, presp)
}

func TestResponseGetBytes(t *testing.T) {
	respValid := &ResponseGet{
		Status: StatusOK,
		Value:  []byte("Bar"),
	}

	b, err := respValid.Bytes()
	assert.NoError(t, err)

	expected := []byte{0x1, 0x3, 0x0, 0x0, 0x0, 0x42, 0x61, 0x72}

	assert.Equal(t, expected, b)
}

func TestResponseGetParse(t *testing.T) {
	resp := &ResponseGet{
		Status: StatusOK,
		Value:  []byte("Bar"),
	}

	b, err := resp.Bytes()
	assert.NoError(t, err)

	r := bytes.NewReader(b)

	presp, err := ParseGetResponse(r)
	assert.NoError(t, err)

	assert.Equal(t, resp, presp)
}

func TestResponseDeleteBytes(t *testing.T) {
	respValid := &ResponseDelete{
		Status: StatusOK,
	}

	b, err := respValid.Bytes()
	assert.NoError(t, err)

	expected := []byte{0x01} // StatusOK = 1
	assert.Equal(t, expected, b)
}

func TestResponseDeleteParse(t *testing.T) {
	resp := &ResponseDelete{
		Status: StatusOK,
	}

	b, err := resp.Bytes()
	assert.NoError(t, err)

	r := bytes.NewReader(b)

	presp, err := ParseDeleteResponse(r)
	assert.NoError(t, err)

	assert.Equal(t, resp, presp)
}

func TestCommandGetBytes(t *testing.T) {
	cmdValid := &CommandGet{
		Key: []byte("Foo"),
	}

	b, err := cmdValid.Bytes()
	assert.NoError(t, err)

	expected := []byte{
		0x02,                   // CmdGet = 2
		0x03, 0x00, 0x00, 0x00, // length of "Foo" in little-endian format
		0x46, 0x6F, 0x6F, // ASCII values for "Foo"
	}
	assert.Equal(t, expected, b)
}

func TestCommandGetParse(t *testing.T) {
	cmd := &CommandGet{
		Key: []byte("Foo"),
	}

	b, err := cmd.Bytes()
	assert.NoError(t, err)

	r := bytes.NewReader(b)

	pcmd, err := ParseCommand(r)
	assert.NoError(t, err)

	pcmdGet, ok := pcmd.(*CommandGet)
	assert.True(t, ok)

	assert.Equal(t, cmd, pcmdGet)
}

func TestCommandDeleteBytes(t *testing.T) {
	cmdValid := &CommandDelete{
		Key: []byte("Foo"),
	}

	b, err := cmdValid.Bytes()
	assert.NoError(t, err)

	expected := []byte{
		0x03,                   // CmdDel = 3
		0x03, 0x00, 0x00, 0x00, // length of "Foo" in little-endian format
		0x46, 0x6F, 0x6F, // ASCII values for "Foo"
	}
	assert.Equal(t, expected, b)
}

func TestCommandDeleteParse(t *testing.T) {
	cmd := &CommandDelete{
		Key: []byte("Foo"),
	}

	b, err := cmd.Bytes()
	assert.Nil(t, err)

	r := bytes.NewReader(b)

	pcmd, err := ParseCommand(r)
	assert.Nil(t, err)

	pcmdDel, ok := pcmd.(*CommandDelete)
	assert.True(t, ok)

	assert.Equal(t, cmd, pcmdDel)
}
