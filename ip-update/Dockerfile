FROM python:3.9


RUN pip install pipenv

RUN mkdir -p /opt/ipu

WORKDIR /opt/ipu

COPY main.py .
COPY Pipfile.lock .
COPY Pipfile .

RUN pipenv install --system --deploy --ignore-pipfile

ENTRYPOINT [ "python", "main.py" ]
