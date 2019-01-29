package query

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/config"
	core "github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

type QueryValidator struct {
	fc   model.ModelFactory
	conf *config.Config
}

func NewQueryValidator(fc model.ModelFactory, conf *config.Config) core.QueryValidator {
	return &QueryValidator{ fc, conf}
}

func containsPublicKey(keys []model.PublicKey, pub model.PublicKey) bool {
	for _, key := range keys {
		if bytes.Equal(key, pub) {
			return true
		}
	}
	return false
}

func (q *QueryValidator) Validate(wsv model.ObjectFinder, query model.Query) error {
	// 署名チェック
	ac := q.fc.NewEmptyAccount()
	authorizer := model.MustAddress(model.MustAddress(query.GetPayload().GetAuthorizerId()).AccountId())
	err := wsv.Query(authorizer, ac)
	if err != nil {
		return errors.Wrapf(core.ErrQueryProcessorNotExistAuthoirizer,
			"authorizer : %s", query.GetPayload().GetAuthorizerId())
	}
	if !containsPublicKey(ac.GetPublicKeys(), query.GetSignature().GetPublicKey()) {
		return errors.Wrapf(core.ErrQueryProcessorNotSignedAuthorizer,
			"authorizer : %s, expect key : %x",
			query.GetPayload().GetAuthorizerId(), query.GetSignature().GetPublicKey())
	}
	return nil
}
