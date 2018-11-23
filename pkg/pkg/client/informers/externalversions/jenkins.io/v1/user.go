// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	time "time"

	jenkinsiov1 "github.com/xcloudnative/xcloud/pkg/pkg/apis/jenkins.io/v1"
	versioned "github.com/xcloudnative/xcloud/pkg/pkg/client/clientset/versioned"
	internalinterfaces "github.com/xcloudnative/xcloud/pkg/pkg/clientient/informers/externalversions/internalinterfaces"
	v1 "github.com/xcloudnative/xcloud/pkg/pkg/clientient/listers/jenkins.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// UserInformer provides access to a shared informer and lister for
// Users.
type UserInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.UserLister
}

type userInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewUserInformer constructs a new informer for User type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewUserInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredUserInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredUserInformer constructs a new informer for User type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredUserInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.JenkinsV1().Users(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.JenkinsV1().Users(namespace).Watch(options)
			},
		},
		&jenkinsiov1.User{},
		resyncPeriod,
		indexers,
	)
}

func (f *userInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredUserInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *userInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&jenkinsiov1.User{}, f.defaultInformer)
}

func (f *userInformer) Lister() v1.UserLister {
	return v1.NewUserLister(f.Informer().GetIndexer())
}