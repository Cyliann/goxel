#version 460
out vec4 frag_color;

uniform vec2 uSize;

int raySphereIntersection(vec3 rayOrigin, vec3 rayDirection, float radius) {
    int hit = 0;
    float a = dot(rayDirection, rayDirection);
    float b = 2.0 * dot(rayOrigin, rayDirection);
    float c = dot(rayOrigin, rayOrigin) - pow(radius, 2);
    float discriminant = pow(b, 2) - 4 * a * c;

    if (discriminant >= 0) {
        hit = 1;
    }

    return hit;
}

void main() {
    vec2 uv = (gl_FragCoord.xy * 2.0 - uSize) / uSize.y;
    vec3 rayOrigin = vec3(0, 0, 5);
    vec3 rayDirection = vec3(uv, 1);

    frag_color = vec4(raySphereIntersection(rayOrigin, rayDirection, 2), 0, 0, 1);
}
