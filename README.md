# K8s Pruner

## Description

k8s-pruner is a command-line tool to list and prune unused Kubernetes resources. It helps to clean up resources in a Kubernetes cluster to save costs and improve performance. This tool provides an easy-to-use interface to manage and prune resources such as Pods, ConfigMaps, Secrets, PVCs, and more.

## Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/manthan-parmar-1998/k8s-pruner.git
    ```
2. Navigate into the project directory:
    ```bash
    cd k8s-pruner
    ```
3. Run `go mod tidy` to install dependencies:
    ```bash
    go mod tidy
    ```
4. Build the binary:
    ```bash
    go build -o k8s-pruner
    ```

## Usage

To list unused resources:

```bash
./k8s-pruner list
```

To prune unused resources:

```bash
./k8s-pruner prune
```

## 🔍 Feature Comparison: `k8s-pruner` vs Alternatives

| Feature                            | k8s-pruner | kubectl-gc  | KubeJanitor     | Pluto | kube-cleanup-operator |
| ---------------------------------- | ---------- | ----------- | --------------- | ----- | --------------------- |
| CLI-based interface                | ✅         | ✅          | ❌              | ✅    | ❌                    |
| Prune unused ConfigMaps            | ✅         | ❌          | ✅              | ❌    | ✅                    |
| Prune unused Secrets               | ✅         | ❌          | ✅              | ❌    | ✅                    |
| Prune stale PVCs                   | ✅         | ❌          | ✅              | ❌    | ✅                    |
| Detect & clean Jobs/Completed Pods | ✅         | ✅          | ✅              | ❌    | ✅                    |
| Namespace targeting                | ✅         | ❌          | Partial         | ❌    | ✅                    |
| `--dry-run` support                | ✅         | ❌          | ❌              | ❌    | ❌                    |
| `--age` flag                       | ✅         | ❌          | ✅              | ❌    | ✅                    |
| Safe delete with confirmation      | ✅         | ❌          | ❌              | ❌    | ❌                    |
| CronJob integration                | Planned    | ❌          | ✅              | ❌    | ✅                    |
| `kubectl` plugin-compatible        | Planned    | ❌          | ❌              | ❌    | ❌                    |
| Actively maintained                | ✅         | ⚠️ Inactive | ⚠️ Low activity | ✅    | ⚠️ Inactive           |

> ✅ = Fully supported | ❌ = Not supported | ⚠️ = Limited or outdated

## Contributing

Contributions are welcome! Please submit a pull request for any improvements or new features.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
