import jaydebeapi as jdb
import os

DB2_DRIVER = 'db2jcc4.jar'
INFORMIX_DRIVER = 'ifxjdbc.jar'

def driver():
    print(os.getenv('JDBC_DRIVERS'))


def main():
    driver()

if __name__ == '__main__':
    main()
