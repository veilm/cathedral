#!/usr/bin/env python3
"""Cathedral CLI - The command gateway to memories."""

import sys
import os

# Add the installation directory to Python path to find our module
INSTALL_DIR = "/usr/local/lib/cathedral"
if os.path.exists(INSTALL_DIR):
    sys.path.insert(0, INSTALL_DIR)

# Import the main module (which is the current cathedral.py)
from cathedral.main import main

if __name__ == "__main__":
    sys.exit(main())
