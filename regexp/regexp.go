package regexp

import "regexp"

type Regexper struct {
	VerifyAccountId   *regexp.Regexp
	VerifyPeerAddress *regexp.Regexp
	VerifyDomainId    *regexp.Regexp
	VerifyStorageId   *regexp.Regexp
	VerifyWalletId    *regexp.Regexp
	SplitAddress      *regexp.Regexp
}

var instance *Regexper

func GetRegexp() *Regexper {
	if instance == nil {
		instance = &Regexper{
			VerifyAccountId:   newVerifyAccountId(),
			VerifyPeerAddress: newVerifyPeerAddress(),
			VerifyDomainId:    newVerifyDomainId(),
			VerifyStorageId:   newVerifyStorageId(),
			VerifyWalletId:    newVerifyWalletId(),
			SplitAddress:      newSplitAddress(),
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

////
// accountId = "account@domain.com"
// domainId = "domain.com"
// storageId = "domain.com/storage"
// walletId = "account@domain.com/storage"
//

func newVerifyDomainId() *regexp.Regexp {
	return regexp.MustCompile(`^[a-z_0-9]{1,32}(\.[a-z_0-9]{1,32}){0,4}$`)
}

func newVerifyStorageId() *regexp.Regexp {
	return regexp.MustCompile(`^[a-z_0-9]{1,32}(\.[a-z_0-9]{1,32}){0,4}/[a-z_0-9]{1,32}$`)
}

func newVerifyWalletId() *regexp.Regexp {
	return regexp.MustCompile(`^[a-z_0-9]{1,32}\@[a-z_0-9]{1,32}(\.[a-z_0-9]{1,32}){0,4}/[a-z_0-9]{1,32}$`)
}

func newSplitAddress() *regexp.Regexp {
	return regexp.MustCompile(`^(?:([a-z_0-9]{1,32})\@)?([a-z_0-9]{1,32}(?:\.[a-z_0-9]{1,32}){0,4})(?:/([a-z_0-9]{1,32}))?$`)
}
