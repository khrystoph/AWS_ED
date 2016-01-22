
def get_credentials():
    import re

    file = open('/<path to user's home directory>/.aws/credentials')

    for line in file:
        if 'default' in line:
            for line2 in file :
                if re.search('^#.*', line2):
                    continue
                elif re.search('aws_secret_access_key', line2):
                    trash, aws_secret_access_key = line2.rsplit(' = ')
                    aws_secret_access_key = aws_secret_access_key.strip()
                elif re.search('aws_access_key_id', line2):
                    trash, aws_access_key_id = line2.rsplit(' = ')
                    aws_access_key_id = aws_access_key_id.strip()
        #if aws_secret_key_id is not None:
            #if aws_secret_access_key is not None:
                #break
    file.close()
    return (aws_access_key_id,aws_secret_access_key)