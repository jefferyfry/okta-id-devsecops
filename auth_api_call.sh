
#! /usr/bin/env nix-shell
#! nix-shell -i bash -p curl -p jq
# auth_api_call.sh
# A shell script which demonstrates how to get an access_token from from Okta using the OAuth 2.0
# Modified from https://github.com/jpf/okta-get-id-token
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

curl="curl"
jq="jq"

org_server=""
auth_server=""
client_id=""
redirect=""
username=""
password=""
verbose=0

while getopts ":o:s:c:r:u:p:v" OPTION
do
    case $OPTION in
    o)
        org_server="$OPTARG"
    ;;
    s)
        auth_server="$OPTARG"
    ;;
    c)
        client_id="$OPTARG"
    ;;
    r)
        redirect="$OPTARG"
    ;;
    u)
        username="$OPTARG"
    ;;
    p)
        password="$OPTARG"
    ;;
    v)
        verbose=1
    ;;
    [?])
        echo "Usage: $0 -o org_server -s auth_server -c client_id -r redirect -u username -p password apiurl" >&2
        echo ""
        echo "Example:"
        echo "$0 -o 'https://dev-73225252.okta.com' -s 'https://dev-73225252.okta.com/oauth2/aus1efvp3jwospP0Y5d7' -c aBCdEf0GhiJkLMno1pq2 -u AzureDiamond -p hunter2 -r 'https://example.net/your_application' http://localhost/api/v1/users"
        exit 1
    ;;
    esac
done

shift $(($OPTIND - 1))
apiurl=$1

redirect_uri=$(curl --silent --output /dev/null --write-out %{url_effective} --get --data-urlencode "$redirect" "" | cut -d '?' -f 2)
if [ $verbose -eq 1 ]; then
    echo "Redirect URI: '${redirect_uri}'"
fi

# authenticate to the org server using username and password. gets a temporary session token.
rv=$(curl --silent "${org_server}/api/v1/authn" \
          -H "Origin: ${redirect}" \
          -H 'Content-Type: application/json' \
          -H 'Accept: application/json' \
          --data-binary $(printf '{"username":"%s","password":"%s"}' $username $password) )
session_token=$(echo $rv | jq -r .sessionToken )

if [ $verbose -eq 1 ]; then
    echo "Authn curl: '${rv}'"
fi
if [ $verbose -eq 1 ]; then
    echo "Session token: '${session_token}'"
fi

if [ -z "$session_token" ]; then exit 1; fi

# now call the authentication server with the session token and get an access token
authorize_url=$(printf "${auth_server}/v1/authorize?sessionToken=%s&client_id=%s&scope=openid&response_type=token&response_mode=fragment&nonce=%s&redirect_uri=%s&state=%s" \
      $session_token \
      $client_id \
      "staticNonce" \
      $redirect_uri \
      "staticState")
if [ $verbose -eq 1 ]; then
    echo "Authorize URL: '${authorize_url}'"
fi

rv=$(curl --silent -I $authorize_url 2>&1)
if [ $verbose -eq 1 ]; then
    echo "Here is the return value: "
    echo $rv
fi
if [ $verbose -eq 1 ]; then
    echo "Here is the location url: "
    echo "$rv" | grep 'location'
fi

access_token=$(echo "$rv" | grep 'location' |  awk -F'[=&]' '{print $2}')
if [ $verbose -eq 1 ]; then
    echo "Here is the access token: "
    echo $access_token
fi

if [ -z "$access_token" ]; then exit 1; fi

if [ $verbose -eq 1 ]; then
    echo "API URL: "
    echo $apiurl
fi

# now call our API using the access token
status_code=$(curl -v -H "Authorization: Bearer ${access_token}" -s -o /dev/null -w "%{http_code}" $apiurl)

if [ $verbose -eq 1 ]; then
    echo "API Response: "
    echo $status_code
fi

if [ $status_code == 200 ]; then
  exit 0
else
  exit 1
fi
