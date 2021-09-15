import glob
import os
import pandas as pd
import re


DATA = 'data\\*'


def list_files(stage: str) -> list:
    files = glob.glob(stage + '\\*')

    return files


def get_stages() -> list:
    stages = []
    base = glob.glob(DATA)
    for b in base:
        for s in glob.glob(b + '\\*'):
            stages.append(s)
    
    return stages


def stage_to_xslx(stage: str):
    files = list_files(stage=stage)
    if not files:
        return

    writer = pd.ExcelWriter(stage.replace('\\', '_').replace('data_', '') + '.xlsx', engine='xlsxwriter')

    for f in files:
        if '.csv' not in f:
            continue
        try:
            df = pd.read_csv(f, dtype=object)
            df.to_excel(writer, sheet_name=re.sub(r'_[0-9]*.csv', r'', os.path.basename(f)), index=False)
        except pd.errors.EmptyDataError:
            continue

    writer.save()


def main():
    stages = get_stages()
    for s in stages:
        stage_to_xslx(stage=s)


if __name__ == '__main__':
    main()
