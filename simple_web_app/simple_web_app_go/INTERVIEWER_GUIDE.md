# INTERVIEWER GUIDE - Complete Assessment Framework

## 🎯 Quick Start for Interviewer

### Before the Interview
1. Review this guide
2. Ensure the skeleton code is working (test with `go run main.go`)
3. Test Docker and Kubernetes setup if doing those portions
4. Decide interview format (live coding vs. take-home)

### During the Interview
1. Give applicant the task (see TASK.md or use the task document provided earlier)
2. Observe their problem-solving approach
3. Use the evaluation rubric below
4. Ask follow-up questions based on skill level

### After the Interview
1. Compare their solution to `handlers_solution.go`
2. Score using rubric below
3. Make hiring decision

---

## 📋 Evaluation Rubric

### Junior Level (Pass) - Score: 50-70%

**Implementation:**
- ✅ Basic implementation works for happy path
- ✅ Makes HTTP request to GitHub API
- ✅ Returns JSON response
- ✅ Code compiles and runs

**What they might miss:**
- ❌ Timeout on HTTP client
- ❌ Proper error handling for all cases
- ❌ Input validation
- ❌ Not closing response body

**Interview questions:**
- "Walk me through your implementation"
- "What does your code do when GitHub is down?"
- "Why did you structure the response this way?"

**Decision:** Consider for junior positions with mentorship

---

### Mid Level (Good) - Score: 71-85%

**Implementation:**
- ✅ All Junior requirements
- ✅ Uses timeout on HTTP client
- ✅ Validates input parameters (empty org name)
- ✅ Closes response body with defer
- ✅ Handles most error cases (404, network failures)
- ✅ Clean, readable code
- ✅ Proper HTTP status codes

**What they might miss:**
- ❌ Doesn't spot NetworkPolicy issue immediately
- ❌ Doesn't discuss rate limiting
- ❌ Limited production considerations

**Interview questions:**
- "What happens if GitHub rate limits us?"
- "How would you test this code?"
- "Walk me through the BasicAuth middleware - what's `subtle.ConstantTimeCompare` doing?"
- "Why use multi-stage Docker build?"

**Decision:** Good fit for mid-level positions

---

### Senior Level (Excellent) - Score: 86-100%

**Implementation:**
- ✅ All Mid Level requirements
- ✅ Comprehensive error handling with appropriate status codes
- ✅ Security considerations (validates input, discusses implications)
- ✅ **Spots and fixes the NetworkPolicy issue in Kubernetes**
- ✅ Discusses production improvements proactively
- ✅ Understands trade-offs (caching, pagination, auth)
- ✅ Can debug issues independently
- ✅ Suggests improvements to existing code

**Production Discussion:**
- ✅ Talks about caching strategy
- ✅ Mentions rate limiting and circuit breakers
- ✅ Discusses observability (logging, metrics, tracing)
- ✅ Considers pagination for large organizations
- ✅ Understands Kubernetes networking and NetworkPolicy

**Interview questions:**
- "How would you make this production-ready?"
- "What observability would you add?"
- "How would you handle thousands of repositories efficiently?"
- "Why isn't the Kubernetes deployment accessible?" (NetworkPolicy test!)
- "What's the security issue with storing credentials in environment variables?"

**Decision:** Strong candidate for senior positions

---

## 🔍 Key Things to Watch For

### 1. HTTP Client Usage (CRITICAL)

**❌ Bad:**
```go
resp, err := http.Get(url)
```
No timeout! This is a **major red flag**.

**✅ Good:**
```go
client := &http.Client{Timeout: 10 * time.Second}
resp, err := client.Get(url)
```

**Why it matters:** Without timeout, slow GitHub responses hang goroutines forever, leading to resource exhaustion.

---

### 2. Response Body Handling (CRITICAL)

**❌ Bad:**
```go
resp, err := client.Get(url)
json.Decode(resp.Body)  // Body never closed!
```

**✅ Good:**
```go
resp, err := client.Get(url)
defer resp.Body.Close()
json.Decode(resp.Body)
```

**Why it matters:** Not closing body causes memory leaks and connection exhaustion.

---

### 3. Error Handling

**❌ Bad:**
```go
resp, _ := client.Get(url)  // Ignoring errors!
var repos []Repo
json.Decode(resp.Body, &repos)
w.Write(repos)
```

**✅ Good:**
```go
resp, err := client.Get(url)
if err != nil {
    http.Error(w, "Failed to connect", http.StatusServiceUnavailable)
    return
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusOK {
    http.Error(w, "GitHub error", http.StatusBadGateway)
    return
}
```

---

### 4. GitHub API Headers (Important)

**❌ Will fail:**
```go
http.Get("https://api.github.com/orgs/golang/repos")
// GitHub returns 403 without User-Agent!
```

**✅ Will work:**
```go
req.Header.Set("User-Agent", "simple-web-app")
```

**Ask:** "Why did you add the User-Agent header?"
**Expected:** "GitHub API requires it or returns 403"

---

### 5. The NetworkPolicy Trap (Senior Test)

In `kubernetes-manifests.yaml`:
```yaml
spec:
  policyTypes:
  - Ingress
  ingress: []  # ⚠️ THIS BLOCKS ALL TRAFFIC
```

**Senior candidates should:**
1. Notice the Service isn't accessible
2. Check NetworkPolicy configuration
3. Identify `ingress: []` as the problem
4. Explain that empty ingress = deny all
5. Fix it or suggest removal

**Questions to ask:**
- "The Kubernetes deployment is running but not accessible. What could be wrong?"
- "What does the NetworkPolicy do?"
- "Why does the homepage-tester pod fail?"

---

## 💬 Discussion Topics by Level

### For All Levels

1. **Code walkthrough**: "Explain your implementation"
2. **Error scenarios**: "What happens if X fails?"
3. **Testing**: "How would you test this?"
4. **Basic security**: "What security concerns exist?"

### For Mid+ Levels

5. **BasicAuth middleware**: "Explain the auth.go implementation"
   - Expected: Discuss `subtle.ConstantTimeCompare` preventing timing attacks
6. **Docker**: "Why multi-stage build?"
   - Expected: Smaller image size, security
7. **Status codes**: "Why 503 vs 502?"
   - Expected: 503 for our service issues, 502 for upstream issues

### For Senior Levels

8. **Production readiness**: "What's missing for production?"
   - Expected: Caching, rate limiting, monitoring, pagination
9. **Kubernetes troubleshooting**: "Deployment won't serve traffic"
   - Expected: Identifies NetworkPolicy issue
10. **Trade-offs**: "Caching vs freshness?"
11. **Scalability**: "How to handle 1000+ repos?"
    - Expected: Pagination, streaming responses, worker pools
12. **Observability**: "What metrics would you track?"
    - Expected: Request latency, error rate, GitHub API latency, cache hit rate

---

## 🎬 Interview Scenarios

### Scenario 1: Live Coding (45-60 min)

```
Time allocation:
- 5 min: Explain task, answer questions
- 30-40 min: Implementation (observe their process)
- 10-15 min: Discussion and follow-ups
```

**Watch for:**
- Do they read existing code first?
- Do they test incrementally?
- How do they debug errors?
- Do they ask clarifying questions?

---

### Scenario 2: Code Review (30-45 min)

Applicant implements beforehand, you review together.

```
Time allocation:
- 10 min: Walkthrough their solution
- 15 min: Deep dive on specific parts
- 10 min: Docker/Kubernetes discussion
```

**Questions:**
- "Why did you choose this approach?"
- "What alternatives did you consider?"
- "What would you change for production?"

---

### Scenario 3: Debugging Exercise (20-30 min)

Give them intentionally broken code, ask them to fix it.

**Bugs to include:**
- No timeout on HTTP client
- Response body not closed
- Missing User-Agent header
- Ignoring errors

**Watch for:**
- How quickly they identify issues
- Their debugging approach
- Tool usage (logging, curl, etc.)

---

## 🚩 Red Flags

### Critical Issues (Immediate rejection)
- ❌ Doesn't handle errors at all
- ❌ Uses panic instead of error returns
- ❌ Can't get basic implementation working
- ❌ Doesn't understand HTTP basics
- ❌ Can't debug simple issues
- ❌ Plagiarizes without understanding

### Warning Signs (Probe deeper)
- ⚠️ No timeout on HTTP client (junior mistake)
- ⚠️ Doesn't close response body (resource leak)
- ⚠️ Ignores status codes (incomplete error handling)
- ⚠️ Over-engineered solution (premature optimization)
- ⚠️ Can't explain their own code

---

## ✅ Green Flags

### Strong Indicators
- ✅ Reads existing code before implementing
- ✅ Tests incrementally during development
- ✅ Asks clarifying questions about requirements
- ✅ Handles errors comprehensively
- ✅ Uses timeouts and closes resources
- ✅ Writes clean, readable code
- ✅ Can explain trade-offs and alternatives
- ✅ Spots issues in existing code (NetworkPolicy!)
- ✅ Discusses production considerations proactively

---

## 📊 Scoring Sheet

| Category | Junior (Pass) | Mid (Good) | Senior (Excellent) |
|----------|---------------|------------|-------------------|
| **Basic Implementation** | Works for happy path | Handles most cases | Comprehensive |
| **Error Handling** | Some errors handled | Most errors handled | All errors, proper codes |
| **Resource Management** | May leak resources | Closes most resources | Perfect cleanup |
| **Security** | Basic awareness | Validates input | Proactive discussion |
| **Code Quality** | Compiles, runs | Clean, readable | Production-ready |
| **Docker** | Can build image | Understands basics | Optimizes & secures |
| **Kubernetes** | Can deploy | Understands manifests | Debugs NetworkPolicy |
| **Production Mindset** | Not mentioned | Some awareness | Comprehensive plan |

---

## 🎯 Decision Matrix

### Hire (Senior)
- ✅ Score 86%+
- ✅ Spots NetworkPolicy issue
- ✅ Discusses production improvements
- ✅ Excellent error handling
- ✅ Strong system design thinking

### Hire (Mid)
- ✅ Score 71-85%
- ✅ Solid implementation
- ✅ Good error handling
- ✅ Clean code
- ⚠️ May miss advanced issues

### Maybe (Junior with mentorship)
- ⚠️ Score 50-70%
- ✅ Basic implementation works
- ⚠️ Needs improvement in error handling
- ⚠️ Limited production awareness
- ✅ Shows potential to learn

### No Hire
- ❌ Score <50%
- ❌ Can't complete basic implementation
- ❌ Poor error handling
- ❌ Can't debug simple issues
- ❌ Doesn't understand fundamentals

---

## 📝 Follow-up Questions Bank

### Error Handling
- "What happens if GitHub returns 429 Too Many Requests?"
- "How would you handle partial failures?"
- "What if JSON parsing fails?"

### Performance
- "How would you handle an org with 1000+ repos?"
- "What about caching?"
- "How would you reduce GitHub API calls?"

### Security
- "What security issues exist in this implementation?"
- "How would you prevent injection attacks?"
- "What's the purpose of `subtle.ConstantTimeCompare`?"

### Architecture
- "How would you structure this for a larger application?"
- "Where would you add middleware?"
- "How would you make this testable?"

### DevOps
- "What metrics would you track?"
- "How would you handle secrets in Kubernetes?"
- "What's wrong with the NetworkPolicy?"

---

## 🎓 Teaching Moments

If candidate is close but missing key concepts, guide them:

**On timeouts:**
> "What happens if GitHub takes 5 minutes to respond?"

**On resource leaks:**
> "What happens to the HTTP connection if we don't close the body?"

**On security:**
> "How does timing-based attack work on password comparison?"

**On NetworkPolicy:**
> "What does an empty ingress list mean in NetworkPolicy?"

---

## 📁 Reference Files

- **`handlers_solution.go`** - Complete reference implementation
- **`SOLUTION_README.md`** - Detailed solution documentation
- **`test_solution.sh`** - Automated testing script
- **`main_solution.go`** - Test runner

---

## ✨ Final Tips

1. **Be flexible**: Adjust difficulty based on candidate level
2. **Listen actively**: Their thought process matters more than perfect code
3. **Give hints**: If stuck, guide them - interview stress is real
4. **Focus on fundamentals**: Error handling, resource management, security
5. **Make it conversational**: This isn't just code - it's about communication

**Remember:** You're looking for someone you want to work with, not just someone who can write code!

---

Good luck with your interview! 🚀
