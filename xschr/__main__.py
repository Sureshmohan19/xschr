"""
xschr.__main__

the main entry point. invoke as `xschr' or `python -m xschr' from the terminal.
"""

import sys

def main():
    """
    Wrapper to handle strict exit codes and keyboard interrupts.
    """
    try:
        # Lazy import: We don't load the heavy logic until we are sure
        # the user actually wants to run the program.
        from xschr.core import main as xschr_main
        exit_status = xschr_main()
    except KeyboardInterrupt:
        # Standard Unix convention: 128 + SIGINT (2) = 130
        exit_status = 130
        print("\n\n[!] xschr execution cancelled.", file=sys.stderr)

    return exit_status

if __name__ == '__main__':
    sys.exit(main())
