{
  description = "Venvi Development Environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils }:
    utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            # Go toolchain
            go
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
          ];

          shellHook = ''
            export GOPATH=$HOME/go
            export PATH=$GOPATH/bin:$PATH
            echo "Venvi PocketBase development environment loaded"
          '';
        };
      }
    );
}
