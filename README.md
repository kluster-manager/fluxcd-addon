https://open-cluster-management.io/developer-guides/addon/
https://github.com/open-cluster-management-io/addon-framework
https://github.com/kluster-management/addon-contrib/tree/main


```bash
> kubebuilder init --domain open-cluster-management.io --skip-go-version-check
> kubebuilder create api --group fluxcd --version v1alpha1 --kind FluxCDConfig
```
