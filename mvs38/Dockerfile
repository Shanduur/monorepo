FROM centos:8

RUN yum install -y perl libnsl unzip 
RUN mkdir -p /usr/mvs38
WORKDIR /usr/mvs38

# install hercules
RUN curl -sSLO http://www.hercules-390.org/hercules-3.07-1.x86_64.rpm && rpm -i ./hercules-3.07-1.x86_64.rpm

# install MVS
RUN curl -sSLO http://wotho.ethz.ch/tk4-/tk4-_v1.00_current.zip && unzip tk4-_v1.00_current.zip
RUN curl -sSLO http://wotho.ethz.ch/tk4-/tk4-source.zip && unzip tk4-source.zip
RUN curl -sSLO http://wotho.ethz.ch/tk4-/tk4-cbt.zip && unzip tk4-cbt.zip

# remove downloaded files
RUN rm *.zip *.rpm

# configure 
RUN sed -i 's/NUMCPU \${NUMCPU:=1}/NUMCPU \${NUMCPU:=2}/g' ./conf/tk4-.cnf
RUN sed -i 's/MAXCPU \${MAXCPU:=1}/MAXCPU \${MAXCPU:=2}/g' ./conf/tk4-.cnf
RUN echo 'DAEMON' > /usr/mvs38/unattended/mode

ENTRYPOINT [ "/usr/mvs38/mvs" ]
