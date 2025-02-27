Breast Cancer Prediction Dashboard
---
MSDS 434 final project

### Project Structure
- prediction-api/ - Rest API written in Go that returns data from SageMaker
endpoint
- frontend/ - Web UI for app
- misc/ - Jupyter Notebook for training model on SageMaker, script for deploying
sagemaker endpoint, Ec2 dependency installation script

### AWS Configuration
There are several steps that need to be taken on AWS in order to train the model
, deploy the model, and create an Ec2 instance
- Create an AWS account or login to an existing account
- Create an IAM user to manage the application if one does not exist

### Model Training

### Deployment
- Create a key-value pair on AWS and download the key to the development machine
- Create an Ec2 instance on AWS
    - Use the latest version of Ubuntu if following this guide
    - The smallest possible instance should be enough
    - Be sure to assign the key value pair and set access for SSH, HTTP, and
HTTPS traffice from the ip ranges you want
- Copy code from the development to the Ec2 instance
```
scp -R -i <path to pem file> <path to project directory> ubuntu@<public DNS of Ec2 instance>
```
- Login to the Ec2 instance and run the dependency installation script
- Use Docker to start the services on the Ec2 instance
- You should now be able to naviage to the public DNS of your Ec2 instance and
view the app

This is **not** a good long term solution for deploying an application, but it
works well enough for this project.

Given more time, a better approach would be to use Ansible to create a Docker
image of the project, transfer the image to a target server, and then run the
image on that server.
