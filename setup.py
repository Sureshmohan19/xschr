#!/usr/bin/env python
# -*- coding: utf-8 -*-

import os
import sys
from shutil import rmtree
from setuptools import find_packages, setup, Command

# --- Configuration ---
# The name of the folder containing your code
PACKAGE_DIR = 'xschr' 
REQUIRES_PYTHON = '>=3.10.0'

# --- Python Version Check ---
# We check this before loading anything else
CURRENT_PYTHON = sys.version_info[:2]
REQUIRED_PYTHON_TUPLE = (3, 10)

if CURRENT_PYTHON < REQUIRED_PYTHON_TUPLE:
    sys.stderr.write("""

==========================
Unsupported Python version
==========================

XSchr requires at least Python {}.{}.
You are currently using Python {}.{}.
""".format(*(REQUIRED_PYTHON_TUPLE + CURRENT_PYTHON)))
    sys.exit(1)

# --- Load Metadata from __version__.py ---
here = os.path.abspath(os.path.dirname(__file__))
about = {}

# This reads xschr/__version__.py and puts the variables into the 'about' dictionary
with open(os.path.join(here, "xschr", "__version__.py"), "r", encoding="utf-8") as f:
    exec(f.read(), about)

# --- Upload Support ---
class UploadCommand(Command):
    """Support setup.py publish."""
    description = 'Build and publish the package.'
    user_options = []

    @staticmethod
    def status(s):
        print('\033[1m{0}\033[0m'.format(s))

    def initialize_options(self): pass
    def finalize_options(self): pass

    def run(self):
        try:
            self.status('Removing previous builds…')
            rmtree(os.path.join(here, 'dist'))
        except OSError:
            pass

        self.status('Building Source and Wheel…')
        os.system('{0} setup.py sdist bdist_wheel --universal'.format(sys.executable))

        self.status('Uploading the package to PyPI via Twine…')
        os.system('twine upload dist/*')

        self.status('Pushing git tags…')
        os.system('git tag v{0}'.format(about['__version__']))
        os.system('git push --tags')
        sys.exit()

# --- Setup ---
setup(
    # Metadata extracted from xschr/__version__.py
    name=about['__title__'],
    version=about['__version__'],
    description=about['__description__'],
    url=about['__url__'],
    author=about['__author__'],
    author_email=about['__author_email__'],
    license=about['__license__'],
    
    # Python constraints
    python_requires=REQUIRES_PYTHON,
    
    # Auto-discovery
    packages=find_packages(exclude=["tests"]),
    
    # Dependencies
    install_requires=[
    "PyYAML>=6.0",
    ],
    
    # The 'xschr' terminal command
    entry_points={
        'console_scripts': ['xschr=xschr.__main__:main'],
    },
    
    # Standard Classifiers
    classifiers=[
        'License :: OSI Approved :: MIT License',
        'Programming Language :: Python :: 3.10',
    ],
    
    # Upload helper
    cmdclass={
        'upload': UploadCommand,
    },
)
