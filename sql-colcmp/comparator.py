import logging as l

class Comparator:
    def __init__(self, file1: str, file2: str) -> None:
        self.__f1 = open(file1)
        self.__f1 = open(file2)

    def compare(self) -> bool:
        if not self.__n_cols_equal():
            l.warn('length not equal, column is dangling')
        pass

    def __n_cols_equal(self) -> bool:
        return False

    def __prepare(self) -> list:
        return []

    def __find_dangling(self, col1: dict, col2: dict) -> dict:
        return {}
