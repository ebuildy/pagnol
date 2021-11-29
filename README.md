# pagnol

Execute a collection of HTTP queries (like postman/newman):
* raw HTTP queries
* elasticsearch: to provision index templates, snapshot repositories, ILM policies
* more is coming!

[![Helm chart](https://img.shields.io/github/v/release/ebuildy/pagnol?display_name=tag&label=Helm%20chart)]()  [![Github](https://img.shields.io/github/issues-raw/ebuildy/pagnol?style=flat-square)](https://github.com/ebuildy/pagnol/issues)

## usage

```
export PAGNOL_TARGET_USERNAME=admin
export PAGNOL_TARGET_PASSWORD=elastic
pagnol --actions actions.yaml --tls-ca /my-ca.crt --url https://hot-es/ --ignore-error
```

## actions.yaml

```
---

- name: get cluster health
  kind: http
  spec:
    method: get
    url: /_cluster/health
  asserts:
  - status: 200

- name: conso
  kind: org.elasticsearch/index_template
  spec:
    index_patterns:
    - consolidation-*
    - conso-*
    template:
      settings:
        index:
          similarity:
            default:
              type: boolean
        refresh_interval: 30s
        number_of_shards: 1
```

## CLI options

```
GLOBAL OPTIONS:
   --debug                    if true, log is verbose (default: false)
   --verbose                  if set, log is very verbose (default: false)
   --dry-run                  if set, nothing is sent (default: false)
   --ignore-error             if set, no stop actions if error occured (default: false)
   --actions value, -a value  YAML actions
   --username value           (default: "admin ") [$PAGNOL_TARGET_USERNAME]
   --password value           (default: "elastic") [$PAGNOL_TARGET_PASSWORD]
   --tls-ca value              [$PAGNOL_TARGET_TLS_CA]
   --url value                 [$PAGNOL_TARGET_URL]
   --tls-no-verify            (default: false) [$PAGNOL_TARGET_TLS_NO_VERIFY]
   --help, -h                 show help (default: false)
```

## packaging

pagnol is coming as a Docker image:

```
docker run -v $(PWD)/actions.yaml:/actions.yaml ghcr.io/ebuildy/pagnol:latest --url hot-es --actions /actions.yaml
```

and as a Helm chart:

```
helm install ebuildy/pagnol
```
