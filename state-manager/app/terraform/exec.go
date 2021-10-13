package terraform

import (
	"context"
	"fmt"
	"github.com/compuzest/zlifecycle-state-manager/app/util"
	"github.com/compuzest/zlifecycle-state-manager/app/zlog"
	"github.com/hashicorp/terraform-exec/tfexec"
	gocache "github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

var cache = gocache.New(15*time.Minute, 30*time.Minute)

func InitTerraform(ctx context.Context, workdir string, execPath string) (*Wrapper, error) {
	deadlineCtx, cancel := deadlineContext(ctx)
	defer cancel()

	tf, err := checkCache(ctx, workdir)
	if err != nil {
		return nil, err
	}
	if tf != nil {
		return &Wrapper{ctx: ctx, tf: tf}, nil
	}

	tf, err = tfexec.NewTerraform(workdir, execPath)
	if err != nil {
		return nil, err
	}

	setTfLogLevel(tf)

	zlog.CtxLogger(ctx).WithFields(
		logrus.Fields{"workdir": workdir},
	).Info("Running terraform init")
	if err := tf.Init(deadlineCtx); err != nil {
		return nil, err
	}

	zlog.CtxLogger(ctx).WithFields(
		logrus.Fields{"workdir": workdir},
	).Info("Caching terraform instance")
	if err := cache.Add(workdir, tf, gocache.DefaultExpiration); err != nil {
		return nil, err
	}

	wrapper := Wrapper{ctx: ctx, tf: tf}
	return &wrapper, nil
}

func checkCache(ctx context.Context, workdir string) (tf *tfexec.Terraform, err error) {
	isInitialized, err := util.DirExists(path.Join(workdir, ".terraform"))
	if err != nil {
		return nil, err
	}
	if v, found := cache.Get(workdir); found && isInitialized {
		zlog.CtxLogger(ctx).WithFields(
			logrus.Fields{"workdir": workdir},
		).Info("Terraform instance exists in cache and is initialized")
		tf, ok := v.(*tfexec.Terraform)
		if !ok {
			return nil, fmt.Errorf("type asssertion error for terraform from cache for key %s", workdir)
		}
		return tf, nil
	}
	zlog.CtxLogger(ctx).WithFields(
		logrus.Fields{"workdir": workdir},
	).Info("Terraform instance does not exist in cache and/or is not initialized")
	return nil, nil
}

func (w *Wrapper) State() (state *StateWrapper, err error) {
	deadlineCtx, cancel := deadlineContext(w.ctx)
	defer cancel()

	setTfLogLevel(w.tf)
	s, err := w.tf.Show(deadlineCtx)
	if err != nil {
		return nil, err
	}

	return NewStateWrapper(s), nil
}

func (w *Wrapper) RemoveStateResources(resources []string) (state *StateWrapper, err error) {
	deadlineCtx, cancel := deadlineContext(w.ctx)
	defer cancel()

	setTfLogLevel(w.tf)
	for _, r := range resources {
		if err := w.tf.StateRm(deadlineCtx, r); err != nil {
			return nil, err
		}
	}

	return w.State()
}

func deadlineContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithDeadline(ctx, time.Now().Add(60*time.Second))
}

func setTfLogLevel(tf *tfexec.Terraform) {
	if os.Getenv("DEV_MODE") == "true" {
		tf.SetStdout(os.Stdout)
	}
}
