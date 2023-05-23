package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	downloadURL     = "https://github.com/taylormonacelli/deliverhalf/releases/latest/download/deliverhalf_Windows_x86_64.zip"
	nssmDownloadURL = "https://nssm.cc/release/nssm-2.24.zip"
	targetDirectory = "C:\\Program Files\\Deliverhalf"
	tempDirectory   = "%TEMP%"
	programdataDir  = "C:\\Programdata\\deliverhalf"
	exeFilePath     = targetDirectory + "\\deliverhalf.exe"
	nssmFilePath    = targetDirectory + "\\nssm.exe"
)

func main() {
	stopService()
	createDirectory(targetDirectory)
	createDirectory(programdataDir)

	checkAndDownloadFile(downloadURL, filepath.Join(tempDirectory, "deliverhalf_Windows_x86_64.zip"), 1)
	extractZipFile(filepath.Join(tempDirectory, "deliverhalf_Windows_x86_64.zip"), targetDirectory)

	fmt.Printf("Deliverhalf executable extracted to: %s\n", exeFilePath)

	checkAndDownloadFile(nssmDownloadURL, filepath.Join(tempDirectory, "nssm-2.24.zip"), 1)
	installNssm(filepath.Join(tempDirectory, "nssm-2.24.zip"), filepath.Join(tempDirectory, "nssm-2.24\\win64\\nssm.exe"), targetDirectory)

	updateConfig()
	installService()
}

func stopService() {
	fmt.Println("Stopping deliverhalf service...")
	cmd := exec.Command("sc", "stop", "deliverhalf")
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error stopping deliverhalf service:", err)
	} else {
		fmt.Println("Stopped deliverhalf service.")
	}
}

func createDirectory(directoryPath string) {
	err := os.MkdirAll(directoryPath, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory:", err)
	}
}

func checkAndDownloadFile(url, filePath string, expiryDays int) {
	fileInfo, err := os.Stat(filePath)
	if err == nil && !fileInfo.ModTime().Before(time.Now().AddDate(0, 0, -expiryDays)) {
		fmt.Println("File already up to date:", filePath)
		return
	}

	fmt.Println("Downloading file:", url)
	err = downloadFile(url, filePath)
	if err != nil {
		fmt.Println("Error downloading file:", err)
	}
}

func downloadFile(url, filePath string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func updateConfig() {
	fmt.Println("Updating deliverhalf config...")
	err := os.Chdir(tempDirectory)
	if err != nil {
		fmt.Println("Error changing directory:", err)
		return
	}

	cmd := exec.Command(exeFilePath, "config", "fetch")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error updating config:", err)
		return
	}

	tempConfigPath := filepath.Join(tempDirectory, ".deliverhalf.yaml")
	programdataConfigPath := filepath.Join(programdataDir, ".deliverhalf.yaml")

	if fileExists(tempConfigPath) {
		oldConfig := programdataConfigPath + "." + time.Now().Format("01-02-2006_15_04_05") + ".yaml"
		err = os.Rename(programdataConfigPath, oldConfig)
		if err != nil {
			fmt.Println("Error renaming config file:", err)
		}

		err = os.Rename(tempConfigPath, programdataConfigPath)
		if err != nil {
			fmt.Println("Error moving config file:", err)
		}

		err = os.MkdirAll(programdataDir, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating directory:", err)
		}
	}

	if fileExists(programdataConfigPath) {
		oldConfig := programdataConfigPath + "." + time.Now().Format("01-02-2006_15_04_05") + ".yaml"
		err = os.Rename(programdataConfigPath, oldConfig)
		if err != nil {
			fmt.Println("Error renaming config file:", err)
		}

		err = os.MkdirAll(programdataDir, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating directory:", err)
		}
	}
}

func extractZipFile(zipFilePath, destinationPath string) error {
	reader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		filePath := filepath.Join(destinationPath, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		writer, err := os.Create(filePath)
		if err != nil {
			return err
		}

		source, err := file.Open()
		if err != nil {
			writer.Close()
			return err
		}

		if _, err := io.Copy(writer, source); err != nil {
			writer.Close()
			source.Close()
			return err
		}

		writer.Close()
		source.Close()
	}

	return nil
}

func installNssm(nssmZipFilePath, nssmSourcePath, targetDirectory string) {
	fmt.Println("Installing NSSM...")
	err := extractZipFile(nssmZipFilePath, tempDirectory)
	if err != nil {
		fmt.Println("Error extracting NSSM zip file:", err)
		return
	}

	copyFile(nssmSourcePath, nssmFilePath)
}

func copyFile(sourcePath, destinationPath string) {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		fmt.Println("Error opening source file:", err)
		return
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		fmt.Println("Error creating destination file:", err)
		return
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		fmt.Println("Error copying file:", err)
	}
}

func installService() {
	fmt.Println("Installing deliverhalf service...")
	cmd := exec.Command(nssmFilePath, "install", "deliverhalf", exeFilePath)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error installing deliverhalf service:", err)
		return
	}

	cmd = exec.Command(nssmFilePath, "set", "deliverhalf", "Start", "SERVICE_AUTO_START")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error setting deliverhalf service start mode:", err)
	}

	cmd = exec.Command(nssmFilePath, "set", "deliverhalf", "DisplayName", "deliverhalf")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error setting deliverhalf service display name:", err)
	}

	cmd = exec.Command(nssmFilePath, "set", "deliverhalf", "AppDirectory", programdataDir)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error setting deliverhalf service application directory:", err)
	}

	cmd = exec.Command(nssmFilePath, "set", "deliverhalf", "AppParameters", fmt.Sprintf(`client send forever --config "%s\.deliverhalf.yaml"`, programdataDir))
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error setting deliverhalf service application parameters:", err)
	}

	cmd = exec.Command("sc", "start", "deliverhalf")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error starting deliverhalf service:", err)
	} else {
		fmt.Println("Deliverhalf service installed and started.")
	}
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
