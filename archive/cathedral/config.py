"""Cathedral configuration - the sacred settings of the memory temple."""

from dataclasses import dataclass
from typing import Optional


# ============================================================================
# Memory Compression Settings - for perfect memory consolidation
# ============================================================================

DEFAULT_COMPRESSION_RATIO = 0.5  # 50% retention (2x compression)
DEFAULT_ROUNDING = 50  # Round to nearest 50 chars/words for clean numbers

# Future compression profiles (from config.md)
COMPRESSION_PROFILES = {
    "default": 0.5,  # Balanced: 50% retention
    "compact": 0.25,  # Aggressive: 25% retention
    "verbose": 0.75,  # Gentle: 75% retention
    "full": 1.0,  # No compression (for testing)
}


# ============================================================================
# Memory Node Settings - structure of our knowledge cathedral
# ============================================================================

DEFAULT_NODE_SIZE = 1000  # Target tokens per node (soft limit)
MAX_NODE_SIZE = 2000  # Maximum before splitting required
INDEX_NODE_SIZE = 3000  # Special size for index.md

# Node types and their characteristics
NODE_TYPES = {
    "episodic": {
        "max_size": DEFAULT_NODE_SIZE,
        "link_to_raw": True,
        "chronological": True,
    },
    "semantic": {
        "max_size": DEFAULT_NODE_SIZE,
        "link_to_source": True,  # Link back to episodic
        "organized_by": "topic",
    },
}


# ============================================================================
# Session Management - tracking our journeys
# ============================================================================

SESSION_ID_ALPHABET = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"  # After Z, continues with AA, AB...
DEFAULT_ROLE_NAMES = {
    "world": "External input (human, system, environment)",
    "self": "Model's own thoughts and responses",
}


# ============================================================================
# File Paths and Templates - where the sacred texts reside
# ============================================================================

DEFAULT_TEMPLATES = {
    "write_memory": "grimoire/write-memory.md",
    "init_session": "grimoire/conv-start-injection.md",
    "blank_index": "grimoire/index-blank.md",
}

# Directory structure
STORE_STRUCTURE = {
    "index": "index.md",
    "episodic_raw": "episodic-raw/",
    "episodic": "episodic/",
    "semantic": "semantic/",
}


# ============================================================================
# Runtime Configuration - for dynamic navigation
# ============================================================================


@dataclass
class CathedralSettings:
    """Runtime settings that can be adjusted per operation."""

    compression_ratio: float = DEFAULT_COMPRESSION_RATIO
    compression_profile: Optional[str] = None

    # Memory writing settings
    preserve_quotes: bool = True  # Keep exact phrasing when significant
    link_all_episodic: bool = False  # Link every episodic mention vs selective

    # Reflection settings (future)
    reflection_enabled: bool = False
    reflection_frequency: int = 5  # After N sessions

    # Debug/testing
    verbose: bool = False
    dry_run: bool = False

    def __post_init__(self):
        """Apply compression profile if specified."""
        if (
            self.compression_profile
            and self.compression_profile in COMPRESSION_PROFILES
        ):
            self.compression_ratio = COMPRESSION_PROFILES[self.compression_profile]


# ============================================================================
# Future Experimental Settings (from config.md considerations)
# ============================================================================

# Option A vs B vs C vs D for conversation continuity
CONTINUITY_MODE = "default"  # Could be: "rolling", "artificial", "explicit", "none"

# Whether all semantic memory must link to episodic source
REQUIRE_SEMANTIC_SOURCES = True  # From config.md default

# How much context to include when splitting nodes
SPLIT_CONTEXT_PRESERVATION = 0.1  # Keep 10% overlap when splitting
