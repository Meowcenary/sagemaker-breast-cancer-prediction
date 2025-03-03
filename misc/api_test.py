"""
Script to test prediction API endpoints. Run with: python3 api_test.py
"""
import csv
import requests

import pandas as pd


def test_api(api_endpoint_url, csv_file):
    df = pd.read_csv(csv_file, header=0)
    df.columns = df.columns.str.strip()  # Strip extra spaces in headers

    query_params = {}
    for index, row in df.iterrows():
        for column in df.columns:
            column_key = column.replace(" ", "_")
            value = row[column]
            query_params[column_key] = value

        # Send request
        try:
            response = requests.get(api_endpoint_url, params=query_params)
            print(f"Row {index} Response: {response.status_code}")
            print(response.json())
        except requests.exceptions.RequestException as e:
            print(f"Error on row {index}: {e}")

if __name__ == "__main__":
    test_api(
        "http://ec2-52-91-186-123.compute-1.amazonaws.com/predict",
        "../data/train_with_headers.csv"
    )
