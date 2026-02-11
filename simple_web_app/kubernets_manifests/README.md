# Deploy workloads to Kubernetes

## Steps:

1. Clone the repository

2. Build the docker image from `simple_web_app_*` folder

3. Tag the image:

   ```bash
   docker tag docker.io/library/simple-web-server localhost:5001/simple_web_app
   ```

3. Create kind cluster with registry using:

   ```bash
   ./kind-with-registry.sh
   ```

4. Deploy manifests to Kubernetes

## Tasks

1. Validate if `homepage-tester` is working, and fix if there are issues

2. Change `homepage-tester` to curl `/protected` service endpoint

3. Implement `/repo-list/{org_name}[?repo_filter=filter]` service endpoint and test it with the `homepage-tester` Pod to list all repositories from `SAP-cloud-infrastructure` github org
