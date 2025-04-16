kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Wait for ArgoCD components to be ready
kubectl wait --namespace argocd \
  --for=condition=Available \
  --timeout=5m \
  deployment/argocd-redis deployment/argocd-repo-server deployment/argocd-applicationset-controller deployment/argocd-server deployment/argocd-dex-server deployment/argocd-notifications-controller

kubectl patch deployment argocd-server -n argocd \
--type='json' \
-p='[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--insecure"}]'

# Wait for ArgoCD components to be ready
kubectl wait --namespace argocd \
  --for=condition=Available \
  --timeout=5m \
  deployment/argocd-server

kubectl apply -f ingress.yaml

PASSWORD=$(kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d)
echo "Username: admin and Password: $PASSWORD"
