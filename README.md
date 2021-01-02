![Go](https://github.com/vleedev/gitlabRegistryController/workflows/Go/badge.svg)
# gitlabRegistryController
This tool is used to delete images in gitlab registry
## Install
Download binary file in releases tab.

This is an example shell script to install. (Please change the file name and version)

    #!/bin/sh
    workDir=$(pwd)
    gitlabRegistryControllerFile="${workDir}/gitlabRegistryController"
    if [ ! -s "${gitlabRegistryControllerFile}" ]; then
        curl --request GET -sL \
              --url 'https://github.com/vleedev/gitlabRegistryController/releases/download/0.4.9/gitlabRegistryController-linux-amd64'\
              --output "${gitlabRegistryControllerFile}"
    fi
    chmod +x "${gitlabRegistryControllerFile}"

## Get help
`./gitlabRegistryController -h`
## Usage of ./gitlabRegistryController:
- `-domain string`
    - a base url of your gitlab with api version, ex: https://gitlab.example.com/api/v4
    - If not set, take `CI_API_V4_URL` env automatically.
- `-authToken string`
    - a token that is used to auth with gitlab
    - If not set, take `AUTH_TOKEN` env automatically.
- `-nameSpace string`
    - a namespace of your project
    - If not set, take `CI_PROJECT_NAMESPACE` env automatically.
- `-projectName string`
    - a project name of your project
    - If not set, take `CI_PROJECT_NAME` env automatically.
- `-specificTag string`
    - a image tag that you want to delete
- `-regex string`
    - a regex pattern to match all images
- `-hold int`
    - a volume of images that you want to keep from latest to older (default 3)
    
### Delete an specific image:
    ./gitlabRegistryController -specificTag nginx-latest
### Delete all images that match with an regex pattern
- Default: keep 3 images from latest to older


    ./gitlabRegistryController -regex=".*-nginx"
    
- Keep n images from latest to older


    ./gitlabRegistryController -regex=".*-nginx" -hold=n