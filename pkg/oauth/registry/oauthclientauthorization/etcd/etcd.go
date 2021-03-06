package etcd

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	kapi "k8s.io/kubernetes/pkg/api"

	"github.com/openshift/origin/pkg/oauth/api"
	"github.com/openshift/origin/pkg/oauth/registry/oauthclient"
	"github.com/openshift/origin/pkg/oauth/registry/oauthclientauthorization"
	"github.com/openshift/origin/pkg/util/restoptions"
)

// rest implements a RESTStorage for oauth client authorizations against etcd
type REST struct {
	*registry.Store
}

// NewREST returns a RESTStorage object that will work against oauth clients
func NewREST(optsGetter restoptions.Getter, clientGetter oauthclient.Getter) (*REST, error) {
	strategy := oauthclientauthorization.NewStrategy(clientGetter)

	store := &registry.Store{
		Copier:            kapi.Scheme,
		NewFunc:           func() runtime.Object { return &api.OAuthClientAuthorization{} },
		NewListFunc:       func() runtime.Object { return &api.OAuthClientAuthorizationList{} },
		PredicateFunc:     oauthclientauthorization.Matcher,
		QualifiedResource: api.Resource("oauthclientauthorizations"),

		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,
	}

	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: oauthclientauthorization.GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}

	return &REST{store}, nil
}
