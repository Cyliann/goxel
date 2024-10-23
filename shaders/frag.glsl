#version 460
out vec4 frag_color;

uniform vec2 uSize;
uniform float uTime;
uniform vec3 uPlayerPos;
const float PI = 3.14159265359;

vec3 raySphereIntersection(vec3 rayOrigin, vec3 rayDirection, float radius, vec3 lightDirection) {
    vec3 color = vec3(0.53, 0.81, 0.94);
    float a = dot(rayDirection, rayDirection);
    float b = 2.0 * dot(rayOrigin, rayDirection);
    float c = dot(rayOrigin, rayOrigin) - pow(radius, 2);
    float discriminant = pow(b, 2) - 4 * a * c;

    if (discriminant < 0) {
        return color;
    }

    float d = (-b - sqrt(discriminant)) / (2 * a);
    vec3 hitPoint = rayOrigin + rayDirection * d;
    vec3 normal = normalize(hitPoint);
    float lighting = max(dot(normal, -lightDirection), 0);

    color = vec3((sin(uTime * 2) + 1) / 2, (cos(uTime) + 1) / 2, (cos(uTime * 2) + 1) / 2) * lighting;
    return color;
}

void main() {
    vec2 uv = (gl_FragCoord.xy * 2.0 - uSize) / uSize.y;
    vec3 rayOrigin = uPlayerPos;
    vec3 rayDirection = normalize(vec3(uv, -1));
    vec3 lightDirection = -1 * normalize(vec3(sin(uTime * 2), 1, cos(uTime * 2)));

    frag_color = vec4(raySphereIntersection(rayOrigin, rayDirection, .7, lightDirection), 1);
    // frag_color = vec4(1);
}
