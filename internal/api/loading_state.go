/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/vmware-tanzu/octant/internal/kubeconfig"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/action"
)

const (
	UploadKubeConfig = "action.octant.dev/uploadKubeConfig"
	CheckLoading     = "action.octant.dev/loading"
)

type LoadingManager struct {
	client     OctantClient
	kubeConfig kubeconfig.TemporaryKubeConfig
	ctx        context.Context
}

var _ StateManager = (*LoadingManager)(nil)

func NewLoadingManager(temporaryKubeConfig kubeconfig.TemporaryKubeConfig) *LoadingManager {
	lm := &LoadingManager{
		kubeConfig: temporaryKubeConfig,
	}

	return lm
}

func (l *LoadingManager) Handlers() []octant.ClientRequestHandler {
	return []octant.ClientRequestHandler{
		{
			RequestType: UploadKubeConfig,
			Handler:     l.UploadKubeConfig,
		},
		{
			RequestType: CheckLoading,
			Handler:     l.CheckLoading,
		},
	}
}

func (l *LoadingManager) Start(ctx context.Context, state octant.State, client OctantClient) {
	l.client = client
	l.ctx = ctx
}

func (l *LoadingManager) CheckLoading(state octant.State, payload action.Payload) error {
	loading, err := payload.Bool("loading")
	if err != nil {
		return fmt.Errorf("getting loading from payload: %w", err)
	}

	if loading {
		l.client.Send(octant.Event{
			Type: octant.EventTypeLoading,
		})
	}

	return nil
}

func (l *LoadingManager) UploadKubeConfig(state octant.State, payload action.Payload) error {
	kubeConfig64, err := payload.String("kubeConfig")
	if err != nil {
		return fmt.Errorf("getting kubeConfig from payload: %w", err)
	}

	kubeConfig, err := base64.StdEncoding.DecodeString(kubeConfig64)
	if err != nil {
		return fmt.Errorf("decode base64 kubeconfig: %w", err)
	}

	// Check if kube config can be converted to config object
	// TODO: Show error elsewhere instead of notifier
	if _, err := clientcmd.Load(kubeConfig); err != nil {
		message := fmt.Sprintf("Error parsing KubeConfig: %v", err)
		alert := action.CreateAlert(action.AlertTypeError, message, action.DefaultAlertExpiration)
		state.SendAlert(alert)
		return err
	}

	l.kubeConfig.KubeConfig <- string(kubeConfig)

	tempFile, err := ioutil.TempFile(os.TempDir(), "kubeconfig")
	if err != nil {
		return err
	}

	if _, err := tempFile.Write(kubeConfig); err != nil {
		return err
	}

	if err := tempFile.Close(); err != nil {
		return err
	}

	l.client.Send(octant.Event{
		Type: octant.EventTypeRefresh,
	})

	l.kubeConfig.Path <- tempFile.Name()
	return nil
}
