"""
The model itself was built using AWS Canvas and then registered with Sagemaker Studio. The
code here is to create a serverless endpoint that can invoke the model for predictions. It
needs to be run from Sagemaker Studio as a Jupyter notebook.
"""
import boto3

sm_client = boto3.client("sagemaker")

# list models that exist - if the model is registered, it will have the listed name and
# some id string appended e.g canvas-predict-breast-cancer-diganosis-1740119171497
# response = sm_client.list_models()
# for model in response["Models"]:
#     print(model["ModelName"])

model_name = "canvas-predict-breast-cancer-diganosis-1740119171497"
config_name = "breast-cancer-prediction-serverless-config"
# create serverless endoint configuration
response = sm_client.create_endpoint_config(
    EndpointConfigName=config_name,
    ProductionVariants=[
        {
            "VariantName": "AllTraffic",
            "ModelName": model_name,
            "ServerlessConfig": {
                "MemorySizeInMB": 1024,  # Options: 1024, 2048, 3072, 4096
                "MaxConcurrency": 2      # Increase depending on load, 2 seems fine
            }
        }
    ]
)
print("Created Endpoint Config:", response)

# create the serverless endpoint
# EnddpointConfigName should match what was passed for the first block
response = sm_client.create_endpoint(
    EndpointName="breast-cancer-prediction-endpoint",
    EndpointConfigName=config_name
)
print("Endpoint creation started")

# list endpoints
response = sm_client.list_endpoints()
for end_point in response["Endpoints"]:
    print(end_point["EndpointName"], "-", end_point["EndpointStatus"])

# detailed failure reason
endpoint_name = "breast-cancer-prediction-endpoint"
# response = sm_client.describe_endpoint(EndpointName=endpoint_name)
response = sm_client.describe_endpoint(EndpointName=endpoint_name)
print("Failure Reason:", response.get("FailureReason", "No failure reason provided."))

print("Endpoint Status:", response["EndpointStatus"])
print("Endpoint ARN:", response["EndpointArn"])
print("Endpoint Config:", response["EndpointConfigName"])
