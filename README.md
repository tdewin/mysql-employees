# mysql-employees

Demo app to get a fake app in kubernetes. Use mysql helm chart to setup the db. Make sure to stick to mysql-demo or change the files

```
helm repo add bitnami https://charts.bitnami.com/bitnami
kubectl create namespace mysql-demo
helm install mysql-demo bitnami/mysql --namespace=mysql-demo
```

create the frontend
```
kubectl -n mysql-demo apply -f https://raw.githubusercontent.com/tdewin/mysql-employees/main/configmap.yaml
kubectl -n mysql-demo apply -f https://raw.githubusercontent.com/tdewin/mysql-employees/main/deployment.yaml
kubectl -n mysql-demo apply -f https://raw.githubusercontent.com/tdewin/mysql-employees/main/svc.yaml
```

check if mysql is deployed by checking that all pods are ready (1/1)
```
kubectl -n mysql-demo get pod
```

now start the db create job
```
kubectl -n mysql-demo apply -f https://raw.githubusercontent.com/tdewin/mysql-employees/main/initjob.yaml
```

you can check the logs now
```
kubectl -n mysql-demo get pod
kubectl -n mysql-demo logs mysql-employees-init-[id]
```

delete the initjob
```
kubectl -n mysql-demo delete -f https://raw.githubusercontent.com/tdewin/mysql-employees/main/initjob.yaml
```

check the svc
```
kubectl -n mysql-demo get svc mysql-employees-svc -o wide 
```


manual db as a reference
```
CREATE DATABASE IF NOT EXISTS employees;
USE employees;

CREATE TABLE IF NOT EXISTS employees (
    emp_no      INT             NOT NULL,
    birth_date  DATE            NOT NULL,
    first_name  VARCHAR(14)     NOT NULL,
    last_name   VARCHAR(16)     NOT NULL,
    gender      ENUM ('M','F','X')  NOT NULL,
    hire_date   DATE            NOT NULL,
    PRIMARY KEY (emp_no)
);
INSERT INTO `employees` VALUES (10001,'1953-09-02','Georgi','Facello','M','1986-06-26'),
(10002,'1964-06-02','Bezalel','Simmel','F','1985-11-21'),
(10003,'1959-12-03','Parto','Bamford','M','1986-08-28'),
(10004,'1954-05-01','Chirstian','Koblick','M','1986-12-01'),
(10005,'1955-01-21','Kyoichi','Maliniak','M','1989-09-12'),
(10006,'1953-04-20','Anneke','Preusig','F','1989-06-02'),
(10007,'1957-05-23','Tzvetan','Zielinski','F','1989-02-10'),
(10008,'1958-02-19','Saniya','Kalloufi','M','1994-09-15'),
(10009,'1952-04-19','Sumant','Peac','F','1985-02-18'),
(10010,'1963-06-01','Duangkaew','Piveteau','F','1989-08-24'),
(10011,'1953-11-07','Mary','Sluis','F','1990-01-22'),
(10012,'1960-10-04','Patricio','Bridgland','M','1992-12-18');
```

cleanup
```
kubectl delete ns mysql-demo
```