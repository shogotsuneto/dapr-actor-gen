#!/bin/bash

# Extract changelog content from the "Unreleased" section
# Usage: ./extract-changelog.sh [path-to-changelog] [output-file]

set -e

CHANGELOG_FILE="${1:-CHANGELOG.md}"
OUTPUT_FILE="${2:-release-notes.md}"

if [[ ! -f "$CHANGELOG_FILE" ]]; then
    echo "❌ Changelog file not found: $CHANGELOG_FILE"
    exit 1
fi

echo "📝 Extracting changelog content from $CHANGELOG_FILE..."

# Extract content between "## [Unreleased]" and the next "## [" section
# This will capture all content under the Unreleased section
UNRELEASED_CONTENT=$(awk '
BEGIN { in_unreleased = 0; found_unreleased = 0 }
/^## \[Unreleased\]/ { 
    in_unreleased = 1
    found_unreleased = 1
    next 
}
/^## \[/ && in_unreleased { 
    in_unreleased = 0 
}
in_unreleased && NF > 0 { 
    print $0 
}
END { 
    if (!found_unreleased) {
        print "No Unreleased section found"
        exit 1
    }
}' "$CHANGELOG_FILE")

if [[ -z "$UNRELEASED_CONTENT" || "$UNRELEASED_CONTENT" == "No Unreleased section found" ]]; then
    echo "⚠️  No unreleased changes found in changelog"
    UNRELEASED_CONTENT="- No changes documented in this release"
fi

# Clean up the content - remove empty lines at start/end
UNRELEASED_CONTENT=$(echo "$UNRELEASED_CONTENT" | sed '/^$/d' | sed -e :a -e '/^\s*$/d;N;ba')

echo "✅ Extracted changelog content:"
echo "$UNRELEASED_CONTENT"

# Save to output file
echo "$UNRELEASED_CONTENT" > "$OUTPUT_FILE"
echo "📄 Saved changelog content to $OUTPUT_FILE"