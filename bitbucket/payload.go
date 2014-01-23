package bitbucket

type Payload struct {
}

func NewPayload(data []byte) Payload {
	var p Payload
	return p
}

func (p *Payload) Url() string {
	// TODO: return p.Foo.Url
	return ""
}
