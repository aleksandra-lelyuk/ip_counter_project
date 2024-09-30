import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
import os


def read_stats(file_name):
    # Set display options to show all rows and columns
    #pd.set_option('display.max_rows', None)
    #pd.set_option('display.max_columns', None)

    df = pd.read_csv(file_name)
    file_map = {
        "small.txt": "small",
        "medium.txt": "medium",
        "large.txt": "large"
    }

    if df['fileSize'].iloc[0] not in file_map.keys():
        df.drop(columns = ['fileSize'], inplace=True)
    else:
        df['fileSize'] = df['fileSize'].map(file_map)

    df = df.reset_index(drop=True)

    df['perc_accuracy'] = ((1-np.abs((df['customEstimatedCount'] - df['realCount']))/df['realCount'])*100).round(2)
    df['perc_error'] = (np.abs((df['customEstimatedCount'] - df['realCount']))/df['realCount']*100).round(2)
    df = df[['numHashes', 'hashType', 'customEstimatedCount', 'elapsedTimeSec', 'realCount', 'perc_accuracy', 'perc_error']]

    print('\n\n', 50*'-', '\n', 'RESULTS\n', 50*'-', '\n')
    print(df)
    return df

# Plot results from csvs with ONE dataset size only
def plot_accuracy(dataset, output_file):
    if isinstance(dataset, str):
        df = read_stats(dataset)
        data_size = df['fileSize'].iloc[0]
        plot_name = f"Results for {data_size.upper()} dataset"

    else:
        plot_name = f"Mean accuracy per Hash Function type"
        df = pd.DataFrame(dataset)

    hash_types = df['hashType'].unique().tolist()
    print(hash_types)
    fig, ax1 = plt.subplots()
    ax2 = ax1.twinx()

    if "parseIP" in df.columns:
        df = df[df['parseIP'] == False]
        df.drop(columns=['parseIP'], inplace=True)
        print('Plotting df:\n', df)

    for this_hash_t in hash_types:
        this_hash_df = df[(df['hashType']==this_hash_t)]
        ax1.plot(this_hash_df['numHashes'].astype(str), this_hash_df['perc_accuracy'], label=this_hash_t)
        ax2.scatter(this_hash_df['numHashes'].astype(str), this_hash_df['elapsedTimeSec'], label=this_hash_t)

    plt.title(plot_name)
    ax1.grid()
    ax1.set_xlabel('Number of Registers (hash functions)')
    ax1.set_ylim(top=110)
    ax1.set_ylim(bottom=20)
    ax1.set_ylabel('Accuracy (%)')
    ax2.set_ylabel('Execution Time (s)')
    ax1.legend(title='accuracy', loc='upper left')
    ax2.legend(title='exec time', loc='lower right')

    try:
        plt.savefig(output_file)
        if os.path.exists(output_file):
            print(f"File saved successfully at {output_file}")
        else:
            print(f"Failed to save file at {output_file}")
    except Exception as e:
        print(f"An error occurred while saving the file: {e}")
    #plt.show()

def merge_all_results(small_file, medium_file, large_file):
    df_small = read_stats(small_file)
    df_medium = read_stats(medium_file)
    df_large = read_stats(large_file)

    df_all = pd.concat([df_small, df_medium, df_large], ignore_index=True)
    df_all = df_all[['numHashes', 'hashType', 'perc_accuracy', 'elapsedTimeSec']]

    mean_accuracies = df_all.groupby(by=["hashType", "numHashes"]).mean()
    mean_accuracies = mean_accuracies.reset_index()
    print(mean_accuracies)
    return mean_accuracies


def plot_mean_accuracies_per_hash_fun(small_file, medium_file, large_file, output_file):
    mean_df = merge_all_results(small_file, medium_file, large_file)
    plot_accuracy(mean_df, output_file)




def main():
    import sys

    # plot_mean_accuracies_per_hash_fun("../data/small_test_result.csv", "../data/medium_test_result.csv", "../data/large_test_result.csv", "../data/graphs/average_hash_res.png")

    if len(sys.argv) < 2:
        print("Usage: python print_results.py <function_name> [<file_name>]")
        sys.exit(1)

    function_name = sys.argv[1]

    if function_name == "read_stats":
        if len(sys.argv) != 3:
            print("Usage: python print_results.py read_stats <file_name>")
            sys.exit(1)
        file_name = sys.argv[2]
        read_stats(file_name)
    elif function_name == "plot_accuracy":
        if len(sys.argv) != 4:
            print("Usage: python print_results.py plot_accuracy <file_name> <output_file_name>")
            sys.exit(1)
        file_name = sys.argv[2]
        output_file_name = sys.argv[3]
        plot_accuracy(file_name, output_file_name)
    elif function_name == "plot_mean_accuracies_per_hash_fun":
        if len(sys.argv) != 6:
            print("Usage: python print_results.py plot_mean_accuracies_per_hash_fun <small_file> <medium_file> <large_file> <output_file>")
            sys.exit(1)
        small_file = sys.argv[2]
        medium_file = sys.argv[3]
        large_file = sys.argv[4]
        output_file = sys.argv[4]
        plot_mean_accuracies_per_hash_fun(small_file, medium_file, large_file, output_file)
    else:
        print(f"Unknown function: {function_name}")
        sys.exit(1)

if __name__ == "__main__":
    main()

