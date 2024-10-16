with import <nixpkgs> {};
  stdenv.mkDerivation {
    name = "goxel";
    buildInputs = [
      wayland
      libxkbcommon
      libGL
      pkg-config
      glfw-wayland
      xorg.libX11
    ];
  }
