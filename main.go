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

package main

import (
	"context"
	"flag"

	"github.com/compuzest/zlifecycle-il-operator/controllers/apm"
	"github.com/compuzest/zlifecycle-il-operator/controllers/gitreconciler"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/log"
	"github.com/newrelic/go-agent/v3/integrations/logcontext/nrlogrusplugin"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/cache"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers"
	// +kubebuilder:scaffold:imports
)

var Version = "0.0.1"

var (
	scheme   = runtime.NewScheme()
	setupLog *logrus.Entry
)

// +kubebuilder:rbac:groups="",resources=configmaps;secrets,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;create;update
func init() { // nolint
	l := logrus.New()
	l.SetFormatter(nrlogrusplugin.ContextFormatter{})
	setupLog = l.WithField("logger", "setup")

	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(stablev1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	mode := env.Config.Mode
	ctrl.SetLogger(zap.New(zap.UseDevMode(mode == "local")))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "ce9255a7.compuzest.com",
		CertDir:            env.Config.KubernetesCertDir,
		NewCache:           cache.MultiNamespacedCacheBuilder(getWatchedNamespaces()),
	})
	if err != nil {
		setupLog.WithError(err).Panic("unable to start manager")
	}

	setupLog.Info("Running zLifecycle IL operator version " + Version)

	// ctx
	ctx := context.Background()

	// Git reconciler
	gitReconciler := gitreconciler.NewReconciler(
		ctx,
		log.NewLogger().WithFields(logrus.Fields{"logger": "GitReconciler", "instance": env.Config.CompanyName, "company": env.Config.CompanyName}),
		mgr.GetClient(),
	)
	if err := gitReconciler.Start(); err != nil {
		setupLog.WithError(err).Error(err, "failed to start git reconciler")
	}

	// new relic
	var _apm apm.APM

	if env.Config.EnableNewRelic == "true" {
		setupLog.Info("Initializing NewRelic APM")
		logrus.SetFormatter(nrlogrusplugin.ContextFormatter{})
		_apm, err = apm.NewNewRelic()
		if err != nil {
			setupLog.WithError(err).Panic("unable to init new relic")
		}
	} else {
		_apm = apm.NewNoop()
		setupLog.Info("Defaulting to no-op APM")
	}

	// company controller init
	if err = (&controllers.CompanyReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Company"),
		LogV2:  log.NewLogger().WithFields(logrus.Fields{"logger": "controller.Company", "instance": env.Config.CompanyName, "company": env.Config.CompanyName}),
		Scheme: mgr.GetScheme(),
		APM:    _apm,
	}).SetupWithManager(mgr); err != nil {
		setupLog.WithError(err).WithField("controller", "Company").Panic("unable to create controller")
	}

	// team controller init
	if err = (&controllers.TeamReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Team"),
		LogV2:  log.NewLogger().WithFields(logrus.Fields{"logger": "controller.Team", "instance": env.Config.CompanyName, "company": env.Config.CompanyName}),
		Scheme: mgr.GetScheme(),
		APM:    _apm,
	}).SetupWithManager(mgr); err != nil {
		setupLog.WithError(err).WithField("controller", "Team").Panic("unable to create controller")
	}

	// environment controller init
	if err = (&controllers.EnvironmentReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Environment"),
		LogV2:  log.NewLogger().WithFields(logrus.Fields{"logger": "controller.Environment", "instance": env.Config.CompanyName, "company": env.Config.CompanyName}),
		Scheme: mgr.GetScheme(),
		APM:    _apm,
	}).SetupWithManager(mgr); err != nil {
		setupLog.WithError(err).WithField("controller", "Environment").Panic("unable to create controller")
	}

	if env.Config.KubernetesDisableWebhooks != "true" {
		setupLog.Info("Initializing webhook service")
		if err = (&stablev1.Environment{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.WithError(err).Panic("unable to create webhook", "webhook", "Environment")
		}
	}

	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.WithError(err).Panic("problem running manager")
	}
}

func getWatchedNamespaces() []string {
	namespaces := make([]string, 0, 2)
	systemNamespace := env.SystemNamespace()
	configNamespace := env.ConfigNamespace()
	namespaces = append(namespaces, systemNamespace)
	if systemNamespace != configNamespace {
		namespaces = append(namespaces, configNamespace)
	}
	return namespaces
}
