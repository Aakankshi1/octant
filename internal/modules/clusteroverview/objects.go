/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package clusteroverview

import (
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"

	"github.com/vmware-tanzu/octant/internal/describer"
	"github.com/vmware-tanzu/octant/pkg/icon"
	"github.com/vmware-tanzu/octant/pkg/store"
)

var (
	customResourcesDescriber = describer.NewCRDSection(
		"/custom-resources",
		"Custom Resources",
		describer.ResourceLink{Title: "Cluster Overview", Url: "/cluster-overview"},
	)

	rbacClusterRoles = describer.NewResource(describer.ResourceOptions{
		Path:           "/rbac/cluster-roles",
		ObjectStoreKey: store.Key{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "ClusterRole"},
		ListType:       &rbacv1.ClusterRoleList{},
		ObjectType:     &rbacv1.ClusterRole{},
		Titles:         describer.ResourceTitle{List: "Cluster Roles", Object: "Cluster Roles"},
		ClusterWide:    true,
		IconName:       icon.ClusterOverviewClusterRole,
		RootPath:       describer.ResourceLink{Title: "Cluster Overview", Url: "/cluster-overview"},
	})

	rbacClusterRoleBindings = describer.NewResource(describer.ResourceOptions{
		Path:           "/rbac/cluster-role-bindings",
		ObjectStoreKey: store.Key{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "ClusterRoleBinding"},
		ListType:       &rbacv1.ClusterRoleBindingList{},
		ObjectType:     &rbacv1.ClusterRoleBinding{},
		Titles:         describer.ResourceTitle{List: "Cluster Role Bindings", Object: "Cluster Role Bindings"},
		ClusterWide:    true,
		IconName:       icon.ClusterOverviewClusterRoleBinding,
		RootPath:       describer.ResourceLink{Title: "Cluster Overview", Url: "/cluster-overview"},
	})

	rbacDescriber = describer.NewSection(
		"/rbac",
		"RBAC",
		rbacClusterRoles,
		rbacClusterRoleBindings,
	)

	nodesDescriber = describer.NewResource(describer.ResourceOptions{
		Path:                  "/nodes",
		ObjectStoreKey:        store.Key{APIVersion: "v1", Kind: "Node"},
		ListType:              &v1.NodeList{},
		ObjectType:            &v1.Node{},
		Titles:                describer.ResourceTitle{List: "Nodes", Object: "Nodes"},
		DisableResourceViewer: true,
		ClusterWide:           true,
		IconName:              icon.ClusterOverviewNode,
		RootPath:       describer.ResourceLink{Title: "Cluster Overview", Url: "/cluster-overview"},
	})

	storagePersistentVolumeDescriber = describer.NewResource(describer.ResourceOptions{
		Path:           "/storage/persistent-volumes",
		ObjectStoreKey: store.Key{APIVersion: "v1", Kind: "PersistentVolume"},
		ListType:       &v1.PersistentVolumeList{},
		ObjectType:     &v1.PersistentVolume{},
		Titles:         describer.ResourceTitle{List: "Persistent Volumes", Object: "Persistent Volumes"},
		ClusterWide:    true,
		IconName:       icon.ClusterOverviewPersistentVolume,
		RootPath:       describer.ResourceLink{Title: "Cluster Overview", Url: "/cluster-overview"},
	})

	storageDescriber = describer.NewSection(
		"/storage",
		"Storage",
		storagePersistentVolumeDescriber,
	)

	namespacesDescriber = describer.NewResource(describer.ResourceOptions{
		Path:                  "/namespaces",
		ObjectStoreKey:        store.Key{APIVersion: "v1", Kind: "Namespace"},
		ListType:              &v1.NamespaceList{},
		ObjectType:            &v1.Namespace{},
		Titles:                describer.ResourceTitle{List: "Namespaces", Object: "Namespaces"},
		DisableResourceViewer: true,
		ClusterWide:           true,
		IconName:              icon.ClusterOverviewNamespace,
		RootPath:       		describer.ResourceLink{Title: "Cluster Overview", Url: "/cluster-overview"},
	})

	portForwardDescriber = NewPortForwardListDescriber()

	rootDescriber = describer.NewSection(
		"/",
		"Cluster Overview",
		namespacesDescriber,
		customResourcesDescriber,
		rbacDescriber,
		nodesDescriber,
		storageDescriber,
		portForwardDescriber,
	)
)
