package terraform

import (
	"context"
	"fmt"
	"github.com/compuzest/zlifecycle-state-manager/zlog"
	"github.com/hashicorp/terraform-exec/tfexec"
	gocache "github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"time"
)

var cache = gocache.New(15*time.Minute, 30*time.Minute)

func InitTerraform(ctx context.Context, workdir string, execPath string) (*Wrapper, error) {
	if v, found := cache.Get(workdir); found {
		zlog.Logger.WithFields(
			logrus.Fields{"workdir": workdir},
		).Info("Terraform instance exists in cache")
		tf, ok := v.(*tfexec.Terraform)
		if !ok {
			return nil, fmt.Errorf("type asssertion error for terraform from cache for key %s", workdir)
		}
		return &Wrapper{ctx: ctx, tf: tf}, nil
	}

	tf, err := tfexec.NewTerraform(workdir, execPath)
	if err != nil {
		return nil, err
	}

	zlog.Logger.WithFields(
		logrus.Fields{"workdir": workdir},
	).Info("Running terraform init")
	if err := tf.Init(ctx); err != nil {
		return nil, err
	}

	zlog.Logger.WithFields(
		logrus.Fields{"workdir": workdir},
	).Info("Caching terraform instance")
	if err := cache.Add(workdir, tf, gocache.DefaultExpiration); err != nil {
		return nil, err
	}

	wrapper := Wrapper{ctx: ctx, tf: tf}
	return &wrapper, nil
}

func (w *Wrapper) State() (state *StateWrapper, err error) {
	s, err := w.tf.Show(w.ctx)
	if err != nil {
		return nil, err
	}

	return NewStateWrapper(s), nil
}



func (w *Wrapper) RemoveStateResources(resources []string) (state *StateWrapper, err error) {
	for _, r := range resources {
		if err := w.tf.StateRm(w.ctx, r); err != nil {
			return nil, err
		}
	}

	return w.State()
}
