{
  lib,
  buildGoModule,
  fetchFromGitHub,
  cosign,
  gnupg,
}:

buildGoModule (finalAttrs: {
  pname = "goreleaser-wizard";
  version = "0.1.0";

  src = fetchFromGitHub {
    owner = "LarsArtmann";
    repo = "GoReleaser-Wizard";
    rev = "v${finalAttrs.version}";
    hash = "sha256-l7mylIZ5iAEkUFEMyqQL5cAvjFBa4gj+n6Qq9fFdUAE=";
  };

  vendorHash = "sha256-CTXzZ/yGuuo5T+CtNROah0AIXp5mKwk2OcoZTyfGJnQ=";

  # v0.1.0 validation treats missing gpg/cosign as a hard error (upstream
  # later downgraded this to a warning); unit tests probe PATH for them.
  nativeCheckInputs = [
    gnupg
    cosign
  ];

  subPackages = [ "cmd/goreleaser-wizard" ];

  ldflags = [
    "-s"
    "-w"
    "-X main.version=${finalAttrs.version}"
    "-X main.commit=${finalAttrs.src.rev}"
    "-X main.date=1970-01-01T00:00:00Z"
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
})
