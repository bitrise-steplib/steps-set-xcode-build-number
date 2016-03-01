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
  echo " [!] No build_version specified!"
  exit 1
fi

# ---------------------
# --- Configs:

CONFIG_project_info_plist_path="${plist_path}"
CONFIG_new_bundle_version="${build_version}"
CONFIG_new_short_bundle_short_version_string="${build_short_version_string}"

echo " (i) Provided Info.plist path: ${CONFIG_project_info_plist_path}"
echo " (i) Build number (bundle version): ${CONFIG_new_bundle_version}"

if [ -n "${CONFIG_new_short_bundle_short_version_string}" ] ; then
  echo " (i) Build number (bundle short version string): ${CONFIG_new_short_bundle_short_version_string}"
else
  echo " (i) Build number (bundle short version string): Not specified"
fi

# ---------------------
# --- Main:

# verbose / debug print commands
set -v

# ---- Current Bundle Version:
/usr/libexec/PlistBuddy -c "Print CFBundleVersion" "${CONFIG_project_info_plist_path}"

# ---- Set Bundle Version:
/usr/libexec/PlistBuddy -c "Set :CFBundleVersion ${CONFIG_new_bundle_version}" "${CONFIG_project_info_plist_path}"

if [ -n "${CONFIG_new_short_bundle_short_version_string}" ] ; then
  # ---- Current Bundle Short Version String:
  /usr/libexec/PlistBuddy -c "Print CFBundleShortVersionString" "${CONFIG_project_info_plist_path}"

  # ---- Set Bundle Short Version String:
  /usr/libexec/PlistBuddy -c "Set :CFBundleShortVersionString ${CONFIG_new_short_bundle_short_version_string}" "${CONFIG_project_info_plist_path}"
fi
# ==> Bundle Version and Short Version changed
