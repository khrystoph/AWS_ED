import pycurl
from StringIO import StringIO
import boto.route53
import socket
from os import path
from os import getenv
from boto.route53.record import ResourceRecordSets
import creds

#configuration section
baseDomain = '<hostname here>'
subDomain = '<sub domain here if you have one>'
if subDomain is not '':
    fullDomain = subDomain + "." + baseDomain
else:
    fullDomain = baseDomain
ttl = 300
record_type = 'AAAA'
if record_type is 'A':
    socktype = socket.AF_INET
elif record_type is "AAAA":
    socktype = socket.AF_INET6
#end configuration section

#pull credentials from the configured aws-cli tool
AWS_ACCESS_KEY_ID,AWS_SECRET_ACCESS_KEY = creds.get_credentials()

amazonIpCheck = StringIO()
c = pycurl.Curl()
c.setopt(c.URL, 'https://icanhazip.com')
if record_type is 'AAAA':
    c.setopt(c.IPRESOLVE, c.IPRESOLVE_V6)
elif record_type is 'A':
    c.setopt(c.IPRESOLVE, c.IPRESOLVE_V4)
c.setopt(c.WRITEDATA, amazonIpCheck)
c.perform()
c.close()

ip = amazonIpCheck.getvalue()

host = socket.getaddrinfo(fullDomain, None, socktype)
currentValue = host[0][4][0]

ip = ip.replace("\n", "")
print(currentValue)
print(ip)

if currentValue != ip:
    print("IP addresses do not match. Jumping into update function")

    r53 = boto.route53.connect_to_region('universal',
        aws_access_key_id=AWS_ACCESS_KEY_ID,
        aws_secret_access_key=AWS_SECRET_ACCESS_KEY)

    zone = r53.get_hosted_zone_by_name(fullDomain)
    zone_id = zone.Id
    zone_id = zone_id.strip('/hostedzone/')

    records = r53.get_all_rrsets(zone_id,record_type,fullDomain,maxitems=1)[0]
    print(records)
    r53rr = ResourceRecordSets(r53, zone_id)
    print(zone_id)
    d_record = r53rr.add_change("DELETE", fullDomain, record_type, ttl)
    d_record.add_value(currentValue)
    c_record = r53rr.add_change("CREATE", fullDomain, record_type, ttl)
    c_record.add_value(ip)
    print(d_record)
    print(c_record)
    r53rr.commit()

else:
    print("Ip addresses match. Not changing anything")
