import requests
import os
import time
import re
import logging
import ipaddress
from pythonjsonlogger import jsonlogger


log = logging.getLogger()
__log_handler = logging.StreamHandler()
__formatter = jsonlogger.JsonFormatter('%(asctime)s %(levelname)s %(message)s')
__log_handler.setFormatter(__formatter)
log.addHandler(__log_handler)
log.setLevel(logging.WARN)


def parse_timeout(timeout: str) -> float:
    regexp = r'(\d+[H|h])|(\d+[M|m])|(\d+[S|s])'

    t = 1.0

    if not timeout:
        raise ValueError('Timeout was not provided')

    m = re.match(regexp, timeout)
    if not m:
        raise ValueError('Timeout has wrong format')
    
    if 'H' in timeout or 'h' in timeout:
        t = t * 60*60
    elif 'M' in timeout or 'm' in timeout:
        t = t * 60
    elif 'S' in timeout or 's' in timeout:
        pass

    return t * float(timeout[:-1])


class Updater():
    def __init__(self):
        self.config = dict()
        self.config['token'] = os.getenv('DUCKDNS_TOKEN')
        self.config['domains'] = os.getenv('DOMAINS')
        

        lvl = os.getenv('LOG_LEVEL')
        if lvl:
            log.setLevel(lvl)

        try:
            self.timeout = parse_timeout(os.getenv('TIMEOUT'))
        except ValueError as e:
            log.warning(f'{e}: using default (6H)')
            self.timeout = parse_timeout('6H')

        self.ip_api_uri = r'https://api64.ipify.org'


    def __get_ip(self):
        resp = requests.get(self.ip_api_uri)
        if resp.status_code == 200:
            self.config['ip'] = resp.text
        else:
            log.error(f'unable to get ipv4: {resp.reason}')


    def __parse_update_uri(self) -> str:
        self.__get_ip()

        root = 'https://www.duckdns.org/update?'

        params = list()

        params.append(f"domains={self.config['domains']}")
        params.append(f"token={self.config['token']}")

        try:
            ip = ipaddress.ip_address(self.config['ip'])
            if type(ip) == ipaddress.IPv4Address:
                params.append(f"ipv4={self.config['ip']}")
            elif type(ip) == ipaddress.IPv6Address:
                params.append(f"ipv6={self.config['ip']}")    
        except ValueError or TypeError:
            raise ValueError('wrong IP provided')
            
        uri = root + "&".join(params)

        log.debug(uri)

        return uri


    def __update(self):
        while True:
            resp = requests.get(self.__parse_update_uri())
            if resp.text == 'KO':
                log.error(f'unable to update ip')
                log.debug(self.config)
            log.info('sleeping...')
            time.sleep(self.timeout)


    def run(self) -> None:
        self.__get_ip()
        self.__update()


def main() -> None:
    try:
        upd = Updater()
        upd.run()
    except KeyboardInterrupt:
        exit(0)


if __name__ == '__main__':
    main()
