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


@app.get("/repo-list/{org_name}")
async def get_repo_list(org_name: str, response: Response, repo_filter: str | None = None):
    pass


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
