
# Technical Interview Task - Simple Web Application

This skeleton of a simple web application (both in Go and Python) is provided to help you get started with the technical interview task.
The additionally Dockerfile and Kubernetes manifests are provided to help you containerize and deploy the application.

## Repository Structure

This folder contains the following:

* `simple_web_app_go`
  * A simple web application implemented in Go using plain net/http package. With a Dockerfile for containerization.
* `simple_web_app_python`
  * A simple web application implemented in Python (poetry project) using FastAPI framework, with a dockerfile for containerization.
* `kubernets_manifests`
  * A Kubernetes manifest file to deploy the web application in a Kubernetes cluster.

