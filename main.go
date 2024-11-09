package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sahl/structs"
	"sahl/utils"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

const (
	rgbGray = "\033[38;2;150;150;150m" // Approximate #bdbdbd in 256-color mode
	reset   = "\033[0m"                 // Reset color
)

var verbose bool

func main() {
	green := color.New(color.FgGreen).SprintFunc()

	// making background white and text black
	black := color.New(color.FgBlack)
	boldBlack := black.Add(color.Bold)
	whiteBackgroundBlackText := boldBlack.Add(color.BgWhite)

	// Define flags using pflag, adding short flags
	installFlag := pflag.BoolP("install", "i", false, "Install packages (use space-separated names for multiple packages)")
	helpFlag := pflag.BoolP("help", "h", false, "Display help")
	listFlag := pflag.BoolP("list", "l", false, "List installed packages")
	uninstallFlag := pflag.StringP("uninstall", "r", "", "Uninstall a package")
	verboseFlag := pflag.BoolP("verbose", "v", false, "Enable verbose output")
	checkFlag := pflag.StringP("check", "d", "", "Check if a package is installed")
	forceFlag := pflag.BoolP("force", "f", false, "Force reinstallation of the package")
	forceDependenciesFlag := pflag.BoolP("force-deps", "F", false, "Force reinstallation of the package and its dependencies")

	pflag.Parse()
	verbose = *verboseFlag

	if *helpFlag {
		whiteBackgroundBlackText.Println(" Usage of sahl: ")
		pflag.PrintDefaults()
		return
	}

	if *installFlag {
        packages := pflag.Args() // Get all arguments after flags
        if len(packages) == 0 {
            fmt.Println("No packages specified for installation.")
            return
        }

        forceReinstall := *forceFlag
        forceReinstallWithDeps := *forceDependenciesFlag
        visited := make(map[string]bool)

        for _, packageName := range packages {
            fmt.Printf("\nInstalling package: %s\n", packageName)
            installed, err := installPackageWithDependencies(packageName, visited, forceReinstall, forceReinstallWithDeps)
            if err != nil {
                fmt.Printf("Failed to install package %s: %v\n", packageName, err)
            } else if installed {
                fmt.Printf("%s[✓] Package %s was successfully installed%s\n", green(""), packageName, reset)
            } else {
                fmt.Printf("Package %s was already installed\n", packageName)
            }
        }
        return
    }

	if *listFlag {
		whiteBackgroundBlackText.Printf("Installed packages:%s", reset)
		fmt.Println()
		utils.PrintInstalledPackages()
		return
	}

	if *checkFlag != "" {
		packageName := *checkFlag
		if utils.CheckIfInstalled(packageName) {
			color.Green("[✓] Package %s: Installed\n", packageName)
		} else {
			fmt.Printf("%sPackage %s: Not installed%s\n", rgbGray, packageName, reset)
		}
		return
	}

	if *uninstallFlag != "" {
		packageName := *uninstallFlag
		fmt.Printf("Uninstalling package: %s\n", packageName)
		err := utils.UninstallPackage(packageName)
		if err != nil {
			color.Red("[✗] Error uninstalling package: %v", err)
		}
		return
	}
}

func installPackageWithDependencies(packageName string, visited map[string]bool, forceReinstall bool, forceReinstallWithDeps bool) (bool, error) {
	if visited[packageName] {
		return false, nil
	}

	yamlURL := fmt.Sprintf("https://raw.githubusercontent.com/abdurahman-harouat/fennec-hub/main/source_files/%s/OOO.yaml", packageName)
	yamlData, err := utils.DownloadFile(yamlURL)
	if err != nil {
		return false, fmt.Errorf("error downloading YAML file for %s: %v", packageName, err)
	}

	var pkgDef structs.PackageDefinition
	err = yaml.Unmarshal(yamlData, &pkgDef)
	if err != nil {
		return false, fmt.Errorf("error parsing package definition for %s: %v", packageName, err)
	}

	if utils.CheckIfInstalled(packageName) && !(forceReinstall || forceReinstallWithDeps) {
		fmt.Printf("%s>>> Package %s is already installed --SKIPPING installation.%s\n", rgbGray, packageName, reset)
		return false, nil
	}

	visited[packageName] = true

	for _, dep := range pkgDef.Dependencies {
		_, err := installPackageWithDependencies(dep, visited, forceReinstallWithDeps, forceReinstallWithDeps)
		if err != nil {
			return false, fmt.Errorf("failed to install dependency %s: %v", dep, err)
		}
	}

	fmt.Printf("%s[*] Installing package: %s...%s\n", rgbGray, packageName, reset)
	err = installPackage(packageName)
	if err != nil {
		return false, fmt.Errorf("failed to install package %s: %v", packageName, err)
	}

	err = updatePackageDatabase(packageName, pkgDef.PkgVersion)
	if err != nil {
		return false, fmt.Errorf("failed to update package database for %s: %v", packageName, err)
	}

	return true, nil
}

func installPackage(packageName string) error {
	green := color.New(color.FgGreen).SprintFunc()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("[✗] Error getting user's home directory: %v", err)
	}

	cacheDir := filepath.Join(homeDir, ".cache", "sahl")
	err = os.MkdirAll(cacheDir, 0755)
	if err != nil {
		return fmt.Errorf("[✗] Error creating cache directory: %v", err)
	}

	yamlURL := fmt.Sprintf("https://raw.githubusercontent.com/abdurahman-harouat/fennec-hub/main/source_files/%s/OOO.yaml", packageName)
	yamlData, err := utils.DownloadFile(yamlURL)
	if err != nil {
		return fmt.Errorf("[✗] Error downloading YAML file: %v", err)
	}

	if verbose {
		fmt.Printf("%s[✓] Package definition (OOO.yaml) downloaded successfully.%s\n", green(""), reset)
	}

	var pkgDef structs.PackageDefinition
	err = yaml.Unmarshal(yamlData, &pkgDef)
	if err != nil {
		return fmt.Errorf("[✗] Error parsing package definition: %v", err)
	}

	cachedFilePath, err := utils.GetOrDownloadPackage(pkgDef.Source.URL, cacheDir, pkgDef.Source.MD5)
	if err != nil {
		return fmt.Errorf("[✗] Error getting package file: %v", err)
	}

	// Download additional packages
	for _, additional := range pkgDef.AdditionalDownloads {
		fmt.Printf("[*] Downloading additional package: %s\n", additional.URL)
		_, err := utils.GetOrDownloadPackage(additional.URL, cacheDir, additional.MD5)
		if err != nil {
			return fmt.Errorf("[✗] Error getting additional package file: %v", err)
		}
	}

	err = utils.UntarFile(cachedFilePath, cacheDir)
	if err != nil {
		return fmt.Errorf("[✗] Error extracting package: -> cachedFilePath: %s %v", cachedFilePath, err)
	}

	var extractedDirPath string
	if pkgDef.ExtractedDir != "" {
		extractedDirPath = filepath.Join(cacheDir, pkgDef.ExtractedDir)
	} else {
		packageDirName := utils.RemoveExtension(filepath.Base(cachedFilePath))
		extractedDirPath = filepath.Join(cacheDir, packageDirName)
	}

	err = os.Chdir(extractedDirPath)
	if err != nil {
		return fmt.Errorf("[✗] Error changing directory to %s: %v", extractedDirPath, err)
	}

	if verbose {
		fmt.Printf("\n%s[✓] Changed directory to: %s%s\n", green(""), extractedDirPath, reset)
	}

	entries, err := os.ReadDir(extractedDirPath)
	if err != nil {
		return fmt.Errorf("[✗] Error reading directory contents: %v", err)
	}

	currentDirName := filepath.Base(extractedDirPath)

	for _, entry := range entries {
		if entry.IsDir() && entry.Name() == currentDirName {
			newPath := filepath.Join(extractedDirPath, entry.Name())
			err = os.Chdir(newPath)
			if err != nil {
				return fmt.Errorf("[✗] Error changing directory to %s: %v", newPath, err)
			}
			if verbose {
				fmt.Printf("%s[✓] Changed directory to subdirectory: %s%s\n", green(""), newPath, reset)
			}
			break
		}
	}

	// Run build commands
	for _, command := range pkgDef.Build {
		if verbose {
			fmt.Printf("%s[*] Executing build command: %s%s\n", green(""), command, reset)
		}
		err := utils.RunCommand(command)
		if err != nil {
			return fmt.Errorf("[✗] Error executing build command '%s': %v", command, err)
		}
	}

	if verbose {
		fmt.Printf("%s[✓] Build process completed successfully%s\n", green(""), reset)
	}

	// Run install commands
	for _, command := range pkgDef.Install {
		if verbose {
			fmt.Printf("%s[*] Executing install command: %s%s\n", green(""), command, reset)
		}
		err := utils.RunCommandWithSudo(command)
		if err != nil {
			return fmt.Errorf("[✗] Error executing install command '%s': %v", command, err)
		}
	}

	if verbose {
		fmt.Printf("%s[✓] All commands installed successfully%s\n", green(""), reset)
	}

	// Run additional commands
	for _, command := range pkgDef.AdditionalCommands {
		if verbose {
			fmt.Printf("[*] Running additional command: %s\n", command)
		}
		err := utils.RunCommand(command)
		if err != nil {
			return fmt.Errorf("[✗] Error executing additional command '%s': %v", command, err)
		}
	}

	if verbose {
		fmt.Printf("%s[✓] Additional commands completed successfully%s\n", green(""), reset)
	}

	// Run additional commands with sudo
	for _, command := range pkgDef.AdditionalCommandsWithSudo {
		if verbose {
			fmt.Printf("[*] Running additional command with sudo: %s\n", command)
		}
		err := utils.RunCommandWithSudo(command)
		if err != nil {
			return fmt.Errorf("[✗] Error executing additional command '%s': %v", command, err)
		}
	}

	if verbose {
		fmt.Printf("%s[✓] Additional sudo commands completed successfully%s\n", green(""), reset)
	}

	return nil
}

func updatePackageDatabase(packageName string, packageVersion string) error {
	logFile := "/var/log/packages.log"

	// Remove existing entry for the package
	removeCmd := exec.Command("sudo", "sed", "-i", fmt.Sprintf("/^%s /d", packageName), logFile)
	removeOutput, err := removeCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error removing previous installation record for %s: %v\nOutput: %s", packageName, err, string(removeOutput))
	}

	// Add new entry for the package
	logEntry := fmt.Sprintf("%s %s %s", packageName, packageVersion, time.Now().Format(time.RFC3339))
	logCmd := exec.Command("sudo", "sh", "-c", fmt.Sprintf("echo '%s' >> %s", logEntry, logFile))
	logOutput, err := logCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error logging installation details for %s: %v\nOutput: %s", packageName, err, string(logOutput))
	}

	if verbose {
		fmt.Printf("[✓] Installation details logged for %s\n", packageName)
	}

	return nil
}