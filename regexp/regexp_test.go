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
