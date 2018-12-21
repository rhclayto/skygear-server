package oauth

import (
	"github.com/sirupsen/logrus"
	"github.com/skygeario/skygear-server/pkg/core/db"
)

type safeProviderImpl struct {
	impl      *providerImpl
	txContext db.SafeTxContext
}

func NewSafeProvider(
	builder db.SQLBuilder,
	executor db.SQLExecutor,
	logger *logrus.Entry,
	txContext db.SafeTxContext,
) Provider {
	return &safeProviderImpl{
		impl:      newProvider(builder, executor, logger),
		txContext: txContext,
	}
}

func (p *safeProviderImpl) GetPrincipalByProviderUserID(providerName string, providerUserID string) (*Principal, error) {
	p.txContext.EnsureTx()
	return p.impl.GetPrincipalByProviderUserID(providerName, providerUserID)
}

func (p *safeProviderImpl) GetPrincipalByUserID(providerName string, userID string) (*Principal, error) {
	p.txContext.EnsureTx()
	return p.impl.GetPrincipalByUserID(providerName, userID)
}

func (p *safeProviderImpl) CreatePrincipal(principal Principal) error {
	p.txContext.EnsureTx()
	return p.impl.CreatePrincipal(principal)
}

func (p *safeProviderImpl) UpdatePrincipal(principal *Principal) error {
	p.txContext.EnsureTx()
	return p.impl.UpdatePrincipal(principal)
}

func (p *safeProviderImpl) DeletePrincipal(principal *Principal) error {
	p.txContext.EnsureTx()
	return p.impl.DeletePrincipal(principal)
}

func (p *safeProviderImpl) GetPrincipalsByUserID(userID string) ([]*Principal, error) {
	p.txContext.EnsureTx()
	return p.impl.GetPrincipalsByUserID(userID)
}
