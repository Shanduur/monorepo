import sys
import os
import re


def main():
    files = []
    if len(sys.argv) >= 1:
        files = [f if f.endswith('.sql') else os.listdir(f) if os.path.isdir(f) else f for f in sys.argv[1:]]
    else:
        files = [f for f in os.listdir('.') if os.path.isfile(f)]

    files = [f for f in files if f.endswith('.sql')]

    for f in files:
        out = None
        with open(f, 'r') as sql_file:
            lines = sql_file.readlines()
            sql_lines = "".join(lines).split(';')
            commit_protected_lines = []
            i = 0
            j = 0
            for l in sql_lines:
                if "COMMIT" in l.upper() or "SET CURRENT" in l.upper() or "CREATE SYNONYM" in l.upper() or "DROP TABLE" in l.upper():
                    commit_protected_lines.append(f'{l}/* dedup protect {i} */')
                    i += 1
                else:
                    commit_protected_lines.append(l)
            dedup_lines = list(dict.fromkeys(commit_protected_lines))
            out = re.sub(r'\/\*\s([\w|\s]+)\s\*\/', '', ";".join(dedup_lines))
        
        with open(f, 'w') as sql_file:
            sql_file.write(out)

if __name__ == '__main__':
    try:
        main()
    except FileNotFoundError as e:
        print(e.strerror)
    except Exception as e:
        print(e)
