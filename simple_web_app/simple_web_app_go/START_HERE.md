# 📦 Complete Interview Solution Package

## ✅ Solution Status: TESTED & WORKING

All endpoints have been tested and are functioning correctly:
- ✅ `/hello-world` - Returns "Hello World" message
- ✅ `/repo-list/{org_name}` - Fetches repositories from GitHub
- ✅ `/repo-list/{org_name}?repo_filter=X` - Filters repositories
- ✅ `/protected` - BasicAuth protected endpoint
- ✅ Error handling (404, 503, 502, 400)

---

## 📁 Files Created for You

### 🎯 For the Interviewer (YOU)

1. **`INTERVIEWER_GUIDE.md`** ⭐ **START HERE**
   - Complete evaluation rubric
   - Red flags and green flags
   - Question bank by skill level
   - Scoring matrix
   - Decision framework

2. **`handlers_solution.go`** (in `internal/handlers/`)
   - Complete implementation with extensive comments
   - Production-ready code
   - Shows all best practices

3. **`main_solution.go`**
   - Test runner that uses the solution
   - For quick testing

4. **`test_solution.sh`**
   - Automated test script
   - Tests all endpoints and error cases
   - Validates responses

5. **`SOLUTION_README.md`**
   - Detailed solution documentation
   - Test results
   - Docker/Kubernetes instructions
   - Production improvements discussion

---

## 🚀 Quick Start for Interview

### Before the Interview (5 min)

1. **Read the interviewer guide:**
   ```bash
   cat INTERVIEWER_GUIDE.md
   ```

2. **Test the solution (optional):**
   ```bash
   ./test_solution.sh
   ```

3. **Prepare your questions** based on the level you're hiring for

---

### During the Interview

#### Give the applicant this task:

> **Task:** Implement the missing `/repo-list/{org_name}` endpoint in `internal/handlers/handlers.go`
>
> **Requirements:**
> - Fetch repositories from GitHub API for the given organization
> - Support optional `?repo_filter=` query parameter to filter by name
> - Return JSON response with repository list
> - Handle errors appropriately
>
> **Endpoints to implement:**
> - `GET /repo-list/{org_name}` - List all repos for an organization
> - `GET /repo-list/{org_name}?repo_filter=go` - Filter repos by name
>
> **Time:** 30-45 minutes for implementation
>
> **Test with:**
> ```bash
> go run main.go
> curl "http://localhost:8080/repo-list/golang?repo_filter=go"
> ```

---

### Evaluation Process

1. **Watch them code** (or review their solution if take-home)
2. **Use the rubric** in `INTERVIEWER_GUIDE.md`
3. **Ask follow-up questions** from the question bank
4. **Compare to** `handlers_solution.go`
5. **Score and decide** using the decision matrix

---

## 📊 Quick Evaluation Checklist

### Critical Must-Haves (All levels)
- [ ] Basic implementation works
- [ ] Makes HTTP request to GitHub
- [ ] Returns JSON response
- [ ] Code compiles and runs

### Mid-Level Expectations
- [ ] Uses HTTP client with timeout
- [ ] Validates input (empty org name)
- [ ] Closes response body with defer
- [ ] Handles errors (network, 404, etc.)
- [ ] Proper HTTP status codes

### Senior-Level Expectations
- [ ] All Mid-Level items
- [ ] Comprehensive error handling
- [ ] Security considerations
- [ ] **Spots NetworkPolicy issue in Kubernetes**
- [ ] Discusses production improvements
- [ ] Clean, production-ready code

---

## 🎯 Key Things to Watch For

### ❌ Critical Red Flags
1. **No timeout on HTTP client** → Goroutine leaks
2. **Response body not closed** → Memory leaks
3. **No error handling** → Production nightmare
4. **Missing User-Agent header** → GitHub API fails

### ✅ Green Flags
1. **Reads existing code first** → Good practice
2. **Tests incrementally** → Good workflow
3. **Asks clarifying questions** → Good communication
4. **Discusses trade-offs** → Senior mindset
5. **Spots NetworkPolicy trap** → Strong Kubernetes knowledge

---

## 🧪 Testing the Solution

### Option 1: Automated Test (Recommended)
```bash
cd simple_web_app_go
./test_solution.sh
```

### Option 2: Manual Test
```bash
# Start server
go run main_solution.go

# Test in another terminal
curl http://localhost:8080/hello-world
curl "http://localhost:8080/repo-list/golang?repo_filter=go"
curl http://localhost:8080/repo-list/nonexistent999
```

### Option 3: Docker Test
```bash
docker build -t simple-web-server .
docker run --rm -p 8080:8080 simple-web-server
curl http://localhost:8080/hello-world
```

---

## 💡 Interview Tips

### For Junior Candidates
- Focus on: Does it work? Do they understand basic HTTP?
- Be patient with errors, guide them if stuck
- Look for potential and willingness to learn

### For Mid-Level Candidates
- Focus on: Error handling, resource management, code quality
- Ask about Docker and Kubernetes basics
- Expect clean, working code

### For Senior Candidates
- Focus on: Production readiness, trade-offs, system design
- **Test with NetworkPolicy issue** in Kubernetes
- Expect proactive discussion of improvements
- Ask about observability, scalability, security

---

## 🎬 Interview Formats

### Format 1: Live Coding (45-60 min)
```
5 min  - Explain task
35 min - Implement (observe)
10 min - Discussion
10 min - Docker/K8s (optional)
```

### Format 2: Code Review (30-45 min)
```
10 min - Walkthrough their solution
15 min - Deep dive questions
10 min - Production discussion
5 min  - Kubernetes debugging
```

### Format 3: Take-Home + Review (60 min)
```
Before: They implement at home
30 min - Code review and questions
20 min - Docker and Kubernetes
10 min - System design discussion
```

---

## 🔥 The NetworkPolicy Trap (Senior Test)

In `kubernetes_manifests/kubernetes-manifests.yaml`:

```yaml
spec:
  policyTypes:
  - Ingress
  ingress: []  # ⚠️ BLOCKS ALL TRAFFIC
```

**The Test:**
> "Your Kubernetes deployment is running, but the service isn't accessible. What could be wrong?"

**Expected Answer (Senior):**
1. Checks pods are running: `kubectl get pods`
2. Checks service: `kubectl get svc`
3. Tests connectivity: `kubectl run -it debug --image=curlimages/curl --rm`
4. **Discovers NetworkPolicy blocks ingress**
5. Explains: `ingress: []` means "deny all"
6. Fixes or suggests removing NetworkPolicy

**If they don't spot it:** They might not be senior-level for Kubernetes.

---

## 📚 Reference Implementation

The complete solution in `handlers_solution.go` demonstrates:

✅ **Go Best Practices**
- HTTP client with timeout
- defer for cleanup
- Proper error handling
- Struct tags for JSON

✅ **Security**
- Input validation
- Timeout prevents DoS
- Discusses rate limiting

✅ **Production Considerations**
- Comprehensive error messages
- Proper HTTP status codes
- Ready for caching/monitoring
- Detailed comments

✅ **Code Quality**
- Clean, readable
- Well-structured
- Extensive documentation
- Easy to maintain

---

## 🎓 Common Candidate Mistakes

### Mistake 1: No Timeout
```go
// ❌ BAD - Will hang forever if GitHub is slow
http.Get(url)

// ✅ GOOD
client := &http.Client{Timeout: 10 * time.Second}
client.Get(url)
```

### Mistake 2: Not Closing Body
```go
// ❌ BAD - Memory leak
resp, _ := client.Get(url)
json.Decode(resp.Body)

// ✅ GOOD
resp, err := client.Get(url)
defer resp.Body.Close()
```

### Mistake 3: Ignoring Errors
```go
// ❌ BAD - Ignoring errors
resp, _ := client.Get(url)

// ✅ GOOD
resp, err := client.Get(url)
if err != nil {
    http.Error(w, "Error", 503)
    return
}
```

### Mistake 4: Wrong Status Codes
```go
// ❌ BAD - Always returns 200
w.Write([]byte(`{"error":"not found"}`))

// ✅ GOOD - Proper status code
w.WriteHeader(http.StatusNotFound)
w.Write([]byte(`{"error":"not found"}`))
```

---

## 🏆 Hiring Decision Framework

| Score | Level | Decision |
|-------|-------|----------|
| 86-100% | Senior | ✅ **Strong Hire** - Production-ready code, spots K8s issues |
| 71-85% | Mid | ✅ **Hire** - Solid fundamentals, clean code |
| 50-70% | Junior | ⚠️ **Maybe** - Basic understanding, needs mentorship |
| <50% | - | ❌ **No Hire** - Lacks fundamentals |

---

## 📞 Need Help?

All reference materials are in this directory:
- `INTERVIEWER_GUIDE.md` - Complete guide
- `SOLUTION_README.md` - Solution details
- `handlers_solution.go` - Reference code
- `test_solution.sh` - Automated tests

---

## ✨ Final Checklist

Before the interview:
- [ ] Read `INTERVIEWER_GUIDE.md`
- [ ] Test the skeleton code works
- [ ] Prepare Docker/K8s environment (optional)
- [ ] Choose interview format
- [ ] Prepare follow-up questions

During the interview:
- [ ] Give clear task description
- [ ] Observe problem-solving approach
- [ ] Take notes on key decisions
- [ ] Ask follow-up questions
- [ ] Test K8s debugging (for senior)

After the interview:
- [ ] Compare to reference solution
- [ ] Score using rubric
- [ ] Make hiring decision
- [ ] Provide feedback

---

**Good luck with your interview! You have everything you need to conduct a thorough, fair evaluation.** 🚀

**Remember:** You're hiring a teammate, not just evaluating code. Look for communication, problem-solving, and growth potential!
