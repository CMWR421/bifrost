package complexity

// Default semantic-routing exemplars are intentionally balanced across tiers.
// Each tier contains 25 coding, 10 general-knowledge, 8 math/reasoning, and
// 7 writing prompts. Topic ladders place related vocabulary in every tier so
// nearest-neighbor classification learns requested work rather than topic.

var defaultSimpleExemplars = []string{
	// Coding and software engineering (25).
	"How do I convert a string to an integer in Go?",
	"Write a Python function that returns the largest number in a non-empty list.",
	"What does write-ahead logging mean in PostgreSQL?",
	"I need the number of paid orders. Give me the SQL query for an `orders` table with a `status` column.",
	"Rename the parameter `n` to `count` in `function double(n) { return n * 2; }`, update the declaration and function body, and keep the behavior unchanged.",
	"How do I read the `Authorization` header from an HTTP request in Go?",
	"Write the PostgreSQL statement to add an index on the `email` column of the `users` table.",
	"Using Go's standard library, how do I check whether an error wraps `context.Canceled`?",
	"What `if` condition makes a GitHub Actions job run only for pull request events?",
	"Explain the difference between a cache hit and a cache miss.",
	"In an OpenAI-compatible chat request, which field specifies the model to use?",
	"Show a Go `slog` statement that records a `request_id` and an error.",
	"Write a shell command that lists every `.json` file under the current directory.",
	"Write a React button that disables itself while a form is submitting.",
	"What HTTP status code should an API return when a requested resource does not exist?",
	"What is an MCP tool call?",
	"Replace email addresses in a string with `[REDACTED]` using a regular expression.",
	"What does p95 latency mean?",
	"Show the npm command to update `lodash` to the latest version allowed by `package.json`.",
	"What is the difference between request latency and token cost?",
	"How do I compute a SHA-256 checksum for a file in Go?",
	"What does a WebSocket ping frame do?",
	"What is a model context window?",
	"Convert `2026-08-01T12:00:00Z` to India Standard Time.",
	"In an OpenAI Chat Completions SSE stream, what does `data: [DONE]` indicate?",

	// General knowledge and explanation (10).
	"What is evaporation, and what causes it?",
	"What was the Industrial Revolution?",
	"What does inflation mean?",
	"What is antibiotic resistance?",
	"What does separation of powers mean?",
	"What is a carbon footprint?",
	"What is spaced practice in learning?",
	"What happens when I tap a card at checkout?",
	"What does gross margin tell a business?",
	"What does end-to-end encryption protect?",

	// Math, logic, and data reasoning (8).
	"An $80 jacket is discounted by 15%. What is its price before tax?",
	"What is the probability of rolling a total of 8 with two fair six-sided dice?",
	"What is the median of 4, 7, 9, 12, and 18?",
	"All administrators are employees. Priya is an administrator. Is Priya an employee?",
	"Monthly active users increased from 4,000 to 4,600. What was the percentage increase?",
	"Weekly sales were 90, 100, and 110 units. Forecast next week using their three-week average.",
	"Vendor A charges $48 for 12 units; Vendor B charges $65 for 20. Which has the lower unit price?",
	"Conversion was 10% for control and 13% for treatment. What is the absolute percentage-point lift?",

	// Writing and content transformation (7).
	"Correct the grammar without changing the meaning: \"The reports is ready and it include last months numbers.\"",
	"Summarize in one sentence: \"The library closes at 6 p.m. Friday for electrical maintenance. The return slot and online catalog remain available, scheduled events move online, and the building reopens at 9 a.m. Saturday.\"",
	"Rewrite for a nontechnical reader: \"The API rejected the request because its authentication token had expired.\"",
	"Turn these notes into three bullets: \"Launch Tuesday. Owner: Maya. Final security review is Monday.\"",
	"Rewrite in plain language: \"Reimbursement requests submitted after thirty calendar days shall be ineligible for processing.\"",
	"Write a short release-note title for: \"Users can now export dashboard results as CSV files.\"",
	"Improve this status update: \"Payments are broke for some users. We are looking into it.\"",
}

var defaultMediumExemplars = []string{
	// Coding and software engineering (25).
	"We need a Go worker pool that runs at most five jobs. After cancellation, reject new jobs, let active jobs finish, and discard queued work.",
	"Fix this React effect so it fetches when `userId` changes and ignores an older response that finishes later: `useEffect(() => { load(userId).then(setUser) })`.",
	"Add cursor pagination to a REST endpoint ordered by ascending ID. The first page has no cursor; later cursors expire after one hour and invalid ones are rejected.",
	"A single service owns a two-million-row table. Plan a PostgreSQL migration that adds nullable `status`, sets `pending` as the default for new writes before backfilling existing nulls in batches, verifies none remain, and then enforces `NOT NULL`.",
	"Build a Python API client that makes one initial request and at most three retries for HTTP 429 or 503, using exponential backoff under one five-second overall deadline.",
	"Add API-key authentication to a Go endpoint. Hash keys before lookup, reject missing, invalid, or revoked keys without logging the raw key, and continue only for active keys.",
	"This query has separate indexes on `customer_id` and `created_at` but still sorts: `SELECT * FROM orders WHERE customer_id = $1 ORDER BY created_at DESC LIMIT 50`. Write a replacement index, name a candidate index for removal, and explain which query-plan and usage evidence must be checked first.",
	"A Go test increments one shared counter from two goroutines and fails under `-race`. Fix the race and make the final count reliably equal 2,000.",
	"Update a Go CI workflow to run three test packages in parallel matrix jobs, cache Go modules, and cancel older runs for the same branch.",
	"Add a five-minute in-memory cache around `load(ctx, key)` in Go. Do not cache failures, and coalesce concurrent misses for the same key.",
	"Add image input conversion for one model provider, preserving text and image order while rejecting unsupported media types with a clear error.",
	"Instrument a checkout handler with latency, status counters, and trace correlation while keeping customer data out of labels and logs.",
	"Build a CLI command that recursively validates YAML files, prints file-and-line errors, and exits nonzero if any file is invalid.",
	"Build a React settings form with client validation, an unsaved-change warning, and server errors attached to the correct fields.",
	"Define one JSON error response for an API, map validation and authorization failures to it, and include a request ID for support.",
	"Expose two auto-executable MCP tools to a model, validate their arguments, and return ordinary execution failures as tool results so the agent loop can continue.",
	"Add a request guardrail that blocks secrets, redacts email addresses, and reports which rule acted without persisting the original values.",
	"Profile a Go endpoint that slowed after a release, determine whether CPU or allocations dominate, and verify one targeted optimization with benchmarks.",
	"Upgrade a Go library with one deprecated API, update its callers, and add tests proving behavior remains unchanged.",
	"Choose between two models using a latency budget and estimated input and output token cost, then record the inputs and reason for each routing decision.",
	"Upload large files to object storage in parts, retry failed parts, and verify the completed object against a client-provided checksum.",
	"Implement a WebSocket client that reconnects with backoff, restores each subscription once, and detects dead connections using ping and pong frames.",
	"Reject requests that exceed a model's context window after counting serialized system, user, assistant, and tool content plus reserved output tokens; return the remaining input budget.",
	"Schedule a report for each user's local 2:30 a.m. On a skipped daylight-saving time run at the next valid time; on a repeated time run once; store the next run in UTC.",
	"Convert one provider's text deltas into OpenAI Chat Completions SSE chunks, honor `stream_options.include_usage`, and terminate with `data: [DONE]`.",

	// General knowledge and explanation (10).
	"Explain how evaporation, condensation, and precipitation interact in the water cycle, then predict how less vegetation could change runoff and groundwater recharge.",
	"How did factory production change work and city life during the Industrial Revolution?",
	"Why can higher interest rates reduce inflation, and why does the effect take time?",
	"Explain how unnecessary antibiotic use increases resistance and why that affects people who did not take the medicine.",
	"How do legislatures, executives, and courts check one another when a law is made and enforced?",
	"Compare reducing emissions directly with buying carbon offsets, including the main limitations of each.",
	"Compare spaced practice with cramming, and explain when retrieval practice should be added.",
	"Explain the difference between payment authorization and settlement, including why a transaction can remain pending.",
	"Compare gross margin with contribution margin and explain when each is useful.",
	"Explain how key exchange and device verification let two people communicate without the service reading their messages.",

	// Math, logic, and data reasoning (8).
	"A loan charges a $120 fee plus 6% simple annual interest on $4,000 for 18 months. Calculate the total repayment and effective cost as a percentage of principal.",
	"A test detects 95% of defective parts and correctly clears 90% of good parts. If 1% are defective, what is the probability that a flagged part is defective?",
	"Delivery times were 2, 3, 3, 4, 18 days before a change and 2, 3, 3, 4, 5 after it. Compare the means and medians, then explain whether the improvement is broad or outlier-driven.",
	"Schedule design, legal review, security review, and launch from Monday through Thursday, one per day. Design and legal must precede launch; security can occur anytime. Design is unavailable Tuesday, legal Monday, and security Wednesday. List every valid schedule and explain the eliminated cases.",
	"January had 80 activations from 200 signups, and 40 of those activated users were retained. February had 135 activations from 300 signups, and 45 were retained. Calculate both activation and post-activation retention rates and explain what changed.",
	"Monthly signups from January through April were 100, 120, 140, and 160. Fit a linear trend, forecast May and June, calculate the residuals, and state one reason the forecast could fail.",
	"Product A earns $30 and uses 2 machine hours; Product B earns $40 and uses 4. With 20 hours and demand capped at 6 units each, find the integer production mix with the highest profit.",
	"Control converted 40 of 200 users and treatment converted 55 of 200. Calculate absolute and relative lift, then explain whether the sample alone justifies a rollout.",

	// Writing and content transformation (7).
	"Rewrite this support reply to acknowledge frustration, explain that retries delay processing, preserve the Friday deadline, give the next action, and state when escalation is appropriate: \"Please stop retrying and wait for our email.\"",
	"Summarize in three bullets: \"Revenue grew 12%, mainly from new customers. Repeat purchases fell 6%. Marketing spend rose 20%, and gross margin declined from 64% to 60%.\" Separate growth, retention, and efficiency, then add a neutral headline without implying causation.",
	"Rewrite for customer administrators: \"API version 1 retires on October 1. Version 2 changes pagination tokens but not authentication. Administrators should update integrations and test before September 15.\" Include impact and next steps.",
	"Turn these notes into a brief decision memo: \"Option A launches in two weeks and costs $20,000 but lacks audit logs. Option B takes six weeks and costs $35,000 with complete auditing. Compliance requires audit logs this quarter.\"",
	"Write a four-question employee FAQ from this policy: \"Remote employees may expense up to $500 for equipment in any rolling 24-month period. Manager approval is required before purchase. Phones and recurring subscriptions are excluded.\"",
	"Write release notes for saved filters, faster CSV exports, and a deletion warning. Group them under Added, Improved, and Safety; state who is affected and any required action without inventing claims.",
	"Draft a customer update from this timeline: \"09:10 UTC login failures began; 09:25 traffic was shifted; 09:40 error rates recovered. Some sessions still require a new login. Root cause is under investigation.\"",
}

var defaultComplexExemplars = []string{
	// Coding and software engineering (25).
	"A service has outgrown a single database and must adopt sharding while traffic continues. How should reads and writes move so rollback remains possible and users retain read-your-writes consistency?",
	"A shared job queue serves tenants with different quotas. How should it enforce limits without starving low-volume tenants when retries and worker failures increase load?",
	"Our API uses JWTs signed by keys replicated across regions. Rotate keys without invalidating unexpired tokens, including how verification keys propagate before any region begins signing with the new key.",
	"Original and fallback providers can race to stream the first chunk. Ensure exactly one wins without buffering the response.",
	"While a rollout is running, both a CLI and web dashboard can edit its deployment configuration. How should conflicting revisions be stored, audited, and restored without leaving rollout state inconsistent?",
	"Migrate a public API from static keys to short-lived JWTs while both credentials work during rollout. Preserve revocation and audit history under a stable client identity.",
	"A multi-tenant orders table must support recent orders by customer, unpaid orders by tenant, and lookup by external ID. Design indexes and any partitioning or workload controls needed to limit write cost and noisy-neighbor impact.",
	"Parallel integration tests share PostgreSQL and Redis, and some depend on wall-clock expiry. Redesign isolation and time control while keeping CI time within 20% of today.",
	"Several repositories deploy services independently while sharing a database schema. Design a CI/CD release workflow that checks version compatibility, orders schema and service rollouts, blocks unsafe combinations, and supports rollback when one deployment fails.",
	"Authorization decisions are cached in three regions, and policy updates carry increasing versions. A disconnected region may fail closed or use cached decisions for at most 60 seconds. Design invalidation and reconnection while preserving tenant isolation.",
	"Define a canonical request and response model for text, images, and tool calls across three incompatible providers, preserving content order, call-result identity, usage, and provider error metadata without breaking existing adapters.",
	"Trace partial failures across a gateway and three services without exposing tenant data or creating unusably high-cardinality metrics.",
	"A CLI updates hundreds of configuration files while other processes may edit them. Design atomic writes, conflict detection, crash recovery, and a reversible batch operation.",
	"A collaborative React editor must merge offline edits after reconnecting while showing remote changes live. Design synchronization, conflict handling, and undo semantics.",
	"Multiple services and SDKs depend on an API error contract. Evolve it without breaking older clients while preserving retry behavior, localization, and incident diagnostics.",
	"MCP clients expose overlapping tools and permissions. Coordinate discovery, filtering, cancellation, execution, and audit across concurrent agent sessions.",
	"Streaming guardrails must inspect partial output, stop unsafe content, and emit only safe chunks while persisting encrypted reversible mappings for authorized log-detail reveal, without restoring originals to the live stream or leaking buffers across requests.",
	"Regional traffic shifts trigger latency spikes while gateway CPU stays flat and downstream 429s rise. Design a safe diagnosis and mitigation plan across load balancers, messaging, databases, and downstream limits.",
	"Upgrade a shared serialization library across independently deployed services when its old and new wire formats are incompatible. Plan compatibility, rollout, rollback, and data repair.",
	"Route across providers with changing prices, rate limits, failures, and quality targets without overspending or oscillating between models.",
	"A versioned object store must satisfy regional residency and legal deletion while retaining disaster-recovery backups. Design lifecycle, replication, auditing, and proof of deletion.",
	"A real-time gateway serves millions of WebSocket clients across regions. Design connection ownership, message fan-out, failover, and session recovery during deployments.",
	"Tool schemas, calls, results, conversation history, and cached prefixes all consume context, while fallback models have different context windows. Truncate safely without dropping required instructions or corrupting tool-call sequences.",
	"A global scheduler runs recurring jobs across regions. Prevent duplicate execution during failover while preserving local-time semantics, missed-run policy, and audit history.",
	"Normalize interleaved text, reasoning, and tool deltas into ordered Responses events with stable item and sequence indices; stop cleanly on client cancellation and emit one valid terminal outcome on provider completion or failure.",

	// General knowledge and explanation (10).
	"A farming region gets similar annual rainfall but longer droughts and heavier storms. Explain how warming, land use, and groundwater pumping could interact, and what evidence would separate their effects.",
	"Historians disagree on whether the Industrial Revolution initially improved living standards. Evaluate the claim across income, health, housing, gender, and time period.",
	"An economy has high inflation, weak growth, and rising unemployment. Compare monetary and fiscal responses, their distributional effects, and the uncertainty around timing.",
	"A hospital has rising resistant infections but limited lab capacity. Evaluate how rapid testing, prescribing rules, isolation, and staffing constraints should shape its response.",
	"During a public emergency, officials need speed while courts and legislatures protect rights and accountability. Compare oversight designs and their likely failure modes.",
	"A city has limited funds for cutting emissions across transport, buildings, and electricity. Evaluate a portfolio that balances impact, equity, feasibility, and uncertain forecasts.",
	"A school wants to adopt spaced practice across subjects. Design a pilot that measures durable learning without disadvantaging students with limited study time or device access.",
	"A merchant sees duplicate pending charges after network interruptions. Trace likely failure points across the terminal, payment network, and bank, then identify useful evidence.",
	"A cash-constrained subscription company subsidizes hardware upfront. Evaluate whether to raise the device price, monthly fee, or both using margin, churn, support cost, inventory risk, and acquisition payback.",
	"A messaging service needs multi-device access, device revocation, and history recovery without weakening end-to-end encryption. Compare designs and the security tradeoffs users would face.",

	// Math, logic, and data reasoning (8).
	"Equipment costs $90,000 now plus $8,000 per used year and resells for $55,000, $42,000, or $30,000 after years one, two, or three. The project ends after year one with 40% probability; otherwise, it ends after year two with 25%. A cancellable lease costs $35,000 yearly. At 8%, compare expected present costs and sensitivity to duration probabilities.",
	"Fraud affects 0.5% of 100,000 daily transactions. A low threshold catches 98% of fraud and flags 10% of legitimate payments; a high threshold catches 80% and flags 1%. Reviews cost $4, missed fraud costs $300, and reviewers can handle 1,500 flags. Choose a threshold and quantify the trade-off.",
	"Checkout A converted 520 of 10,000 users and B converted 570 of 10,000. On mobile, A was 400 of 8,000 and B was 420 of 7,000; on desktop, A was 120 of 2,000 and B was 150 of 3,000. Within 30 days, 20 converted A users and 45 converted B users received refunds. Assess whether B should launch.",
	"Access requires an active workforce relationship and a trusted device. Employees qualify while employed; contractors qualify only with an active sponsor. Break-glass access may bypass device trust but never the workforce requirement, while a sanctions match always denies access. Define precedence, identify conflicts, and produce the final decision logic.",
	"Treatment A succeeded for 81 of 87 small cases and 192 of 263 large cases. Treatment B succeeded for 234 of 270 small cases and 55 of 80 large cases. Compare overall and segment rates, explain the reversal, and decide whether to deploy one treatment globally or choose by case size, stating what additional evidence is needed.",
	"Weekly demand was 100, 102, 98, 101, 140, and 145 units; a promotion started in week five. For week seven, demand has a 60% chance of staying at 145 and otherwise returning to 100. Capacity is 125, excess costs $2 per unit, and shortages cost $8. Choose production and explain the forecast risk.",
	"Staff day, evening, and night with two nurses each and at least one ICU nurse. Per-shift options are A(ICU, day/evening, $220), B(ICU, day/night, $240), C(all, $180), D(ICU, evening/night, $230), and E(all, $170). Day-evening and evening-night are adjacent; day-night is not. No nurse may work adjacent shifts or more than two. Minimize cost.",
	"A feature was first enabled at 20 high-volume stores while 80 others remained untreated. Customers visit multiple stores, and a holiday campaign begins in week four. Design an analysis that separates feature impact from store selection, spillovers, and seasonality.",

	// Writing and content transformation (7).
	"Legal rejects claims that a beta is completely secure or can never lose data; engineering says closing the app before sync can lose changes; support needs backup steps; product wants sync delays disclosed. Produce a customer notice and internal support guidance that keep guarantees, risks, and mitigations consistent.",
	"Write an executive decision brief from these facts: \"Sales rose 12%; new customers rose 25%; repeat purchases fell 6%; promotional spend doubled; gross margin fell four points.\" Separate facts from hypotheses, explain the conflicting signals, and define evidence gates before extending promotions.",
	"Draft messages for employees, customers, and local officials from these facts: \"A plant closes December 31; 120 roles are affected; severance depends on tenure; customer orders transfer to another site; union negotiations remain open.\" Keep commitments and unknowns consistent while tailoring detail and action.",
	"Write a recommendation from these notes: \"Sales wants launch in May; security needs four more weeks; a June competitor launch is expected; the pilot improved conversion 8% but doubled support tickets; phased rollout is available.\" Compare options, state assumptions, and define rollback criteria.",
	"Exports are limited to approved regions; contracts get 90 days from renewal; regions may change; legal may grant time-limited exceptions; blocked exports need audit records and appeals. Create a rollout announcement, support decision tree, and legal escalation guide with consistent rules.",
	"Pilot summaries reduced review time 30% on average but varied widely. They are English-only, may omit details, require approval, retain data seven days, and can be disabled. Write launch, administrator, and in-product messages that qualify the claim and keep controls consistent.",
	"Monitoring links checkout errors to a 14:05 rollout, while support reports onset at 13:58. Rollback ended 14:31; 3% of attempts failed, with no duplicates. Draft customer and internal reports that share facts, surface timeline uncertainty, separate cause from hypothesis, and leave missing owners or deadlines explicit.",
}

func sharedTierDefaults(keywords, exemplars []string) []string {
	combined := make([]string, 0, len(keywords)+len(exemplars))
	combined = append(combined, keywords...)
	return append(combined, exemplars...)
}
