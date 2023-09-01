# dependency-track-policy-applier

Manage OWASP Dependency-Track policy idempotent.

Create policy with projects and tags, and add conditions.

```shell
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
```

Apply the difference.

```shell
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
```

## Policy format

### JSON

```json
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
```

- `subject`
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
- `operator`
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


```shell
$ curl https://www.cisa.gov/sites/default/files/feeds/known_exploited_vulnerabilities.json | \
    jq '[.vulnerabilities[] | {subject: "VULNERABILITY_ID", operator: "IS", value: .cveID}]' | \
    DT_API_KEY="..." dependency-track-policy-applier --policy-name kev --policy-violation-state FAIL --policy-operator ANY

2023/09/01 12:49:48 apply policy: create policy: kev
2023/09/01 12:49:48 apply policyConditions: add policyCondition: VULNERABILITY_ID IS "CVE-2021-27104"
2023/09/01 12:49:48 apply policyConditions: add policyCondition: VULNERABILITY_ID IS "CVE-2021-27102"
2023/09/01 12:49:48 apply policyConditions: add policyCondition: VULNERABILITY_ID IS "CVE-2021-27101"
2023/09/01 12:49:48 apply policyConditions: add policyCondition: VULNERABILITY_ID IS "CVE-2021-27103"
2023/09/01 12:49:48 apply policyConditions: add policyCondition: VULNERABILITY_ID IS "CVE-2021-21017"
2023/09/01 12:49:48 apply policyConditions: add policyCondition: VULNERABILITY_ID IS "CVE-2021-28550"
2023/09/01 12:49:48 apply policyConditions: add policyCondition: VULNERABILITY_ID IS "CVE-2018-4939"
:
```