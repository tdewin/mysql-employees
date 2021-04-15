FROM ubuntu
WORKDIR /usr/sbin/
COPY mysql-employees /usr/sbin/mysql-employees
RUN chmod 755 /usr/sbin/mysql-employees
EXPOSE 8080
CMD /usr/sbin/mysql-employees