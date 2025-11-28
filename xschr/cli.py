"""
xschr.cli

command-line interface for the xschr package.
"""

import argparse
import sys
import os
import textwrap
from dataclasses import dataclass
from typing import TextIO, Optional

from .__version__ import __version__

# --- 1. The Environment Context ---
@dataclass
class Environment:
    """
    Represents the I/O environment for CLI operations.
    """
    stdin: TextIO = sys.stdin
    stdout: TextIO = sys.stdout
    stderr: TextIO = sys.stderr
    is_terminal: bool = sys.stdout.isatty()

    def log_error(self, message: str):
       """Write a formatted error message to stderr."""
       self.stderr.write(f"\033[91m[Error]\033[0m {message}\n")

# --- 2. The Custom Parser ---
class XSchrArgumentParser(argparse.ArgumentParser):
    """
    Argument parser with custom validation and error reporting.
    """
    def __init__(self, env: Environment = None, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.env = env or Environment()

    def parse_args(self, args=None, namespace=None):
        """Parse command-line arguments and apply post-processing hooks."""

        # 1. Standard parsing
        parsed_args = super().parse_args(args, namespace)
        
        # 2. Post-processing logic
        self._validate_path(parsed_args)
        self._process_verbosity(parsed_args)
        
        return parsed_args

    def _validate_path(self, args):
        """Ensure the config file actually exists before proceeding."""
        
        if not os.path.exists(args.path):
            self.error(f"Configuration file not found: '{args.path}'")

    def _process_verbosity(self, args):
        """Map generic flags to specific internal logic settings."""
 
        if args.debug:
            args.verbose = True

    # Override error to use our custom environment printer
    def error(self, message):
        self.print_usage(self.env.stderr)
        self.env.log_error(message)
        sys.exit(2)

# --- 3. The Definition (Groups & formatting) ---
def get_parser(env: Environment = None) -> XSchrArgumentParser:
    """Create and return the standard XSchr argument parser."""

    parser = XSchrArgumentParser(
        env=env,
        formatter_class=argparse.RawDescriptionHelpFormatter,
        description=textwrap.dedent("""

            XSchr: simple job scheduler for personal deep-learning research
            ---------------------------------------------------------------
            Runs a sequence of experiments defined in a JSON file,
            managing logs, errors, and resources automatically.
        """)
    )

    # -- Group: Essential Inputs --
    source_group = parser.add_argument_group(title="Source")
    source_group.add_argument(
        "-p", "--path",
        metavar="FILE",
        required=True,
        help="Path to the experiment configuration file (.json)"
    )

    # -- Group: Execution Control --
    exec_group = parser.add_argument_group(title="Execution Control")
    exec_group.add_argument(
        "--fail-fast",
        action="store_true",
        default=False,
        help="Stop the entire queue immediately if a single run fails."
    )
    exec_group.add_argument(
        "--dry-run",
        action="store_true",
        default=False,
        help="Simulate the execution plan without running any scripts."
    )

    # -- Group: Troubleshooting & Info --
    debug_group = parser.add_argument_group(title="Troubleshooting")
    debug_group.add_argument(
        "--debug",
        action="store_true",
        help="Print detailed diagnostic information for bug reports."
    )
    debug_group.add_argument(
        "-v", "--version",
        action="version",
        version=f"XSchr {__version__}",
        help="Show version and exit."
    )

    return parser
