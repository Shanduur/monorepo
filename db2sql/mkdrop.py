import sys
import os


def main():
    cen = os.getenv('SRODOWISKO_CENTRALNE')
    roczne = os.getenv('SRODOWISKA_ROCZNE').split(',')

    files = []
    if len(sys.argv) >= 1:
        files = [f if f.endswith('.sql') else os.listdir(f) if os.path.isdir(f) else f for f in sys.argv[1:]]
    else:
        files = [f for f in os.listdir('.') if os.path.isfile(f)]

    files = [f for f in files if f.endswith('.sql')]

    for f in files:
        tab_names = []
        lines = []
        with open(f, 'r') as sql_file:
            lines = sql_file.readlines()
            sql_lines = "".join(lines).split(';')
            for l in sql_lines:
                sl = l.split(' ')
                if len(sl) >= 3 and 'CREATE' in sl[0].upper() and 'TABLE' in sl[1].upper():
                    ssl = sl[2].split('(', 1) if "(" in sl[2] else sl[2]
                    if len(ssl) >= 2:
                        tab_names.append(ssl)

        tab_names = list(dict.fromkeys(tab_names))
        with open(f, 'w') as sql_file:
            sql_file.write(f"SET CURRENT SQLID = '{cen}';\n")
            for t in tab_names:
                sql_file.write(f"DROP TABLE {t}; COMMIT;\n")
            for l in lines:
                sql_file.write(l)

if __name__ == '__main__':
    try:
        main()
    except FileNotFoundError as e:
        print(e.strerror)
    except Exception as e:
        print(e)
