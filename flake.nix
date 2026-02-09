{
  description = "Venvi Development Environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
    pre-commit-hooks.url = "github:cachix/pre-commit-hooks.nix";
  };

  outputs = { self, nixpkgs, utils, pre-commit-hooks, ... } @ inputs:
    utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };

        pre-commit-check = inputs.pre-commit-hooks.lib.${system}.run {
          src = ./.;
          hooks = {
            # Code quality
            gofmt.enable = true;
            nixpkgs-fmt.enable = true;

            # Custom check for fmt.Sprintf
            no-fmt-sprintf = {
              enable = true;
              name = "no-fmt-sprintf";
              description = "Ban fmt.Sprintf in favor of safe alternatives";
              entry = "${pkgs.lib.getBin pkgs.bash}/bin/bash -c 'if grep -r \"fmt.Sprintf\" . --include=*.go --exclude-dir=vendor; then echo \"Error: fmt.Sprintf is banned. Use strings.ReplaceAll, text/template, or concatenation.\"; exit 1; fi'";
              pass_filenames = false;
            };

            # Security
            gitleaks = {
              enable = true;
              name = "gitleaks";
              description = "Detect hardcoded secrets";
              entry = "${pkgs.gitleaks}/bin/gitleaks detect --source . --no-git -v";
            };
          };
        };
      in
      {
        checks = {
          pre-commit-check = pre-commit-check;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            # Go toolchain
            go_1_25
            gopls
            golangci-lint

            # Build tools
            gcc
            pkg-config

            # Node.js for Playwright & Frontend tooling
            nodejs_20

            # CI/CD & Deployment
            github-cli

            # SQLite (for PocketBase)
            sqlite

            # Security tools
            gitleaks
          ];

          shellHook = ''
            export GOPATH=$HOME/go
            export PATH=$GOPATH/bin:$PATH
            echo "Venvi PocketBase development environment loaded"

            # Install pre-commit hooks
            ${pre-commit-check.shellHook}
          '';
        };
      }
    );
}
