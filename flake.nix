{
  description = "GoReleaser-Wizard — Interactive GoReleaser configuration wizard";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

    systems.url = "github:nix-systems/default";

    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };

    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs@{ self, flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import inputs.systems;

      imports = [ inputs.treefmt-nix.flakeModule ];

      perSystem =
        {
          config,
          pkgs,
          lib,
          ...
        }:
        let
          version = self.rev or self.dirtyRev or "dev";

          vendorHash = "sha256-XGVIudRhIVD0py9gBXxfPSgaVZAd5mCfY1sw9fSuESc=";

          goreleaser-wizard = pkgs.buildGoModule {
            pname = "goreleaser-wizard";
            inherit version vendorHash;

            src = lib.fileset.toSource {
              root = ./.;
              fileset = lib.fileset.unions [
                ./cmd
                ./internal
                ./assets
                ./go.mod
                ./go.sum
              ];
            };

            subPackages = [ "cmd/goreleaser-wizard" ];

            ldflags = [
              "-s"
              "-w"
              "-X main.version=${version}"
              "-X main.commit=${self.rev or "dirty"}"
              "-X main.date=${toString self.lastModified}"
            ];

            meta = {
              description = "Interactive GoReleaser configuration wizard";
              homepage = "https://github.com/LarsArtmann/GoReleaser-Wizard";
              license = lib.licenses.mit;
              mainProgram = "goreleaser-wizard";
              maintainers = [
                {
                  name = "Lars Artmann";
                  github = "LarsArtmann";
                }
              ];
              platforms = lib.platforms.linux ++ lib.platforms.darwin;
            };
          };
        in
        {
          packages = {
            default = goreleaser-wizard;
            inherit goreleaser-wizard;
          };

          devShells = {
            default = pkgs.mkShell {
              inputsFrom = [ config.packages.default ];

              packages = with pkgs; [
                go
                gotools
                golangci-lint
                gofumpt
                govulncheck
                goreleaser
                jq
                yq-go
                git
                gh
              ];

              GOFLAGS = "-mod=mod";
            };

            ci = pkgs.mkShellNoCC {
              inputsFrom = [ config.packages.default ];

              packages = with pkgs; [
                go
                gotools
                golangci-lint
                gofumpt
                govulncheck
                jq
              ];
            };
          };

          apps.default = {
            type = "app";
            program = lib.getExe config.packages.default;
          };

          checks = {
            format = config.treefmt.build.check self;
            build = config.packages.default;

            test = config.packages.default.overrideAttrs (_: {
              doCheck = true;
            });
          };

          treefmt.config = {
            projectRootFile = "go.mod";
            programs = {
              nixfmt.enable = true;
              templ.enable = true;
              gofmt.enable = true;
            };
          };
        };

      flake.overlays.default = final: _prev: {
        goreleaser-wizard = final.callPackage ./nix/package.nix { };
      };
    };
}
