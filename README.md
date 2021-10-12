# calendar-module-kube

フロントエンドUIがカレンダー情報を参照・更新するためのバックエンドサービスです。


# デプロイ手順

## Git clone
```bash
$ git clone git@bitbucket.org:latonaio/calendar-module-kube.git  
```

## Dockerイメージの作成
```
$ cd /path/to/calendar-module-kube
$ make docker-build
```

## K8sにデプロイ

`k8s/calendar-module-kube.yaml`をエディタで開き、`<username>`、`<password>`の値をセットアップ環境に合わせて修正
```
            - name: DB_USER
              value: <username> # <- 変更
            - name: DB_PASSWORD
              value: <password> # <- 変更
```

修正後、kubectlコマンドでデプロイする
```
$ kubectl apply -f k8s/calendar-module-kube.yaml
```
