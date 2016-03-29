package core

type Purview interface {
	Verify(req []byte, target string) bool
}

func (this *VMs) SetAuther(p Purview) {
	this.Auth = p
}
