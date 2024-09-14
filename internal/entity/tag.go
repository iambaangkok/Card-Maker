package entity

type Tag struct {
	Name string
}

var Tags = []string{"Kinetic", "Mag-Fed"}

var TagNameMap = getTagNameMap()

func getTagNameMap() map[string]Tag {
	m := map[string]Tag{}
	for _, n := range Tags {
		m[n] = Tag{
			Name: n,
		}
	}
	return m
}