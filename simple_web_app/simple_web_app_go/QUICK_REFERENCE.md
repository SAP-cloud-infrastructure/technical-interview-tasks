# Quick Reference Card for Interviewer

## 📋 3-Minute Prep

### The Task
Implement `/repo-list/{org_name}` endpoint that calls GitHub API and returns repos as JSON.

### What to Look For
1. ✅ Uses HTTP client with **timeout**
2. ✅ **Closes response body** with defer
3. ✅ Handles **errors** properly
4. ✅ Returns correct **HTTP status codes**

### Critical Red Flags
- ❌ No timeout → goroutine leaks
- ❌ No body close → memory leaks
- ❌ No error handling → production disaster

---

## 🎯 Quick Scoring

| Level | Score | Key Indicator |
|-------|-------|---------------|
| **Senior** | 86%+ | Spots NetworkPolicy issue, discusses production |
| **Mid** | 71-85% | Clean code, good error handling |
| **Junior** | 50-70% | Basic working implementation |
| **No hire** | <50% | Can't complete basic task |

---

## 💬 Top 5 Questions

1. "Walk me through your implementation"
2. "What happens if GitHub is down?"
3. "Why did you add the User-Agent header?"
4. "What does `subtle.ConstantTimeCompare` do in auth.go?"
5. "The K8s deployment won't serve traffic. Why?" (NetworkPolicy trap!)

---

## 🧪 Quick Test Commands

```bash
# Start server
go run main.go

# Test basic
curl http://localhost:8080/hello-world

# Test repo-list
curl "http://localhost:8080/repo-list/golang?repo_filter=go"

# Test 404
curl http://localhost:8080/repo-list/fake999
```

---

## 🔑 Must-Have Code Elements

### HTTP Client with Timeout
```go
client := &http.Client{Timeout: 10 * time.Second}
```

### Close Response Body
```go
defer resp.Body.Close()
```

### Check Status Code
```go
if resp.StatusCode != http.StatusOK {
    // handle error
}
```

### User-Agent Header
```go
req.Header.Set("User-Agent", "simple-web-app")
```

---

## 🎓 Level Detection

### Junior Signs
- Basic implementation works
- May forget timeout or close body
- Limited error handling

### Mid Signs
- Clean, working code
- Good error handling
- Understands Docker basics

### Senior Signs
- **Spots NetworkPolicy issue**
- Discusses caching, rate limits
- Production-ready thinking
- Security awareness

---

## 🚨 The NetworkPolicy Trap

File: `kubernetes_manifests/kubernetes-manifests.yaml`

```yaml
ingress: []  # ⚠️ This blocks ALL traffic!
```

Ask: "Why isn't the K8s service accessible?"

Senior candidate should identify this!

---

## 📁 Solution Files

- `START_HERE.md` - Overview (this location)
- `INTERVIEWER_GUIDE.md` - Full guide
- `handlers_solution.go` - Reference code
- `test_solution.sh` - Run tests

---

## ⚡ Time Management

**45-min interview:**
- 5 min: Explain task
- 30 min: Code (observe)
- 10 min: Discussion

**30-min code review:**
- 10 min: Walkthrough
- 15 min: Deep questions
- 5 min: K8s trap

---

## ✅ Decision Framework

**Hire as Senior:** 
- Production-ready code + spots K8s issue + system design thinking

**Hire as Mid:**
- Clean working code + good error handling + understands basics

**Maybe (Junior):**
- Basic implementation + shows learning potential

**No Hire:**
- Can't complete task + poor fundamentals + can't debug

---

**Pro Tip:** Focus on problem-solving approach over perfect syntax. Interview stress is real!
