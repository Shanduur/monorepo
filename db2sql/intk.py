import sys
import os


def main():
    files = []
    if len(sys.argv) >= 1:
        files = [f if f.endswith('.sql') else os.listdir(f) if os.path.isdir(f) else f for f in sys.argv[1:]]
    else:
        files = [f for f in os.listdir('.') if os.path.isfile(f)]

    files = [f for f in files if f.endswith('.sql')]

    for f in files:
        out = []
        with open(f, 'r') as sql_file:
            lines = sql_file.readlines()
            sql_lines = "".join(lines).split(';')
            out_lines = []
            for l in sql_lines:
                if 'CREATE TABLE' in l.upper() and 'IN DBKT2.TSKT1S04' not in l.upper():
                    out_lines.append(f'{l} IN DBKT2.TSKT1S04')
                else:
                    out_lines.append(l)

            out = ";".join(out_lines)
        
        with open(f, 'w') as sql_file:
            sql_file.write(out)

if __name__ == '__main__':
    try:
        main()
    except FileNotFoundError as e:
        print(e.strerror)
    except Exception as e:
        print(e)
