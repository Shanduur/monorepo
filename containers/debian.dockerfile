FROM debian:latest

RUN apt-get update && apt-get upgrade -y

RUN apt-get install -y openssh-server

RUN mkdir /root/.ssh && echo > /root/.ssh/authorized_keys

COPY ./configs/etc /etc
COPY ./configs/.profile /root/.profile
COPY ./scripts/entrypoint.sh /bin/entrypoint.sh

ENTRYPOINT [ "/bin/entrypoint.sh" ]
