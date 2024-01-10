{
  description = "Fleet - Ephemeral HTTP API tunnel";
  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixos-23.11";
  outputs = { self, nixpkgs, ... }@inputs:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in
    {
      formatter.${system} = pkgs.nixpkgs-fmt;
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
