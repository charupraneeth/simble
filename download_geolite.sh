#!/bin/bash
# Replace YOUR_LICENSE_KEY with an env variable later
wget -O GeoLite2-City.tar.gz "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=$MAXMIND_LICENSE_KEY&suffix=tar.gz"
tar -xvzf GeoLite2-City.tar.gz -C GeoLite2-City --strip-components=1
rm GeoLite2-City.tar.gz