package regexp_test

import (
	. "github.com/proskenion/proskenion/regexp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRegexp(t *testing.T) {
	regexp := GetRegexp()
	assert.True(t, regexp.VerifyAccountId.MatchString("account@domain.com"))

	regexp = GetRegexp()
	assert.True(t, regexp.VerifyPeerAddress.MatchString("peer:50051"))
}

func TestVerifyAccountId(t *testing.T) {
	for _, c := range []struct {
		name      string
		accountId string
		ok        bool
	}{
		{
			"case 1 correct",
			"acccount@domain.com",
			true,
		},
		{
			"case 2 correct",
			"account@domain.com.a.b.c",
			true,
		},
		{
			"case 3 over domain",
			"account@domain.com.a.b.c.d",
			false,
		},
		{
			"case 4 a little account",
			"@",
			false,
		},
		{
			"case 5 no domain",
			"account@",
			false,
		},
		{
			"case 6 no accountt",
			"@domain.com",
			false,
		},
		{
			"case 7 over account",
			"adfsfsfsdfsadfdsfaasdfdsfsdfsdfsdfdsfsdffsfsdfsdf@a",
			false,
		},
		{
			"case 8 ng char",
			"Str@dcom",
			false,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			if c.ok {
				assert.True(t, GetRegexp().VerifyAccountId.MatchString(c.accountId))
			} else {
				assert.False(t, GetRegexp().VerifyAccountId.MatchString(c.accountId))
			}
		})
	}
}

func TestNewVerifyDomainId(t *testing.T) {
	for _, c := range []struct {
		name     string
		domainId string
		ok       bool
	}{
		{
			"case 1 correct",
			"domain.com",
			true,
		},
		{
			"case 2 correct",
			"domain.com.a.b.c",
			true,
		},
		{
			"case 3 correct",
			"domain",
			true,
		},
		{
			"case 4 over domain",
			"domain.com.a.b.c.d",
			false,
		},
		{
			"case 5 a little account",
			"@",
			false,
		},
		{
			"case 6 no domain",
			"account@",
			false,
		},
		{
			"case 7 accountId",
			"account@domain.com",
			false,
		},
		{
			"case 8 over domain",
			"adfsfsfsdfsadfdsfaasdfdsfsdfsdfsdfdsfsdffsfsdfsdf",
			false,
		},
		{
			"case 9 ng char",
			"Dcom",
			false,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			if c.ok {
				assert.True(t, GetRegexp().VerifyDomainId.MatchString(c.domainId))
			} else {
				assert.False(t, GetRegexp().VerifyDomainId.MatchString(c.domainId))
			}
		})
	}
}

func TestNewVerifyStorageId(t *testing.T) {
	for _, c := range []struct {
		name      string
		storageId string
		ok        bool
	}{
		{
			"case 1 correct",
			"domain.com/storage",
			true,
		},
		{
			"case 2 correct",
			"domain.com.a.b.c/storage",
			true,
		},
		{
			"case 3 over domain",
			"domain.com.a.b.c.d/storage",
			false,
		},
		{
			"case 4 no domain",
			"/storage",
			false,
		},
		{
			"case 5 no storage",
			"/",
			false,
		},
		{
			"case 6 over domain",
			"adfsfsfsdfsadfdsfaasdfdsfsdfsdfsdfdsfsdffsfsdfsdf/a",
			false,
		},
		{
			"case 6 over storage",
			"a/adfsfsfsdfsadfdsfaasdfdsfsdfsdfsdfdsfsdffsfsdfsdf",
			false,
		},
		{
			"case 8 ng char",
			"dcom/Storage",
			false,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			if c.ok {
				assert.True(t, GetRegexp().VerifyStorageId.MatchString(c.storageId))
			} else {
				assert.False(t, GetRegexp().VerifyStorageId.MatchString(c.storageId))
			}
		})
	}
}

func TestNewVerifyWalletId(t *testing.T) {
	for _, c := range []struct {
		name     string
		walletId string
		ok       bool
	}{
		{
			"case 1 correct",
			"accoun@domain.com/storage",
			true,
		},
		{
			"case 2 correct",
			"account@domain.com.a.b.c/storage",
			true,
		},
		{
			"case 3 over domain",
			"a@domain.com.a.b.c.d/storage",
			false,
		},
		{
			"case 4 no domain",
			"a@/storage",
			false,
		},
		{
			"case 5 no storage",
			"a@domain/",
			false,
		},
		{
			"case 6 over domain",
			"a@adfsfsfsdfsadfdsfaasdfdsfsdfsdfsdfdsfsdffsfsdfsdf/a",
			false,
		},
		{
			"case 6 over storage",
			"a@a/adfsfsfsdfsadfdsfaasdfdsfsdfsdfsdfdsfsdffsfsdfsdf",
			false,
		},
		{
			"case 8 ng char",
			"a@dcom/Storage",
			false,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			if c.ok {
				assert.True(t, GetRegexp().VerifyWalletId.MatchString(c.walletId))
			} else {
				assert.False(t, GetRegexp().VerifyWalletId.MatchString(c.walletId))
			}
		})
	}
}

func TestNewSplitAddress(t *testing.T) {
	for _, c := range []struct {
		name    string
		id      string
		storage string
		domain  string
		account string
	}{
		{
			"case 1 all",
			"account@domain.com/storage",
			"storage",
			"domain.com",
			"account",
		},
		{
			"case 2 correct",
			"account@domain.com.a.b.c/storage",
			"storage",
			"domain.com.a.b.c",
			"account",
		},
		{
			"case 3 correct",
			"a@domain/storage",
			"storage",
			"domain",
			"a",
		},
		{
			"case 4 no account",
			"domain/storage",
			"storage",
			"domain",
			"",
		},
		{
			"case 5 no storage",
			"a@domain",
			"",
			"domain",
			"a",
		},
		{
			"case 6 only domain",
			"domain.com.a.b",
			"",
			"domain.com.a.b",
			"",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			res := GetRegexp().SplitAddress.FindStringSubmatch(c.id)
			assert.Equal(t, 4, len(res))
			assert.Equal(t, c.id, res[0])
			assert.Equal(t, c.account, res[1])
			assert.Equal(t, c.domain, res[2])
			assert.Equal(t, c.storage, res[3])
		})
	}
}
