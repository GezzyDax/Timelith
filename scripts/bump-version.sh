#!/bin/bash
# Manual version bump script for local development
set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Get current version
CURRENT_VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
echo -e "${YELLOW}Current version: $CURRENT_VERSION${NC}"
echo ""

# Ask for bump type
echo "Select version bump type:"
echo "  1) patch (0.0.X) - bug fixes"
echo "  2) minor (0.X.0) - new features"
echo "  3) major (X.0.0) - breaking changes"
echo "  4) custom - specify version manually"
echo ""
read -p "Enter choice [1-4]: " CHOICE

# Parse current version
VERSION_NUMBER=${CURRENT_VERSION#v}
IFS='.' read -r MAJOR MINOR PATCH <<< "$VERSION_NUMBER"

case $CHOICE in
  1)
    PATCH=$((PATCH + 1))
    BUMP_TYPE="patch"
    ;;
  2)
    MINOR=$((MINOR + 1))
    PATCH=0
    BUMP_TYPE="minor"
    ;;
  3)
    MAJOR=$((MAJOR + 1))
    MINOR=0
    PATCH=0
    BUMP_TYPE="major"
    ;;
  4)
    read -p "Enter custom version (without 'v' prefix): " CUSTOM_VERSION
    IFS='.' read -r MAJOR MINOR PATCH <<< "$CUSTOM_VERSION"
    BUMP_TYPE="custom"
    ;;
  *)
    echo -e "${RED}Invalid choice${NC}"
    exit 1
    ;;
esac

NEW_VERSION="v${MAJOR}.${MINOR}.${PATCH}"

echo ""
echo -e "${YELLOW}New version will be: $NEW_VERSION${NC}"
read -p "Continue? (y/n): " CONFIRM

if [ "$CONFIRM" != "y" ]; then
    echo "Cancelled"
    exit 0
fi

echo ""
echo "📝 Updating version files..."

# Update Go version file
mkdir -p go-backend/internal/version
cat > go-backend/internal/version/version.go << EOF
package version

// Version is the current version of the application
const Version = "${MAJOR}.${MINOR}.${PATCH}"
EOF

echo -e "${GREEN}✓ Updated go-backend/internal/version/version.go${NC}"

# Update package.json
cd web-ui
npm version ${MAJOR}.${MINOR}.${PATCH} --no-git-tag-version
cd ..

echo -e "${GREEN}✓ Updated web-ui/package.json${NC}"

echo ""
echo "📋 Generating changelog..."

# Generate changelog
CHANGELOG_FILE="CHANGELOG_${NEW_VERSION}.md"
echo "# Release $NEW_VERSION" > $CHANGELOG_FILE
echo "" >> $CHANGELOG_FILE
echo "## Changes" >> $CHANGELOG_FILE
echo "" >> $CHANGELOG_FILE

git log ${CURRENT_VERSION}..HEAD --pretty=format:"%s" | while read line; do
    if echo "$line" | grep -qiE "^feat|^feature"; then
        echo "- ✨ $line" >> $CHANGELOG_FILE
    elif echo "$line" | grep -qiE "^fix"; then
        echo "- 🐛 $line" >> $CHANGELOG_FILE
    elif echo "$line" | grep -qiE "^docs"; then
        echo "- 📚 $line" >> $CHANGELOG_FILE
    elif echo "$line" | grep -qiE "^refactor"; then
        echo "- ♻️ $line" >> $CHANGELOG_FILE
    elif echo "$line" | grep -qiE "^test"; then
        echo "- 🧪 $line" >> $CHANGELOG_FILE
    elif echo "$line" | grep -qiE "^chore|^ci"; then
        echo "- 🔧 $line" >> $CHANGELOG_FILE
    else
        echo "- $line" >> $CHANGELOG_FILE
    fi
done

echo ""
echo "Changelog saved to: $CHANGELOG_FILE"
cat $CHANGELOG_FILE

echo ""
read -p "Create git tag and commit? (y/n): " CREATE_TAG

if [ "$CREATE_TAG" = "y" ]; then
    # Commit version changes
    git add -A
    git commit -m "chore: bump version to $NEW_VERSION"

    # Create tag
    git tag -a $NEW_VERSION -m "Release $NEW_VERSION"

    echo -e "${GREEN}✓ Created commit and tag${NC}"
    echo ""
    echo "To push changes and tag:"
    echo "  git push origin main"
    echo "  git push origin $NEW_VERSION"
else
    echo "Version files updated but not committed"
    echo ""
    echo "To commit manually:"
    echo "  git add -A"
    echo "  git commit -m 'chore: bump version to $NEW_VERSION'"
    echo "  git tag -a $NEW_VERSION -m 'Release $NEW_VERSION'"
fi

# Cleanup
rm -f $CHANGELOG_FILE

echo ""
echo -e "${GREEN}✅ Version bump complete!${NC}"
