{
  description = "PDM Flake";

  inputs = {
  	flake-utils.url = "github:numtide/flake-utils";
  	nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
	gomod2nix.url = "github:tweag/gomod2nix";
	};

  outputs = { self, nixpkgs, flake-utils, gomod2nix }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = import nixpkgs {
      	inherit system;
	overlays = [ gomod2nix.overlay];
		
	};

	pdm-cli = pkgs.buildGoApplication {
	  pname = "pdm-cli";
	  version = "0.1";
	  src = ./.;
	  modules = ./gomod2nix.toml;


  	  nativeBuildInputs = [ pkgs.installShellFiles ];

	  installPhase = ''
          	runHook preInstall
          	mkdir -p $out
          	dir="$GOPATH/bin"
          	[ -e "$dir" ] && cp -r $dir $out
	  	
    	        installShellCompletion --cmd pdm \
    	          --bash <($out/bin/pdm completion bash) \
    	          --zsh <($out/bin/pdm completion zsh)
          	runHook postInstall
	  '';
	};
      in rec {

	packages = {
		pdm-cli = pdm-cli;
	};

        defaultPackage = packages.pdm-cli;
        overlay = final: prev: {
          pdm-cli = pdm-cli;
        };

        devShell = pkgs.mkShell {
		buildInputs = [ pkgs.gomod2nix pkgs.go packages.pdm-cli ]; 
		};

      });
}
