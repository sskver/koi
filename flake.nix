{
  description = "koi - waveshare e-paper display library in go";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };

        vendorHash = "sha256-7wtOC5xhk0nNfFhnzK4h1NGRVgLZOGYl5OR/DuI9tJY=";

        # per-platform build settings for everything under cmd/<platform>/.
        # add a new board by adding its directory under cmd/ plus an entry here.
        platformArch = {
          rpi = { goos = "linux"; goarch = "arm64"; };
        };

        isDir = name: type: type == "directory";

        dirNames = path: builtins.attrNames
          (pkgs.lib.filterAttrs isDir (builtins.readDir path));

        platforms = dirNames ./cmd;

        mkDemoPackage = platform: name: arch: pkgs.buildGoModule {
          pname = "koi";
          version = "0.1.0";
          src = ./.;
          inherit vendorHash;

          subPackages = [ "cmd/${platform}/${name}" ];

          env = {
            CGO_ENABLED = "0";
            GOOS = arch.goos;
            GOARCH = arch.goarch;
          };

          ldflags = [ "-s" "-w" ];

          # buildGoModule runs tests by default
          # skip, since it would try to run binaries for other archs on the build host
          doCheck = false;

          installPhase = ''
            mkdir -p $out/bin
            cp $GOPATH/bin/${name} $out/bin/koi-${platform}-${name}-${arch.goarch}
          '';
        };

        demoPackages = pkgs.lib.foldl'
          (acc: platform:
            let
              arch = platformArch.${platform} or (throw
                "flake.nix: no GOOS/GOARCH configured for platform '${platform}' - add it to platformArch");

              names = dirNames (./cmd + "/${platform}");

              platformPackages = pkgs.lib.listToAttrs (map
                (name: {
                  name = "koi-${platform}-${name}-${arch.goarch}";
                  value = mkDemoPackage platform name arch;
                })
                names);
            in
            acc // platformPackages)
          { }
          platforms;
      in
      {
        packages = demoPackages // {
          default = pkgs.buildGoModule {
            pname = "koi";
            version = "0.1.0";
            src = ./.;
            inherit vendorHash;

            subPackages = [ "cmd/" ];
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
            go-tools # staticcheck etc.
            gnumake
          ];

          shellHook = ''
            echo "koi dev shell: go $(go version | cut -d' ' -f3), make $(make --version | head -n1)"
          '';
        };
      });
}
