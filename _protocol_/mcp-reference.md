# MDPlanner MCP Reference

79+ tools across 31 domains. All return `{ content: [{ type: "text", text:
JSON.stringify(data) }] }` on success or `{ isError: true, content: [...] }`
on failure.

## Response Shape

```
ok  → content[0].text = JSON.stringify(data, null, 2)
err → content[0].text = "Error: <message>", isError: true
```

## Global Conventions

- IDs: auto-generated timestamp-based strings, returned as `id` on create
- Dates: `YYYY-MM-DD` strings throughout
- Required fields: marked `(req)` below; all others optional
- Input names: snake_case (e.g. `start_date`); stored as camelCase
- Enums: exact case-sensitive match when filtering
- Array filters: OR logic (any match includes the record)
- `completed: true` is ALWAYS set by the human, never by Claude

---

## Tasks

```
list_tasks           section? project? milestone? assignee? priority? completed? tags?
get_task             id(req)
create_task          title(req) section? description? assignee? due_date? priority?
                     effort? tags? milestone? project? blocked_by? planned_start?
                     planned_end? parentId?
update_task          id(req) + any create fields + completed? blocked_by?
add_task_comment     id(req) comment(req) author?(default:"Claude")
add_task_attachments id(req) attachments(req)[str]  ← adds file paths to task
move_task            id(req) section(req) position?(integer)
delete_task          id(req)
```

Sections (standard): `Backlog` `Todo` `In Progress` `Done`
Tags (standard): `Bug` `Content` `Design` `Doc` `Feature` `Idea` `QOL`
`Question` `Refactor` `Scope` `Tech_debt` `Test`
Priority: integer 1–5 (1 = highest)

**Workflow rule**: always move to `In Progress` + assign Claude before starting;
move to `Done` + assign owner when complete; never set `completed: true`.
`section` + `assignee` can be combined in a single `update_task` call.

---

## Notes

```
list_notes   search?(substring match on title)
get_note     id(req)
create_note  title(req) content?(markdown)
update_note  id(req) title? content?(full replacement)
delete_note  id(req)
```

Notes are the primary knowledge store. Title prefix = discoverable tag.
See `CLAUDE.md` for the `[project]` `[architecture]` `[decision]` etc. convention.

---

## Goals

```
list_goals   status?(planning|on-track|at-risk|late|success|failed) type?(enterprise|project)
get_goal     id(req)
create_goal  title(req) description? status?(default:planning) type?(default:project)
             kpi? startDate? endDate?
update_goal  id(req) + any create fields
delete_goal  id(req)
```

---

## Milestones

```
list_milestones   status?(open|completed) project?
get_milestone     id(req)
create_milestone  name(req) project(req) description? target?(YYYY-MM-DD) status?(default:open)
update_milestone  id(req) + any create fields
delete_milestone  id(req)
```

**CRITICAL**: `project` is REQUIRED on create. (project, name) must be unique —
returns duplicate error with existing ID if not.

---

## People

```
list_people           department?
get_person            id(req)
create_person         name(req) title? email? phone? departments?[str] reportsTo?(person_id)
                      startDate? notes? role? hoursPerDay?(number) workingDays?[str]
update_person         id(req) + any create fields
delete_person         id(req)
get_people_tree       (no params) ← returns org hierarchy for org chart
get_people_summary    (no params) ← aggregate stats (count, departments, etc.)
get_people_departments (no params) ← distinct department list
get_person_reports    id(req) ← direct reports of a person
```

`role`: team role distinct from title (e.g. "Tech Lead").
`hoursPerDay` + `workingDays`: used by capacity planning (`["Mon","Tue","Wed","Thu","Fri"]`).

---

## Meetings

```
list_meetings   date_from?(YYYY-MM-DD) date_to? open_actions_only?(bool)
get_meeting     id(req)
create_meeting  title(req) date(req) attendees?[str] agenda? notes?
                actions?[{description, owner?, due?, status?}]
update_meeting  id(req) + any create fields  (actions = full replacement)
delete_meeting  id(req)
```

Action status: `open` (default) | `done`
List returns summary: id, title, date, attendees[], actionCount, openActions.

---

## Journal

```
list_journal_entries   (no params — sorted by date desc)
get_journal_entry      id(req)
create_journal_entry   date(req)(YYYY-MM-DD) title? mood? tags?[str] body?(markdown)
update_journal_entry   id(req) + any create fields
delete_journal_entry   id(req)
```

Mood: `great` `good` `neutral` `bad` `terrible`

---

## Ideas

```
list_ideas   category?
get_idea     id(req)
create_idea  title(req) category? priority?(low|medium|high) description?
             start_date? end_date? resources?
update_idea  id(req) + any create fields
delete_idea  id(req)
```

Status `new` is set automatically on create.

---

## Project Config

```
get_project_config    (no params)
update_project_config name? description? features?[str] settings?(object)
get_analytics         (no params)
```

`get_project_config` returns: name, description, serverVersion, projectPath,
cacheEnabled, plus full project config (features, settings, etc.).

`get_analytics` returns aggregated cross-entity health stats: task completion
rates, open action items, goal progress, storage usage, people/meeting counts.
No `--cache` required — computed from markdown files directly.

---

## Search

```
search   query(req)(min 1 char) types?[str] limit?(1-100, default:20)
```

**Requires `--cache` flag** on server start. Returns explicit error if cache
disabled. Types: `task` `note` `goal` `meeting` `person` etc.

---

## MoSCoW

```
list_moscow    (no params)
get_moscow     id(req)
create_moscow  title(req) date?(default:today)
update_moscow  id(req) title? description? must?[str] should?[str]
               could?[str] wont?[str]
delete_moscow  id(req)
```

---

## Eisenhower Matrix

```
list_eisenhower    (no params)
get_eisenhower     id(req)
create_eisenhower  title(req) date?(default:today)
update_eisenhower  id(req) title? urgentImportant?[str] notUrgentImportant?[str]
                   urgentNotImportant?[str] notUrgentNotImportant?[str]
delete_eisenhower  id(req)
```

---

## Retrospectives

```
list_retrospectives    (no params)
get_retrospective      id(req)
create_retrospective   title(req) date?(default:today) status?(default:open)
update_retrospective   id(req) title? date? status?(open|closed)
                       continue?[str] stop?[str] start?[str]
delete_retrospective   id(req)
```

---

## SWOT

```
list_swot    (no params)
get_swot     id(req)
create_swot  title(req) date?(default:today)
update_swot  id(req) title? strengths?[str] weaknesses?[str]
             opportunities?[str] threats?[str]
delete_swot  id(req)
```

---

## Risk Analysis

```
list_risks    (no params)
get_risk      id(req)
create_risk   title(req) date?(default:today)
update_risk   id(req) title? highImpactHighProb?[str] highImpactLowProb?[str]
              lowImpactHighProb?[str] lowImpactLowProb?[str]
delete_risk   id(req)
```

---

## Project Brief

```
list_briefs    (no params)
get_brief      id(req)
create_brief   title(req) date? summary?[str] mission?[str]
               responsible?[str] accountable?[str] consulted?[str] informed?[str]
               highLevelBudget?[str] highLevelTimeline?[str]
               culture?[str] changeCapacity?[str] guidingPrinciples?[str]
update_brief   id(req) + any create fields
delete_brief   id(req)
```

---

## SAFE Agreements

```
list_safe    (no params)
get_safe     id(req)
create_safe  investor(req) amount(req) valuation_cap?(default:0) discount?(0-100)
             type?(pre-money|post-money|mfn, default:post-money)
             status?(draft|signed|converted, default:draft) date? notes?
update_safe  id(req) + any create fields
delete_safe  id(req)
```

---

## Finances

```
list_finances    (no params)
get_finance      id(req)
create_finance   period(req)(e.g."2026-03") cash_on_hand?
                 revenue?[{category,amount}] expenses?[{category,amount}] notes?
update_finance   id(req) + any create fields
delete_finance   id(req)
```

---

## KPIs

```
list_kpis    (no params)
get_kpi      id(req)
create_kpi   period(req)(e.g."2026-Q1") mrr? churn_rate? ltv? cac?
             growth_rate? active_users? nrr? gross_margin? notes?
update_kpi   id(req) + any create fields
delete_kpi   id(req)
```

ARR auto-calculated as `mrr * 12`. All numeric fields default to 0.

---

## Investors

```
list_investors    status?(not_started|in_progress|term_sheet|passed|invested)
get_investor      id(req)
create_investor   name(req) type?(vc|angel|family_office|corporate|accelerator)
                  stage?(lead|associate|partner|passed) status?(default:not_started)
                  amount_target? contact? intro_date? last_contact? notes?
update_investor   id(req) + any create fields
delete_investor   id(req)
```

---

## Onboarding

```
list_onboarding    (no params)
get_onboarding     id(req)
create_onboarding  employeeName(req) role(req) startDate?(default:today)
                   personId? notes?
update_onboarding  id(req) + any create fields
delete_onboarding  id(req)

list_onboarding_templates    (no params)
get_onboarding_template      id(req)
create_onboarding_template   name(req) description?
                              steps?[{title, category}]
delete_onboarding_template   id(req)
```

Step categories: `equipment` `accounts` `docs` `training` `intro` `other`

---

## CRM

```
list_companies    (no params)
get_company       id(req)
create_company    name(req) industry? website? phone? notes?
update_company    id(req) + any create fields
delete_company    id(req)

list_contacts     company_id?
get_contact       id(req)
create_contact    first_name(req) last_name(req) company_id(req) email? phone?
                  title? is_primary? notes?
update_contact    id(req) + any create fields
delete_contact    id(req)

list_deals        stage?(lead|qualified|proposal|negotiation|won|lost) company_id?
get_deal          id(req)
create_deal       title(req) company_id(req) value? probability?(0-100)
                  stage?(default:lead) expected_close_date? notes?
update_deal       id(req) + any create fields
delete_deal       id(req)
```

---

## Billing

```
list_customers    (no params)
get_customer      id(req)
create_customer   name(req) email? company? address? notes?
update_customer   id(req) + any create fields
delete_customer   id(req)

list_quotes       (no params)
get_quote         id(req)
create_quote      customer_id(req) items(req)[{description,quantity,unit_price}]
                  notes? valid_until?(YYYY-MM-DD) status?(draft|sent|accepted|rejected)
update_quote      id(req) + any create fields
delete_quote      id(req)

list_invoices     (no params)
get_invoice       id(req)
create_invoice    customer_id(req) items(req)[{description,quantity,unit_price}]
                  notes? due_date?(YYYY-MM-DD) status?(draft|sent|paid|overdue|cancelled)
update_invoice    id(req) + any create fields
delete_invoice    id(req)
```

---

## Portfolio

```
list_portfolio              status?(active|completed|archived|production|maintenance|cancelled)
get_portfolio_item          id(req)
create_portfolio_item       name(req) description? status?(default:active) category?
                            license? start_date? end_date? team?[str] tech_stack?[str]
                            client? revenue?(number) expenses?(number) progress?(0-100)
                            kpis?[str] urls?[{label,href}] logo?(path|URL)
                            billingCustomerId? githubRepo?(owner/repo)
update_portfolio_item       id(req) + any create fields
delete_portfolio_item       id(req)
get_portfolio_summary       (no params) ← aggregated stats (count by status, avg progress, etc.)
add_portfolio_status_update id(req) status(req) note?(markdown)
delete_portfolio_status_update id(req) update_id(req)
```

Status is not a fixed enum; any string is accepted.
`urls`: pass `[]` (explicit empty array) to clear all URLs — omitting the field keeps existing.

---

## Habits

```
list_habits           (no params)
get_habit             id(req)
create_habit          name(req) description? frequency?(daily|weekly, default:daily)
                      target_days?[str] notes?
update_habit          id(req) + any create fields + dayNotes?(Record<YYYY-MM-DD,string>)
mark_habit_complete   id(req) date?(default:today) note?(per-day note string)
unmark_habit_complete id(req) date(req)
delete_habit          id(req)
```

Mark/unmark auto-recalculate `streakCount`. Returns `{ success, streakCount }`.
`dayNotes` in `update_habit`: bulk update of date-keyed notes (full replacement).

---

## DNS

```
list_dns_domains    (no params)
get_dns_domain      id(req)
create_dns_domain   domain(req) expiry_date? auto_renew? renewal_cost_usd?
                    provider? nameservers?[str] notes? project?
update_dns_domain   id(req) + any create fields
delete_dns_domain   id(req)
sync_cloudflare_dns (no params — requires Cloudflare token in settings)

list_dns_records    domain_id(req)
add_dns_record      domain_id(req) type(req)(A|CNAME|MX|TXT|NS|SRV|AAAA)
                    name(req) value(req) ttl?(default:3600)
update_dns_record   domain_id(req) record_id(req) type? name? value? ttl?
delete_dns_record   domain_id(req) record_id(req)
```

---

## GitHub

All tools require GitHub PAT token configured in Settings > Integrations.
Returns error response (not throw) if token missing.

```
github_list_repos      query?(substring filter on owner/repo name)
github_get_repo        owner(req) repo(req)
github_get_issue       owner(req) repo(req) number(req)
github_create_issue    owner(req) repo(req) title(req) body?
github_set_issue_state owner(req) repo(req) number(req) state(req)(open|closed)
github_get_pr          owner(req) repo(req) number(req)
```

---

## Fishbone Diagrams

```
list_fishbones    (no params)
get_fishbone      id(req)
create_fishbone   title(req)(the problem/effect) description?
                  causes?[{category, subcauses:[str]}]
update_fishbone   id(req) title? description? causes?(full replacement)
delete_fishbone   id(req)
```

Default cause categories: People, Process, Machine, Material, Method, Measurement.

---

## Mindmaps

```
list_mindmaps    (no params)
get_mindmap      id(req)
create_mindmap   title(req) nodes?[{id,text,level,children:[],parent?}]
update_mindmap   id(req) title? nodes?(full replacement)
delete_mindmap   id(req)
```

Stored in projectInfo (not separate files). Nodes use flexible schema.

---

## Canvas / Sticky Notes

```
list_sticky_notes    (no params)
create_sticky_note   content(req) color?(yellow|pink|blue|green|purple|orange)
                     x?(default:100) y?(default:100)
delete_sticky_note   id(req)
```

---

## Capacity Planning

```
list_capacity_plans      (no params)
get_capacity_plan        id(req)
create_capacity_plan     title(req) date?(default:today) budgetHours?
update_capacity_plan     id(req) title? budgetHours? notes?
delete_capacity_plan     id(req)
add_capacity_member      plan_id(req) person_id(req) role? hoursPerDay? workingDays?[str]
remove_capacity_member   plan_id(req) person_id(req)
add_capacity_allocation  plan_id(req) title(req) hours(req) assignee? start_date? end_date?
remove_capacity_allocation plan_id(req) allocation_id(req)
```

---

## MCP Resources (read-only URIs)

```
mdplanner://project   → project config + metadata (JSON)
mdplanner://tasks     → all tasks (JSON)
mdplanner://notes     → all notes (JSON)
mdplanner://goals     → all goals (JSON)
```

---

## Known Gotchas

| Situation | What happens | Fix |
|-----------|-------------|-----|
| `create_milestone` without `project` | Error | Always pass `project` |
| `create_milestone` duplicate (project+name) | Error with existing ID | Check `list_milestones` first |
| `search` without `--cache` | Explicit error response | Start server with `--cache` flag |
| GitHub tools without token | Explicit error response | Configure PAT in Settings |
| `update_meeting` with `actions` | Full replacement, new IDs generated | Pass complete action list |
| `update_fishbone` with `causes` | Full replacement | Pass complete causes list |
| `update_mindmap` with `nodes` | Full replacement | Pass complete node tree |
| `update_portfolio_item` with `urls: undefined` | Omitted key — existing URLs kept | Pass `urls: []` to clear |
| `completed: true` on task | Should only be set by human after testing | Never set this from Claude |
| `update_task` section+assignee combined | Supported since v0.8.1 | Single call is fine |
