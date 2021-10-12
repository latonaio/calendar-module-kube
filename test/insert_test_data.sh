#!/bin/sh
SCRIPT_DIR=$(cd $(dirname $0); pwd)
echo $SCRIPT_DIR
kubectl exec -it $(kubectl get pods | grep mysql | cut -d ' ' -f 1) -- /bin/sh -c "mysql -u root -ppassword --default-character-set=utf8" < ${SCRIPT_DIR}/../sql/calendar-module-kube.sql
kubectl exec -it $(kubectl get pods | grep mysql | cut -d ' ' -f 1) -- /bin/sh -c "mysql -u root -ppassword --default-character-set=utf8" < ${SCRIPT_DIR}/sql/test_dataset.sql
