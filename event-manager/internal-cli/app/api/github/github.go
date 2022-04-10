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

package github

import (
	"strconv"

	"github.com/compuzest/zlifecycle-internal-cli/app/util"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func GetAppInstallationID(
	log *logrus.Entry,
	client API,
	org string,
) (installationID *int64, appID *int64, err error) {
	log.WithFields(logrus.Fields{
		"org": org,
	}).Infof("Finding GitHub App installation ID for organization %s", org)
	installation, resp, err := client.FindOrganizationInstallation(org)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "error finding repository installation ID for org %s", org)
	}
	defer util.CloseBody(resp.Body)

	if resp.StatusCode != 200 {
		return nil, nil, errors.Errorf("find repository installation returned non-OK status: %d", resp.StatusCode)
	}

	log.WithFields(logrus.Fields{
		"org":            org,
		"installationId": strconv.FormatInt(*installation.ID, 10),
		"appId":          strconv.FormatInt(*installation.AppID, 10),
	}).Infof("Found installation ID for organization %s", org)

	return installation.ID, installation.AppID, nil
}

func GenerateInstallationToken(log *logrus.Entry, client API, org string) (token string, err error) {
	installationID, _, err := GetAppInstallationID(log, client, org)
	if err != nil {
		return "", errors.Wrapf(err, "error getting installation id for org [%s]", org)
	}
	log.WithFields(logrus.Fields{
		"org": org,
	}).Infof("Creating installation token for organization %s", org)
	installationToken, resp, err := client.CreateInstallationToken(*installationID)
	if err != nil {
		return "", errors.Wrapf(
			err,
			"error creating installation token for org [%s] and installation ID [%s]",
			strconv.FormatInt(*installationID, 10), org,
		)
	}
	defer util.CloseBody(resp.Body)

	if resp.StatusCode != 201 {
		return "", errors.Errorf("create installation token returned non-Created status: %d", resp.StatusCode)
	}

	return installationToken.GetToken(), nil
}
