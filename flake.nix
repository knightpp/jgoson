{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-23.05";
  };

  outputs = { self, nixpkgs }:
    let
      # Systems supported
      allSystems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];

      # Helper to provide system-specific attributes
      forAllSystems = f: nixpkgs.lib.genAttrs allSystems (system: f {
        pkgs = import nixpkgs { inherit system; };
      });
    in
    {
      # Development environment output
      devShells = forAllSystems ({ pkgs }: {
        default = pkgs.mkShell {
          # The Nix packages provided in the environment
          packages = with pkgs; [
            go_1_20
            gotools # Go tools like goimports, godoc, and others
          ];
        };
      });

      packages = forAllSystems
        ({ pkgs }: {
          default = pkgs.buildGoModule
            {
              name = "jgoson";
              src = pkgs.nix-gitignore.gitignoreSource [ ] ./.;
              # vendorSha256 = pkgs.lib.fakeSha256; # uncomment to get error with real sha256
              vendorSha256 = "sha256-gJ2gcscEIFaeRztTR1kCC8N8uiNmbo5pZDoAyQoFz7c=";
            };
        });

      meta = forAllSystems
        ({ pkgs }: with pkgs.lib; {
          description = "Generate Go structs from JSON";
          homepage = "https://github.com/knightpp/jgoson";
          license = licenses.mit;
          maintainers = with maintainers; [ knightpp ];
        });
    };
}
