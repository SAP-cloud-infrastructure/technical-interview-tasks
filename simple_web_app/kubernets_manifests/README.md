# Deploy workload to Kubernetes

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

5. Validate if Pod `homepage-tester` is working
