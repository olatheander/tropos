# Tropos
Cloud-native development on Kubernetes utilizing Visual Studio Code's Remote Development features

## Pod security policies
```
securityContext:
    privileged: true
    capabilities:
    add:
        - SYS_ADMIN
```
