package validators

import (
	"context"
	"github.com/sirupsen/logrus"
	"testing"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestValidateNames(t *testing.T) {
	t.Parallel()

	e1 := v1.Environment{
		Spec: v1.EnvironmentSpec{EnvName: "validenv1", TeamName: "1-valid-team"},
	}

	errList1 := validateNames(&e1)
	assert.Len(t, errList1, 0)

	e2 := v1.Environment{
		Spec: v1.EnvironmentSpec{EnvName: "invalid_env", TeamName: "team1"},
	}

	errList2 := validateNames(&e2)
	assert.Len(t, errList2, 1)
	assert.Equal(t, errList2[0].Field, "spec.envName")
	assert.Equal(t, errList2[0].Type, field.ErrorTypeInvalid)

	e3 := v1.Environment{
		Spec: v1.EnvironmentSpec{EnvName: "reallyreallyreallyreallyreallyreallyreallyreallyreallyreallyreallylongenvname", TeamName: "team_name"},
	}

	errList3 := validateNames(&e3)
	assert.Len(t, errList3, 2)
	assert.Equal(t, errList3[0].Field, "spec.envName")
	assert.Equal(t, errList3[0].Type, field.ErrorTypeInvalid)
	assert.Equal(t, errList3[1].Field, "spec.teamName")
	assert.Equal(t, errList3[1].Type, field.ErrorTypeInvalid)
}

func TestValidateTeamExists(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockKClient := mocks.NewMockClient(mockCtrl)

	namespace := "test"
	teamName := "some-team"
	team := v1.Team{
		ObjectMeta: metav1.ObjectMeta{Namespace: namespace},
		Spec:       v1.TeamSpec{TeamName: teamName},
	}
	env := v1.Environment{
		ObjectMeta: metav1.ObjectMeta{Namespace: namespace},
		Spec:       v1.EnvironmentSpec{TeamName: teamName},
	}

	teamList := v1.TeamList{}
	mockKClient.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).Do(func(context.Context, client.ObjectList, ...client.ListOption) {
		teamList.Items = []v1.Team{team}
	})

	log := logrus.NewEntry(logrus.New())
	err := validateTeamExists(context.Background(), &env, mockKClient, &teamList, log)
	assert.Nil(t, err)
	assert.Len(t, teamList.Items, 1)
	assert.Equal(t, teamList.Items[0].Spec.TeamName, teamName)
}
