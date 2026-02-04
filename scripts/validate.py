import os
import subprocess
import sys


def run_step(name: str, command: list[str], show_output: bool = False) -> None:
    """Run a validation step and handle output."""
    print(f" -> {name}...")
    try:
        if show_output:
            subprocess.run(command, check=True)
        else:
            subprocess.run(command, check=True, capture_output=True, text=True)
    except subprocess.CalledProcessError as e:
        print(f"❌ {name} failed!")
        if not show_output:
            if e.stdout:
                print("\nStandard Output:")
                print(e.stdout)
            if e.stderr:
                print("\nError Output:")
                print(e.stderr)
        sys.exit(e.returncode)


def main() -> None:
    """Run all validation steps."""
    # Ensure we are in the project root
    try:
        repo_root = subprocess.check_output(
            ["git", "rev-parse", "--show-toplevel"], text=True
        ).strip()
        os.chdir(repo_root)
    except Exception as e:
        print(f"Error finding project root: {e}")
        sys.exit(1)

    print("🛡️  Starting full project validation...")

    # 1. Database initialization
    run_step("[1/6] Initializing test database", ["bash", "scripts/init_db.sh"])

    # 2. Formatting
    run_step("[2/6] Checking formatting (Ruff)", ["ruff", "format", "--check", "."])

    # 3. Linting
    run_step("[3/6] Checking linting (Ruff)", ["ruff", "check", "."])

    # 4. Type Checking
    run_step("[4/6] Checking types (Mypy)", ["mypy", "."])

    # 5. Tests
    # We show output for tests to let pytest show its progress
    run_step("[5/6] Running tests (Pytest)", ["pytest"], show_output=True)

    # 6. Documentation
    run_step("[6/6] Verifying documentation (MkDocs)", ["mkdocs", "build", "--strict"])

    print("\n✅ All validations passed!")


if __name__ == "__main__":
    main()
