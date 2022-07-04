package main

import (
	"os"
	"path/filepath"

	"github.com/bitrise-io/go-xcode/xcodeproject/xcodeproj"
)

func main() {
	bundleVersion := 11
	shortVersion := "9.0"
	projectPath := "/Users/szabi/Dev/misc/TestAPP2/TestAPP2.xcodeproj"
	//projectPath := "/Users/szabi/Dev/bitrise-io/sample-apps-ios-simple-objc/ios-simple-objc/ios-simple-objc.xcodeproj"

	project, err := xcodeproj.Open(projectPath)
	if err != nil {
		os.Exit(1)
	}

	if generatesInfoPlist(project) {
		updateVersionNumbersInProject(project, bundleVersion, shortVersion)
	} else {
		updateVersionNumbersInInfoPlist(project, projectPath, bundleVersion, shortVersion)
	}

	os.Exit(0)
}

func generatesInfoPlist(project xcodeproj.XcodeProj) bool {
	buildSettings, _ := project.TargetBuildSettings("TestAPP2", "Debug")
	generatesInfoPlist, _ := buildSettings.String("GENERATE_INFOPLIST_FILE")

	return generatesInfoPlist == "YES"
}

func updateVersionNumbersInProject(project xcodeproj.XcodeProj, bundleVersion int, shortVersion string) {
	err := project.UpdateBuildSetting("TestAPP2", "", "CURRENT_PROJECT_VERSION", bundleVersion)
	if err != nil {
		os.Exit(2)
	}

	err = project.UpdateBuildSetting("TestAPP2", "", "MARKETING_VERSION", shortVersion)
	if err != nil {
		os.Exit(3)
	}
}

func updateVersionNumbersInInfoPlist(project xcodeproj.XcodeProj, projectPath string, bundleVersion int, shortVersion string) {
	buildSettings, _ := project.TargetBuildSettings("ios-simple-objc", "Debug")
	infoPlistPath, _ := buildSettings.String("INFOPLIST_FILE")

	absoluteInfoPlistPath := filepath.Join(filepath.Dir(projectPath), infoPlistPath)

	infoPlist, format, _ := xcodeproj.ReadPlistFile(absoluteInfoPlistPath)

	infoPlist["CFBundleVersion"] = bundleVersion
	infoPlist["CFBundleShortVersionString"] = shortVersion

	_ = xcodeproj.WritePlistFile(absoluteInfoPlistPath, infoPlist, format)
}
