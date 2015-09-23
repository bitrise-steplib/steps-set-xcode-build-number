#!/bin/bash

#
# Don't forget to replace the my_plist_path environment
#  in your Workflow with your project's Info.plist path!
#

THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# exit if a command fails
set -e
# verbose / debug print commands
set -v

#
# Required parameters
if [ -z "${plist_path}" ] ; then
	echo " [!] Missing required input: plist_path"
	exit 1
fi
if [ ! -f "${plist_path}" ] ; then
	echo " [!] Specified Info.plist path is not a file: ${plist_path}"
	exit 1
fi

if [ -z "${build_version}" ] ; then
	echo " [!] You didn't set the version number manually, using build_version"
  exit 1
fi

# ---------------------
# --- Configs:

CONFIG_project_info_plist_path="${plist_path}"

CONFIG_new_bundle_version="${BITRISE_BUILD_NUMBER}"

echo " (i) Info.plist path: ${CONFIG_project_info_plist_path}"
echo " (i) Build number (bundle version): ${CONFIG_new_bundle_version}"


# ---------------------
# --- Main:

echo "---- Current Bundle Version:"
/usr/libexec/PlistBuddy -c "Print CFBundleVersion" "${CONFIG_project_info_plist_path}"

echo "---- Set Bundle Version:"
echo " * to: ${CONFIG_new_bundle_version}"
/usr/libexec/PlistBuddy -c "Set :CFBundleVersion ${CONFIG_new_bundle_version}" "${CONFIG_project_info_plist_path}"

echo "-------------------------"

echo
echo "==> Finished with success"
echo
