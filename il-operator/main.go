package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"

	"github.com/compuzest/zlifecycle-il-operator/controller/services/webhooks/mutating"
	"github.com/compuzest/zlifecycle-il-operator/controller/services/webhooks/validating"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/file"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/apm"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/cloudknitservice"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/eventservice"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/log"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/compuzest/zlifecycle-il-operator/controller/services/gitreconciler"
	"github.com/compuzest/zlifecycle-il-operator/controller/validator"

	"github.com/compuzest/zlifecycle-il-operator/controller/env"
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
	"github.com/compuzest/zlifecycle-il-operator/controller"
	// +kubebuilder:scaffold:imports
)

//go:embed .version
var Version string

var (
	scheme   = runtime.NewScheme()
	setupLog *logrus.Entry
)

//go:generate mockgen --build_flags=--mod=mod -destination=./mocks/mock_kclient.go -package=mocks "sigs.k8s.io/controller-runtime/pkg/client" Client

// +kubebuilder:rbac:groups="",resources=configmaps;secrets,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;create;update
// +kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions=get;list;watch
// nolint
func init() {
	l := logrus.New()
	l.SetFormatter(nrlogrusplugin.ContextFormatter{})
	setupLog = l.WithField("logger", "setup")

	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(apiextensionsv1.AddToScheme(scheme))
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
	gitReconciler, err := gitreconciler.NewReconciler(
		ctx,
		log.NewLogger().WithFields(logrus.Fields{"logger": "GitReconciler", "instance": env.Config.CompanyName, "company": env.Config.CompanyName, "version": Version}),
		mgr.GetClient(),
	)
	if err != nil {
		setupLog.WithError(err).Panic(err, "failed to instantiate")
	}
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
	if err = (&controller.CompanyReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controller").WithName("Company"),
		LogV2:  log.NewLogger().WithFields(logrus.Fields{"logger": "controller.Company", "instance": env.Config.CompanyName, "company": env.Config.CompanyName, "version": Version}),
		Scheme: mgr.GetScheme(),
		APM:    _apm,
	}).SetupWithManager(mgr); err != nil {
		setupLog.WithError(err).WithField("controller", "Company").Panic("unable to create controller")
	}

	// team controller init
	if err = (&controller.TeamReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controller").WithName("Team"),
		LogV2:  log.NewLogger().WithFields(logrus.Fields{"logger": "controller.Team", "instance": env.Config.CompanyName, "company": env.Config.CompanyName, "version": Version}),
		Scheme: mgr.GetScheme(),
		APM:    _apm,
	}).SetupWithManager(mgr); err != nil {
		setupLog.WithError(err).WithField("controller", "Team").Panic("unable to create controller")
	}

	// environment controller init
	if err = (&controller.EnvironmentReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controller").WithName("Environment"),
		LogV2: log.NewLogger().WithFields(
			logrus.Fields{
				"logger": "controller.Environment", "instance": env.Config.CompanyName, "company": env.Config.CompanyName, "version": Version,
			},
		),
		Scheme:        mgr.GetScheme(),
		APM:           _apm,
		GitReconciler: gitReconciler,
	}).SetupWithManager(mgr); err != nil {
		setupLog.WithError(err).WithField("controller", "Environment").Panic("unable to create controller")
	}

	if env.Config.KubernetesDisableWebhooks != "true" {
		registerWebhooks(mgr)
	}

	// +kubebuilder:scaffold:builder
	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.WithError(err).Panic("error running manager")
	}
}

func registerWebhooks(mgr manager.Manager) {
	environmentValidator := validator.NewEnvironmentValidatorImpl(mgr.GetClient(), file.NewOSFileService(), cloudknitservice.NewService(env.Config.ZLifecycleAPIURL))
	teamValidator := validator.NewTeamValidatorImpl(mgr.GetClient(), file.NewOSFileService(), eventservice.NewService(env.Config.ZLifecycleEventServiceURL))

	hs := mgr.GetWebhookServer()
	es := eventservice.NewService(env.Config.ZLifecycleEventServiceURL)
	hs.Register(
		fmt.Sprintf("/%s/validate-stable-compuzest-com-v1-environment", env.Config.CompanyName),
		&webhook.Admission{Handler: validating.NewEnvironmentValidatingWebhook(environmentValidator, es, setupLog.WithField("logger", "controllers.EnvironmentValidatingWebhook"))},
	)
	hs.Register(
		fmt.Sprintf("/%s/validate-stable-compuzest-com-v1-team", env.Config.CompanyName),
		&webhook.Admission{Handler: validating.NewTeamValidatingWebhook(teamValidator, es, setupLog.WithField("logger", "controllers.TeamValidatingWebhook"))},
	)
	hs.Register(
		fmt.Sprintf("/%s/mutate-stable-compuzest-com-v1-environment", env.Config.CompanyName),
		&webhook.Admission{Handler: mutating.NewEnvironmentMutatingWebhook(mgr.GetClient(), es, setupLog.WithField("logger", "controllers.EnvironmentMutatingWebhook"))},
	)
	hs.Register(
		fmt.Sprintf("/%s/mutate-stable-compuzest-com-v1-team", env.Config.CompanyName),
		&webhook.Admission{Handler: mutating.NewTeamMutatingWebhook(mgr.GetClient(), es, setupLog.WithField("logger", "controllers.TeamMutatingWebhook"))},
	)
}

func getWatchedNamespaces() []string {
	namespaces := make([]string, 0, 2)
	systemNamespace := env.SystemNamespace()
	configNamespace := env.ConfigNamespace()
	executorNamespace := env.ExecutorNamespace()
	namespaces = append(namespaces, systemNamespace)
	if systemNamespace != configNamespace {
		namespaces = append(namespaces, configNamespace)
	}
	if systemNamespace != executorNamespace {
		namespaces = append(namespaces, executorNamespace)
	}
	if systemNamespace != env.ArgocdNamespace() {
		namespaces = append(namespaces, env.ArgocdNamespace())
	}
	return namespaces
}
