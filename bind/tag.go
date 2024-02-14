package bind

import (
	"strings"
)

type tag struct {
	Name     string
	Default  string
	Required bool
	CommaSep bool
}

var TagPartSep = ","

func parseTag(tagValue string) tag {
	tg := tag{}

	parts := strings.Split(tagValue, TagPartSep)

	tg.Name = parts[0]

	for _, part := range parts[1:] {
		switch part {
		case "required", "require", "req":
			tg.Required = true
		case "comma":
			tg.CommaSep = true
		default:
			kv := strings.SplitN(part, "=", 2)
			k := kv[0]

			switch k {
			case "default":
				if len(kv) != 0 {
					tg.Default = kv[1]
				}
			}
		}
	}

	return tg
}
