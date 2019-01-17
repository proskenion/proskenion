package query

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	. "github.com/proskenion/proskenion/regexp"
)

type QueryVerifier struct{}

func NewQueryVerifier() core.QueryVerifier {
	return &QueryVerifier{}
}

var (
	ErrQueryVerifyAccountTargetIdNotAccountId = fmt.Errorf("Failed Query Verify targetId is not accountId when get account object")
	ErrQueryVerifyPeerTargetIdNotPeerAddress  = fmt.Errorf("Failed Query Verify targetId is not PeerAddress when get peer object")
	ErrQueryVerifyAuthorizerIdNotAccountId    = fmt.Errorf("Failed Query Verify authorizerId is not accountId")
	ErrQueryVerifyFromIdNotIdFormat           = fmt.Errorf("Failed Query Verify fromId is not valid format")
)

func (q *QueryVerifier) Verify(query model.Query) error {
	qp := query.GetPayload()

	if ok := GetRegexp().VerifyAccountId.MatchString(qp.GetAuthorizerId()); !ok {
		return errors.Wrapf(ErrQueryVerifyAuthorizerIdNotAccountId,
			"authorizerId : %s, must be : %s", qp.GetAuthorizerId(), GetRegexp().VerifyAccountId.String())
	}
	if _, err := model.NewAddress(qp.GetFromId()); err != nil {
		return errors.Wrapf(ErrQueryVerifyFromIdNotIdFormat,
			"fromId : %s, not invalid id format", qp.GetFromId())
	}

	/*
		switch qp.GetRequestCode() {
		case model.AccountObjectCode:
			return q.accountObjectVerify(qp)
		case model.PeerObjectCode:
			return q.peerObjectVerify(qp)
		}
	*/
	return nil
}

// TODO これ Verify
func (q *QueryVerifier) accountObjectVerify(qp model.QueryPayload) error {
	if ok := GetRegexp().VerifyAccountId.MatchString(qp.GetFromId()); !ok {
		return errors.Wrapf(ErrQueryVerifyAccountTargetIdNotAccountId,
			"targetId : %s, must be : %s", qp.GetFromId(), GetRegexp().VerifyAccountId.String())
	}
	return nil
}

func (q *QueryVerifier) peerObjectVerify(qp model.QueryPayload) error {
	if ok := GetRegexp().VerifyPeerAddress.MatchString(qp.GetFromId()); !ok {
		return errors.Wrapf(ErrQueryVerifyPeerTargetIdNotPeerAddress,
			"targetid : %s, must be : %s", qp.GetFromId(), GetRegexp().VerifyPeerAddress.String())
	}
	return nil
}
