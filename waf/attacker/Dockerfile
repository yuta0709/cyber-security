FROM python:3.13-bullseye

WORKDIR /app

COPY Pipfile Pipfile.lock ./

RUN pip install pipenv

RUN pipenv install  --system --deploy

COPY scripts/ ./

CMD ["bash"]