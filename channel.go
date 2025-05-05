package main

type Channel struct {
	Name       string `json:name`
	Desc       string `json:desc`
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

func FindInChannels2(property string, value string) (Channel, bool) {
	channel := Channel{}
	ok := false
	for _, v := range channels {
		bindValue := ""
		switch property {
		case "name":
			bindValue = v.Name
		case "desc":
			bindValue = v.Desc
		case "url":
			bindValue = v.Url
		}
		if bindValue != "" && bindValue == value {
			ok = true
			channel = v

			break
		}
	}

	return channel, ok
}
