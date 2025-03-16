Breast Cancer Prediction API
---
This project provides code for training and deploying a predictive model with
Amazon Web Services' SageMaker, a simple REST API written in Go to allow for
live predictions, resource and traffic monitoring with Promteheus and Grafana,
and a docker-compose.yml file to orchestrate each component as a microservice.

### Project Structure
- docker-compose.yml - Defines the services the project uses
- prediction-api/ - Rest API written in Go that returns data from SageMaker
endpoint
- misc/ - Jupyter Notebook for training model on SageMaker, script for deploying
sagemaker endpoint, Ec2 dependency installation script
- data/ - The data used to train the model in libsvm and csv format. The
SageMaker model training script imports the data from a Python library, but it
is included here for debugging or general viewing

### AWS Configuration
There are several steps that need to be taken on AWS in order to train the model
, deploy the model, and create an Ec2 instance
- Create an AWS account or login to an existing account
    - Assign two-factor authentication device for root device
- Create an IAM user to manage the application if one does not exist
- Assign policies/roles to IAM user for all AWS services such SageMaker, Ec2, S3
, etc.

### Model Training
To save on costs, the demonstration project deploys the SageMaker model to an
endpoint that will only incur charges when called directly from the API. To
deploy the model to an endpoint, a domain must be created.
- Navigate to Amazon SageMaker AI dashboard
- Select "Domains" from the menu on the left under "Admin configurations" and
click "Create domain"
- Create the domain for a single user

Once the domain has been created, the model can be trained and deployed.
- Navigate to Amazon SageMaker AI dashboard
- Select "Studio" from the menu on the left under "Applications and IDEs" and
open the studio
- Select "JuptyerLab" at the the top left Applications dashboard and click
"Create JupyterLab space"
- Copy the code from "misc/create_and_deploy_model.py" to the Juptyer Notebook
updating the resource names and other configuration values to reflect the
current instance
- Run the code to train and deploy the model

When the model has deployed successfully, move on to deployment.

### Deployment
- Create a key-value pair on AWS and download the key to the development machine
- Create an Ec2 instance on AWS
    - Use the latest version of Ubuntu if following this guide
    - The smallest possible instance should be enough for demonstrations, but a
larger instance may be necessary for production
    - Be sure to assign the key value pair and set access for SSH, HTTP, and
HTTPS traffice from the ip ranges you want. For demonstrations, Prometheus and
Grafana were run on port 9090 and 3000 respectively. This can be changed to any
other unused port and the configuration updated accordingly
- Copy code from the development to the Ec2 instance
```
scp -R -i <path to pem file> <path to project directory> ubuntu@<public DNS of Ec2 instance>
```
- Login to the Ec2 instance and run the dependency installation script
- Use Docker compose to start the services on the Ec2 instance
- You should now be able to naviage to the public DNS of your Ec2 instance and
view the app

### Updating the Public Domain
To update the domain from the default Ec2 public domain, it is recommended to
use the AWS tool Route 53.

The general process for registering a domain with Route 53 and assigning it to
an Ec2 is as follows:
- Navigate to the Route 53 dashboard on AWS
- Under "Domains", click "Registered domains"
- Click the "Register domains" button and find a suitable domain that is
available
- Navigate to "Hosted zones" and create a hosted zone with the domain you would
like to use e.g "spsdatascience.com"
- Within the freshly created hosted zone, click the "Create record" button
- Select "Simple routing", then click the "Define simple record" button
- On the form, select a simple subdomain such as "www", set the record type to A
, and set "Value/Route traffic" to "IP address or another value..." with the IP
value set to the public IPv4 address of the Ec2 instance as listed on the Ec2
dashboard

### Monitoring
- Prometheus is used to monitor resource usage and gather usage statistics. It
runs on port 9090, but this can be changed in docker-compose.yml.
- Grafana is available on port 3000. This can also be changed in
docker-compose.yml

Prometheus can be added to Grafana as a data source to visualize metrics that it
has gathered. To do this, login to Grafana, navigate to dashboard creation,
create a new dashboard, create a new visuzalition, when prompted to add a data
source, select Prometheus, and finally set the endpoint to
http://prometheus:9090 to reflect the name it will have on the Docker network.

Maintenance
---
Use the below command to clear Docker resources to free up space
```
docker system prune --volumes -a
```
