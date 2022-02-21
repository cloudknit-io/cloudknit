/* Copyright (C) 2020 CompuZest, Inc. - All Rights Reserved
 *
 * Unauthorized copying of this file, via any medium, is strictly prohibited
 * Proprietary and confidential
 *
 * NOTICE: All information contained herein is, and remains the property of
 * CompuZest, Inc. The intellectual and technical concepts contained herein are
 * proprietary to CompuZest, Inc. and are protected by trade secret or copyright
 * law. Dissemination of this information or reproduction of this material is
 * strictly forbidden unless prior written permission is obtained from CompuZest, Inc.
 */

package argocd

import (
	"github.com/compuzest/zlifecycle-il-operator/controllers/env"
	"github.com/pkg/errors"
)

func getArgocdCredentialsFromEnv() (*Credentials, error) {
	username := env.Config.ArgocdUsername
	password := env.Config.ArgocdPassword
	if username == "" || password == "" {
		return nil, errors.New("missing 'ARGOCD_USERNAME' or 'ARGOCD_PASSWORD' env variables")
	}

	creds := Credentials{Username: username, Password: password}

	return &creds, nil
}
