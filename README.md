# Set Xcode Project Build Number

[![Step changelog](https://shields.io/github/v/release/bitrise-io/set-xcode-build-number?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-io/set-xcode-build-number/releases)

Set the value of your iOS app's build number to the specified version number.

<details>
<summary>Description</summary>

Set the value of your iOS app's build number to the specified version number. A great way to keep track of versions 
when submitting bug reports.

If your IPA contains multiple build targets, they need to have the same version number as your app's main target has.
In that case, you need to add this Step to your Workflow for each build target: if you have, say, three targets, you need to have three instances of this Step in your Workflow.
If there are targets with different version numbers the app cannot be submitted for App Review or Beta App Review.

### Configuring the Step 

The step can handle if versions numbers are specified in the project file (default configuration since Xcode 13) and the old style
where the version numbers appear in the **Info.plist** file. It can automatically detect which style is used and act accordingly.

For the simple projects you do not need to do anything because the step uses the previously defined $BITRISE_PROJECT_PATH
and $BITRISE_SCHEME env vars to detect the target settings.

For better customisation the step can be also instructed to look for a specific target and even for specific target configurations
to update the version numbers.

### Useful links 

- [Build numbering and app versioning](https://devcenter.bitrise.io/builds/build-numbering-and-app-versioning/#setting-the-cfbundleversion-and-cfbundleshortversionstring-of-an-ios-app)
- [Current Project Version in Apple documentation](https://developer.apple.com/documentation/xcode/build-settings-reference#Current-Project-Version)
- [Marketing Version in Apple documentation](https://developer.apple.com/documentation/xcode/build-settings-reference#Marketing-Version)
- [CFBundleversion in Apple documentation](https://developer.apple.com/documentation/bundleresources/information_property_list/cfbundleversion)

### Related Steps 

- [Xcode Archive & Export for iOS](https://www.bitrise.io/integrations/steps/xcode-archive)
- [Set Android Manifest Version Code and Name](https://www.bitrise.io/integrations/steps/set-android-manifest-versions)
</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://devcenter.bitrise.io/steps-and-workflows/steps-and-workflows-index/).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `project_path` | Xcode Project (`.xcodeproj`) or Workspace (`.xcworkspace`) path. | required | `$BITRISE_PROJECT_PATH` |
| `scheme` | Xcode Scheme name. | required | `$BITRISE_SCHEME` |
| `target` | Xcode Target name.  It is optional and if specified then the step will find the given target and update the version numbers for it.   If it is left empty then the step will use the scheme's default target to update the version numbers. |  |  |
| `configuration` | Xcode Configuration name.  It is optional and if specified then the step will only update the configuration with the given name.   If it is left empty then the step will update all of the target's configurations with the build and version number. |  |  |
| `build_version` | This will be either the CFBundleVersion in the Info.plist file or the CURRENT_PROJECT_VERSION in the project file. | required | `$BITRISE_BUILD_NUMBER` |
| `build_version_offset` | This offset will be added to `build_version` input's value. It must be a positive number. |  |  |
| `build_short_version_string` | This will be either the CFBundleShortVersionString in the Info.plist file or the MARKETING_VERSION in the project file.  If it is empty then the step will not modify the existing value. |  |  |
</details>

<details>
<summary>Outputs</summary>

| Environment Variable | Description |
| --- | --- |
| `XCODE_BUNDLE_VERSION` | The bundle version used in either in Info.plist or project file |
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-io/set-xcode-build-number/pulls) and [issues](https://github.com/bitrise-io/set-xcode-build-number/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://devcenter.bitrise.io/bitrise-cli/run-your-first-build/).

Learn more about developing steps:

- [Create your own step](https://devcenter.bitrise.io/contributors/create-your-own-step/)
- [Testing your Step](https://devcenter.bitrise.io/contributors/testing-and-versioning-your-steps/)
