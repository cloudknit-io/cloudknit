package argocd

import (
	"fmt"
	"regexp"
	"strings"
)

func parsePolicyCsv(policyCsv string) (*RbacMap, error) {
	rbacMap := newRbacMap()
	for _, record := range strings.Split(policyCsv, "\n") {
		if record == "" {
			continue
		}
		id := string(record[0])
		switch id {
		case "p":
			policy, err := parsePolicy(record)
			if err != nil {
				return nil, err
			}
			rbacMap.Policies[policy.Subject] = append(rbacMap.Policies[policy.Subject], policy)
		case "g":
			group, err := parseGroup(record)
			if err != nil {
				return nil, err
			}
			rbacMap.Groups[group.Group] = append(rbacMap.Groups[group.Group], group)
		}
	}

	return rbacMap, nil
}

func (rbacMap *RbacMap) generatePolicyCsv() string {
	var csv string
	for _, subjectPolicies := range rbacMap.Policies {
		for _, policy := range subjectPolicies {
			csv += fmt.Sprintf("%s\n", policy.toCsv())
		}
	}
	for _, subjectGroups := range rbacMap.Groups {
		for _, group := range subjectGroups {
			csv += fmt.Sprintf("%s\n", group.toCsv())
		}
	}
	return csv
}

func (rbacMap *RbacMap) updateRbac(subject string, projects []string, oidcGroup string) {
	var policies []*RbacPolicy
	policies = append(policies, newRepositoryPolicy(subject))
	for _, project := range projects {
		policies = append(policies, newApplicationPolicy(subject, "*", project))
	}
	rbacMap.Policies[subject] = policies
	groups := []*RbacGroup{newGroup(oidcGroup, subject)}
	rbacMap.Groups[oidcGroup] = groups
}

func newApplicationPolicy(subject string, action string, project string) *RbacPolicy {
	return newPolicy(subject, "applications", action, fmt.Sprintf("%s/*", project), Allow)
}

func newRepositoryPolicy(subject string) *RbacPolicy {
	return newPolicy(subject, "repositories", "get", "*", Allow)
}

func newPolicy(subject string, resource string, action string, object string, permission Permission) *RbacPolicy {
	p := RbacPolicy{
		Identifier: Policy,
		Subject: subject,
		Resource: resource,
		Action: action,
		Object: object,
		Permission: permission,
	}
	return &p
}

func newGroup(group string, role string) *RbacGroup {
	g := RbacGroup{
		Identifier: Group,
		Group: group,
		Role: role,
	}
	return &g
}

func newRbacMap() *RbacMap {
	rbac := RbacMap{
		Policies: make(map[string][]*RbacPolicy),
		Groups: make(map[string][]*RbacGroup),
	}
	return &rbac
}

func parsePolicy(p string) (*RbacPolicy, error) {
	rgx := regexp.MustCompile("^(.+),(.+),(.+),(.+),(.+),(.+)$")
	matched := rgx.FindAllStringSubmatch(p, -1)
	if matched == nil || len(matched) != 1 || len(matched[0]) != 7 {
		return nil, fmt.Errorf("error file parsing policy: %s", p)
	}
	parsed := RbacPolicy{
		Identifier: matched[0][1],
		Subject: matched[0][2],
		Resource: matched[0][3],
		Action: matched[0][4],
		Object: matched[0][5],
		Permission: matched[0][6],
	}

	return &parsed, nil
}

func parseGroup(g string) (*RbacGroup, error) {
	rgx := regexp.MustCompile("^(.+),(.+),(.+)$")
	matched := rgx.FindAllStringSubmatch(g, -1)
	if matched == nil || len(matched) != 1 || len(matched[0]) != 4 {
		return nil, fmt.Errorf("error parsing group: %s", g)
	}
	parsed := RbacGroup{
		Identifier: matched[0][1],
		Group: matched[0][2],
		Role: matched[0][3],
	}

	return &parsed, nil
}

func (p *RbacPolicy) toCsv() string {
	return fmt.Sprintf("%s,%s,%s,%s,%s,%s", p.Identifier, p.Subject, p.Resource, p.Action, p.Object, p.Permission)
}

func (g *RbacGroup) toCsv() string {
	return fmt.Sprintf("%s,%s,%s", g.Identifier, g.Group, g.Role)
}
