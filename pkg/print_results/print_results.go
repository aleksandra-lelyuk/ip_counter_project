package print_results

import (
	"log"
	"os"
	"os/exec"
)

// StartReadData starts the read_data() function in the pake_plots.py given a file name
func PrintResultDataExperiment(fileName string) error {
	cmd := exec.Command("python", "../../print_results/print_results.py", "read_stats", fileName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Error running command: %v", err)
		return err
	}

	//log.Println(cmd.Run())
	return nil
}

func PlotResultDataExperiment(fileName string, outputFile string) error {
	cmd := exec.Command("python", "../../print_results/print_results.py", "plot_accuracy", fileName, outputFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Error running command: %v", err)
		return err
	}

	//log.Println(cmd.Run())
	return nil
}
