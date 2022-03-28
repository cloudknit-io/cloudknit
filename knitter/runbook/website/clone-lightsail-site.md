# Clone Lightsail compuzest website

## Overview

Steps to clone the compuzest website so that we can make changes to it

## Initial Steps Overview

* Go to [Snapshots page](https://lightsail.aws.amazon.com/ls/webapp/us-east-1/instances/compuzest-website/snapshot) & create new instance using latest snapshot
* SSH into the new instance
* Run following commands:

```bash
vim ~/apps/wordpress/htdocs/wp-config.php
```

Update `DOMAIN_CURRENT_SITE` to `qa.compuzest.com`

```bash
sudo ~/apps/wordpress/bnconfig --machine_hostname qa.compuzest.com
```

* Go to qa.compuzest.com and update the hoe page to say "QA Site"
