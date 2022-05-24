package validator

import (
	"context"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"

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

	invalidName1 := v1.Environment{
		Spec: v1.EnvironmentSpec{EnvName: "validenv1", TeamName: "invalid--team"},
	}

	invalidErrorList := validateNames(&invalidName1)
	assert.Len(t, invalidErrorList, 1)

	invalidName2 := v1.Environment{
		Spec: v1.EnvironmentSpec{EnvName: "validenv1", TeamName: "invalid-team-"},
	}

	invalidErrorList1 := validateNames(&invalidName2)
	assert.Len(t, invalidErrorList1, 1)

	e1 := v1.Environment{
		Spec: v1.EnvironmentSpec{EnvName: "validenv1", TeamName: "1-invalid-team"},
	}

	errList1 := validateNames(&e1)
	assert.Len(t, errList1, 1)

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

func TestCheckEnvironmentComponentNames(t *testing.T) {
	t.Parallel()

	errList1 := checkEnvironmentComponentName("some-component", 0)
	assert.Len(t, errList1, 0)

	errList2 := checkEnvironmentComponentName("some_component", 0)
	assert.Len(t, errList2, 1)
	assert.Equal(t, errList2[0].Field, "spec.components[0].name")
	assert.Equal(t, errList2[0].Type, field.ErrorTypeInvalid)

	errList3 := checkEnvironmentComponentName("reallyreallyreallyreallyreallyreallyreallyreallyreallyreallyreallylongenvname", 1)
	assert.Len(t, errList3, 1)
	assert.Equal(t, errList3[0].Field, "spec.components[1].name")
	assert.Equal(t, errList3[0].Type, field.ErrorTypeInvalid)
}

func TestIsUniqueEnvAndTeam(t *testing.T) {
	t.Parallel()

	envName := "test"
	teamName := "some-team"

	env := v1.Environment{
		Spec: v1.EnvironmentSpec{TeamName: teamName, EnvName: envName},
	}

	envList := v1.EnvironmentList{
		Items: []v1.Environment{{
			Spec: v1.EnvironmentSpec{TeamName: teamName, EnvName: "diff-env"},
		}},
	}

	envListDuplicate := v1.EnvironmentList{
		Items: []v1.Environment{{
			Spec: v1.EnvironmentSpec{TeamName: teamName, EnvName: envName},
		}},
	}

	err := isUniqueEnvAndTeam(&env, envListDuplicate)
	assert.Contains(t, err.Detail, fmt.Sprintf("the environment %s already exists within team %s", envName, teamName))

	err1 := isUniqueEnvAndTeam(&env, envList)
	assert.Nil(t, err1)
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

func TestCheckEnvironmentComponentDuplicateDependencies(t *testing.T) {
	t.Parallel()

	err1 := checkEnvironmentComponentDuplicateDependencies([]string{"here", "are", "duplicate", "duplicate", "entries"}, 5)
	assert.Contains(t, err1.Detail, "dependsOn cannot contain duplicates: [duplicate]")

	err2 := checkEnvironmentComponentDuplicateDependencies([]string{"here", "are", "duplicate", "duplicate", "entries", "entries", "entries", "duplicate"}, 5)
	assert.Contains(t, err2.Detail, "dependsOn cannot contain duplicates: [duplicate entries entries duplicate]")

	err := checkEnvironmentComponentDuplicateDependencies([]string{"here", "are", "duplicate", "entries"}, 5)
	assert.Nil(t, err)
}

func TestCheckValueFromsExist(t *testing.T) {
	t.Parallel()

	ec := v1.EnvironmentComponent{
		Variables: []*v1.Variable{
			{Name: "name", Value: "some-value"},
			{Name: "should-match", ValueFrom: "context.context"},
		},
	}

	ecs := []*v1.EnvironmentComponent{
		{
			Name: "unused",
			Outputs: []*v1.Output{
				{Name: "unused", Sensitive: false},
				{Name: "blah-unused", Sensitive: false},
			},
		},
		{
			Name: "context",
			Outputs: []*v1.Output{
				{Name: "context"},
				{Name: "doesnt-match"},
			},
		},
	}

	errs := checkValueFromsExist(&ec, ecs)
	assert.Nil(t, errs)
}

func TestCheckValueFromsExistBadValueFrom(t *testing.T) {
	t.Parallel()

	ecBadValueFrom := v1.EnvironmentComponent{
		Variables: []*v1.Variable{
			{Name: "bad-value-from", ValueFrom: "blah-context"},
			{Name: "unused", Value: "some-value"},
			{Name: "no-match", ValueFrom: "blah.context"},
			{Name: "should-match", ValueFrom: "context.context"},
			{Name: "good-component", ValueFrom: "context.badVariable"},
		},
	}

	ecs := []*v1.EnvironmentComponent{
		{
			Name: "unused",
			Outputs: []*v1.Output{
				{Name: "unused", Sensitive: false},
				{Name: "blah-unused", Sensitive: false},
			},
		},
		{
			Name: "context",
			Outputs: []*v1.Output{
				{Name: "context"},
				{Name: "doesnt-match"},
			},
		},
	}

	errs := checkValueFromsExist(&ecBadValueFrom, ecs)
	assert.Len(t, errs, 3)
	assert.Contains(t, errs[0].Detail, "valueFrom must be 'componentName.componentOutputName'")
	assert.Contains(t, errs[1].Detail, "valueFrom blah.context references component blah which does not exist")
	assert.Contains(t, errs[2].Detail, "valueFrom context.badVariable does not match any outputs")
}
