package git

type Tag struct {
	name      string
	ref       string
	commitSha string
}

type Tags []*Tag

// Constructor for Tag
func NewTag(name string, ref string, commitSha string) *Tag {
	return &Tag{
		name:      name,
		ref:       ref,
		commitSha: commitSha,
	}
}

// Getter methods
func (t *Tag) GetName() string {
	return t.name
}

func (t *Tag) GetRef() string {
	return t.ref
}

func (t *Tag) GetCommitSha() string {
	return t.commitSha
}
