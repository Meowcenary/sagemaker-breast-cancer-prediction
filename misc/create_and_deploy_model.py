"""
This Python script is the result of combining the codeblocks from a Jupyter Notebook used
with AWS Sagemaker Studio to train a breast cancer prediction model. The code blocks are
marked by comments with "---" and a quick description of the block.
"""
# --- Imports
import boto3
import pandas
import sagemaker
from sklearn.datasets import dump_svmlight_file, load_breast_cancer
from sklearn.model_selection import train_test_split

# --- Load dataset
print("Loading Breast Cancer Wisconsin dataset from sklearn.datasets")
data = load_breast_cancer()
df = pandas.DataFrame(data.data, columns=data.feature_names)
df["target"] = data.target  # Binary classification (0: Malignant, 1: Benign)

# Split data
train, test = train_test_split(df, test_size=0.2, random_state=42)

# Save to CSV (SageMaker requires files in S3)
# train.to_csv("train.csv", index=False, header=False)
dump_svmlight_file(train, "train.libsvm")
# test.to_csv("test.csv", index=False, header=False)
dump_svmlight_file(train, "test.libsvm")
print("Saved training and testing data to csv files")

# --- Save dataset to S3
# Save data to s3 bucket - only necessary once
s3_bucket = "breast-cancer-ml-prediction"
prefix = "breast-cancer-xgboost"

# Initialize s3 client
s3 = boto3.client("s3")

# Upload files
s3.upload_file("train.libsvm", s3_bucket, f"{prefix}/train.libsvm")
s3.upload_file("test.libsvm", s3_bucket, f"{prefix}/test.libsvm")
print(f"Data uploaded to s3://{s3_bucket}/{prefix}/")

# --- Create container
# Get SageMaker session and role
print("Creating Sagemaker session")
sagemaker_session = sagemaker.Session()
role = "arn:aws:iam::476114121159:role/service-role/AmazonSageMaker-ExecutionRole-20250213T120500"

# Define the built-in XGBoost container
print("Creating sagemaker container")
container = sagemaker.image_uris.retrieve(
    framework="xgboost",
    region="us-east-1",
    version="1.5-1"
)

# Define the training job using XGBoost
print("Creating training job")
xgb_estimator = sagemaker.estimator.Estimator(
    image_uri=container,
    role=role,
    instance_count=1,
    instance_type="ml.m5.large",  # Use ml.m5.xlarge for larger datasets
    output_path=f"s3://{s3_bucket}/{prefix}/output",
    sagemaker_session=sagemaker_session
)

# Set XGBoost hyperparameters
xgb_estimator.set_hyperparameters(
    objective="binary:logistic",  # Binary classification
    num_round=100,  # Number of boosting rounds
    max_depth=5,  # Depth of each tree
    eta=0.2,  # Learning rate
    eval_metric="auc"  # Use AUC for performance evaluation
)

# --- Train model
# Specify input data location in S3
train_input = f"s3://{s3_bucket}/{prefix}/train.csv"

# Start training
train_input = sagemaker.inputs.TrainingInput(f"s3://{s3_bucket}/{prefix}/train.csv", content_type="text/csv")
xgb_estimator.fit({"train": train_input})
print("Model is training...")

# --- Deploy endpoint
predictor = xgb_estimator.deploy(
    initial_instance_count=1,
    instance_type="ml.m5.large",  # Or use serverless deployment
)
print("Model deployed")
