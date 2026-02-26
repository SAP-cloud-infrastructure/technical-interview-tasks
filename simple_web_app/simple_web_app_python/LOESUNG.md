# Python Interview Lösung - EINFACH ERKLÄRT

## 📋 Die Aufgabe

Implementiere den fehlenden `/repo-list/{org_name}` Endpoint der:
1. GitHub API aufruft
2. Repositories einer Organisation holt
3. Optional nach Name filtert
4. JSON zurückgibt

---

## ✅ Lösung: main_solution.py

### Die wichtigsten Punkte:

#### 1. **Timeout ist PFLICHT**
```python
resp = requests.get(url, timeout=10)
```
❌ **Ohne timeout** → App hängt wenn GitHub langsam ist!

#### 2. **User-Agent Header setzen**
```python
headers = {"User-Agent": "simple-web-app"}
```
❌ **Ohne User-Agent** → GitHub gibt 403 zurück!

#### 3. **Fehler behandeln**
```python
try:
    resp = requests.get(...)
    if resp.status_code == 404:
        return {"error": "Not found"}
except requests.Timeout:
    return {"error": "Timeout"}
```

#### 4. **Filtern (optional)**
```python
if repo_filter:
    repos = [r for r in repos if repo_filter.lower() in r["name"].lower()]
```

---

## 🧪 Testen

### Automatisch:
```bash
chmod +x test_solution.sh
./test_solution.sh
```

### Manuell:
```bash
# Server starten
poetry install
poetry run python -m simple_web_app.main_solution

# In anderem Terminal testen:
curl http://localhost:8080/hello-world
curl "http://localhost:8080/repo-list/golang?repo_filter=go"
```

---

## 🎯 Worauf achten beim Bewerber?

### ✅ MUSS-Kriterien:
1. **Timeout gesetzt** → sonst hängt die App
2. **User-Agent gesetzt** → sonst funktioniert GitHub nicht
3. **Fehler behandelt** → try/except für Netzwerk-Fehler
4. **Status Codes geprüft** → 404, 503, etc.

### ❌ Häufige Fehler:
- Kein timeout → **KRITISCH**
- Kein User-Agent → funktioniert nicht
- Keine Error-Handling → App crasht
- Status Codes ignoriert

---

## 📊 Bewertung

### Junior (50-70%)
- ✅ Grundimplementierung funktioniert
- ⚠️ Fehlt: timeout, Error-Handling

### Mid (71-85%)
- ✅ Alles funktioniert
- ✅ Gutes Error-Handling
- ✅ Sauberer Code

### Senior (86%+)
- ✅ Alles perfekt
- ✅ Spricht über Caching, Rate Limits
- ✅ Kennt Production Best Practices

---

## 💡 Diskussionsfragen

1. **"Warum timeout=10?"**
   - Antwort: Sonst hängt die App bei langsamen Requests

2. **"Warum User-Agent Header?"**
   - Antwort: GitHub braucht das, sonst 403

3. **"Was bei 1000+ Repos?"**
   - Antwort: Pagination, GitHub gibt nur 30 zurück

4. **"Production-ready machen?"**
   - Antwort: Caching, Logging, Monitoring, Rate-Limit-Handling

---

## 🚀 Erwartete Antwort

```python
@app.get("/repo-list/{org_name}")
async def get_repo_list(org_name: str, response: Response,
                        repo_filter: str | None = None):
    url = f"https://api.github.com/orgs/{org_name}/repos"
    headers = {"User-Agent": "simple-web-app"}

    try:
        resp = requests.get(url, headers=headers, timeout=10)

        if resp.status_code == 404:
            response.status_code = 404
            return {"error": "Not found"}

        if resp.status_code != 200:
            response.status_code = 502
            return {"error": "GitHub error"}

        repos = resp.json()

        if repo_filter:
            repos = [r for r in repos
                    if repo_filter.lower() in r["name"].lower()]

        return {
            "organization": org_name,
            "count": len(repos),
            "repositories": repos
        }

    except requests.Timeout:
        response.status_code = 503
        return {"error": "Timeout"}
```

---

## ✨ Fertig!

Die Lösung ist **einfach, klar und production-ready**. Alle wichtigen Punkte sind abgedeckt.

**Viel Erfolg beim Interview!** 🎉
