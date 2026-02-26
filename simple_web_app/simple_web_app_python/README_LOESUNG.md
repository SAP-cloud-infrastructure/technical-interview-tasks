# 🐍 Python Interview Lösung - KOMPAKT

## ✅ Was wurde erstellt

### 📁 Dateien:
1. **`main_solution.py`** - Komplette Lösung mit Kommentaren
2. **`LOESUNG.md`** - Diese Dokumentation
3. **`test_solution.sh`** - Test-Script (falls poetry installiert)

---

## 🎯 Die Lösung (Kern-Code)

```python
@app.get("/repo-list/{org_name}")
async def get_repo_list(org_name: str, response: Response,
                        repo_filter: str | None = None):
    # 1. URL bauen
    url = f"https://api.github.com/orgs/{org_name}/repos"

    # 2. Headers setzen (WICHTIG!)
    headers = {
        "User-Agent": "simple-web-app",  # GitHub braucht das!
        "Accept": "application/vnd.github.v3+json"
    }

    try:
        # 3. Request mit TIMEOUT (WICHTIG!)
        resp = requests.get(url, headers=headers, timeout=10)

        # 4. Fehler behandeln
        if resp.status_code == 404:
            response.status_code = 404
            return {"error": f"Organization '{org_name}' not found"}

        if resp.status_code != 200:
            response.status_code = 502
            return {"error": f"GitHub API error: {resp.status_code}"}

        # 5. JSON parsen
        repos = resp.json()

        # 6. Optional filtern
        if repo_filter:
            repos = [r for r in repos
                    if repo_filter.lower() in r["name"].lower()]

        # 7. Antwort zurückgeben
        return {
            "organization": org_name,
            "count": len(repos),
            "repositories": [
                {
                    "name": r["name"],
                    "full_name": r["full_name"],
                    "description": r.get("description", ""),
                    "html_url": r["html_url"],
                    "stargazers_count": r["stargazers_count"],
                    "language": r.get("language", "")
                }
                for r in repos
            ]
        }

    except requests.Timeout:
        response.status_code = 503
        return {"error": "GitHub API timeout"}

    except requests.RequestException as e:
        response.status_code = 503
        return {"error": f"Connection failed: {str(e)}"}
```

---

## 🧪 Getestet & funktioniert!

```bash
✓ /hello-world → {"message": "Hello World"}
✓ /repo-list/golang → 30 Repositories
✓ /repo-list/golang?repo_filter=tools → 1 Repository gefiltert
✓ /repo-list/fake999 → 404 Error korrekt
```

---

## 🚨 KRITISCHE Punkte beim Bewerber

### ❌ Ohne das → FAIL:
1. **Kein `timeout`** → App hängt ewig
2. **Kein `User-Agent`** → GitHub gibt 403
3. **Kein Error-Handling** → App crashed

### ✅ Mit das → PASS:
1. **timeout=10** gesetzt
2. **User-Agent** Header gesetzt
3. **try/except** für Netzwerk-Fehler
4. **Status Codes** geprüft (404, 503, 502)

---

## 💬 Interview-Fragen

**Frage 1:** "Warum timeout=10?"
→ **Antwort:** Ohne timeout hängt die App bei langsamen Requests

**Frage 2:** "Warum User-Agent?"
→ **Antwort:** GitHub API braucht das, sonst 403 Forbidden

**Frage 3:** "Was bei 1000 Repos?"
→ **Antwort:** Pagination nötig, GitHub gibt nur 30 pro Seite

**Frage 4:** "Production-ready?"
→ **Antwort:** Caching hinzufügen, Rate-Limits beachten, Logging

---

## 📊 Schnell-Bewertung

| Bewerber Typ | Was sie machen |
|--------------|----------------|
| **Junior** | Funktioniert, aber fehlt timeout/error-handling |
| **Mid** | Alles funktioniert, sauberer Code |
| **Senior** | Perfekt + spricht über Caching, Monitoring, Scale |

---

## 🚀 Zum Testen

```bash
# Dependencies installieren
pip install fastapi uvicorn requests

# Server starten
python3 -m simple_web_app.main_solution

# Testen (anderes Terminal)
curl http://localhost:8080/hello-world
curl "http://localhost:8080/repo-list/golang?repo_filter=go"
```

---

## ✨ Das war's!

**Einfach, klar, getestet. Fertig für's Interview!** 🎉

### Wichtigste Dateien:
- `main_solution.py` - Die Lösung
- `LOESUNG.md` - Diese Doku

**Viel Erfolg!**
