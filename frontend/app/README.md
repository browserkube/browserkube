# Browserkube frontend

### How-to use .evn file 
1. Navigate to the /frontend/app directory in your terminal or IDE
2. Create .env file
3. Copy everything from .env.example and paste in created file .evn
4. Save .env file

### How-to use Eslint
1. Open Eslint configuration in your IDE settings
2. Choose Eslint package: node_modules/eslint
3. Choose Eslint configuration file: .eslintrc.js
4. Select running Eslint on save 

### for Development
1. clean the docker builder and image memory with
docker builder prune
docker image prune
in the main terminal
2. check the pods kubectl get pods -n browserkube
there should be only running pods
3. check wether you have already the installed nginx
if not use this command also in the main terminal
helm upgrade --install ingress-nginx ingress-nginx \
  --repo https://kubernetes.github.io/ingress-nginx \
  --namespace ingress-nginx --create-namespace
4. to run [BE] locally use skaffold dev command in the /browserkube folder
5. run npm run build/dev in the frontend/app folder
6. start developing new features.