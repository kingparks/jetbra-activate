# shellcheck disable=SC2019
# shellcheck disable=SC2018
software=$(echo "$1"|tr A-Z a-z)
softvm="${software}.vmoptions"
currCrackPath=$(
  cd $(dirname $0)
  pwd
)
targetFilePath="/Users/${USER}/Library/Application Support"
cpDir="${targetFilePath}/JetBrains"
echo $cpDir
if [ ! -d "${cpDir}" ]; then
  mkdir -p "${cpDir}"
fi
jarFile="${currCrackPath}/active-agt.jar"
plugins="${currCrackPath}/plugins"
config="${currCrackPath}/config"
if [ -f "${jarFile}" ]; then
  $(cp "${jarFile}" "${cpDir}")
  if [ ! -d "${cpDir}/plugins" ]; then
    $(mkdir "${cpDir}/plugins")
  fi
  $(cp -rf ${plugins}/* "${cpDir}/plugins")
  if [ ! -d "${cpDir}/config" ]; then
    $(mkdir "${cpDir}/config")
  fi
  $(cp -rf ${config}/* "${cpDir}/config")
else
  echo "active-agt.jar is missing, ${software} crack failed!"
  exit
fi
softwareInstall="false"
for file in $(ls -a "$cpDir"); do
  if [ -d "$cpDir/$file" ]; then
    result=$(echo $file | grep -i $software)
    if [ ${result}x != ""x ]; then
      softwareInstall="true"
      echo "Success! Activate ${software} to 2099"
      keyInstall="${currCrackPath}/key/${software}.key"
      cpxDir="$cpDir/$file"
      if [ -f "${keyInstall}" ]; then
        $(cp "${keyInstall}" "${cpxDir}")
      fi
      versionInstall="$cpDir/$file/$softvm"
      if [ -f "$versionInstall" ]; then
        $(echo "" >"$versionInstall")
      else
        $(touch "$versionInstall")
      fi
      $(echo "-javaagent:${cpDir}/active-agt.jar" >"${versionInstall}")
      $(echo "--add-opens=java.base/jdk.internal.org.objectweb.asm=ALL-UNNAMED" >>"${versionInstall}")
      $(echo "--add-opens=java.base/jdk.internal.org.objectweb.asm.tree=ALL-UNNAMED" >>"${versionInstall}")
    fi
  fi
done
if [ "$softwareInstall" = "false" ]; then
  echo "Please start ${software} first!"
fi
