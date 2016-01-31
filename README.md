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
*/6 * * * * /usr/bin/python2 ~/github/elasticDNS/edd.py
```

Due to the current TTL of the domain setting, you should use a cron that runs every 6 minutes to make sure that your change has time to propagate before you make a second call to avoid double calls.

Here are the relevant lines you need to configure manually to get the program to work in edd.py:

```python
baseDomain = '<enter base domain here>'
subDomain = '<enter subdomain here>'
record_type = '<enter type here. A or AAAA>'
```

You also need to modify this line in creds.py:

```python
file = open('/<path to user's home directory>/.aws/credentials')
```

As of this particular writing, AWS_R53_ADDR_2 is not being used...in fact, this is likely going to be removed as I was GOING to manually update two resource records, but I opted against it at this time. Eventually, this will be implemented in a loop based on an array which tracks the number of Addresses you want to update. This will all be done at config time.

You do not need to add the period "." at the end of the base domain. Route53 will allow for the domain name without the FQDN "." at the end and will still find the resource record. I'm also going to remove some of these variables as they are redundant, but it was better to make the first push and then clean it up as I go.

I have now converted this over to use boto3 instead of boto2. If you would like to use boto2, you can uncomment the boto2 section and comment out the boto3 section. I will likely split this out into its own python file based on which one you have installed.
