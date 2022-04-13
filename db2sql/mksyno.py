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
        with open(f, 'a') as sql_file:
            out_format = """SET CURRENT sqlid = '{srod}';
DROP TABLE {tab};
COMMIT;
CREATE SYNONYM {tab} FOR {cen}.{tab};
COMMIT;
"""
            for t in tab_names:
                for r in roczne:
                    sql_file.write(out_format.format(srod=r,tab=t,cen=cen))

if __name__ == '__main__':
    try:
        main()
    except FileNotFoundError as e:
        print(e.strerror)
    except Exception as e:
        print(e)
