# Set Xcode Project Build Number

[![Step changelog](https://shields.io/github/v/release/bitrise-io/set-xcode-build-number?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-io/set-xcode-build-number/releases)

Set the value of your iOS app's bundle version in the `Info.plist` file to the specified version number.

<details>
<summary>Description</summary>

Set the value of your iOS app's bundle version in the `Info.plist` file to the specified version number. A great
way to keep track of versions when submitting bug reports.

If your IPA contains multiple build targets, they need to have the same version number as your app's main target has.
In that case, you need to add this Step to your Workflow for each build target: if you have, say, three targets, you need to have three instances of this Step in your Workflow.
If there are targets with different version numbers the app cannot be submitted for App Review or Beta App Review.

### Configuring the Step 

1. In your Xcode project, set the Generate Info.plist File to No, under PROJECT and TARGETS on the Build Settings tab.
1. Manually create the `Info.plist` file and check it into source control. Make sure you have all the necessary keys defined in the file.
1. Configure this step by pointing the **Info.plist file path** input to the `Info.plist` file in the source repo.
1. Add a value in the Build Number input. 
   This sets the CFBundleVersion key to the specified value in the `Info.plist` file. The default value is the `$BITRISE_BUILD_NUMBER` Environment Variable.
1. Optionally, add a value in the Version Number input. This will set the `CFBundleShortVersionString` key to the specified value in the `Info.plist` file. This input is not required.

### Useful links 

- [Build numbering and app versioning](https://devcenter.bitrise.io/builds/build-numbering-and-app-versioning/#setting-the-cfbundleversion-and-cfbundleshortversionstring-of-an-ios-app)
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
| `plist_path` | Path to the given target's Info.plist file. You need to use this Step for each archivable target of your project.  **NOTE:**<br/> If your IPA contains multiple build targets, they would need to have the same version number as your app's main target has.<br/> You need to add this Step to your Workflow for each build target: if you have, say, three targets, you need to have three instances of this Step in your Workflow. If there are targets with different version numbers the app cannot be submitted for App Review or Beta App Review.  | required |  |
| `build_version` | Set the CFBundleVersion to this value. You can find this in Xcode: - Select your project in the **Project navigator** - Go to the **General** tab and then the **Identity** section - **Build field**  | required | `$BITRISE_BUILD_NUMBER` |
| `build_version_offset` | This offset will be added to `build_version` input's value.  |  |  |
| `build_short_version_string` | Set the CFBundleShortVersionString to this value. You can find this in Xcode: - Select your project in the **Project navigator** - Go to the **General** tab and then the **Identity** section - **Version field**  |  |  |
</details>

<details>
<summary>Outputs</summary>

| Environment Variable | Description |
| --- | --- |
| `XCODE_BUNDLE_VERSION` | The bundle version used in the Info.plist file |
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-io/set-xcode-build-number/pulls) and [issues](https://github.com/bitrise-io/set-xcode-build-number/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://devcenter.bitrise.io/bitrise-cli/run-your-first-build/).

Learn more about developing steps:

- [Create your own step](https://devcenter.bitrise.io/contributors/create-your-own-step/)
- [Testing your Step](https://devcenter.bitrise.io/contributors/testing-and-versioning-your-steps/)
