package step

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/bitrise-io/go-steputils/v2/export"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-xcode/v2/autocodesign/projectmanager"
	"github.com/bitrise-io/go-xcode/xcodeproject/serialized"
	"github.com/bitrise-io/go-xcode/xcodeproject/xcodeproj"
)

const (
	infoPlistFileKey = "INFOPLIST_FILE"
	envVarRegex      = `^.*\$\(.+\).*$`
)

type Updater struct {
	inputParser stepconf.InputParser
	exporter    export.Exporter
	logger      log.Logger
}

func NewUpdater(inputParser stepconf.InputParser, exporter export.Exporter, logger log.Logger) Updater {
	return Updater{
		inputParser: inputParser,
		exporter:    exporter,
		logger:      logger,
	}
}

func (u Updater) ProcessConfig() (Config, error) {
	var input Input
	err := u.inputParser.Parse(&input)
	if err != nil {
		return Config{}, err
	}

	u.logger.EnableDebugLog(input.Verbose)

	stepconf.Print(input)
	u.logger.Println()

	return Config{
		ProjectPath:             input.ProjectPath,
		Scheme:                  input.Scheme,
		Target:                  input.Target,
		Configuration:           input.Configuration,
		BuildVersion:            input.BuildVersion,
		BuildVersionOffset:      input.BuildVersionOffset,
		BuildShortVersionString: input.BuildShortVersionString,
	}, nil
}

func (u Updater) Run(config Config) (Result, error) {
	helper, err := projectmanager.NewProjectHelper(config.ProjectPath, config.Scheme, config.Configuration)
	if err != nil {
		return Result{}, err
	}

	generated, err := generatesInfoPlist(helper, config.Target, config.Configuration)
	if err != nil {
		return Result{}, err
	}

	config.BuildVersion, err = incrementBuildVersion(u.logger, config.BuildVersion, config.BuildVersionOffset)
	if err != nil {
		return Result{}, err
	}

	if generated {
		u.logger.Printf("The version numbers are stored in the project file.")

		err := u.updateVersionNumbersInProject(helper, config.Target, config.Configuration, config.BuildVersion, config.BuildShortVersionString)
		if err != nil {
			return Result{}, err
		}
	} else {
		u.logger.Printf("The version numbers are stored in the plist file.")

		err := u.updateVersionNumbersInInfoPlist(helper, config.Scheme, config.Target, config.Configuration, config.BuildVersion, config.BuildShortVersionString)
		if err != nil {
			return Result{}, err
		}
	}

	u.logger.Donef("Version numbers successfully updated.")

	return Result{BuildVersion: config.BuildVersion}, nil
}

func (u Updater) Export(result Result) error {
	return u.exporter.ExportOutput("XCODE_BUNDLE_VERSION", result.BuildVersion)
}

func generatesInfoPlist(helper *projectmanager.ProjectHelper, targetName, configuration string) (bool, error) {
	buildConfig, err := buildConfiguration(helper, targetName, configuration)
	if err != nil {
		return false, err
	}

	generatesInfoPlist := buildConfig.BuildSettings["GENERATE_INFOPLIST_FILE"] == "YES"

	return generatesInfoPlist, err
}

func incrementBuildVersion(logger log.Logger, buildVersion string, offset int64) (string, error) {
	// Check if build version is numeric
	parsedBuildVersion, err := strconv.ParseInt(buildVersion, 10, 64)
	if err != nil {
		logger.Infof("Provided build version is not numeric (%s), using it as-is without incrementing", buildVersion)
		if offset > 0 {
			return "", fmt.Errorf("build version offset (%d) cannot be applied to non-numeric build version (%s), use 0 as the offset to use the build version as-is", offset, buildVersion)
		}
		return buildVersion, nil
	}

	// Numeric build version provided, increment it
	if offset >= 0 {
		return strconv.FormatInt(parsedBuildVersion+offset, 10), nil
	}
	logger.Infof("Build version offset is negative (%d), skipping version increment.", offset)
	return buildVersion, nil
}

func (u Updater) updateVersionNumbersInProject(helper *projectmanager.ProjectHelper, targetName, configuration string, bundleVersion, shortVersion string) error {
	if targetName == "" {
		targetName = helper.MainTarget.Name
	}

	for _, target := range helper.XcProj.Proj.Targets {
		if target.Name != targetName {
			continue
		}

		for _, buildConfig := range target.BuildConfigurationList.BuildConfigurations {
			if configuration != "" && buildConfig.Name != configuration {
				continue
			}

			u.logger.Printf("Updating build settings for the %s target", target.Name)

			oldProjectVersion := buildConfig.BuildSettings["CURRENT_PROJECT_VERSION"]
			buildConfig.BuildSettings["CURRENT_PROJECT_VERSION"] = bundleVersion

			u.logger.Debugf("CURRENT_PROJECT_VERSION %s -> %s", oldProjectVersion, bundleVersion)

			if shortVersion != "" {
				oldMarketingVersion := buildConfig.BuildSettings["MARKETING_VERSION"]
				buildConfig.BuildSettings["MARKETING_VERSION"] = shortVersion

				u.logger.Debugf("MARKETING_VERSION %s -> %s", oldMarketingVersion, shortVersion)
			}
		}
	}

	err := helper.XcProj.Save()
	if err != nil {
		return err
	}

	return nil
}

func (u Updater) updateVersionNumbersInInfoPlist(helper *projectmanager.ProjectHelper, schemeName, targetName, configuration string, bundleVersion, shortVersion string) error {
	buildConfig, err := buildConfiguration(helper, targetName, configuration)
	if err != nil {
		return err
	}

	infoPlistPath, err := buildConfig.BuildSettings.String(infoPlistFileKey)
	// If the path is extracted into a xcconfig file, then it will not appear here in the build settings.
	// We need to use xcodebuild to resolve the path.
	if err != nil {
		if !serialized.IsKeyNotFoundError(err) {
			return err
		}

		u.logger.Printf("Info.plist path was not found in the project\n")
		u.logger.Printf("Using xcodebuild to resolve it\n")

		infoPlistPath, err = extractInfoPlistPathWithXcodebuild(helper.XcProj.Path, schemeName, targetName, configuration)
		if err != nil {
			return err
		}
	}

	// By default, the setting for the Info.plist file path is a relative path from the project file. Of course,
	// developers can override this with something more custom to their setup. They can also use Xcode env vars as part
	// of their path.
	//
	// An example from a SWAT ticket: `$(SRCROOT)/path/to/Info.plist`.
	//
	// The problem with this is that it is not a real path until the env var is resolved. And in this case, Xcode
	// will define this env var, so we only know its value during an xcodebuild execution. So if we see an env var in
	// the path, then we have to list the build settings with `xcodebuild -showBuildSettings` to get a valid path value.
	if hasEnvVars(infoPlistPath) {
		u.logger.Printf("Info.plist path contains env var: %s\n", infoPlistPath)
		u.logger.Printf("Using xcodebuild to resolve it\n")

		infoPlistPath, err = extractInfoPlistPathWithXcodebuild(helper.XcProj.Path, schemeName, targetName, configuration)
		if err != nil {
			return err
		}
	}

	if pathutil.IsRelativePath(infoPlistPath) {
		infoPlistPath = filepath.Join(filepath.Dir(helper.XcProj.Path), infoPlistPath)
	}

	u.logger.Printf("Updating Info.plist at %s", infoPlistPath)

	infoPlist, format, err := xcodeproj.ReadPlistFile(infoPlistPath)
	if err != nil {
		return err
	}

	oldVersion := infoPlist["CFBundleVersion"]
	infoPlist["CFBundleVersion"] = bundleVersion

	u.logger.Debugf("CFBundleVersion %s -> %s", oldVersion, bundleVersion)

	if shortVersion != "" {
		oldVersionString := infoPlist["CFBundleShortVersionString"]
		infoPlist["CFBundleShortVersionString"] = shortVersion

		u.logger.Debugf("CFBundleShortVersionString %s -> %s", oldVersionString, shortVersion)
	}

	err = xcodeproj.WritePlistFile(infoPlistPath, infoPlist, format)
	if err != nil {
		return err
	}

	return nil
}

func hasEnvVars(path string) bool {
	re := regexp.MustCompile(envVarRegex)
	containsEnvVar := re.Match([]byte(path))

	return containsEnvVar
}

func extractInfoPlistPathWithXcodebuild(projectPath, scheme, target, configuration string) (string, error) {
	args := []string{"-project", projectPath}

	if target != "" {
		args = append(args, "-target", target)
	} else if scheme != "" {
		args = append(args, "-scheme", scheme)
	}

	if configuration != "" {
		args = append(args, "-configuration", configuration)
	}

	args = append(args, "-showBuildSettings")

	cmd := command.NewFactory(env.NewRepository()).Create("xcodebuild", args, nil)
	output, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return "", err
	}

	path := infoPlistPathFromOutput(output)
	if path == "" {
		return "", fmt.Errorf("missing Info.plist file path")
	}

	return path, nil
}

func infoPlistPathFromOutput(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		split := strings.Split(line, " = ")

		if len(split) < 2 {
			continue
		}

		if strings.TrimSpace(split[0]) != infoPlistFileKey {
			continue
		}

		return strings.TrimSpace(split[1])
	}

	return ""
}

func buildConfiguration(helper *projectmanager.ProjectHelper, targetName, configuration string) (*xcodeproj.BuildConfiguration, error) {
	if targetName == "" {
		targetName = helper.MainTarget.Name
	}

	if configuration == "" {
		configuration = helper.MainTarget.BuildConfigurationList.DefaultConfigurationName
	}

	var xcodeprojTarget *xcodeproj.Target
	for _, target := range helper.XcProj.Proj.Targets {
		if target.Name == targetName {
			xcodeprojTarget = &target
			break
		}
	}

	if xcodeprojTarget == nil {
		return nil, fmt.Errorf("target '%s' not found in project: %s", targetName, helper.XcProj.Path)
	}

	var xcodeprojBuildConfiguration *xcodeproj.BuildConfiguration
	for _, buildConfig := range xcodeprojTarget.BuildConfigurationList.BuildConfigurations {
		if buildConfig.Name == configuration {
			xcodeprojBuildConfiguration = &buildConfig
			break
		}
	}

	if xcodeprojBuildConfiguration == nil {
		return nil, fmt.Errorf("build configuration '%s' not found for target '%s' in project: %s", configuration, targetName, helper.XcProj.Path)
	}

	return xcodeprojBuildConfiguration, nil
}
