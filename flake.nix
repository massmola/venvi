{
  description = "Venvi Development Environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.05";
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
            python312
            poetry
            postgresql
            zlib
            gcc
            pkg-config
            ruff
          ];

          shellHook = ''
            export LD_LIBRARY_PATH=${pkgs.lib.makeLibraryPath [ pkgs.stdenv.cc.cc pkgs.postgresql pkgs.zlib ]}
            export POETRY_VIRTUALENVS_IN_PROJECT=true
            export PYTHON_KEYRING_BACKEND=keyring.backends.null.Keyring
          '';
        };
      }
    );
}
