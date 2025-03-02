"""
Calculate the median value of each column of a CSV file
"""
import pandas as pd

def compute_statistics(csv_file):
    df = pd.read_csv(csv_file)

    # Compute median and mean for each column
    stats = df.median().to_frame(name='Median')
    stats['Mean'] = df.mean()

    # Determine column width for formatting
    col_width = max(len(col) for col in stats.index) + 2

    print(f"{'Column'.ljust(col_width)} {'Median'.ljust(15)} {'Mean'.ljust(15)}")
    print("-" * (col_width + 30))
    for column in stats.index:
        print(f"{column.ljust(col_width)} {str(stats.at[column, 'Median']).ljust(15)} {str(stats.at[column, 'Mean']).ljust(15)}")

if __name__ == "__main__":
    csv_file = "../data/train_with_headers.csv"  # Change this to your actual CSV file
    compute_statistics(csv_file)
