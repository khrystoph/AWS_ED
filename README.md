# Brief overview

This project is intended to allow you to leverage your own Domain name to keep track of your public IP address. If you are a
residential customer or you are a small/medium business that does not have a static IP at your office and you want to be able
to access resources remotely, this is something you may want to leverage.

This tool will track your IP address locally and if it detects a change in your IP, it will make a call to route53 using boto
and update your resource record (your "A" name record) for your domain. This currently assumes that you are using a subdomain instead of using the Apex record.

To get this to function on a repeating basis, you are going to need to create a cron job to run as often as you need to update your record. For now, there are a lot of things that are hard-coded into the application that will be changed over time (ttl, hostname(s), read from config file, initial configuration, external or ec2, etc.).

In the future, this will apply to ec2 instances on stop/start so that you can create a private hosted zone (or use a public DNS entry if you like) to track your IP address of an instance if you perform a stop/start. Once that part is implemented, you will not need to use EIPs for instances you wish to access.

# Basic instructions

You need to fill in the sections that have `<>` with the relevant information. Also, you need to make sure you have set up your hosted zone within route53 prior to running this application. Also, this is written in python2, not python 3. So, if you have python 3 installed, you need to make sure you also have python 2 installed on your system and you will need to directly call python2 to run the project. Here is what my cron file looks like:

```cron
*/5 * * * * <username> /usr/bin/python2 ~/github/elasticDNS/edd.py
```

Here are the relevant lines you need to configure manually to get the program to work:

```python
baseDomain = '<Enter base domain here>'
subDomain = '<enter subdomain here>'
AWS_ACCESS_KEY_ID = '<Enter AWS_ACCESS_KEY_ID HERE>'
AWS_SECRET_ACCESS_KEY = '<Enter AWS_SECRET_ACCESS_KEY HERE>'
AWS_R53_ADDR_1 = "<Enter first address to update>"
AWS_R53_ADDR_2 = "<Enter second address to update>"
```

As of this particular writing, AWS_R53_ADDR_2 is not being used...in fact, this is likely going to be removed as I was GOING to manually update two resource records, but I opted against it at this time. Eventually, this will be implemented in a loop based on an array which tracks the number of Addresses you want to update. This will all be done at config time.

You do not need to add the period "." at the end of the base domain. Route53 will allow for the domain name without the FQDN "." at the end and will still find the resource record. I'm also going to remove some of these variables as they are redundant, but it was better to make the first push and then clean it up as I go.
