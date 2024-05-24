{
  inputs = {
    # Candidate channels
    #   - https://github.com/kachick/anylang-template/issues/17
    #   - https://discourse.nixos.org/t/differences-between-nix-channels/13998
    # How to update the revision
    #   - `nix flake update --commit-lock-file` # https://nixos.org/manual/nix/stable/command-ref/new-cli/nix3-flake-update.html
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        updaterVersion = if (self ? shortRev) then self.shortRev else "dev";
      in
      rec {
        formatter = pkgs.nixfmt-rfc-style;
        devShells.default =
          with pkgs;
          mkShell {
            buildInputs = [
              # https://github.com/NixOS/nix/issues/730#issuecomment-162323824
              # https://github.com/kachick/dotfiles/pull/228
              bashInteractive

              nil
              nixfmt-rfc-style
              dprint
              yamlfmt
              typos
              go-task
              go_1_22
            ];
          };

        packages.gwurl = pkgs.buildGo122Module {
          pname = "gwurl";
          version = updaterVersion;
          src = pkgs.lib.cleanSource self;

          # When updating go.mod or go.sum, update this sha together as following
          # vendorHash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";
          # (`pkgs.lib.fakeSha256` returns invalid string in thesedays... :<)
          vendorHash = "sha256-i07qu5jt4XTD2YxorJlAY/Kq1zk4yfiUSlSr3toBBGA=";
        };

        packages.default = packages.gwurl;

        # `nix run`
        apps.default = {
          type = "app";
          program = "${packages.gwurl}/bin/gwurl";
        };
      }
    );
}
