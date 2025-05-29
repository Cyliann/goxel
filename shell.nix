{pkgs ? import <nixpkgs> {}}:
pkgs.mkShell {
  packages = with pkgs; [
    wayland
    libxkbcommon
    libGL
    pkg-config
    glfw-wayland
    xorg.libX11
  ];

  shellHook = "nu";
}
