package main

type Channel struct {
	Name       string `json:name`
	Url        string `json:url`
	PlayPrefix string `json:play_prefix`
}

var (
	channels []Channel
)

func FindInChannels(name string) (Channel, bool) {
	channel := Channel{}
	ok := false
	for _, v := range channels {
		if name == v.Name {
			ok = true
			channel = v

			break
		}
	}

	return channel, ok
}
