# Installation du reverse proxy nginx

```sh
* sudo apt install nginx
* sudo apt install certbot python3-certbot-nginx
* sudo certbot --nginx -d flow.linkafric.com
* sudo certbot renew --dry-run

```
# Configuration de Nginx

* sudo touch /etc/nginx/sites-available/flow.linkafric.com
* sudo cp nginx.conf /etc/nginx/sites-available/flow.linkafric.com
* sudo rm /etc/nginx/sites-available/default
* sudo rm /etc/nginx/sites-enabled/default
* sudo ln -s /etc/nginx/sites-available/flow.linkafric.com /etc/nginx/sites-enabled/
* sudo nginx -t
* sudo systemctl restart nginx

# Verification de l'API

* curl https://flow.linkafric.com/api/v1/flows

