# pm0

Process manager inspired by pm2 (but faster)

## Install script

```Shell
curl -s https://api.github.com/repos/TrixiS/pm0/releases \
| grep "browser_download_url" \
| cut -d : -f 2,3 \
| tr -d \" \
| wget -qi -
```
