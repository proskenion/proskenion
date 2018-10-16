package regexp

import "regexp"

type Regexper struct {
	VerifyAccountId   *regexp.Regexp
	VerifyPeerAddress *regexp.Regexp
}

var instance *Regexper

func GetRegexp() *Regexper {
	if instance == nil {
		instance = &Regexper{
			VerifyAccountId:   newVerifyAccountId(),
			VerifyPeerAddress: newVerifyPeerAddress(),
		}
	}
	return instance
}

func newVerifyAccountId() *regexp.Regexp {
	return regexp.MustCompile(`^[a-z_0-9]{1,32}\@[a-z_0-9]{1,32}(\.[a-z_0-9]{1,32}){0,4}$`)
}

func newVerifyPeerAddress() *regexp.Regexp {
	return regexp.MustCompile(`[a-z_0-9]{1,32}:\d{5}`)
}
