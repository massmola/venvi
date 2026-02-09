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
            export PATH=$PATH:$GOPATH/bin
            echo "Venvi PocketBase development environment loaded"

            # Install pre-commit hooks
            ${pre-commit-check.shellHook}
          '';
        };
      }
    );
}
