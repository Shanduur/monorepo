import argparse
import comparator as cmp

def main():
    parser = argparse.ArgumentParser(description='Check SQL files for col differences.')
    parser.add_argument('f1', metavar='F1', type=str,
                        help='an integer for the accumulator')
    parser.add_argument('f2', metavar='F2', type=str,
                        help='an integer for the accumulator')
    args = parser.parse_args()
    
    c = cmp.Comparator(args.f1, args.f2)
    c.compare()


if __name__ == '__main__':
    main()
