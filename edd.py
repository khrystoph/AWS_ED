import pycurl
from StringIO import StringIO
import boto.route53
import socket
from os import path
from os import getenv
from boto.route53.record import ResourceRecordSets


baseDomain = '<Enter base domain here>'
subDomain = '<enter subdomain here>'
fullDomain = subDomain + "." + baseDomain
FQDN = baseDomain + "."

#This section will be added in a future release once I figure out how to parse the aws credentials file the way I want to.
#credsFile = file(path.abspath(getenv("HOME") + "/.aws/credentials"))
#for line in credsFile:
#    if 'AWS_ACCESS_KEY_ID':
#        AWS_ACCESS_KEY_ID = line.strip('AWS_ACCESS_KEY_ID = ')
#        print(AWS_ACCESS_KEY_ID)

AWS_ACCESS_KEY_ID = '<Enter AWS_ACCESS_KEY_ID HERE>'
AWS_SECRET_ACCESS_KEY = '<Enter AWS_SECRET_ACCESS_KEY HERE>'
AWS_R53_ADDR_1 = "<Enter first address to update>"
AWS_R53_ADDR_2 = "<Enter second address to update>"

amazonIpCheck = StringIO()
c = pycurl.Curl()
c.setopt(c.URL, 'http://checkip.amazonaws.com')
c.setopt(c.IPRESOLVE, c.IPRESOLVE_V4)
c.setopt(c.WRITEDATA, amazonIpCheck)
c.perform()
c.close()

ip = amazonIpCheck.getvalue()

currentValue = socket.gethostbyname(subDomain + "." + baseDomain)
currentValue = currentValue.replace("'", "")
ip = ip.replace("\n", "")
print(currentValue)
print(ip)

if currentValue != ip:
    print("IP addresses do not match. Jumping into update function")

    r53 = boto.route53.connect_to_region('universal',
        aws_access_key_id=AWS_ACCESS_KEY_ID,
        aws_secret_access_key=AWS_SECRET_ACCESS_KEY)

    zone = r53.get_hosted_zone_by_name(FQDN)
    zone_id = zone.Id
    zone_id = zone_id.strip('/hostedzone/')

    records = r53.get_all_rrsets(zone_id,'A',AWS_R53_ADDR_1,maxitems=1)[0]
    print(records)
    r53rr = ResourceRecordSets(r53, zone_id)
    print(zone_id)
    d_record = r53rr.add_change("DELETE", AWS_R53_ADDR_1, "A", 300)
    d_record.add_value(currentValue)
    c_record = r53rr.add_change("CREATE", AWS_R53_ADDR_1, "A", 300)
    c_record.add_value(ip)
    print(d_record)
    print(c_record)
    r53rr.commit()

else:
    print("Ip addresses match. Not changing anything")
