package cmd

import (
	"context"
	"reflect"
	"testing"

	dtrack "github.com/DependencyTrack/client-go"
	"github.com/google/uuid"
	"github.com/takumakume/dependency-track-policy-applier/dependencytrack"
	"github.com/takumakume/dependency-track-policy-applier/pkg"
)

func Test_applyPolicy(t *testing.T) {
	type args struct {
		ctx           context.Context
		client        dependencytrack.DependencyTrackClient
		desierdPolicy dtrack.Policy
	}
	tests := []struct {
		name       string
		args       args
		wantPolicy dtrack.Policy
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPolicy, err := applyPolicy(tt.args.ctx, tt.args.client, tt.args.desierdPolicy)
			if (err != nil) != tt.wantErr {
				t.Errorf("applyPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPolicy, tt.wantPolicy) {
				t.Errorf("applyPolicy() = %v, want %v", gotPolicy, tt.wantPolicy)
			}
		})
	}
}

func Test_applyTags(t *testing.T) {
	type args struct {
		ctx    context.Context
		client dependencytrack.DependencyTrackClient
		policy dtrack.Policy
		tags   []dtrack.Tag
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := applyTags(tt.args.ctx, tt.args.client, tt.args.policy, tt.args.tags); (err != nil) != tt.wantErr {
				t.Errorf("applyTags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_applyProjects(t *testing.T) {
	type args struct {
		ctx          context.Context
		client       dependencytrack.DependencyTrackClient
		policy       dtrack.Policy
		projectUUIDs []uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := applyProjects(tt.args.ctx, tt.args.client, tt.args.policy, tt.args.projectUUIDs); (err != nil) != tt.wantErr {
				t.Errorf("applyProjects() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_applyPolicyConditions(t *testing.T) {
	type args struct {
		ctx        context.Context
		client     dependencytrack.DependencyTrackClient
		policy     dtrack.Policy
		conditions []dtrack.PolicyCondition
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := applyPolicyConditions(tt.args.ctx, tt.args.client, tt.args.policy, tt.args.conditions); (err != nil) != tt.wantErr {
				t.Errorf("applyPolicyConditions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_desierdPolicy(t *testing.T) {
	type args struct {
		policyName     string
		operator       string
		violationState string
	}
	tests := []struct {
		name string
		args args
		want dtrack.Policy
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := desierdPolicy(tt.args.policyName, tt.args.operator, tt.args.violationState); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("desierdPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_desierdTags(t *testing.T) {
	type args struct {
		tagSlice []string
	}
	tests := []struct {
		name string
		args args
		want []dtrack.Tag
	}{
		{
			name: "Empty slice",
			args: args{
				tagSlice: []string{},
			},
			want: []dtrack.Tag{},
		},
		{
			name: "Non-empty slice",
			args: args{
				tagSlice: []string{"tag1", "tag2"},
			},
			want: []dtrack.Tag{{Name: "tag1"}, {Name: "tag2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := desierdTags(tt.args.tagSlice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("desierdTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_desierdProjectUUIDs(t *testing.T) {
	type args struct {
		ctx                 context.Context
		client              dependencytrack.DependencyTrackClient
		projectNameVersions []string
	}
	tests := []struct {
		name      string
		args      args
		wantUuids []uuid.UUID
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUuids, err := desierdProjectUUIDs(tt.args.ctx, tt.args.client, tt.args.projectNameVersions)
			if (err != nil) != tt.wantErr {
				t.Errorf("desierdProjectUUIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUuids, tt.wantUuids) {
				t.Errorf("desierdProjectUUIDs() = %v, want %v", gotUuids, tt.wantUuids)
			}
		})
	}
}

func Test_desierdPolicyConditions(t *testing.T) {
	type args struct {
		policyConditions pkg.PolicyConditions
	}
	tests := []struct {
		name      string
		args      args
		wantConds []dtrack.PolicyCondition
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotConds := desierdPolicyConditions(tt.args.policyConditions); !reflect.DeepEqual(gotConds, tt.wantConds) {
				t.Errorf("desierdPolicyConditions() = %v, want %v", gotConds, tt.wantConds)
			}
		})
	}
}

func Test_compareTags(t *testing.T) {
	type args struct {
		aa []dtrack.Tag
		bb []dtrack.Tag
	}
	tests := []struct {
		name       string
		args       args
		wantRemove []dtrack.Tag
		wantAdd    []dtrack.Tag
	}{
		{
			name: "No differences",
			args: args{
				aa: []dtrack.Tag{
					{Name: "tag1"},
					{Name: "tag2"},
				},
				bb: []dtrack.Tag{
					{Name: "tag1"},
					{Name: "tag2"},
				},
			},
			wantRemove: []dtrack.Tag(nil),
			wantAdd:    []dtrack.Tag(nil),
		},
		{
			name: "Differences",
			args: args{
				aa: []dtrack.Tag{
					{Name: "tag1"},
					{Name: "tag2"},
				},
				bb: []dtrack.Tag{
					{Name: "tag1"},
					{Name: "tag3"},
				},
			},
			wantRemove: []dtrack.Tag{{Name: "tag2"}},
			wantAdd:    []dtrack.Tag{{Name: "tag3"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRemove, gotAdd := compareTags(tt.args.aa, tt.args.bb)
			if !reflect.DeepEqual(gotRemove, tt.wantRemove) {
				t.Errorf("compareTags() gotRemove = %v, want %v", gotRemove, tt.wantRemove)
			}
			if !reflect.DeepEqual(gotAdd, tt.wantAdd) {
				t.Errorf("compareTags() gotAdd = %v, want %v", gotAdd, tt.wantAdd)
			}
		})
	}
}

func Test_compareUUIDs(t *testing.T) {
	type args struct {
		aa []uuid.UUID
		bb []uuid.UUID
	}
	tests := []struct {
		name       string
		args       args
		wantRemove []uuid.UUID
		wantAdd    []uuid.UUID
	}{
		{
			name: "No differences",
			args: args{
				aa: []uuid.UUID{
					uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				},
				bb: []uuid.UUID{
					uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				},
			},
			wantRemove: []uuid.UUID(nil),
			wantAdd:    []uuid.UUID(nil),
		},
		{
			name: "Differences",
			args: args{
				aa: []uuid.UUID{
					uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				},
				bb: []uuid.UUID{
					uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				},
			},
			wantRemove: []uuid.UUID{uuid.MustParse("00000000-0000-0000-0000-000000000002")},
			wantAdd:    []uuid.UUID{uuid.MustParse("00000000-0000-0000-0000-000000000003")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRemove, gotAdd := compareUUIDs(tt.args.aa, tt.args.bb)
			if !reflect.DeepEqual(gotRemove, tt.wantRemove) {
				t.Errorf("compareUUIDs() gotRemove = %#v, want %#v", gotRemove, tt.wantRemove)
			}
			if !reflect.DeepEqual(gotAdd, tt.wantAdd) {
				t.Errorf("compareUUIDs() gotAdd = %+v, want %+v", gotAdd, tt.wantAdd)
			}
		})
	}
}

func Test_comparePolicyConditions(t *testing.T) {
	type args struct {
		aa []dtrack.PolicyCondition
		bb []dtrack.PolicyCondition
	}
	tests := []struct {
		name        string
		args        args
		wantRemoved []dtrack.PolicyCondition
		wantAdded   []dtrack.PolicyCondition
	}{
		{
			name: "No differences",
			args: args{
				aa: []dtrack.PolicyCondition{
					{
						Subject:  dtrack.PolicyConditionSubjectVulnerabilityID,
						Operator: dtrack.PolicyConditionOperatorIs,
						Value:    "CVE-2020-1234",
					},
					{
						Subject:  dtrack.PolicyConditionSubjectVulnerabilityID,
						Operator: dtrack.PolicyConditionOperatorIs,
						Value:    "CVE-2021-1234",
					},
				},
				bb: []dtrack.PolicyCondition{
					{
						Subject:  dtrack.PolicyConditionSubjectVulnerabilityID,
						Operator: dtrack.PolicyConditionOperatorIs,
						Value:    "CVE-2020-1234",
					},
					{
						Subject:  dtrack.PolicyConditionSubjectVulnerabilityID,
						Operator: dtrack.PolicyConditionOperatorIs,
						Value:    "CVE-2021-1234",
					},
				},
			},
			wantRemoved: []dtrack.PolicyCondition(nil),
			wantAdded:   []dtrack.PolicyCondition(nil),
		},
		{
			name: "Differences",
			args: args{
				aa: []dtrack.PolicyCondition{
					{
						Subject:  dtrack.PolicyConditionSubjectVulnerabilityID,
						Operator: dtrack.PolicyConditionOperatorIs,
						Value:    "CVE-2020-1234",
					},
					{
						Subject:  dtrack.PolicyConditionSubjectVulnerabilityID,
						Operator: dtrack.PolicyConditionOperatorIs,
						Value:    "CVE-2021-1234",
					},
				},
				bb: []dtrack.PolicyCondition{
					{
						Subject:  dtrack.PolicyConditionSubjectVulnerabilityID,
						Operator: dtrack.PolicyConditionOperatorIs,
						Value:    "CVE-2021-1234",
					},
					{
						Subject:  dtrack.PolicyConditionSubjectVulnerabilityID,
						Operator: dtrack.PolicyConditionOperatorIs,
						Value:    "CVE-2022-1234",
					},
				},
			},
			wantRemoved: []dtrack.PolicyCondition{
				{
					Subject:  dtrack.PolicyConditionSubjectVulnerabilityID,
					Operator: dtrack.PolicyConditionOperatorIs,
					Value:    "CVE-2020-1234",
				},
			},
			wantAdded: []dtrack.PolicyCondition{
				{
					Subject:  dtrack.PolicyConditionSubjectVulnerabilityID,
					Operator: dtrack.PolicyConditionOperatorIs,
					Value:    "CVE-2022-1234",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRemoved, gotAdded := comparePolicyConditions(tt.args.aa, tt.args.bb)
			if !reflect.DeepEqual(gotRemoved, tt.wantRemoved) {
				t.Errorf("comparePolicyConditions() gotRemoved = %+v, want %+v", gotRemoved, tt.wantRemoved)
			}
			if !reflect.DeepEqual(gotAdded, tt.wantAdded) {
				t.Errorf("comparePolicyConditions() gotAdded = %+v, want %+v", gotAdded, tt.wantAdded)
			}
		})
	}
}
