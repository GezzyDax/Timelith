#!/bin/bash
# Setup git hooks for automated pre-commit checks
set -e

echo "🔧 Setting up git hooks..."
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Get the git directory
GIT_DIR=$(git rev-parse --git-dir 2>/dev/null)

if [ -z "$GIT_DIR" ]; then
    echo "❌ Error: Not a git repository"
    exit 1
fi

HOOKS_DIR="$GIT_DIR/hooks"

# Create hooks directory if it doesn't exist
mkdir -p "$HOOKS_DIR"

# Create pre-commit hook
PRE_COMMIT_HOOK="$HOOKS_DIR/pre-commit"

echo "Creating pre-commit hook..."

cat > "$PRE_COMMIT_HOOK" << 'EOF'
#!/bin/bash
# Auto-generated pre-commit hook
# This hook runs quality checks before allowing a commit

echo ""
echo "🔍 Running pre-commit checks..."
echo ""

# Run the pre-commit script
if ! make pre-commit; then
    echo ""
    echo "❌ Pre-commit checks failed!"
    echo ""
    echo "Options:"
    echo "  1. Fix the issues and try again"
    echo "  2. Run 'make pre-commit' to see detailed output"
    echo "  3. Skip hook with 'git commit --no-verify' (not recommended)"
    echo ""
    exit 1
fi

echo "✅ Pre-commit checks passed!"
echo ""
EOF

# Make the hook executable
chmod +x "$PRE_COMMIT_HOOK"

echo -e "${GREEN}✓ Pre-commit hook created${NC}"
echo ""

# Create commit-msg hook for commit message validation
COMMIT_MSG_HOOK="$HOOKS_DIR/commit-msg"

echo "Creating commit-msg hook..."

cat > "$COMMIT_MSG_HOOK" << 'EOF'
#!/bin/bash
# Auto-generated commit-msg hook
# Validates commit message format

commit_msg_file=$1
commit_msg=$(cat "$commit_msg_file")

# Regex for conventional commits
conventional_commit_regex="^(feat|fix|docs|style|refactor|test|chore|ci|perf|build|revert)(\(.+\))?: .{1,100}"

# Check if commit message follows conventional commit format
if ! echo "$commit_msg" | grep -qE "$conventional_commit_regex"; then
    echo ""
    echo "❌ Invalid commit message format!"
    echo ""
    echo "Commit messages should follow the Conventional Commits specification:"
    echo ""
    echo "  <type>[optional scope]: <description>"
    echo ""
    echo "Types:"
    echo "  feat:     New feature"
    echo "  fix:      Bug fix"
    echo "  docs:     Documentation changes"
    echo "  style:    Code style changes (formatting, etc.)"
    echo "  refactor: Code refactoring"
    echo "  test:     Test changes"
    echo "  chore:    Build process or tooling changes"
    echo "  ci:       CI/CD changes"
    echo "  perf:     Performance improvements"
    echo ""
    echo "Examples:"
    echo "  feat: add user authentication"
    echo "  fix: resolve database connection issue"
    echo "  docs: update API documentation"
    echo "  refactor(api): simplify error handling"
    echo ""
    echo "Your commit message:"
    echo "  $commit_msg"
    echo ""
    echo "To skip this check, use: git commit --no-verify"
    echo ""
    exit 1
fi
EOF

chmod +x "$COMMIT_MSG_HOOK"

echo -e "${GREEN}✓ Commit-msg hook created${NC}"
echo ""

# Summary
echo "╔══════════════════════════════════════════╗"
echo "║   ✅ Git Hooks Installed                ║"
echo "╚══════════════════════════════════════════╝"
echo ""
echo "Installed hooks:"
echo -e "  ${BLUE}pre-commit${NC}  - Runs quality checks before commit"
echo -e "  ${BLUE}commit-msg${NC}  - Validates commit message format"
echo ""
echo "These hooks will run automatically on:"
echo "  • git commit"
echo ""
echo "To skip hooks (not recommended):"
echo "  git commit --no-verify"
echo ""
echo -e "${YELLOW}Note: These hooks only affect your local repository${NC}"
echo ""
