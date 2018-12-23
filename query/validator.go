package query

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	. "github.com/proskenion/proskenion/regexp"
)

type QueryValidator struct {}

func NewQueryValidator() core.QueryValidator {
	return &QueryValidator{}
}

var (
	ErrQueryValidateAccountTargetIdNotAccountId = fmt.Errorf("Failed Query Validate targetId is not accountId when get account object")
	ErrQueryValidatePeerTargetIdNotPeerAddress  = fmt.Errorf("Failed Query Validate targetId is not PeerAddress when get peer object")
	ErrQueryValidateAuthorizerIdNotAccountId    = fmt.Errorf("Failed Query Validate authorizerId is not accountId")
)

func (q *QueryValidator) Validate(query model.Query) error {
	qp := query.GetPayload()

	if ok := GetRegexp().VerifyAccountId.MatchString(qp.GetAuthorizerId()); !ok {
		return errors.Wrapf(ErrQueryValidateAuthorizerIdNotAccountId,
			"authorizerId : %s, must be : %s", qp.GetAuthorizerId(), GetRegexp().VerifyAccountId.String())
	}

	switch qp.GetRequestCode() {
	case model.AccountObjectCode:
		return q.accountObjectValidate(qp)
	case model.PeerObjectCode:
		return q.peerObjectValidate(qp)
	}
	return nil
}

// TODO これ Verify
func (q *QueryValidator) accountObjectValidate(qp model.QueryPayload) error {
	if ok := GetRegexp().VerifyAccountId.MatchString(qp.GetFromId()); !ok {
		return errors.Wrapf(ErrQueryValidateAccountTargetIdNotAccountId,
			"targetId : %s, must be : %s", qp.GetFromId(), GetRegexp().VerifyAccountId.String())
	}
	return nil
}

func (q *QueryValidator) peerObjectValidate(qp model.QueryPayload) error {
	if ok := GetRegexp().VerifyPeerAddress.MatchString(qp.GetFromId()); !ok {
		return errors.Wrapf(ErrQueryValidatePeerTargetIdNotPeerAddress,
			"targetid : %s, must be : %s", qp.GetFromId(), GetRegexp().VerifyPeerAddress.String())
	}
	return nil
}
