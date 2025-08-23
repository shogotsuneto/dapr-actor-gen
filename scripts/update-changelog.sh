#!/bin/bash

# Update changelog by converting "Unreleased" to a versioned section and adding a new "Unreleased" section
# Usage: ./update-changelog.sh <version> [path-to-changelog]

set -e

VERSION="$1"
CHANGELOG_FILE="${2:-CHANGELOG.md}"

if [[ -z "$VERSION" ]]; then
    echo "❌ Version is required"
    echo "Usage: $0 <version> [path-to-changelog]"
    exit 1
fi

if [[ ! -f "$CHANGELOG_FILE" ]]; then
    echo "❌ Changelog file not found: $CHANGELOG_FILE"
    exit 1
fi

# Remove 'v' prefix if present for the changelog
CLEAN_VERSION="${VERSION#v}"

# Get current date in YYYY-MM-DD format
CURRENT_DATE=$(date +%Y-%m-%d)

echo "📝 Updating changelog: $CHANGELOG_FILE"
echo "🏷️  Version: $CLEAN_VERSION"
echo "📅 Date: $CURRENT_DATE"

# Create a backup
cp "$CHANGELOG_FILE" "${CHANGELOG_FILE}.backup"

# Use awk to update the changelog
awk -v version="[$CLEAN_VERSION]" -v date="$CURRENT_DATE" '
BEGIN { 
    updated = 0
    print "# Changelog"
    print ""
    print "All notable changes to this project will be documented in this file."
    print ""
    print "## [Unreleased]"
    print ""
}
/^# Changelog/ { next }
/^$/ && NR <= 4 { next }
/^All notable changes/ { next }
/^## \[Unreleased\]/ { 
    if (!updated) {
        print "## " version " - " date
        updated = 1
    }
    next 
}
{ print }
' "$CHANGELOG_FILE" > "${CHANGELOG_FILE}.tmp"

# Replace original file
mv "${CHANGELOG_FILE}.tmp" "$CHANGELOG_FILE"

echo "✅ Changelog updated successfully"
echo "📄 Backup saved as ${CHANGELOG_FILE}.backup"