#version 460
out vec4 frag_color;

uniform vec2 uSize;

void main() {
    vec2 color = gl_FragCoord.xy / uSize;
    frag_color = vec4(color.xy, pow(1 - color.x * color.y, 2), 1);
}
