{
  description = "Fleet - Ephemeral HTTP API tunnel";
  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixos-23.11";
  outputs = { self, nixpkgs, ... }@inputs:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
      version = "24.01";
    in
    {
      formatter.${system} = pkgs.nixpkgs-fmt;
      packages.${system} = rec {
        fleet = pkgs.buildGoModule {
          pname = "fleet";
          version = version;
          src = ./.;
          vendorHash = "sha256-TM9ZfQ/YfSYdLbdhDgJUMHbhFfm+iOc70zVRKclAjf0=";
          CGO_ENABLED = 0;
          subPackages = [ "cmd/fleet" ];
          nativeBuildInputs = [ pkgs.installShellFiles ];
          postInstall = ''
            for shell in bash zsh; do
              installShellCompletion --$shell ./cli/completion/$shell/fleet
            done
          '';
        };
        dockerImage = pkgs.dockerTools.buildImage {
          name = "murtazau/fleet";
          tag = version;
          config = {
            Cmd = [ "${fleet}/bin/fleet" "run" ];
            WorkingDir = "/data";
          };
        };
        default = fleet;
      };
      devShells.${system}.default = pkgs.mkShell {
        packages = with pkgs; [
          go
          go-tools
          gopls
          protobuf
          protoc-gen-go
          protoc-gen-go-grpc
          grpcurl
          openssl
        ];
      };
    };
}
