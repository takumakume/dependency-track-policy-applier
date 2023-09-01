package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	dtrack "github.com/DependencyTrack/client-go"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/takumakume/dependency-track-policy-applier/config"
	"github.com/takumakume/dependency-track-policy-applier/dependencytrack"
	"github.com/takumakume/dependency-track-policy-applier/pkg"
)

var rootCmd = &cobra.Command{
	Use:   "dependency-track-policy-applier",
	Short: "Manage OWASP Dependency-Track policy idempotent.",
	Long: `Manage OWASP Dependency-Track policy idempotent.

    $ echo '[
        {
            "subject": "VULNERABILITY_ID",
            "operator": "IS",
            "value": "CVE-2023-11111"
        },
        {
            "subject": "VULNERABILITY_ID",
            "operator": "IS",
            "value": "CVE-2023-22222"
        }
    ]' | DT_API_KEY="..." dependency-track-policy-applier --policy-name myPolicy --policy-projects test:latest --policy-tags foo

    2023/09/01 12:00:01 apply policy: create policy: myPolicy
    2023/09/01 12:00:01 apply tags: add tag "foo"
    2023/09/01 12:00:01 apply projects: add project 451b427e-cd46-45f0-98eb-63705c4dc624
    2023/09/01 12:00:01 apply policyConditions: add policyCondition: VULNERABILITY_ID IS "CVE-2023-11111"
    2023/09/01 12:00:01 apply policyConditions: add policyCondition: VULNERABILITY_ID IS "CVE-2023-22222"

Apply the difference.

    $ echo '[
        {
        "subject": "VULNERABILITY_ID",
        "operator": "IS",
        "value": "CVE-2023-11111"
        }
    ]' | DT_API_KEY="..." dependency-track-policy-applier --policy-name myPolicy 

    2023/09/01 12:42:46 apply tags: remove tag "foo"
    2023/09/01 12:42:46 apply projects: remove project 451b427e-cd46-45f0-98eb-63705c4dc624
    2023/09/01 12:42:46 apply policyConditions: remove policyCondition: VULNERABILITY_ID IS "CVE-2023-22222"

## Policy format

### JSON

    [
      {
        "subject": "VULNERABILITY_ID",
        "operator": "IS",
        "value": "CVE-2023-11111"
      },
      {
        "subject": "VULNERABILITY_ID",
        "operator": "IS",
        "value": "CVE-2023-22222"
      }
    ]

    - subject
      - "AGE"
      - "COORDINATES"
      - "CPE"
      - "LICENSE"
      - "LICENSE_GROUP"
      - "PACKAGE_URL"
      - "SEVERITY"
      - "SWID_TAGID"
      - "VERSION"
      - "COMPONENT_HASH"
      - "CWE"
      - "VULNERABILITY_ID"
    - operator
      - "IS"
      - "IS_NOT"
      - "MATCHES"
      - "NO_MATCH"
      - "NUMERIC_GREATER_THAN"
      - "NUMERIC_LESS_THAN"
      - "NUMERIC_EQUAL"
      - "NUMERIC_NOT_EQUAL"
      - "NUMERIC_GREATER_THAN_OR_EQUAL"
      - "NUMERIC_LESSER_THAN_OR_EQUAL"
      - "CONTAINS_ALL"
      - "CONTAINS_ANY"

## case: Generate Policy based on KEV (Known Exploited Vulnerabilities)

    $ curl https://www.cisa.gov/sites/default/files/feeds/known_exploited_vulnerabilities.json | \
        jq '[.vulnerabilities[] | {subject: "VULNERABILITY_ID", operator: "IS", value: .cveID}]' | \
        DT_API_KEY="..." dependency-track-policy-applier --policy-name kev --policy-violation-state FAIL --policy-operator ANY

    2023/09/01 12:49:48 apply policy: create policy: kev
    2023/09/01 12:49:48 apply policyConditions: add policyCondition: VULNERABILITY_ID IS "CVE-2021-27104"
    2023/09/01 12:49:48 apply policyConditions: add policyCondition: VULNERABILITY_ID IS "CVE-2021-27102"
    :
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c := config.New(
			viper.GetString("base-url"),
			viper.GetString("api-key"),
			viper.GetString("policy-name"),
			viper.GetString("policy-operator"),
			viper.GetString("policy-violation-state"),
			viper.GetStringSlice("policy-projects"),
			viper.GetStringSlice("policy-tags"),
		)
		if err := c.Validate(); err != nil {
			return err
		}

		b, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		pcs, err := pkg.NewPolicyConditions(b)
		if err != nil {
			return err
		}

		dtrackClient, err := dependencytrack.New(c.BaseURL, c.APIKey, 10*time.Second)
		if err != nil {
			return err
		}

		return run(ctx, dtrackClient, c, pcs)
	},
}

func init() {
	flags := rootCmd.PersistentFlags()
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetEnvPrefix("DT")

	flags.StringP("base-url", "u", "http://127.0.0.1:8081/", "Dependency Track base URL (env: DT_BASE_URL)")
	flags.StringP("api-key", "k", "", "Dependency Track API key (env: DT_API_KEY)")
	flags.StringP("policy-name", "", "", "Dependency Track policy name")
	flags.StringP("policy-operator", "", string(dtrack.PolicyOperatorAny), fmt.Sprintf("Dependency Track policy operator %v", []dtrack.PolicyOperator{dtrack.PolicyOperatorAny, dtrack.PolicyOperatorAll}))
	flags.StringP("policy-violation-state", "", string(dtrack.PolicyViolationStateFail), fmt.Sprintf("Dependency Track policy violationState %v", []dtrack.PolicyViolationState{dtrack.PolicyViolationStateFail, dtrack.PolicyViolationStateWarn, dtrack.PolicyViolationStateInfo}))
	flags.StringSliceP("policy-projects", "", []string{}, "Dependency Track policy projects")
	flags.StringSliceP("policy-tags", "", []string{}, "Dependency Track policy tags")

	viper.BindPFlag("base-url", flags.Lookup("base-url"))
	viper.BindPFlag("api-key", flags.Lookup("api-key"))
	viper.BindPFlag("policy-name", flags.Lookup("policy-name"))
	viper.BindPFlag("policy-operator", flags.Lookup("policy-operator"))
	viper.BindPFlag("policy-violation-state", flags.Lookup("policy-violation-state"))
	viper.BindPFlag("policy-projects", flags.Lookup("policy-projects"))
	viper.BindPFlag("policy-tags", flags.Lookup("policy-tags"))
}

func Execute() error {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	return rootCmd.Execute()
}

func run(ctx context.Context, client dependencytrack.DependencyTrackClient, config *config.Config, policyConditions pkg.PolicyConditions) error {
	desierdPolicy := desierdPolicy(config.PolicyName, config.PolicyOperator, config.PolicyViolationState)
	policy, err := applyPolicy(ctx, client, desierdPolicy)
	if err != nil {
		return err
	}

	tags := desierdTags(config.PolicyTags)
	if err := applyTags(ctx, client, policy, tags); err != nil {
		return err
	}

	projectUUIDs, err := desierdProjectUUIDs(ctx, client, config.PolicyProjects)
	if err != nil {
		return err
	}
	if err := applyProjects(ctx, client, policy, projectUUIDs); err != nil {
		return err
	}

	desierdPolicyConditions := desierdPolicyConditions(policyConditions)
	if err := applyPolicyConditions(ctx, client, policy, desierdPolicyConditions); err != nil {
		return err
	}

	return nil
}

func applyPolicy(ctx context.Context, client dependencytrack.DependencyTrackClient, desierdPolicy dtrack.Policy) (policy dtrack.Policy, err error) {
	if policy, err = client.GetPolicyForName(ctx, desierdPolicy.Name); err != nil {
		if dependencytrack.IsNotFound(err) {
			log.Printf("apply policy: create policy: %s", desierdPolicy.Name)

			policy, err = client.CreatePolicy(ctx, desierdPolicy)
			if err != nil {
				return policy, err
			}
			// FIXME: https://github.com/DependencyTrack/dependency-track/issues/2365
			policy.Operator = dtrack.PolicyOperator(desierdPolicy.Operator)
			policy.ViolationState = dtrack.PolicyViolationState(desierdPolicy.ViolationState)
			policy, err = client.UpdatePolicy(ctx, policy)
			if err != nil {
				return policy, err
			}
		} else {
			return policy, err
		}

	} else {
		if client.NeedsUpdatePolicy(policy, desierdPolicy) {
			log.Printf("apply policy: update policy: %s", desierdPolicy.Name)

			policy.Operator = desierdPolicy.Operator
			policy.ViolationState = desierdPolicy.ViolationState
			policy, err = client.UpdatePolicy(ctx, policy)
			if err != nil {
				return policy, err
			}
		}
	}

	return policy, err
}

func applyTags(ctx context.Context, client dependencytrack.DependencyTrackClient, policy dtrack.Policy, tags []dtrack.Tag) error {
	remove, add := compareTags(policy.Tags, tags)
	for _, o := range remove {
		log.Printf("apply tags: remove tag %v", o)

		_, err := client.DeleteTag(ctx, policy.UUID, o.Name)
		if err != nil {
			if dependencytrack.IsNotFound(err) {
				log.Printf("WARN: apply tags: remove tag: not found %v", o)

				continue
			}
			return err
		}
	}
	for _, o := range add {
		log.Printf("apply tags: add tag %v", o)

		_, err := client.AddTag(ctx, policy.UUID, o.Name)
		if err != nil {
			if dependencytrack.IsNotFound(err) {
				log.Printf("WARN: apply tags: add tag: not found %v", o)

				continue
			}
			return err
		}
	}
	return nil
}

func applyProjects(ctx context.Context, client dependencytrack.DependencyTrackClient, policy dtrack.Policy, projectUUIDs []uuid.UUID) error {
	currentProjectUUIDs := []uuid.UUID{}
	for _, p := range policy.Projects {
		currentProjectUUIDs = append(currentProjectUUIDs, p.UUID)
	}

	remove, add := compareUUIDs(currentProjectUUIDs, projectUUIDs)
	for _, o := range remove {
		log.Printf("apply projects: remove project %s", o)

		_, err := client.DeleteProject(ctx, policy.UUID, o)
		if err != nil {
			return err
		}
	}
	for _, o := range add {
		log.Printf("apply projects: add project %s", o)

		_, err := client.AddProject(ctx, policy.UUID, o)
		if err != nil {
			return err
		}
	}
	return nil
}

func applyPolicyConditions(ctx context.Context, client dependencytrack.DependencyTrackClient, policy dtrack.Policy, conditions []dtrack.PolicyCondition) error {
	remove, add := comparePolicyConditions(policy.PolicyConditions, conditions)
	for _, o := range remove {
		log.Printf("apply policyConditions: remove policyCondition: %s %s %q", o.Subject, o.Operator, o.Value)

		if err := client.DeletePolicyCondition(ctx, o.UUID); err != nil {
			return err
		}
	}
	for _, o := range add {
		log.Printf("apply policyConditions: add policyCondition: %s %s %q", o.Subject, o.Operator, o.Value)

		_, err := client.CreatePolicyCondition(ctx, policy.UUID, o)
		if err != nil {
			return err
		}
	}
	return nil
}

func desierdPolicy(policyName, operator, violationState string) dtrack.Policy {
	return dtrack.Policy{
		Name:           policyName,
		Operator:       dtrack.PolicyOperator(operator),
		ViolationState: dtrack.PolicyViolationState(violationState),
	}
}

func desierdTags(tagSlice []string) []dtrack.Tag {
	m := make(map[string]bool)
	uniq := []string{}

	for _, ele := range tagSlice {
		if !m[ele] {
			m[ele] = true
			uniq = append(uniq, ele)
		}
	}

	tags := make([]dtrack.Tag, len(uniq))
	for i, s := range uniq {
		tags[i] = dtrack.Tag{Name: s}
	}

	return tags
}

func desierdProjectUUIDs(ctx context.Context, client dependencytrack.DependencyTrackClient, projectNameVersions []string) (uuids []uuid.UUID, err error) {
	projects := []dtrack.Project{}
	for _, nv := range projectNameVersions {
		projectNameVersion := strings.SplitN(nv, ":", 2)
		if len(projectNameVersion) == 2 {
			p, err := client.GetProjectForNameVersion(ctx, projectNameVersion[0], projectNameVersion[1], true, true)
			if err != nil {
				if dependencytrack.IsNotFound(err) {
					log.Printf("WARN: desierdProjectUUIDs: GetProjectForNameVersion: project version not found %q", nv)

					continue
				}
				return uuids, err
			}
			projects = append(projects, p)
		} else {
			pp, err := client.GetProjectsForName(ctx, projectNameVersion[0], true, true)
			if err != nil {
				if dependencytrack.IsNotFound(err) {
					log.Printf("WARN: desierdProjectUUIDs: GetProjectsForName: project not found %q", nv)

					continue
				}
				return uuids, err
			}
			projects = append(projects, pp...)
		}
	}

	seen := make(map[uuid.UUID]bool)
	for _, project := range projects {
		if _, ok := seen[project.UUID]; ok {
			continue
		}
		seen[project.UUID] = true
		uuids = append(uuids, project.UUID)
	}

	return uuids, nil
}

func desierdPolicyConditions(policyConditions pkg.PolicyConditions) (conds []dtrack.PolicyCondition) {
	for _, c := range policyConditions {
		conds = append(conds, dtrack.PolicyCondition{
			Subject:  c.Subject,
			Operator: c.Operator,
			Value:    c.Value,
		})
	}
	return conds
}

func compareTags(aa, bb []dtrack.Tag) (remove, add []dtrack.Tag) {
	aaMap := make(map[string]dtrack.Tag)
	for _, a := range aa {
		aaMap[a.Name] = a
	}

	for _, b := range bb {
		_, ok := aaMap[b.Name]
		if ok {
			delete(aaMap, b.Name)
			continue
		}
		add = append(add, b)
	}

	for _, a := range aaMap {
		remove = append(remove, a)
	}

	return remove, add
}

func compareUUIDs(aa, bb []uuid.UUID) (remove, add []uuid.UUID) {
	aaMap := make(map[uuid.UUID]uuid.UUID)
	for _, a := range aa {
		aaMap[a] = a
	}

	for _, b := range bb {
		_, ok := aaMap[b]
		if ok {
			delete(aaMap, b)
			continue
		}
		add = append(add, b)
	}

	for _, a := range aaMap {
		remove = append(remove, a)
	}

	return remove, add
}

func comparePolicyConditions(aa, bb []dtrack.PolicyCondition) (removed, added []dtrack.PolicyCondition) {
	aaMap := make(map[string]dtrack.PolicyCondition)
	for _, a := range aa {
		aaMap[a.Value] = a
	}

	for _, b := range bb {
		_, ok := aaMap[b.Value]
		if ok {
			delete(aaMap, b.Value)
			continue
		}
		added = append(added, b)
	}

	for _, a := range aaMap {
		removed = append(removed, a)
	}

	return removed, added
}
