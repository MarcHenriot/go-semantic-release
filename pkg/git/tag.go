package git

import "github.com/go-git/go-git/v5/plumbing"

type Tag struct {
	hash plumbing.Hash
	name string
	ref  string
}

type Tags []*Tag

func NewTag(name, ref string, hash plumbing.Hash) *Tag {
	return &Tag{
		name: name,
		ref:  ref,
		hash: hash,
	}
}

func (t *Tag) GetName() string {
	return t.name
}

func (t *Tag) GetRef() string {
	return t.ref
}

func (t *Tag) GetHash() plumbing.Hash {
	return t.hash
}
