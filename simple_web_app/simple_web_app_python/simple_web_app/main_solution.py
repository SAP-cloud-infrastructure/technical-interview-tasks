import secrets
from fastapi import FastAPI, Response
import uvicorn
import requests
import os
from fastapi.security import HTTPBasic, HTTPBasicCredentials
from typing import Annotated

from fastapi import Depends, HTTPException, status


app = FastAPI()
security = HTTPBasic()


@app.get("/hello-world")
async def hello_world():
    return {"message": "Hello World"}


# ============================================================================
# LÖSUNG: /repo-list/{org_name} Endpoint
# ============================================================================
@app.get("/repo-list/{org_name}")
async def get_repo_list(org_name: str, response: Response, repo_filter: str | None = None):
    """
    Holt Repositories von GitHub für eine Organisation.

    Parameter:
    - org_name: Name der GitHub Organisation (z.B. "golang")
    - repo_filter: Optional - Filter für Repository-Namen

    Beispiel: /repo-list/golang?repo_filter=go
    """

    # SCHRITT 1: GitHub API URL erstellen
    url = f"https://api.github.com/orgs/{org_name}/repos"

    # SCHRITT 2: HTTP Request mit Timeout senden
    # Wichtig: Immer timeout setzen, sonst hängt die App!
    headers = {
        "User-Agent": "simple-web-app",  # GitHub braucht User-Agent!
        "Accept": "application/vnd.github.v3+json"
    }

    try:
        # timeout=10 -> nach 10 Sekunden abbrechen
        resp = requests.get(url, headers=headers, timeout=10)

        # SCHRITT 3: Fehler behandeln
        if resp.status_code == 404:
            response.status_code = 404
            return {"error": f"Organization '{org_name}' not found"}

        if resp.status_code != 200:
            response.status_code = 502  # Bad Gateway
            return {"error": f"GitHub API returned status: {resp.status_code}"}

        # SCHRITT 4: JSON parsen
        repos = resp.json()

        # SCHRITT 5: Optional filtern
        if repo_filter:
            # Case-insensitive Suche
            filter_lower = repo_filter.lower()
            repos = [r for r in repos if filter_lower in r["name"].lower()]

        # SCHRITT 6: Antwort zurückgeben
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
        response.status_code = 503  # Service Unavailable
        return {"error": "GitHub API timeout"}

    except requests.RequestException as e:
        response.status_code = 503
        return {"error": f"Failed to connect to GitHub: {str(e)}"}


def get_current_username(
    credentials: Annotated[HTTPBasicCredentials, Depends(security)],
):
    current_username_bytes = credentials.username.encode("utf8")
    correct_username_bytes = os.environ["AUTH_USERNAME"].encode("utf8")
    is_correct_username = secrets.compare_digest(
        current_username_bytes, correct_username_bytes
    )
    current_password_bytes = credentials.password.encode("utf8")
    correct_password_bytes = os.environ["AUTH_PASSWORD"].encode("utf8")
    is_correct_password = secrets.compare_digest(
        current_password_bytes, correct_password_bytes
    )
    if not (is_correct_username and is_correct_password):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Incorrect username or password",
            headers={"WWW-Authenticate": "Basic"},
        )
    return credentials.username


@app.get("/protected")
def protected(username: Annotated[HTTPBasicCredentials, Depends(get_current_username)]):
    return {"message": f"Login successful for {username}"}


def start():
    uvicorn.run(app="simple_web_app.main:app", host="0.0.0.0", port=8080, reload=True)


if __name__ == "__main__":
    start()
