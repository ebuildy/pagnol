**Execute a collection of HTTP queries** (like postman/newman):
* raw HTTP queries
* elasticsearch: to provision index templates, snapshot repositories, ILM policies
* more is coming!


## usage

```
pagnol --debug --actions ./actions.yaml
```

## actions.yaml

```
connection:
  url: http://localhost:9200

actions:

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

## packaging

pagnol is coming as a Docker image:

```
docker run -v $(PWD)/actions.yaml:/actions.yaml ebuildy/pagnol:v0.1.0 --actions /actions.yaml
```

and as a Helm chart:

```
helm install ebuildy/pagnol 
```
