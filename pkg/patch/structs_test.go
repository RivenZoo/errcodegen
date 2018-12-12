package patch

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCoverStructsField(t *testing.T) {
	type A struct {
		Name string
		ID   int
		Desc *string
		Data []byte
	}
	desc := "abc"
	src := A{
		Name: "test",
		ID:   0,
		Desc: &desc,
		Data: []byte("data"),
	}
	dst := A{}
	err := CoverStructsField(src, &dst)
	assert.Nil(t, err)

	assert.EqualValues(t, src, dst)

	src = A{
		Name: "test",
		ID:   0,
	}
	dst = A{
		Name: "test2",
		ID:   2,
		Desc: &desc,
	}
	err = CoverStructsField(src, &dst)
	assert.Nil(t, err)
	assert.EqualValues(t, src.Name, dst.Name)
	assert.EqualValues(t, 2, dst.ID)
	assert.EqualValues(t, &desc, dst.Desc)

	err = CoverStructsField(src, dst)
	assert.NotNil(t, err)

	err = CoverStructsField(struct {
		ID string
	}{ID: "1"}, &dst)
	assert.NotNil(t, err)
}
