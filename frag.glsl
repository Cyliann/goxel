#version 460
out vec4 frag_color;

uniform vec2 uSize;

bool circleIntersection(in vec2 uv) {
    return length(uv) < .5;
}
void main() {
    vec2 uv = (gl_FragCoord.xy * 2.0 - uSize) / uSize.y;
    frag_color = vec4(float(circleIntersection(uv)), 0, 0, 1);
}
