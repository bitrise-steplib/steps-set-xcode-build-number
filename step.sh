#!/bin/bash

# exit if a command fails
set -e

#
# Required parameters
if [ -z "${plist_path}" ] ; then
  echo " [!] Missing required input: plist_path"
  exit 1
fi
if [ ! -f "${plist_path}" ] ; then
  echo " [!] File doesn't exist at specified Info.plist path: ${plist_path}"
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

echo " (i) Provided Info.plist path: ${CONFIG_project_info_plist_path}"
echo " (i) Build number (bundle version): ${CONFIG_new_bundle_version}"


# ---------------------
# --- Main:

# verbose / debug print commands
set -v

# ---- Current Bundle Version:
/usr/libexec/PlistBuddy -c "Print CFBundleVersion" "${CONFIG_project_info_plist_path}"

# ---- Set Bundle Version:
/usr/libexec/PlistBuddy -c "Set :CFBundleVersion ${CONFIG_new_bundle_version}" "${CONFIG_project_info_plist_path}"

# ==> Bundle Version / Build Number changed
