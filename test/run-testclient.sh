#!/bin/sh
SCRIPT_DIR=$(cd $(dirname $0); pwd)
echo $SCRIPT_DIR
kubectl exec -it $(kubectl get pods | grep mysql | cut -d ' ' -f 1) -- /bin/sh -c "mysql -u root -ppassword --default-character-set=utf8" < ${SCRIPT_DIR}/../sql/calendar-module-kube.sql
kubectl delete -f ${SCRIPT_DIR}/k8s/deployment.yaml
kubectl apply -f ${SCRIPT_DIR}/k8s/deployment.yaml
