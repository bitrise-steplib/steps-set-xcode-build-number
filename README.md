# Set Xcode Project Build Number

Sets the Build Number (bundle version) to the specified value,
in the target `Info.plist` file for the next build.

## Inputs

- plist_path: "" __(required)__
    > Path to the given target's Info.plist file. You need to use this step for each archivable target of your project.
- build_version: "$BITRISE_BUILD_NUMBER" __(required)__
    > Set the CFBundleVersion to this value.
- build_version_offset: ""
    > This offset will be added to `build_version` input's value.
- build_short_version_string: ""
    > Set the CFBundleShortVersionString to this value.

## Outputs

### Exported Environment variables

- XCODE_BUNDLE_VERSION: The bundle version used in the Info.plist file

## Contribute

1. Fork this repository
1. Make changes
1. Submit a PR

## How to run this step from source

1. Clone this repository
1. `cd` to the cloned repository's root
1. Create a bitrise.yml (if not yet created)
1. Prepare a workflow that contains a step with the id: `path::./`
    > For example:
    > ```yaml
    > format_version: "6"
    > default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
    > 
    > workflows:
    >   my-workflow:
    >     steps:
    >     - path::./:
    >         inputs: 
    >         - my_input: "my input value"
    > ```
1. Run the workflow: `bitrise run my-workflow`

## About
This is an official Step managed by Bitrise.io and is available in the [Workflow Editor](https://www.bitrise.io/features/workflow-editor) and in our [Bitrise CLI](https://github.com/bitrise-io/bitrise) tool. If you seen something in this readme that never before please visit some of our knowledge base to read more about that:
  - devcenter.bitrise.io
  - discuss.bitrise.io
  - blog.bitrise.io
