# log-based-metrics-exporter

log based metric exporter for kubenetes. Significantly reduce cloud costs by computing the metrics within your cluster.

log based metric exporter will find all pods for the containers, in the namespaces that you provide. It will then measure the occurence of the conditions. All conditions must match for a metric to be incremented, they match based on sub string matches. For example `failed to read database` would match a log entry like `epic service 4923. failed to read database`. It also updates every 15 seconds and addes/removes pods as they change.

See [example](./example/) for how to deploy to kubenetes

## Configuration

```
[
{
  "name": string,
  "metric": string,
  "namespace": []string,
  "container": []string,
  "condition": []string
},
...
]
```

mount your rules directory to `/rules`

## Example configuration

```json
[
  {
    "name": "Epic Service failing to read database",
    "metric": "epic_service_database_read_fail",
    "namespace": ["default"],
    "container": ["my_epic_service"],
    "condition": ["error", "failed to read database"]
  }
]
```

## Example of the metrics

```
log_based_metric{metric="log_entry_4",name="log entry 4",namespace="default",pod="dummy-59b956b77c-2h78c"} 0

log_based_metric{metric="log_entry_4",name="log entry 4",namespace="default",pod="dummy-59b956b77c-99rvv"} 73
```
