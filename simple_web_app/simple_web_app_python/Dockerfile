FROM python:3.13-bullseye

RUN pip install poetry

WORKDIR /app

COPY ./pyproject.toml ./poetry.lock* README.md /app/
COPY ./simple_web_app /app/simple_web_app


RUN poetry install

ENTRYPOINT ["poetry", "run", "simple_web_app"]