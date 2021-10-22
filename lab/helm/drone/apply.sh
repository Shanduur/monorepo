helm install --name drone-release stable/drone
helm upgrade drone --reuse-values \
    --set 'service.type=NodePort' \
    --set 'service.nodePort=32000' \
    --set 'sourceControl.provider=github' \
    --set 'sourceControl.github.clientID=0cb6e52e26befec35a3e' \
    --set 'sourceControl.secret=drone-server-secrets'  \
    stable/drone
