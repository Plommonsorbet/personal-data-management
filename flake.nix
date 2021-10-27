{
  description = "PDM Flake";

  inputs = {
  	nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
	gomod2nix.url = "github:tweag/gomod2nix";
  	flake-utils.url = "github:numtide/flake-utils";
	};

  outputs = { self, nixpkgs, flake-utils, gomod2nix }: 
  	let
	system = "x86_64-linux";
      	pkgs = import nixpkgs {
		inherit system;		
		overlays = [ gomod2nix.overlay ];
	};

	in rec {
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
		packages.${system}.pdm-cli = pdm-cli;
		defaultPackage.${system} = pdm-cli;
		evShell = pkgs.mkShell { buildInputs = [pkgs.gomod2nix];};

      		overlay = final: prev: {
      			pdm-cli = pdm-cli;
      		};

	};
  	

}

#{
#  description = "PDM Flake";
#
#  inputs = {
#  	nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
#	gomod2nix.url = "github:tweag/gomod2nix";
#  	flake-utils.url = "github:numtide/flake-utils";
#	};
#
#  outputs = { self, nixpkgs, flake-utils, gomod2nix }:
#
#      let 
#	system = "x86_64-linux";
#      	pkgs = import nixpkgs {
#		system = "x86_64-linux";
#		overlays = [ gomod2nix.overlay ];
#	};
#
#	pdm-cli = pkgs.buildGoApplication {
#	  pname = "pdm-cli";
#	  version = "0.1";
#	  src = ./.;
#	  modules = ./gomod2nix.toml;
#
#  	  nativeBuildInputs = [ pkgs.installShellFiles ];
#
#	  installPhase = ''
#          	runHook preInstall
#          	mkdir -p $out
#          	dir="$GOPATH/bin"
#          	[ -e "$dir" ] && cp -r $dir $out
#	  	
#    	        installShellCompletion --cmd pdm \
#    	          --bash <($out/bin/pdm completion bash) \
#    	          --zsh <($out/bin/pdm completion zsh)
#          	runHook postInstall
#	  '';
#	};
#    in flake-utils.lib.simpleFlake {
#    	inherit self nixpkgs;
#	name = "simple";
#	packages.pdm-cli = pdm-cli;
#
#        defaultPackage = pdm-cli;
#
#        devShell = pkgs.mkShell {
#		buildInputs = [ pkgs.gomod2nix pkgs.go pdm-cli ]; 
#		};
#
#
#      };
#
#}
